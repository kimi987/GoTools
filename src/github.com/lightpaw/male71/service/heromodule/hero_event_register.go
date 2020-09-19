package heromodule

import (
	"github.com/lightpaw/male7/entity"
	"github.com/lightpaw/male7/service/conflict/heroservice/herolock"
	"github.com/lightpaw/male7/pb/shared_proto"
	"github.com/lightpaw/male7/util/call"
)

// 玩家下线事件监听
type HeroEventFunc func(hero *entity.Hero, result herolock.LockResult, event shared_proto.HeroEvent)

type HeroEventHandler struct {
	name string
	f    HeroEventFunc
}

var heroEventHandlers []*HeroEventHandler

func RegisterHeroEventHandler(name string, f HeroEventFunc) {
	heroEventHandlers = append(heroEventHandlers, &HeroEventHandler{
		name: name,
		f:    f,
	})
}

func OnHeroEvent(hero *entity.Hero, result herolock.LockResult, event shared_proto.HeroEvent, skipEvents ...shared_proto.HeroEvent) {

	if len(skipEvents) > 0 {
		for _, e := range skipEvents {
			if event == e {
				return
			}
		}
	}

	for _, handler := range heroEventHandlers {
		call.CatchPanic(func() {
			handler.f(hero, result, event)
		}, handler.name)
	}
}

// 玩家
type HeroEventWithSubTypeFunc func(hero *entity.Hero, result herolock.LockResult, event shared_proto.HeroEvent, subType uint64)

type HeroEventWithSubTypeHandler struct {
	name string
	f    HeroEventWithSubTypeFunc
}

var heroEventWithSubTypeHandlers []*HeroEventWithSubTypeHandler

func RegisterHeroEventWithSubTypeHandler(name string, f HeroEventWithSubTypeFunc) {
	heroEventWithSubTypeHandlers = append(heroEventWithSubTypeHandlers, &HeroEventWithSubTypeHandler{
		name: name,
		f:    f,
	})
}

func OnHeroEventWithSubType(hero *entity.Hero, result herolock.LockResult, event shared_proto.HeroEvent, subType uint64, skipEvents ...shared_proto.HeroEvent) {

	if len(skipEvents) > 0 {
		for _, e := range skipEvents {
			if event == e {
				return
			}
		}
	}

	for _, handler := range heroEventWithSubTypeHandlers {
		call.CatchPanic(func() {
			handler.f(hero, result, event, subType)
		}, handler.name)

	}
}
