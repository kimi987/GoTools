package realm

import (
	"github.com/lightpaw/logrus"
	"github.com/lightpaw/male7/config/basedata"
	"github.com/lightpaw/male7/config/blockdata"
	"github.com/lightpaw/male7/config/data"
	"github.com/lightpaw/male7/config/domestic_data"
	"github.com/lightpaw/male7/config/i18n"
	"github.com/lightpaw/male7/config/maildata"
	"github.com/lightpaw/male7/config/monsterdata"
	"github.com/lightpaw/male7/config/pushdata"
	"github.com/lightpaw/male7/config/regdata"
	"github.com/lightpaw/male7/config/singleton"
	"github.com/lightpaw/male7/entity"
	"github.com/lightpaw/male7/entity/cb"
	"github.com/lightpaw/male7/entity/hexagon"
	"github.com/lightpaw/male7/entity/npcid"
	"github.com/lightpaw/male7/gen/iface"
	"github.com/lightpaw/male7/gen/pb/domestic"
	"github.com/lightpaw/male7/gen/pb/guild"
	"github.com/lightpaw/male7/gen/pb/misc"
	"github.com/lightpaw/male7/gen/pb/region"
	"github.com/lightpaw/male7/module/realm/realmerr"
	"github.com/lightpaw/male7/module/realm/realmevent"
	"github.com/lightpaw/male7/module/realm/realmface"
	"github.com/lightpaw/male7/module/xiongnu/xiongnuface"
	"github.com/lightpaw/male7/pb/server_proto"
	"github.com/lightpaw/male7/pb/shared_proto"
	"github.com/lightpaw/male7/service/conflict/heroservice/herolock"
	"github.com/lightpaw/male7/service/heromodule"
	"github.com/lightpaw/male7/service/monitor/metrics"
	"github.com/lightpaw/male7/service/operate_type"
	"github.com/lightpaw/male7/service/ticker/tickdata"
	"github.com/lightpaw/male7/util"
	"github.com/lightpaw/male7/util/atomic"
	"github.com/lightpaw/male7/util/call"
	"github.com/lightpaw/male7/util/i64"
	"github.com/lightpaw/male7/util/idbytes"
	"github.com/lightpaw/male7/util/imath"
	"github.com/lightpaw/male7/util/must"
	"github.com/lightpaw/male7/util/timer"
	"github.com/lightpaw/male7/util/timeutil"
	"github.com/lightpaw/male7/util/u64"
	"github.com/lightpaw/pbutil"
	"math"
	"runtime/debug"
	"sync"
	"time"
	"github.com/lightpaw/male7/gamelogs"
	"github.com/lightpaw/male7/constants"
	"github.com/lightpaw/male7/entity/heroid"
)

var (
	queueTimeoutWheel = timer.NewTimingWheel(500*time.Millisecond, 8) // 最大timeout是4秒
)

const (
	queueTimeout = 3 * time.Second
)

func newRealm(id int64, service *services, levelData *regdata.RegionData) *Realm {

	r := &Realm{
		id:         id,
		levelData:  levelData,
		services:   service,
		dep:        service.dep,
		conflict:   newConfllict(levelData),
		eventQueue: realmevent.NewEventQueue(),

		funcChan:       make(chan *action, 1024),
		closeNotify:    make(chan struct{}),
		loopExitNotify: make(chan struct{}),

		//bases:                   make(map[int64]*baseWithData, 64),
		baseManager:             &BaseManager{},
		resourceConflictHeroMap: make(map[cb.Cube][]int64),

		conflictBaseCount: atomic.NewUint64(0),
		blockManager:      newBlockManager(levelData.Block, levelData.InitRadius),

		astDefendLogs: NewAstDefendLogs(),
	}

	r.ruinsBasePosInfoMap = NewRuinsBasePosInfoMap(r)

	// loop要创建者手动启动, 不然开服时还没准备好就开始处理消息什么的了

	if len(levelData.GetMonsters()) > 0 {
		proto := &region.S2CNpcBaseInfoProto{
			MapId: i64.Int32(r.id),
		}

		for _, mon := range levelData.GetMonsters() {
			id := npcid.GetNpcId(mon.Id, mon.Base.Id, npcid.NpcType_Monster)
			base := r.newBasicNpcBaseWithData(id, mon.Base, mon.BaseX, mon.BaseY)

			r.addBaseToMap(base)
			r.conflict.doAddBase(base.BaseX(), base.BaseY())

			proto.NpcId = append(proto.NpcId, base.IdBytes())
			proto.BaseX = append(proto.BaseX, int32(base.BaseX()))
			proto.BaseY = append(proto.BaseY, int32(base.BaseY()))
			proto.DataId = append(proto.DataId, u64.Int32(mon.Base.Id))
		}

		r.npcBaseMsg = region.NewS2cNpcBaseInfoProtoMsg(proto).Static()
	}

	// 联盟主城地区，把联盟堡垒的位置设置上去
	if levelData.RegionType == shared_proto.RegionType_GUILD {
		r.conflict.addBaseIfCanAdd(
			int(service.datas.RegionConfig().GuildRegionCenterX),
			int(service.datas.RegionConfig().GuildRegionCenterY),
		)
	}

	return r
}

//gogen:iface entity
type Realm struct {
	id        int64
	levelData *regdata.RegionData
	services  *services
	dep       iface.ServiceDep
	*conflict

	// --- event ---
	eventQueue    realmevent.EventQueue // events 按事件要发生的时间顺序排序
	nextEventTime <-chan time.Time

	dailyTicker       tickdata.TickTime
	per10MinuteTicker tickdata.TickTime

	lastRefreshBaoZangNpcTime time.Time // 最后一次刷新宝藏Npc时间
	lastRefreshJunTuanNpcTime time.Time // 最后一次刷新军团Npc时间

	// 事件处理
	funcChan       chan *action
	directExecFunc bool // true表示不进入funcChan, 直接调用（for test,默认false）
	closeOnce      sync.Once
	closeNotify    chan struct{}
	loopExitNotify chan struct{}
	isClosed       bool // 场景线程内使用

	// 只loop可读数据, 由loop管理
	//bases               map[int64]*baseWithData // 老家, 行营 数据. 英雄id
	baseManager         *BaseManager
	basePosInfoMap      basePosInfoMap       // key: 地图x、y坐标，value: 城池id
	ruinsBasePosInfoMap *ruinsBasePosInfoMap // key: 地图x、y坐标，value: 废墟的时间，玩家id

	conflictBaseCount *atomic.Uint64 // 记录当前占格子的主城数（用于计算是否需要扩展地图）
	blockManager      *block_manager

	// 所有地图格子被多少人影响. 方便找到冲突的势力, 再方便找到无冲突的地块.
	resourceConflictHeroMap map[cb.Cube][]int64 // 每个地图块，资源点冲突的玩家id列表

	npcBaseMsg pbutil.Buffer // npc主城信息缓存

	astDefendLogs *server_proto.AllAstDefendLogProto // 盟友驻扎恢复繁荣度日志
}

func (r *Realm) config() *singleton.RegionConfig {
	return r.services.datas.RegionConfig()
}

func (r *Realm) genConfig() *singleton.RegionGenConfig {
	return r.services.datas.RegionGenConfig()
}

func (r *Realm) getBase(id int64) *baseWithData {
	return r.baseManager.getBase(id)
}

func (r *Realm) addBase(base *baseWithData) {
	r.baseManager.addBase(base)
}

func (r *Realm) deleteBase(id int64) {
	r.baseManager.removeBase(id)
}

func (r *Realm) rangeBases(f func(base *baseWithData) (toContinue bool)) {
	r.baseManager.rangeBases(f)
}

type BaseManager struct {
	baseMap sync.Map

	// 玩家召唤殷墟
	heroBaozMap            sync.Map
	heroBaozNextExpireTime int64
}

func (r *BaseManager) getHeroBaozId(heroId int64) int64 {
	if bid, ok := r.heroBaozMap.Load(heroId); ok {
		return bid.(int64)
	}
	return 0
}

func (r *Realm) GetHeroBaozRoBase(heroId int64) *server_proto.RoBaseProto {
	if id := r.baseManager.getHeroBaozId(heroId); id != 0 {
		return r.GetRoBase(id)
	}
	return nil
}

func (r *BaseManager) addHeroBaozId(heroId, baozBaseId int64) {
	r.heroBaozMap.Store(heroId, baozBaseId)
}

func (r *BaseManager) removeHeroBaozId(heroId int64) {
	r.heroBaozMap.Delete(heroId)
}

func (r *BaseManager) rangeHeroBaozIds(f func(heroId, baozBaseId int64) (toContinue bool)) {
	r.heroBaozMap.Range(func(key, value interface{}) bool {
		return f(key.(int64), value.(int64))
	})
}

func (r *BaseManager) getBase(id int64) *baseWithData {
	if b, _ := r.baseMap.Load(id); b != nil {
		return b.(*baseWithData)
	}
	return nil
}

func (r *BaseManager) addBase(base *baseWithData) {
	r.baseMap.Store(base.Id(), base)
}

func (r *BaseManager) removeBase(id int64) {
	r.baseMap.Delete(id)
}

func (r *BaseManager) rangeBases(f func(base *baseWithData) (toContinue bool)) {
	r.baseMap.Range(func(key, value interface{}) bool {
		return f(value.(*baseWithData))
	})
}

func (r *Realm) GetRoBase(id int64) *server_proto.RoBaseProto {
	if b := r.getBase(id); b != nil {
		return b.getRoBase()
	}
	return nil
}

func (r *Realm) GetRoBaseByPos(x, y int) *server_proto.RoBaseProto {
	if baseId := r.basePosInfoMap.GetBase(x, y); baseId != 0 {
		return r.GetRoBase(baseId)
	}
	return nil
}

func (r *Realm) GetRuinsBase(x, y int) int64 {
	return r.ruinsBasePosInfoMap.GetRuinsBase(x, y)
}

func (r *Realm) addBaseToMap(base *baseWithData) {
	r.addBase(base)
	r.basePosInfoMap.AddBase(base)
	r.ruinsBasePosInfoMap.OnBaseChanged(base.Base())

	//Kimi 添加base到索引
	r.blockManager.AddBaseIndex(base.BaseX(), base.BaseY(), base)

	r.blockManager.addBase(base.BaseX(), base.BaseY(), base.isHeroHomeBase())
	if baoz := GetBaoZangBase(base); baoz != nil {
		b := r.blockManager.getOrCreateBlockByBasePos(base.BaseX(), base.BaseY())
		b.increseBaozCount(baoz.Data())

		if baoz.heroType == HeroTypeCreater && baoz.heroId != 0 {
			r.baseManager.addHeroBaozId(baoz.heroId, base.Id())

			if baoz.heroEndTime > 0 {
				r.baseManager.heroBaozNextExpireTime = i64.Min(r.baseManager.heroBaozNextExpireTime, int64(baoz.heroEndTime))
			}
		}
	}

	if juntuan := GetJunTuanBase(base); juntuan != nil {
		b := r.blockManager.getOrCreateBlockByBasePos(base.BaseX(), base.BaseY())
		b.increseJunTuanCount(juntuan.Data())
	}

	r.conflictBaseCount.Inc()
}

func (r *Realm) removeBaseFromMap(base *baseWithData) {
	r.deleteBase(base.Id())
	r.basePosInfoMap.RemoveBase(base.Base())

	//Kimi 从索引中移除base
	r.blockManager.RemoveBaseIndex(base.BaseX(), base.BaseY(), base)

	// 有主城阻挡，移除主城阻挡
	r.conflict.removeBase(base.BaseX(), base.BaseY())

	r.blockManager.removeBase(base.BaseX(), base.BaseY(), base.isHeroHomeBase())
	if baoz := GetBaoZangBase(base); baoz != nil {
		b := r.blockManager.getOrCreateBlockByBasePos(base.BaseX(), base.BaseY())
		b.decreseBaozCount(baoz.Data())

		if baoz.heroType == HeroTypeCreater && baoz.heroId != 0 {
			r.baseManager.removeHeroBaozId(baoz.heroId)

			// TODO 告诉玩家，召唤的殷墟怪没有了
			//r.services.world.Send(heroId, )
		}
	}

	if juntuan := GetJunTuanBase(base); juntuan != nil {
		b := r.blockManager.getOrCreateBlockByBasePos(base.BaseX(), base.BaseY())
		b.decreseJunTuanCount(juntuan.Data())
	}

	r.conflictBaseCount.Dec()

}

func (r *Realm) Id() int64 {
	return r.id
}

func (r *Realm) GetMapData() *blockdata.StitchedBlocks {
	return r.levelData.Block
}

// 预定一个随机主城坐标, 返回的是个可以建主城的位置
func (r *Realm) ReserveNewHeroHomePos(country uint64) (ok bool, randomX, randomY int) {
	return r.randomNewBasePos(country)
}

// 预定一个随机主城坐标, 返回的是个可以建主城的位置
func (r *Realm) ReserveRandomHomePos(t realmface.RandomPointType) (ok bool, randomX, randomY int) {
	switch t {
	case realmface.RPTReborn:
		ok, randomX, randomY = r.randomRebornBasePos()
	case realmface.RPTRandom:
		ok, randomX, randomY = r.randomBasePos()
	default:
		ok, randomX, randomY = r.randomBasePos()
	}
	return
}

func (r *Realm) IsPosOpened(x, y int) bool {
	blockX, blockY := r.GetMapData().GetBlockByPos(x, y)
	return r.GetMapData().GetRadiusBlock(r.GetRadius()).ContainsBlock(blockX, blockY)
}

// 预定一个坐标
func (r *Realm) ReservePos(x, y int) (ok bool) {
	// 判断当前是否开放这个坐标

	ok = r.conflict.addBaseIfCanAdd(x, y)
	return
}

// 在同一个场景迁城，预定一个坐标
func (r *Realm) ReservePosForMoveBase(oldX, oldY, x, y int) (ok bool) {
	if oldX == x && oldY == y {
		return false
	} else {
		ok = r.conflict.moveBaseIfCanAdd(oldX, oldY, x, y)
	}
	return
}

// 取消预定的坐标. 由于某种原因, 哥来不了了.
func (r *Realm) CancelReservedPos(x, y int) {
	r.conflict.removeBase(x, y)
}

func (r *Realm) queueSaveRealmNpcBases(proto *server_proto.RegionModuleProto) {
	r.queueFunc(true, func() {
		r.saveRealmNpcBases(proto)
	})
	return
}

func (r *Realm) saveRealmNpcBases(proto *server_proto.RegionModuleProto) {
	proto.BaoZangNpcRefreshTime = timeutil.Marshal64(r.lastRefreshBaoZangNpcTime)
	proto.JunTuanNpcRefreshTime = timeutil.Marshal64(r.lastRefreshJunTuanNpcTime)

	r.rangeBases(func(base *baseWithData) (toContinue bool) {
		if !npcid.IsNpcId(base.Id()) {
			return true
		}

		// 部队
		npcProto := base.internalBase.encodeNpcProto(base, r)
		if npcProto != nil {
			proto.NpcBaseList = append(proto.NpcBaseList, npcProto)
			for _, troop := range base.selfTroops {
				if troop.startingBase == base && troop.State() != realmface.Defending {
					proto.TroopList = append(proto.TroopList, troop.doEncodeToServer(r))
				}
			}
		}
		return true
	})
	return
}

// 获得跟我土地有冲突的玩家id
func (r *Realm) GetConflictHeroIds(hc iface.HeroController) (suc bool, conflictIds []int64) {
	suc = r.queueFunc(true, func() {
		base := r.getBase(hc.Id())
		if base == nil {
			return
		}

		for _, cube := range r.services.datas.RegionConfig().GetEvenOffsetCubesIncludeLowLevel(base.BaseLevel()) {
			targetX, targetY := cube.XY()
			evenOffset := hexagon.ShiftEvenOffset(base.BaseX(), base.BaseY(), targetX, targetY)
			if !r.isResourcePointConflicted(hc.Id(), evenOffset) {
				continue
			}

			heroIds := r.resourceConflictHeroMap[evenOffset]
			for _, heroId := range heroIds {
				conflictIds = i64.AddIfAbsent(conflictIds, heroId)
			}
		}
	})

	return
}

func (r *Realm) OnHeroLogin(hc iface.HeroController) {

	r.queueFunc(false, func() {
		base := r.getBase(hc.Id())
		if base == nil {
			logrus.Debug("英雄登陆, 本地区没有base, 刚流亡?")
		}

		var addFigntNpcFunc func()
		hc.FuncWithSend(func(hero *entity.Hero, result herolock.LockResult) {

			if base == nil {
				if hero.BaseRegion() == r.Id() {
					logrus.Errorf("英雄登陆 %d %s，hero.BaseRegion()指向的场景不存在，删除主城，region: %d", hero.Id(), hero.Name(), hero.BaseRegion())

					// 场景中没有这个主城
					hero.ClearBase()
					result.Changed()

					// 踢下线
					hc.Disconnect(misc.ErrDisconectReasonFailKick)
					return
				}
			} else {

				home := GetHomeBase(base)
				if home == nil {
					// 防御性
					logrus.Error("英雄登陆, home == nil")

					// 踢下线
					hc.Disconnect(misc.ErrDisconectReasonFailKick)
					return
				}

				if dt := hero.GetNewHeroMianDisappearTime().Unix(); dt > 0 {
					ctime := r.services.timeService.CurrentTime()
					if ctime.Unix() < dt && base.MianDisappearTime() == int32(dt) {
						// 发消息
						result.Add(region.NewS2cUpdateNewHeroMianDisappearTimeMsg(int32(dt)))
					} else {
						hero.SetNewHeroMianDisappearTime(time.Time{})
						result.Changed()
					}
				}

				// 更新
				addFigntNpcFunc = r.updateHeroRealmInfo(hero, result, base, home, true, true, true)
			}
		})

		if base == nil {
			return
		}

		// 自己部队军情
		for _, t := range base.selfTroops {
			if t.assembly != nil && t.assembly.self != t && t.assembly.self.State() != realmface.Assembly {
				protoBytes := t.assembly.self.getProtoBytes(r)
				hc.Send(region.NewS2cUpdateSelfMilitaryInfoMsg(u64.Int32(entity.GetTroopIndex(t.Id())+1), t.IdBytes(), protoBytes))
			} else {
				hc.Send(t.getUpdateMsgToSelf(r))
			}
		}

		if addFigntNpcFunc != nil {
			addFigntNpcFunc()
		}

	})

}

func (r *Realm) AddXiongNuBase(info xiongnuface.RResistXiongNuInfo, originX, originY, minRange, maxRange int) (baseId int64, baseX int32, baseY int32) {
	r.queueFunc(true, func() {
		ok, x, y := r.randomXiongNuBasePos(originX, originY, minRange, maxRange)
		if !ok {
			logrus.Errorf("添加匈奴主营，刷不出一个有效的点")
			return
		}

		baseX, baseY = int32(x), int32(y)

		baseId = npcid.NewXiongNuNpcId(info.GuildId(), info.Data().Level)

		base := r.newXiongNuNpcBase(baseId, info, x, y)
		r.addBaseToMap(base)
		r.conflict.doAddBase(base.BaseX(), base.BaseY())

		// 加部队

		// 加怪物了
		for _, monster := range info.Data().AssistMonsters {
			// 通过了, 部队出发
			troopId := base.nextNpcTroopId()
			if troopId == 0 {
				logrus.Errorf("创建协助部队的时候，创建不出id")
				continue
			}

			troop := r.newDefendingNpcTroop(troopId, base, base.BaseLevel(), base, base.BaseLevel(), monster)
			troop.onChanged()

			logrus.WithField("troopid", troop.Id()).WithField("realmid", r.id).Debug("匈奴协助加入地图成功")

			// 构建成功
			base.selfTroops[troop.Id()] = troop
			base.targetingTroops[troop.Id()] = troop
			base.remindAttackOrRobCountChanged(r)

			// 广播消息
			r.broadcastMaToCared(troop, addTroopTypeInvate, 0)
		}

		// 通知所有正在看这张地图的人
		r.broadcastBaseInfoToCared(base, addBaseTypeNewHero, 0)

	})

	return
}

var removeUseMianGoodsTimeMsg = region.NewS2cUseMianGoodsMsg(0, 0).Static()

