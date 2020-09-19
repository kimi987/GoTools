package region

import (
	"fmt"
	"github.com/lightpaw/logrus"
	"github.com/lightpaw/male7/config/data"
	"github.com/lightpaw/male7/config/goods"
	"github.com/lightpaw/male7/config/guild_data"
	"github.com/lightpaw/male7/config/pushdata"
	"github.com/lightpaw/male7/config/regdata"
	"github.com/lightpaw/male7/config/resdata"
	"github.com/lightpaw/male7/constants"
	"github.com/lightpaw/male7/entity"
	"github.com/lightpaw/male7/entity/cb"
	"github.com/lightpaw/male7/entity/hexagon"
	"github.com/lightpaw/male7/entity/npcid"
	"github.com/lightpaw/male7/gen/iface"
	"github.com/lightpaw/male7/gen/pb/guild"
	"github.com/lightpaw/male7/gen/pb/misc"
	"github.com/lightpaw/male7/gen/pb/region"
	"github.com/lightpaw/male7/module/realm/realmerr"
	"github.com/lightpaw/male7/module/realm/realmface"
	"github.com/lightpaw/male7/pb/server_proto"
	"github.com/lightpaw/male7/pb/shared_proto"
	"github.com/lightpaw/male7/service/conflict/guilddataservice/sharedguilddata"
	"github.com/lightpaw/male7/service/conflict/heroservice/herolock"
	"github.com/lightpaw/male7/service/heromodule"
	"github.com/lightpaw/male7/service/operate_type"
	"github.com/lightpaw/male7/util"
	"github.com/lightpaw/male7/util/i32"
	"github.com/lightpaw/male7/util/i64"
	"github.com/lightpaw/male7/util/idbytes"
	"github.com/lightpaw/male7/util/imath"
	"github.com/lightpaw/male7/util/must"
	"github.com/lightpaw/male7/util/recovtimes"
	"github.com/lightpaw/male7/util/sortkeys"
	"github.com/lightpaw/male7/util/timeutil"
	"github.com/lightpaw/male7/util/u64"
	"github.com/lightpaw/pbutil"
	"math"
	"sort"
	"time"
)

func NewRegionModule(dep iface.ServiceDep, datas iface.ConfigDatas, timeService iface.TimeService, tlog iface.TlogService,
	worldService iface.WorldService, heroDataServive iface.HeroDataService, extraTimesService iface.ExtraTimesService,
	realmService iface.RealmService, guildSnapshotService iface.GuildSnapshotService, mingcWarService iface.MingcWarService,
	heroSnapshotServive iface.HeroSnapshotService, baizhanService iface.BaiZhanService, pushService iface.PushService,
	rankModule iface.RankModule, mailModule iface.MailModule) *RegionModule {

	// 初始化region信息

	m := &RegionModule{}

	m.dep = dep
	m.datas = datas
	m.time = timeService
	m.world = worldService
	m.heroDataServive = heroDataServive
	m.extraTimesService = extraTimesService
	m.mingcWarService = mingcWarService
	//m.sharedRegionDataService = sharedRegionDataService
	m.tlog = tlog

	m.realmService = realmService
	m.guildSnapshotService = guildSnapshotService
	m.heroSnapshotServive = heroSnapshotServive
	m.baizhanService = baizhanService
	m.rankModule = rankModule
	m.mailModule = mailModule
	m.pushService = pushService

	// 注册处理野怪解锁函数
	var unlockMLNFunc uint64 = constants.FunctionType_TYPE_MULTI_LEVEL_MONSTER
	heromodule.RegisterHeroEventWithSubTypeHandler("RegionModule.handleUnlockMonster", func(hero *entity.Hero, result herolock.LockResult, event shared_proto.HeroEvent, subType uint64) {
		if event == shared_proto.HeroEvent_HERO_EVENT_UNLOCK_FUNC && subType == unlockMLNFunc {
			// 如果是解锁讨伐野怪，那么重新设置初始化次数
			ctime := timeService.CurrentTime()

			multiLevelNpcTimes := hero.GetMultiLevelNpcTimes()

			extraTime := extraTimesService.MultiLevelNpcMaxTimes().TotalTimes()
			multiLevelNpcTimes.SetTimes(datas.RegionConfig().MultiLevelNpcInitTimes+extraTime, ctime, extraTime)

			result.Add(region.NewS2cUpdateMultiLevelNpcTimesMsg(multiLevelNpcTimes.StartTimeUnix32(), nil))
		}
	})

	return m
}

//gogen:iface
type RegionModule struct {
	dep iface.ServiceDep

	datas iface.ConfigDatas

	time iface.TimeService

	world iface.WorldService

	pushService iface.PushService

	heroDataServive iface.HeroDataService

	extraTimesService iface.ExtraTimesService

	realmService iface.RealmService

	guildSnapshotService iface.GuildSnapshotService

	heroSnapshotServive iface.HeroSnapshotService

	baizhanService iface.BaiZhanService

	mingcWarService iface.MingcWarService

	rankModule iface.RankModule

	mailModule iface.MailModule

	tlog iface.TlogService
}

func (m *RegionModule) getMaxLostProsperity(prosperityCapcity uint64) uint64 {
	config := m.datas.RegionConfig()
	return u64.Max(config.MaxLostProsperity, u64.MultiCoef(prosperityCapcity, config.LostProsperityCoef))
}

func (m *RegionModule) InitHeroBase(hc iface.HeroController, ctime time.Time, country uint64, addBaseType realmface.AddBaseType) (hasError bool) {
	realm, baseX, baseY := m.realmService.ReserveNewHeroHomePos(country)
	logrus.WithField("x", baseX).WithField("y", baseY).Debugf("初始化玩家主城")

	processed, err := realm.AddBase(hc.Id(), baseX, baseY, addBaseType)
	if !processed || err != nil {
		hasError = true
		logrus.WithError(err).WithField("processed", processed).Error("创建英雄时, realm.AddBase失败")
		return
	}

	return
}

//gogen:iface
func (m *RegionModule) ProcessCreateBase(proto *region.C2SCreateBaseProto, hc iface.HeroController) {
	hasError := hc.FuncWithSend(func(hero *entity.Hero, result herolock.LockResult) {
		if hero.BaseLevel() > 0 {
			// 原先有基地的，创建个P啊
			logrus.Debugf("新建基地，不是流亡状态")
			result.Add(region.ERR_CREATE_BASE_FAIL_BASE_EXIST)
			return
		}
		if hero.BaseRegion() != 0 {
			logrus.Error("重建时, baseLevel为0, 但是有baseRegion...")
			hero.ClearBase()
		}

		result.Ok()
	})

	if hasError {
		return
	}

	realm, baseX, baseY := m.realmService.ReserveRandomHomePos(realmface.RPTReborn)
	processed, err := realm.AddBase(hc.Id(), baseX, baseY, realmface.AddBaseHomeReborn)

	if !processed {
		hc.Send(region.ERR_CREATE_BASE_FAIL_SERVER_ERROR)
		return
	}

	switch err {
	case nil:
		if giftData := m.datas.EventLimitGiftConfig().GetRebornGift(); giftData != nil {
			hc.FuncWithSend(func(hero *entity.Hero, result herolock.LockResult) {
				heromodule.ActivateEventLimitGift(hero, result, giftData, m.time.CurrentTime())
				result.Ok()
			})
		}
	case realmerr.ErrLockHeroErr:
		hc.Send(region.ERR_CREATE_BASE_FAIL_SERVER_ERROR)
	case realmerr.ErrAddBaseHomeAlive:
		hc.Send(region.ERR_CREATE_BASE_FAIL_BASE_EXIST)
	case realmerr.ErrAddBaseAlreadyHasRealm:
		hc.Send(region.ERR_CREATE_BASE_FAIL_BASE_EXIST)
	default:
		logrus.WithError(err).Error("Realm.ProcessCreateBase, 返回了个没有判断的err")
		hc.Send(region.ERR_CREATE_BASE_FAIL_BASE_EXIST)
	}
}

//gogen:iface
func (m *RegionModule) ProcessFastMoveBase(proto *region.C2SFastMoveBaseProto, hc iface.HeroController) {

	goodsId := u64.FromInt32(proto.GoodsId)
	goodsData := m.datas.GetGoodsData(goodsId)
	if goodsData == nil {
		logrus.Debugf("迁移基地，物品没找到")
		hc.Send(region.ERR_FAST_MOVE_BASE_FAIL_INVALID_GOODS)
		return
	}

	if goodsData.GoodsEffect == nil {
		logrus.Debugf("迁移基地，物品不是迁城令")
		hc.Send(region.ERR_FAST_MOVE_BASE_FAIL_INVALID_GOODS)
		return
	}

	moveType := goodsData.GoodsEffect.MoveBaseType

	// 高级迁城允许自己的部队在外面
	removeSelfTroop := moveType == shared_proto.GoodsMoveBaseType_MOVE_BASE_POINT
	isFixedPos := moveType == shared_proto.GoodsMoveBaseType_MOVE_BASE_POINT

	switch moveType {
	case shared_proto.GoodsMoveBaseType_MOVE_BASE_POINT:
	case shared_proto.GoodsMoveBaseType_MOVE_BASE_RANDOM:
	case shared_proto.GoodsMoveBaseType_MOVE_BASE_GUILD:
	default:
		logrus.Debugf("迁移基地，物品不是迁城令")
		hc.Send(region.ERR_FAST_MOVE_BASE_FAIL_INVALID_GOODS)
		return
	}

	newX, newY := 0, 0
	newRealm := m.realmService.GetBigMap()
	if isFixedPos {
		newX, newY = int(proto.GetNewX()), int(proto.GetNewY())

		if !newRealm.GetMapData().IsValidHomePosition(newX, newY) {
			logrus.Debugf("迁移主城，无效的主城坐标, %v,%v", newX, newY)
			hc.Send(region.ERR_FAST_MOVE_BASE_FAIL_INVALID_POS)
			return
		}

		if !newRealm.IsPosOpened(newX, newY) {
			logrus.Debugf("迁移主城，坐标还未开放, %v,%v", newX, newY)
			hc.Send(region.ERR_FAST_MOVE_BASE_FAIL_INVALID_POS)
			return
		}

		if newRealm.IsEdgeNotHomePos(newX, newY) {
			logrus.Debugf("迁移主城，边界坐标, %v,%v", newX, newY)
			hc.Send(region.ERR_FAST_MOVE_BASE_FAIL_INVALID_POS)
			return
		}
	}

	hctx := heromodule.NewContext(m.dep, operate_type.RegionFastMoveBase)
	ctime := m.time.CurrentTime()

	var guildId int64
	var originBaseRegion int64
	var originBaseX, originBaseY int = -1, -1 // 在origin场景的x，y
	var reserveResult *entity.ReserveResult
	if hc.FuncWithSend(func(hero *entity.Hero, result herolock.LockResult) {

		if !hero.Depot().HasEnoughGoods(goodsId, 1) {
			logrus.Debugf("迁移基地，迁城令个数不足")
			result.Add(region.ERR_FAST_MOVE_BASE_FAIL_GOODS_NOT_ENOUGH)
			return
		}

		if hero.BaseLevel() <= 0 || hero.BaseRegion() == 0 {
			// 原先有基地的，创建个P啊
			logrus.Debugf("迁移基地，当前是流亡状态")
			result.Add(region.ERR_FAST_MOVE_BASE_FAIL_BASE_NOT_EXIST)
			return
		}

		if isFixedPos && hero.BaseX() == newX && hero.BaseY() == newY {
			// 坐标一致，移动个毛线撒
			logrus.Debugf("迁移基地，无效的坐标, 新坐标跟目前基地坐标一致, %v,%v", newX, newY)
			result.Add(region.ERR_FAST_MOVE_BASE_FAIL_SELF_POS)
			return
		}

		guildId = hero.GuildId()

		originBaseRegion = hero.BaseRegion()

		originBaseX = hero.BaseX()
		originBaseY = hero.BaseY()

		if originBaseRegion != 0 {
			if !removeSelfTroop {
				// 应该只有从这个城出去的才不能迁移
				for _, t := range hero.Troops() {
					if t.GetInvateInfo() != nil && t.GetInvateInfo().RegionID() == originBaseRegion {
						logrus.Debugf("迁移基地，当前有出征的武将")
						result.Add(region.ERR_FAST_MOVE_BASE_FAIL_CAPTAIN_OUT_SIDE)
						return
					}
				}
			}
		}

		// 预约扣物品
		var hasEnoughGoods bool
		if hasEnoughGoods, reserveResult = heromodule.ReserveGoods(hctx, hero, result, goodsData, 1, ctime); !hasEnoughGoods {
			logrus.Debug("迁移基地，预约扣除物品，个数不足")
			result.Add(region.ERR_FAST_MOVE_BASE_FAIL_GOODS_NOT_ENOUGH)
			return
		}

		result.Ok()
	}) {
		return
	}

	heromodule.ConfirmReserveResult(hctx, hc, reserveResult, func() (success bool) {
		var leaderX, leaderY int
		if moveType == shared_proto.GoodsMoveBaseType_MOVE_BASE_GUILD {
			// 验证联盟迁城令
			guild := m.guildSnapshotService.GetSnapshot(guildId)
			if guildId == 0 || guild == nil {
				logrus.Debug("迁移基地，不在联盟中")
				hc.Send(region.ERR_FAST_MOVE_BASE_FAIL_GUILD_MOVE_BASE_NOT_IN_GUILD)
				return
			}

			if guild.IsNpcGuild {
				logrus.Debug("迁移基地，Npc联盟")
				hc.Send(region.ERR_FAST_MOVE_BASE_FAIL_GUILD_MOVE_BASE_LEADER_IS_NPC)
				return
			}

			if guild.LeaderId == hc.Id() {
				logrus.Debug("迁移基地，自己是盟主")
				hc.Send(region.ERR_FAST_MOVE_BASE_FAIL_GUILD_MOVE_BASE_IS_LEADER)
				return
			}

			leader := m.heroSnapshotServive.Get(guild.LeaderId)
			if leader == nil {
				logrus.Debug("迁移基地，盟主不存在")
				hc.Send(region.ERR_FAST_MOVE_BASE_FAIL_GUILD_MOVE_BASE_LEADER_IS_NPC)
				return
			}

			if leader.BaseLevel == 0 || leader.BaseRegion == 0 {
				logrus.Debug("迁移基地，盟主已流亡")
				hc.Send(region.ERR_FAST_MOVE_BASE_FAIL_GUILD_MOVE_BASE_LEADER_NO_BASE)
				return
			}

			leaderX, leaderY = leader.BaseX, leader.BaseY
			if m.realmService.GetBigMap().AroundBase(originBaseX, originBaseY, leaderX, leaderY) {
				hc.Send(region.ERR_FAST_MOVE_BASE_FAIL_GUILD_MOVE_BASE_ALREADY_IN_LEADER_AROUND)
				return
			}
		}

		switch moveType {
		case shared_proto.GoodsMoveBaseType_MOVE_BASE_POINT:
			// 具体位置，尝试预留空位
			if ok := newRealm.ReservePosForMoveBase(originBaseX, originBaseY, newX, newY); !ok {
				logrus.Debug("迁移基地, 位置冲突")
				hc.Send(region.ERR_FAST_MOVE_BASE_FAIL_TOO_CLOSE_OTHER)
				return
			}

		case shared_proto.GoodsMoveBaseType_MOVE_BASE_RANDOM:
			// 随机位置，并且预留了空位
			newRealm, newX, newY = m.realmService.ReserveRandomHomePos(realmface.RPTRandom)

		case shared_proto.GoodsMoveBaseType_MOVE_BASE_GUILD:
			newRealm = m.realmService.GetBigMap()
			x, y, ok := newRealm.RandomAroundBase(leaderX, leaderY)
			if !ok {
				logrus.Warnf("使用联盟随机迁城令，获得随机坐标失败")
				hc.Send(region.ERR_FAST_MOVE_BASE_FAIL_GUILD_MOVE_BASE_FULL)
				return
			}
			newX, newY = x, y

		default:
			logrus.Error("迁移基地，物品不是迁城令（前面不是判断过了吗）")
			hc.Send(region.ERR_FAST_MOVE_BASE_FAIL_INVALID_GOODS)
			return
		}

		success = m.realmService.DoMoveBase(moveType, newRealm, hc, originBaseX, originBaseY, newX, newY, removeSelfTroop)
		return
	})
}

//gogen:iface c2s_upgrade_base
func (m *RegionModule) ProcessUpgradeBase(hc iface.HeroController) {
	var realmId int64

	hasError := hc.FuncWithSend(func(hero *entity.Hero, result herolock.LockResult) {
		if hero.BaseLevel() <= 0 || hero.Prosperity() <= 0 {
			logrus.Debugf("升级基地，当前是流亡状态")
			result.Add(region.ERR_UPGRADE_BASE_FAIL_BASE_NOT_EXIST)
			return
		}

		nextLevelData := m.datas.GetBaseLevelData(hero.BaseLevel() + 1)
		if nextLevelData == nil {
			logrus.Debugf("升级基地，已经达到最大等级")
			result.Add(region.ERR_UPGRADE_BASE_FAIL_BASE_MAX_LEVEL)
			return
		}

		if hero.Prosperity() < nextLevelData.Prosperity {
			logrus.Debugf("升级基地，繁荣度不足")
			result.Add(region.ERR_UPGRADE_BASE_FAIL_PROSPRITY_NOT_ENOUGH)
			return
		}

		realmId = hero.BaseRegion()
		result.Ok()
		return
	})

	if hasError {
		return
	}

	realm := m.realmService.GetBigMap()
	if realm == nil {
		logrus.WithField("realmid", realmId).Error("升级基地, 没找到英雄所属的realm")
		hc.Send(region.ERR_UPGRADE_BASE_FAIL_BASE_NOT_EXIST)
		return
	}

	processed, err := realm.UpgradeBase(hc)
	if !processed {
		logrus.Error("升级基地, 没有process")
		hc.Send(region.ERR_UPGRADE_BASE_FAIL_SERVER_ERROR)
		return
	}

	switch err {
	case nil:
		hc.Send(region.UPGRADE_BASE_S2C)
	case realmerr.ErrUpgradeBaseNotEnoughProsperity:
		hc.Send(region.ERR_UPGRADE_BASE_FAIL_PROSPRITY_NOT_ENOUGH)
	case realmerr.ErrUpgradeBaseAlreadyMax:
		hc.Send(region.ERR_UPGRADE_BASE_FAIL_BASE_MAX_LEVEL)
	case realmerr.ErrUpgradeBaseHomeNotAlive:
		hc.Send(region.ERR_UPGRADE_BASE_FAIL_BASE_NOT_EXIST)
	case realmerr.ErrUpgradeBaseNotMyRealm:
		hc.Send(region.ERR_UPGRADE_BASE_FAIL_BASE_NOT_EXIST)
	default:
		hc.Send(region.ERR_UPGRADE_BASE_FAIL_SERVER_ERROR)
		logrus.WithError(err).Error("region.ProcessUpgradeBase有err没有处理")
	}
}

var (
	msgSwitchActionOpen  = region.NewS2cSwitchActionMsg(true).Static()
	msgSwitchActionClose = region.NewS2cSwitchActionMsg(false).Static()
)

//gogen:iface
func (m *RegionModule) ProcessUpdateSelfView(proto *region.C2SUpdateSelfViewProto, hc iface.HeroController) {
	posX := int(proto.PosX)
	posY := int(proto.PosY)
	lenX := int(proto.LenX)
	lenY := int(proto.LenY)

	current := hc.GetViewArea()
	if current != nil {

		//if m.time.CurrentTime().Sub(current.UpdateTime) < 500 * time.Millisecond {
		//	hc.Disconnect(misc.ErrDisconectReasonFailMsgRate)
		//	return
		//}

		if current.CenterX == posX && current.CenterY == posY {
			// 假装设置成功了
			hc.Send(region.NewS2cUpdateSelfViewMsg(
				int32(current.MinX), int32(current.MinY),
				int32(current.MaxX), int32(current.MaxY)))
			return
		}
	}

	realm := m.realmService.GetBigMap()

	realm.StartCareRealm(hc, posX, posY, lenX, lenY)

	if current := hc.GetViewArea(); current != nil {
		hc.Send(region.NewS2cUpdateSelfViewMsg(
			int32(current.MinX), int32(current.MinY),
			int32(current.MaxX), int32(current.MaxY)))
	} else {
		hc.Send(region.CLOSE_VIEW_S2C)
	}

}

//gogen:iface c2s_close_view
func (m *RegionModule) ProcessCloseView(hc iface.HeroController) {
	current := hc.GetViewArea()
	if current == nil {
		// 假装设置成功了
		hc.Send(region.CLOSE_VIEW_S2C)
		return
	}

	realm := m.realmService.GetBigMap()

	realm.StopCareRealm(hc)
	hc.Send(region.CLOSE_VIEW_S2C)
}

