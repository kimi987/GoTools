package discover

import (
	"github.com/lightpaw/go-zookeeper/zk"
	"github.com/lightpaw/logrus"
	"github.com/pkg/errors"
	"strconv"
	"strings"
	"time"
)

var (
	slash = "/"[0]
)

func (c *Conn) CreateEphemeral(path string, data []byte) error {
	if len(path) == 0 || path[0] != slash {
		return errors.New("Path must starts with /")
	}
	c.createParentIfNecessary(path) // ignore err. try create node anyway

	var err error
outer:
	for i := 0; i < 3 && c.ServerStatus() != Expired; i++ {
		_, err = c.conn.Create(path, data, zk.FlagEphemeral, zk.WorldACL(zk.PermAll))
		switch err {
		case zk.ErrNoNode:
			c.createParentIfNecessary(path)
		case zk.ErrNodeExists:
			return errors.Wrap(err, "zk中节点已存在")
		case nil:
			break outer
		default:
			// other error. try later
			logrus.WithError(err).WithField("path", path).Error("Create节点失败")
			time.Sleep(1 * time.Second)
		}
	}

	return err
}

// 在zk中创建个节点表示自己. path必须以斜杠/开头
func (c *Conn) CreateEphemeralSequential(path string, data []byte) (*RegisterdNode, error) {
	if len(path) == 0 || path[0] != slash {
		return nil, errors.New("Path must starts with /")
	}
	c.createParentIfNecessary(path) // ignore err. try create node anyway

	var result string
	var err error
	for i := 0; i < 3 && c.ServerStatus() != Expired; i++ {
		result, err = c.conn.Create(path, data, zk.FlagSequence|zk.FlagEphemeral, zk.WorldACL(zk.PermAll))
		if err != nil && err == zk.ErrNoNode {
			c.createParentIfNecessary(path)
		} else if err == nil {
			break
		} else {
			// other error. try later
			time.Sleep(1 * time.Second)
		}
	}

	if err != nil {
		return nil, err
	}

	if !strings.HasPrefix(result, path) {
		return nil, errors.New("sequential node must have the same prefix: " + result)
	}

	generatedPath := result[len(path):]
	id, err := strconv.ParseInt(generatedPath, 10, 0)
	if err != nil {
		return nil, errors.Wrap(err, "sequential node cannot be cast to int64: "+generatedPath)
	}

	return &RegisterdNode{
		fullPath:      result,
		generatedPath: generatedPath,
		id:            id,
		data:          data,
		conn:          c,
	}, nil
}

// 代表zk中一个节点
type RegisterdNode struct {
	fullPath      string // full path of the node
	generatedPath string // generated id
	id            int64  // generated id casted to int64
	data          []byte // last data saved into zk
	conn          *Conn
}

// 获得此节点在zk自增后的id
func (r RegisterdNode) Id() int64 {
	return r.id
}

// 获得此节点最近一次保存在zk中的数据
func (r RegisterdNode) Data() []byte {
	return r.data
}

// 更新节点在zk中的数据
func (r *RegisterdNode) UpdateData(newData []byte) (err error) {
	for i := 0; i < 10 && r.conn.ServerStatus() != Expired; i++ {
		_, err = r.conn.conn.Set(r.fullPath, newData, -1)
		if err == nil {
			r.data = newData
			return
		}

		if err == zk.ErrClosing {
			return
		}

		time.Sleep(1 * time.Second)
	}

	if err == nil {
		r.data = newData
	}
	return
}

func (r RegisterdNode) Remove() (err error) {
	for i := 0; i < 10 && r.conn.ServerStatus() != Expired; i++ {
		err = r.conn.conn.Delete(r.fullPath, -1)
		if err == nil {
			return
		}

		if err == zk.ErrClosing {
			return
		}

		if err == zk.ErrNoNode {
			return nil
		}

		time.Sleep(1 * time.Second)
	}

	return
}

func (c *Conn) createParentIfNecessary(path string) {
	needCheck := true
	for i := 1; i < len(path); i++ {
		if path[i] == slash {
			parentPath := path[0:i]
			exists := false
			var err error
			if needCheck {
				for i := 0; i < 10 && c.ServerStatus() != Expired; i++ {
					exists, _, err = c.conn.Exists(parentPath)
					if err != nil {
						logrus.WithError(err).WithField("path", parentPath).Error("Zk check parent path exists fail")
						time.Sleep(1 * time.Second)
					} else {
						break
					}
				}

				if !exists {
					needCheck = false
				}
			}
			if !exists {
				logrus.WithField("path", parentPath).Debug("Zk creating node parent")

				for i := 0; i < 10 && c.ServerStatus() != Expired; i++ {

					_, err = c.conn.Create(parentPath, []byte{}, 0, zk.WorldACL(zk.PermAll))
					if err != nil && err != zk.ErrNodeExists {
						logrus.WithField("path", parentPath).WithError(err).Error("Zk create node parent fail")
						time.Sleep(1 * time.Second)
					} else {
						logrus.WithField("path", parentPath).Debug("Zk create node parent success")
						break
					}
				}
			}
		}
	}
	return
}