// 把玩家的基地或行营加入进来, 可以是新英雄, 可以是随机迁城, 可以是高级快速迁城. 此时英雄的城必须不是流亡状态.
// 英雄的主城或行营当前必须不能已经属于其他地区管理.
// 加入前必须已预定坐标. processed为true的话, 就算err也不需要取消预定
// isHome表示是不是主城.
// return processed 是否已处理. err 表示是否处理有错. 根据err判断错误的类型
func (r *Realm) AddBase(heroId int64, x, y int, addType realmface.AddBaseType) (processed bool, err error) {
	processed = r.queueFunc(true, func() {
		logrus.WithField("realmid", r.id).WithField("heroid", heroId).WithField("addType", addType).Debug("准备把英雄加入地图")

		if npcid.IsNpcId(heroId) {
			logrus.WithField("heroId", heroId).Debug("英雄加入地区失败, 发送的ID是个NPC")
			return
		}

		if base := r.getBase(heroId); base != nil {
			err = realmerr.ErrAddBaseAlreadyHasRealm
			logrus.WithField("heroId", heroId).
				WithField("level", base.BaseLevel()).
				WithField("prosperity", base.Prosperity()).
				WithField("region", base.RegionID()).
				WithField("x", base.BaseX()).
				WithField("y", base.BaseY()).
				WithError(err).Debug("英雄加入地区失败, 本地区已经有base了")

			// 处理一下异常情况，防止玩家初始化城池失败
			if base.Prosperity() <= 0 || base.BaseLevel() <= 0 {
				if base.Prosperity() <= 0 {
					// 加1点繁荣度
					base.AddProsperity(1)
				} else {
					// 设置成1级主城
					base.SetBaseLevel(1)
				}
			}

			r.heroFuncWithSend(heroId, func(hero *entity.Hero, result herolock.LockResult) {
				hero.UpdateBase(base.Base())
			})
			return
		}

		ctime := r.services.timeService.CurrentTime()

		var base *baseWithData
		var homeNpcBase []*baseWithData
		var sendMianMsgFunc func()
		if r.heroFuncWithSend(heroId, func(hero *entity.Hero, result herolock.LockResult) {
			if hero.BaseRegion() != 0 {
				err = realmerr.ErrAddBaseAlreadyHasRealm
				logrus.WithError(err).Debug("英雄加入地区失败")
				return
			}

			if hero.BaseLevel() <= 0 || hero.Prosperity() <= 0 {
				switch addType {
				case realmface.AddBaseHomeTransfer:
					// 迁城时, 城必须活着
					err = realmerr.ErrAddBaseHomeNotAlive
					logrus.WithError(err).Debug("英雄加入地区失败")
					return

				case realmface.AddBaseHomeNewHero:
					hero.SetBaseLevel(1)
					hero.SetProsperity(r.services.datas.HeroCreateData().Prosperity)

				case realmface.AddBaseHomeReborn:
					// 重生, 给1/3 max 繁荣度?
					prosperity := u64.Max(1, hero.Domestic().ProsperityCapcity()/3)
					hero.SetBaseLevel(1)
					hero.SetProsperity(prosperity)

				default:
					err = realmerr.ErrAddBaseHomeNotAlive
					logrus.WithField("addType", addType).WithField("baseLevel", hero.BaseLevel()).WithField("prosperity", hero.Prosperity()).Error("Realm.AddBase时, baseLevel为0, 但是竟然没有case对应")
					return
				}
			} else {
				// 城活着
				switch addType {
				case realmface.AddBaseHomeNewHero:
					logrus.WithField("baseLevel", hero.BaseLevel()).WithField("prosperity", hero.Prosperity()).Error("realm.AddBase时, type是newHero, 但是竟然有baseLevel和prosperity")

					hero.SetBaseLevel(1)
					hero.SetProsperity(r.services.datas.HeroCreateData().Prosperity)

				case realmface.AddBaseHomeReborn:
					// 可能点了很多下之类的
					err = realmerr.ErrAddBaseHomeAlive
					logrus.WithError(err).Debug("英雄加入地区失败")
					return
				}
			}

			hero.SetBaseXY(r.id, x, y)
			base = r.newBase(hero)
			r.addBaseToMap(base)
			// 不要调用conflict.addBase. 预定坐标的时候已经调用过了

			// 设置迁城恢复间隔
			if addType == realmface.AddBaseHomeReborn {
				r.giveRestoreProsperityBuf(hero, base)

				// 流亡重建免战
				disappearTime := ctime.Add(r.config().RebornMianDuration)
				hero.SetMianDisappearTime(ctime, disappearTime)
				base.SetMianDisappearTime(i64.Int32(disappearTime.Unix()))

				sendMianMsgFunc = func() {
					r.services.world.SendFunc(heroId, func() pbutil.Buffer {
						return region.NewS2cUpdateSelfMianDisappearTimeMsg(base.MianDisappearTime(), timeutil.Marshal32(ctime))
					})
				}

				// 流亡重建，清掉免战物品使用CD
				hero.SetNextUseMianGoodsTime(time.Time{})
				result.Add(removeUseMianGoodsTimeMsg)
			} else if addType == realmface.AddBaseHomeNewHero {
				// 新手免战
				disappearTime := ctime.Add(r.config().NewHeroMianDuration)
				hero.SetNewHeroMianDisappearTime(disappearTime)
				hero.SetMianDisappearTime(ctime, disappearTime)
				base.SetMianDisappearTime(i64.Int32(disappearTime.Unix()))
			}

			if home := GetHomeBase(base); home != nil {
				// 添加Npc基地
				homeNpcBase = r.newAllHomeNpcBase(hero, x, y)
				for _, v := range homeNpcBase {
					home.homeNpcBase[v.Id()] = v
				}
			}

			heromodule.UpdateTaskProgress(hero, result, shared_proto.TaskTargetType_TASK_TARGET_HOME_IN_GUILD_REGION)

			result.Changed()
			result.Ok()

			return
		}) {
			// 下面有判断 err != nil，里面去取消reservedPos
			if err == nil {
				err = realmerr.ErrLockHeroErr
				logrus.WithField("stack", string(debug.Stack())).WithError(err).Error("英雄加入地区失败")
			}
		}

		// still in loop
		if err != nil {
			return
		}

		if base == nil {
			logrus.Error("realm.AddBase, 没有err但是也没有base")
			err = realmerr.ErrAddBaseAlreadyHasRealm
			return
		}
		logrus.WithField("realmid", r.id).WithField("heroid", heroId).Debug("英雄加入地图成功")

		t := addBaseTypeUpdate
		switch addType {
		case realmface.AddBaseHomeNewHero:
			t = addBaseTypeNewHero
		case realmface.AddBaseHomeReborn:
			t = addBaseTypeReborn
		case realmface.AddBaseHomeTransfer:
			t = addBaseTypeTransfer
		}

		// 通知所有正在看这张地图的人
		r.broadcastBaseInfoToCared(base, t, heroId)

		// 发条成功消息
		switch addType {
		case realmface.AddBaseHomeReborn:
			r.services.world.Send(heroId, region.NewS2cCreateBaseMsg(i64.Int32(r.id), int32(x), int32(y), int32(base.BaseLevel()), int32(base.Prosperity())))
		case realmface.AddBaseHomeTransfer:
			r.services.world.Send(heroId, region.NewS2cFastMoveBaseMsg(i64.Int32(r.id), int32(x), int32(y), false))
		}

		// npc基地推送消息
		r.services.world.FuncHero(heroId, func(id int64, hc iface.HeroController) {
			if area := hc.GetViewArea(); area != nil {
				for _, base := range homeNpcBase {
					if base.CanSee(area) {
						hc.Send(base.NewUpdateBaseInfoMsg(r, addBaseTypeCanSee))
					}
				}
			}
		})

		// 先发建城消息，后发更新免战消息
		if sendMianMsgFunc != nil {
			sendMianMsgFunc()
		}

		// 添加资源点占用
		r.addHomeResourcePointBlock(base, true)

		// 扩展地图
		r.updateMapRadius()
	})

	if !processed || err != nil {
		r.CancelReservedPos(x, y)
	}

	return
}

func (r *Realm) updateMapRadius() {

	// 如果当前主城个数 >= 上限，则扩展地图
	count := r.conflictBaseCount.Load()

	radius := r.GetRadius()
	for i := radius; i <= r.mapData.MaxRadius(); i++ {
		rb := r.mapData.GetRadiusBlock(i)
		if count < rb.AutoExpandBaseCount {
			// 扩展到这个位置
			if i > radius {
				r.blockManager.radius.Store(i)

				// 广播给所有人，地图扩展了
				toSend := rb.GetUpdateMapRadiusMsg()
				r.services.world.WalkHero(func(id int64, hc iface.HeroController) {
					hc.Send(toSend)
				})

				// 添加Npc对象进去
				r.addMultiLevelMonsterIfAbsent(radius+1, i)
			}
			return
		}
	}
}

func (r *Realm) addMultiLevelMonsterIfAbsent(oldRadius, newRadius uint64) {

	for i := oldRadius; i <= newRadius; i++ {
		for _, xy := range r.mapData.GetRadiusBlock(i).GetRingBlockXYs() {
			x, y := xy.XY()
			ux, uy := uint64(x), uint64(y)
			sequence := regdata.BlockSequence(ux, uy)

			startX := ux * r.levelData.Block.BlockData().XLen
			startY := uy * r.levelData.Block.BlockData().YLen

			for _, data := range r.levelData.GetMultiLevelMonsters() {
				id := npcid.GetNpcId(sequence, data.Id, npcid.NpcType_MultiLevelMonster)

				baseX := startX + data.OffsetBaseX
				baseY := startY + data.OffsetBaseY

				base := r.newMultiLevelNpcBase(id, data, int(baseX), int(baseY))
				r.addBaseToMap(base)
				r.conflict.doAddBase(base.BaseX(), base.BaseY())
			}

		}
	}

}

// 给恢复繁荣度的buf
func (r *Realm) giveRestoreProsperityBuf(hero *entity.Hero, base *baseWithData) {
	if hero.Prosperity() >= hero.ProsperityCapcity() {
		return
	}

	guanFu := hero.Domestic().GetBuilding(shared_proto.BuildingType_GUAN_FU)
	if guanFu == nil {
		logrus.Errorln("官府没找到")
		return
	}

	guanFuLevelData := r.services.datas.GetGuanFuLevelData(guanFu.Level)
	if guanFuLevelData == nil {
		logrus.Errorln("官府等级没找到")
		return
	}

	endTime := r.services.timeService.CurrentTime().Add(guanFuLevelData.MoveBaseRestoreHomeProsperityDuration)
	hero.SetMoveBaseRestoreProsperityBufEndTime(endTime)
	r.services.world.Send(hero.Id(), region.NewS2cProsperityBufMsg(timeutil.Marshal32(endTime)))

	if home, ok := base.Base().(*heroHome); ok {
		home.UpdateMoveBaseBuf(hero)
	}
}

func (r *Realm) CheckCanMoveBase(heroId int64, oldX, oldY int, removeSelfTroop bool) (err error) {
	r.queueFunc(true, func() {
		err = checkCanMoveBase(r, heroId, oldX, oldY, removeSelfTroop)
	})
	return
}

func checkCanMoveBase(r *Realm, heroId int64, oldX, oldY int, removeSelfTroop bool) (err error) {
	// 新的坐标已经
	base := r.getBase(heroId)
	if base == nil {
		err = realmerr.ErrFastMoveBaseSelfNoBase
		logrus.WithError(err).Debugf("迁移同地图内基地")
		return
	}

	if base.BaseX() != oldX || base.BaseY() != oldY {
		err = realmerr.ErrFastMoveBasePosChanged
		logrus.WithError(err).Debugf("迁移同地图内基地，位置变更了")
		return
	}

	if !removeSelfTroop {
		// 不能有自己的出征部队（高级迁城可以）
		for _, t := range base.selfTroops {
			if t.state != realmface.Defending {
				err = realmerr.ErrFastMoveBaseOutside
				logrus.WithError(err).Debugf("迁移同地图内基地")
				return
			}
		}
	}

	return
}

// 在同一个地图中移动基地
// 移动前必须已预定坐标. processed为true的话, 就算err也不需要取消预定
func (r *Realm) MoveBase(hc iface.HeroController, oldX, oldY, x, y int, removeSelfTroop bool) (processed bool, err error) {

	processed = r.queueFunc(true, func() {
		if err = checkCanMoveBase(r, hc.Id(), oldX, oldY, removeSelfTroop); err != nil {
			return
		}

		if targetBaseId := r.basePosInfoMap.GetBase(x, y); targetBaseId != 0 {
			err = realmerr.ErrFastMoveBasePosChanged
			logrus.WithError(err).Debugf("迁移同地图内基地")
			return
		}

		// 新的坐标已经
		base := r.getBase(hc.Id())
		ctime := r.services.timeService.CurrentTime()

		originX, originY := base.BaseX(), base.BaseY()
		var homeNpcBaseMsgs []pbutil.Buffer
		var homeNpcBases []*baseWithData
		var guildId int64
		if r.heroFuncWithSend(base.Id(), func(hero *entity.Hero, result herolock.LockResult) {
			// 设置新的坐标
			switch t := base.Base().(type) {
			case realmface.Home:
				// 主城
				base.SetBaseXY(x, y)
				hero.UpdateHome(t)

				r.basePosInfoMap.ChangeBasePos(base.Id(), originX, originY, x, y)

				// 设置迁城恢复间隔
				r.giveRestoreProsperityBuf(hero, base)
			default:
				err = realmerr.ErrFastMoveBaseSelfNoBase
				logrus.Errorf("迁移同地图内基地，lock hero失败")
				return
			}

			// 移除原来的占位(新的位置，在外面已经申请占用了)
			r.conflict.removeBase(originX, originY)
			r.blockManager.moveBase(originX, originY, x, y, base.isHeroHomeBase())

			// 移动野怪Npc
			if home := GetHomeBase(base); home != nil {
				for _, v := range home.homeNpcBase {
					if nb, ok := v.Base().(*homeNpcBase); ok {
						oldCanSee := v.CanSee(hc.GetViewArea())

						basePos := hexagon.ShiftEvenOffset(x, y, nb.data.EvenOffsetX, nb.data.EvenOffsetY)
						v.SetBaseXY(basePos.XY())
						nb.ClearUpdateBaseInfoMsg()

						newCanSee := v.CanSee(hc.GetViewArea())

						if newCanSee {
							homeNpcBaseMsgs = append(homeNpcBaseMsgs, nb.NewUpdateBaseInfoMsg(r, addBaseTypeTransfer, v))
						} else if oldCanSee {
							homeNpcBaseMsgs = append(homeNpcBaseMsgs, region.NewS2cRemoveBaseUnitMsg(removeBaseTypeTransfer, nb.IdBytes()))
						}

						homeNpcBases = append(homeNpcBases, v)
					}
				}
			}

			guildId = hero.GuildId()

			result.Changed()
			result.Ok()

			return
		}) {
			if err == nil {
				err = realmerr.ErrLockHeroErr
			}
			logrus.Errorf("迁移同地图内基地，lock hero失败")
			return
		}

		if err != nil {
			return
		}

		base.updateRoBase()
		for _, v := range homeNpcBases {
			v.updateRoBase()
		}

		// 不是自己的部队回家
		r.updateTroopsOnRemoveBase(base, r.getTextHelp().MRDRMoveBase4a.Text, r.getTextHelp().MRDRMoveBase4d.Text, ctime)
		r.ruinsBasePosInfoMap.OnBaseChanged(base.Base())

		// 原来看得见的要移除
		r.broadcastToCaredPos(originX, originY, base.GetCantSeeMeMsg(), 0)

		// 通知所有正在看这张地图的人
		r.broadcastBaseInfoToCared(base, addBaseTypeTransfer, hc.Id())

		// 通知野怪Npc位置改变
		for _, v := range homeNpcBaseMsgs {
			hc.Send(v)
		}

		switch base.BaseType() {
		case realmface.BaseTypeHome:
			// 添加资源点占用
			r.updateHomeResourcePointBlockWhenPosChanged(base.Base(), originX, originY)
		default:
			logrus.WithField("baseType", base.BaseType()).Error("realm.MoveBase unkown BaseType")
		}

		hc.Send(region.NewS2cFastMoveBaseMsg(i64.Int32(r.id), int32(x), int32(y), false))

		// 更新控制的宝藏
		if home := GetHomeBase(base); home != nil {
			home.TryRemoveKeepBaozFarAway(r, ctime, x, y, r.config().KeepBaozMaxDistance)
		}

		// 联盟联盟信息
		if g := r.services.guildService.GetSnapshot(guildId); g != nil {
			r.services.guildService.ClearSelfGuildMsgCache(guildId)
			r.services.world.MultiSend(g.UserMemberIds, guild.SELF_GUILD_CHANGED_S2C)
		}
	})

	if !processed || err != nil {
		r.CancelReservedPos(x, y)
	}

	return
}

func (r *Realm) getTextHelp() *data.TextHelp {
	return r.services.datas.TextHelp()
}

func (r *Realm) getMailHelp() *maildata.MailHelp {
	return r.services.datas.MailHelp()
}

// 把玩家在这个地图上的基地或者行营移除. 玩家的部队必须都已不在外面. 而且不能是流亡状态且归这个地图管. (流亡状态的话, 都不归这里管)
func (r *Realm) RemoveBase(heroId int64, checkTroopsOutside bool, reason4a, reason4d *i18n.I18nRef) (processed bool, err error, baseX, baseY int) {

	processed = r.queueFunc(true, func() {
		logrus.WithField("realmid", r.id).WithField("heroid", heroId).Debug("准备把城池移出地图")

		base := r.getBase(heroId)
		if base == nil {
			err = realmerr.ErrRemoveBaseSelfNoBase
			logrus.WithError(err).Debug("移除基地")
			return
		}

		baseX, baseY = base.BaseX(), base.BaseY()
		err = r.removeBase(base, checkTroopsOutside, reason4a, reason4d, removeBaseTypeTransfer)
	})

	return
}

func (r *Realm) removeBase(base *baseWithData, checkTroopsOutside bool, reason4a, reason4d *i18n.I18nRef, removeType int32) (err error) {
	if checkTroopsOutside {
		// 不能有自己的出征部队
		for _, t := range base.selfTroops {
			if t.state != realmface.Defending {
				err = realmerr.ErrRemoveBaseOutside
				logrus.WithError(err).Debugf("移除基地")
				return
			}
		}
	}

	ctime := r.services.timeService.CurrentTime()
	if base.isHeroBaseType() {
		r.services.heroDataService.Func(base.Id(), func(hero *entity.Hero, lockErr error) (heroChanged bool) {

			if lockErr != nil {
				err = lockErr
				logrus.WithError(err).Debugf("移除基地")
				return
			}

			switch base.BaseType() {
			case realmface.BaseTypeHome:
				// 将自己的主城坐标设置成0
				hero.ClearBase()
			default:
				err = realmerr.ErrRemoveBaseUnkownBaseType
				logrus.WithError(err).Debugf("移除基地")
				return
			}

			heroChanged = true
			return
		})

		if err != nil {
			return
		}
	}

	r.removeRealmBase(base, reason4a, reason4d, removeType, ctime)

	return
}

func (r *Realm) removeRealmBaseNoReason(target *baseWithData, removeType int32, ctime time.Time) {
	r.removeRealmBase(target, nil, nil, removeType, ctime)
}

func (r *Realm) removeRealmBase(target *baseWithData, reason4a, reason4d *i18n.I18nRef, removeType int32, ctime time.Time) {

	// 更新目标队伍
	r.updateTroopsOnRemoveBase(target, reason4a, reason4d, ctime)

	if npcid.IsHomeNpcId(target.Id()) {
		if homeNpcBase := GetHomeNpcBase(target); homeNpcBase != nil {
			r.services.world.Send(homeNpcBase.ownerHeroId, region.NewS2cRemoveBaseUnitMsg(removeType, target.IdBytes()))

			if owner := r.getBase(homeNpcBase.ownerHeroId); owner != nil {
				if home := GetHomeBase(owner); home != nil {
					delete(home.homeNpcBase, target.Id())
				}
			}
		}
		return
	}

	// 删除基地
	r.removeBaseFromMap(target)

	// 发消息告诉周围的人，城没了
	broadcastMsg := region.NewS2cRemoveBaseUnitMsg(removeType, target.IdBytes()).Static()
	r.broadcastBaseToCared(target, broadcastMsg, 0)

	// 删除资源点占用
	r.removeHomeResourcePointBlock(target)

	// 玩家主城，同时检查是不是有关联的Npc主城需要移除
	if home := GetHomeBase(target); home != nil {
		home.TryRemoveAllKeepBaoz(r, ctime)
	}

	// 宝藏Npc，删除玩家的击杀数据
	if b := GetBaoZangBase(target); b != nil && b.heroId != 0 {
		if killer := r.getBase(b.heroId); killer != nil {
			if home := GetHomeBase(killer); home != nil {
				home.TryRemoveKeepBaoz(b.Id())
			}
		}
	}

	r.heroBaseFuncWithSend(target.Base(), func(hero *entity.Hero, result herolock.LockResult) {
		hero.HistoryAmount().IncreaseOne(server_proto.HistoryAmountType_AccumBaseDead)
		heromodule.UpdateTaskProgress(hero, result, shared_proto.TaskTargetType_TASK_TARGET_ACCUM_BASE_DEAD)
	})
}

func (r *Realm) getSelfTroopReason(reason *i18n.I18nRef) *i18n.I18nRef {
	if reason != nil {
		switch reason {
		case r.services.datas.TextHelp().MRDRMoveBase4a.Text:
			return r.services.datas.TextHelp().MRDRAttMoveBase4a.Text
		case r.services.datas.TextHelp().MRDRMoveBase4d.Text:
			return r.services.datas.TextHelp().MRDRAttMoveBase4d.Text
		case r.services.datas.TextHelp().MRDRBroken4a.Text:
			return r.services.datas.TextHelp().MRDRAttBroken4a.Text
		case r.services.datas.TextHelp().MRDRBroken4d.Text:
			return r.services.datas.TextHelp().MRDRAttBroken4d.Text
		}
	}
	return nil
}

// 所有往我家跑的马，回家
func (r *Realm) updateTroopsOnRemoveBase(target *baseWithData, reason4a, reason4d *i18n.I18nRef, ctime time.Time) {

	// 以这个基地为目标的马，回家
	isSetWhiteFlag := false
	for _, t := range target.targetingTroops {
		if !isSetWhiteFlag && r.trySetWhiteFlag(t) {
			// 第一个持续掠夺的玩家，插白旗
			isSetWhiteFlag = true
		}

		// 由于匈奴入侵的协助部队的startingbase跟targetBase都是匈奴主城，导致需要处理成直接移除
		if t.startingBase == t.targetBase {
			t.removeWithoutReturnCaptain(r)
			r.broadcastRemoveMaToCared(t)
			t.clearMsgs()

			// lock hero
			r.heroBaseFuncNotError(t.startingBase.Base(), func(hero *entity.Hero) (heroChanged bool) {
				hero.UpdateTroop(t, false)
				return true
			})
		} else {
			if t.State() != realmface.AssemblyArrived {
				// 掠夺结束邮件
				r.trySendTroopDoneMail(t, reason4a, reason4d, ctime)
				t.backHome(r, ctime, true, true)

				if t.assembly != nil {
					t.assembly.broadcast(r, misc.NewS2cScreenShowWordsMsg(
						r.getTextHelp().RealmAssemblyTargetDestroy.New().
							JsonString()).Static())
				}
			}
		}
	}

	if len(target.selfTroops) > 0 {
		// 处理自己的部队
		r.heroBaseFuncWithSend(target.Base(), func(hero *entity.Hero, result herolock.LockResult) {
			for _, t := range target.selfTroops {
				// 自己的主城被打爆了，处理出征部队
				t.leaveTarget(hero, result, ctime)
				r.doRemoveHeroTroop(hero, result, t, nil)
			}

			result.Ok()

			return
		})

		for _, t := range target.selfTroops {

			if assembly := t.GetAssembly(); assembly != nil {

				if assemblyTargetId := assembly.self.originTargetId;
					!npcid.IsNpcId(assemblyTargetId) &&
						assembly.self.State() == realmface.Assembly {
					// 集结等待阶段，返还次数
					assembly.self.walkAll(func(st *troop) (toContinue bool) {
						r.returnInvaseTimes(st.startingBase.Id(), assemblyTargetId, ctime)
						return true
					})
				}

				// 里面会处理成员
				r.trySendTroopDoneMail(assembly.self, r.getSelfTroopReason(reason4a), r.getSelfTroopReason(reason4a), ctime)

				// 集结部队
				if assembly.self == t {
					// 如果自己是集结创建者，那么所有的马，全部瞬间回家
					t.assembly = nil

					t.removeWithoutReturnCaptain(r)
					r.broadcastRemoveMaToCared(t)
					t.clearMsgs()

					// 集结
					for i, t := range assembly.member {
						if t == nil {
							continue
						}
						assembly.member[i] = nil

						t.removeWithoutReturnCaptain(r)
						r.broadcastRemoveMaToCared(t)
						t.clearMsgs()

						r.services.world.SendFunc(t.startingBase.Id(), func() pbutil.Buffer {
							return misc.NewS2cScreenShowWordsMsg(
								r.getTextHelp().RealmAssemblyMemberDestroy.New().
									JsonString())
						})
					}
				} else {

					// 移除自己的马，别人的马，回家
					for i, member := range assembly.member {
						if member == t {
							assembly.member[i] = nil
							t.removeWithoutReturnCaptain(r)
							r.broadcastRemoveMaToCared(t)
							t.clearMsgs()
							break
						}
					}

					assembly.broadcast(r, misc.NewS2cScreenShowWordsMsg(
						r.getTextHelp().RealmAssemblyMemberDestroy.New().
							JsonString()).Static())

					if assembly.self.State() == realmface.Assembly {
						// 集结等待状态
						fromBase := assembly.self.startingBase
						assembly.destroyAndTroopBackHome(r, fromBase.Id(), fromBase.BaseX(), fromBase.BaseY(), ctime, true)
					} else {
						assembly.self.backHome(r, ctime, true, true)
					}
				}

			} else {
				r.trySendTroopDoneMail(t, r.getSelfTroopReason(reason4a), r.getSelfTroopReason(reason4a), ctime)
				t.removeWithoutReturnCaptain(r)
				r.broadcastRemoveMaToCared(t)
				t.clearMsgs()
			}
		}

		for _, t := range target.targetingTroops {
			// 还有马没有回去，弄回家去
			t.backHome(r, ctime, true, true)
		}
	}

}