//gogen:iface
func (m *RegionModule) ProcessRequestTroopUnit(proto *region.C2SRequestTroopUnitProto, hc iface.HeroController) {

	troopId, ok := idbytes.ToId(proto.TroopId)
	if !ok {
		logrus.Debug("查询军情，无效的id")
		hc.Send(region.ERR_REQUEST_TROOP_UNIT_FAIL_INVALID_ID)
		return
	}

	heroId := entity.GetTroopHeroId(troopId)
	//if npcid.IsNpcId(heroId) {
	//	// 查询Npc军情？
	//	logrus.Debug("查询军情，NpcId")
	//	hc.Send(region.ERR_REQUEST_TROOP_UNIT_FAIL_INVALID_ID)
	//	return
	//}

	realm := m.realmService.GetBigMap()

	processed, err := realm.QueryTroopUnit(hc, heroId, troopId)
	if !processed {
		logrus.Error("查询军情，!processed")
		hc.Send(region.ERR_REQUEST_TROOP_UNIT_FAIL_SERVER_ERROR)
		return
	}

	switch err {
	case nil:
		// nothing
	case realmerr.ErrGetMilitaryBaseNotFound:
		hc.Send(region.ERR_REQUEST_TROOP_UNIT_FAIL_INVALID_ID)
		return
	case realmerr.ErrGetMilitaryTroopNotFound:
		hc.Send(region.ERR_REQUEST_TROOP_UNIT_FAIL_INVALID_ID)
		return
	default:
		logrus.Error("查询军情，遇到未处理的错误")
		hc.Send(region.ERR_REQUEST_TROOP_UNIT_FAIL_SERVER_ERROR)
		return
	}

}

// action
//gogen:iface
func (m *RegionModule) ProcessSwitchAction(proto *region.C2SSwitchActionProto, hc iface.HeroController) {
	if !proto.Open {
		hc.Send(msgSwitchActionClose)
		hc.SetCareCondition(nil)
		return
	}

	hc.Send(msgSwitchActionOpen)
	hc.SetCareCondition(&server_proto.MilitaryConditionProto{})
	m.realmService.StartCareMilitary(hc)
}

//出征相关
//gogen:iface
func (m *RegionModule) ProcessPreInvasionTarget(proto *region.C2SPreInvasionTargetProto, hc iface.HeroController) {

	// 检查id
	targetId, ok := idbytes.ToId(proto.Target)
	if !ok {
		logrus.Errorf("RegionModule.ProcessInvasion 解析Target id失败, %s", proto.Target)
		hc.Send(region.ERR_PRE_INVASION_TARGET_FAIL_INVALID_TARGET)
		return
	}

	if npcid.IsNpcId(targetId) {
		var found bool

		if npcid.IsHomeNpcId(targetId) {
			hc.FuncWithSend(func(hero *entity.Hero, result herolock.LockResult) {
				hero.WalkHomeNpcBase(func(base *entity.HomeNpcBase) bool {
					if base.Id() == targetId {
						npcData := base.GetData().Data.Npc
						result.Add(region.NewS2cPreInvasionTargetMsg(npcData.Icon.Id, u64.Int32(npcData.Level), 0, 0))
						found = true
						return true
					}
					return false
				})
			})
		}

		if !found {
			hc.Send(region.ERR_PRE_INVASION_TARGET_FAIL_INVALID_TARGET)
		}

		return
	}

	var toSend pbutil.Buffer
	m.heroDataServive.Func(targetId, func(hero *entity.Hero, err error) (heroChanged bool) {

		if err != nil {
			toSend = region.ERR_PRE_INVASION_TARGET_FAIL_INVALID_TARGET
			return
		}

		toSend = region.NewS2cPreInvasionTargetMsg(hero.Head(), u64.Int32(hero.Level()), u64.Int32(hero.Tower().HistoryMaxFloor()), u64.Int32(m.baizhanService.GetJunXianLevel(targetId)))
		return
	})

	hc.Send(toSend)
}

//gogen:iface
func (m *RegionModule) ProcessWatchBaseUnit(proto *region.C2SWatchBaseUnitProto, hc iface.HeroController) {

	targetId, ok := idbytes.ToId(proto.Target)
	if !ok {
		logrus.Debug("观察野外主城，id无效")
		hc.Send(region.ERR_WATCH_BASE_UNIT_FAIL_INVALID_TARGET)
		return
	}

	if npcid.IsNpcId(targetId) {

		if npcid.IsHomeNpcId(targetId) {
			var found bool
			hc.FuncWithSend(func(hero *entity.Hero, result herolock.LockResult) {
				hero.WalkHomeNpcBase(func(base *entity.HomeNpcBase) bool {
					if base.Id() == targetId {

						prosperity := u64.Int32(base.GetData().Data.ProsperityCapcity)
						npcData := base.GetData().Data.Npc

						result.Add(region.NewS2cWatchBaseUnitMsg(proto.Target, "", u64.Int32(npcData.FightAmount), prosperity, prosperity, npcData.Icon.Id, u64.Int32(npcData.Level), 0, 0, 0, nil, nil, 0, 0))
						found = true
						return true
					}
					return false
				})
			})

			if !found {
				hc.Send(region.ERR_WATCH_BASE_UNIT_FAIL_INVALID_TARGET)
			}

			return
		}

		if npcid.IsBaoZangNpcId(targetId) {
			if base := m.realmService.GetBigMap().GetRoBase(targetId); base != nil {

				var heroBytes []byte
				if base.HeroId != 0 {
					if hero := m.heroSnapshotServive.Get(base.HeroId); hero != nil {
						heroBytes = hero.EncodeBasic4ClientBytes()
					}
				}

				hc.Send(region.NewS2cWatchBaseUnitMsg(proto.Target,
					"",
					u64.Int32(base.GetFightAmount()),
					u64.Int32(base.GetProsperity()),
					u64.Int32(base.GetProsperityCapcity()),
					"",
					u64.Int32(base.GetLevel()),
					0,
					0,
					u64.Int32(base.GetSoldier()),
					u64.Int32Array(base.GetCaptainSoldier()),
					heroBytes,
					base.HeroEndTime,
					base.HeroType,
				))
				return
			}
		}

		if npcid.IsJunTuanNpcId(targetId) {
			dataId := npcid.GetNpcDataId(targetId)
			data := m.datas.GetJunTuanNpcData(dataId)
			if data == nil {
				hc.Send(region.ERR_WATCH_BASE_UNIT_FAIL_INVALID_TARGET)
				return
			}

			if base := m.realmService.GetBigMap().GetRoBase(targetId); base != nil {

				hc.Send(region.NewS2cWatchBaseUnitProtoMsg(&region.S2CWatchBaseUnitProto{
					Target:           proto.Target,
					FightAmount:      u64.Int32(data.Npc.Npc.FightAmount),
					Prosprity:        u64.Int32(base.GetProsperity()),
					ProsprityCapcity: u64.Int32(base.GetProsperityCapcity()),
					Soldier:          u64.Int32(base.GetSoldier()),
				}))
				return
			}
		}

		hc.Send(region.ERR_WATCH_BASE_UNIT_FAIL_INVALID_TARGET)
		return
	}

	var toSend pbutil.Buffer
	m.heroDataServive.Func(targetId, func(hero *entity.Hero, err error) (heroChanged bool) {

		if err != nil {
			toSend = region.ERR_WATCH_BASE_UNIT_FAIL_INVALID_TARGET
			return
		}

		var guildName string
		if g := m.guildSnapshotService.GetSnapshot(hero.GuildId()); g != nil {
			guildName = g.Name
		}

		defenserFightAmount := hero.GetHomeDefenserFightAmount()
		if copyDefenser := hero.GetCopyDefenser(); copyDefenser != nil {
			if t := hero.GetHomeDefenser(); t == nil || t.IsOutside() {
				var captainFightAmounts []uint64
				for _, cis := range copyDefenser.GetCaptains() {
					if cis != nil {
						c := hero.Military().Captain(cis.GetId())
						if c != nil && cis.GetSoldier() > 0 {
							captainFightAmounts = append(captainFightAmounts, cis.GetFightAmount())
						}
					}
				}
				defenserFightAmount = data.TroopFightAmount(captainFightAmounts...)
			}
		}

		toSend = region.NewS2cWatchBaseUnitMsg(hero.IdBytes(), guildName, u64.Int32(defenserFightAmount), u64.Int32(hero.Prosperity()), u64.Int32(hero.ProsperityCapcity()), hero.Head(), u64.Int32(hero.Level()), u64.Int32(hero.Tower().HistoryMaxFloor()), u64.Int32(m.baizhanService.GetJunXianLevel(targetId)), 0, nil, nil, 0, 0)
		return
	})

	hc.Send(toSend)
}

//gogen:iface
func (m *RegionModule) ProcessRequestMilitaryPush(proto *region.C2SRequestMilitaryPushProto, hc iface.HeroController) {

	cond := &server_proto.MilitaryConditionProto{}
	cond.IsOr = true

	switch {
	case proto.MainMilitary:
		// 我的军情
		cond.Attributes = append(cond.Attributes, &server_proto.MilitaryAttributeProto{
			StartBaseId:  hc.Id(), // 我的军情（从我家出发的）
			TargetBaseId: hc.Id(), // 我的军情(往我家跑的)
		})

		// 盟友被持续掠夺的军情
		if guildId, ok := hc.LockGetGuildId(); ok && guildId != 0 {
			cond.Conditions = append(cond.Conditions, &server_proto.MilitaryConditionProto{
				IsOr: false,
				Attributes: []*server_proto.MilitaryAttributeProto{
					&server_proto.MilitaryAttributeProto{
						TargetBaseGuildId: guildId, // 我的联盟军情（往盟友家跑的）
					},
					&server_proto.MilitaryAttributeProto{
						TroopState: 5, // 持续掠夺
					},
				},
			})

			cond.Conditions = append(cond.Conditions, &server_proto.MilitaryConditionProto{
				IsOr: false,
				Attributes: []*server_proto.MilitaryAttributeProto{
					&server_proto.MilitaryAttributeProto{
						StartBaseGuildId: guildId, // 我的联盟军情（往盟友家跑的）
					},
					&server_proto.MilitaryAttributeProto{
						JoinAssemblyHeroId: hc.Id(), // 我加入的集结
					},
				},
			})
		}

	case proto.GuildMilitary:

		if guildId, ok := hc.LockGetGuildId(); ok && guildId != 0 {
			cond.Attributes = append(cond.Attributes, &server_proto.MilitaryAttributeProto{
				StartBaseGuildId:  guildId, // 我的联盟军情（从盟友家出发的）
				TargetBaseGuildId: guildId, // 我的联盟军情（往盟友家跑的）
			})
		} else {
			logrus.Debug("请求推送军情，自己没有联盟，不能请求联盟军情")
			hc.Send(region.ERR_REQUEST_MILITARY_PUSH_FAIL_NOT_IN_GUILD)
			return
		}

	case len(proto.ToTarget) > 0:
		targetId, ok := idbytes.ToId(proto.ToTarget)
		if !ok {
			logrus.Debug("请求推送军情，无效的ToTarget")
			hc.Send(region.ERR_REQUEST_MILITARY_PUSH_FAIL_INVALID_ID)
			return
		}

		// 我的军情
		cond.Attributes = append(cond.Attributes, &server_proto.MilitaryAttributeProto{
			TargetBaseId: targetId, // 我的军情(往我家跑的)
		})

		if proto.ToTargetBase {
			// TODO 关注主城详情变化
		}
	default:

		if len(proto.FromTarget) <= 0 {
			// 玩家不再关心推送信息
			hc.SetCareCondition(nil)
			hc.Send(region.NewS2cRequestMilitaryPushMsg(proto.MainMilitary, proto.GuildMilitary, proto.ToTarget, proto.ToTargetBase, proto.FromTarget))
			return
		}
	}

	if len(proto.FromTarget) > 0 {
		targetId, ok := idbytes.ToId(proto.FromTarget)
		if !ok {
			logrus.Debug("请求推送军情，无效的FromTarget")
			hc.Send(region.ERR_REQUEST_MILITARY_PUSH_FAIL_INVALID_ID)
			return
		}

		cond.Attributes = append(cond.Attributes, &server_proto.MilitaryAttributeProto{
			StartBaseId: targetId, // 我的军情(往我家跑的)
		})
	}

	// 设置care数据
	hc.SetCareCondition(cond)
	hc.Send(region.NewS2cRequestMilitaryPushMsg(proto.MainMilitary, proto.GuildMilitary, proto.ToTarget, proto.ToTargetBase, proto.FromTarget))

	m.realmService.StartCareMilitary(hc)
}

//gogen:iface
func (m *RegionModule) ProcessInvasion(proto *region.C2SInvasionProto, hc iface.HeroController) {

	// 检查id
	targetId, ok := idbytes.ToId(proto.Target)
	if !ok {
		logrus.Errorf("RegionModule.ProcessInvasion 解析Target id失败, %s", proto.Target)
		hc.Send(region.ERR_INVASION_FAIL_INVALID_TARGET)
		return
	}

	heroId := hc.Id()
	if targetId == heroId {
		logrus.Debugf("出征，目标id跟自己的id一样 %s", proto.Target)
		hc.Send(region.ERR_INVASION_FAIL_INVALID_TARGET)
		return
	}

	if npcid.IsJunTuanNpcId(targetId) {
		logrus.Debugf("出征，目标不能是军团怪")
		hc.Send(region.ERR_INVASION_FAIL_INVALID_TARGET)
		return
	}

	//realmId := int64(proto.MapId)
	//
	//if realmId == 0 {
	//	// 临时处理一下, 向前兼容
	//	logrus.Debug("客户端的invasion必须要带map_id, 临时取他主城作为map_id")
	//	hc.FuncNotError(func(hero *entity.Hero) (heroChanged bool) {
	//		realmId = hero.BaseRegion()
	//		return false
	//	})
	//}

	//realm := m.realmService.GetRealm(realmId)
	realm := m.realmService.GetBigMap()
	if realm == nil {
		logrus.Debugf("出征，出征地图没找到")
		hc.Send(region.ERR_INVASION_FAIL_NOT_SAME_MAP)
		return
	}

	operate := shared_proto.TroopOperate(proto.Operate)

	troopIndex := u64.FromInt32(proto.TroopIndex - 1)

	targetLevel := u64.FromInt32(proto.TargetLevel)

	switch operate {
	case shared_proto.TroopOperate_ToAssist:
		// 什么都不干
	case shared_proto.TroopOperate_ToInvasion:
		// 什么都不干
	case shared_proto.TroopOperate_ToWorkshopBuild, shared_proto.TroopOperate_ToWorkshopProd, shared_proto.TroopOperate_ToWorkshopPrize:

		if npcid.GetNpcIdType(targetId) != npcid.NpcType_Guild {
			logrus.WithField("operate", operate).Debugf("出征目标不是联盟工坊，不能使用工坊的操作类型")
			hc.Send(region.ERR_INVASION_FAIL_INVALID_TARGET)
			return
		}

	case shared_proto.TroopOperate_ToInvestigate:
		if npcid.IsNpcId(targetId) {
			logrus.WithField("operate", operate).Debugf("出征目标不是英雄，不能使用侦查的操作类型")
			hc.Send(region.ERR_INVASION_FAIL_INVALID_TARGET)
			return
		}
	default:
		logrus.WithField("operate", operate).Debugf("野外出征，无效的出征类型")
		hc.Send(region.ERR_INVASION_FAIL_INVALID_TARGET_INVATION)
		return
	}

	isInvasionOperate := operate == shared_proto.TroopOperate_ToInvasion ||
		operate == shared_proto.TroopOperate_ToInvestigate

	if !npcid.IsNpcId(targetId) {
		if isInvasionOperate {
			if _, joined := m.mingcWarService.JoiningFightMingc(hc.Id()); joined {
				hc.Send(region.ERR_INVASION_FAIL_IN_MC_WAR_FIGHT)
				return
			}
		}
	}

	//todo 讨伐令数量 vip 验证
	var npcTimes uint64

	// 各种检查
	hasError := hc.FuncWithSend(func(hero *entity.Hero, result herolock.LockResult) {

		// 检查主城是否存在
		if hero.BaseRegion() == 0 || hero.BaseLevel() <= 0 || hero.Prosperity() <= 0 {
			logrus.Debugf("出征，主城流亡了")
			result.Add(region.ERR_INVASION_FAIL_NO_BASE_IN_MAP)
			return
		}

		ctime := m.time.CurrentTime()

		if !npcid.IsNpcId(targetId) {
			if isInvasionOperate {
				// 自己处于新手免战，不能出征打别人
				if ctime.Before(hero.GetNewHeroMianDisappearTime()) &&
					hero.GetNewHeroMianDisappearTime().Equal(hero.GetMianDisappearTime()) {
					logrus.Debugf("出征，新手免战期间不能出征玩家")
					result.Add(region.ERR_INVASION_FAIL_MIAN)
					return
				}
			}
		}

		t := hero.GetTroopByIndex(troopIndex)
		if t == nil {
			logrus.Debugf("出征，无效的部队编号")
			result.Add(region.ERR_INVASION_FAIL_INVALID_TROOP_INDEX)
			return
		}

		if t.IsOutside() {
			logrus.Debugf("出征，队伍出征中")
			result.Add(region.ERR_INVASION_FAIL_OUTSIDE)
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
			logrus.Debugf("出征，队伍没有士兵")
			result.Add(region.ERR_INVASION_FAIL_NO_SOLDIER)
			return
		}

		if !npcid.IsNpcId(targetId) {
			// 出征次数
			if hero.GetInvaseHeroTimes().Times(ctime, 0) <= 0 {

				// 自动使用物品，检查是否有物品
				if proto.GoodsId == 0 {
					logrus.Debugf("出征玩家，讨伐次数不足")
					result.Add(region.ERR_INVASION_FAIL_MLN_TIMES_LIMIT)
					return
				} else {
					goodsData := m.datas.GetGoodsData(u64.FromInt32(proto.GoodsId))
					if goodsData == nil {
						logrus.Debug("出征玩家使用讨伐令，无效的物品id")
						result.Add(region.ERR_INVASION_FAIL_INVALID_GOODS)
						return
					}

					if goodsData.GoodsEffect == nil || !goodsData.GoodsEffect.AddInvaseHeroTimes {
						logrus.Debug("出征玩家使用讨伐令，不是讨伐令")
						result.Add(region.ERR_INVASION_FAIL_INVALID_GOODS)
						return
					}

					if !m.useInvaseHeroTimesGoods(hero, result, goodsData, proto.AutoBuy, ctime,
						region.ERR_INVASION_FAIL_MLN_TIMES_LIMIT,
						region.ERR_INVASION_FAIL_COST_NOT_ENOUGH,
						region.ERR_INVASION_FAIL_COST_NOT_ENOUGH, ) {

						logrus.Debug("出征玩家，使用讨伐令失败")
						return
					}

					// 使用物品加次数成功
				}
			}
		}

		// 特殊野怪判定
		if npcid.IsMultiLevelMonsterNpcId(targetId) {
			// 功能是否已经解锁
			if !hero.Function().IsFunctionOpened(constants.FunctionType_TYPE_MULTI_LEVEL_MONSTER) {
				logrus.Debugf("出征野怪，但是功能还未开启")
				result.Add(region.ERR_INVASION_FAIL_MLN_FUNC_LOCKED)
				return
			}

			monsterData := m.datas.GetRegionMultiLevelNpcData(npcid.GetNpcDataId(targetId))
			if monsterData == nil {
				logrus.Debugf("出征，根据id找不到MultiLevelNpcData")
				result.Add(region.ERR_INVASION_FAIL_INVALID_TARGET)
				return
			}

			if hero.GetMultiLevelNpcPassLevel()+1 < targetLevel {
				logrus.Debugf("出征，选择的野怪未解锁")
				result.Add(region.ERR_INVASION_FAIL_TARGET_LEVEL_LOCKED)
				return
			}

			// 讨伐野怪次数，兼容旧版本，默认1次
			npcTimes = u64.Max(1, u64.FromInt32(proto.MultiLevelMonsterCount))
			//if npcTimes > 1+m.datas.GetVipLevelData(hero.VipLevel()).InvadeMultiLevelMonsterOnceCount {
			//	result.Add(region.ERR_INVASION_FAIL_MULTI_LEVEL_MONSTER_COUNT_VIP_LIMIT)
			//	return
			//}

			// 同一只野怪，只能一个部队出征
			var useTimes uint64 = 0
			for _, t := range hero.Troops() {
				if info := t.GetInvateInfo(); info != nil {
					isInvadeOrRobbing := info.State() == realmface.MovingToInvade || info.State() == realmface.Robbing
					if info.OriginTargetID() == targetId && isInvadeOrRobbing {

						logrus.Debugf("出征，存在另一个队伍对野怪进行出征")
						result.Add(region.ERR_INVASION_FAIL_DUPLICATE_TARGET)
						return
					}

					if npcid.IsMultiLevelMonsterNpcId(info.OriginTargetID()) && isInvadeOrRobbing {
						useTimes += u64.Max(1, info.NpcTimes())
					}
				}
			}

			needGoodsCount := u64.Sub(useTimes+npcTimes, hero.GetMultiLevelNpcTimes().Times(ctime, m.extraTimesService.MultiLevelNpcMaxTimes().TotalTimes()))

			// 出征次数
			if needGoodsCount > 0 {
				// 自动使用物品，检查是否有物品
				if proto.GoodsId == 0 {
					logrus.Debugf("出征野怪，讨伐野怪次数不足")
					result.Add(region.ERR_INVASION_FAIL_MLN_TIMES_LIMIT)
					return
				} else {
					goodsData := m.datas.GetGoodsData(u64.FromInt32(proto.GoodsId))
					if goodsData == nil {
						logrus.Debug("出征使用讨伐令，无效的物品id")
						result.Add(region.ERR_INVASION_FAIL_INVALID_GOODS)
						return
					}

					if goodsData.GoodsEffect == nil || !goodsData.GoodsEffect.AddMultiLevelNpcTimes {
						logrus.Debug("出征使用讨伐令，不是讨伐令")
						result.Add(region.ERR_INVASION_FAIL_INVALID_GOODS)
						return
					}

					if !m.useMultiLevelNpcTimesGoods(hero, result, goodsData, needGoodsCount, proto.AutoBuy, ctime,
						region.ERR_INVASION_FAIL_MLN_TIMES_LIMIT,
						region.ERR_INVASION_FAIL_COST_NOT_ENOUGH,
						region.ERR_INVASION_FAIL_COST_NOT_ENOUGH, ) {

						logrus.Debug("出征野怪，使用讨伐令失败")
						return
					}

					// 使用物品加次数成功
				}

			}
		} else if npcid.IsBaoZangNpcId(targetId) {
			data := m.datas.GetBaozNpcData(npcid.GetNpcDataId(targetId))
			if data == nil {
				logrus.Debugf("出征，根据id找不到BaozNpcData")
				result.Add(region.ERR_INVASION_FAIL_INVALID_TARGET)
				return
			}

			if hero.Level() < data.RequiredHeroLevel {
				logrus.Debugf("出征，君主等级不足")
				result.Add(region.ERR_INVASION_FAIL_REQUIRED_HERO_LEVEL)
				return
			}

		}

		result.Ok()
	})

	if hasError {
		return
	}

	processed, err := realm.Invasion(hc, operate, targetId, targetLevel, troopIndex, npcTimes)
	if !processed {
		hc.Send(region.ERR_INVASION_FAIL_SERVER_ERROR)
		return
	}

	switch err {
	case nil:
		hc.Send(region.NewS2cInvasionMsg(proto.Target, proto.TroopIndex))

	case realmerr.ErrLockHeroErr:
		hc.Send(region.ERR_INVASION_FAIL_SERVER_ERROR)

	case realmerr.ErrInvasionInvalidTarget:
		hc.Send(region.ERR_INVASION_FAIL_INVALID_TARGET)

	case realmerr.ErrInvasionTargetNotExist:
		hc.Send(region.ERR_INVASION_FAIL_TARGET_NOT_EXIST)

	case realmerr.ErrInvasionSelfNoBase:
		hc.Send(region.ERR_INVASION_FAIL_NO_BASE_IN_MAP)

	case realmerr.ErrInvasionEmptyGeneral:
		hc.Send(region.ERR_INVASION_FAIL_INVALID_TROOP_INDEX)

	case realmerr.ErrInvasionGeneralOutside:
		hc.Send(region.ERR_INVASION_FAIL_OUTSIDE)

	case realmerr.ErrInvasionNoSoldier:
		hc.Send(region.ERR_INVASION_FAIL_NO_SOLDIER)

	case realmerr.ErrInvasionInvalidRelation:
		hc.Send(region.ERR_INVASION_FAIL_INVALID_TARGET_INVATION)

	case realmerr.ErrInvasionInvalidTroopIndex:
		hc.Send(region.ERR_INVASION_FAIL_INVALID_TROOP_INDEX)

	case realmerr.ErrInvasionMian:
		hc.Send(region.ERR_INVASION_FAIL_MIAN)

	case realmerr.ErrInvasionTodayJoinXiongNu:
		hc.Send(region.ERR_INVASION_FAIL_TODAY_JOIN_XIONG_NU)

	default:
		logrus.WithError(err).Error("region.ProcessInvasion有err没有处理")
		hc.Send(region.ERR_INVASION_FAIL_SERVER_ERROR)
	}
}

