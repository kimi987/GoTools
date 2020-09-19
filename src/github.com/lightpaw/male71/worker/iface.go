package worker

import (
	"github.com/lightpaw/muxface"
)

//type Conn interface {
//	io.Writer
//	Close()
//	WriteIfFree([]byte) error
//	MsgChan() <-chan interface{}
//	ClosedNotify() <-chan struct{}
//	GetLoginToken() logintoken.LoginToken
//}

type Conn = muxface.Conn