func getTroopDialogueIdByTarget(base *baseWithData) uint64 {
	if base.isHeroHomeBase() {
		return regdata.GetTroopDialogueId(shared_proto.BaseTargetType_Hero, 0)
	}

	npcType := npcid.GetNpcIdType(base.Id())
	var subType uint64

	//switch npcType {
	//case npcid.NpcType_MultiLevelMonster:
	if mln := GetMultiLevelNpcBase(base); mln != nil {
		subType = uint64(mln.data.TypeData.Type)
	}
	//}

	return regdata.GetTroopDialogueId(shared_proto.BaseTargetType(npcType), subType)
}

// 出发侦察
func (r *Realm) InvasionInvestigate(hc iface.HeroController, targetId int64) (processed bool, err error) {
	processed = r.queueFunc(true, func() {
		logrus.WithField("realmid", r.id).WithField("heroid", hc.Id()).WithField("targetid", targetId).Debug("准备把部队加入地图")

		startingBase := r.getBase(hc.Id())
		if startingBase == nil {
			err = realmerr.ErrInvasionSelfNoBase
			logrus.WithError(err).Debug("把部队加入地图失败")
			return
		}
		var targetBase *baseWithData
		if npcid.IsHomeNpcId(targetId) {
			if home := GetHomeBase(startingBase); home != nil {
				targetBase = home.homeNpcBase[targetId]
			}
		} else {
			targetBase = r.getBase(targetId)
		}
		// 对方是否为空
		if targetBase == nil {
			err = realmerr.ErrInvasionTargetNotExist
			logrus.WithError(err).Debug("把部队加入地图失败")
			return
		}
		//是否是盟友
		if targetBase.GuildId() != 0 && targetBase.GuildId() == startingBase.GuildId() {
			err = realmerr.ErrInvasionInvalidRelation
			logrus.WithError(err).Debug("把部队加入地图失败")
			return
		}
		ctime := r.services.timeService.CurrentTime()

		// 是否免战
		if targetBase.MianDisappearTime() != 0 && timeutil.Marshal32(ctime) <= targetBase.MianDisappearTime() {
			err = realmerr.ErrInvasionMian
			logrus.WithError(err).Debug("把部队加入地图失败")
			return
		}

		dialogueData := r.services.datas.GetTroopDialogueData(getTroopDialogueIdByTarget(targetBase))
		var troop *troop
		if hc.FuncWithSend(func(hero *entity.Hero, result herolock.LockResult) {
			t := hero.GetInvestigateTroop(r.services.datas.RegionConfig().InvestigateTroopId, hero.Id())

			if t == nil {
				err = realmerr.ErrInvasionInvalidTroopIndex
				logrus.WithError(err).Debug("把部队加入地图失败")
				return
			}

			troop = startingBase.selfTroops[t.Id()]
			if troop != nil {
				err = realmerr.ErrInvasionInvalidTroopIndex
				logrus.WithError(err).Debug("把部队加入地图失败")
				return
			}

			var dialogueId uint64
			if dialogueData != nil && dialogueData.IsValid(hero.TaskList().GetCompletedBaYeStage(), hero.Level()) {
				dialogueId = dialogueData.Id
			}

			encodeCaptain := hero.EncodeInvaseCaptainInfo

			// 通过了, 部队出发
			troop = r.newTroop(t.Id(), startingBase, startingBase.BaseLevel(), targetBase, targetBase.BaseLevel(),
				realmface.MovingToInvesigate, encodeCaptain, nil, 0, dialogueId, hero.GetMoveSpeedRate()+r.services.datas.RegionConfig().InvestigateSpeedup)
			// troop.npcTimes = npcTimes

			//保存数据
			t.InitInvateInfo(r.id, troop)

			result.Changed()
			result.Ok()

			return
		}) {
			if err == nil {
				err = realmerr.ErrLockHeroErr
			}
			logrus.WithField("stack", string(debug.Stack())).WithError(err).Error("把部队加入地图失败")
			return
		}

		if err != nil {
			return
		}

		if troop == nil {
			logrus.Error("realm.InvasionInvestigate, 没有err却也没有troop")
			err = realmerr.ErrInvasionEmptyGeneral
			return
		}

		logrus.WithField("troopid", troop.Id()).WithField("realmid", r.id).WithField("heroid", hc.Id()).WithField("targetid", targetId).WithField("operate", shared_proto.TroopOperate_ToInvestigate).Debug("部队加入地图成功")

		troop.onChanged()

		// 构建成功
		startingBase.selfTroops[troop.Id()] = troop
		targetBase.targetingTroops[troop.Id()] = troop

		targetBase.remindAttackOrRobCountChanged(r)

		//添加观察者
		troop.AddWatchListOnCreate(r)
		// 广播消息
		r.broadcastMaToCared(troop, addTroopTypeInvate, 0)

		// 取消免战
		// if isInvasionOperate && !npcid.IsNpcId(targetId) {
		// 	r.doRemoveBaseMian(startingBase, true)
		// }

		// 设置部队出征状态
		hc.Send(region.NewS2cUpdateSelfTroopsOutsideMsg(u64.Int32(entity.GetTroopIndex(troop.Id()))+1, true))

		// if targetBase.isHeroHomeBase() {
		// 	attackerName := r.toFlagHeroNameByGuildId(guildId, heroName)
		// 	r.services.pushService.PushFunc(shared_proto.SettingType_ST_BEEN_ATTACKING, targetBase.Id(), func(d *pushdata.PushData) (title, content string) {
		// 		return d.Title, d.ReplaceContent("{{attacker}}", attackerName)
		// 	})

		// 	if snapshot := r.services.guildService.GetSnapshot(targetBase.GuildId()); snapshot != nil {
		// 		friendName := r.toBaseFlagHeroName(targetBase, targetLevel)

		// 		r.services.pushService.MultiPushFunc(shared_proto.SettingType_ST_GUILD_MEMBER_BEEN_ATTACKING, snapshot.UserMemberIds, targetBase.Id(), func(d *pushdata.PushData) (title, content string) {
		// 			return d.Title, d.ReplaceContent("{{attacker}}", attackerName, "{{friend}}", friendName)
		// 		})
		// 	}
		// }

	})
	return
}

// 出发攻打/帮忙驱逐
// 没有err的话调用者还需要发送成功消息
func (r *Realm) Invasion(hc iface.HeroController, operate shared_proto.TroopOperate, targetId int64, targetLevel, troopIndex uint64, npcTimes uint64) (processed bool, err error) {
	processed = r.queueFunc(true, func() {
		logrus.WithField("realmid", r.id).WithField("heroid", hc.Id()).WithField("targetid", targetId).WithField("operate", operate).Debug("准备把部队加入地图")

		startingBase := r.getBase(hc.Id())
		if startingBase == nil {
			err = realmerr.ErrInvasionSelfNoBase
			logrus.WithError(err).Debug("把部队加入地图失败")
			return
		}

		var targetBase *baseWithData
		if npcid.IsHomeNpcId(targetId) {
			if home := GetHomeBase(startingBase); home != nil {
				targetBase = home.homeNpcBase[targetId]
			}
		} else {
			targetBase = r.getBase(targetId)
		}

		if targetBase == nil {
			err = realmerr.ErrInvasionTargetNotExist
			logrus.WithError(err).Debug("把部队加入地图失败")
			return
		}

		ctime := r.services.timeService.CurrentTime()

		isInvasionOperate := operate == shared_proto.TroopOperate_ToInvasion ||
			operate == shared_proto.TroopOperate_ToInvestigate

		// 检查关系
		if isInvasionOperate {
			if targetBase.GuildId() != 0 && targetBase.GuildId() == startingBase.GuildId() {
				err = realmerr.ErrInvasionInvalidRelation
				logrus.WithError(err).Debug("把部队加入地图失败")
				return
			}

			//if targetBase.isHeroTentBase() && startingBase.isHeroTentBase() {
			//	err = realmerr.ErrInvasionCannotAttackTentFromTent
			//	logrus.WithError(err).Debug("把部队加入地图失败")
			//	return
			//}

			// 免战
			if targetBase.MianDisappearTime() != 0 && timeutil.Marshal32(ctime) <= targetBase.MianDisappearTime() {
				err = realmerr.ErrInvasionMian
				logrus.WithError(err).Debug("把部队加入地图失败")
				return
			}

			if npcid.IsXiongNuNpcId(targetBase.Id()) {
				if xiongNuNpcBase, ok := targetBase.internalBase.(*xiongNuNpcBase); ok {
					if startingBase.GuildId() != xiongNuNpcBase.Info().GuildId() {
						err = realmerr.ErrInvasionInvalidTarget
						logrus.WithError(err).Debug("无法出征本联盟之外的匈奴主城")
						return
					}

					if !i64.Contains(xiongNuNpcBase.Info().GivePrizeMembers(), hc.Id()) {
						err = realmerr.ErrInvasionTodayJoinXiongNu
						logrus.WithError(err).Debug("今日已经参加过了，无法再反击匈奴主城")
						return
					}
				}
			}
		} else {
			if targetBase.GuildId() == 0 || targetBase.GuildId() != startingBase.GuildId() {
				err = realmerr.ErrInvasionInvalidRelation
				logrus.WithError(err).Debug("把部队加入地图失败")
				return
			}
		}

		// Npc出征检查
		//if targetBase.OnlyOwnerCanSee() && targetBase.OwnerHeroId() != hc.Id() {
		//	err = realmerr.ErrInvasionInvalidTarget
		//	logrus.WithError(err).Debug("把部队加入地图失败，不是自己的Npc目标")
		//	return
		//}

		dialogueData := r.services.datas.GetTroopDialogueData(getTroopDialogueIdByTarget(targetBase))

		var troop *troop
		var heroName string
		var guildId int64
		if hc.FuncWithSend(func(hero *entity.Hero, result herolock.LockResult) {
			t := hero.GetTroopByIndex(troopIndex)
			if t == nil {
				err = realmerr.ErrInvasionInvalidTroopIndex
				logrus.WithError(err).Debug("把部队加入地图失败")
				return
			}

			if t.IsOutside() {
				err = realmerr.ErrInvasionGeneralOutside
				logrus.WithError(err).Debug("把部队加入地图失败")
				return
			}

			var dialogueId uint64
			if dialogueData != nil && dialogueData.IsValid(hero.TaskList().GetCompletedBaYeStage(), hero.Level()) {
				dialogueId = dialogueData.Id
			}

			// 如果改了检查的内容, 也需要把region中调用时的检查也改了

			validCaptainCount := 0
			for _, pos := range t.Pos() {
				captain := pos.Captain()
				if captain == nil {
					continue
				}

				if captain.Soldier() <= 0 || captain.Soldier() >= math.MaxInt32 {
					continue
				}

				validCaptainCount++
			}

			if validCaptainCount == 0 {
				err = realmerr.ErrInvasionNoSoldier
				logrus.WithError(err).Debug("把部队加入地图失败")
				return
			}

			var encodeCaptain CaptainEncoder
			if isInvasionOperate {
				encodeCaptain = hero.EncodeInvaseCaptainInfo
			} else {
				encodeCaptain = hero.EncodeAssistCaptainInfo
			}

			state := realmface.NewTroopState(operate, realmface.MoveForward)

			// 通过了, 部队出发
			targetBaseLevel := targetBase.internalBase.getBaseInfoByLevel(targetLevel).GetBaseLevel()
			troop = r.newTroop(t.Id(), startingBase, startingBase.BaseLevel(), targetBase, targetBaseLevel,
				state, encodeCaptain, t.Pos(), 0, dialogueId, hero.GetMoveSpeedRate())
			troop.npcTimes = npcTimes
			t.InitInvateInfo(r.id, troop)

			switch troop.state {
			case realmface.MovingToAssist:
				heromodule.IncreTaskProgressOne(hero, result, shared_proto.TaskTargetType_TASK_TARGET_ASSIST_GUILD_MEMBER)
			case realmface.MovingToInvade:
				if npcid.IsMultiLevelMonsterNpcId(targetBase.Id()) {
					hero.HistoryAmount().IncreaseOne(server_proto.HistoryAmountType_InvadeMultiLevelMonster)
					heromodule.UpdateTaskProgress(hero, result, shared_proto.TaskTargetType_TASK_TARGET_INVADE_MULTI_LEVEL_MONSTER)
				}
			}

			if npcid.IsNpcId(targetId) {
				t := npcid.GetNpcIdType(targetId)
				switch t {
				case npcid.NpcType_HomeNpc:
					if hero.Bools().TrySet(shared_proto.HeroBoolType_BOOL_INVASE_HOME_NPC) {
						result.Add(misc.NewS2cSetHeroBoolMsg(int32(shared_proto.HeroBoolType_BOOL_INVASE_HOME_NPC)))
					}
				case npcid.NpcType_MultiLevelMonster:
					if hero.Bools().TrySet(shared_proto.HeroBoolType_BOOL_INVASE_ML_NPC) {
						result.Add(misc.NewS2cSetHeroBoolMsg(int32(shared_proto.HeroBoolType_BOOL_INVASE_ML_NPC)))

						heromodule.UpdateTaskProgress(hero, result, shared_proto.TaskTargetType_TASK_TARGET_BOOL)
					}
				}
			}

			heroName = hero.Name()
			guildId = hero.GuildId()

			if operate == shared_proto.TroopOperate_ToInvasion && targetBase.isHeroHomeBase() {
				// 打人，扣出征次数
				hero.GetInvaseHeroTimes().ReduceOneTimes(ctime, 0)
				result.Add(region.NewS2cUpdateInvaseHeroTimesMsg(hero.GetInvaseHeroTimes().StartTimeUnix32(), nil))
			}

			result.Changed()
			result.Ok()

			return
		}) {
			if err == nil {
				err = realmerr.ErrLockHeroErr
			}
			logrus.WithField("stack", string(debug.Stack())).WithError(err).Error("把部队加入地图失败")
			return
		}

		// still in loop
		if err != nil {
			return
		}

		if troop == nil {
			logrus.Error("realm.Invasion, 没有err却也没有troop")
			err = realmerr.ErrInvasionEmptyGeneral
			return
		}
		logrus.WithField("troopid", troop.Id()).WithField("realmid", r.id).WithField("heroid", hc.Id()).WithField("targetid", targetId).WithField("operate", operate).Debug("部队加入地图成功")

		troop.onChanged()

		// 构建成功
		startingBase.selfTroops[troop.Id()] = troop
		targetBase.targetingTroops[troop.Id()] = troop

		targetBase.remindAttackOrRobCountChanged(r)

		//添加观察者
		troop.AddWatchListOnCreate(r)
		// 广播消息
		r.broadcastMaToCared(troop, addTroopTypeInvate, 0)

		// 取消免战
		if isInvasionOperate && !npcid.IsNpcId(targetId) {
			r.doRemoveBaseMian(startingBase, true)
		}

		// 设置部队出征状态
		hc.Send(region.NewS2cUpdateSelfTroopsOutsideMsg(u64.Int32(entity.GetTroopIndex(troop.Id()))+1, true))

		if targetBase.isHeroHomeBase() {
			attackerName := r.toFlagHeroNameByGuildId(guildId, heroName)
			r.services.pushService.PushFunc(shared_proto.SettingType_ST_BEEN_ATTACKING, targetBase.Id(), func(d *pushdata.PushData) (title, content string) {
				return d.Title, d.ReplaceContent("{{attacker}}", attackerName)
			})

			if snapshot := r.services.guildService.GetSnapshot(targetBase.GuildId()); snapshot != nil {
				friendName := r.toBaseFlagHeroName(targetBase, targetLevel)

				r.services.pushService.MultiPushFunc(shared_proto.SettingType_ST_GUILD_MEMBER_BEEN_ATTACKING, snapshot.UserMemberIds, targetBase.Id(), func(d *pushdata.PushData) (title, content string) {
					return d.Title, d.ReplaceContent("{{attacker}}", attackerName, "{{friend}}", friendName)
				})
			}
		}
	})
	return
}

// 创建集结
func (r *Realm) CreateAssembly(hc iface.HeroController, targetId int64, targetLevel, troopIndex uint64, assemblyCount uint64, waitDuration time.Duration) (processed bool, err error) {

	if assemblyCount <= 1 || waitDuration <= 0 {
		logrus.WithField("count", assemblyCount).
			WithField("wait", waitDuration).
			Error("创建集结，无效的集结数据")
		return true, realmerr.ErrCreateAssemblyInvalidInput
	}

	processed = r.queueFunc(true, func() {
		logrus.WithField("realmid", r.id).WithField("heroid", hc.Id()).WithField("targetid", targetId).Debug("创建集结")

		startingBase := r.getBase(hc.Id())
		if startingBase == nil {
			err = realmerr.ErrCreateAssemblySelfNoBase
			logrus.Debug("创建集结，自己没有主城，流亡了？")
			return
		}

		if startingBase.GuildId() == 0 {
			err = realmerr.ErrCreateAssemblySelfNoGuild
			logrus.Debug("创建集结，自己没有联盟不能创建集结")
			return
		}

		targetBase := r.getBase(targetId)
		if targetBase == nil {
			err = realmerr.ErrCreateAssemblyTargetNotExist
			logrus.Debug("创建集结，目标不存在")
			return
		}

		ctime := r.services.timeService.CurrentTime()

		// 检查关系
		if targetBase.GuildId() != 0 && targetBase.GuildId() == startingBase.GuildId() {
			err = realmerr.ErrCreateAssemblyInvalidRelation
			logrus.Debug("创建集结，不能对盟友进行集结")
			return
		}

		// 免战
		if targetBase.MianDisappearTime() != 0 && timeutil.Marshal32(ctime) <= targetBase.MianDisappearTime() {
			err = realmerr.ErrCreateAssemblyMian
			logrus.WithError(err).Debug("把部队加入地图失败")
			return
		}

		if npcid.IsXiongNuNpcId(targetBase.Id()) {
			if xiongNuNpcBase, ok := targetBase.internalBase.(*xiongNuNpcBase); ok {
				if startingBase.GuildId() != xiongNuNpcBase.Info().GuildId() {
					err = realmerr.ErrCreateAssemblyInvalidTarget
					logrus.Debug("创建集结，无法对本联盟之外的匈奴主城进行集结")
					return
				}

				if !i64.Contains(xiongNuNpcBase.Info().GivePrizeMembers(), hc.Id()) {
					err = realmerr.ErrCreateAssemblyTodayJoinXiongNu
					logrus.Debug("创建集结，今日已经参加过了，无法再反击匈奴主城")
					return
				}
			}
		}

		// Npc出征检查
		//if targetBase.OnlyOwnerCanSee() && targetBase.OwnerHeroId() != hc.Id() {
		//	err = realmerr.ErrInvasionInvalidTarget
		//	logrus.Debug("把部队加入地图失败，不是自己的Npc目标")
		//	return
		//}

		var troop *troop
		if hc.FuncWithSend(func(hero *entity.Hero, result herolock.LockResult) {
			t := hero.GetTroopByIndex(troopIndex)
			if t == nil {
				err = realmerr.ErrCreateAssemblyInvalidTroopIndex
				logrus.WithField("index", troopIndex).Error("创建集结，t == nil")
				return
			}

			if t.IsOutside() {
				err = realmerr.ErrCreateAssemblyGeneralOutside
				logrus.Debug("创建集结，队伍出征中")
				return
			}

			// 如果改了检查的内容, 也需要把region中调用时的检查也改了

			validCaptainCount := 0
			for _, pos := range t.Pos() {
				captain := pos.Captain()
				if captain == nil {
					continue
				}

				if captain.Soldier() <= 0 || captain.Soldier() >= math.MaxInt32 {
					continue
				}

				validCaptainCount++
			}

			if validCaptainCount == 0 {
				err = realmerr.ErrCreateAssemblyNoSoldier
				logrus.Debug("创建集结，派出的队伍没有士兵")
				return
			}

			// 通过了, 部队出发
			targetBaseLevel := targetBase.internalBase.getBaseInfoByLevel(targetLevel).GetBaseLevel()
			troop = r.newTroop(t.Id(), startingBase, startingBase.BaseLevel(), targetBase, targetBaseLevel,
				realmface.Assembly, hero.EncodeInvaseCaptainInfo, t.Pos(), assemblyCount, 0, hero.GetMoveSpeedRate())

			// 设置事件
			troop.moveStartTime = ctime
			troop.moveArriveTime = ctime.Add(waitDuration)
			troop.event = r.newEvent(troop.moveArriveTime, troop.assemblyTimesUpEvent)

			troop.assembly.updateAddedStat(r, 0)

			t.InitInvateInfo(r.id, troop)

			if targetBase.isHeroHomeBase() {
				// 打人，扣出征次数
				hero.GetInvaseHeroTimes().ReduceOneTimes(ctime, 0)
				result.Add(region.NewS2cUpdateInvaseHeroTimesMsg(hero.GetInvaseHeroTimes().StartTimeUnix32(), nil))
			}

			result.Changed()
			result.Ok()

			return
		}) {
			if err == nil {
				err = realmerr.ErrLockHeroErr
			}
			logrus.WithField("stack", string(debug.Stack())).Error("创建集结，lock hero失败")
			return
		}

		// still in loop
		if err != nil {
			return
		}

		if troop == nil {
			logrus.Error("创建集结, 没有err却也没有troop")
			err = realmerr.ErrCreateAssemblyEmptyGeneral
			return
		}
		logrus.WithField("troopid", troop.Id()).WithField("realmid", r.id).WithField("heroid", hc.Id()).WithField("targetid", targetId).Debug("创建集结成功")

		troop.onChanged()

		// 构建成功
		startingBase.selfTroops[troop.Id()] = troop
		targetBase.targetingTroops[troop.Id()] = troop

		targetBase.remindAttackOrRobCountChanged(r)
		
		//添加观察者列表
		troop.AddWatchListOnCreate(r)
		// 广播消息
		r.broadcastMaToCared(troop, addTroopTypeInvate, 0)

		// 取消免战
		if !npcid.IsNpcId(targetId) {
			r.doRemoveBaseMian(startingBase, true)
		}

		// 设置部队出征状态
		hc.Send(region.NewS2cUpdateSelfTroopsOutsideMsg(u64.Int32(entity.GetTroopIndex(troop.Id()))+1, true))

		// 推送
		//if targetBase.isHeroHomeBase() {
		//	attackerName := r.toFlagHeroNameByGuildId(guildId, heroName)
		//	r.services.pushService.PushFunc(shared_proto.SettingType_ST_BEEN_ATTACKING, targetBase.Id(), func(d *pushdata.PushData) (title, content string) {
		//		return d.Title, d.ReplaceContent("{{attacker}}", attackerName)
		//	})
		//
		//	if snapshot := r.services.guildService.GetSnapshot(targetBase.GuildId()); snapshot != nil {
		//		friendName := r.toBaseFlagHeroName(targetBase, targetLevel)
		//
		//		r.services.pushService.MultiPushFunc(shared_proto.SettingType_ST_GUILD_MEMBER_BEEN_ATTACKING, snapshot.UserMemberIds, targetBase.Id(), func(d *pushdata.PushData) (title, content string) {
		//			return d.Title, d.ReplaceContent("{{attacker}}", attackerName, "{{friend}}", friendName)
		//		})
		//	}
		//}

	})
	return
}

