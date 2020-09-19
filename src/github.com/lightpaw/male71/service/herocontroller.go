package service

import (
	"github.com/lightpaw/male7/entity"
	"github.com/lightpaw/male7/face"
	"github.com/lightpaw/male7/module/realm/realmface"
	"github.com/lightpaw/male7/pb/server_proto"
	"github.com/lightpaw/male7/service/conflict/heroservice/herolock"
	"github.com/lightpaw/male7/service/sender"
	"github.com/lightpaw/male7/util/atomic"
	"github.com/lightpaw/male7/util/idbytes"
	"github.com/lightpaw/male7/util/msg"
	"github.com/lightpaw/pbutil"
	"sync"
	"time"
	"github.com/lightpaw/male7/util/timeutil"
	"github.com/lightpaw/male7/entity/heroid"
	"github.com/lightpaw/male7/constants"
)

//gogen:iface entity
type HeroController struct {
	id      int64
	idBytes []byte
	//hero *entity.Hero

	sender     sender.ClosableSender
	clientIp   string
	clientIp32 uint32 /* 127.0.0.1 = 0x100007f */
	pf         uint32

	heroLocker herolock.HeroLocker

	// 地区关注
	careCondition atomic.Value

	// window area
	viewAreaValue atomic.Value

	// kimi 反向索引所看地块
	viewBlockIndex interface{}
	// kimi 观察对象列表
	watchObjList map[interface{}]int

	careWaterTimesMap map[int64]uint64

	nextSearchNoGuildHeros time.Time // 下次模糊搜索无联盟玩家的时间

	nextRefreshRecommendHeroTime time.Time

	nextSearchHeros time.Time // 下次模糊搜索玩家的时间

	lastClickTime time.Time // 上次鼠标点击的时间（各种模块定义各种Time太麻烦了，又不是外挂，玩家怎么可能同时操作N个业务，就判定这一个变量与dt比较足够了）

	isInBackgroud    atomic.Bool
	backgroudEndTime time.Time // 后台结束时间

	// 下次记录登陆日志时间
	nextWriteOnlineLogTime time.Time

	funcMux sync.RWMutex
	funcs   []face.BFunc
}

func (hc *HeroController) TryNextWriteOnlineLogTime(ctime time.Time, duration time.Duration) bool {
	if hc.nextWriteOnlineLogTime.Before(ctime) {
		isZero := timeutil.IsZero(hc.nextWriteOnlineLogTime)
		hc.nextWriteOnlineLogTime = ctime.Add(duration)
		if !isZero {
			return true
		}
	}
	return false
}

func (hc *HeroController) AddTickFunc(f face.BFunc) {
	hc.funcMux.Lock()
	defer hc.funcMux.Unlock()

	hc.funcs = append(hc.funcs, f)
}

func (hc *HeroController) TickFunc() {

	hc.funcMux.Lock()
	fs := hc.funcs
	hc.funcs = nil
	hc.funcMux.Unlock()

	if len(fs) <= 0 {
		return
	}

	var newFuncs []face.BFunc
	for _, f := range fs {
		if f != nil {
			if !f() {
				newFuncs = append(newFuncs, f)
			}
		}
	}

	if len(newFuncs) <= 0 {
		return
	}

	// 如果还有剩，加回原来的func列表后面
	hc.funcMux.Lock()
	defer hc.funcMux.Unlock()

	if len(hc.funcs) > 0 {
		hc.funcs = append(hc.funcs, newFuncs...)
	} else {
		hc.funcs = newFuncs
	}

}

func (hc *HeroController) TotalOnlineTime() time.Duration {
	return 0
}

func (hc *HeroController) GetClientIp() string {
	return hc.clientIp
}

func (hc *HeroController) GetClientIp32() uint32 {
	return hc.clientIp32
}

func (hc *HeroController) GetPf() uint32 {
	return hc.pf
}

func (hc *HeroController) GetIsInBackgroud() bool {
	return hc.isInBackgroud.Load()
}

func (hc *HeroController) UpdateIsInBackgroud(ctime time.Time) {
	if hc.isInBackgroud.Load() {
		if hc.backgroudEndTime.Before(ctime) {
			hc.isInBackgroud.Store(false)
		}
	}
}

func (hc *HeroController) SetIsInBackgroud(toSetTime time.Time, isInBackgroud bool) {
	hc.backgroudEndTime = toSetTime
	hc.isInBackgroud.Store(isInBackgroud)
}

func (hc *HeroController) GetCareWaterTimesMap() map[int64]uint64 {
	if hc.careWaterTimesMap == nil {
		hc.careWaterTimesMap = make(map[int64]uint64)
	}
	return hc.careWaterTimesMap
}

func (hc *HeroController) SetCareWaterTimesMap(toSet map[int64]uint64) {
	hc.careWaterTimesMap = toSet
}

func (hc *HeroController) GetCareCondition() *server_proto.MilitaryConditionProto {
	c := hc.careCondition.Load()
	if c != nil {
		return c.(*server_proto.MilitaryConditionProto)
	}

	return nil
}

