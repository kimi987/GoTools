package nonmux

import (
	"net"
	"github.com/lightpaw/muxface"
)

func Wrap(l *Listener, err error) (muxface.Listener, error) {
	return &face{
		l: l,
	}, err
}

type face struct {
	l *Listener
}

func (l *face) Close() error                  { return l.l.Close() }
func (l *face) Addr() net.Addr                { return l.l.listener.Addr() }
func (l *face) Accept() (muxface.Conn, error) { return l.l.Accept() }