var showNotExistMsg = region.NewS2cShowAssemblyMsg(true, nil, 0, nil).Static()

// 查看集结
func (r *Realm) ShowAssembly(hc iface.HeroController, targetId, targetTroopId int64, version int32) {

	r.queueFunc(false, func() {
		assemblyBase := r.getBase(targetId)
		if assemblyBase == nil {
			logrus.WithField("target", targetId).WithField("troop", targetTroopId).Debug("展示集结，目标不存在")
			hc.Send(showNotExistMsg)
			return
		}

		assemblyTroop := assemblyBase.selfTroops[targetTroopId]
		if assemblyTroop == nil {
			logrus.WithField("target", targetId).WithField("troop", targetTroopId).Debug("展示集结，目标队伍不存在")
			hc.Send(showNotExistMsg)
			return
		}

		assembly := assemblyTroop.assembly
		if assembly == nil {
			logrus.WithField("target", targetId).WithField("troop", targetTroopId).Debug("展示集结，目标队伍没有集结")
			hc.Send(showNotExistMsg)
			return
		}

		t := assembly.self

		proto := &shared_proto.AssemblyInfoProto{}
		_, moveType := stateToActionMoveType(t.State())
		proto.MoveType = moveType.Int32()
		proto.MoveStartTime = timeutil.Marshal32(t.moveStartTime)
		proto.MoveArrivedTime = timeutil.Marshal32(t.moveArriveTime)
		proto.RobbingEndTime = timeutil.Marshal32(t.robbingEndTime)

		proto.Self = t.getStartingBaseSnapshot(r)
		proto.Self.JunXianLevel = u64.Int32(r.services.baiZhanService.GetJunXianLevel(t.startingBase.Id()))
		proto.SelfFightAmount = u64.Int32(t.FightAmount())

		if t.targetBase != nil {
			baseInfo := t.getTargetBaseInfo()
			if baseInfo != nil {
				proto.Target = baseInfo.EncodeAsHeroBasicSnapshot(r.services.heroSnapshotService.Get)
			}

			npcDataId, npcType := t.targetBase.internalBase.getNpcConfig()
			proto.TargetNpcDataId, proto.TargetNpcType = u64.Int32(npcDataId), npcType

			var targetTroop, targetTotalTroop int32
			if xiongNu := GetXiongNuBase(t.targetBase); xiongNu != nil {
				proto.TargetMorale = u64.Int32(xiongNu.info.Morale())

				targetTroop = int32(len(t.targetBase.getAssisterTroops()) + 1)
				targetTotalTroop = int32(len(xiongNu.Data().AssistMonsters) + 1)
			}

			if junTuan := GetJunTuanBase(t.targetBase); junTuan != nil {
				targetTroop = int32(len(t.targetBase.getAssisterTroops()))
				targetTotalTroop = int32(junTuan.Data().TroopCount)
			}

			proto.TargetTroop = targetTroop
			proto.TargetTotalTroop = targetTotalTroop
		}
		proto.TargetBaseX = int32(t.originTargetX)
		proto.TargetBaseY = int32(t.originTargetY)

		for _, member := range assembly.member {
			if member != nil {
				mp := &shared_proto.AssemblyMemberProto{
					TroopId:     member.IdBytes(),
					Hero:        member.getStartingBaseSnapshot(r),
					FightAmount: u64.Int32(member.FightAmount()),
				}

				if member.State() != realmface.AssemblyArrived {
					mp.MoveStartTime = timeutil.Marshal32(member.moveStartTime)
					mp.MoveArrivedTime = timeutil.Marshal32(member.moveArriveTime)
				}

				mp.Hero.JunXianLevel = u64.Int32(r.services.baiZhanService.GetJunXianLevel(member.startingBase.Id()))

				proto.Member = append(proto.Member, mp)
			}
		}
		proto.TotalCount = u64.Int32(assembly.TotalCount())

		version := int32(time.Now().Unix()) // TODO
		hc.Send(region.NewS2cShowAssemblyMsg(false, idbytes.ToBytes(targetTroopId), version, must.Marshal(proto)))
	})

}

// 加入集结
func (r *Realm) JoinAssembly(hc iface.HeroController, targetId, targetTroopId int64, troopIndex uint64) (processed bool, err error) {

	if hc.Id() == targetId {
		logrus.Debug("加入集结，不能加入自己创建的集结")
		return true, realmerr.ErrJoinAssemblyInvalidTarget
	}

	processed = r.queueFunc(true, func() {
		logrus.WithField("realmid", r.id).WithField("heroid", hc.Id()).WithField("targetid", targetId).Debug("加入集结")

		startingBase := r.getBase(hc.Id())
		if startingBase == nil {
			err = realmerr.ErrJoinAssemblySelfNoBase
			logrus.Debug("加入集结，自己没有主城，流亡了？")
			return
		}

		targetBase := r.getBase(targetId)
		if targetBase == nil {
			err = realmerr.ErrJoinAssemblyTargetNotExist
			logrus.Debug("加入集结，目标不存在")
			return
		}

		ctime := r.services.timeService.CurrentTime()

		// 检查关系
		if targetBase.GuildId() == 0 || targetBase.GuildId() != startingBase.GuildId() {
			err = realmerr.ErrJoinAssemblyInvalidRelation
			logrus.Debug("加入集结，只能加入盟友的集结")
			return
		}

		if !targetBase.isHeroHomeBase() {
			err = realmerr.ErrJoinAssemblyInvalidRelation
			logrus.Debug("加入集结，只能加入玩家的集结")
			return
		}

		targetTroop := targetBase.selfTroops[targetTroopId]
		if targetTroop == nil {
			err = realmerr.ErrJoinAssemblyTargetNotExist
			logrus.Debug("加入集结，集结部队不存在")
			return
		}

		if targetTroop.State() != realmface.Assembly {
			err = realmerr.ErrJoinAssemblyStarted
			logrus.Debug("加入集结，集结已经出发")
			return
		}

		assembly := targetTroop.GetAssembly()
		if assembly == nil || assembly.self != targetTroop {
			err = realmerr.ErrJoinAssemblyTargetNotExist
			logrus.Debug("加入集结，不是部队集结者")
			return
		}

		if assembly.Count() >= assembly.TotalCount() {
			err = realmerr.ErrJoinAssemblyFull
			logrus.Debug("加入集结，集结已满")
			return
		}

		for _, m := range assembly.member {
			if m != nil && m.startingBase == startingBase {
				err = realmerr.ErrJoinAssemblyMultiJoin
				logrus.Debug("加入集结，不能多个队伍加入同一个集结")
				return
			}
		}

		if targetBaseOfTarget := targetTroop.targetBase; targetBaseOfTarget != nil {
			targetBase := targetBaseOfTarget

			if npcid.IsXiongNuNpcId(targetBase.Id()) {
				if xiongNuNpcBase, ok := targetBase.internalBase.(*xiongNuNpcBase); ok {
					if startingBase.GuildId() != xiongNuNpcBase.Info().GuildId() {
						err = realmerr.ErrJoinAssemblyInvalidTarget
						logrus.Debug("加入集结，无法对本联盟之外的匈奴主城进行集结")
						return
					}

					if !i64.Contains(xiongNuNpcBase.Info().GivePrizeMembers(), hc.Id()) {
						err = realmerr.ErrJoinAssemblyTodayJoinXiongNu
						logrus.Debug("加入集结，今日已经参加过了，无法再反击匈奴主城")
						return
					}
				}
			}
		}

		var troop *troop
		if hc.FuncWithSend(func(hero *entity.Hero, result herolock.LockResult) {
			t := hero.GetTroopByIndex(troopIndex)
			if t == nil {
				err = realmerr.ErrJoinAssemblyInvalidTroopIndex
				logrus.WithField("index", troopIndex).Error("加入集结，t == nil")
				return
			}

			if t.IsOutside() {
				err = realmerr.ErrJoinAssemblyGeneralOutside
				logrus.Debug("加入集结，队伍出征中")
				return
			}

			// 如果改了检查的内容, 也需要把region中调用时的检查也改了

			validCaptainCount := 0
			for _, pos := range t.Pos() {
				captain := pos.Captain()
				if captain == nil {
					continue
				}

				if captain.Soldier() <= 0 || captain.Soldier() >= math.MaxInt32 {
					continue
				}

				validCaptainCount++
			}

			if validCaptainCount == 0 {
				err = realmerr.ErrJoinAssemblyNoSoldier
				logrus.Debug("加入集结，派出的队伍没有士兵")
				return
			}

			// 通过了, 部队出发
			troop = r.newTroop(t.Id(), startingBase, startingBase.BaseLevel(), targetBase, targetBase.BaseLevel(),
				realmface.MovingToAssembly, hero.EncodeInvaseCaptainInfo, t.Pos(), 0, 0, hero.GetMoveSpeedRate())
			troop.assembly = assembly
			assembly.addTroop(troop)
			t.InitInvateInfo(r.id, troop)

			if !npcid.IsNpcId(targetTroop.originTargetId) {
				// 打人，扣出征次数
				hero.GetInvaseHeroTimes().ReduceOneTimes(ctime, 0)
				result.Add(region.NewS2cUpdateInvaseHeroTimesMsg(hero.GetInvaseHeroTimes().StartTimeUnix32(), nil))
			}

			result.Changed()
			result.Ok()

			return
		}) {
			if err == nil {
				err = realmerr.ErrLockHeroErr
			}
			logrus.WithField("stack", string(debug.Stack())).Error("加入集结，lock hero失败")
			return
		}

		// still in loop
		if err != nil {
			return
		}

		if troop == nil {
			logrus.Error("加入集结, 没有err却也没有troop")
			err = realmerr.ErrJoinAssemblyEmptyGeneral
			return
		}
		logrus.WithField("troopid", troop.Id()).WithField("realmid", r.id).WithField("heroid", hc.Id()).WithField("targetid", targetId).Debug("加入集结成功")

		troop.onChanged()
		targetTroop.onChanged()

		// 构建成功
		startingBase.selfTroops[troop.Id()] = troop
		targetBase.targetingTroops[troop.Id()] = troop

		//targetBase.remindAttackOrRobCountChanged(r)

		//添加观察者列表
		troop.AddWatchListOnCreate(r)
		// 广播消息
		r.broadcastMaToCared(troop, addTroopTypeInvate, 0)

		// 取消免战
		if !npcid.IsNpcId(targetTroop.originTargetId) {
			r.doRemoveBaseMian(startingBase, true)
		}

		// 设置部队出征状态
		hc.Send(region.NewS2cUpdateSelfTroopsOutsideMsg(u64.Int32(entity.GetTroopIndex(troop.Id()))+1, true))

		assembly.broadcastChanged(r)

		// 推送
		//if targetBase.isHeroHomeBase() {
		//	attackerName := r.toFlagHeroNameByGuildId(guildId, heroName)
		//	r.services.pushService.PushFunc(shared_proto.SettingType_ST_BEEN_ATTACKING, targetBase.Id(), func(d *pushdata.PushData) (title, content string) {
		//		return d.Title, d.ReplaceContent("{{attacker}}", attackerName)
		//	})
		//
		//	if snapshot := r.services.guildService.GetSnapshot(targetBase.GuildId()); snapshot != nil {
		//		friendName := r.toBaseFlagHeroName(targetBase, targetLevel)
		//
		//		r.services.pushService.MultiPushFunc(shared_proto.SettingType_ST_GUILD_MEMBER_BEEN_ATTACKING, snapshot.UserMemberIds, targetBase.Id(), func(d *pushdata.PushData) (title, content string) {
		//			return d.Title, d.ReplaceContent("{{attacker}}", attackerName, "{{friend}}", friendName)
		//		})
		//	}
		//}
	})
	return

}

func (r *Realm) AddInvasionMonster(heroId int64, npcType shared_proto.MultiLevelNpcType, level uint64) {

	r.queueFunc(false, func() {

		// 怪物攻城，添加出征部队
		targetBase := r.getBase(heroId)
		if targetBase == nil {
			logrus.Debug("怪物攻城添加野怪，玩家主城找不到")
			return
		}

		// 根据玩家当前的坐标，找到一个出征的Npc
		centerX, centerY := r.mapData.GetBlockByPos(targetBase.BaseX(), targetBase.BaseY())

		blockXYs := blockdata.GetSpiralBlockXYs(int(centerX), int(centerY), 2, r.mapData.IsValidIntBlock)
		cb.Mix(blockXYs) // 打乱顺序

		rb := r.mapData.GetRadiusBlock(r.GetRadius())

		// 循环找怪
		var startingBase *baseWithData
		var newTroopId int64
	out:
		for _, c := range blockXYs {
			x, y := c.XY()
			ux, uy := uint64(x), uint64(y)
			if !rb.ContainsBlock(ux, uy) {
				// 未开放
				continue
			}

			sequence := regdata.BlockSequence(ux, uy)

			for _, data := range r.levelData.GetMultiLevelMonsters() {
				if data.TypeData.Type == npcType {
					// 找这只怪
					id := npcid.GetNpcId(sequence, data.Id, npcid.NpcType_MultiLevelMonster)

					base := r.getBase(id)
					if base == nil {
						logrus.WithField("x", ux).WithField("y", uy).
							WithField("data", data.Id).
							Error("怪物攻城添加野怪，根据id找不到野怪")
						continue
					}

					// 怪物太多了，换一个
					if len(base.selfTroops) >= npcid.TroopMaxSequence {
						continue
					}

					// 小于最小距离，换一个
					//if u64.FromInt(hexagon.Distance(base.BaseX(), base.BaseY(), targetBase.BaseX(), targetBase.BaseY())) < data.TypeData.FightMustDistance {
					//	continue
					//}
					// 使用客户端的计算方式，开方
					if util.IsInRange(base.BaseX(), base.BaseY(), targetBase.BaseX(), targetBase.BaseY(), int(data.TypeData.FightMustDistance)) {
						// 200距离以内
						continue
					}

					troopId := base.nextNpcTroopId()
					if troopId == 0 {
						continue
					}

					// 找到了，就这个怪了
					startingBase = base
					newTroopId = troopId

					break out
				}
			}
		}

		if startingBase == nil || newTroopId == 0 {
			// 找不到怪，算了
			logrus.Error("怪物攻城添加野怪，找不到符合条件的怪物主城")
			return
		}

		// 找到出征开始的地方
		monsterLevelData := startingBase.internalBase.getBaseInfoByLevel(level).getHateData()
		if monsterLevelData == nil {
			logrus.Error("怪物攻城添加野怪，野怪城池找不到这个等级的数据")
			return
		}

		// 添加出征野怪
		r.addInvasionFightNpcTroop(newTroopId, startingBase, targetBase, level, monsterLevelData.FightNpc, 0)

	})

}

func (r *Realm) addInvasionFightNpcTroop(newTroopId int64, startingBase, targetBase *baseWithData, startingBaseLevel uint64, mon *monsterdata.MonsterMasterData, arriveOffset time.Duration) {

	troop := r.newMoveNpcTroop(newTroopId, startingBase, startingBaseLevel, targetBase, targetBase.BaseLevel(), mon, true, arriveOffset)

	troop.onChanged()

	// 构建成功
	startingBase.selfTroops[troop.Id()] = troop
	targetBase.targetingTroops[troop.Id()] = troop

	targetBase.remindAttackOrRobCountChanged(r)

	//添加观察者列表
	troop.AddWatchListOnCreate(r)
	// 广播消息
	r.broadcastMaToCared(troop, addTroopTypeInvate, 0)
}

func (r *Realm) CheckIsFucked(id int64) bool {
	base := r.getBase(id)
	if base == nil || len(base.targetingTroops) <= 0 {
		return false
	}
	for _, troop := range base.targetingTroops {
		if troop.state == realmface.Robbing {
			return true
		}
	}
	return false
}

// 自己驱逐自己城里的坏人
func (r *Realm) Expel(hc iface.HeroController, targetId int64, troopIndex uint64) (processed, expelSuccess bool, link string, err error) {
	processed = r.queueFunc(true, func() {
		logrus.WithField("realmid", r.id).WithField("heroid", hc.Id()).WithField("targetid", targetId).Debug("驱逐掠夺者")

		startingBase := r.getBase(hc.Id())
		if startingBase == nil {
			err = realmerr.ErrExpelSelfNoBase
			logrus.WithError(err).Debug("驱逐敌人")
			return
		}

		robbing := startingBase.targetingTroops[targetId]
		if robbing == nil {
			err = realmerr.ErrExpelTroopsNotFound
			logrus.WithError(err).Debug("驱逐敌人")
			return
		}

		if robbing.State() != realmface.Robbing {
			err = realmerr.ErrExpelTroopsNoRobbing
			logrus.WithError(err).Debug("驱逐敌人")
			return
		}

		var troop *troop
		r.services.heroDataService.Func(hc.Id(), func(hero *entity.Hero, lockErr error) (heroChanged bool) {
			if lockErr != nil {
				err = lockErr
				return
			}

			t := hero.GetTroopByIndex(troopIndex)
			if t == nil {
				err = realmerr.ErrExpelInvalidTroopIndex
				logrus.WithError(err).Debug("驱逐敌人")
				return
			}

			if t.IsOutside() {
				err = realmerr.ErrExpelCaptainOutside
				logrus.WithError(err).Debug("驱逐敌人")
				return
			}

			validCaptainCount := 0
			for _, pos := range t.Pos() {
				captain := pos.Captain()
				if captain == nil {
					continue
				}

				if captain.Soldier() <= 0 || captain.Soldier() >= math.MaxInt32 {
					continue
				}

				validCaptainCount++
			}

			if validCaptainCount == 0 {
				err = realmerr.ErrExpelNoSoldier
				logrus.WithError(err).Debug("驱逐敌人")
				return
			}

			// 通过了, 部队出发
			troop = r.newTempTroop(t.Id(), startingBase, startingBase.BaseLevel(), hero.EncodeDefenseCaptainInfo, t.Pos())
			return
		})

		if err != nil {
			return
		}

		if assembly := robbing.assembly; assembly == nil {
			// 生成一只自己的部队，跟那只部队干一架
			ctx := &fightContext{}
			fightErr, fightSuccess, response, _ := r.fight(ctx, robbing, troop, fightTypeExpel, false)
			if fightErr {
				err = realmerr.ErrExpelFightError
				logrus.WithError(err).Debug("驱逐敌人")
				return
			}

			link = response.Link
			expelSuccess = !fightSuccess
		} else {
			// 驱逐的目标是个集结目标，跟全部人干架
			defenserBase := startingBase
			fightErr, isAttackerAlive, isDefenserAlive, _ := r.fightAssembly(robbing, defenserBase, nil, troop, fightTypeExpel)
			if fightErr {
				err = realmerr.ErrExpelFightError
				logrus.WithError(err).Debug("驱逐集结敌人")
				return
			}

			if isDefenserAlive {
				// 进攻方输了，看下进攻方是否有兵活着，如果有，回家
				if isAttackerAlive {
					ctime := r.services.timeService.CurrentTime()
					// 解散集结，返回自己家
					assembly.destroyAndTroopBackHome(r, defenserBase.Id(), defenserBase.BaseX(), defenserBase.BaseY(), ctime, true)
					return
				}
			}

			expelSuccess = isDefenserAlive
		}

	})
	return
}

