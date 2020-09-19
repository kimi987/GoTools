package iface

import (
	"github.com/lightpaw/pbutil"
	"net"
	"github.com/lightpaw/muxface"
	"github.com/lightpaw/male7/service/ticker/tickdata"
)

type CountryHeroWalker func(id int64, hc HeroController)

type HeroWalker func(id int64, hc HeroController)

type UserWalker func(id int64, cu ConnectedUser)

type MsgFunc func() pbutil.Buffer

type ServeListener func(listener net.Listener) (muxface.Listener, error)

type ConnHandler func(conn muxface.Conn)

type Func func()

type TickFunc func(tick tickdata.TickTime)