// 召回
//gogen:iface
func (m *RegionModule) ProcessCancelInvasion(proto *region.C2SCancelInvasionProto, hc iface.HeroController) {
	troopId, ok := idbytes.ToId(proto.Id)
	if !ok {
		logrus.Debugf("召回，无效的军情id")
		hc.Send(region.ERR_CANCEL_INVASION_FAIL_INVALID_ID)
		return
	}

	var realmId int64
	hasError := hc.FuncWithSend(func(hero *entity.Hero, result herolock.LockResult) {
		troop := hero.Troop(troopId)
		if troop == nil {
			logrus.Debugf("召回，没找到部队id")
			result.Add(region.ERR_CANCEL_INVASION_FAIL_INVALID_ID)
			return
		}

		invateInfo := troop.GetInvateInfo()
		if invateInfo == nil {
			logrus.Debugf("召回，部队没有出征")
			result.Add(region.ERR_CANCEL_INVASION_FAIL_INVALID_ID)
			return
		}

		if invateInfo.State() == realmface.InvadeMovingBack ||
			invateInfo.State() == realmface.AssistMovingBack ||
			invateInfo.State() == realmface.AssemblyMovingBack ||
			invateInfo.State() == realmface.InvesigateMovingBack {
			logrus.Debugf("召回，部队已经在回来路上了")
			result.Add(region.ERR_CANCEL_INVASION_FAIL_BACKING)
			return
		}

		realmId = invateInfo.RegionID()
		result.Ok()
		return
	})

	if hasError {
		return
	}

	realm := m.realmService.GetBigMap()
	if realm == nil {
		logrus.WithField("realmid", realmId).Error("CancelInvasion时, hero里拿出来的部队所在realmid竟然不存在")
		hc.Send(region.ERR_CANCEL_INVASION_FAIL_SERVER_ERROR)
		return
	}

	processed, err := realm.CancelInvasion(hc, troopId)
	if !processed {
		hc.Send(region.ERR_CANCEL_INVASION_FAIL_SERVER_ERROR)
		return
	}

	switch err {
	case nil:
		hc.Send(region.NewS2cCancelInvasionMsg(proto.Id))
	case realmerr.ErrCancelInvasionTroopAlreadyHome:
		hc.Send(region.ERR_CANCEL_INVASION_FAIL_INVALID_ID)
	case realmerr.ErrCancelInvasionTroopAlreadyBacking:
		hc.Send(region.ERR_CANCEL_INVASION_FAIL_BACKING)
	case realmerr.ErrCancelInvasionTroopNotFound:
		hc.Send(region.ERR_CANCEL_INVASION_FAIL_INVALID_ID)
	case realmerr.ErrCancelInvasionTroopAssemblyStarted:
		hc.Send(region.ERR_CANCEL_INVASION_FAIL_ASSEMBLY_STARTED)
	case realmerr.ErrLockHeroErr:
		hc.Send(region.ERR_CANCEL_INVASION_FAIL_SERVER_ERROR)
	default:
		logrus.WithError(err).Error("region.CancelInvasion有err没有处理")
		hc.Send(region.ERR_CANCEL_INVASION_FAIL_SERVER_ERROR)
	}
}

// 遣返
//gogen:iface
func (m *RegionModule) ProcessRepatriate(proto *region.C2SRepatriateProto, hc iface.HeroController) {
	troopId, ok := idbytes.ToId(proto.Id)
	if !ok {
		logrus.Debugf("遣返，无效的军情id")
		hc.Send(region.ERR_REPATRIATE_FAIL_ID_NOT_FOUND)
		return
	}

	realm := m.realmService.GetBigMap()
	if realm == nil {
		logrus.Error("Repatriate时, realm == nil")
		hc.Send(region.ERR_REPATRIATE_FAIL_SERVER_ERROR)
		return
	}

	processed, err := realm.Repatriate(hc, troopId)
	if !processed {
		logrus.Debugf("遣返，超时")
		hc.Send(region.ERR_REPATRIATE_FAIL_SERVER_ERROR)
		return
	}

	switch err {
	case nil:
		hc.Send(region.NewS2cRepatriateMsg(proto.Id, proto.IsTent))
	case realmerr.ErrRepatriateTroopNotFound:
		hc.Send(region.ERR_REPATRIATE_FAIL_ID_NOT_FOUND)
	case realmerr.ErrRepatriateTroopNoDefending:
		hc.Send(region.ERR_REPATRIATE_FAIL_NO_DEFENDING)
	case realmerr.ErrRepatriateAssemblyStarted:
		hc.Send(region.ERR_REPATRIATE_FAIL_ASSEMBLY_STARTED)
	case realmerr.ErrRepatriateNotAssemblyCreater:
		hc.Send(region.ERR_REPATRIATE_FAIL_NO_ASSEMBLY_CREATER)
	default:
		logrus.WithError(err).Error("region.ProcessRepatriate有err没有处理")
		hc.Send(region.ERR_REPATRIATE_FAIL_SERVER_ERROR)
	}
}

// 宝藏遣返
//gogen:iface
func (m *RegionModule) ProcessBaozRepatriate(proto *region.C2SBaozRepatriateProto, hc iface.HeroController) {

	baseId, ok := idbytes.ToId(proto.BaseId)
	if !ok {
		logrus.Debugf("宝藏遣返，无效的宝藏id")
		hc.Send(region.ERR_BAOZ_REPATRIATE_FAIL_ID_NOT_FOUND)
		return
	}

	troopId, ok := idbytes.ToId(proto.TroopId)
	if !ok {
		logrus.Debugf("宝藏遣返，无效的队伍id")
		hc.Send(region.ERR_BAOZ_REPATRIATE_FAIL_ID_NOT_FOUND)
		return
	}

	realm := m.realmService.GetBigMap()
	if realm == nil {
		logrus.Error("宝藏遣返, realm == nil")
		hc.Send(region.ERR_BAOZ_REPATRIATE_FAIL_SERVER_ERROR)
		return
	}

	processed, err := realm.BaozRepatriate(hc, baseId, troopId)
	if !processed {
		logrus.Debugf("宝藏遣返，超时")
		hc.Send(region.ERR_BAOZ_REPATRIATE_FAIL_SERVER_ERROR)
		return
	}

	switch err {
	case nil:
		hc.Send(region.NewS2cBaozRepatriateMsg(proto.BaseId, proto.TroopId))
	case realmerr.ErrBaozRepatriateTroopNotFound:
		hc.Send(region.ERR_BAOZ_REPATRIATE_FAIL_ID_NOT_FOUND)
	default:
		logrus.WithError(err).Error("宝藏遣返，有err没有处理")
		hc.Send(region.ERR_BAOZ_REPATRIATE_FAIL_SERVER_ERROR)
	}

}

// 行军加速
//gogen:iface
func (m *RegionModule) ProcessSpeedUp(proto *region.C2SSpeedUpProto, hc iface.HeroController) {
	troopId, ok := idbytes.ToId(proto.Id)
	if !ok {
		logrus.Debugf("行军加速，无效的军情id")
		hc.Send(region.ERR_SPEED_UP_FAIL_ID_NOT_FOUND)
		return
	}

	otherTroopId, ok := idbytes.ToId(proto.OtherId)
	if !ok {
		logrus.Debugf("行军加速，无效的军情id")
		hc.Send(region.ERR_SPEED_UP_FAIL_ID_NOT_FOUND)
		return
	}

	goodsData := m.datas.GetGoodsData(u64.FromInt32(proto.GoodsId))
	if goodsData == nil {
		logrus.WithField("goods", proto.GoodsId).Debug("行军加速，物品不存在")
		hc.Send(region.ERR_SPEED_UP_FAIL_INVALID_GOODS)
		return
	}

	if goodsData.EffectType != shared_proto.GoodsEffectType_EFFECT_SPEED_UP {
		logrus.WithField("goods", goodsData.Name).Debug("行军加速，物品不是行军加速物品")
		hc.Send(region.ERR_SPEED_UP_FAIL_INVALID_GOODS)
		return
	}

	hctx := heromodule.NewContext(m.dep, operate_type.RegionSpeedUp)

	var speedUpCostType uint64
	var reserveResult *entity.ReserveResult
	if hc.FuncWithSend(func(hero *entity.Hero, result herolock.LockResult) {

		var troop *entity.Troop
		if entity.GetTroopHeroId(troopId) == hc.Id() {
			troop = hero.Troop(troopId)
		} else {
			for _, t := range hero.Troops() {
				if t != nil {
					if t.Id() == troopId {
						troop = t
						break
					}

					invateInfo := t.GetInvateInfo()
					if invateInfo != nil && invateInfo.AssemblyId() == troopId {
						troop = t
						break
					}
				}
			}
		}

		if troop == nil {
			logrus.Debugf("行军加速，没找到部队id")
			result.Add(region.ERR_SPEED_UP_FAIL_ID_NOT_FOUND)
			return
		}

		invateInfo := troop.GetInvateInfo()
		if invateInfo == nil {
			logrus.Debugf("行军加速，部队没有出征")
			result.Add(region.ERR_SPEED_UP_FAIL_ID_NOT_FOUND)
			return
		}

		if otherTroopId == 0 && invateInfo.State() != realmface.AssemblyArrived && !invateInfo.State().IsMoving() {
			logrus.WithField("state", invateInfo.State()).Debugf("行军加速，部队不是行军中")
			result.Add(region.ERR_SPEED_UP_FAIL_NO_MOVING)
			return
		}

		ctime := m.time.CurrentTime()

		var hasEnoughGoods bool
		if proto.Money {
			speedUpCostType = operate_type.SpeedUpMoney
			if goodsData.YuanbaoPrice > 0 {
				if hasEnoughGoods, reserveResult = heromodule.ReserveYuanbao(hctx, hero, result, goodsData.YuanbaoPrice, ctime); !hasEnoughGoods {
					logrus.Debug("行军加速，购买，元宝不足")
					result.Add(region.ERR_SPEED_UP_FAIL_COST_NOT_ENOUGH)
					return
				}
			} else if goodsData.DianquanPrice > 0 {
				if hasEnoughGoods, reserveResult = heromodule.ReserveDianquan(hctx, hero, result, goodsData.DianquanPrice, ctime); !hasEnoughGoods {
					logrus.Debug("行军加速，购买，点券不足")
					result.Add(region.ERR_SPEED_UP_FAIL_COST_NOT_ENOUGH)
					return
				}
			} else {
				logrus.Debugf("行军加速，不支持元宝或点券购买")
				result.Add(region.ERR_SPEED_UP_FAIL_COST_NOT_SUPPORT)
				return
			}
		} else {
			speedUpCostType = operate_type.SpeedUpItem
			// 预约扣物品
			var hasEnoughGoods bool
			if hasEnoughGoods, reserveResult = heromodule.ReserveGoods(hctx, hero, result, goodsData, 1, ctime); !hasEnoughGoods {
				logrus.Debug("行军加速，预约扣除物品，个数不足")
				result.Add(region.ERR_SPEED_UP_FAIL_GOODS_NOT_ENOUGH)
				return
			}
		}

		result.Changed()
		result.Ok()
	}) {
		return
	}

	heromodule.ConfirmReserveResult(hctx, hc, reserveResult, func() (success bool) {
		realm := m.realmService.GetBigMap()
		if realm == nil {
			logrus.Error("行军加速, hero里拿出来的realmid竟然不存在")
			hc.Send(region.ERR_SPEED_UP_FAIL_SERVER_ERROR)
			return
		}

		processed, err := realm.SpeedUp(hc, troopId, otherTroopId, goodsData.GoodsEffect.TroopSpeedUpRate, speedUpCostType)
		if !processed {
			logrus.Debugf("行军加速，超时")
			hc.Send(region.ERR_SPEED_UP_FAIL_SERVER_ERROR)
			return
		}

		switch err {
		case nil:
			success = true
			hc.Send(region.NewS2cSpeedUpMsg(proto.Id))
		case realmerr.ErrSpeedUpTroopNotFound:
			hc.Send(region.ERR_SPEED_UP_FAIL_ID_NOT_FOUND)
		case realmerr.ErrSpeedUpTroopNoMoving:
			hc.Send(region.ERR_SPEED_UP_FAIL_NO_MOVING)
		case realmerr.ErrSpeedUpOtherTroopNotFound:
			hc.Send(region.ERR_SPEED_UP_FAIL_OTHER_ID_NOT_FOUND)
		case realmerr.ErrSpeedUpAssemblyWait:
			hc.Send(region.ERR_SPEED_UP_FAIL_ASSEMBLY_WAIT)
		default:
			logrus.WithError(err).Error("region.ProcessSpeedUp有err没有处理")
			hc.Send(region.ERR_SPEED_UP_FAIL_SERVER_ERROR)
		}

		return
	})

}

//gogen:iface
func (m *RegionModule) ProcessExpel(proto *region.C2SExpelProto, hc iface.HeroController) {

	troopId, ok := idbytes.ToId(proto.Id)
	if !ok {
		logrus.Debugf("驱逐，无效的军情id")
		hc.Send(region.ERR_EXPEL_FAIL_INVALID_ID)
		return
	}

	realmId := int64(proto.Mapid)

	troopIndex := u64.FromInt32(proto.TroopIndex - 1)

	// 驱逐CD中
	ctime := m.time.CurrentTime()
	if hc.FuncWithSend(func(hero *entity.Hero, result herolock.LockResult) {
		if ctime.Before(hero.Military().NextExpelTime()) {
			logrus.Debugf("驱逐，CD中")
			result.Add(region.ERR_EXPEL_FAIL_COOLDOWN)
			return
		}

		t := hero.GetTroopByIndex(troopIndex)
		if t == nil {
			logrus.Debugf("驱逐，无效的部队编号")
			result.Add(region.ERR_EXPEL_FAIL_INVALID_TROOP_INDEX)
			return
		}

		if t.IsOutside() {
			logrus.Debugf("驱逐，队伍出征中")
			result.Add(region.ERR_EXPEL_FAIL_OUTSIDE)
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
			logrus.Debugf("出征，队伍没有士兵")
			hc.Send(region.ERR_EXPEL_FAIL_NO_SOLDIER)
			return
		}

		if realmId != hero.BaseRegion() {
			logrus.Debugf("驱逐，地图id无效（必须是主城地图或行营地图id）")
			result.Add(region.ERR_EXPEL_FAIL_INVALID_MAP)
			return
		}

		result.Ok()
	}) {
		return
	}

	realm := m.realmService.GetBigMap()
	if realm == nil {
		logrus.WithField("realmid", realmId).Error("驱逐, hero里拿出来的部队所在realmid竟然不存在")
		hc.Send(region.ERR_EXPEL_FAIL_SERVER_ERROR)
		return
	}

	processed, success, link, err := realm.Expel(hc, troopId, troopIndex)
	if !processed {
		logrus.Errorf("驱逐，timeout")
		hc.Send(region.ERR_EXPEL_FAIL_SERVER_ERROR)
		return
	}

	switch err {
	case nil:
		var nextExpelTime int32
		if !success {
			hc.FuncNotError(func(hero *entity.Hero) (heroChanged bool) {
				// 驱逐失败了，进入CD
				toSet := ctime.Add(m.datas.RegionConfig().ExpelDuration)
				hero.Military().SetNextExpelTime(toSet)
				nextExpelTime = timeutil.Marshal32(toSet)

				return true
			})
		}

		hc.Send(region.NewS2cExpelMsg(proto.Id, nextExpelTime, link))

		if success {
			hc.Send(misc.NewS2cScreenShowWordsMsg(
				m.datas.TextHelp().RealmExpelledSuccess.New().
					WithTroopIndex(troopIndex + 1).
					JsonString()))
		} else {
			hc.Send(misc.NewS2cScreenShowWordsMsg(
				m.datas.TextHelp().RealmExpelledFail.New().
					WithTroopIndex(troopIndex + 1).
					JsonString()))
		}

	case realmerr.ErrExpelSelfNoBase:
		hc.Send(region.ERR_EXPEL_FAIL_INVALID_MAP)
	case realmerr.ErrExpelTroopsNotFound:
		hc.Send(region.ERR_EXPEL_FAIL_NOT_SELF)
	case realmerr.ErrExpelTroopsNoRobbing:
		hc.Send(region.ERR_EXPEL_FAIL_NOT_ARRIVED)
	case realmerr.ErrExpelFightError:
		hc.Send(region.ERR_EXPEL_FAIL_SERVER_ERROR)
	case realmerr.ErrExpelInvalidTroopIndex:
		hc.Send(region.ERR_EXPEL_FAIL_INVALID_TROOP_INDEX)
	case realmerr.ErrExpelCaptainOutside:
		hc.Send(region.ERR_EXPEL_FAIL_OUTSIDE)
	case realmerr.ErrExpelNoSoldier:
		hc.Send(region.ERR_EXPEL_FAIL_NO_SOLDIER)
	default:
		logrus.WithError(err).Error("region.Expel有err没有处理")
		hc.Send(region.ERR_EXPEL_FAIL_SERVER_ERROR)
	}
}