// 班师回朝
func (r *Realm) CancelInvasion(hc iface.HeroController, troopId int64) (processed bool, err error) {
	processed = r.queueFunc(true, func() {
		base := r.getBase(hc.Id())
		if base == nil {
			err = realmerr.ErrCancelInvasionTroopNotFound
			return
		}

		troop := base.selfTroops[troopId]
		if troop == nil {
			err = realmerr.ErrCancelInvasionTroopNotFound
			return
		}

		ctime := r.services.timeService.CurrentTime()

		assembly := troop.GetAssembly()

		if troop.State() == realmface.Assembly {
			// 解散集结
			logrus.WithField("troop", troopId).Debug("集结解散")

			startingBase := base
			if assembly != nil {
				createrName := assembly.self.getStartingBaseInfo().getName()

				r.trySendTroopDoneMail(troop, r.getTextHelp().MRDRRecall4a.Text, r.getTextHelp().MRDRRecall4d.Text, ctime)

				for i, t := range assembly.member {
					if t == nil {
						continue
					}

					assembly.removeTroopByIndex(i)
					t.updateCaptainStat(assembly.addedStat, nil)

					rate := timeutil.Rate(t.moveStartTime, t.moveArriveTime, ctime)
					t.backHomeFrom(r, startingBase.Id(), startingBase.BaseX(), startingBase.BaseY(), rate, ctime)

					r.broadcastMaToCared(t, addTroopTypeUpdate, 0)

					r.heroBaseFuncNotError(t.startingBase.Base(), func(hero *entity.Hero) (heroChanged bool) {
						hero.UpdateTroop(t, false)
						return true
					})

					r.services.world.SendFunc(t.startingBase.Id(), func() pbutil.Buffer {
						return misc.NewS2cScreenShowWordsMsg(
							r.getTextHelp().RealmAssemblyCanceled.New().
								WithHero(createrName).
								WithTroopIndex(entity.GetTroopIndex(t.Id()) + 1).
								JsonString())
					})
				}
			}

			troop.removeWithoutReturnCaptain(r)
			r.broadcastRemoveMaToCared(troop)
			troop.clearMsgs()

			r.heroBaseFuncWithSend(startingBase.Base(), func(hero *entity.Hero, result herolock.LockResult) {
				r.doRemoveHeroTroop(hero, result, troop, nil)
				result.Ok()
				return
			})

			assembly.broadcastChanged(r)
			return
		}

		destroyAssembly := false
		switch troop.state {
		case realmface.InvadeMovingBack, realmface.AssistMovingBack, realmface.AssemblyMovingBack, realmface.InvesigateMovingBack:
			err = realmerr.ErrCancelInvasionTroopAlreadyBacking
			return
		case realmface.Defending:
			if troop.targetBase == nil {
				err = realmerr.ErrCancelInvasionTroopAlreadyHome
				return
			}
			destroyAssembly = true
		case realmface.Robbing:
			// 持续掠夺中
			destroyAssembly = true
		case realmface.AssemblyArrived:
			// 如果集结已经出发，则不能召回
			if assembly == nil || assembly.self.State() != realmface.Assembly {
				err = realmerr.ErrCancelInvasionTroopAssemblyStarted
				return
			}
			fallthrough
		case realmface.MovingToAssembly:
			// 还次数
			returnInvaseTimes(hc, troop.originTargetId, ctime)

		default:
			// 要么正在移动去攻击, 要么正在移动去帮忙. 反正都是在路上, 还没有到
		}

		logrus.WithField("troop", troopId).Debug("召回")

		r.trySetWhiteFlag(troop)
		r.trySendTroopDoneMail(troop, r.getTextHelp().MRDRRecall4a.Text, r.getTextHelp().MRDRRecall4d.Text, ctime)

		troop.backHome(r, ctime, destroyAssembly, true)

		r.services.world.SendFunc(troop.startingBase.Id(), func() pbutil.Buffer {
			return misc.NewS2cScreenShowWordsMsg(
				r.getTextHelp().RealmTroopRecall.New().
					WithTroopIndex(entity.GetTroopIndex(troop.Id()) + 1).
					JsonString())
		})

		if assembly != nil && assembly.self != troop && assembly.self.State() == realmface.Assembly {
			assembly.self.onChanged()
			r.broadcastMaToCared(assembly.self, addTroopTypeUpdate, 0)
		}
	})
	return
}

func (r *Realm) returnInvaseTimes(heroId, targetId int64, ctime time.Time) {
	returnFunc := getReturnInvaseTimesFunc(targetId, ctime)
	if returnFunc != nil && !npcid.IsNpcId(heroId) {
		r.services.heroDataService.FuncWithSend(heroId, returnFunc)
	}
}

// 返还出征次数
func returnInvaseTimes(hc iface.HeroController, targetId int64, ctime time.Time) {

	returnFunc := getReturnInvaseTimesFunc(targetId, ctime)
	if returnFunc != nil {
		hc.FuncWithSend(returnFunc)
	}
}

func getReturnInvaseTimesFunc(targetId int64, ctime time.Time) func(hero *entity.Hero, result herolock.LockResult) {

	if !npcid.IsNpcId(targetId) {
		return func(hero *entity.Hero, result herolock.LockResult) {
			hero.GetInvaseHeroTimes().AddTimes(1, ctime, 0)
			result.Add(region.NewS2cUpdateInvaseHeroTimesMsg(hero.GetInvaseHeroTimes().StartTimeUnix32(), nil))
		}
	}

	return nil
}

// 遣返
func (r *Realm) Repatriate(hc iface.HeroController, troopId int64) (processed bool, err error) {
	processed = r.queueFunc(true, func() {
		base := r.getBase(hc.Id())
		if base == nil {
			err = realmerr.ErrRepatriateTroopNotFound
			logrus.WithError(err).Debugf("遣返盟友部队")
			return
		}

		troop := base.targetingTroops[troopId]
		if troop == nil {
			err = realmerr.ErrRepatriateTroopNotFound
			logrus.WithError(err).Debugf("遣返盟友部队")
			return
		}

		var returnInvaseTimesTargetId int64
		assembly := troop.GetAssembly()
		if assembly != nil {
			if assembly.self.startingBase != base {
				err = realmerr.ErrRepatriateNotAssemblyCreater
				logrus.Debug("遣返盟友集结部队，你不是集结创建者")
				return
			}

			if assembly.self.State() != realmface.Assembly {
				err = realmerr.ErrRepatriateAssemblyStarted
				logrus.Debug("遣返盟友集结部队，集结已出发")
				return
			}

			returnInvaseTimesTargetId = assembly.self.originTargetId
		} else {
			if troop.State() != realmface.Defending {
				err = realmerr.ErrRepatriateTroopNoDefending
				logrus.WithError(err).Debugf("遣返盟友部队")
				return
			}
		}

		ctime := r.services.timeService.CurrentTime()
		troop.backHome(r, ctime, false, true)

		r.services.world.SendFunc(troop.startingBase.Id(), func() pbutil.Buffer {
			return misc.NewS2cScreenShowWordsMsg(
				r.getTextHelp().RealmTroopRepatriate.New().
					WithTroopIndex(entity.GetTroopIndex(troop.Id()) + 1).
					JsonString())
		})

		if returnInvaseTimesTargetId != 0 {
			// 返还次数
			r.returnInvaseTimes(troop.startingBase.Id(), returnInvaseTimesTargetId, ctime)
		}

		if assembly != nil {
			assembly.self.onChanged()
			r.broadcastMaToCared(assembly.self, addTroopTypeUpdate, 0)
		}
	})
	return
}

// 宝藏遣返
func (r *Realm) BaozRepatriate(hc iface.HeroController, baozBaseId, troopId int64) (processed bool, err error) {

	processed = r.queueFunc(true, func() {
		base := r.getBase(baozBaseId)
		if base == nil {
			err = realmerr.ErrBaozRepatriateTroopNotFound
			logrus.Debugf("宝藏遣返，base == nil")
			return
		}

		baoz := GetBaoZangBase(base)
		if baoz == nil {
			err = realmerr.ErrBaozRepatriateTroopNotFound
			logrus.Debugf("宝藏遣返，不是宝藏baseId")
			return
		}

		ctime := r.services.timeService.CurrentTime()
		switch baoz.heroType {
		case HeroTypeKiller:
			if baoz.heroId != hc.Id() {
				err = realmerr.ErrBaozRepatriateBaozNotKeep
				logrus.Debugf("宝藏遣返，不是你控制的宝藏，不能遣返")
				return
			}

			if baoz.heroEndTime < timeutil.Marshal32(ctime) {
				err = realmerr.ErrBaozRepatriateBaozNotKeep
				logrus.Debugf("宝藏遣返，宝藏控制时间已过期")
				return
			}
		case HeroTypeCreater:
			if baoz.heroId != hc.Id() {
				err = realmerr.ErrBaozRepatriateBaozNotKeep
				logrus.Debugf("宝藏遣返，不是你召唤的宝藏")
				return
			}
		default:
			err = realmerr.ErrBaozRepatriateBaozNotKeep
			logrus.Debugf("宝藏遣返，未知的heroType")
			return
		}

		troop := base.targetingTroops[troopId]
		if troop == nil {
			err = realmerr.ErrBaozRepatriateTroopNotFound
			logrus.Debugf("宝藏遣返，troop == nil")
			return
		}

		if troop.startingBase.Id() == hc.Id() {
			err = realmerr.ErrBaozRepatriateTroopNotFound
			logrus.Debugf("宝藏遣返，不能遣返自己的部队")
			return
		}

		if !troop.State().IsInvateState() {
			err = realmerr.ErrBaozRepatriateTroopNotFound
			logrus.Debugf("宝藏遣返，队伍不是出征或者掠夺状态")
			return
		}

		if troop.targetBase == nil {
			err = realmerr.ErrBaozRepatriateTroopNotFound
			logrus.Error("宝藏遣返，troop.targetBase == nil")
			return
		}

		isRobbing := troop.State() == realmface.Robbing

		troop.backHome(r, ctime, false, true)

		killerFlagName := r.toFlagHeroNameByHeroId(hc.Id())

		attackerId := troop.startingBase.Id()

		// 飘字
		r.services.world.SendFunc(attackerId, func() pbutil.Buffer {
			return misc.NewS2cScreenShowWordsMsg(
				r.getTextHelp().RealmTroopBaozRepatriate.New().
					WithTroopIndex(entity.GetTroopIndex(troop.Id()) + 1).
					WithKiller(killerFlagName).
					JsonString())
		})

		// 邮件
		data := r.services.datas.MailHelp().ReportBaozRepatriateMoving
		if isRobbing {
			data = r.services.datas.MailHelp().ReportBaozRepatriateRobber
		}

		if data != nil {
			attackerName := r.toBaseFlagHeroName(troop.startingBase, 0)
			defenserName := baoz.data.Npc.Name

			mailProto := data.NewTextMail(shared_proto.MailType_MailNormal)
			mailProto.Text = data.NewTextFields().
				WithAttacker(attackerName).
				WithDefenser(defenserName).
				WithKiller(killerFlagName).
				JsonString()
			mailProto.ReportTag = maildata.TagPve
			r.services.mail.SendReportMail(attackerId, mailProto, ctime)
		}

	})
	return
}

// 加速
func (r *Realm) SpeedUp(hc iface.HeroController, troopId, otherTroopId int64, speedUpRate float64, speedUpCostType uint64) (processed bool, err error) {
	processed = r.queueFunc(true, func() {
		base := r.getBase(hc.Id())
		if base == nil {
			err = realmerr.ErrSpeedUpTroopNotFound
			logrus.Debugf("部队加速行军，自己主城没找到")
			return
		}

		selfTroop := base.selfTroops[troopId]
		if selfTroop == nil {
			for _, t := range base.selfTroops {
				if t.AssemblyId() == troopId {
					selfTroop = t
					break
				}
			}

			if selfTroop == nil {
				err = realmerr.ErrSpeedUpTroopNotFound
				logrus.Debugf("部队加速行军，自己队伍没找到")
				return
			}
		}

		speedUpTroop := selfTroop
		if assembly := selfTroop.GetAssembly(); assembly != nil {
			if otherTroopId == 0 {
				if assembly.self.State() != realmface.Assembly {
					// 给主队加速
					speedUpTroop = assembly.self
				} else {
					if selfTroop.State() == realmface.AssemblyArrived {
						err = realmerr.ErrSpeedUpAssemblyWait
						logrus.Debugf("部队加速行军，集结等待中")
						return
					}
				}

			} else {
				otherTroop := assembly.self
				if otherTroop.Id() != otherTroopId {
					otherTroop = nil
					for _, m := range assembly.member {
						if m == nil {
							continue
						}

						if m.Id() == otherTroopId {
							otherTroop = m
							break
						}
					}

					if otherTroop == nil {
						err = realmerr.ErrSpeedUpOtherTroopNotFound
						logrus.Debugf("部队加速行军，有otherTroopId，没找到目标队伍")
						return
					}

					if otherTroop.State() == realmface.AssemblyArrived {
						err = realmerr.ErrSpeedUpAssemblyWait
						logrus.Debugf("部队加速行军，给别的集结队伍加速，目标已经到达集结地")
						return
					}
				} else {
					if otherTroop.State() == realmface.Assembly {
						err = realmerr.ErrSpeedUpAssemblyWait
						logrus.Debugf("部队加速行军，集结等待中，不能给主队加速")
						return
					}
				}

				speedUpTroop = otherTroop
			}
		}

		if !speedUpTroop.State().IsMoving() {
			err = realmerr.ErrSpeedUpTroopNoMoving
			logrus.WithField("state", speedUpTroop.State()).Debugf("部队加速行军，目标队伍不是行军中")
			return
		}

		oldArrivedTime := speedUpTroop.moveArriveTime
		ctime := r.services.timeService.CurrentTime()
		if !speedUpTroop.speedUp(r, ctime, speedUpRate) {
			err = realmerr.ErrSpeedUpTroopNoMoving
			logrus.Debugf("部队加速行军，已经不能再加速")
			return
		}

		r.broadcastMaToCared(speedUpTroop, addTroopTypeUpdate, 0)

		// lock hero
		r.heroBaseFuncWithSend(speedUpTroop.startingBase.Base(), func(hero *entity.Hero, result herolock.LockResult) {
			hero.UpdateTroop(speedUpTroop, false)

			result.Changed()
			result.Ok()

			r.services.tlogService.TlogSpeedUpFlow(hero, speedUpCostType, uint64(oldArrivedTime.Sub(speedUpTroop.moveArriveTime).Seconds()), uint64(oldArrivedTime.Sub(ctime).Seconds()), uint64(speedUpTroop.moveArriveTime.Sub(ctime).Seconds()), speedUpCostType, 0)
		})

		r.services.world.SendFunc(speedUpTroop.startingBase.Id(), func() pbutil.Buffer {
			return misc.NewS2cScreenShowWordsMsg(
				r.getTextHelp().RealmTroopSpeedUp.New().
					WithTroopIndex(entity.GetTroopIndex(speedUpTroop.Id()) + 1).
					JsonString())
		})

		if assembly := selfTroop.GetAssembly(); assembly != nil {
			assembly.broadcastChanged(r)
		}
	})
	return
}

func (r *Realm) GmSpeedUpFightMe(id int64) {

	r.queueFunc(true, func() {
		base := r.getBase(id)
		if base == nil {
			return
		}

		ctime := r.services.timeService.CurrentTime()
		for _, t := range base.targetingTroops {
			if t.speedUp(r, ctime, 1) {
				r.broadcastMaToCared(t, addTroopTypeUpdate, 0)
			}
		}
	})
	return
}

// 取消缓慢迁城
// 取消缓慢迁城，快速迁城，流亡，等等会自动取消缓慢迁移
func (r *Realm) CancelSlowMoveBase(hc iface.HeroController) (processed bool, err error) {
	processed = r.queueFunc(true, func() {

		base := r.getBase(hc.Id())
		if base == nil {
			err = realmerr.ErrCancelSlowMoveBaseSelfNoBase
			logrus.WithError(err).Debugf("取消缓慢迁城")
			return
		}

		heroHome := GetHomeBase(base)
		if heroHome == nil {
			err = realmerr.ErrCancelSlowMoveBaseTent
			logrus.WithError(err).Debugf("取消缓慢迁城")
			return
		}

		r.heroBaseFuncNotError(base.Base(), func(hero *entity.Hero) (heroChanged bool) {
			hero.UpdateHome(heroHome)
			heroChanged = true
			return
		})

		return
	})

	return
}

func (r *Realm) Mian(heroId int64, disappearTime time.Time, overwrite bool) (processed bool, err error) {

	processed = r.queueFunc(true, func() {
		ctime := r.services.timeService.CurrentTime()
		if !disappearTime.After(ctime) {
			logrus.Error("设置的免战时间，比当前时间要小")
			return
		}

		base := r.getBase(heroId)
		if base == nil {
			err = realmerr.ErrMianSelfNoBase
			logrus.WithError(err).Debugf("设置免战")
			return
		}

		if !base.isHeroHomeBase() {
			err = realmerr.ErrMianTent
			logrus.WithError(err).Debugf("设置免战")
			return
		}

		toSet := timeutil.Marshal32(disappearTime)
		if base.MianDisappearTime() != 0 {
			if overwrite {
				if toSet <= base.MianDisappearTime() {
					err = realmerr.ErrMianCantOverwrite
					logrus.WithError(err).Debugf("设置免战")
					return
				}
			} else {
				if timeutil.Marshal32(ctime) <= base.MianDisappearTime() {
					err = realmerr.ErrMianExist
					logrus.WithError(err).Debugf("设置免战")
					return
				}
			}
		}

		base.SetMianDisappearTime(toSet)

		// 更新到英雄身上
		r.heroBaseFuncNotError(base.Base(), func(hero *entity.Hero) (heroChanged bool) {
			hero.SetMianDisappearTime(ctime, disappearTime)
			return true
		})

		// 通知所有正在看这张地图的人
		r.broadcastBaseInfoToCared(base, addBaseTypeUpdate, 0)

		// 发更新消息给自己
		r.services.world.SendFunc(heroId, func() pbutil.Buffer {
			return region.NewS2cUpdateSelfMianDisappearTimeMsg(toSet, timeutil.Marshal32(ctime))
		})

		// 正在持续掠夺的马，回家
		for _, troop := range base.targetingTroops {
			if troop.State() == realmface.Robbing {
				r.trySendTroopDoneMail(troop, r.getTextHelp().MRDRMian4a.Text, r.getTextHelp().MRDRMian4d.Text, ctime)
				troop.backHome(r, ctime, true, true)
			}
		}
	})

	return
}

func (r *Realm) TryRemoveBaseMian(heroId int64) (processed bool) {

	if npcid.IsNpcId(heroId) {
		logrus.Error("移除主城免战，但是id是个npcid")
		return true
	}

	processed = r.queueFunc(false, func() {

		base := r.getBase(heroId)
		if base == nil {
			logrus.Debug("移除主城免战，base == nil")
			return
		}

		r.doRemoveBaseMian(base, true)

	})

	return
}

func (r *Realm) doRemoveBaseMian(base *baseWithData, shouldAddMianGoodsCd bool) {

	if base.MianDisappearTime() == 0 {
		return
	}
	base.SetMianDisappearTime(0)

	if base.BaseType() != realmface.BaseTypeHome {
		// 主城才有免战
		logrus.Error("不是主城也有免战？")
		return
	}

	heroId := base.Id()

	ctime := r.services.timeService.CurrentTime()

	// 更新到英雄身上
	r.heroBaseFuncWithSend(base.Base(), func(hero *entity.Hero, result herolock.LockResult) {

		if shouldAddMianGoodsCd {
			if ctime.Before(hero.GetMianDisappearTime()) {
				maxEndTime := ctime.Add(r.services.datas.RegionGenConfig().UseGoodsMianMaxDuraion)
				nextUseTime := timeutil.Min(hero.GetMianDisappearTime(), maxEndTime)
				hero.SetNextUseMianGoodsTime(nextUseTime)

				result.AddFunc(func() pbutil.Buffer {
					return region.NewS2cUseMianGoodsMsg(0, timeutil.Marshal32(nextUseTime))
				})
			}
		}

		hero.SetMianDisappearTime(time.Time{}, time.Time{})

		result.Ok()
	})

	// 通知所有正在看这张地图的人
	r.broadcastBaseInfoToCared(base, addBaseTypeUpdate, 0)

	// 发更新消息给自己
	r.services.world.SendFunc(heroId, func() pbutil.Buffer {
		return region.NewS2cUpdateSelfMianDisappearTimeMsg(0, 0)
	})
}

// 改变英雄基础信息（含帮派）
func (r *Realm) UpdateHeroBasicInfoNoBlock(heroId int64) {
	r.queueFuncNoBlock(func() {
		r.updateHeroBasicInfo(heroId)
	})

	return
}

func (r *Realm) updateHeroBasicInfo(heroId int64) {

	logrus.WithField("realmid", r.id).WithField("heroid", heroId).Debug("更新英雄基础信息")

	var err error

	base := r.getBase(heroId)
	if base == nil {
		err = realmerr.ErrChangeGuildSelfNoBase
		logrus.WithError(err).Debug("更新英雄基础信息, 本地区没有base（刚流亡?）")
		return
	}

	if !base.isHeroBaseType() {
		err = realmerr.ErrChangeGuildSelfNoBase
		logrus.WithError(err).Debug("更新英雄基础信息, 传入的heroId是个Npc?")
		return
	}

	originGuildId := base.GuildId()

	var isRemoveWhiteFlag bool
	r.services.heroDataService.FuncNotError(base.Id(), func(hero *entity.Hero) (heroChanged bool) {
		base.internalBase.UpdateHeroBasicInfo(hero)

		// 移除白旗
		if hero.GetWhiteFlagGuildId() != 0 {
			isRemoveWhiteFlag = hero.GetWhiteFlagGuildId() == hero.GuildId()
			if isRemoveWhiteFlag {
				hero.RemoveWhiteFlag()
			}
		}

		if homeBase := GetHomeBase(base); homeBase != nil {
			if guanfu := hero.Domestic().GetBuilding(shared_proto.BuildingType_GUAN_FU); guanfu != nil {
				homeBase.guanFuLevel = guanfu.Level
			}
			homeBase.outerCityUnlockBit = hero.Domestic().OuterCities().UnlockBit()
		}

		heroChanged = true
		return
	})

	base.updateRoBase()

	newGuildId := base.GuildId()
	if originGuildId != newGuildId {
		ctime := r.services.timeService.CurrentTime()
		for _, troop := range base.targetingTroops {

			if newGuildId == 0 {
				// 所有援助的马回家
				if !troop.state.IsAssistState() {
					continue
				}
			} else {
				if troop.startingBase.GuildId() == newGuildId {
					if !troop.state.IsInvateState() {
						continue
					}
					// 帮派相同的进攻马回家
				} else {
					if !troop.state.IsAssistState() {
						continue
					}
					// 帮派不同的援助马回家
				}
			}

			r.trySendTroopDoneMail(troop, r.getTextHelp().MRDRFinished4a.Text, r.getTextHelp().MRDRFinished4d.Text, ctime)
			troop.backHome(r, ctime, true, true)
		}

		for _, troop := range base.selfTroops {

			if newGuildId == 0 {
				// 所有援助的马回家
				if !troop.state.IsAssistState() {
					continue
				}
			} else {
				if troop.targetBase != nil && troop.targetBase.GuildId() == newGuildId {
					if !troop.state.IsInvateState() {
						continue
					}
					// 帮派相同的进攻马回家
				} else {
					if !troop.state.IsAssistState() {
						continue
					}
					// 帮派不同的援助马回家
				}
			}

			r.trySendTroopDoneMail(troop, r.getTextHelp().MRDRFinished4a.Text, r.getTextHelp().MRDRFinished4d.Text, ctime)
			troop.backHome(r, ctime, true, true)
		}

		// 如果新联盟跟白旗联盟一致，则移除白旗
		if newGuildId != 0 {
			r.removeHeroWhiteFlagIfSame(base, newGuildId, true)
		}

		base.remindAttackOrRobCountChanged(r)

		r.updateHomeResourcePointBlockWhenGuildChanged(base)
	}

	// 通知所有正在看这张地图的人
	r.broadcastBaseInfoToCared(base, addBaseTypeUpdate, 0)

	// 更新这个玩家控制的宝藏
	if home := GetHomeBase(base); home != nil {
		if hero := r.services.heroSnapshotService.Get(home.Id()); hero != nil {
			heroBytes := hero.EncodeBasic4ClientBytes()
			for baozId := range home.keepBaozMap {
				if baozBase := r.getBase(baozId); baozBase != nil {
					if baoz := GetBaoZangBase(baozBase); baoz != nil {
						if baoz.heroId == base.Id() {
							baoz.heroBytes = heroBytes

							baoz.ClearUpdateBaseInfoMsg()
							r.broadcastBaseInfoToCared(baozBase, addBaseTypeUpdate, 0)
						}
					}
				}
			}
		}
	}
}