func (hc *HeroController) SetCareCondition(toSet *server_proto.MilitaryConditionProto) {
	hc.careCondition.Store(toSet)
}

func (hc *HeroController) GetViewArea() *realmface.ViewArea {
	c := hc.viewAreaValue.Load()
	if c != nil {
		return c.(*realmface.ViewArea)
	}

	return nil
}

func (hc *HeroController) SetViewArea(toSet *realmface.ViewArea) {
	hc.viewAreaValue.Store(toSet)
}

//GetBlockIndex 获取地块索引
func (hc *HeroController) GetBlockIndex() interface{} {
	return hc.viewBlockIndex
}

//SetBlockIndex 设置地块索引(AOI)
func (hc *HeroController) SetBlockIndex(index interface{}) interface{} {
	temp := hc.viewBlockIndex
	hc.viewBlockIndex = index
	return temp
}

//RemoveBlockIndex 退出时 移除出场景
func (hc *HeroController) RemoveBlockIndex() {
	hc.viewBlockIndex = nil
}

//GetWatchObjList 获取观察列表
func (hc *HeroController) GetWatchObjList() map[interface{}]int {
	if hc.watchObjList == nil {
		return make(map[interface{}]int)
	}
	return hc.watchObjList
}

//AddWatchObjList 设置观察对象列表 如果设置nil,表示清空
func (hc *HeroController) SetWatchObjList(watchObjList map[interface{}]int) {
	hc.watchObjList = watchObjList
}

func (hc *HeroController) Id() int64 {
	return hc.id
}

func (hc *HeroController) IdBytes() []byte {
	return hc.idBytes
}

func (hc *HeroController) Pid() uint32 {
	return constants.PID
}

func (hc *HeroController) Sid() uint32 {
	return heroid.GetSid(hc.id)
}

func (hc *HeroController) IsClosed() bool {
	return hc.sender.IsClosed()
}

// 发送消息.
func (hc *HeroController) Send(msg pbutil.Buffer) {
	hc.sender.Send(msg)
}

// 发送消息.
func (hc *HeroController) SendAll(msg []pbutil.Buffer) {
	hc.sender.SendAll(msg)
}

// 发送在线路繁忙时可以被丢掉的消息
func (hc *HeroController) SendIfFree(msg pbutil.Buffer) {
	hc.sender.SendIfFree(msg)
}

func (hc *HeroController) Disconnect(err msg.ErrMsg) {
	// 不能等，因为大部分情况下都是在自己线程调用，一等就死掉了
	hc.sender.Disconnect(err)
}

func NewHeroController(id int64, sender sender.ClosableSender, clientIp string, clientIp32, pf uint32, locker herolock.HeroLocker) *HeroController {
	hc := &HeroController{
		id:                           id,
		idBytes:                      idbytes.ToBytes(id),
		sender:                       sender,
		clientIp:                     clientIp,
		pf:                           pf,
		heroLocker:                   locker,
		nextRefreshRecommendHeroTime: time.Time{},
	}

	return hc
}

func (hc *HeroController) Func(f herolock.Func) {
	hc.heroLocker.Func(f)
}

func (hc *HeroController) FuncNotError(f herolock.FuncNotError) (hasError bool) {
	return hc.heroLocker.FuncNotError(f)
}

func (hc *HeroController) FuncWithSend(f herolock.SendFunc) (hasError bool) {
	return hc.heroLocker.FuncWithSend(f, hc.sender)
}

func (hc *HeroController) LockGetGuildId() (guildId int64, ok bool) {
	ok = !hc.FuncNotError(func(hero *entity.Hero) (heroChanged bool) {
		guildId = hero.GuildId()
		return
	})

	return
}

func (hc *HeroController) LockHeroCountry() (countryId uint64) {
	hc.FuncNotError(func(hero *entity.Hero) (heroChanged bool) {
		countryId = hero.CountryId()
		return
	})
	return
}

func (hc *HeroController) NextSearchNoGuildHeros() time.Time {
	return hc.nextSearchNoGuildHeros
}

func (hc *HeroController) SetNextSearchNoGuildHeros(toSet time.Time) {
	hc.nextSearchNoGuildHeros = toSet
}

func (r *HeroController) NextRefreshRecommendHeroTime() time.Time {
	return r.nextRefreshRecommendHeroTime
}

func (r *HeroController) UpdateNextRefreshRecommendHeroTime(newTime time.Time) {
	r.nextRefreshRecommendHeroTime = newTime
}

func (r *HeroController) NextSearchHeroTime() time.Time {
	return r.nextSearchHeros
}

func (r *HeroController) UpdateNextSearchHeroTime(newTime time.Time) {
	r.nextSearchHeros = newTime
}

func (r *HeroController) LastClickTime() time.Time {
	return r.lastClickTime
}

func (r *HeroController) SetLastClickTime(newTime time.Time) {
	r.lastClickTime = newTime
}