//gogen:iface
func (m *RegionModule) ProcessGetWhiteFlagDetail(proto *region.C2SWhiteFlagDetailProto, hc iface.HeroController) {

	heroId := hc.Id()
	if len(proto.HeroId) > 0 {
		id, ok := idbytes.ToId(proto.HeroId)
		if !ok {
			logrus.Debugf("查看玩家白旗详情，id无效")
			hc.Send(region.ERR_WHITE_FLAG_DETAIL_FAIL_NO_FLAG)
			return
		}
		heroId = id
	}

	var whiteFlagHeroId, whiteFlagGuildId int64
	var disappearTime time.Time
	var errMsg pbutil.Buffer
	m.heroDataServive.Func(heroId, func(hero *entity.Hero, err error) (heroChanged bool) {

		if err != nil {
			logrus.WithError(err).Debugf("查看玩家白旗详情，lock hero 失败")
			errMsg = region.ERR_WHITE_FLAG_DETAIL_FAIL_NO_FLAG
			return
		}

		if hero.GetWhiteFlagGuildId() == 0 || hero.GetWhiteFlagHeroId() == 0 {
			logrus.Debugf("查看玩家白旗详情，玩家当前没有白旗")
			errMsg = region.ERR_WHITE_FLAG_DETAIL_FAIL_NO_FLAG
			return
		}

		ctime := m.time.CurrentTime()
		disappearTime = hero.GetWhiteFlagDisappearTime()
		if ctime.After(disappearTime) {
			logrus.Debugf("查看玩家白旗详情，玩家当前没有白旗")
			errMsg = region.ERR_WHITE_FLAG_DETAIL_FAIL_NO_FLAG
			return
		}

		whiteFlagHeroId = hero.GetWhiteFlagHeroId()
		whiteFlagGuildId = hero.GetWhiteFlagGuildId()

		return
	})

	if errMsg != nil {
		hc.Send(errMsg)
		return
	}

	g := m.guildSnapshotService.GetSnapshot(whiteFlagGuildId)
	if g == nil {
		logrus.Debugf("查看玩家白旗详情，获取帮派快照失败")
		hc.Send(region.ERR_WHITE_FLAG_DETAIL_FAIL_NO_FLAG)
		return
	}

	target := m.heroSnapshotServive.Get(whiteFlagHeroId)
	if target == nil {
		logrus.Debugf("查看玩家白旗详情，获取插旗英雄快照失败")
		hc.Send(region.ERR_WHITE_FLAG_DETAIL_FAIL_NO_FLAG)
		return
	}

	var countryId int32
	if country := g.Country; country != nil {
		countryId = u64.Int32(country.Id)
	}

	hc.Send(region.NewS2cWhiteFlagDetailMsg(idbytes.ToBytes(heroId), target.IdBytes, target.Name, i64.Int32(g.Id), g.Name, timeutil.Marshal32(disappearTime), countryId))
}

//gogen:iface
func (m *RegionModule) ProcessUseMianGoods(proto *region.C2SUseMianGoodsProto, hc iface.HeroController) {
	m.UseMianGoods(u64.FromInt32(proto.Id), proto.Buy, hc)
}

func (m *RegionModule) UseMianGoods(goodsId uint64, buy bool, hc iface.HeroController) (succ bool) {

	// 使用免战物品

	goodsData := m.datas.GetGoodsData(goodsId)
	if goodsData == nil {
		logrus.Debug("使用免战物品，无效的物品id")
		hc.Send(region.ERR_USE_MIAN_GOODS_FAIL_INVALID_ID)
		return
	}

	if goodsData.GoodsEffect == nil || goodsData.GoodsEffect.MianDuration <= 0 {
		logrus.Debug("使用免战物品，不是免战物品")
		hc.Send(region.ERR_USE_MIAN_GOODS_FAIL_INVALID_GOODS)
		return
	}

	ctime := m.time.CurrentTime()
	disappearTime := ctime.Add(goodsData.GoodsEffect.MianDuration)
	hctx := heromodule.NewContext(m.dep, operate_type.RegionUseMianGoods)

	var baseRegion int64
	var reserveResult *entity.ReserveResult
	if hc.FuncWithSend(func(hero *entity.Hero, result herolock.LockResult) {
		if !buy && hero.Depot().GetGoodsCount(goodsData.Id) <= 0 {
			logrus.Debug("使用免战物品，物品个数不足")
			result.Add(region.ERR_USE_MIAN_GOODS_FAIL_COUNT_NOT_ENOUGH)
			return
		}

		// 免战中
		if disappearTime.Before(hero.GetMianDisappearTime()) {
			logrus.Debug("使用免战物品，当前免战中")
			result.Add(region.ERR_USE_MIAN_GOODS_FAIL_MIAN)
			return
		}

		// CD
		if ctime.Before(hero.GetNextUseMianGoodsTime()) {
			logrus.Debug("使用免战物品，免战物品CD中")
			result.Add(region.ERR_USE_MIAN_GOODS_FAIL_COOLDOWN)
			return
		}

		if hero.BaseRegion() == 0 || hero.BaseLevel() <= 0 || hero.Prosperity() <= 0 {
			logrus.Debug("使用免战物品，主城流亡了")
			result.Add(region.ERR_USE_MIAN_GOODS_FAIL_HOME_NOT_ALIVE)
			return
		}

		// 部队出征中（如果是主城出去的，打怪的可以，帮人的可以，其他都不可以...）
		for _, t := range hero.Troops() {
			if v := t.GetInvateInfo(); v != nil {
				if npcid.IsNpcId(v.OriginTargetID()) {
					continue
				}

				//if v.RegionID() != hero.BaseRegion() {
				//	// 不是主城出去的，不可以（行营都没有，怎么出去...）
				//	logrus.
				//		WithField("base_region", hero.BaseRegion()).
				//		WithField("troop_region", v.RegionID()).
				//		WithField("target", v.TargetBaseID()).
				//		Error("使用免战物品，行营在家的情况，居然有部队不在主城所在的地图")
				//	result.Add(region.ERR_USE_MIAN_GOODS_FAIL_TROOP_OUTSIDE)
				//	return
				//}

				switch v.State() {
				case realmface.MovingToAssist,
					realmface.Defending,
					realmface.AssistMovingBack:
					// 主城出去援助的可以，其他的不行
				default:
					logrus.WithField("state", v.State()).Debug("使用免战物品，部队出征中")
					result.Add(region.ERR_USE_MIAN_GOODS_FAIL_TROOP_OUTSIDE)
					return
				}
			}
		}

		var ok bool
		if ok, reserveResult = heromodule.ReserveGoodsOrBuy(hctx, hero, result, goodsData, 1, buy, ctime); !ok {
			logrus.Debug("使用免战物品，预约扣除购买物品，扣除失败")
			result.Add(region.ERR_USE_MIAN_GOODS_FAIL_COUNT_NOT_ENOUGH)
			return
		}

		result.Ok()
	}) {
		return
	}

	heromodule.ConfirmReserveResultWithHeroFunc(hctx, hc, reserveResult, func() (success bool, heroFunc herolock.SendFunc) {
		realm := m.realmService.GetBigMap()
		if realm == nil {
			logrus.WithField("base_region", baseRegion).Error("使用免战物品，主城所在的realm不存在")
			hc.Send(region.ERR_USE_MIAN_GOODS_FAIL_HOME_NOT_ALIVE)
			return
		}

		// 加免战结束时间
		processed, err := realm.Mian(hc.Id(), disappearTime, true)
		if !processed {
			logrus.WithField("base_region", baseRegion).Error("使用免战物品，主城所在的realm不存在")
			hc.Send(region.ERR_USE_MIAN_GOODS_FAIL_SERVER_ERROR)
			return
		}

		switch err {
		case nil:
			// 成功（CD要检查，但是免战CD在破除免战的时候加，所以这里不加免战CD）
			//nextUseTime := disappearTime.Add(goodsData.Cd)
			//hc.Send(region.NewS2cUseMianGoodsMsg(proto.Id, timeutil.Marshal32(nextUseTime)))
			//heroFunc = func(hero *entity.Hero, result herolock.LockResult) {
			//	hero.SetNextUseMianGoodsTime(nextUseTime)
			//}

			hc.Send(region.NewS2cUseMianGoodsMsg(u64.Int32(goodsId), 0))

			success = true
			succ = true
		case realmerr.ErrMianSelfNoBase,
			realmerr.ErrMianTent:
			hc.Send(region.ERR_USE_MIAN_GOODS_FAIL_HOME_NOT_ALIVE)
		case realmerr.ErrMianExist,
			realmerr.ErrMianCantOverwrite:
			hc.Send(region.ERR_USE_MIAN_GOODS_FAIL_MIAN)
		default:
			logrus.WithError(err).Error("使用免战物品，存在位置的错误类型")
			hc.Send(region.ERR_USE_MIAN_GOODS_FAIL_SERVER_ERROR)
			return
		}

		return
	})

	return
}

//gogen:iface c2s_favorite_pos
func (m *RegionModule) ProcessFavoritePos(proto *region.C2SFavoritePosProto, hc iface.HeroController) {
	hc.FuncWithSend(func(hero *entity.Hero, result herolock.LockResult) {
		favoritePoses := hero.FavoritePoses()
		if proto.Add {
			realm := m.realmService.GetBigMap()
			if realm == nil {
				logrus.Debugf("添加收藏点失败，没找到场景")
				result.Add(region.ERR_FAVORITE_POS_FAIL_SCENE_NOT_FOUND)
				return
			}

			if uint64(proto.PosX) >= realm.GetMapData().XLen {
				logrus.Debugf("添加收藏点失败，X越界了")
				result.Add(region.ERR_FAVORITE_POS_FAIL_POS_INVALID)
				return
			}

			if uint64(proto.PosY) >= realm.GetMapData().YLen {
				logrus.Debugf("添加收藏点失败，Y越界了")
				result.Add(region.ERR_FAVORITE_POS_FAIL_POS_INVALID)
				return
			}

			full, exist := favoritePoses.Add(proto.Id, proto.PosX, proto.PosY)
			if full {
				logrus.Debugf("添加收藏点失败，收藏点列表已满")
				result.Add(region.ERR_FAVORITE_POS_FAIL_FULL)
				return
			}

			if exist {
				logrus.Debugf("添加收藏点失败，收藏点已经存在了")
				result.Add(region.ERR_FAVORITE_POS_FAIL_EXIST)
				return
			}
		} else {
			if !favoritePoses.Del(proto.Id, proto.PosX, proto.PosY) {
				// 没找到，删除失败
				logrus.Debugf("删除收藏点失败，没找到")
				result.Add(region.ERR_FAVORITE_POS_FAIL_NOT_FOUND)
				return
			}
		}

		result.Add(region.NewS2cFavoritePosMsg(proto.Add, proto.Id, proto.PosX, proto.PosY))
		result.Changed()
		result.Ok()
	})
}

//gogen:iface c2s_favorite_pos_list
func (m *RegionModule) ProcessFavoritePosList(hc iface.HeroController) {
	listProto := &shared_proto.FavoritePosListProto{}

	hc.FuncWithSend(func(hero *entity.Hero, result herolock.LockResult) {
		favoritePoses := hero.FavoritePoses()

		listProto.Detail = make([]*shared_proto.FavoritePosDetailProto, 0, favoritePoses.PosCount())
		favoritePoses.Walk(func(pos *shared_proto.FavoritePosProto) {
			detail := &shared_proto.FavoritePosDetailProto{
				Pos: pos,
			}

			if int64(pos.Id) == hero.BaseRegion() {
				// 检查是否有npc
				hero.WalkHomeNpcBase(func(base *entity.HomeNpcBase) bool {
					basePos := hexagon.ShiftEvenOffset(hero.BaseX(), hero.BaseY(), base.GetData().EvenOffsetX, base.GetData().EvenOffsetY)
					if basePos == cb.XYCube(int(pos.X), int(pos.Y)) {
						detail.HaveHero = true
						detail.Hero = base.GetData().Data.Npc.EncodeHeroBasicProto(base.IdBytes())
					}

					return false
				})
			}

			listProto.Detail = append(listProto.Detail, detail)
		})
	})

	for _, detail := range listProto.Detail {
		if detail.HaveHero {
			continue
		}

		realm := m.realmService.GetBigMap()
		if realm == nil {
			continue
		}

		info := realm.GetRoBaseByPos(int(detail.Pos.X), int(detail.Pos.Y))
		if info == nil {
			// 没人
			continue
		}

		if info.GetBaseType() == realmface.BaseTypeNpc {
			switch info.GetNpcType() {
			case npcid.NpcType_MultiLevelMonster:
				// npc
				npcBaseData := m.datas.GetRegionMultiLevelNpcData(info.GetNpcDataId())
				if npcBaseData != nil {
					detail.HaveHero = true
					detail.Hero = npcBaseData.GetFirstLevel().Npc.Npc.EncodeHeroBasicProto(idbytes.ToBytes(info.GetId()))
				} else {
					logrus.WithField("id", info.GetId()).Errorln("请求玩家收藏列表时，竟然没找到npc的信息，难道存的不是monster而是homeNpc?")
				}
			case npcid.NpcType_XiongNu:
				xiongNuData := m.datas.GetResistXiongNuData(info.GetNpcDataId())
				if xiongNuData != nil {
					detail.HaveHero = true
					detail.Hero = xiongNuData.NpcBaseData.Npc.EncodeHeroBasicProto(idbytes.ToBytes(info.GetId()))
				} else {
					logrus.WithField("id", info.GetId()).Errorln("请求玩家收藏列表时，竟然没找到npc的信息，难道存的不是monster而是homeNpc?")
				}
			case npcid.NpcType_BaoZang:
				data := m.datas.GetBaozNpcData(info.GetNpcDataId())
				if data != nil {
					detail.HaveHero = true
					detail.Hero = data.Npc.Npc.EncodeHeroBasicProto(idbytes.ToBytes(info.GetId()))
				} else {
					logrus.WithField("id", info.GetId()).Errorln("请求玩家收藏列表时，竟然没找到npc的信息，难道存的不是monster而是homeNpc?")
				}
			default:
				// npc
				npcBaseData := m.datas.GetNpcBaseData(info.GetNpcDataId())
				if npcBaseData != nil {
					detail.HaveHero = true
					detail.Hero = npcBaseData.Npc.EncodeHeroBasicProto(idbytes.ToBytes(info.GetId()))
				} else {
					logrus.WithField("id", info.GetId()).Errorln("请求玩家收藏列表时，竟然没找到npc的信息，难道存的不是monster而是homeNpc?")
				}
			}

			continue
		}

		snapshot := m.heroSnapshotServive.Get(info.GetId())
		if snapshot == nil {
			logrus.Debugf("请求玩家收藏列表，没找到玩家的snapshot")
			continue
		}

		detail.HaveHero = true
		detail.Hero = snapshot.EncodeBasic4Client()
		detail.IsHome = info.GetBaseType() == realmface.BaseTypeHome
	}

	hc.Send(region.NewS2cFavoritePosListMsg(must.Marshal(listProto)))
}

var buyProsperity0CostMsg = region.NewS2cGetBuyProsperityCostMsg(0).Static()

//gogen:iface c2s_get_buy_prosperity_cost
func (m *RegionModule) ProcessGetBuyProsperityCost(hc iface.HeroController) {
	hc.FuncWithSend(func(hero *entity.Hero, result herolock.LockResult) {
		if hero.BaseRegion() == 0 || hero.BaseLevel() <= 0 || hero.Prosperity() <= 0 {
			result.Add(buyProsperity0CostMsg)
			return
		}

		guanFu := hero.Domestic().GetBuilding(shared_proto.BuildingType_GUAN_FU)
		if guanFu == nil {
			logrus.Errorln("获取购买繁荣度消耗，没有官府")
			result.Add(buyProsperity0CostMsg)
			return
		}

		guanFuLevelData := m.datas.GetGuanFuLevelData(guanFu.Level)
		if guanFuLevelData == nil {
			logrus.Errorln("获取购买繁荣度消耗，没有官府等级")
			result.Add(buyProsperity0CostMsg)
			return
		}

		buyProsperity := u64.Sub(hero.ProsperityCapcity(), hero.Prosperity())
		if buyProsperity <= 0 {
			result.Add(buyProsperity0CostMsg)
			return
		}

		costAmount := calculateBuyProsperityCost(buyProsperity, guanFuLevelData.RestoreProsperity,
			m.datas.RegionConfig().RestoreHomeProsperityDuration,
			guanFuLevelData.BuyProsperityRestoreDurationWith1Cost)

		result.Add(region.NewS2cGetBuyProsperityCostMsg(u64.Int32(costAmount)))

		result.Ok()
	})
}

//gogen:iface c2s_buy_prosperity
func (m *RegionModule) ProcessBuyProsperity(hc iface.HeroController) {
	ctime := m.time.CurrentTime()
	var buyProsperity uint64

	hctx := heromodule.NewContext(m.dep, operate_type.RegionBuyProsperity)

	var baseRegion int64
	var reserveResult *entity.ReserveResult
	if hc.FuncWithSend(func(hero *entity.Hero, result herolock.LockResult) {
		if hero.BaseRegion() == 0 || hero.BaseLevel() <= 0 || hero.Prosperity() <= 0 {
			logrus.Debug("购买繁荣度，主城流亡了")
			result.Add(region.ERR_BUY_PROSPERITY_FAIL_HOME_NOT_ALIVE)
			return
		}

		guanFu := hero.Domestic().GetBuilding(shared_proto.BuildingType_GUAN_FU)
		if guanFu == nil {
			logrus.Errorln("购买繁荣度，没有官府")
			result.Add(region.ERR_BUY_PROSPERITY_FAIL_SERVER_ERROR)
			return
		}

		guanFuLevelData := m.datas.GetGuanFuLevelData(guanFu.Level)
		if guanFuLevelData == nil {
			logrus.Errorln("购买繁荣度，没有找到官府等级")
			result.Add(region.ERR_BUY_PROSPERITY_FAIL_SERVER_ERROR)
			return
		}

		buyProsperity = u64.Sub(hero.ProsperityCapcity(), hero.Prosperity())
		if buyProsperity <= 0 {
			logrus.Debug("购买繁荣度，繁荣度已满")
			result.Add(region.ERR_BUY_PROSPERITY_FAIL_PROSPERITY_FULL)
			return
		}

		if vipData := m.datas.GetVipLevelData(hero.VipLevel()); vipData == nil || !vipData.BuyProsperity {
			logrus.Debug("购买繁荣度，vip等级不够")
			result.Add(region.ERR_BUY_PROSPERITY_FAIL_VIP_LEVEL_LIMIT)
			return
		}

		costAmount := calculateBuyProsperityCost(buyProsperity, guanFuLevelData.RestoreProsperity,
			m.datas.RegionConfig().RestoreHomeProsperityDuration,
			guanFuLevelData.BuyProsperityRestoreDurationWith1Cost)
		if costAmount <= 0 {
			logrus.Debug("购买繁荣度，计算消耗结果<=0")
			result.Add(region.ERR_BUY_PROSPERITY_FAIL_SERVER_ERROR)
			return
		}

		var hasEnough bool
		if hasEnough, reserveResult = heromodule.ReserveYuanbao(hctx, hero, result, costAmount, ctime); !hasEnough {
			logrus.Debug("购买繁荣度，预约扣除消耗不足")
			result.Add(region.ERR_BUY_PROSPERITY_FAIL_COST_NOT_ENOUGH)
			return
		}

		baseRegion = hero.BaseRegion()

		result.Ok()
	}) {
		return
	}

	heromodule.ConfirmReserveResult(hctx, hc, reserveResult, func() (success bool) {
		realm := m.realmService.GetBigMap()
		if realm == nil {
			logrus.WithField("base_region", baseRegion).Error("购买繁荣度，主城所在的realm不存在")
			hc.Send(region.ERR_BUY_PROSPERITY_FAIL_HOME_NOT_ALIVE)
			return
		}

		// 加繁荣度
		processed, err := realm.AddProsperity(hc.Id(), buyProsperity)
		if !processed {
			logrus.WithField("base_region", baseRegion).Error("购买繁荣度，主城所在的realm不存在")
			hc.Send(region.ERR_BUY_PROSPERITY_FAIL_SERVER_ERROR)
			return
		}

		switch err {
		case nil:
			// 成功
			hc.Send(region.NewS2cBuyProsperityMsg(u64.Int32(buyProsperity)))
			success = true
		case realmerr.ErrAddProsperitySelfNoBase:
			hc.Send(region.ERR_BUY_PROSPERITY_FAIL_HOME_NOT_ALIVE)
		default:
			logrus.WithError(err).Error("购买繁荣度，存在未知的错误类型")
			hc.Send(region.ERR_BUY_PROSPERITY_FAIL_SERVER_ERROR)
			return
		}

		return
	})

}

