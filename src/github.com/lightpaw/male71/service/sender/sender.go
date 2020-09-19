package sender

import (
	"github.com/lightpaw/male7/util/msg"
	"github.com/lightpaw/pbutil"
)

type Sender interface {
	Id() int64

	// 发送消息.
	SendAll(msgs []pbutil.Buffer)

	Send(msg pbutil.Buffer)

	SendIfFree(msg pbutil.Buffer)
}

type ClosableSender interface {
	Sender
	//Close()
	//
	//CloseAndWait()
	IsClosed() bool

	Disconnect(err msg.ErrMsg)
	DisconnectAndWait(err msg.ErrMsg)
}
