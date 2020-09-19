package nonmux

import (
	"encoding/binary"
	"github.com/lightpaw/logintoken"
	"github.com/lightpaw/logrus"
	"io"
	"net"
	"sync"
	"time"
)

const (
	rewriterSize       = 16384
	msgChanSize        = 64
	clientRewriterSize = 32768
)

type Conn struct {
	id          uint64
	listener    *Listener
	conn        net.Conn
	token       logintoken.LoginToken
	unmarshaler msgUnmarshaller

	writeBytes uint64
	readBytes  uint64
	rewriter   *rewriter

	writeLock sync.Mutex
	closed    bool

	disconnectFunc *time.Timer
	closedCounter  uint64

	closedNotify chan struct{}
	msgChan      chan interface{}
}

func newConn(id uint64, listener *Listener, c net.Conn, token logintoken.LoginToken) *Conn {
	result := &Conn{
		id:           id,
		listener:     listener,
		conn:         c,
		token:        token,
		rewriter:     &rewriter{data: make([]byte, rewriterSize)},
		unmarshaler:  listener.config.Unmarshaller,
		closedNotify: make(chan struct{}),
		msgChan:      make(chan interface{}, msgChanSize),
	}

	token.GameServerRandomID = listener.randomID
	token.GenerateTime = time.Now()

	go result.recvLoop(c)

	return result
}

func (c *Conn) ID() uint64 {
	return c.id
}

func (c *Conn) MsgChan() <-chan interface{} {
	return c.msgChan
}

func (c *Conn) ClosedNotify() <-chan struct{} {
	return c.closedNotify
}

func (c *Conn) GetLoginToken() logintoken.LoginToken {
	return c.token
}

func (c *Conn) WriteIfFree(data []byte) error {
	_, err := c.Write(data)
	return err
}

// 服务器主动关闭连接, 不等待重连
func (c *Conn) Close() {
	c.writeLock.Lock()
	defer c.writeLock.Unlock()

	if c.closed {
		return
	}

	c.doClose()
}

func (c *Conn) recvLoop(currentConn net.Conn) {
	headerBuf := make([]byte, 2)
	for {
		if _, err := io.ReadFull(currentConn, headerBuf); err != nil {
			c.markRst(currentConn)
			return
		}

		size := int(binary.LittleEndian.Uint16(headerBuf))

		buf := make([]byte, size)
		if _, err := io.ReadFull(currentConn, buf); err != nil {
			c.markRst(currentConn)
			return
		}

		msg, err := c.unmarshaler.Unmarshal(buf)
		if err != nil {
			logrus.WithError(err).Error("消息unmarshal失败")
			c.Close()
			return
		}

		c.writeLock.Lock()
		if c.conn != currentConn {
			c.writeLock.Unlock()
			return
		}

		select {
		case c.msgChan <- msg:
			c.readBytes += uint64(msgOriginalSize(size + 1)) // 每条消息的index被连接服过滤了
			c.writeLock.Unlock()
		default:
			c.writeLock.Unlock()
			logrus.Error("连接的msgChan满了")
			c.Close()
			return
		}
	}
}

func msgOriginalSize(size int) int {
	switch {
	case size <= 127:
		return size + 1
	default:
		return size + 2
	}
}

func (c *Conn) Write(data []byte) (int, error) {
	c.writeLock.Lock()

	if c.closed {
		c.writeLock.Unlock()
		return 0, io.ErrClosedPipe
	}

	c.rewriter.Push(data)
	c.writeBytes += uint64(len(data))

	if c.conn == nil {
		c.writeLock.Unlock()
		return len(data), nil
	}

	if n, err := c.conn.Write(data); err != nil {
		currentConn := c.conn
		c.writeLock.Unlock()
		c.markRst(currentConn) // unlock before call
		return n, err
	} else {
		c.writeLock.Unlock()
		return n, err
	}
}

func (c *Conn) markRst(currentConn net.Conn) {
	c.writeLock.Lock()
	defer c.writeLock.Unlock()

	if c.conn != currentConn {
		return
	}

	c.conn = nil

	c.closedCounter++
	counter := c.closedCounter

	c.disconnectFunc = time.AfterFunc(c.listener.config.CloseWaitTime, func() {
		c.writeLock.Lock()
		defer c.writeLock.Unlock()

		if c.conn == nil && c.closedCounter == counter && !c.closed {
			c.disconnectFunc = nil
			c.doClose()
		}
	})
}

// called under lock
func (c *Conn) doClose() {
	c.closed = true
	close(c.closedNotify)
	c.listener.removeConn(c)

	if c.conn != nil {
		c.conn.Close()
		c.conn = nil
	}

	if c.disconnectFunc != nil {
		c.disconnectFunc.Stop()
		c.disconnectFunc = nil
	}
}

func (c *Conn) tryReconn(conn net.Conn, readBytes, writeBytes uint64, token logintoken.LoginToken) (ok bool) {
	c.writeLock.Lock()
	defer c.writeLock.Unlock()

	if c.closed {
		return
	}

	if token.GameServerRandomID != c.token.GameServerRandomID || token.UserSelfID != c.token.UserSelfID || token.GameServerID != c.token.GameServerID {
		return
	}

	if writeBytes < c.readBytes {
		return
	}

	if c.writeBytes < readBytes {
		return
	}

	if int(c.writeBytes-readBytes) > c.rewriter.length {
		return
	}

	if writeBytes-c.readBytes > clientRewriterSize {
		return
	}

	canRewrite, b1, b2 := c.rewriter.Rewrite(c.writeBytes, readBytes)
	if !canRewrite {
		return
	}

	c.token.GenerateTime = time.Now()
	buf := make([]byte, logintoken.UNENCRYPTED_SIZE+16)
	c.token.Marshal(buf[:logintoken.UNENCRYPTED_SIZE])
	binary.LittleEndian.PutUint64(buf[logintoken.UNENCRYPTED_SIZE:], c.readBytes)
	binary.LittleEndian.PutUint64(buf[logintoken.UNENCRYPTED_SIZE+8:], c.writeBytes)

	if _, err := conn.Write(buf); err != nil {
		return
	}

	if b1 != nil {
		if _, err := conn.Write(b1); err != nil {
			return
		}

		if b2 != nil {
			if _, err := conn.Write(b2); err != nil {
				return
			}
		}
	}

	if c.conn != nil {
		c.conn.Close()
	}

	if c.disconnectFunc != nil {
		c.disconnectFunc.Stop()
		c.disconnectFunc = nil
	}
	c.conn = conn

	go c.recvLoop(conn)
	return true
}