func calculateBuyProsperityCost(buyAmount, restoreAmountPerTimes uint64, restoreDurationPerTimes, buyDurationPerCost time.Duration) uint64 {

	// 计算当前需要多少个消耗回复间隔
	multi := u64.Division2Float64(buyAmount, restoreAmountPerTimes)

	if buyDurationPerCost > 0 && restoreDurationPerTimes != buyDurationPerCost {
		// 计算总共所需的回复时间
		multi = multi * float64(restoreDurationPerTimes)

		// 除于每单位价格购买时间，得到花费价格
		multi = multi / float64(buyDurationPerCost)
	}

	return u64.FromInt64(int64(math.Ceil(multi)))
}

//gogen:iface c2s_request_ruins_base
func (m *RegionModule) ProcessRequestRuinsBase(proto *region.C2SRequestRuinsBaseProto, hc iface.HeroController) {
	realm := m.realmService.GetBigMap()
	if realm == nil {
		logrus.WithField("realm", proto.RealmId).Debugln("请求废墟信息，场景没找到")
		hc.Send(region.ERR_REQUEST_RUINS_BASE_FAIL_REALM_NOT_FOUND)
		return
	}

	if proto.PosX < 0 || uint64(proto.PosX) >= realm.GetMapData().XLen {
		logrus.WithField("x", proto.PosX).Debugln("请求废墟信息，x坐标非法")
		hc.Send(region.ERR_REQUEST_RUINS_BASE_FAIL_INVALID_X_OR_Y)
		return
	}

	if proto.PosY < 0 || uint64(proto.PosY) >= realm.GetMapData().YLen {
		logrus.WithField("y", proto.PosY).Debugln("请求废墟信息，y坐标非法")
		hc.Send(region.ERR_REQUEST_RUINS_BASE_FAIL_INVALID_X_OR_Y)
		return
	}

	heroId := realm.GetRuinsBase(int(proto.PosX), int(proto.PosY))
	if heroId == 0 {
		logrus.Debugln("请求废墟信息，没找到废墟")
		hc.Send(region.ERR_REQUEST_RUINS_BASE_FAIL_NO_RUINS)
		return
	}

	snapshot := m.heroSnapshotServive.Get(heroId)
	if snapshot == nil {
		logrus.WithField("id", heroId).Errorln("请求废墟信息，没load到玩家信息")
		hc.Send(region.ERR_REQUEST_RUINS_BASE_FAIL_SERVER_ERROR)
		return
	}

	hc.Send(region.NewS2cRequestRuinsBaseMsg(proto.RealmId, proto.PosX, proto.PosY, must.Marshal(snapshot.EncodeBasic4Client())))
}

//gogen:iface
func (m *RegionModule) ProcessGetPrevInvestigate(proto *region.C2SGetPrevInvestigateProto, hc iface.HeroController) {

	targetId, ok := idbytes.ToId(proto.HeroId)
	if !ok {
		hc.Send(region.NewS2cGetPrevInvestigateMsg(proto.HeroId, nil, 0))
		return
	}

	ctime := m.time.CurrentTime()
	hc.FuncWithSend(func(hero *entity.Hero, result herolock.LockResult) {
		if inv := hero.GetInvestigation(targetId); inv != nil && ctime.Before(inv.GetExpireTime()) {
			hc.Send(region.NewS2cGetPrevInvestigateMsg(proto.HeroId, i64.ToBytesU64(inv.GetMailId()), timeutil.Marshal32(inv.GetExpireTime())))
		} else {
			hc.Send(region.NewS2cGetPrevInvestigateMsg(proto.HeroId, nil, 0))
		}
		result.Ok()
	})

}

//gogen:iface
func (m *RegionModule) ProcessInvestigate(proto *region.C2SInvestigateProto, hc iface.HeroController) {

	targetId, ok := idbytes.ToId(proto.HeroId)
	if !ok {
		logrus.WithField("id", targetId).Debug("侦查城池，玩家id无效")
		hc.Send(region.ERR_INVESTIGATE_FAIL_HERO_NOT_FOUND)
		return
	}

	if hc.Id() == targetId {
		logrus.Debug("侦查城池，不能侦查自己")
		hc.Send(region.ERR_INVESTIGATE_FAIL_SELF_ID)
		return
	}

	if npcid.IsNpcId(targetId) {
		logrus.Debug("侦查城池，目标是Npc")
		hc.Send(region.ERR_INVESTIGATE_FAIL_HERO_NOT_FOUND)
		return
	}

	// 名城战期间，不能集结
	if _, joined := m.mingcWarService.JoiningFightMingc(hc.Id()); joined {
		hc.Send(region.ERR_INVESTIGATE_FAIL_IN_MC_WAR_FIGHT)
		return
	}

	ctime := m.time.CurrentTime()
	// 检查CD
	var selfGuildId int64
	var selfBaseX, selfBaseY int
	var actProsperityCapcity uint64
	var attacker *shared_proto.ReportHeroProto
	selfIsMian := false
	if hc.FuncWithSend(func(hero *entity.Hero, result herolock.LockResult) {
		selfIsMian = ctime.Before(hero.GetMianDisappearTime())
		// 现在不判断免战状态，免战的话，直接去除免战
		//if ctime.Before(hero.GetMianDisappearTime()) {
		//	logrus.Debug("侦查城池，自己免战中")
		//	result.Add(region.ERR_INVESTIGATE_FAIL_SELF_MIAN)
		//	return
		//}

		// 自己处于新手免战，不能瞭望别人
		if ctime.Before(hero.GetNewHeroMianDisappearTime()) &&
			hero.GetNewHeroMianDisappearTime().Equal(hero.GetMianDisappearTime()) {
			logrus.Debugf("侦查城池，新手免战期间不能侦查玩家")
			result.Add(region.ERR_INVESTIGATE_FAIL_SELF_MIAN)
			return
		}

		selfGuildId = hero.GuildId()
		selfBaseX, selfBaseY = hero.BaseX(), hero.BaseY()
		if ctime.Before(hero.GetNextInvestigateTime()) {
			if proto.Cost {
				// 检查钱够不够
				if !heromodule.HasEnoughCost(hero, m.datas.RegionConfig().MiaoInvestigateCdCost) {
					logrus.Debug("侦查城池，有CD，消耗不足")
					result.Add(region.ERR_INVESTIGATE_FAIL_COST_NOT_ENOUGH)
					return
				}
			} else {
				logrus.Debug("侦查城池，侦查CD中")
				result.Add(region.ERR_INVESTIGATE_FAIL_COOLDOWN)
				return
			}
		}

		actProsperityCapcity = hero.ProsperityCapcity()

		attacker = &shared_proto.ReportHeroProto{
			Id:                hero.IdBytes(),
			Name:              hero.Name(),
			Level:             int32(hero.Level()),
			Head:              hero.Head(),
			BaseRegion:        i64.Int32(hero.BaseRegion()),
			BaseX:             imath.Int32(hero.BaseX()),
			BaseY:             imath.Int32(hero.BaseY()),
			Prosperity:        u64.Int32(hero.Prosperity()),
			ProsperityCapcity: u64.Int32(hero.ProsperityCapcity()),
			Country:           u64.Int32(hero.CountryId()),
		}

		result.Ok()
	}) {
		return
	}

	report := &shared_proto.FightReportProto{}

	hctx := heromodule.NewContext(m.dep, operate_type.RegionInvestigate)

	var toSend pbutil.Buffer
	var defenserGuildId int64

	if m.heroDataServive.FuncWithSend(targetId, func(hero *entity.Hero, result herolock.LockResult) {
		if selfGuildId != 0 && selfGuildId == hero.GuildId() {
			logrus.Debug("侦查城池，不能侦查盟友")
			toSend = region.ERR_INVESTIGATE_FAIL_SAME_GUILD
			return
		}

		if !util.IsInRange(selfBaseX, selfBaseY, hero.BaseX(), hero.BaseY(),
			int(m.datas.RegionConfig().InvestigateMaxDistance)) {
			logrus.Debug("侦查城池，超出最大距离限制")
			toSend = region.ERR_INVESTIGATE_FAIL_DISTANCE
			return
		}

		if ctime.Before(hero.GetMianDisappearTime()) {
			logrus.Debug("侦查城池，目标免战中")
			toSend = region.ERR_INVESTIGATE_FAIL_TARGET_MIAN
			return
		}

		// 收税收
		heromodule.TryUpdateTax(hero, result, ctime, m.datas.MiscGenConfig().TaxDuration, m.datas.GetBuffEffectData)

		defenserGuildId = hero.GuildId()

		unsafe := hero.GetUnsafeResource()

		report.AttackerSide = true
		report.Defenser = &shared_proto.ReportHeroProto{
			Id:                hero.IdBytes(),
			Name:              hero.Name(),
			Level:             int32(hero.Level()),
			Head:              hero.Head(),
			BaseRegion:        i64.Int32(hero.BaseRegion()),
			BaseX:             imath.Int32(hero.BaseX()),
			BaseY:             imath.Int32(hero.BaseY()),
			Prosperity:        u64.Int32(hero.Prosperity()),
			ProsperityCapcity: u64.Int32(hero.ProsperityCapcity()),
			Country:           u64.Int32(hero.CountryId()),
		}

		report.Defenser.BaseLevel = u64.Int32(hero.BaseLevel())

		if bd := hero.Domestic().GetBuilding(shared_proto.BuildingType_CHENG_QIANG); bd != nil {
			report.Defenser.WallLevel = u64.Int32(bd.Level)
		}

		// 防守阵容
		if t := hero.GetHomeDefenser(); t != nil && !t.IsOutside() && t.HasSoldier() {
			for _, pos := range t.Pos() {
				c := pos.Captain()
				race := shared_proto.Race_InvalidRace
				if c != nil {
					race = c.Race().Race
				}
				report.Defenser.Race = append(report.Defenser.Race, race)
			}

			report.Defenser.TotalFightAmount = u64.Int32(t.CalDefenseFightAmount(hero))
		} else if copyDefenser := hero.GetCopyDefenser(); copyDefenser != nil {

			var captainFightAmounts []uint64
			for _, cis := range copyDefenser.GetCaptains() {
				race := shared_proto.Race_InvalidRace
				if cis != nil {
					c := hero.Military().Captain(cis.GetId())
					if c != nil {
						race = c.Race().Race

						if cis.GetSoldier() > 0 {
							captainFightAmounts = append(captainFightAmounts, cis.GetFightAmount())
						}
					}
				}
				report.Defenser.Race = append(report.Defenser.Race, race)
			}

			report.Defenser.TotalFightAmount = u64.Int32(data.TroopFightAmount(captainFightAmounts...))
		}

		// 衰减系数= arctan（（M守－M攻）/（M守＋M攻））/ π ＋1    π为圆周率
		defProsperityCapcity := hero.ProsperityCapcity()
		weakCoef := math.Atan2(
			u64.Sub2Float64(defProsperityCapcity, actProsperityCapcity),
			float64(defProsperityCapcity+actProsperityCapcity),
		)/math.Pi + 1
		//systemCoef := math.Min(m.datas.RegionConfig().RobberCoef, 1) // 系统损耗

		// 可掠夺资源 =（该资源当前储量-仓库保护值）*衰减系数*系统损耗
		coef := weakCoef //* systemCoef
		report.ShowPrize = &shared_proto.PrizeProto{}
		report.ShowPrize.Gold = i32.MultiF64(u64.Sub(unsafe.Gold(), hero.BuildingEffect().ProtectedCapcity()), coef)
		report.ShowPrize.Food = i32.MultiF64(u64.Sub(unsafe.Food(), hero.BuildingEffect().ProtectedCapcity()), coef)
		report.ShowPrize.Wood = i32.MultiF64(u64.Sub(unsafe.Wood(), hero.BuildingEffect().ProtectedCapcity()), coef)
		report.ShowPrize.Stone = i32.MultiF64(u64.Sub(unsafe.Stone(), hero.BuildingEffect().ProtectedCapcity()), coef)

		// 宝物
		topN := m.datas.RegionConfig().InvestigationBaowuCount
		if topN > 0 {

			top := sortkeys.NewU64TopN(topN)

			hero.Depot().RangeBaowu(func(id, count uint64) (toContinue bool) {
				if count <= 0 {
					return true
				}

				data := m.datas.GetBaowuData(id)
				if data == nil {
					return true
				}

				if data.CantRob {
					return true
				}

				top.Add(data.Level, data)

				return true
			})

			for _, kv := range top.SortDesc() {
				data := kv.V.(*resdata.BaowuData)

				count := hero.Depot().GetBaowuCount(data.Id)
				report.ShowPrize.BaowuId = append(report.ShowPrize.BaowuId, u64.Int32(data.Id))
				report.ShowPrize.BaowuCount = append(report.ShowPrize.BaowuCount, u64.Int32(count))
			}
		}

		resdata.SetPrizeProtoIsNotEmpty(report.ShowPrize)

		// 加被瞭望次数
		hero.HistoryAmount().IncreaseOne(server_proto.HistoryAmountType_BeenInverstigation)
		heromodule.UpdateTaskProgress(hero, result, shared_proto.TaskTargetType_TASK_TARGET_ACCUM_BEEN_INVESTIGATION)

		result.Changed()
		result.Ok()
	}) {
		//if err == lock.ErrEmpty {
		//	logrus.WithField("id", targetId).Debug("侦查城池，英雄不存在")
		//	toSend = region.ERR_INVESTIGATE_FAIL_HERO_NOT_FOUND
		//	return
		//}

		if toSend == nil {
			logrus.Debug("侦查城池，lock英雄错误")
			toSend = region.ERR_INVESTIGATE_FAIL_SERVER_ERROR
		}
	}

	if toSend != nil {
		hc.Send(toSend)
		return
	}

	attackerFlagName := ""
	if attacker != nil {
		if selfGuildId != 0 {
			g := m.guildSnapshotService.GetSnapshot(selfGuildId)
			if g != nil {
				attacker.GuildFlagName = g.FlagName
			}
		}
		attackerFlagName = m.toFlagHeroName(attacker.GuildFlagName, attacker.Name)
	}

	defenserFlagName := ""
	if report.Defenser != nil {
		if defenserGuildId != 0 {
			g := m.guildSnapshotService.GetSnapshot(defenserGuildId)
			if g != nil {
				report.Defenser.GuildFlagName = g.FlagName
			}
		}
		defenserFlagName = m.toFlagHeroName(report.Defenser.GuildFlagName, report.Defenser.Name)
	}

	// 战报邮件
	attackerId := hc.Id()
	defenserId := targetId
	var attackerMailId uint64
	if data := m.datas.MailHelp().ReportWatchAttacker; data != nil {
		mailProto := data.NewTextMail(shared_proto.MailType_MailInvestigation)
		mailProto.SubTitle = data.NewSubTitleFields().WithFields("attacker", attackerFlagName).WithFields("defenser", defenserFlagName).JsonString()
		mailProto.Text = data.NewTextFields().WithFields("attacker", attackerFlagName).WithFields("defenser", defenserFlagName).JsonString()
		mailProto.Report = report
		m.mailModule.SendReportMail(attackerId, mailProto, ctime)

		attackerMailId, _ = i64.FromBytesU64(mailProto.Id)
	}

	if data := m.datas.MailHelp().ReportWatchDefenser; data != nil {
		mailProto := data.NewTextMail(shared_proto.MailType_MailBeenInvestigation)
		mailProto.SubTitle = data.NewSubTitleFields().WithFields("attacker", attackerFlagName).WithFields("defenser", defenserFlagName).JsonString()
		mailProto.Text = data.NewTextFields().WithFields("attacker", attackerFlagName).WithFields("defenser", defenserFlagName).JsonString()
		mailProto.Report = &shared_proto.FightReportProto{
			AttackerSide: false,
			Attacker:     attacker,
		}

		m.mailModule.SendReportMail(defenserId, mailProto, ctime)
	}

	// 成功
	hc.FuncWithSend(func(hero *entity.Hero, result herolock.LockResult) {

		// 加瞭望次数
		hero.HistoryAmount().IncreaseOne(server_proto.HistoryAmountType_Inverstigation)
		heromodule.UpdateTaskProgress(hero, result, shared_proto.TaskTargetType_TASK_TARGET_ACCUM_INVESTIGATION)

		nextTime := hero.GetNextInvestigateTime()
		if ctime.Before(hero.GetNextInvestigateTime()) && proto.Cost {
			// 扣点券
			heromodule.ReduceCostAnyway(hctx, hero, result, m.datas.RegionConfig().MiaoInvestigateCdCost)

		} else {
			nextTime = timeutil.Max(nextTime, ctime)

			// 加CD
			nextTime = nextTime.Add(timeutil.MaxDuration(m.datas.RegionConfig().InvestigateCd, 10*time.Second))
			hero.SetNextInvestigateTime(nextTime)
		}

		result.Add(region.NewS2cInvestigateMsg(proto.HeroId, timeutil.Marshal32(nextTime)))

		if attackerMailId > 0 {
			expireTime := ctime.Add(m.datas.RegionConfig().InvestigationExpireDuration)
			hero.AddInvestigation(targetId, expireTime, attackerMailId)
			if hero.GetInvestigationCount() >= m.datas.RegionConfig().InvestigationLimit {
				hero.ClearExpiredInvestigation(ctime, true)
			}

			result.Add(region.NewS2cGetPrevInvestigateMsg(idbytes.ToBytes(targetId), i64.ToBytesU64(attackerMailId), timeutil.Marshal32(expireTime)))
		}

		result.Ok()
	})

	m.pushService.PushFunc(shared_proto.SettingType_ST_INVESTIGATE, targetId, func(d *pushdata.PushData) (title, content string) {
		return d.Title, d.ReplaceContent("{{attacker}}", attackerFlagName)
	})

	// 破除免战
	if selfIsMian {
		m.realmService.GetBigMap().TryRemoveBaseMian(hc.Id())
	}
}

