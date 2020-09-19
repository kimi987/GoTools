package herolock

import (
	"github.com/lightpaw/male7/entity"
	"github.com/lightpaw/male7/service/sender"
	"github.com/lightpaw/pbutil"
)

type HeroLocker interface {
	Func(Func)
	FuncNotError(FuncNotError) (hasError bool)

	FuncWithSend(SendFunc, sender.ClosableSender) (hasError bool)
}

type Func func(hero *entity.Hero, err error) (heroChanged bool)
type FuncNotError func(hero *entity.Hero) (heroChanged bool)
type SendFunc func(hero *entity.Hero, result LockResult)
type SendFuncWithError func(hero *entity.Hero, result LockResult, err error)

type LockResult interface {
	Add(pbutil.Buffer)
	AddFunc(func() pbutil.Buffer)

	AddBroadcast(pbutil.Buffer)

	Changed()
	Ok()
}