func (r *Realm) onGuildChanged(guildId int64, name, flagName string, isRemoved bool) (processed bool) {

	processed = r.queueFunc(false, func() {
		if isRemoved {
			// 帮派删除了，所有跟这个帮派有关的白旗，拔掉(新建基地的时候，基地加入场景时候，也判断一次是否存在，不存在，拔旗)
			r.rangeBases(func(base *baseWithData) (toContinue bool) {
				r.removeHeroWhiteFlagIfSame(base, guildId, false)
				return true
			})

			// 删除联盟工坊
			baseId := npcid.NewGuildWorkshopId(guildId)
			base := r.getBase(baseId)
			if base == nil {
				return
			}

			ctime := r.services.timeService.CurrentTime()
			r.removeRealmBase(base, nil, nil, removeBaseTypeTransfer, ctime)
		} else {
			// 更新联盟数据
			r.rangeBases(func(base *baseWithData) (toContinue bool) {
				if base.GuildId() == guildId {
					r.broadcastBaseInfoToCared(base, addBaseTypeUpdate, 0)
				}
				return true
			})
		}
	})
	return
}

var noMoveBaseRestoreProsperityBufEndTimeMsg = region.NewS2cProsperityBufMsg(0).Static()

// 英雄操作导致繁荣度增加. 传入增加的量, 由这里执行具体增加的操作
// 升级建筑在英雄线程不要直接增加繁荣度, 要调这个方法来修改繁荣度.
func (r *Realm) AddProsperity(heroId int64, toAdd uint64) (processed bool, err error) {
	toAddAmount := int64(toAdd)
	if toAddAmount <= 0 {
		return true, nil
	}

	processed = r.queueFunc(false, func() {
		r.handleUpdateProsperity(heroId, toAddAmount)
	})

	return
}

func (r *Realm) UpdateProsperity(heroId int64) (processed bool) {
	processed = r.queueFunc(false, func() {
		r.handleUpdateProsperity(heroId, 0)
	})

	return
}

func (r *Realm) ReduceProsperity(heroId int64, toAdd uint64) (processed bool) {
	if toAdd <= 0 {
		return true
	}

	processed = r.queueFunc(false, func() {
		r.handleUpdateProsperity(heroId, -int64(toAdd))
	})

	return
}

func (r *Realm) GmReduceProsperity(heroId int64, toAdd uint64) (processed bool) {
	toAddAmount := int64(toAdd)
	if toAddAmount <= 0 {
		return true
	}

	processed = r.queueFunc(false, func() {
		r.handleUpdateProsperity(heroId, -toAddAmount)
	})

	return
}

func (r *Realm) handleUpdateProsperity(heroId, toChanged int64) (err error) {
	logrus.WithField("realmid", r.id).WithField("heroid", heroId).Debug("给英雄加繁荣度")

	base := r.getBase(heroId)
	if base == nil {
		err = realmerr.ErrAddProsperitySelfNoBase
		logrus.WithError(err).Debug("英雄加繁荣度, 本地区没有base（刚流亡?）")
		return
	}

	return r.handleUpdateBaseProsperity(base, toChanged)
}

func (r *Realm) handleUpdateBaseProsperity(base *baseWithData, toChanged int64) (err error) {

	originLevel := base.BaseLevel()

	var oldProsperity, newProsperity, prosperityCapcity uint64
	var toCleanAstDefendLog bool
	updateHomeInfo := false
	r.heroBaseFuncWithSend(base.Base(), func(hero *entity.Hero, result herolock.LockResult) {
		oldProsperity = hero.Prosperity()
		// 同步一次主城繁荣度上限
		if base.BaseType() == realmface.BaseTypeHome {
			base.SetProsperityCapcity(hero.ProsperityCapcity())

			if homeBase := GetHomeBase(base); homeBase != nil {
				if guanfu := hero.Domestic().GetBuilding(shared_proto.BuildingType_GUAN_FU); guanfu != nil {
					updateHomeInfo = updateHomeInfo || homeBase.guanFuLevel != guanfu.Level
					homeBase.guanFuLevel = guanfu.Level
				}
				updateHomeInfo = updateHomeInfo || homeBase.outerCityUnlockBit != hero.Domestic().OuterCities().UnlockBit()
				homeBase.outerCityUnlockBit = hero.Domestic().OuterCities().UnlockBit()
			}
		}

		if toChanged >= 0 {
			base.AddProsperity(uint64(toChanged))
		} else {
			base.ReduceProsperityDontKill(uint64(-toChanged))

			// 减少之后，主城可能降级
			newBaseLevel := calculateBaseLevelByProsperity(base.Prosperity(), r.services.datas.BaseLevelData().Must(originLevel))
			// 根据繁荣度计算等级，如果有变化，更新
			if originLevel != newBaseLevel {
				base.SetBaseLevel(newBaseLevel)

				updateHomeInfo = true

				if !npcid.IsNpcId(base.Id()) {
					gamelogs.UpdateBaseLevelLog(constants.PID, heroid.GetSid(base.Id()), base.Id(), base.BaseLevel())
				}
			}
		}

		prosperity := base.Prosperity()
		hero.UpdateBase(base.Base())

		switch base.BaseType() {
		case realmface.BaseTypeHome:
			prosperityCapcity := hero.ProsperityCapcity()
			result.AddFunc(func() pbutil.Buffer {
				return domestic.NewS2cHeroUpdateProsperityMsg(u64.Int32(prosperity), u64.Int32(prosperityCapcity))
			})

			if hero.Prosperity() >= prosperityCapcity && !timeutil.IsZero(hero.GetMoveBaseRestoreProsperityBufEndTime()) {
				hero.SetMoveBaseRestoreProsperityBufEndTime(time.Time{})
				result.Add(noMoveBaseRestoreProsperityBufEndTimeMsg)
			}

			if home, ok := base.Base().(*heroHome); ok {
				home.UpdateMoveBaseBuf(hero)

				// 恢复满，取消恢复状态
				if home.Prosperity() >= home.ProsperityCapcity() {
					home.isRestore = false
					hero.MiscData().SetIsRestoreProsperity(false)
					toCleanAstDefendLog = true
				}
			}

			if toChanged > 0 {
				hero.HistoryAmount().Increase(server_proto.HistoryAmountType_RestoreProsperity, uint64(toChanged))
				heromodule.UpdateTaskProgress(hero, result, shared_proto.TaskTargetType_TASK_TARGET_ACCUM_RECOVER_PROSPERITY)
			}
		}

		newProsperity = hero.Prosperity()
		prosperityCapcity = hero.ProtectedCapcity()
		result.Changed()
		result.Ok()
	})

	if toCleanAstDefendLog {
		r.ClearAstDefendLog(base.Base().Id()) // 清除援助驻扎加繁荣度日志
	}

	if originLevel != base.BaseLevel() {
		// 更新占用的资源点
		r.updateHomeResourcePointBlockWhenLevelChanged(base, originLevel)

		r.services.world.SendFunc(base.Id(), func() pbutil.Buffer {
			return region.NewS2cSelfUpdateBaseLevelMsg(u64.Int32(base.BaseLevel()))
		})
	}

	// 更新base繁荣度
	if updateHomeInfo {
		r.broadcastBaseInfoToCared(base, addBaseTypeUpdate, 0)
	} else {
		r.updateWatchBaseProsperity(base)
	}

	// 飘字，目前减少的时候飘，其他不飘
	if toChanged < 0 {
		r.broadcastShowWords(base, nil, uint64(-toChanged), 0, 0, 0, 0)
	}

	// 更新繁荣度buff
	if base.isHeroHomeBase() {
		r.UpdateProsperityBuff(base.Id(), oldProsperity, newProsperity, prosperityCapcity)
	}

	base.updateRoBase()

	return
}

// 手动升级老家等级
func (r *Realm) UpgradeBase(hc iface.HeroController) (processed bool, err error) {
	processed = r.queueFunc(true, func() {
		base := r.getBase(hc.Id())
		if base == nil || !base.isHeroHomeBase() {
			err = realmerr.ErrUpgradeBaseNotMyRealm
			logrus.WithError(err).Debug("升级主城失败")
			return
		}

		originLevel := base.BaseLevel()
		nextLevelData := r.GetBaseLevel(base.Prosperity())
		if nextLevelData == nil {
			err = realmerr.ErrUpgradeBaseAlreadyMax
			logrus.WithError(err).Debug("升级主城失败")
			return
		}

		if base.Prosperity() < nextLevelData.Prosperity {
			err = realmerr.ErrUpgradeBaseNotEnoughProsperity
			logrus.WithError(err).Debug("升级主城失败")
			return
		}

		ctime := r.services.timeService.CurrentTime()

		// 主城在我这, lock英雄, 再查看一次
		var npcOffsets []cb.Cube
		var heroId int64
		var baseLevel uint64
		var baseX, baseY int
		var toSendProtectRemoveMailFunc func()
		hc.FuncWithSend(func(hero *entity.Hero, result herolock.LockResult) {
			//if hero.BaseRegion() != r.id {
			//	err = realmerr.ErrUpgradeBaseNotMyRealm
			//	logrus.WithError(err).Error("升级主城在lock中失败......")
			//	return
			//}
			//
			//if hero.BaseLevel() != base.BaseLevel() || hero.Prosperity() != base.Prosperity() {
			//	err = realmerr.ErrUpgradeBaseNotEnoughProsperity
			//	logrus.WithError(err).Error("升级主城在lock中失败......hero的baseLevel或prosperity和realm中的不同")
			//	return
			//}
			oldMaxBaseLevel := hero.HomeHistoryMaxLevel() // 记录老的最高主城等级历史记录
			heroId = hero.Id()
			baseX = hero.BaseX()
			baseY = hero.BaseY()
			hero.SetBaseLevel(nextLevelData.Level)
			baseLevel = hero.BaseLevel()
			result.Add(nextLevelData.UpdateBaseLevelMsg)

			hctx := heromodule.NewContext(r.dep, operate_type.HeroUpgradeBaseLevel)
			if oldMaxBaseLevel < hero.HomeHistoryMaxLevel() { // 主城首次升到这个等级的时候
				// 系统广播
				if d := hctx.BroadcastHelp().BaseLevel; d != nil {
					hctx.AddBroadcast(d, hero, result, 0, hero.BaseLevel(), func() *i18n.Fields {
						text := d.NewTextFields()
						text.WithClickHeroFields(data.KeySelf, hctx.GetFlagHeroName(hero), hero.Id())
						text.WithFields(data.KeyNum, hero.BaseLevel())
						return text
					})
				}
				if giftData := r.dep.Datas().EventLimitGiftConfig().GetHomeBaseGift(hero.BaseLevel()); giftData != nil {
					heromodule.ActivateEventLimitGift(hero, result, giftData, ctime)
				}
			}

			// tlog
			hctx.Tlog().TlogCityExpFlow(hero, originLevel, baseLevel, hctx.OperId())

			heromodule.UpdateTaskProgress(hero, result, shared_proto.TaskTargetType_TASK_TARGET_BASE_LEVEL)

			base.SetProsperityCapcity(hero.ProsperityCapcity())

			if nextLevelData.Level >= r.config().NewHeroRemoveMianBaseLevel {
				if dt := hero.GetNewHeroMianDisappearTime().Unix(); dt > 0 {
					hero.SetNewHeroMianDisappearTime(time.Time{})
					result.Add(region.NewS2cUpdateNewHeroMianDisappearTimeMsg(0))

					hero.SetMianDisappearTime(time.Time{}, time.Time{})

					if int64(base.MianDisappearTime()) == dt {
						// 移除新手免战
						base.SetMianDisappearTime(0)

						toSendProtectRemoveMailFunc = heromodule.TrySendMailFunc(r.services.mail, hero, shared_proto.HeroBoolType_BOOL_PROTECT_REMOVED,
							r.services.datas.MailHelp().FirstProtectRemoved, ctime)
					}
				}
			}

			// 再看一下是不是有野怪要加到场景
			if home := GetHomeBase(base); home != nil {
				var homeNpcDatas []*basedata.HomeNpcBaseData
				for _, data := range r.services.datas.GetHomeNpcBaseDataArray() {
					if data.HomeBaseLevel > 0 && data.HomeBaseLevel == nextLevelData.Level {
						homeNpcDatas = append(homeNpcDatas, data)
					}
				}

				npcOffsets = r.addHomeNpc(hero, result, hc.GetViewArea(), base, home, homeNpcDatas)
			}

			heromodule.OnHeroEventWithSubType(hero, result, shared_proto.HeroEvent_HERO_EVENT_BASE_LEVEL, nextLevelData.Level)

			if hero.Bools().TrySet(shared_proto.HeroBoolType_BOOL_UPGRADE_BASE) {
				result.Add(misc.NewS2cSetHeroBoolMsg(int32(shared_proto.HeroBoolType_BOOL_UPGRADE_BASE)))
			}

			result.Changed()
			result.Ok()
			return
		})

		if toSendProtectRemoveMailFunc != nil {
			toSendProtectRemoveMailFunc()
		}

		// still in loop
		if err != nil {
			return
		}
		base.SetBaseLevel(nextLevelData.Level)
		r.broadcastBaseInfoToCared(base, addBaseTypeUpdate, 0)

		r.ruinsBasePosInfoMap.OnBaseChanged(base.Base())

		if originLevel != base.BaseLevel() {
			// 添加资源点占用
			r.updateHomeResourcePointBlockWhenLevelChanged(base, originLevel)
		}

		// 农场野怪冲突
		r.updateFarmWithNpc(heroId, baseLevel, baseX, baseY, npcOffsets, true, ctime)

		base.updateRoBase()

		if !npcid.IsNpcId(base.Id()) {
			gamelogs.UpdateBaseLevelLog(constants.PID, heroid.GetSid(base.Id()), base.Id(), base.BaseLevel())
		}

		return
	})

	return
}

func (r *Realm) AddHomeNpc(hc iface.HeroController, homeNpcDatas []*basedata.HomeNpcBaseData) (processed bool, err error) {

	if len(homeNpcDatas) <= 0 {
		return
	}

	processed = r.queueFunc(true, func() {
		base := r.getBase(hc.Id())
		if base == nil {
			err = realmerr.ErrAddHomeNpcNotHome
			logrus.WithField("reason", err).Debug("添加HomeNpc，base == nil")
			return
		}

		home := GetHomeBase(base)
		if home == nil {
			err = realmerr.ErrAddHomeNpcNotHome
			logrus.WithField("reason", err).Debug("添加HomeNpc，home == nil")
			return
		}

		ctime := r.services.timeService.CurrentTime()

		var npcOffsets []cb.Cube
		hc.FuncWithSend(func(hero *entity.Hero, result herolock.LockResult) {
			npcOffsets = r.addHomeNpc(hero, result, hc.GetViewArea(), base, home, homeNpcDatas)
		})

		// 农场野怪冲突
		r.updateFarmWithNpc(hc.Id(), base.BaseLevel(), base.BaseX(), base.BaseY(), npcOffsets, true, ctime)
	})

	return
}

func (r *Realm) addHomeNpc(hero *entity.Hero, result herolock.LockResult, heroViewArea *realmface.ViewArea,
	base *baseWithData, home *heroHome, homeNpcDatas []*basedata.HomeNpcBaseData) (npcOffsets []cb.Cube) {

	for _, data := range homeNpcDatas {
		if !hero.HasCreateHomeNpcBase(data.Id) {
			//id := r.services.homeNpcBaseIdGen.Next(data.Data.Id)
			heroHomeNpcBase := hero.CreateHomeNpcBase(data)
			homeNpcBase := r.newHomeNpcBase(hero.Id(), heroHomeNpcBase, base.BaseX(), base.BaseY())
			home.homeNpcBase[homeNpcBase.Id()] = homeNpcBase

			if base.CanSee(heroViewArea) {
				result.Add(homeNpcBase.NewUpdateBaseInfoMsg(r, addBaseTypeCanSee))
			}

			npcOffsets = append(npcOffsets, cb.XYCube(heroHomeNpcBase.GetData().EvenOffsetX, heroHomeNpcBase.GetData().EvenOffsetY))
		}
	}

	return
}

func (r *Realm) StartCareMilitary(hc iface.HeroController) (processed bool) {
	processed = r.queueFunc(true, func() {
		if hc.IsClosed() {
			return
		}

		if cond := hc.GetCareCondition(); cond != nil {
			// 所有的马都告诉你吧
			r.rangeBases(func(base *baseWithData) (toContinue bool) {
				for _, troop := range base.selfTroops {
					if !troop.ownerCanSeeTarget || troop.startingBase.Id() == hc.Id() {

						if troop.isMatchCondition(r, cond) {
							hc.Send(troop.getUpdateMsgMa(r))

							// 现在看得见的，设置careMe
							troop.setCareMeHeroId(hc.Id())
						} else {
							// 现在看不见，之前看得见的，发送移除消息
							if troop.tryRemoveCareMeHeroId(hc.Id()) {
								hc.Send(troop.getRemoveMsgMa())
							}
						}
					}
				}
				return true
			})
		}

	})
	return
}

func (r *Realm) newViewArea(posX, posY, lenX, lenY int) *realmface.ViewArea {
	radiusX := imath.Min(imath.Max(lenX, r.config().MinViewXLen), r.config().MaxViewXLen) / 2
	radiusY := imath.Min(imath.Max(lenY, r.config().MinViewYLen), r.config().MaxViewYLen) / 2

	minX := posX - radiusX
	minY := posY - radiusY
	maxX := posX + radiusX
	maxY := posY + radiusY

	return &realmface.ViewArea{
		CenterX:    posX,
		CenterY:    posY,
		MinX:       minX,
		MinY:       minY,
		MaxX:       maxX,
		MaxY:       maxY,
		UpdateTime: r.services.timeService.CurrentTime(),
	}
}

func (r *Realm) QueryTroopUnit(hc iface.HeroController, heroId, troopId int64) (processed bool, err error) {

	processed = r.queueFunc(true, func() {

		base := r.getBase(heroId)
		if base == nil {
			logrus.Debug("查询军情，找不到英雄")
			err = realmerr.ErrGetMilitaryBaseNotFound
			return
		}

		t := base.selfTroops[troopId]
		if t == nil {
			logrus.Debug("查询军情，找不到部队")
			err = realmerr.ErrGetMilitaryTroopNotFound
			return
		}

		hc.Send(region.NewS2cRequestTroopUnitMsg(t.getProtoBytes(r)))
	})

	return
}

// 变更个人签名
func (r *Realm) ChangeSign(heroId int64, sign string) {
	r.queueFunc(false, func() {
		base := r.getBase(heroId)
		if base == nil {
			logrus.Debug("变更签名，找不到英雄")
			return
		}

		if base.BaseType() != realmface.BaseTypeHome {
			logrus.Debug("变更签名，不是主城")
			return
		}

		home := GetHomeBase(base)
		if home == nil {
			return
		}

		home.SetSign(sign)

		// 通知所有正在看这张地图的人
		r.broadcastBaseInfoToCared(base, addBaseTypeUpdate, 0)
	})

	return
}

// 变更个人签名
func (r *Realm) ChangeTitle(heroId int64, title uint64) {
	r.queueFunc(false, func() {
		base := r.getBase(heroId)
		if base == nil {
			logrus.Debug("变更称号，找不到英雄")
			return
		}

		if base.BaseType() != realmface.BaseTypeHome {
			logrus.Debug("变更称号，不是主城")
			return
		}

		home := GetHomeBase(base)
		if home == nil {
			return
		}

		home.SetTitle(title)

		// 通知所有正在看这张地图的人
		r.broadcastBaseInfoToCared(base, addBaseTypeUpdate, 0)
	})

	return
}

// 开始关心这个地图, 获得地图中所有主城的信息
func (r *Realm) StartCareRealm(hc iface.HeroController, posX, posY, lenX, lenY int) (processed bool) {
	toSet := r.newViewArea(posX, posY, lenX, lenY)
	return r.updateViewArea(hc, toSet)
}

func (r *Realm) StopCareRealm(hc iface.HeroController) (processed bool) {
	return r.updateViewArea(hc, nil)
}