//gogen:iface
func (m *RegionModule) ProcessInvestigateInvade(proto *region.C2SInvestigateInvadeProto, hc iface.HeroController) {

	targetId, ok := idbytes.ToId(proto.HeroId)
	if !ok {
		logrus.WithField("id", targetId).Debug("侦查城池，玩家id无效")
		hc.Send(region.ERR_INVESTIGATE_INVADE_FAIL_HERO_NOT_FOUND)
		return
	}

	if hc.Id() == targetId {
		logrus.Debug("侦查城池，不能侦查自己")
		hc.Send(region.ERR_INVESTIGATE_INVADE_FAIL_SELF_ID)
		return
	}

	if npcid.IsNpcId(targetId) {
		logrus.Debug("侦查城池，目标是Npc")
		hc.Send(region.ERR_INVESTIGATE_INVADE_FAIL_HERO_NOT_FOUND)
		return
	}

	// 名城战期间，不能侦察
	if _, joined := m.mingcWarService.JoiningFightMingc(hc.Id()); joined {
		hc.Send(region.ERR_INVESTIGATE_INVADE_FAIL_IN_MC_WAR_FIGHT)
		return
	}

	realm := m.realmService.GetBigMap()
	if realm == nil {
		logrus.Debugf("出征，出征地图没找到")
		hc.Send(region.ERR_INVESTIGATE_INVADE_FAIL_MAP_NOT_FOUND)
		return
	}

	ctime := m.time.CurrentTime()
	// 检查CD
	var selfGuildId int64

	selfIsMian := false
	if hc.FuncWithSend(func(hero *entity.Hero, result herolock.LockResult) {
		selfIsMian = ctime.Before(hero.GetMianDisappearTime())

		// 自己处于新手免战，不能侦察别人
		if ctime.Before(hero.GetNewHeroMianDisappearTime()) &&
			hero.GetNewHeroMianDisappearTime().Equal(hero.GetMianDisappearTime()) {
			logrus.Debugf("侦查城池，新手免战期间不能侦查玩家")
			result.Add(region.ERR_INVESTIGATE_INVADE_FAIL_SELF_MIAN)
			return
		}

		// 检查主城是否存在
		if hero.BaseRegion() == 0 || hero.BaseLevel() <= 0 || hero.Prosperity() <= 0 {
			logrus.Debugf("侦察，主城流亡了")
			result.Add(region.ERR_INVESTIGATE_INVADE_FAIL_BASE_DESTROY)
			return
		}

		selfGuildId = hero.GuildId()

		if m.datas.RegionConfig().InvestigateCost != nil {
			// 检查侦察的消耗
			if !heromodule.HasEnoughCost(hero, m.datas.RegionConfig().InvestigateCost) {
				logrus.Debug("侦查城池, 消耗不足")
				result.Add(region.ERR_INVESTIGATE_INVADE_FAIL_COST_NOT_ENOUGH)
				return
			}
		}

		result.Ok()
	}) {
		return
	}

	hctx := heromodule.NewContext(m.dep, operate_type.RegionInvestigate)

	var toSend pbutil.Buffer
	// var defenserGuildId int64

	if m.heroDataServive.FuncWithSend(targetId, func(hero *entity.Hero, result herolock.LockResult) {
		if selfGuildId != 0 && selfGuildId == hero.GuildId() {
			logrus.Debug("侦查城池，不能侦查盟友")
			toSend = region.ERR_INVESTIGATE_INVADE_FAIL_SAME_GUILD
			return
		}

		if ctime.Before(hero.GetMianDisappearTime()) {
			logrus.Debug("侦查城池，目标免战中")
			toSend = region.ERR_INVESTIGATE_INVADE_FAIL_TARGET_MIAN
			return
		}
		// defenserGuildId = hero.GuildId()

		result.Changed()
		result.Ok()
	}) {
		if toSend == nil {
			logrus.Debug("侦查城池，lock英雄错误")
			toSend = region.ERR_INVESTIGATE_FAIL_SERVER_ERROR
		}
	}

	if toSend != nil {
		hc.Send(toSend)
		return
	}
	///成功 出发
	processed, err := realm.InvasionInvestigate(hc, targetId)
	if !processed {
		hc.Send(region.ERR_INVASION_FAIL_SERVER_ERROR)
		return
	}

	switch err {
	case nil:
		hc.Send(region.NewS2cInvestigateInvadeMsg(proto.HeroId))
	default:
		logrus.WithError(err)
		return
	}
	// 成功
	hc.FuncWithSend(func(hero *entity.Hero, result herolock.LockResult) {
		// 加侦察次数
		hero.HistoryAmount().IncreaseOne(server_proto.HistoryAmountType_Inverstigation)
		heromodule.UpdateTaskProgress(hero, result, shared_proto.TaskTargetType_TASK_TARGET_ACCUM_INVESTIGATION)

		// nextTime := hero.GetNextInvestigateTime()
		if m.datas.RegionConfig().InvestigateCost != nil {
			// 扣货币
			heromodule.ReduceCostAnyway(hctx, hero, result, m.datas.RegionConfig().InvestigateCost)
		}
		result.Add(region.NewS2cInvestigateInvadeMsg(proto.HeroId))

		result.Ok()
	})

	// m.pushService.PushFunc(shared_proto.SettingType_ST_INVESTIGATE, targetId, func(d *pushdata.PushData) (title, content string) {
	// 	return d.Title, d.ReplaceContent("{{attacker}}", attackerFlagName)
	// })

	// 破除免战
	if selfIsMian {
		m.realmService.GetBigMap().TryRemoveBaseMian(hc.Id())
	}
}

// 计算移动速度
//gogen:iface c2s_calc_move_speed
func (m *RegionModule) ProcessCalcMoveSpeed(proto *region.C2SCalcMoveSpeedProto, hc iface.HeroController) {
	id, ok := idbytes.ToId(proto.Id)
	if !ok {
		return
	}

	var rate float64
	hc.FuncNotError(func(hero *entity.Hero) (heroChanged bool) {
		rate = hero.GetMoveSpeedRate()
		return
	})

	hc.Send(region.NewS2cCalcMoveSpeedMsg(proto.Id, i32.MultiF64(1000, m.realmService.GetBigMap().CalcMoveSpeed(id, rate))))
}

func (r *RegionModule) toFlagHeroNameByGuildId(guildId int64, heroName string) string {
	return r.toFlagHeroName(r.getFlagName(guildId), heroName)
}

func (r *RegionModule) toFlagHeroName(flagName, heroName string) string {
	return r.datas.MiscConfig().FlagHeroName.FormatIgnoreEmpty(flagName, heroName)
}

func (r *RegionModule) getFlagName(guildId int64) string {
	if guildId != 0 {
		g := r.guildSnapshotService.GetSnapshot(guildId)
		if g != nil {
			return g.FlagName
		}
	}
	return ""
}

//gogen:iface c2s_list_enemy_pos
func (m *RegionModule) ProcessListEnemyPos(hc iface.HeroController) {

	var enemyIds []int64
	hc.FuncWithSend(func(hero *entity.Hero, result herolock.LockResult) {
		enemyIds = hero.Relation().EnemyIds()

		result.Ok()
	})

	result := &region.S2CListEnemyPosProto{}
	for _, id := range enemyIds {
		base := m.realmService.GetBigMap().GetRoBase(id)
		if base != nil {
			result.PosX = append(result.PosX, base.GetX())
			result.PosY = append(result.PosY, base.GetY())

			if len(result.PosX) >= m.datas.RegionConfig().ListEnemyPosCount {
				// 最多显示X个坐标
				break
			}
		}
	}

	hc.Send(region.NewS2cListEnemyPosProtoMsg(result))
}

//gogen:iface
func (m *RegionModule) ProcessSearchBaozNpc(proto *region.C2SSearchBaozNpcProto, hc iface.HeroController) {

	data := m.datas.GetBaozNpcData(u64.FromInt32(proto.DataId))
	if data == nil {
		logrus.Debug("搜索宝藏Npc，无效的配置id")
		hc.Send(region.ERR_SEARCH_BAOZ_NPC_FAIL_INVALID_DATA_ID)
		return
	}

	// 获取玩家自己主城所在所在的块的坐标
	base := m.realmService.GetBigMap().GetRoBase(hc.Id())
	if base == nil {
		logrus.Debug("搜索宝藏Npc，自己已经流亡")
		hc.Send(region.ERR_SEARCH_BAOZ_NPC_FAIL_HOME_NOT_ALIVE)
		return
	}

	blockData := m.realmService.GetBigMap().GetMapData()

	selfX, selfY := int(base.GetX()), int(base.GetY())

	var kvs []*sortkeys.U64KV
	for _, v := range blockData.GetRound4BlockByPos(selfX, selfY) {
		blockX, blockY := v.XY()
		blockSequence := regdata.BlockSequence(uint64(blockX), uint64(blockY))

		for i := uint64(0); i < data.KeepCount; i++ {
			id := npcid.NewBaoZangNpcId(blockSequence, i, data.Id)
			base := m.realmService.GetBigMap().GetRoBase(id)
			if base != nil {
				distance := hexagon.Distance(selfX, selfY, int(base.GetX()), int(base.GetY()))
				kvs = append(kvs, sortkeys.NewU64KV(uint64(distance), base))
			}
		}
	}

	// 距离近的排在前面
	sort.Sort(sortkeys.U64KVSlice(kvs))

	breakCount := m.datas.RegionConfig().SearchBaozNpcCount

	toSendProto := &region.S2CSearchBaozNpcProto{}
	toSendProto.DataId = proto.DataId
	for _, kv := range kvs {
		if base, _ := kv.V.(*server_proto.RoBaseProto); base != nil {
			toSendProto.BaseId = append(toSendProto.BaseId, i64.ToBytes(base.GetId()))
			toSendProto.BaseX = append(toSendProto.BaseX, base.GetX())
			toSendProto.BaseY = append(toSendProto.BaseY, base.GetY())

			if len(toSendProto.BaseId) >= breakCount {
				break
			}
		}
	}

	hc.Send(region.NewS2cSearchBaozNpcProtoMsg(toSendProto))
}

//gogen:iface
func (m *RegionModule) ProcessUseMultiLevelNpcTimesGoods(proto *region.C2SUseMultiLevelNpcTimesGoodsProto, hc iface.HeroController) {

	// 使用物品

	goodsData := m.datas.GetGoodsData(u64.FromInt32(proto.Id))
	if goodsData == nil {
		logrus.Debug("使用讨伐令，无效的物品id")
		hc.Send(region.ERR_USE_MULTI_LEVEL_NPC_TIMES_GOODS_FAIL_INVALID_ID)
		return
	}

	if goodsData.GoodsEffect == nil || !goodsData.GoodsEffect.AddMultiLevelNpcTimes {
		logrus.Debug("使用讨伐令，不是讨伐令")
		hc.Send(region.ERR_USE_MULTI_LEVEL_NPC_TIMES_GOODS_FAIL_INVALID_GOODS)
		return
	}

	ctime := m.time.CurrentTime()
	hc.FuncWithSend(func(hero *entity.Hero, result herolock.LockResult) {

		if !m.useMultiLevelNpcTimesGoods(hero, result, goodsData, 1, proto.Buy, ctime,
			region.ERR_USE_MULTI_LEVEL_NPC_TIMES_GOODS_FAIL_FULL_TIMES,
			region.ERR_USE_MULTI_LEVEL_NPC_TIMES_GOODS_FAIL_COST_NOT_ENOUGH,
			region.ERR_USE_MULTI_LEVEL_NPC_TIMES_GOODS_FAIL_COUNT_NOT_ENOUGH, ) {
			return
		}

		result.Add(region.NewS2cUseMultiLevelNpcTimesGoodsMsg(proto.Id, proto.Buy))

		result.Ok()
	})

}

//gogen:iface
func (m *RegionModule) ProcessUseInvaseHeroTimesGoods(proto *region.C2SUseInvaseHeroTimesGoodsProto, hc iface.HeroController) {

	// 使用物品

	goodsData := m.datas.GetGoodsData(u64.FromInt32(proto.Id))
	if goodsData == nil {
		logrus.Debug("使用攻城令，无效的物品id")
		hc.Send(region.ERR_USE_INVASE_HERO_TIMES_GOODS_FAIL_INVALID_ID)
		return
	}

	if goodsData.GoodsEffect == nil || !goodsData.GoodsEffect.AddInvaseHeroTimes {
		logrus.Debug("使用攻城令，不是讨伐令")
		hc.Send(region.ERR_USE_INVASE_HERO_TIMES_GOODS_FAIL_INVALID_GOODS)
		return
	}

	ctime := m.time.CurrentTime()
	hc.FuncWithSend(func(hero *entity.Hero, result herolock.LockResult) {

		if !m.useInvaseHeroTimesGoods(hero, result, goodsData, proto.Buy, ctime,
			region.ERR_USE_INVASE_HERO_TIMES_GOODS_FAIL_FULL_TIMES,
			region.ERR_USE_INVASE_HERO_TIMES_GOODS_FAIL_COST_NOT_ENOUGH,
			region.ERR_USE_INVASE_HERO_TIMES_GOODS_FAIL_COUNT_NOT_ENOUGH, ) {
			return
		}

		result.Add(region.NewS2cUseInvaseHeroTimesGoodsMsg(proto.Id, proto.Buy))

		result.Ok()
	})

}

func (m *RegionModule) useMultiLevelNpcTimesGoods(hero *entity.Hero, result herolock.LockResult,
	goodsData *goods.GoodsData, count uint64, autoBuy bool, ctime time.Time,
	fullTimesMsg, costNotEnoughMsg, countNotEnoughMsg pbutil.Buffer) bool {

	hctx := heromodule.NewContext(m.dep, operate_type.RegionUseMultiLevelNpcTimesGoods)
	rt := hero.GetMultiLevelNpcTimes()
	success := useTimesGoodsMulti("使用讨伐令", hctx, hero, result,
		rt, m.extraTimesService.MultiLevelNpcMaxTimes().TotalTimes(),
		goodsData, count, autoBuy, ctime, fullTimesMsg, costNotEnoughMsg, countNotEnoughMsg)

	if success {
		result.Add(region.NewS2cUpdateMultiLevelNpcTimesMsg(rt.StartTimeUnix32(), nil))
	}

	return success
}

func (m *RegionModule) useInvaseHeroTimesGoods(hero *entity.Hero, result herolock.LockResult,
	goodsData *goods.GoodsData, autoBuy bool, ctime time.Time,
	fullTimesMsg, costNotEnoughMsg, countNotEnoughMsg pbutil.Buffer) bool {

	hctx := heromodule.NewContext(m.dep, operate_type.RegionUseInvaseHeroTimesGoods)
	rt := hero.GetInvaseHeroTimes()
	success := useTimesGoods("使用攻城令", hctx, hero, result,
		rt, 0,
		goodsData, autoBuy, ctime, fullTimesMsg, costNotEnoughMsg, countNotEnoughMsg)

	if success {
		result.Add(region.NewS2cUpdateInvaseHeroTimesMsg(rt.StartTimeUnix32(), nil))
	}

	return success
}

func (m *RegionModule) useJunTuanNpcTimesGoods(hero *entity.Hero, result herolock.LockResult,
	goodsData *goods.GoodsData, autoBuy bool, ctime time.Time,
	fullTimesMsg, costNotEnoughMsg, countNotEnoughMsg pbutil.Buffer) bool {

	hctx := heromodule.NewContext(m.dep, operate_type.RegionUseJunTuanNpcTimesGoods)
	rt := hero.GetJunTuanNpcTimes()
	success := useTimesGoods("使用攻城令", hctx, hero, result,
		rt, 0,
		goodsData, autoBuy, ctime, fullTimesMsg, costNotEnoughMsg, countNotEnoughMsg)

	if success {
		result.Add(region.NewS2cUpdateJunTuanNpcTimesMsg(rt.StartTimeUnix32(), nil))
	}

	return success
}

func useTimesGoods(operateString string, hctx *heromodule.HeroContext, hero *entity.Hero, result herolock.LockResult,
	recoverTimes *recovtimes.ExtraRecoverTimes, extraTimes uint64,
	goodsData *goods.GoodsData, autoBuy bool, ctime time.Time,
	fullTimesMsg, costNotEnoughMsg, countNotEnoughMsg pbutil.Buffer) bool {
	return useTimesGoodsMulti(operateString, hctx, hero, result, recoverTimes, extraTimes, goodsData, 1, autoBuy, ctime, fullTimesMsg, costNotEnoughMsg, countNotEnoughMsg)
}

func useTimesGoodsMulti(operateString string, hctx *heromodule.HeroContext, hero *entity.Hero, result herolock.LockResult,
	recoverTimes *recovtimes.ExtraRecoverTimes, extraTimes uint64,
	goodsData *goods.GoodsData, count uint64, autoBuy bool, ctime time.Time,
	fullTimesMsg, costNotEnoughMsg, countNotEnoughMsg pbutil.Buffer) bool {

	currentTimes := recoverTimes.Times(ctime, extraTimes)
	maxTimes := recoverTimes.MaxTimes(extraTimes)

	if currentTimes+count > maxTimes {
		logrus.Debugf("%s，次数已满", operateString)
		result.Add(fullTimesMsg)
		return false
	}

	if !heromodule.TryReduceOrBuyGoods(hctx, hero, result, goodsData, count, autoBuy) {
		if autoBuy {
			logrus.Debug("%s，购买消耗不足", operateString)
			result.Add(costNotEnoughMsg)
			return false
		} else {
			logrus.Debug("%s，物品个数不足", operateString)
			result.Add(countNotEnoughMsg)
			return false
		}
	}

	recoverTimes.AddTimes(count, ctime, extraTimes)

	return true
}

//gogen:iface c2s_home_ast_defending_info
func (m *RegionModule) ProcessHomeAstDefendingInfo(hc iface.HeroController) {
	if hc.FuncWithSend(func(hero *entity.Hero, result herolock.LockResult) {
		if hero.Prosperity() >= hero.ProsperityCapcity() {
			result.Add(region.ERR_HOME_AST_DEFENDING_INFO_FAIL_PROSPERITY_FULL)
			return
		}
		result.Ok()
	}) {
		return
	}

	heros := m.realmService.GetBigMap().GetAstDefendHeros(hc.Id())
	logs := m.realmService.GetBigMap().GetAstDefendLogsByHero(hc.Id())

	hc.Send(region.NewS2cHomeAstDefendingInfoMsg(heros, logs))
}

//gogen:iface c2s_guild_please_help_me
func (m *RegionModule) ProcessGuildPleaseHelpMe(hc iface.HeroController) {
	if hc.FuncWithSend(func(hero *entity.Hero, result herolock.LockResult) {
		if hero.Prosperity() >= hero.ProsperityCapcity() {
			result.Add(region.ERR_GUILD_PLEASE_HELP_ME_FAIL_PROSPERITY_FULL)
			return
		}
		result.Ok()
	}) {
		return
	}

	if m.realmService.GetBigMap().GetDefendingTroopCount(hc.Id()) >= m.datas.RegionConfig().MaxAssist {
		hc.Send(region.ERR_GUILD_PLEASE_HELP_ME_FAIL_TROOPS_LIMIT)
		return
	}

	gid, ok := hc.LockGetGuildId()
	if !ok || gid <= 0 {
		hc.Send(region.ERR_GUILD_PLEASE_HELP_ME_FAIL_NO_GUILD)
		return
	}

	if d := m.datas.TextHelp().GuildPleaseHelpMe; d != nil {
		m.dep.Chat().SysChat(0, gid, shared_proto.ChatType_ChatGuild, d.Text.New().JsonString(), shared_proto.ChatMsgType_ChatMsgGuildLog, false, true, true, false)
	}
	hc.Send(region.GUILD_PLEASE_HELP_ME_S2C)
}

