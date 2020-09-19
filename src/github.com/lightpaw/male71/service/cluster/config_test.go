package cluster

import (
	"testing"
	"github.com/lightpaw/go-zookeeper/zk"
	. "github.com/onsi/gomega"
	"fmt"
	"github.com/lightpaw/rpc7"
	"context"
	"time"
	"github.com/lightpaw/logrus"
)

func TestName(t *testing.T) {

	//addr:="login.lightpaw.com:7890"
	addr := "127.0.0.1:9527"
	loginClient, err := rpc7.NewClient(addr)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(loginClient.CheckAlive(context.TODO()))
	fmt.Println(loginClient.Handle(context.TODO(), "test", "", 0, &GameServerInfoProto{}))
}

//func TestSave(t *testing.T) {
//	RegisterTestingT(t)
//
//	zkConn, _, err := zk.Connect([]string{"127.0.0.1:2181"}, sessionTimeout, zk.WithSessionExpireAndQuit())
//	Ω(err).Should(Succeed())
//
//	cosPath := "/m7/config/cos"
//	//b, _, err := zkConn.Get(cosPath)
//	//Ω(err).Should(Equal(zk.ErrNoNode))
//	//Ω(b).Should(BeEmpty())
//
//	appid := ""
//	secretid := ""
//	secretKey := ""
//	region := ""
//	bucketName := ""
//	dataMap := map[string]interface{}{
//		"appid":      appid,
//		"secretid":   secretid,
//		"secretKey":  secretKey,
//		"region":     region,
//		"bucketName": bucketName,
//	}
//	data, err := jsoniter.Marshal(dataMap)
//	Ω(err).Should(Succeed())
//	fmt.Println(string(data))
//
//	Ω(err).Should(Succeed())
//	Ω(encryptAndSave(zkConn, data, cosPath)).Should(Succeed())
//	Ω(loadAndUnencrypt(zkConn, cosPath)).Should(Equal(data))
//
//	// delete
//	//zkConn.Delete(cosPath, 0)
//}

func TestGetW(t *testing.T) {

	RegisterTestingT(t)

	zkConn, _, err := zk.Connect([]string{"127.0.0.1:2181"}, sessionTimeout, zk.WithSessionExpireAndQuit())
	Ω(err).Should(Succeed())

	path := "/a/b/c"

	closeChan := make(chan struct{})

	go Watch(zkConn, path, closeChan, func(data []byte) {
		logrus.Info("data", len(data), string(data))
	})

	data1 := []byte("haha")
	createNode(zkConn, data1, path)

	time.Sleep(1 * time.Second)
	data2 := []byte("b2")
	createNode(zkConn, data2, path)

	time.Sleep(1 * time.Second)
	data3 := []byte("b3")
	createNode(zkConn, data3, path)

	zkConn.Delete(path, -1)
	time.Sleep(1 * time.Second)

	createNode(zkConn, data1, path)
	time.Sleep(1 * time.Second)

	zkConn.Delete(path, -1)
	time.Sleep(1 * time.Second)

	createNode(zkConn, data1, path)
	time.Sleep(1 * time.Second)

	zkConn.Delete(path, -1)
	time.Sleep(1 * time.Second)

	close(closeChan)
	<-closeChan
}