func (r *Realm) updateViewArea(hc iface.HeroController, toSet *realmface.ViewArea) (processed bool) {

	if hc.IsClosed() {
		return
	}

	oldArea := hc.GetViewArea()

	// if oldArea != nil && toSet != nil {
	// 	//判定有限的距离内不予刷新
	// 	if oldArea.GetCenterDistant(toSet) < 64 {
	// 		return
	// 	}
	// }

	hc.SetViewArea(toSet)

	index := 0

	watchObjList := hc.GetWatchObjList()
	
	if toSet != nil {
		//kimi 添加英雄AOI索引
		r.blockManager.AddHeroIndex(toSet.CenterX,toSet.CenterY,hc)

		newWatchObjList := make(map[interface{}]int)
		//找到所有附近9宫格的对象,并发送在屏幕范围内的对象
		for i := -constants.RealmIndexBlockSize; i <= constants.RealmIndexBlockSize; i += constants.RealmIndexBlockSize {
			for j := -constants.RealmIndexBlockSize; j <= constants.RealmIndexBlockSize; j += constants.RealmIndexBlockSize {
				x, y := toSet.CenterX+i, toSet.CenterY+j
				r.blockManager.rangeBase(x, y, func(base *baseWithData) (toContinue bool) {

					// oldCanSee := base.CanSee(oldArea)
					newCanSee := base.CanSee(toSet)

					if newCanSee {
						if _, ok := watchObjList[base];ok {
							delete(watchObjList, base)
						} else {
							base.AddWatcher(hc)
							//出现,添加消息
							hc.Send(base.NewUpdateBaseInfoMsg(r, addBaseTypeCanSee))
						}
						newWatchObjList[base] = 1
					}

					index++

					// if oldCanSee != newCanSee {
					// 	if oldCanSee {
					// 		//  消失，移除消息
					// 		if toSet != nil {
					// 			hc.Send(base.GetCantSeeMeMsg())
					// 		}
					// 	} else {
					// 		// 出现，添加消息
					// 		hc.Send(base.NewUpdateBaseInfoMsg(r, addBaseTypeCanSee))
					// 	}
					// }

					// for _, troop := range base.selfTroops {
					// 	if !troop.ownerCanSeeTarget || troop.startingBase.Id() == hc.Id() {
					// 		oldCanSee := troop.CanSee(oldArea)
					// 		newCanSee := troop.CanSee(toSet)

					// 		if oldCanSee != newCanSee {
					// 			if oldCanSee {
					// 				//  消失，移除消息
					// 				if toSet != nil {
					// 					hc.Send(troop.GetCantSeeMeMsg())
					// 				}
					// 			} else {
					// 				// 出现，添加消息
					// 				hc.Send(troop.getAddSeeMeMsg(r))
					// 			}
					// 		}
					// 	}
					// 	index++
					// }
					return true
				})

				r.blockManager.rangeTroop(x,y, func(troop *troop) (toContinue bool) {
					if !troop.ownerCanSeeTarget || troop.startingBase.Id() == hc.Id() {
						// oldCanSee := troop.CanSee(oldArea)
						newCanSee := troop.CanSee(toSet)
						
						if newCanSee {
							if _, ok := watchObjList[troop];ok {
								delete(watchObjList, troop)
							} else {
								troop.AddWatcher(hc)
								//出现,添加消息
								hc.Send(troop.getAddSeeMeMsg(r))
							}
							newWatchObjList[troop] = 1
						}
						// if oldCanSee != newCanSee {
						// 	if oldCanSee {
						// 		//  消失，移除消息
						// 		if toSet != nil {
						// 			hc.Send(troop.GetCantSeeMeMsg())
						// 		}
						// 	} else {
						// 		// 出现，添加消息
						// 		hc.Send(troop.getAddSeeMeMsg(r))
						// 	}
						// }
					}
					index++
					return true
				})

				r.blockManager.rangeRuin(x, y, func(info *ruinsBase) (toContinue bool) {
					if info.proto == nil {
						return true
					}
					// oldCanSee := info.CanSee(oldArea)
					newCanSee := info.CanSee(toSet)

					if newCanSee {
						if _, ok := watchObjList[info];ok {
							delete(watchObjList, info)
						} else {
							info.AddWatcher(hc)
							//出现,添加消息
							hc.Send(info.addMsg)
						}
						newWatchObjList[info] = 1
					}

					// if oldCanSee != newCanSee {
					// 	if oldCanSee {
					// 		//  消失，移除消息
					// 		if toSet != nil {
					// 			hc.Send(info.removeMsg)
					// 		}
					// 	} else {
					// 		// 出现，添加消息
					// 		hc.Send(info.addMsg)
					// 	}
					// }
					index++
					return true
				})
			}
		}
		
		hc.SetWatchObjList(newWatchObjList)

		//移除老的 未看到的对象
		for k := range watchObjList {
			obj,ok := k.(view_object)
			if ok {
				hc.Send(obj.GetCantSeeMeMsg())
			}
		}

		if r.npcBaseMsg != nil {
			hc.Send(r.npcBaseMsg)
		}
	} else {
		//kimi 移除英雄AOI索引
		r.blockManager.RemoveHeroIndex(hc)

		//移除老的 未看到的对象
		for k := range watchObjList {
			obj,ok := k.(view_object)
			if ok {
				obj.RemoveWatcher(hc)
			}
		}
		hc.SetWatchObjList(nil)
	}

	// 主城Npc
	if selfBase := r.getBase(hc.Id()); selfBase != nil {
		if home := GetHomeBase(selfBase); home != nil {
			for _, base := range home.homeNpcBase {
				oldCanSee := base.CanSee(oldArea)
				newCanSee := base.CanSee(toSet)
				if oldCanSee != newCanSee {
					if oldCanSee {
						//  消失，移除消息
						if toSet != nil {
							hc.Send(base.GetCantSeeMeMsg())
						}
					} else {
						// 出现，添加消息
						hc.Send(base.NewUpdateBaseInfoMsg(r, addBaseTypeCanSee))
					}
				}
			}
		}
	}

	
	logrus.Debugf("一次遍历次数为[%d]", index)
	return
}

func (r *Realm) broadcastBaseInfoToCared(base *baseWithData, addType int32, except int64) {
	broadcastMsg := base.NewUpdateBaseInfoMsg(r, addType)
	r.broadcastBaseToCared(base, broadcastMsg, except)

	//TODO
	// 这里 broadcastMsg 有两种可能，一种是缓存消息，一种是非缓存消息，所以不可以调用 DoFreeEvenItsStaticAndFuckMeIfItExplodes，后面再来优化
	//broadcastMsg.DoFreeEvenItsStaticAndFuckMeIfItExplodes()
}

func (r *Realm) broadcastBaseToCared(base *baseWithData, msg pbutil.Buffer, except int64) {

	r.broadcastToCared(base, msg, except)
}

type CanSeeFunc func(area *realmface.ViewArea) bool

func (f CanSeeFunc) CanSee(area *realmface.ViewArea) bool {
	return f(area)
}
/*** 
代码纯粹是为了实现接口方法, 暂时看起来没有其他方法被调用     
***/
func (f CanSeeFunc)GetPos()[][2]int {
	return nil
}

func (f CanSeeFunc) GetCantSeeMeMsg() pbutil.Buffer {
	return nil
}

func (f CanSeeFunc) AddWatcher(hc iface.HeroController) {

}


func (f CanSeeFunc) RemoveWatcher(hc iface.HeroController) {

}

func (f CanSeeFunc) RangeWatchList(func(hc iface.HeroController)) {

}



type view_object interface {
	CanSee(area *realmface.ViewArea) bool
	GetPos() [][2]int

	//统一接口 发送移除消息
	GetCantSeeMeMsg() pbutil.Buffer
	//kimi 添加观察者(这里是为了玩家在野外视野移动中,给观察对象加一个观察者索引)
	AddWatcher(hc iface.HeroController)
	//移除观察者
	RemoveWatcher(hc iface.HeroController)
	//遍历观察者
	RangeWatchList(func(hc iface.HeroController))
}

func (r *Realm) sendIfCared(heroId int64, obj view_object, msg pbutil.Buffer) {
	if heroId == 0 {
		return
	}

	r.services.world.FuncHero(heroId, func(id int64, hc iface.HeroController) {
		if obj.CanSee(hc.GetViewArea()) {
			hc.Send(msg)
		}
	})
}
//遍历兴趣英雄
func (r *Realm) rangeCareHeros(centerX, centerY int ,f func(hc iface.HeroController)) {
	for i := -constants.RealmIndexBlockSize; i <= constants.RealmIndexBlockSize; i += constants.RealmIndexBlockSize {
		for j := -constants.RealmIndexBlockSize; j <= constants.RealmIndexBlockSize; j += constants.RealmIndexBlockSize {
			x, y := centerX+i, centerY+j
			r.blockManager.rangeHero(x, y, func(hc iface.HeroController) (toContinue bool) {
				f(hc)
				logrus.Debugf("遍历兴趣英雄[%d]", hc.Id())
				return true
			})
		}
	}
}

// 广播给关心这个地图的人
func (r *Realm) broadcastToCared(obj view_object, msg pbutil.Buffer, except int64) {
	poses := obj.GetPos()
	if except == 0 {
		if poses == nil {
			r.services.world.WalkHero(func(id int64, hc iface.HeroController) {
				if obj.CanSee(hc.GetViewArea()) {
					hc.Send(msg)
				}
			})
		} else {
			// for _,v := range poses {
			// 	r.rangeCareHeros(v[0],v[1], func(hc iface.HeroController){
			// 		if obj.CanSee(hc.GetViewArea()) {
			// 			hc.Send(msg)
			// 		}
			// 	})
			// }
			obj.RangeWatchList(func(hc iface.HeroController){
				hc.Send(msg)
			})
		}
		// r.services.world.WalkHero(func(id int64, hc iface.HeroController) {
		// 	if obj.CanSee(hc.GetViewArea()) {
		// 		hc.Send(msg)
		// 	}
		// })
	} else {
		if poses == nil {
			r.services.world.WalkHero(func(id int64, hc iface.HeroController) {
				if id != except && obj.CanSee(hc.GetViewArea()) {
					hc.Send(msg)
				}
			})
		} else {
			// for _,v := range poses {
			// 	r.rangeCareHeros(v[0],v[1], func(hc iface.HeroController){
			// 		if hc.Id() != except && obj.CanSee(hc.GetViewArea()) {
			// 			hc.Send(msg)
			// 		}
			// 	})
			// }
			obj.RangeWatchList(func(hc iface.HeroController){
				hc.Send(msg)
			})
		}
	}
}

// 广播给关心这个地图的人
func (r *Realm) broadcastToCaredPos(posX, posY int, msg pbutil.Buffer, except int64) {
	if except == 0 {
		r.rangeCareHeros(posX,posY, func(hc iface.HeroController){
			if area := hc.GetViewArea(); area != nil && area.CanSeePos(posX, posY) {
				hc.Send(msg)
			}
		})
		// r.services.world.WalkHero(func(id int64, hc iface.HeroController) {
		// 	if area := hc.GetViewArea(); area != nil && area.CanSeePos(posX, posY) {
		// 		hc.Send(msg)
		// 	}
		// })
	} else {
		r.rangeCareHeros(posX,posY, func(hc iface.HeroController){
			if hc.Id() != except {
				if area := hc.GetViewArea(); area != nil && area.CanSeePos(posX, posY) {
					hc.Send(msg)
				}
			}
		})
		// r.services.world.WalkHero(func(id int64, hc iface.HeroController) {
		// 	if id != except {
		// 		if area := hc.GetViewArea(); area != nil && area.CanSeePos(posX, posY) {
		// 			hc.Send(msg)
		// 		}
		// 	}
		// })
	}
}

func (r *Realm) broadcastMaToCared(t *troop, addType int32, except int64) {

	f := func(id int64, hc iface.HeroController) {
		//  if t.CanSee(hc.GetViewArea()) {
		// 	hc.Send(t.newAddSeeMeMsg(r, addType))
		//  }
		 hc.Send(t.newAddSeeMeMsg(r, addType))

		if cond := hc.GetCareCondition(); cond != nil {
			careMillitary := t.isMatchCondition(r, cond)

			if careMillitary {
				hc.Send(t.getUpdateMsgMa(r))

				// 现在看得见的，设置careMe
				t.setCareMeHeroId(id)
			} else {

				// 现在看不见，之前看得见的，发送移除消息
				if t.tryRemoveCareMeHeroId(id) {
					hc.Send(t.getRemoveMsgMa())
				}
			}
		}
	}

	if t.ownerCanSeeTarget {
		r.services.world.FuncHero(t.startingBase.Id(), f)
	} else {
		poses := t.GetPos()
		if poses == nil {
			r.services.world.WalkHero(f)
		} else {
			t.RangeWatchList(func(hc iface.HeroController){
				f(hc.Id(),hc)
			})
			// for _,v := range poses {
			// 	r.rangeCareHeros(v[0],v[1], func(hc iface.HeroController){
			// 		f(hc.Id(),hc)
			// 	})
			// }
		}
	}

	r.services.world.SendFunc(t.startingBase.Id(), func() pbutil.Buffer {
		return t.getUpdateMsgToSelf(r)
	})

	// 处理集结的军情
	if t.assembly != nil && t.assembly.self == t && t.State() != realmface.Assembly {
		protoBytes := t.getProtoBytes(r)
		for _, t := range t.assembly.member {
			if t == nil {
				continue
			}
			r.services.world.SendFunc(t.startingBase.Id(), func() pbutil.Buffer {
				return region.NewS2cUpdateSelfMilitaryInfoMsg(u64.Int32(entity.GetTroopIndex(t.Id())+1), t.IdBytes(), protoBytes)
			})
		}
	}

}

func (r *Realm) broadcastRemoveSeeMeMsg(t *troop) {

	f := func(id int64, hc iface.HeroController) {

		if t.CanSee(hc.GetViewArea()) {
			hc.Send(t.newRemoveSeeMeMsg(removeTroopTypeDestroy))
		}

		if cond := hc.GetCareCondition(); cond != nil {
			if t.isMatchCondition(r, cond) {
				hc.Send(t.getRemoveMsgMa())
			}
		}
	}

	if t.ownerCanSeeTarget {
		r.services.world.FuncHero(t.startingBase.Id(), f)
	} else {
		// r.services.world.WalkHero(f)
		poses := t.GetPos()
		if poses == nil {
			r.services.world.WalkHero(f)
		} else {
			t.RangeWatchList(func(hc iface.HeroController){
				f(hc.Id(),hc)
			})
			// for _,v := range poses {
			// 	r.rangeCareHeros(v[0],v[1], func(hc iface.HeroController){
			// 		f(hc.Id(),hc)
			// 	})
			// }
		}
	}
}

func (r *Realm) broadcastRemoveMaToCared(t *troop) {

	r.broadcastRemoveSeeMeMsg(t)

	// 部队被干掉了
	r.services.world.SendFunc(t.startingBase.Id(), func() pbutil.Buffer {
		return t.getRemoveMsgToSelf()
	})
}

// 插白旗
func (r *Realm) trySetWhiteFlag(t *troop) (success bool) {

	if t.state != realmface.Robbing || t.startingBase.GuildId() == 0 {
		return
	}

	g := r.services.guildService.GetSnapshot(t.startingBase.GuildId())
	if g == nil {
		return
	}

	// 持续掠夺，给target插白旗
	if t.targetBase == nil {
		logrus.Error("尝试设置白旗，t.targetBase == nil")
		return
	}

	home := GetHomeBase(t.targetBase)
	if home == nil {
		return
	}

	ctime := r.services.timeService.CurrentTime()
	disappearTime := ctime.Add(r.services.datas.RegionConfig().WhiteFlagDuration)

	if r.heroBaseFuncNotError(t.targetBase.Base(), func(hero *entity.Hero) (heroChanged bool) {
		hero.SetWhiteFlag(t.startingBase.Id(), t.startingBase.GuildId(), disappearTime)
		return true
	}) {
		logrus.Error("尝试设置白旗，lock hero 失败")
		return
	}

	home.SetWhiteFlag(g.Id, timeutil.Marshal32(disappearTime))

	// 设置成功，发消息更新
	broadcastMsg := region.NewS2cUpdateWhiteFlagMsg(t.targetBase.IdBytes(), i64.Int32(g.Id), g.FlagName, home.WhiteFlagDisappearTime()).Static()
	r.services.world.Send(t.targetBase.Id(), broadcastMsg) // 玩家自己一定会收到
	r.broadcastBaseToCared(t.targetBase, broadcastMsg, t.targetBase.Id())

	return true
}

func (r *Realm) removeHeroWhiteFlagIfSameOnOtherRealm(heroId, targetGuildId int64) {
	r.services.otherRealmEventQueue.TryFunc(func() {
		r.queueFunc(false, func() {
			base := r.getBase(heroId)
			if base == nil {
				return
			}

			r.removeHeroWhiteFlagIfSame(base, targetGuildId, false)
		})
	})
}

func (r *Realm) removeHeroWhiteFlagIfSame(base *baseWithData, targetGuildId int64, heroRemoved bool) {

	home := GetHomeBase(base)
	if home == nil {
		return
	}

	if targetGuildId == 0 || home.WhiteFlagGuildId() != targetGuildId {
		return
	}

	heroId := base.Id()

	if !heroRemoved {
		if r.services.heroDataService.FuncNotError(heroId, func(hero *entity.Hero) (heroChanged bool) {
			hero.RemoveWhiteFlag()
			return true
		}) {
			logrus.Error("移除英雄白旗，lock hero 失败")
			return
		}
	}

	home.SetWhiteFlag(0, 0)

	// 设置成功，发消息更新
	broadcastMsg := region.NewS2cUpdateWhiteFlagMsg(base.IdBytes(), 0, "", 0).Static()
	r.services.world.Send(heroId, broadcastMsg) // 玩家自己一定会收到
	r.broadcastBaseToCared(base, broadcastMsg, heroId)
}

// 资源点
func (r *Realm) isResourcePointConflicted(heroId int64, cube cb.Cube) bool {
	// 联盟场景永不冲突
	if r.levelData.RegionType == shared_proto.RegionType_GUILD {
		return false
	}

	ids := r.resourceConflictHeroMap[cube]

	// 超过(>=)2人占着这个位置，并且包含自己，属于冲突，其他的都不属于冲突
	if len(ids) < 2 || !i64.Contains(ids, heroId) {
		return false
	}

	base := r.getBase(heroId)
	if base == nil {
		// 没在这个场景了，我去
		return false
	}

	if base.GuildId() == 0 {
		// 没有联盟
		return true
	}

	for _, conflictId := range ids {
		if conflictId == heroId {
			// 自己
			continue
		}

		conflictBase := r.getBase(conflictId)
		if conflictBase == nil {
			// 目标没在这个场景了
			continue
		}

		if conflictBase.GuildId() == base.GuildId() {
			// 同一个联盟
			continue
		}

		return true
	}

	return false
}

// 取一个冲突玩家ID
func (r *Realm) getOneConflictedHeroId(heroId int64, cube cb.Cube) int64 {
	ids := r.resourceConflictHeroMap[cube]
	base := r.getBase(heroId)
	if base == nil {
		// 没在这个场景了，我去
		return 0
	}

	for _, conflictId := range ids {
		if conflictId == heroId {
			// 自己
			continue
		}

		conflictBase := r.getBase(conflictId)
		if conflictBase == nil {
			// 目标没在这个场景了
			continue
		}

		if base.GuildId() > 0 && conflictBase.GuildId() == base.GuildId() {
			// 同一个联盟
			continue
		}

		return conflictId
	}
	return 0
}

func (r *Realm) queueFuncNoBlock(f func()) {

	e := &action{f: f, called: make(chan struct{})}

	select {
	case r.funcChan <- e:
		// 放进去了
		return
	default:
	}

	// 没进去，只能开个gorutinue再来一遍
	go call.CatchPanic(func() {
		if !r.queueFuncAction(false, e) {
			logrus.Error("Realm.queueFuncNoBlock 进入队列失败")
		}
	}, "Realm.queueNoBlock")
}

func (r *Realm) queueFunc(waitResult bool, f func()) (funcCalled bool) {
	if r.directExecFunc {
		call.CatchPanic(f, "Realm.directExecFunc")
		return true
	}

	e := &action{f: f, called: make(chan struct{})}

	return r.queueFuncAction(waitResult, e)
}

func (r *Realm) queueFuncAction(waitResult bool, e *action) (funcCalled bool) {

	select {
	case r.funcChan <- e:
		if waitResult {
			select {
			case <-r.loopExitNotify:
				return false // main loop exit

			case <-e.called:
				return true
			}
		} else {
			return true // put success
		}

	case <-queueTimeoutWheel.After(queueTimeout):
		return false

	case <-r.closeNotify:
		return false
	}
}

func (r *Realm) loop() {
	defer close(r.loopExitNotify)

	// 恢复繁荣度ticker
	secondTicker := time.NewTicker(time.Second)
	r.dailyTicker = r.services.tickService.GetDailyTickTime()
	r.per10MinuteTicker = r.services.tickService.GetPer10MinuteTickTime()

	// 尝试刷新怪物
	call.CatchPanic(func() {
		r.refreshJunTuanNpc(r.per10MinuteTicker.GetPrevTickTime())
	}, "Realm.refreshJunTuanNpc()")

	call.CatchPanic(func() {
		r.refreshBaoZangNpc(r.per10MinuteTicker.GetPrevTickTime())
	}, "Realm.refreshBaoZangNpc()")

	for {
		if r.handleFunc(secondTicker) {
			return
		}
	}
}

func (r *Realm) handleFunc(secondTicker *time.Ticker) (stop bool) {
	defer func() {
		if r := recover(); r != nil {
			// 严重错误. 英雄线程这里不能panic
			logrus.WithField("err", r).WithField("stack", string(debug.Stack())).Error("Realm.handle recovered from panic!!! SERIOUS PROBLEM")
			metrics.IncPanic()
		}
	}()

	if r.isClosed {
		// 在线程内关闭了场景，r.closeNotify一定是被close
		return true
	}

	select {
	case f := <-r.funcChan:
		call.CatchPanic(f.f, "Realm.handleFunc")
		close(f.called) // notify caller

	case <-r.nextEventTime:

		ctime := r.services.timeService.CurrentTime()

		lastEvent := r.eventQueue.Peek()

		for {
			if lastEvent == nil {
				r.nextEventTime = nil
				break
			}

			t := lastEvent.Time()
			if diff := t.Sub(ctime); diff > 0 {
				r.nextEventTime = time.After(diff)
				break
			}

			call.CatchPanic(func() {
				// 加保护
				lastEvent.Data().(func(*Realm))(r)
			}, "Realm.handleEvent")

			nextEvent := r.eventQueue.Peek()
			if nextEvent == lastEvent && !nextEvent.Time().After(t) {
				logrus.WithField("event", nextEvent).Error("realm.Event的func中, 既没有删除event, 也没有把event设到个更晚的时间")
				nextEvent.RemoveFromQueue()

				lastEvent = r.eventQueue.Peek()
			} else {
				lastEvent = nextEvent
			}
		}

	case <-secondTicker.C:
		// TODO 优化
		ctime := r.services.timeService.CurrentTime()

		r.updateRestoreOrLossProsperity()
		r.ruinsBasePosInfoMap.Update(ctime)

		// 更新联盟工坊
		r.tickGuildWorkshop(ctime)

		// 移除过期的召唤殷墟
		r.tryRemoveExpiredHeroBaozBase(ctime)

	case <-r.per10MinuteTicker.Tick():
		r.per10MinuteTicker = r.services.tickService.GetPer10MinuteTickTime()

		call.CatchPanic(func() {
			r.refreshJunTuanNpc(r.per10MinuteTicker.GetPrevTickTime())
		}, "Realm.refreshJunTuanNpc()")

		call.CatchPanic(func() {
			r.refreshBaoZangNpc(r.per10MinuteTicker.GetPrevTickTime())
		}, "Realm.refreshBaoZangNpc()")

	case <-r.dailyTicker.Tick():
		r.dailyTicker = r.services.tickService.GetDailyTickTime()
		call.CatchPanic(r.updatePerDay, "Realm.updatePerDay()")

	case <-r.closeNotify:
		r.isClosed = true
		return true // quit loop
	}

	return false
}