//gogen:iface
func (m *RegionModule) ProcessCreateAssembly(proto *region.C2SCreateAssemblyProto, hc iface.HeroController) {
	// 检查id
	targetId, ok := idbytes.ToId(proto.Target)
	if !ok {
		logrus.Errorf("RegionModule.ProcessCreateAssembly 解析Target id失败, %s", proto.Target)
		hc.Send(region.ERR_CREATE_ASSEMBLY_FAIL_INVALID_TARGET)
		return
	}

	heroId := hc.Id()
	if targetId == heroId {
		logrus.Debugf("创建集结，目标id跟自己的id一样 %s", proto.Target)
		hc.Send(region.ERR_CREATE_ASSEMBLY_FAIL_INVALID_TARGET)
		return
	}

	assemblyData := m.datas.GetAssemblyData(regdata.GetAssemblyTypeByTarget(targetId))
	if assemblyData == nil {
		logrus.Debugf("创建集结，assemblyData == nil")
		hc.Send(region.ERR_CREATE_ASSEMBLY_FAIL_INVALID_TARGET_ASSEMBLY)
		return
	}

	if !npcid.IsNpcId(targetId) {
		// 名城战期间，不能集结
		if _, joined := m.mingcWarService.JoiningFightMingc(hc.Id()); joined {
			hc.Send(region.ERR_CREATE_ASSEMBLY_FAIL_IN_MC_WAR_FIGHT)
			return
		}
	}

	waitIndex := int(proto.WaitIndex)
	if waitIndex < 0 || waitIndex >= len(assemblyData.WaitDuration) {
		logrus.Debugf("创建集结，invalid wait index")
		hc.Send(region.ERR_CREATE_ASSEMBLY_FAIL_INVALID_WAIT_INDEX)
		return
	}
	waitDuration := assemblyData.WaitDuration[waitIndex]

	//realmId := int64(proto.MapId)
	//
	//if realmId == 0 {
	//	// 临时处理一下, 向前兼容
	//	logrus.Debug("客户端的invasion必须要带map_id, 临时取他主城作为map_id")
	//	hc.FuncNotError(func(hero *entity.Hero) (heroChanged bool) {
	//		realmId = hero.BaseRegion()
	//		return false
	//	})
	//}

	//realm := m.realmService.GetRealm(realmId)
	realm := m.realmService.GetBigMap()
	if realm == nil {
		logrus.Debugf("创建集结，出征地图没找到")
		hc.Send(region.ERR_CREATE_ASSEMBLY_FAIL_INVALID_TARGET)
		return
	}

	troopIndex := u64.FromInt32(proto.TroopIndex - 1)

	targetLevel := u64.FromInt32(proto.TargetLevel)

	// 各种检查
	var troopId, guildId int64
	hasError := hc.FuncWithSend(func(hero *entity.Hero, result herolock.LockResult) {

		// 检查主城是否存在
		if hero.BaseRegion() == 0 || hero.BaseLevel() <= 0 || hero.Prosperity() <= 0 {
			logrus.Debugf("创建集结，主城流亡了")
			result.Add(region.ERR_CREATE_ASSEMBLY_FAIL_SELF_NOT_EXIST)
			return
		}

		if hero.GuildId() == 0 {
			logrus.Debugf("创建集结，没有联盟不能创建集结")
			result.Add(region.ERR_CREATE_ASSEMBLY_FAIL_SELF_NOT_GUILD)
			return
		}

		ctime := m.time.CurrentTime()

		if !npcid.IsNpcId(targetId) {
			// 自己处于新手免战，不能出征打别人
			if ctime.Before(hero.GetNewHeroMianDisappearTime()) &&
				hero.GetNewHeroMianDisappearTime().Equal(hero.GetMianDisappearTime()) {
				logrus.Debugf("创建集结，新手免战期间不能出征玩家")
				result.Add(region.ERR_CREATE_ASSEMBLY_FAIL_MIAN)
				return
			}
		}

		// 不能对同一个目标进行集结
		for _, t := range hero.Troops() {
			if ii := t.GetInvateInfo(); ii != nil && ii.AssemblyId() == targetId {
				logrus.Debugf("创建集结，已经有队伍对这个目标进行集结了")
				result.Add(region.ERR_CREATE_ASSEMBLY_FAIL_SAME_TARGET)
				return
			}
		}

		t := hero.GetTroopByIndex(troopIndex)
		if t == nil {
			logrus.Debugf("创建集结，无效的部队编号")
			result.Add(region.ERR_CREATE_ASSEMBLY_FAIL_INVALID_TROOP_INDEX)
			return
		}

		if t.IsOutside() {
			logrus.Debugf("创建集结，队伍出征中")
			result.Add(region.ERR_CREATE_ASSEMBLY_FAIL_OUTSIDE)
			return
		}
		troopId = t.Id()

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
			logrus.Debugf("创建集结，队伍没有士兵")
			result.Add(region.ERR_CREATE_ASSEMBLY_FAIL_NO_SOLDIER)
			return
		}

		if !npcid.IsNpcId(targetId) {
			// 出征次数
			if hero.GetInvaseHeroTimes().Times(ctime, 0) <= 0 {
				if proto.GoodsId == 0 {
					logrus.Debugf("创建集结，讨伐次数不足")
					result.Add(region.ERR_CREATE_ASSEMBLY_FAIL_TIMES_LIMIT)
					return
				} else {
					goodsData := m.datas.GetGoodsData(u64.FromInt32(proto.GoodsId))
					if goodsData == nil {
						logrus.Debug("创建集结使用讨伐令，无效的物品id")
						result.Add(region.ERR_CREATE_ASSEMBLY_FAIL_INVALID_GOODS)
						return
					}

					if goodsData.GoodsEffect == nil || !goodsData.GoodsEffect.AddInvaseHeroTimes {
						logrus.Debug("创建集结使用讨伐令，不是讨伐令")
						result.Add(region.ERR_CREATE_ASSEMBLY_FAIL_INVALID_GOODS)
						return
					}

					if !m.useInvaseHeroTimesGoods(hero, result, goodsData, proto.AutoBuy, ctime,
						region.ERR_CREATE_ASSEMBLY_FAIL_TIMES_LIMIT,
						region.ERR_CREATE_ASSEMBLY_FAIL_COST_NOT_ENOUGH,
						region.ERR_CREATE_ASSEMBLY_FAIL_COST_NOT_ENOUGH, ) {

						logrus.Debug("创建集结，使用讨伐令失败")
						return
					}

					// 使用物品加次数成功
				}

			}
		} else {
			targetNpcType := npcid.GetNpcIdType(targetId)
			switch targetNpcType {
			case npcid.NpcType_XiongNu:
				guildId := npcid.GetXiongNuGuildId(targetId)
				if hero.GuildId() != guildId {
					logrus.Debugf("创建集结，不是你联盟的匈奴")
					result.Add(region.ERR_CREATE_ASSEMBLY_FAIL_TIMES_LIMIT)
					return
				}

			case npcid.NpcType_JunTuan:
				data := m.datas.GetJunTuanNpcData(npcid.GetNpcDataId(targetId))
				if data == nil {
					logrus.Debugf("创建集结，军团怪Dataid不存在")
					result.Add(region.ERR_CREATE_ASSEMBLY_FAIL_INVALID_TARGET)
					return
				}

				if data.RequiredHeroLevel > 0 && hero.Level() < data.RequiredHeroLevel {
					logrus.Debugf("创建集结，军团怪要求的君主等级不足")
					result.Add(region.ERR_CREATE_ASSEMBLY_FAIL_REQUIRED_HERO_LEVEL)
					return
				}

				var useTime uint64
				for _, t := range hero.Troops() {
					if ii := t.GetInvateInfo(); ii != nil {
						if npcid.IsJunTuanNpcId(ii.AssemblyTargetId()) &&
							ii.State() == realmface.MovingToInvade &&
							ii.State() == realmface.Assembly &&
							ii.State() == realmface.MovingToAssembly &&
							ii.State() == realmface.AssemblyArrived {
							useTime++
						}
					}
				}

				// 出征次数
				if hero.GetJunTuanNpcTimes().Times(ctime, 0) <= useTime {
					if proto.GoodsId == 0 {
						logrus.Debugf("创建集结，讨伐军团Npc次数不足")
						result.Add(region.ERR_CREATE_ASSEMBLY_FAIL_TIMES_LIMIT)
						return
					} else {
						goodsData := m.datas.GetGoodsData(u64.FromInt32(proto.GoodsId))
						if goodsData == nil {
							logrus.Debug("创建集结使用军团讨伐令，无效的物品id")
							result.Add(region.ERR_CREATE_ASSEMBLY_FAIL_INVALID_GOODS)
							return
						}

						if goodsData.GoodsEffect == nil || goodsData.GoodsEffect.AmountType != shared_proto.GoodsAmountType_AmountJunTuan {
							logrus.Debug("创建集结使用军团讨伐令，不是讨伐令")
							result.Add(region.ERR_CREATE_ASSEMBLY_FAIL_INVALID_GOODS)
							return
						}

						if !m.useJunTuanNpcTimesGoods(hero, result, goodsData, proto.AutoBuy, ctime,
							region.ERR_CREATE_ASSEMBLY_FAIL_TIMES_LIMIT,
							region.ERR_CREATE_ASSEMBLY_FAIL_COST_NOT_ENOUGH,
							region.ERR_CREATE_ASSEMBLY_FAIL_COST_NOT_ENOUGH, ) {

							logrus.Debug("创建集结，使用军团讨伐令失败")
							return
						}

						// 使用物品加次数成功
					}
				}

			default:
				logrus.Error("创建集结，无效的目标类型（上面不是判断过了吗）")
				result.Add(region.ERR_CREATE_ASSEMBLY_FAIL_TIMES_LIMIT)
				return
			}
		}

		guildId = hero.GuildId()

		result.Ok()
	})

	if hasError {
		return
	}

	processed, err := realm.CreateAssembly(hc, targetId, targetLevel, troopIndex, assemblyData.MemberCount, waitDuration)
	if !processed {
		hc.Send(region.ERR_CREATE_ASSEMBLY_FAIL_SERVER_ERROR)
		return
	}

	switch err {
	case nil:
		hc.Send(region.NewS2cCreateAssemblyMsg(proto.TroopIndex, proto.Target, idbytes.ToBytes(troopId)))

		if guildId != 0 {

			var targetName, targetFlagName string
			var targetCountry uint64
			if !npcid.IsNpcId(targetId) {
				if target := m.heroSnapshotServive.Get(targetId); target != nil {
					targetName = target.Name
					targetFlagName = target.GuildFlagName()
					targetCountry = target.CountryId
				} else {
					targetName = idbytes.PlayerName(targetId)
				}
			} else {
				if b := m.realmService.GetBigMap().GetRoBase(targetId); b != nil {
					targetName = b.Name
				} else {
					targetName = idbytes.PlayerName(targetId)
				}
			}

			showJson := fmt.Sprintf(`{"room_id":"%v","target_name":"%s","target_country_id":%d,"target_guild_flag":"%s"}`, troopId, targetName, targetCountry, targetFlagName)
			m.dep.Chat().SysChat(hc.Id(), guildId, shared_proto.ChatType_ChatGuild, showJson, shared_proto.ChatMsgType_ChatMsgCreateAssembly, false, true, false, true)
		}

	case realmerr.ErrLockHeroErr:
		hc.Send(region.ERR_CREATE_ASSEMBLY_FAIL_SERVER_ERROR)

	case realmerr.ErrCreateAssemblyInvalidInput:
		hc.Send(region.ERR_CREATE_ASSEMBLY_FAIL_INVALID_TARGET)

	case realmerr.ErrCreateAssemblyInvalidTarget:
		hc.Send(region.ERR_CREATE_ASSEMBLY_FAIL_INVALID_TARGET)

	case realmerr.ErrCreateAssemblyTargetNotExist:
		hc.Send(region.ERR_CREATE_ASSEMBLY_FAIL_TARGET_NOT_EXIST)

	case realmerr.ErrCreateAssemblySelfNoBase:
		hc.Send(region.ERR_CREATE_ASSEMBLY_FAIL_SELF_NOT_EXIST)

	case realmerr.ErrCreateAssemblySelfNoGuild:
		hc.Send(region.ERR_CREATE_ASSEMBLY_FAIL_INVALID_TARGET)

	case realmerr.ErrCreateAssemblyEmptyGeneral:
		hc.Send(region.ERR_CREATE_ASSEMBLY_FAIL_INVALID_TROOP_INDEX)

	case realmerr.ErrCreateAssemblyGeneralOutside:
		hc.Send(region.ERR_CREATE_ASSEMBLY_FAIL_OUTSIDE)

	case realmerr.ErrCreateAssemblyNoSoldier:
		hc.Send(region.ERR_CREATE_ASSEMBLY_FAIL_NO_SOLDIER)

	case realmerr.ErrCreateAssemblyInvalidRelation:
		hc.Send(region.ERR_CREATE_ASSEMBLY_FAIL_INVALID_TARGET)

	case realmerr.ErrCreateAssemblyInvalidTroopIndex:
		hc.Send(region.ERR_CREATE_ASSEMBLY_FAIL_INVALID_TROOP_INDEX)

	case realmerr.ErrCreateAssemblyMian:
		hc.Send(region.ERR_CREATE_ASSEMBLY_FAIL_MIAN)

	case realmerr.ErrCreateAssemblyTodayJoinXiongNu:
		hc.Send(region.ERR_CREATE_ASSEMBLY_FAIL_TODAY_JOIN_XIONG_NU)

	default:
		logrus.WithError(err).Error("region.ProcessCreateAssembly有err没有处理")
		hc.Send(region.ERR_CREATE_ASSEMBLY_FAIL_SERVER_ERROR)
	}
}

var showNotExistMsg = region.NewS2cShowAssemblyMsg(true, nil, 0, nil).Static()

//gogen:iface
func (m *RegionModule) ProcessShowAssembly(proto *region.C2SShowAssemblyProto, hc iface.HeroController) {

	targetTroopId, ok := idbytes.ToId(proto.Id)
	if !ok {
		logrus.Errorf("RegionModule.ProcessShowAssembly 解析Target troop id失败, %s", proto.Id)
		hc.Send(showNotExistMsg)
		return
	}

	targetId := entity.GetTroopHeroId(targetTroopId)

	m.realmService.GetBigMap().ShowAssembly(hc, targetId, targetTroopId, proto.Version)
}

//gogen:iface
func (m *RegionModule) ProcessJoinAssembly(proto *region.C2SJoinAssemblyProto, hc iface.HeroController) {
	// 检查id
	joinTargetTroopId, ok := idbytes.ToId(proto.Id)
	if !ok {
		logrus.Errorf("RegionModule.ProcessJoinAssembly 解析Target troop id失败, %s", proto.Id)
		hc.Send(region.ERR_JOIN_ASSEMBLY_FAIL_INVALID_TARGET)
		return
	}

	joinTargetId := entity.GetTroopHeroId(joinTargetTroopId)
	if npcid.IsNpcId(joinTargetId) {
		logrus.Errorf("加入集结，目标是个npc")
		hc.Send(region.ERR_JOIN_ASSEMBLY_FAIL_INVALID_TARGET)
		return
	}

	heroId := hc.Id()
	if joinTargetId == heroId {
		logrus.Debugf("加入集结，目标id跟自己的id一样 %s", proto.Id)
		hc.Send(region.ERR_JOIN_ASSEMBLY_FAIL_INVALID_TARGET)
		return
	}

	// 名城战期间，不能集结
	if _, joined := m.mingcWarService.JoiningFightMingc(hc.Id()); joined {
		hc.Send(region.ERR_JOIN_ASSEMBLY_FAIL_IN_MC_WAR_FIGHT)
		return
	}

	var errMsg pbutil.Buffer
	var assemblyTargetId int64
	var joinTargetGuildId int64
	if hasErr := m.heroDataServive.FuncNotError(joinTargetId, func(hero *entity.Hero) (heroChanged bool) {
		t := hero.Troop(joinTargetTroopId)
		if t == nil {
			logrus.Debugf("加入集结，目标队伍为空")
			errMsg = region.ERR_JOIN_ASSEMBLY_FAIL_INVALID_TARGET
			return
		}

		ii := t.GetInvateInfo()
		if ii == nil {
			logrus.Debugf("加入集结，目标队伍没有集结")
			errMsg = region.ERR_JOIN_ASSEMBLY_FAIL_INVALID_TARGET
			return
		}

		if ii.State() != realmface.Assembly {
			logrus.Debugf("加入集结，目标队伍不是集结等待状态")
			errMsg = region.ERR_JOIN_ASSEMBLY_FAIL_INVALID_TARGET
			return
		}

		joinTargetGuildId = hero.GuildId()
		assemblyTargetId = ii.OriginTargetID()
		return
	}); hasErr {
		hc.Send(region.ERR_JOIN_ASSEMBLY_FAIL_SERVER_ERROR)
		return
	}

	if errMsg != nil {
		hc.Send(errMsg)
		return
	}

	//realmId := int64(proto.MapId)
	//
	//if realmId == 0 {
	//	// 临时处理一下, 向前兼容
	//	logrus.Debug("客户端的invasion必须要带map_id, 临时取他主城作为map_id")
	//	hc.FuncNotError(func(hero *entity.Hero) (heroChanged bool) {
	//		realmId = hero.BaseRegion()
	//		return false
	//	})
	//}

	//realm := m.realmService.GetRealm(realmId)
	realm := m.realmService.GetBigMap()
	if realm == nil {
		logrus.Debugf("加入集结，出征地图没找到")
		hc.Send(region.ERR_JOIN_ASSEMBLY_FAIL_INVALID_TARGET)
		return
	}

	troopIndex := u64.FromInt32(proto.TroopIndex - 1)

	// 各种检查
	var troopId int64
	hasError := hc.FuncWithSend(func(hero *entity.Hero, result herolock.LockResult) {

		// 检查主城是否存在
		if hero.BaseRegion() == 0 || hero.BaseLevel() <= 0 || hero.Prosperity() <= 0 {
			logrus.Debugf("加入集结，主城流亡了")
			result.Add(region.ERR_JOIN_ASSEMBLY_FAIL_SELF_NOT_EXIST)
			return
		}

		if hero.GuildId() == 0 {
			logrus.Debugf("加入集结，没有联盟不能加入集结")
			result.Add(region.ERR_JOIN_ASSEMBLY_FAIL_SELF_NOT_GUILD)
			return
		}

		if hero.GuildId() != joinTargetGuildId {
			logrus.Debugf("加入集结，不能加入别的联盟的集结")
			result.Add(region.ERR_JOIN_ASSEMBLY_FAIL_SELF_NOT_GUILD)
			return
		}

		ctime := m.time.CurrentTime()

		if !npcid.IsNpcId(assemblyTargetId) {
			// 自己处于新手免战，不能出征打别人
			if ctime.Before(hero.GetNewHeroMianDisappearTime()) &&
				hero.GetNewHeroMianDisappearTime().Equal(hero.GetMianDisappearTime()) {
				logrus.Debugf("加入集结，新手免战期间不能出征玩家")
				result.Add(region.ERR_JOIN_ASSEMBLY_FAIL_MIAN)
				return
			}
		}

		// 不能有多个队伍加入同一个集结
		for _, t := range hero.Troops() {
			if ii := t.GetInvateInfo(); ii != nil && ii.AssemblyId() == joinTargetTroopId {
				logrus.Debugf("加入集结，已经有队伍加入这个集结了")
				result.Add(region.ERR_JOIN_ASSEMBLY_FAIL_MULTI_JOIN)
				return
			}
		}

		t := hero.GetTroopByIndex(troopIndex)
		if t == nil {
			logrus.Debugf("加入集结，无效的部队编号")
			result.Add(region.ERR_JOIN_ASSEMBLY_FAIL_INVALID_TROOP_INDEX)
			return
		}

		if t.IsOutside() {
			logrus.Debugf("加入集结，队伍出征中")
			result.Add(region.ERR_JOIN_ASSEMBLY_FAIL_OUTSIDE)
			return
		}
		troopId = t.Id()

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
			logrus.Debugf("加入集结，队伍没有士兵")
			result.Add(region.ERR_JOIN_ASSEMBLY_FAIL_NO_SOLDIER)
			return
		}

		if !npcid.IsNpcId(assemblyTargetId) {
			// 出征次数
			if hero.GetInvaseHeroTimes().Times(ctime, 0) <= 0 {
				if proto.GoodsId == 0 {
					logrus.Debugf("加入集结，讨伐次数不足")
					result.Add(region.ERR_JOIN_ASSEMBLY_FAIL_TIMES_LIMIT)
					return
				} else {
					goodsData := m.datas.GetGoodsData(u64.FromInt32(proto.GoodsId))
					if goodsData == nil {
						logrus.Debug("加入集结使用讨伐令，无效的物品id")
						result.Add(region.ERR_JOIN_ASSEMBLY_FAIL_INVALID_GOODS)
						return
					}

					if goodsData.GoodsEffect == nil || !goodsData.GoodsEffect.AddInvaseHeroTimes {
						logrus.Debug("加入集结使用讨伐令，不是讨伐令")
						result.Add(region.ERR_JOIN_ASSEMBLY_FAIL_INVALID_GOODS)
						return
					}

					if !m.useInvaseHeroTimesGoods(hero, result, goodsData, proto.AutoBuy, ctime,
						region.ERR_JOIN_ASSEMBLY_FAIL_TIMES_LIMIT,
						region.ERR_JOIN_ASSEMBLY_FAIL_COST_NOT_ENOUGH,
						region.ERR_JOIN_ASSEMBLY_FAIL_COST_NOT_ENOUGH, ) {

						logrus.Debug("加入集结，使用讨伐令失败")
						return
					}

					// 使用物品加次数成功
				}
			}
		} else {
			targetNpcType := npcid.GetNpcIdType(assemblyTargetId)
			switch targetNpcType {
			case npcid.NpcType_XiongNu:
			case npcid.NpcType_JunTuan:
				data := m.datas.GetJunTuanNpcData(npcid.GetNpcDataId(assemblyTargetId))
				if data == nil {
					logrus.Debugf("加入集结，军团怪Dataid不存在")
					result.Add(region.ERR_JOIN_ASSEMBLY_FAIL_INVALID_TARGET)
					return
				}

				if data.RequiredHeroLevel > 0 && hero.Level() < data.RequiredHeroLevel {
					logrus.Debugf("加入集结，军团怪要求的君主等级不足")
					result.Add(region.ERR_JOIN_ASSEMBLY_FAIL_REQUIRED_HERO_LEVEL)
					return
				}

				var useTime uint64
				for _, t := range hero.Troops() {
					if ii := t.GetInvateInfo(); ii != nil {
						if npcid.IsJunTuanNpcId(ii.AssemblyTargetId()) &&
							ii.State() == realmface.MovingToInvade &&
							ii.State() == realmface.Assembly &&
							ii.State() == realmface.MovingToAssembly &&
							ii.State() == realmface.AssemblyArrived {
							useTime++
						}
					}
				}

				// 出征次数
				if hero.GetJunTuanNpcTimes().Times(ctime, 0) <= useTime {
					if proto.GoodsId == 0 {
						logrus.Debugf("加入集结，讨伐军团Npc次数不足")
						result.Add(region.ERR_JOIN_ASSEMBLY_FAIL_TIMES_LIMIT)
						return
					} else {
						goodsData := m.datas.GetGoodsData(u64.FromInt32(proto.GoodsId))
						if goodsData == nil {
							logrus.Debug("加入集结使用军团讨伐令，无效的物品id")
							result.Add(region.ERR_JOIN_ASSEMBLY_FAIL_INVALID_GOODS)
							return
						}

						if goodsData.GoodsEffect == nil || goodsData.GoodsEffect.AmountType != shared_proto.GoodsAmountType_AmountJunTuan {
							logrus.Debug("加入集结使用军团讨伐令，不是讨伐令")
							result.Add(region.ERR_JOIN_ASSEMBLY_FAIL_INVALID_GOODS)
							return
						}

						if !m.useJunTuanNpcTimesGoods(hero, result, goodsData, proto.AutoBuy, ctime,
							region.ERR_JOIN_ASSEMBLY_FAIL_TIMES_LIMIT,
							region.ERR_JOIN_ASSEMBLY_FAIL_COST_NOT_ENOUGH,
							region.ERR_JOIN_ASSEMBLY_FAIL_COST_NOT_ENOUGH, ) {

							logrus.Debug("加入集结，使用军团讨伐令失败")
							return
						}

						// 使用物品加次数成功
					}
				}

			}
		}

		result.Ok()
	})

	if hasError {
		return
	}

	processed, err := realm.JoinAssembly(hc, joinTargetId, joinTargetTroopId, troopIndex)
	if !processed {
		hc.Send(region.ERR_JOIN_ASSEMBLY_FAIL_SERVER_ERROR)
		return
	}

	switch err {
	case nil:
		hc.Send(region.NewS2cJoinAssemblyMsg(idbytes.ToBytes(troopId), proto.TroopIndex))

	case realmerr.ErrLockHeroErr:
		hc.Send(region.ERR_JOIN_ASSEMBLY_FAIL_SERVER_ERROR)

	case realmerr.ErrJoinAssemblyInvalidTarget:
		hc.Send(region.ERR_JOIN_ASSEMBLY_FAIL_INVALID_TARGET)

	case realmerr.ErrJoinAssemblyTargetNotExist:
		hc.Send(region.ERR_JOIN_ASSEMBLY_FAIL_TARGET_NOT_EXIST)

	case realmerr.ErrJoinAssemblySelfNoBase:
		hc.Send(region.ERR_JOIN_ASSEMBLY_FAIL_SELF_NOT_EXIST)

	case realmerr.ErrJoinAssemblyEmptyGeneral:
		hc.Send(region.ERR_JOIN_ASSEMBLY_FAIL_INVALID_TROOP_INDEX)

	case realmerr.ErrJoinAssemblyGeneralOutside:
		hc.Send(region.ERR_JOIN_ASSEMBLY_FAIL_OUTSIDE)

	case realmerr.ErrJoinAssemblyNoSoldier:
		hc.Send(region.ERR_JOIN_ASSEMBLY_FAIL_NO_SOLDIER)

	case realmerr.ErrJoinAssemblyInvalidRelation:
		hc.Send(region.ERR_JOIN_ASSEMBLY_FAIL_INVALID_TARGET)

	case realmerr.ErrJoinAssemblyInvalidTroopIndex:
		hc.Send(region.ERR_JOIN_ASSEMBLY_FAIL_INVALID_TROOP_INDEX)

	case realmerr.ErrJoinAssemblyTodayJoinXiongNu:
		hc.Send(region.ERR_JOIN_ASSEMBLY_FAIL_TODAY_JOIN_XIONG_NU)

	default:
		logrus.WithError(err).Error("region.ProcessJoinAssembly有err没有处理")
		hc.Send(region.ERR_JOIN_ASSEMBLY_FAIL_SERVER_ERROR)
	}
}

