package nonmux

import (
	"crypto/rc4"
	"encoding/binary"
	"github.com/lightpaw/logintoken"
	"github.com/lightpaw/logrus"
	"github.com/pkg/errors"
	"io"
	"math/rand"
	"net"
	"sync"
	"sync/atomic"
	"time"
)

var (
	handshakeKey = []byte("malimalihong")
)

type Config struct {
	Unmarshaller  msgUnmarshaller
	CloseWaitTime time.Duration
	DontEncrypt   bool
}

type Listener struct {
	listener net.Listener
	config   *Config

	randomID      uint32
	connIdCounter uint64

	connMap     map[uint64]*Conn
	connMapLock sync.Mutex
}

type msgUnmarshaller interface {
	Unmarshal([]byte) (interface{}, error)
}

func ListenAddr(addr string, config *Config) (*Listener, error) {
	lsn, err := net.Listen("tcp", addr)
	if err != nil {
		return nil, err
	}
	return Listen(lsn, config)
}

func Listen(lsn net.Listener, config *Config) (*Listener, error) {
	if config.CloseWaitTime < 0 {
		logrus.WithField("CloseWaitTime", config.CloseWaitTime).Panic("CloseWaitTime must >= 0")
	}
	if config.Unmarshaller == nil {
		logrus.Panic("Unmarshaller cannot be nil")
	}

	logrus.WithField("addr", lsn.Addr().String()).Info("nonmux开始监听端口")

	result := &Listener{
		listener: lsn,
		config:   config,
		connMap:  make(map[uint64]*Conn),
		randomID: rand.Uint32(),
	}
	return result, nil
}

func (l *Listener) Close() error {
	return l.listener.Close()
}

func (l *Listener) Addr() net.Addr {
	return l.listener.Addr()
}

func (l *Listener) Accept() (*Conn, error) {
	for {
		c, err := l.listener.Accept()
		if err != nil {
			return nil, errors.Wrap(err, "从listener接收连接出错")
		}
		c.(*net.TCPConn).SetLinger(5)

		logrus.Debug("nonmux接收到新连接")

		if conn, err := l.handshake(c); err != nil {
			logrus.WithError(err).Error("Conn handshake出错")
			c.Close()
			continue
		} else {
			logrus.WithField("conn", conn).Debug("nonmux新连接握手成功")
			if conn == nil {
				// reconn success
				continue
			}

			return conn, nil
		}
	}
}

func (l *Listener) getConn(id uint64) *Conn {
	l.connMapLock.Lock()
	defer l.connMapLock.Unlock()

	return l.connMap[id]
}

func (l *Listener) putConn(conn *Conn) {
	l.connMapLock.Lock()
	defer l.connMapLock.Unlock()

	l.connMap[conn.id] = conn
}

func (l *Listener) removeConn(conn *Conn) {
	l.connMapLock.Lock()
	defer l.connMapLock.Unlock()

	delete(l.connMap, conn.id)
}

// 握手, 先进行连接握手, 再读取客户端信息
// 如果有错, err != nil
// 如果是新连接, 则conn != nil
// 如果是重连, 则conn == nil
func (l *Listener) handshake(c net.Conn) (*Conn, error) {
	c.SetDeadline(time.Now().Add(5 * time.Second))
	defer c.SetDeadline(time.Time{})

	buf := make([]byte, 16)

	if _, err := io.ReadFull(c, buf); err != nil {
		return nil, err
	}

	rc4Key, err := rc4.NewCipher(handshakeKey)
	if err != nil {
		return nil, err
	}
	rc4Key.XORKeyStream(buf, buf)

	first8 := binary.LittleEndian.Uint64(buf)

	myFirst8 := rand.Uint64() | 1
	if l.config.DontEncrypt {
		// 第二位是0表示加密, 1 表示不加密
		myFirst8 = myFirst8 | (1 << 1)
	} else {
		myFirst8 = myFirst8 & (^(uint64(1) << 1))
	}
	binary.LittleEndian.PutUint64(buf, myFirst8)

	rc4Key.XORKeyStream(buf, buf)
	if _, err := c.Write(buf); err != nil {
		return nil, err
	}

	if _, err := io.ReadFull(c, buf); err != nil {
		return nil, err
	}
	rc4Key.XORKeyStream(buf, buf)

	if first8 != binary.LittleEndian.Uint64(buf) {
		return nil, errors.New("first 8 not same")
	}

	if myFirst8 != binary.LittleEndian.Uint64(buf[8:]) {
		return nil, errors.New("my first 8 not match")
	}

	// handshake success

	buf = make([]byte, logintoken.UNENCRYPTED_SIZE+24)
	if _, err := io.ReadFull(c, buf); err != nil {
		return nil, errors.Wrap(err, "读取logintoken出错")
	}

	token := logintoken.Unmarshal(buf[:logintoken.UNENCRYPTED_SIZE])
	reconnID := binary.LittleEndian.Uint64(buf[logintoken.UNENCRYPTED_SIZE:])
	readBytes := binary.LittleEndian.Uint64(buf[logintoken.UNENCRYPTED_SIZE+8:])
	writeBytes := binary.LittleEndian.Uint64(buf[logintoken.UNENCRYPTED_SIZE+16:])

	if reconnID == 0 {
		token.GenerateTime = time.Now()
		token.GameServerRandomID = l.randomID

		conn := newConn(atomic.AddUint64(&l.connIdCounter, 1), l, c, token)

		buf := make([]byte, logintoken.UNENCRYPTED_SIZE+8)
		token.Marshal(buf[:logintoken.UNENCRYPTED_SIZE])
		binary.LittleEndian.PutUint64(buf[logintoken.UNENCRYPTED_SIZE:], conn.id)

		if _, err := c.Write(buf); err != nil {
			return nil, errors.Wrap(err, "写入logintoken出错")
		}
		l.putConn(conn)
		return conn, nil
	}

	oldConn := l.getConn(reconnID)
	if oldConn == nil {
		logrus.Debug("reconn id不存在")
		for i := 0; i < 64; i++ {
			buf[i] = 0
		}
		c.Write(buf[:64])
		return nil, errors.New("reconn Conn id not exist")
	}

	if !oldConn.tryReconn(c, readBytes, writeBytes, token) {
		logrus.Debug("reconn 失败")
		for i := 0; i < 64; i++ {
			buf[i] = 0
		}
		c.Write(buf[:64])
		return nil, errors.New("reconn info not match")
	}

	logrus.Debug("reconn 成功")
	return nil, nil
}
