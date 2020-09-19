package muxface

import (
	"net"
	"github.com/lightpaw/logintoken"
	"io"
)

type Listener interface {
	Close() error
	Addr() net.Addr
	Accept() (Conn, error)
}

type Conn interface {
	io.Writer
	Close()
	WriteIfFree([]byte) error
	MsgChan() <-chan interface{}
	ClosedNotify() <-chan struct{}
	GetLoginToken() logintoken.LoginToken
}
