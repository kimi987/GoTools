package cluster

import (
	"github.com/lightpaw/logrus"
	"github.com/pkg/errors"
	"time"
	"github.com/lightpaw/go-zookeeper/zk"
	"crypto/cipher"
	"crypto/aes"
	"encoding/binary"
	"io"
	"crypto/rand"
	mrand "math/rand"
)

var (
	errMalformData  = errors.New("malform data loaded")
	configSecretKey = "e27d9563cd8dbccd656126d8ea2ba738"
	unencryptKey, _ = aes.NewCipher([]byte(configSecretKey))
	slash           = "/"[0]
)

func loadAndUnencrypt(zkConn *zk.Conn, path string) ([]byte, error) {
	encrypted, _, err := zkConn.Get(path)
	if err != nil {
		return nil, errors.Wrap(err, "从zk读取失败")
	}

	if len(encrypted) <= 16 {
		return nil, errMalformData
	}

	iv := encrypted[:16]
	data := encrypted[16:]

	stream := cipher.NewCFBDecrypter(unencryptKey, iv)
	stream.XORKeyStream(data, data)
	return data, nil
}

func encryptAndSave(zkConn *zk.Conn, data []byte, path string) error {
	encrypted := make([]byte, 16+len(data))

	iv := encrypted[:16]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		binary.LittleEndian.PutUint64(iv, mrand.Uint64())
		binary.LittleEndian.PutUint64(iv[8:], mrand.Uint64())
	}

	stream := cipher.NewCFBEncrypter(unencryptKey, iv)
	stream.XORKeyStream(encrypted[16:], data)

	return createNode(zkConn, encrypted, path)
}

func createNode(zkConn *zk.Conn, data []byte, path string) error {
	parentCreated := false
CREATE:
	_, err := zkConn.Create(path, data, 0, zk.WorldACL(zk.PermAll))
	switch err {
	case nil:
		return nil

	case zk.ErrNoNode:
		if parentCreated {
			return err
		}
		parentCreated = true
		createParentIfNecessary(zkConn, path)
		goto CREATE

	case zk.ErrNodeExists:
		_, err = zkConn.Set(path, data, -1)
		if err != nil {
			return errors.Wrap(err, "无法更新现有节点数据")
		}

		return nil

	default:
		return errors.Wrap(err, "无法创建节点")
	}
}

func createParentIfNecessary(conn *zk.Conn, path string) {
	needCheck := true
	for i := 1; i < len(path); i++ {
		if path[i] == slash {
			parentPath := path[0:i]
			exists := false
			var err error
			if needCheck {
				for i := 0; i < 10; i++ {
					exists, _, err = conn.Exists(parentPath)
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

				for i := 0; i < 10; i++ {

					_, err = conn.Create(parentPath, []byte{}, 0, zk.WorldACL(zk.PermAll))
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