// 联盟工坊

//gogen:iface
func (m *RegionModule) ProcessCreateGuildWorkshop(proto *region.C2SCreateGuildWorkshopProto, hc iface.HeroController) {

	r := m.realmService.GetBigMap()

	// 检查坐标
	posX, posY := int(proto.PosX), int(proto.PosY)
	if !r.GetMapData().IsValidHomePosition(posX, posY) {
		logrus.Debugf("创建联盟工坊，无效的主城坐标, %v,%v", posX, posY)
		hc.Send(region.ERR_CREATE_GUILD_WORKSHOP_FAIL_INVALID_POS)
		return
	}

	if !r.IsPosOpened(posX, posY) {
		logrus.Debugf("创建联盟工坊，坐标还未开放, %v,%v", posX, posY)
		hc.Send(region.ERR_CREATE_GUILD_WORKSHOP_FAIL_INVALID_POS)
		return
	}

	if r.IsEdgeNotHomePos(posX, posY) {
		logrus.Debugf("创建联盟工坊，边界坐标, %v,%v", posX, posY)
		hc.Send(region.ERR_CREATE_GUILD_WORKSHOP_FAIL_INVALID_POS)
		return
	}

	var guildId int64
	var heroName string
	if hc.FuncWithSend(func(hero *entity.Hero, result herolock.LockResult) {
		if hero.BaseRegion() == 0 || hero.BaseLevel() <= 0 || hero.Prosperity() <= 0 {
			logrus.Debugf("创建联盟工坊，主城流亡了")
			result.Add(region.ERR_CREATE_GUILD_WORKSHOP_FAIL_DISTANCE_LIMIT)
			return
		}

		if u64.FromInt(hexagon.Distance(hero.BaseX(), hero.BaseY(), posX, posY)) > m.datas.GuildGenConfig().WorkshopDistanceLimit {
			logrus.Debugf("创建联盟工坊，距离自己主城太远")
			result.Add(region.ERR_CREATE_GUILD_WORKSHOP_FAIL_DISTANCE_LIMIT)
			return
		}

		guildId = hero.GuildId()
		heroName = hero.Name()

		result.Ok()
	}) {
		return
	}

	if guildId == 0 {
		logrus.Debug("创建联盟工坊，英雄没有联盟")
		hc.Send(region.ERR_CREATE_GUILD_WORKSHOP_FAIL_NOT_IN_GUILD)
		return
	}

	// 权限
	var errMsg pbutil.Buffer
	var classLevelName string
	m.dep.Guild().FuncGuild(guildId, func(g *sharedguilddata.Guild) {
		if g == nil {
			logrus.Errorf("创建联盟工坊，g == nil")
			errMsg = region.ERR_CREATE_GUILD_WORKSHOP_FAIL_NOT_IN_GUILD
			return
		}

		member := g.GetMember(hc.Id())
		if member == nil {
			logrus.Errorf("创建联盟工坊，member == nil")
			errMsg = region.ERR_CREATE_GUILD_WORKSHOP_FAIL_NOT_IN_GUILD
			return
		}

		if !member.HasPermission(func(permission *guild_data.GuildPermissionData) bool {
			return permission.Workshop
		}) {
			logrus.Debug("创建联盟工坊，你没有权限")
			errMsg = region.ERR_CREATE_GUILD_WORKSHOP_FAIL_DENY
			return
		}

		if g.GetWorkshop() != nil {
			baseId := npcid.NewGuildWorkshopId(guildId)
			if b := r.GetRoBase(baseId); b != nil {
				logrus.Debug("创建联盟工坊，联盟工坊已经存在")
				errMsg = region.ERR_CREATE_GUILD_WORKSHOP_FAIL_BASE_EXIST
				return
			}
		}
		classLevelName = member.ClassLevelData().Name
	})

	if errMsg != nil {
		hc.Send(errMsg)
		return
	}

	if !r.ReservePos(posX, posY) {
		logrus.Debug("创建联盟工坊，申请坐标失败，太近了")
		hc.Send(region.ERR_CREATE_GUILD_WORKSHOP_FAIL_INVALID_POS)
		return
	}

	ctime := m.time.CurrentTime()
	startTime := timeutil.Marshal32(ctime)
	endTime := startTime + timeutil.DurationMarshal32(m.datas.GuildGenConfig().WorkshopBuildDuration)

	processed := r.AddGuildWorkshop(guildId, posX, posY, startTime, endTime)
	if !processed {
		logrus.Errorf("创建联盟工坊，消息没处理")
		hc.Send(region.ERR_CREATE_GUILD_WORKSHOP_FAIL_SERVER_ERROR)
		return
	}

	hc.Send(region.NewS2cCreateGuildWorkshopMsg(proto.PosX, proto.PosY))

	// add log
	if d := m.datas.TextHelp().GuildGongfangBuilding; d != nil {
		text := d.New().WithClassName(classLevelName)
		text.WithClickHeroFields(data.KeyName, heroName, hc.Id()).JsonString()
		m.dep.Guild().FuncGuild(guildId, func(g *sharedguilddata.Guild) {
			if g == nil {
				return
			}
			if g.GetWorkshop() != nil {
				g.GetWorkshop().AddLog(text.JsonString(), ctime)
			}
		})
	}

	if data := m.datas.TextHelp().GuildWorkshopCreatedChat; data != nil {
		baseId := npcid.NewGuildWorkshopId(guildId)
		showJson := fmt.Sprintf(`{"workshop_id":%d,"pos_x":%d,"pos_y":%d}`, baseId, posX, posY)
		m.dep.Chat().SysChat(hc.Id(), guildId, shared_proto.ChatType_ChatGuild, showJson, shared_proto.ChatMsgType_ChatMsgGuildWorkshopCreated, false, true, false, true)
	}
}

var emptyGuildWorkshopMsg = region.NewS2cShowGuildWorkshopMsg(nil, 0, 0, 0, 0, 0).Static()

//gogen:iface
func (m *RegionModule) ProcessShowGuildWorkshop(proto *region.C2SShowGuildWorkshopProto, hc iface.HeroController) {

	targetId, ok := idbytes.ToId(proto.BaseId)
	if !ok {
		logrus.Errorf("查看联盟工坊，无效的id")
		hc.Send(emptyGuildWorkshopMsg)
		return
	}

	targetGuildId := npcid.GetWorkshopGuildId(targetId)

	var toSend pbutil.Buffer
	m.dep.Guild().FuncGuild(targetGuildId, func(g *sharedguilddata.Guild) {
		if g == nil {
			logrus.Errorf("查看联盟工坊，g == nil")
			toSend = emptyGuildWorkshopMsg
			return
		}

		w := g.GetWorkshop()
		if w == nil {
			logrus.Errorf("查看联盟工坊，w == nil")
			toSend = emptyGuildWorkshopMsg
			return
		}

		var prizeCount int32
		if member := g.GetMember(hc.Id()); member != nil {
			prizeCount = u64.Int32(member.GetWorkshopPrizeCount())
		}

		var output, totalOutput uint64
		if g.GetWorkshopOutputPrizeCount() < len(m.datas.GuildGenConfig().WorkshopMaxOutput) {
			output = g.GetWorkshopOutput()
			totalOutput = m.datas.GuildGenConfig().WorkshopMaxOutput[g.GetWorkshopOutputPrizeCount()]
		}
		g.GetWorkshopOutput()

		toSend = region.NewS2cShowGuildWorkshopMsg(proto.BaseId, i64.Int32(targetGuildId),
			u64.Int32(output), u64.Int32(totalOutput), prizeCount, u64.Int32(g.GetTodayWorkshopBeenHurtTimes()))
	})

	hc.Send(toSend)
}

//gogen:iface
func (m *RegionModule) ProcessHurtGuildWorkshop(proto *region.C2SHurtGuildWorkshopProto, hc iface.HeroController) {

	targetId, ok := idbytes.ToId(proto.BaseId)
	if !ok {
		logrus.Errorf("破坏联盟工坊，无效的id")
		hc.Send(region.ERR_HURT_GUILD_WORKSHOP_FAIL_INVALID_BASE_ID)
		return
	}

	targetGuildId := npcid.GetWorkshopGuildId(targetId)

	ctime := m.time.CurrentTime()

	if hc.FuncWithSend(func(hero *entity.Hero, result herolock.LockResult) {

		if hero.GuildWorkshop().GetDailyHurtTimes() >= m.datas.GuildGenConfig().WorkshopHurtHeroTimesLimit {
			logrus.Errorf("破坏联盟工坊，破坏次数已达上限")
			result.Add(region.ERR_HURT_GUILD_WORKSHOP_FAIL_HURT_TIMES_NOT_ENOUGH)
			return
		}

		if ctime.Before(hero.GuildWorkshop().GetNextHurtTime()) {
			logrus.Debug("破坏联盟工坊，CD未到")
			result.Add(region.ERR_HURT_GUILD_WORKSHOP_FAIL_COOLDOWN)
			return
		}

		if hero.GuildId() == targetGuildId {
			logrus.Errorf("破坏联盟工坊，自己的联盟？")
			result.Add(region.ERR_HURT_GUILD_WORKSHOP_FAIL_INVALID_BASE_ID)
			return
		}

		result.Ok()
	}) {
		return
	}

	r := m.realmService.GetBigMap()

	var errMsg pbutil.Buffer
	m.dep.Guild().FuncGuild(targetGuildId, func(g *sharedguilddata.Guild) {
		if g == nil {
			logrus.Errorf("破坏联盟工坊，g == nil")
			errMsg = region.ERR_HURT_GUILD_WORKSHOP_FAIL_INVALID_BASE_ID
			return
		}

		if g.GetWorkshop() == nil {
			logrus.Debug("破坏联盟工坊，g.workshop == nil")
			errMsg = region.ERR_HURT_GUILD_WORKSHOP_FAIL_INVALID_BASE_ID
			return
		}

		if g.GetTodayWorkshopBeenHurtTimes() >= m.datas.GuildGenConfig().WorkshopHurtTotalTimesLimit {
			logrus.Debug("破坏联盟工坊，目标联盟被破坏次数已达上限")
			errMsg = region.ERR_HURT_GUILD_WORKSHOP_FAIL_BEEN_HURT_TIMES_LIMIT
			return
		}

	})

	if errMsg != nil {
		hc.Send(errMsg)
		return
	}

	processed, baseNotExist := r.HurtGuildWorkshop(hc, targetGuildId)
	if !processed {
		logrus.Debug("破坏联盟工坊，processed = false")
		hc.Send(region.ERR_HURT_GUILD_WORKSHOP_FAIL_SERVER_ERROR)
		return
	}

	if baseNotExist {
		hc.Send(region.ERR_HURT_GUILD_WORKSHOP_FAIL_INVALID_BASE_ID)
		return
	} else {
		var heroName string
		hc.FuncWithSend(func(hero *entity.Hero, result herolock.LockResult) {

			heroName = hero.Name()
			newHurtTimes := hero.GuildWorkshop().IncDailyHurtTimes()

			nextHurtTime := ctime.Add(m.datas.GuildGenConfig().WorkshopHurtCooldown)
			hero.GuildWorkshop().SetNextHurtTime(nextHurtTime)

			result.Add(region.NewS2cHurtGuildWorkshopMsg(proto.BaseId, timeutil.Marshal32(nextHurtTime), u64.Int32(newHurtTimes)))
			result.Ok()
		})

		m.dep.Guild().FuncGuild(targetGuildId, func(g *sharedguilddata.Guild) {
			if g == nil {
				return
			}
			workshop := g.GetWorkshop()
			if workshop == nil {
				return
			}
			// add log
			if workshop.IsComplete() {
				if d := m.datas.TextHelp().GuildGongfangEfficiencyReduce; d != nil {
					text := d.New().WithClickHeroFields(data.KeyName, heroName, hc.Id())
					text.WithNum(m.datas.GuildGenConfig().WorkshopHurtProsperity)
					workshop.AddLog(text.JsonString(), ctime)
				}
			} else if d := m.datas.TextHelp().GuildGongfangBuildTimeAdd; d != nil {
				text := d.New().WithClickHeroFields(data.KeyName, heroName, hc.Id())
				minutes := int64(m.datas.GuildGenConfig().WorkshopHurtDuration / time.Minute)
				text.WithNum(minutes)
				workshop.AddLog(text.JsonString(), ctime)
			}
		})
	}
}

var showWorkshopTrueMsg = guild.NewS2cShowWorkshopNotExistMsg(true).Static()

//gogen:iface c2s_remove_guild_workshop
func (m *RegionModule) ProcessRemoveGuildWorkshop(hc iface.HeroController) {

	guildId, _ := hc.LockGetGuildId()
	if guildId == 0 {
		logrus.Debug("移除联盟工坊，英雄没有联盟")
		hc.Send(region.ERR_REMOVE_GUILD_WORKSHOP_FAIL_NOT_IN_GUILD)
		return
	}

	r := m.realmService.GetBigMap()
	baseId := npcid.NewGuildWorkshopId(guildId)

	// 权限
	var errMsg pbutil.Buffer
	m.dep.Guild().FuncGuild(guildId, func(g *sharedguilddata.Guild) {
		if g == nil {
			logrus.Errorf("移除联盟工坊，g == nil")
			errMsg = region.ERR_REMOVE_GUILD_WORKSHOP_FAIL_NOT_IN_GUILD
			return
		}

		member := g.GetMember(hc.Id())
		if member == nil {
			logrus.Errorf("移除联盟工坊，member == nil")
			errMsg = region.ERR_REMOVE_GUILD_WORKSHOP_FAIL_NOT_IN_GUILD
			return
		}

		if !member.HasPermission(func(permission *guild_data.GuildPermissionData) bool {
			return permission.Workshop
		}) {
			logrus.Debug("移除联盟工坊，你没有权限")
			errMsg = region.ERR_REMOVE_GUILD_WORKSHOP_FAIL_DENY
			return
		}

		if g.GetWorkshop() == nil {
			if b := r.GetRoBase(baseId); b == nil {
				logrus.Debug("移除联盟工坊，联盟工坊不存在")
				errMsg = region.ERR_REMOVE_GUILD_WORKSHOP_FAIL_BASE_NOT_EXIST
				return
			}
		}
	})

	if errMsg != nil {
		hc.Send(errMsg)
		return
	}

	processed, err, _, _ := r.RemoveBase(baseId, false, nil, nil)
	if !processed {
		logrus.Errorf("移除联盟工坊，消息没处理")
		hc.Send(region.ERR_REMOVE_GUILD_WORKSHOP_FAIL_SERVER_ERROR)
		return
	}

	if err != nil {
		logrus.WithError(err).Errorf("移除联盟工坊，野外返回报错")
		hc.Send(region.ERR_REMOVE_GUILD_WORKSHOP_FAIL_SERVER_ERROR)
		return
	}

	var memberIds []int64
	m.dep.Guild().FuncGuild(guildId, func(g *sharedguilddata.Guild) {
		if g == nil {
			logrus.Errorf("移除联盟工坊，g == nil（刚好这么巧？）")
			errMsg = region.ERR_REMOVE_GUILD_WORKSHOP_FAIL_NOT_IN_GUILD
			return
		}

		g.SetWorkshop(nil)

		memberIds = g.AllUserMemberIds()
	})

	hc.Send(region.REMOVE_GUILD_WORKSHOP_S2C)

	if len(memberIds) > 0 {
		m.world.MultiSend(memberIds, showWorkshopTrueMsg)
	}
}

//gogen:iface c2s_catch_guild_workshop_logs
func (m *RegionModule) ProcessCatchGuildWorkshopLogs(proto *region.C2SCatchGuildWorkshopLogsProto, hc iface.HeroController) {
	guildId, _ := hc.LockGetGuildId()
	if guildId == 0 {
		logrus.Debug("获取联盟工坊日志，没有联盟")
		hc.Send(region.ERR_CATCH_GUILD_WORKSHOP_LOGS_FAIL_NO_GUILD)
		return
	}

	r := m.realmService.GetBigMap()
	baseId := npcid.NewGuildWorkshopId(guildId)

	var sendMsg pbutil.Buffer
	m.dep.Guild().FuncGuild(guildId, func(g *sharedguilddata.Guild) {
		if g == nil {
			logrus.Errorf("获取联盟工坊日志，g == nil")
			sendMsg = region.ERR_CATCH_GUILD_WORKSHOP_LOGS_FAIL_NO_GUILD
			return
		}
		workShop := g.GetWorkshop()
		if workShop == nil {
			if b := r.GetRoBase(baseId); b == nil {
				logrus.Debug("获取联盟工坊日志，没有联盟工坊")
				sendMsg = region.ERR_CATCH_GUILD_WORKSHOP_LOGS_FAIL_NOT_EXIST
				return
			}
		}
		sendMsg = workShop.GetMsg(u64.FromInt32(proto.GetVersion()))
	})

	hc.Send(sendMsg)
}

var selfBaozNotExist = region.NewS2cGetSelfBaozMsg(false, nil, 0, 0, 0).Static()

//gogen:iface c2s_get_self_baoz
func (m *RegionModule) ProcessGetSelfBaoz(hc iface.HeroController) {

	base := m.realmService.GetBigMap().GetHeroBaozRoBase(hc.Id())
	if base == nil {
		hc.Send(selfBaozNotExist)
		return
	}

	hc.Send(region.NewS2cGetSelfBaozMsg(true, idbytes.ToBytes(base.Id), base.X, base.Y, base.HeroEndTime))
}