func (r *Realm) tryRemoveExpiredHeroBaozBase(ctime time.Time) {

	ctimeUnix := ctime.Unix()

	if ctimeUnix <= r.baseManager.heroBaozNextExpireTime {
		// 下次刷新时间未到
		return
	}

	var nextExpireTime int64 = math.MaxInt64
	r.baseManager.rangeHeroBaozIds(func(heroId, baozBaseId int64) (toContinue bool) {
		toContinue = true

		base := r.getBase(baozBaseId)
		if base == nil {
			r.baseManager.removeHeroBaozId(heroId)
			return
		}

		baoz := GetBaoZangBase(base)
		if baoz == nil {
			logrus.WithField("id", base.Id()).WithField("name", base.internalBase.HeroName()).Error("更新玩家召唤的殷墟，发现不是个殷墟对象")
			r.baseManager.removeHeroBaozId(heroId)
			return
		}

		if ctimeUnix < int64(baoz.heroEndTime) {
			// 还未过期
			nextExpireTime = i64.Min(nextExpireTime, int64(baoz.heroEndTime))
			return
		}

		// 已过期，移除这个宝藏
		r.removeRealmBaseNoReason(base, removeBaseTypeTransfer, ctime)
		r.baseManager.removeHeroBaozId(heroId)

		return
	})

	r.baseManager.heroBaozNextExpireTime = nextExpireTime
}

func (r *Realm) updateRestoreOrLossProsperity() {

	ctime := r.services.timeService.CurrentTime()

	season := r.services.seasonService.SeasonByTime(ctime)

	// 遍历所有的玩家主城行营
	r.rangeBases(func(b *baseWithData) (toContinue bool) {
		switch b.BaseType() {
		case realmface.BaseTypeHome:
			// 主城，定时恢复繁荣度

			home := GetHomeBase(b)
			if home == nil {
				break
			}

			// 繁荣度满了，不恢复
			if b.Prosperity() >= b.ProsperityCapcity() {
				if b.Base() != nil {
					r.ClearAstDefendLog(b.Base().Id())
				}
				break
			}

			// 是否恢复阶段
			if !home.isRestore {
				// 检查是否处于恢复阶段（在线表示恢复阶段）
				if !r.services.world.FuncHero(home.Id(), func(id int64, hc iface.HeroController) {
					if hc.FuncNotError(func(hero *entity.Hero) (heroChanged bool) {
						home.isRestore = true
						hero.MiscData().SetIsRestoreProsperity(true)
						return
					}) {
						return
					}
				}) {
					return
				}

			}

			if !home.tryRestoreProsperity(ctime, r.services.datas.RegionConfig().RestoreHomeProsperityDuration) {
				// 更新时间未到
				break
			}

			restoreAmount := r.services.datas.RegionConfig().RestoreHomeProsperity

			levelData := home.GetGuanFuLevelData(r)
			if levelData != nil {
				if ctime.Before(home.moveBaseRestoreProsperityBufEndTime) {
					// 迁城buff
					restoreAmount = levelData.MoveBaseRestoreHomeProsperity
				} else {
					// 正常恢复值
					restoreAmount = levelData.RestoreProsperity
				}
			}

			// 盟友驻扎加繁荣
			astToAdd := r.services.datas.RegionConfig().AstDefendRestoreHomeProsperityAmount.Calc(b.ProsperityCapcity())
			restoreAmount += astToAdd * b.getAssistDefendingTroopCount()

			restoreAmount = u64.MultiF64(restoreAmount, 1+season.IncProsperityMultiple)
			if amt := int64(restoreAmount); amt > 0 {
				r.handleUpdateBaseProsperity(b, amt)
			}

		}
		return true
	})
}

func (r *Realm) AddXiongNuTroop(id int64, targets []int64, monsters []*monsterdata.MonsterMasterData) (processed bool) {
	return r.queueFunc(true, func() {
		// 添加队伍
		startingBase := r.getBase(id)
		if startingBase == nil {
			logrus.Debugf("没找到匈奴大营")
			return
		}

		if !npcid.IsXiongNuNpcId(startingBase.Id()) {
			logrus.Debugf("npc竟然不是匈奴类型")
			return
		}

		if _, ok := startingBase.internalBase.(*xiongNuNpcBase); !ok {
			logrus.Debugf("npc竟然不是匈奴大营")
			return
		}

		// 加怪物了
		for idx, targetId := range targets {
			monster := monsters[idx]

			targetBase := r.getBase(targetId)
			if targetBase == nil {
				logrus.Debugf("添加匈奴出征部队，目标没在场景中找到: %d", targetId)
				continue
			}

			// 通过了, 部队出发
			troopId := startingBase.nextNpcTroopId()
			if troopId == 0 {
				continue
			}

			troop := r.newMoveNpcTroop(troopId, startingBase, startingBase.BaseLevel(), targetBase, targetBase.BaseLevel(), monster, true, 0)

			troop.onChanged()

			logrus.WithField("troopid", troop.Id()).WithField("realmid", r.id).WithField("xiongNuId", startingBase.Id()).WithField("targetid", targetId).Debug("部队加入地图成功")

			// 构建成功
			startingBase.selfTroops[troop.Id()] = troop
			targetBase.targetingTroops[troop.Id()] = troop

			targetBase.remindAttackOrRobCountChanged(r)

			//添加观察者列表
			troop.AddWatchListOnCreate(r)
			// 广播消息
			r.broadcastMaToCared(troop, addTroopTypeInvate, 0)
		}

		startingBase.updateXiongNuTarget()
	})
}

func (r *Realm) GetMaxXiongNuTroopFightingAmount(baseId, guildId int64) (suc bool, fightingAmount uint64) {
	r.queueFunc(true, func() {
		// 添加队伍
		startingBase := r.getBase(baseId)
		if startingBase == nil {
			logrus.Debugf("没找到匈奴大营")
			return
		}

		if !npcid.IsXiongNuNpcId(startingBase.Id()) {
			logrus.Debugf("npc竟然不是匈奴类型")
			return
		}

		base, ok := startingBase.internalBase.(*xiongNuNpcBase)
		if !ok {
			logrus.Debugf("npc竟然不是匈奴大营")
			return
		}

		if base.Info().GuildId() != guildId {
			logrus.Debugf("城池不是自己联盟的匈奴大营")
			return
		}

		fightingAmount = base.defenser.FightAmount()
		for _, troop := range startingBase.selfTroops {
			if troop.State() == realmface.Defending {
				// 防守中
				fightingAmount = u64.Max(fightingAmount, troop.FightAmount())
			}
		}

		suc = true
	})

	return
}

func (r *Realm) GetXiongNuTroopInfo(baseId, guildId int64) (proto *shared_proto.XiongNuBaseTroopProto) {
	r.queueFunc(true, func() {
		// 添加队伍
		startingBase := r.getBase(baseId)
		if startingBase == nil {
			logrus.Debugf("没找到匈奴大营")
			return
		}

		if !npcid.IsXiongNuNpcId(startingBase.Id()) {
			logrus.Debugf("npc竟然不是匈奴类型")
			return
		}

		base, ok := startingBase.internalBase.(*xiongNuNpcBase)
		if !ok {
			logrus.Debugf("npc竟然不是匈奴大营")
			return
		}

		if base.Info().GuildId() != guildId {
			logrus.Debugf("城池不是自己联盟的匈奴大营")
			return
		}

		proto = &shared_proto.XiongNuBaseTroopProto{}

		proto.Defender = base.defenser.doEncodeToXiongNu(r)
		for _, troop := range startingBase.selfTroops {
			if troop.State() == realmface.Defending {
				// 防守中
				proto.Assistors = append(proto.Assistors, troop.doEncodeToXiongNu(r))
			}
		}
	})

	return
}

func (r *Realm) GmRefreshBaoZangNpc() {
	r.queueFunc(false, func() {
		r.refreshBaoZangNpc(r.services.timeService.CurrentTime())
	})
}

// 刷新宝藏Npc
func (r *Realm) refreshBaoZangNpc(refreshTime time.Time) {

	if !r.lastRefreshBaoZangNpcTime.Before(refreshTime) {
		return
	}
	r.lastRefreshBaoZangNpcTime = refreshTime

	// 遍历所有的block
	r.blockManager.rangeBlock(func(b *block) (toContinue bool) {

		limit := r.levelData.Block.BlockData().BaseCountLimit
		n := u64.Sub(limit, b.GetBaseCount())
		for i := uint64(0); i < n; i++ {
			if b.GetBaseCount() >= limit {
				break
			}

			data, count := b.getNextToAddBaozNpc(r.services.datas.GetBaozNpcDataArray())
			if data == nil {
				break
			}

			if !r.addBaoZangMonsterIfAbsent(b, data, count) {
				break
			}
		}

		return true
	})

}

func (r *Realm) addBaoZangMonsterIfAbsent(block *block, data *regdata.BaozNpcData, currentCount uint64) bool {

	blockSequence := regdata.BlockSequence(block.blockX, block.blockY)

	// 生成id，倒序从currentCount开始找
	var baseId int64
	for i := uint64(0); i <= currentCount; i++ {
		id := npcid.NewBaoZangNpcId(blockSequence, u64.Sub(currentCount, i), data.Id)
		if b := r.getBase(id); b == nil {
			baseId = id
			break
		}
	}

	if baseId == 0 {
		logrus.Error("Realm野外添加宝藏怪物，但是没有目标id")
		return false
	}

	ok, baseX, baseY := r.randomBlockBasePos(block.blockX, block.blockY)
	if !ok {
		logrus.Debug("Realm野外添加宝藏怪物，block已经找不到空位，跳过")
		return false
	}

	base := r.newBaozNpcBase(baseId, data, baseX, baseY, 0, 0)
	r.addBaseToMap(base)
	r.conflict.doAddBase(base.BaseX(), base.BaseY())

	// 通知所有正在看这张地图的人
	r.broadcastBaseInfoToCared(base, addBaseTypeNewHero, 0)

	return true
}

func (r *Realm) AddHeroBaoZangMonster(data *regdata.BaozNpcData, baseX, baseY int, heroId int64, expireTime int32) (processed, success bool) {
	processed = r.queueFunc(true, func() {
		logrus.WithField("heroId", heroId).WithField("baseX", baseX).WithField("baseY", baseY).Debug("Realm野外添加召唤殷墟怪物")
		success = r.addHeroBaoZangMonster(data, baseX, baseY, heroId, expireTime)
	})

	if !processed || !success {
		r.CancelReservedPos(baseX, baseY)
	}
	return
}

// 召唤殷墟
func (r *Realm) addHeroBaoZangMonster(data *regdata.BaozNpcData, baseX, baseY int, heroId int64, expireTime int32) bool {

	if r.GetHeroBaozRoBase(heroId) != nil {
		logrus.Debug("存在已经召唤的殷墟，不能再次召唤")
		return false
	}

	blockX, blockY := r.regionData.Block.GetBlockByPos(baseX, baseY)
	block := r.blockManager.getBlock(cb.XYCube(int(blockX), int(blockY)))
	if block == nil {
		return false
	}

	currentCount := block.getBaozCount(data)
	blockSequence := regdata.BlockSequence(blockX, blockY)

	// 生成id，倒序从currentCount开始找
	var baseId int64
	for i := uint64(0); i <= currentCount; i++ {
		id := npcid.NewBaoZangNpcId(blockSequence, u64.Sub(currentCount, i), data.Id)
		if b := r.getBase(id); b == nil {
			baseId = id
			break
		}
	}

	if baseId == 0 {
		logrus.Error("Realm野外添加英雄宝藏怪物，但是没有目标id")
		return false
	}

	base := r.newBaozNpcBase(baseId, data, baseX, baseY, heroId, expireTime)
	r.addBaseToMap(base)

	// 通知所有正在看这张地图的人
	r.broadcastBaseInfoToCared(base, addBaseTypeNewHero, 0)

	return true
}

// 刷新军团Npc
func (r *Realm) refreshJunTuanNpc(refreshTime time.Time) {

	if !r.lastRefreshJunTuanNpcTime.Before(refreshTime) {
		return
	}
	r.lastRefreshJunTuanNpcTime = refreshTime

	diff := refreshTime.Sub(r.services.dep.SvrConf().GetServerStartTime())
	day := timeutil.DivideTimes(diff, timeutil.Day)

	group := r.services.datas.JunTuanNpcPlaceConfig().Must(day)

	ctime := r.services.timeService.CurrentTime()

	// 遍历所有的block
	r.blockManager.rangeBlock(func(b *block) (toContinue bool) {

		// 根据区块找到这个区块的刷新列表，循环这个列表，将超出限制数量的，删除，不足的补齐
		groupDataMap := group.GetPlaceData(cb.XYCube(int(b.blockX), int(b.blockY)))

		limit := r.levelData.Block.BlockData().BaseCountLimit

		blockSequence := regdata.BlockSequence(b.blockX, b.blockY)

		// 超出限制数量的，删掉
		for group, count := range b.junTuanMap {
			data := groupDataMap[group]
			if data == nil {
				// TODO 暂时没想到怎么处理，不处理
				continue
			}

			removeCount := u64.Sub(count, data.KeepCount)
			if removeCount > 0 {
			out:
				for i := uint64(0); i <= limit; i++ {
					id := npcid.NewJunTuanNpcId(blockSequence, i, data.Id)
					if b := r.getBase(id); b != nil {
						for _, t := range b.targetingTroops {
							// 当前有部队出征，则不删除这个主城
							if t.startingBase != b {
								continue out
							}
						}

						// 移除这个主城
						r.removeRealmBaseNoReason(b, removeBaseTypeBroken, ctime)

						removeCount = u64.Sub(removeCount, 1)
						if removeCount <= 0 {
							break out
						}
					}
				}
			}
		}

		n := u64.Sub(limit, b.GetBaseCount())
		for i := uint64(0); i < n; i++ {
			if b.GetBaseCount() >= limit {
				break
			}

			data, count := b.getNextToAddJunTuanNpc(groupDataMap)
			if data == nil {
				break
			}

			toAdd := data.Random()

			if !r.addJunTuanMonsterIfAbsent(b, toAdd, count) {
				break
			}
		}

		return true
	})

}

func (r *Realm) addJunTuanMonsterIfAbsent(block *block, data *regdata.JunTuanNpcData, currentCount uint64) bool {

	blockSequence := regdata.BlockSequence(block.blockX, block.blockY)

	// 生成id，倒序从currentCount开始找
	var baseId int64
	for i := uint64(0); i <= currentCount; i++ {
		id := npcid.NewJunTuanNpcId(blockSequence, u64.Sub(currentCount, i), data.Id)
		if b := r.getBase(id); b == nil {
			baseId = id
			break
		}
	}

	if baseId == 0 {
		logrus.Error("Realm野外添加军团怪物，但是没有目标id")
		return false
	}

	ok, baseX, baseY := r.randomBlockBasePos(block.blockX, block.blockY)
	if !ok {
		logrus.Debug("Realm野外添加军团怪物，block已经找不到空位，跳过")
		return false
	}

	base := r.newJunTuanNpcBase(baseId, data, baseX, baseY)
	// 设置协助者
	for i := uint64(0); i < data.TroopCount; i++ {

		troopId := base.nextNpcTroopId()
		if troopId == 0 {
			logrus.Errorf("创建军团怪协助部队的时候，创建不出id")
			continue
		}

		troop := r.newDefendingNpcTroop(troopId, base, base.BaseLevel(), base, base.BaseLevel(), data.Npc.Npc)
		troop.onChanged()

		// 构建成功
		base.selfTroops[troop.Id()] = troop
		base.targetingTroops[troop.Id()] = troop
		// base.remindAttackOrRobCountChanged(r)
	}

	base.updateRoBase()
	r.addBaseToMap(base)
	r.conflict.doAddBase(base.BaseX(), base.BaseY())

	// 通知所有正在看这张地图的人
	r.broadcastBaseInfoToCared(base, addBaseTypeNewHero, 0)

	return true
}

var viewNotNil = CanSeeFunc(func(area *realmface.ViewArea) bool {
	return area != nil
})

var removeStopLostProsperityMsg = region.NewS2cUpdateStopLostProsperityMsg(nil).Static()

func (r *Realm) updatePerDay() {

	ctime := r.services.timeService.CurrentTime()

	r.rangeBases(func(v *baseWithData) (toContinue bool) {
		v.TrySetStopLostProsperity(false)

		if home := GetHomeBase(v); home != nil {
			home.tryRemoveExpireKeepBaoz(ctime)
		}
		return true
	})

	r.broadcastToCared(viewNotNil, removeStopLostProsperityMsg, 0)

}

func (r *Realm) newEvent(t time.Time, f func(*Realm)) realmevent.Event {
	result := r.eventQueue.NewEvent(t, f)

	if peek := r.eventQueue.Peek(); peek == result {
		d := peek.Time().Sub(r.services.timeService.CurrentTime())
		r.nextEventTime = time.After(d)
	}

	return result
}

func (r *Realm) updateNextEventTime(e realmevent.Event) {
	if peek := r.eventQueue.Peek(); peek == e || e.Time().Before(peek.Time()) {
		d := e.Time().Sub(r.services.timeService.CurrentTime())
		r.nextEventTime = time.After(d)
	}
}

func (r *Realm) CalcMoveSpeed(targetId int64, heroRate float64) float64 {
	season := r.services.seasonService.Season()
	regionConfig := r.services.datas.RegionConfig()

	var speed float64

	rate := 1 - season.DecTroopSpeedRate + heroRate

	if targetId != 0 && !npcid.IsNpcId(targetId) {
		speed = regionConfig.BasicTroopMoveVelocityPerSecond * rate
	} else {
		speed = regionConfig.BasicTroopMoveToNpcVelocityPerSecond * rate
	}

	return math.Max(speed, regionConfig.MinTroopMoveVelocityPerSecond)
}

func (r *Realm) stop() {
	r.closeOnce.Do(func() {
		close(r.closeNotify)
	})

	<-r.loopExitNotify
}

type action struct {
	f      func()
	called chan struct{}
}

func NewAstDefendLogs() *server_proto.AllAstDefendLogProto {
	return &server_proto.AllAstDefendLogProto{Logs: make(map[int64]*server_proto.AstDefendLogListProto)}
}

func (r *Realm) AddAstDefendLog(heroId int64, ctime, startTime time.Time, name string, toAdd uint64) (succ bool) {
	if timeutil.IsZero(startTime) || ctime.Before(startTime) {
		return
	}

	logs := r.astDefendLogs.Logs[heroId]
	if logs == nil {
		logs = &server_proto.AstDefendLogListProto{HeroId: heroId}
		r.astDefendLogs.Logs[heroId] = logs
	}

	log := &shared_proto.AstDefendLogProto{}
	log.LogTime = timeutil.Marshal32(ctime)
	log.HeroName = name
	log.AddProsperity = u64.Int32(toAdd)
	log.DefendingDuration = timeutil.DurationMarshal32(ctime.Sub(startTime))
	logs.Log = append(logs.Log, log)
	if len(logs.Log) > int(r.dep.Datas().RegionConfig().AstDefendLogLimit) {
		logs.Log = logs.Log[1:]
	}

	succ = true
	return
}

func (r *Realm) ClearAstDefendLog(heroId int64) {
	delete(r.astDefendLogs.Logs, heroId)
}

func (r *Realm) GetAstDefendLogs() *server_proto.AllAstDefendLogProto {
	return r.astDefendLogs
}

func (r *Realm) GetAstDefendLogsByHero(heroId int64) (result []*shared_proto.AstDefendLogProto) {
	if l := r.astDefendLogs.Logs[heroId]; l != nil {
		return l.Log
	}
	return
}

func (r *Realm) GetAstDefendHeros(heroId int64) (heros []*shared_proto.HeroBasicProto) {
	b := r.getBase(heroId)
	if b == nil {
		return
	}
	for _, t := range b.targetingTroops {
		if t.state != realmface.Defending {
			continue
		}
		if t.startingBase == nil {
			continue
		}
		if basicProto := r.services.heroSnapshotService.GetBasicProto(t.startingBase.internalBase.Id()); basicProto != nil {
			heros = append(heros, basicProto)
		}
	}
	return
}

func (r *Realm) GetAstDefendingTroopCount(heroId int64) (count uint64) {
	b := r.getBase(heroId)
	if b == nil {
		return
	}
	return b.getAssistDefendingTroopCount()
}

func (r *Realm) GetDefendingTroopCount(heroId int64) (count uint64) {
	b := r.getBase(heroId)
	if b == nil {
		return
	}
	return b.getDefendingTroopCount()
}

func (r *Realm) UpdateProsperityBuff(heroId int64, oldAmount, newAmount, amountCapacity uint64) {
	oldBuff := r.prosperityBuffData(oldAmount, amountCapacity)
	newBuff := r.prosperityBuffData(newAmount, amountCapacity)

	if newBuff != nil {
		if oldBuff != nil && oldBuff.Id == newBuff.Id {
			return
		}
		r.services.buffService.AddBuffToSelf(newBuff, heroId)
		return
	}

	if oldBuff != nil {
		r.services.buffService.CancelGroup(heroId, oldBuff.Group)
		return
	}
}

func (r *Realm) prosperityBuffData(amount, capacity uint64) (buffData *data.BuffEffectData) {
	per := u64.DivideTimes(u64.Sub(capacity, amount)*domestic_data.Percent, capacity)

	datas := r.services.datas.GetProsperityDamageBuffDataArray()
	for _, d := range datas {
		if d.Contains(per) {
			buffData = d.BuffData
			return
		}
	}

	return
}
