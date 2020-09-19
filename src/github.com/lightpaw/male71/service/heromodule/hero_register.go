package heromodule

import (
	"github.com/lightpaw/male7/entity"
	"github.com/lightpaw/male7/gen/iface"
	"github.com/lightpaw/male7/pb/shared_proto"
	"github.com/lightpaw/male7/service/conflict/heroservice/herolock"
)

// 玩家注册事件

// 获得物品事件监听
type AddGoodsListener interface {
	OnAddGoodsEvent(hero *entity.Hero, result herolock.LockResult, goodsId uint64, addCount uint64)
}

var addGoodsListeners []AddGoodsListener = make([]AddGoodsListener, 0, 3)

func RegisterAddGoodsListener(event AddGoodsListener) {
	addGoodsListeners = append(addGoodsListeners, event)
}

func onAddGoodsEvent(hero *entity.Hero, result herolock.LockResult, goodsId uint64, addCount uint64) {
	for _, listener := range addGoodsListeners {
		listener.OnAddGoodsEvent(hero, result, goodsId, addCount)
	}
}

// 玩家上线事件监听
type HeroOnlineListener interface {
	OnHeroOnline(hc iface.HeroController)
}

var heroOnlineListeners []HeroOnlineListener = make([]HeroOnlineListener, 0, 3)

func RegisterHeroOnlineListener(lstr HeroOnlineListener) {
	heroOnlineListeners = append(heroOnlineListeners, lstr)
}

func OnHeroOnlineEvent(hc iface.HeroController) {
	for _, lstr := range heroOnlineListeners {
		lstr.OnHeroOnline(hc)
	}
}

// 玩家下线事件监听
type HeroOfflineListener interface {
	OnHeroOffline(hc iface.HeroController)
}

var heroOfflineListeners []HeroOfflineListener = make([]HeroOfflineListener, 0, 3)

func RegisterHeroOfflineListener(lstr HeroOfflineListener) {
	heroOfflineListeners = append(heroOfflineListeners, lstr)
}

func OnHeroOfflineEvent(hc iface.HeroController) {
	for _, lstr := range heroOfflineListeners {
		lstr.OnHeroOffline(hc)
	}
}

// 玩家下线事件监听
type HeroPveTroopChangeEventHandler func(id int64, troopType shared_proto.PveTroopType)

var heroPveTroopChangeEventHandlers = make([]HeroPveTroopChangeEventHandler, 0, 2)

func RegisterHeroPveTroopChangeEventHandlers(handler HeroPveTroopChangeEventHandler) {
	heroPveTroopChangeEventHandlers = append(heroPveTroopChangeEventHandlers, handler)
}

func OnHeroPveTroopChange(id int64, troopType shared_proto.PveTroopType) {
	for _, handle := range heroPveTroopChangeEventHandlers {
		handle(id, troopType)
	}
}
