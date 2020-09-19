package hebi

import (
	"github.com/lightpaw/male7/gen/iface"
	"github.com/lightpaw/male7/gen/pb/hebi"
	"github.com/lightpaw/male7/service/conflict/heroservice/herolock"
	"github.com/lightpaw/male7/entity"
	"github.com/lightpaw/male7/util/u64"
	"github.com/lightpaw/male7/pb/shared_proto"
	"github.com/lightpaw/male7/service/heromodule"
	"github.com/lightpaw/logrus"
	"github.com/lightpaw/male7/pb/server_proto"
	"github.com/lightpaw/male7/config/goods"
	"github.com/lightpaw/male7/service/conflict/guilddataservice/guildsnapshotdata"
	"github.com/lightpaw/male7/service/operate_type"
	hebiConf "github.com/lightpaw/male7/config/hebi"
)

func NewHebiModule(dep iface.ServiceDep, mail iface.MailModule, fightService iface.FightService, ticker iface.TickerService) *HebiModule {
	m := &HebiModule{}
	m.dep = dep
	m.guildSnapshot = dep.GuildSnapshot()
	m.heroService = dep.HeroData()
	m.fightService = fightService
	m.mail = mail
	m.hebiManager = NewHebiManager(dep, mail, ticker)

	return m
}

//gogen:iface
type HebiModule struct {
	dep           iface.ServiceDep
	guildSnapshot iface.GuildSnapshotService
	heroService   iface.HeroDataService
	fightService  iface.FightService
	mail          iface.MailModule
	hebiManager   *HebiManager
}

func (m *HebiModule) Close() {
	m.hebiManager.close()
}

//gogen:iface
func (m *HebiModule) ProcessRoomList(proto *hebi.C2SRoomListProto, hc iface.HeroController) {
	v := proto.V
	if m.hebiManager.msgVersion.Load() == u64.FromInt32(v) {
		hc.Send(m.hebiManager.emptyMsg)
		return
	}

	hc.Send(m.hebiManager.msg)
}

//gogen:iface
func (m *HebiModule) ProcessChangeCaptain(proto *hebi.C2SChangeCaptainProto, hc iface.HeroController) {

	var captainProto *shared_proto.HebiCaptainProto
	if hc.FuncWithSend(func(hero *entity.Hero, result herolock.LockResult) {
		cid := u64.FromInt32(proto.CaptainId)
		captain := hero.Military().Captain(cid)
		if captain == nil {
			result.Add(hebi.ERR_CHANGE_CAPTAIN_FAIL_INVALID_CAPTAIN_ID)
			return
		}

		heroHebi := hero.Hebi()
		if heroHebi.CaptainId == cid {
			result.Add(hebi.ERR_CHANGE_CAPTAIN_FAIL_INVALID_CAPTAIN_ID)
			return
		}

		heroHebi.CaptainId = cid
		captainProto = captain.EncodeHebiCaptain()

		result.Add(hebi.NewS2cChangeCaptainMsg(proto.CaptainId))
		result.Changed()
		result.Ok()
	}) {
		return
	}

	if captainProto == nil {
		return
	}

	// 房间换武将
	m.hebiManager.Func("ProcessChangeCaptain", func() {
		room, exist := m.hebiManager.HeroRoom(hc.Id())
		if !exist {
			return
		}

		if room.ChangeCaptain(hc.Id(), captainProto) {
			ctime := m.dep.Time().CurrentTime()
			m.hebiManager.UpdateMsg(ctime)
			hc.Send(hebi.NewS2cChangeRoomCaptainMsg(u64.Int32(room.roomId), captainProto))
		}
	}, emptyFunc)
}

//gogen:iface
func (m *HebiModule) ProcessCheckInRoom(proto *hebi.C2SCheckInRoomProto, hc iface.HeroController) {
	roomId := u64.FromInt32(proto.RoomId)
	goodsId := u64.FromInt32(proto.GoodsId)

	if proto.RoomId < 0 || roomId >= m.dep.Datas().HebiMiscData().RoomsMaxSize {
		hc.Send(hebi.ERR_CHECK_IN_ROOM_FAIL_INVALID_ROOM_ID)
		return
	}

	goodsData := m.dep.Datas().GetGoodsData(goodsId)
	if goodsData == nil || goodsData.SpecType != shared_proto.GoodsSpecType_GAT_HEBI {
		hc.Send(hebi.ERR_CHECK_IN_ROOM_FAIL_GOODS_INVALID)
		return
	}

	hctx := heromodule.NewContext(m.dep, operate_type.HebiCheckInRoom)

	m.hebiManager.Func("ProcessNewRoom", func() {
		// 验证房间
		r := m.hebiManager.GetRoom(roomId)
		if r != nil {
			if r.state == shared_proto.HebiRoomState_HebiRoomRobProtect {
				if r.hebiRobProtectId != hc.Id() {
					logrus.Debugf("创建房间，房间在保护状态，但不是自己抢来的")
					hc.Send(hebi.ERR_CHECK_IN_ROOM_FAIL_ROOM_ID_NOT_EMPTY)
					return
				}
			} else if r.state != shared_proto.HebiRoomState_HebiRoomEmpty {
				logrus.Debugf("创建房间，房间不在保护状态，但不是空的")
				hc.Send(hebi.ERR_CHECK_IN_ROOM_FAIL_ROOM_ID_NOT_EMPTY)
				return
			}
		}

		if m.hebiManager.InAnyRoom(hc.Id()) {
			hc.Send(hebi.ERR_CHECK_IN_ROOM_FAIL_ALREADY_IN_ROOM)
			return
		}

		var heroBasicProto *shared_proto.HeroBasicProto
		if heroSnapshot := m.dep.HeroSnapshot().Get(hc.Id()); heroSnapshot != nil {
			heroBasicProto = heroSnapshot.EncodeBasic4Client()
		} else {
			logrus.Debugf("ProcessCheckInRoom, 找不到 herosnapshot，heroid:%v", hc.Id())
			hc.Send(hebi.ERR_CHECK_IN_ROOM_FAIL_SERVER_ERR)
			return
		}

		var captainProto *shared_proto.HebiCaptainProto
		if hc.FuncWithSend(func(hero *entity.Hero, result herolock.LockResult) {
			// 武将
			c := hero.Military().Captain(hero.Hebi().CaptainId)
			if c == nil {
				result.Add(hebi.ERR_CHECK_IN_ROOM_FAIL_CAPTAIN_ERR)
				return
			}
			captainProto = c.EncodeHebiCaptain()

			// 扣物品
			if !heromodule.TryReduceGoods(hctx, hero, result, goodsData, 1) {
				result.Add(hebi.ERR_CHECK_IN_ROOM_FAIL_GOODS_INVALID)
				return
			}
			hero.Hebi().CurrentGoodsId = goodsId

			result.Changed()
			result.Ok()
		}) {
			return
		}

		// 创建房间
		ctime := m.dep.Time().CurrentTime()
		room := m.hebiManager.CheckInHebiRoom(roomId, hc.Id(), heroBasicProto, captainProto, goodsData, ctime)
		if room == nil {
			logrus.Debugf("合璧ProcessNewRoom，CheckInHebiRoom 失败")
			hc.Send(hebi.ERR_CHECK_IN_ROOM_FAIL_SERVER_ERR)
			return
		}
		m.hebiManager.UpdateMsg(ctime)

		hc.Send(hebi.NewS2cCheckInRoomMsg(proto.RoomId, proto.GoodsId))
	}, func() {
		hc.Send(hebi.ERR_CHECK_IN_ROOM_FAIL_SERVER_ERR)
	})
}

//gogen:iface
func (m *HebiModule) ProcessJoinRoom(proto *hebi.C2SJoinRoomProto, hc iface.HeroController) {
	roomId := u64.FromInt32(proto.RoomId)
	goodsId := u64.FromInt32(proto.GoodsId)

	if proto.RoomId < 0 || roomId >= m.dep.Datas().HebiMiscData().RoomsMaxSize {
		hc.Send(hebi.ERR_JOIN_ROOM_FAIL_INVALID_ROOM_ID)
		return
	}

	goodsData := m.dep.Datas().GetGoodsData(goodsId)
	if goodsData == nil || goodsData.SpecType != shared_proto.GoodsSpecType_GAT_HEBI {
		hc.Send(hebi.ERR_JOIN_ROOM_FAIL_GOODS_INVALID)
		return
	}

	hctx := heromodule.NewContext(m.dep, operate_type.HebiJoinRoom)

	var hostId, guestId int64
	var hostGoods, guestGoods *goods.GoodsData
	m.hebiManager.Func("ProcessJoinRoom", func() {
		// 验证房间
		room := m.hebiManager.GetRoom(roomId)
		if room == nil || room.state != shared_proto.HebiRoomState_HebiRoomWait {
			hc.Send(hebi.ERR_JOIN_ROOM_FAIL_INVALID_ROOM_ID)
			return
		}

		if m.hebiManager.InAnyRoom(hc.Id()) {
			hc.Send(hebi.ERR_JOIN_ROOM_FAIL_ALREADY_IN_ROOM)
			return
		}

		if goodsData.HebiSubType == room.HostPos() {
			hc.Send(hebi.ERR_JOIN_ROOM_FAIL_GOODS_INVALID)
			return
		}

		if room.hostGoodsData.HebiType != goodsData.HebiType {
			logrus.Debugf("ProcessJoinRoom, 玉璧不是同一种类型")
			hc.Send(hebi.ERR_JOIN_ROOM_FAIL_GOODS_INVALID)
			return
		}

		if room.hostGoodsData.HebiType == m.dep.Datas().HebiMiscData().HeShiBiType {
			guildId, ok := hc.LockGetGuildId()
			if !ok || guildId != int64(room.host.GuildId) {
				logrus.Debugf("ProcessJoinRoom, 和氏璧，不是同一联盟")
				hc.Send(hebi.ERR_JOIN_ROOM_FAIL_HESHIBI_NOT_SAME_GUILD)
				return
			}
		}

		var heroBasicProto *shared_proto.HeroBasicProto
		if heroSnapshot := m.dep.HeroSnapshot().Get(hc.Id()); heroSnapshot != nil {
			heroBasicProto = heroSnapshot.EncodeBasic4Client()
		} else {
			logrus.Debugf("合璧 ProcessJoinRoom, 找不到 herosnapshot，heroid:%v", hc.Id())
			hc.Send(hebi.ERR_JOIN_ROOM_FAIL_SERVER_ERR)
			return
		}

		var captainProto *shared_proto.HebiCaptainProto
		if hc.FuncWithSend(func(hero *entity.Hero, result herolock.LockResult) {
			// 武将
			c := hero.Military().Captain(hero.Hebi().CaptainId)
			if c == nil {
				result.Add(hebi.ERR_JOIN_ROOM_FAIL_CAPTAIN_ERR)
				return
			}
			captainProto = c.EncodeHebiCaptain()

			// 扣物品
			if !heromodule.TryReduceGoods(hctx, hero, result, goodsData, 1) {
				result.Add(hebi.ERR_JOIN_ROOM_FAIL_GOODS_INVALID)
				return
			}
			hero.Hebi().CurrentGoodsId = goodsId

			result.Changed()
			result.Ok()
		}) {
			return
		}

		// 加入房间
		ctime := m.dep.Time().CurrentTime()
		if !room.Join(hc.Id(), heroBasicProto, captainProto, goodsData, ctime) {
			logrus.Debugf("合璧 ProcessJoinRoom，room.Join 失败")
			hc.Send(hebi.ERR_JOIN_ROOM_FAIL_SERVER_ERR)
			return
		}
		m.hebiManager.UpdateMsg(ctime)

		hc.Send(hebi.NewS2cJoinRoomMsg(proto.RoomId, proto.GoodsId, u64.Int32(room.prizeId)))

		m.dep.World().Send(room.hostId, hebi.NewS2cSomeoneJoinedRoomMsg(proto.RoomId))

		hostId, guestId = room.GetHostId(), room.GetGuestId()
		hostGoods, guestGoods = room.GetHostGoods(), room.GetGuestGoods()
	}, func() {
		hc.Send(hebi.ERR_JOIN_ROOM_FAIL_SERVER_ERR)
	})

	// 完成合璧任务
	if hostId != 0 && guestId != 0 {
		completeFunc := func(heroId int64, goodsData *goods.GoodsData) {
			if heroId == 0 || goodsData == nil {
				return
			}

			m.heroService.FuncWithSend(heroId, func(hero *entity.Hero, result herolock.LockResult) {
				hero.HistoryAmount().IncreaseOneWithSubType(server_proto.HistoryAmountType_StartHebi, uint64(goodsData.Quality))
				hero.HistoryAmount().IncreaseOne(server_proto.HistoryAmountType_StartHebi)

				heromodule.UpdateTaskProgress(hero, result, shared_proto.TaskTargetType_TASK_TARGET_HEBI)
			})
		}

		completeFunc(hostId, hostGoods)
		if hostId != guestId {
			completeFunc(guestId, guestGoods)
		}
	}
}

//gogen:iface
func (m *HebiModule) ProcessCopySelf(proto *hebi.C2SCopySelfProto, hc iface.HeroController) {
	roomId := u64.FromInt32(proto.RoomId)
	goodsId := u64.FromInt32(proto.GoodsId)

	if proto.RoomId < 0 || roomId >= m.dep.Datas().HebiMiscData().RoomsMaxSize {
		hc.Send(hebi.ERR_COPY_SELF_FAIL_INVALID_ROOM_ID)
		return
	}

	goodsData := m.dep.Datas().GetGoodsData(goodsId)
	if goodsData == nil || goodsData.SpecType != shared_proto.GoodsSpecType_GAT_HEBI_COPY_SELF {
		hc.Send(hebi.ERR_COPY_SELF_FAIL_GOODS_INVALID)
		return
	}

	hctx := heromodule.NewContext(m.dep, operate_type.HebiCopySelf)

	m.hebiManager.Func("ProcessCopySelf", func() {
		// 验证房间
		room := m.hebiManager.GetRoom(roomId)
		if room.hostId != hc.Id() {
			hc.Send(hebi.ERR_COPY_SELF_FAIL_INVALID_ROOM_ID)
			return
		}
		if room.state != shared_proto.HebiRoomState_HebiRoomWait {
			hc.Send(hebi.ERR_COPY_SELF_FAIL_ROOM_STATE_NOT_WAIT)
			return
		}

		var captainProto *shared_proto.HebiCaptainProto
		if hc.FuncWithSend(func(hero *entity.Hero, result herolock.LockResult) {
			// 扣物品
			if !heromodule.TryReduceGoods(hctx, hero, result, goodsData, 1) {
				result.Add(hebi.ERR_COPY_SELF_FAIL_GOODS_INVALID)
				return
			}

			hero.HistoryAmount().IncreaseOneWithSubType(server_proto.HistoryAmountType_StartHebi, uint64(room.hostGoodsData.Quality))
			hero.HistoryAmount().IncreaseOne(server_proto.HistoryAmountType_StartHebi)

			heromodule.UpdateTaskProgress(hero, result, shared_proto.TaskTargetType_TASK_TARGET_HEBI)

			// 武将
			if c := hero.Military().Captain(hero.Hebi().CaptainId); c != nil {
				captainProto = c.EncodeHebiCaptain()
			}

			result.Changed()
			result.Ok()
		}) {
			hc.Send(hebi.ERR_COPY_SELF_FAIL_SERVER_ERR)
			return
		}

		// 另一半玉璧
		partnerGoodsData := room.hostGoodsData.PartnerHebiGoodsData
		ctime := m.dep.Time().CurrentTime()
		room.CopySelf(ctime, captainProto, partnerGoodsData)

		m.hebiManager.UpdateMsg(ctime)
		hc.Send(hebi.NewS2cCopySelfMsg(proto.RoomId))
	}, func() {
		hc.Send(hebi.ERR_COPY_SELF_FAIL_SERVER_ERR)
	})
}

//gogen:iface
func (m *HebiModule) ProcessLeave(proto *hebi.C2SLeaveRoomProto, hc iface.HeroController) {
	hctx := heromodule.NewContext(m.dep, operate_type.HebiLeave)

	roomId := u64.FromInt32(proto.RoomId)
	m.hebiManager.Func("ProcessLeave", func() {
		// 验证房间
		room := m.hebiManager.GetRoom(roomId)
		if room == nil || room.host == nil {
			hc.Send(hebi.ERR_LEAVE_ROOM_FAIL_NOT_IN_ROOM)
			return
		}
		if room.hostId != hc.Id() {
			hc.Send(hebi.ERR_LEAVE_ROOM_FAIL_NOT_IN_ROOM)
			return
		}
		if room.state != shared_proto.HebiRoomState_HebiRoomWait {
			hc.Send(hebi.ERR_LEAVE_ROOM_FAIL_INVALID_ROOM_STATE)
			return
		}

		succ, goodsData := room.Leave()
		if !succ {
			hc.Send(hebi.ERR_LEAVE_ROOM_FAIL_INVALID_ROOM_STATE)
			return
		}

		ctime := m.dep.Time().CurrentTime()
		m.hebiManager.UpdateMsg(ctime)
		hc.Send(hebi.NewS2cLeaveRoomMsg(proto.RoomId))
		// 物品补回去
		hc.FuncWithSend(func(hero *entity.Hero, result herolock.LockResult) {
			heromodule.AddGoods(hctx, hero, result, goodsData, 1)
			result.Changed()
			result.Ok()
		})

	}, func() {
		hc.Send(hebi.ERR_LEAVE_ROOM_FAIL_SERVER_ERR)
	})
}

//gogen:iface
func (m *HebiModule) ProcessRobPos(proto *hebi.C2SRobPosProto, hc iface.HeroController) {
	roomId := u64.FromInt32(proto.RoomId)
	var defenserId int64
	var defenserMail *shared_proto.MailProto
	m.hebiManager.Func("ProcessRobPos", func() {
		// 验证房间
		room := m.hebiManager.GetRoom(roomId)
		if room == nil || room.host == nil {
			hc.Send(hebi.ERR_ROB_POS_FAIL_INVALID_ROOM_ID)
			return
		}

		if room.state != shared_proto.HebiRoomState_HebiRoomWait {
			hc.Send(hebi.ERR_ROB_POS_FAIL_INVALID_ROOM_STATE)
			return
		}

		// 是否已经在房间里
		if m.hebiManager.InAnyRoom(hc.Id()) {
			hc.Send(hebi.ERR_ROB_POS_FAIL_ALREADY_IN_ROOM)
			return
		}

		miscData := m.dep.Datas().HebiMiscData()
		ctime := m.dep.Time().CurrentTime()

		var hostGuildId int64
		hostHero := m.dep.HeroSnapshot().Get(room.hostId)
		if hostHero != nil {
			hostGuildId = hostHero.GuildId
		}

		hctx := heromodule.NewContext(m.dep, operate_type.HebiBeRobPos)

		// 战斗
		var atkCaptainProto *shared_proto.HebiCaptainProto
		var attacker *shared_proto.CombatPlayerProto
		if hc.FuncWithSend(func(hero *entity.Hero, result herolock.LockResult) {
			if ctime.Before(hero.Hebi().RobPosCdEndTime) {
				result.Add(hebi.ERR_ROB_POS_FAIL_IN_ROB_CD)
				return
			}

			if hero.GuildId() > 0 && hero.GuildId() == hostGuildId {
				result.Add(hebi.ERR_ROB_POS_FAIL_SAME_GUILD)
				return
			}

			var hasYubi bool
			for _, g := range miscData.HebiGoods {
				if hero.Depot().HasEnoughGoods(g.Id, 1) {
					// 找到一个玉璧
					hasYubi = true
					break
				}
			}
			if !hasYubi {
				result.Add(hebi.ERR_ROB_POS_FAIL_NO_YUBI)
				return
			}

			captain := hero.Military().Captain(hero.Hebi().CaptainId)
			if captain == nil {
				result.Add(hebi.ERR_ROB_POS_FAIL_CAPTAIN_ERR)
				return
			}

			attacker = hero.GenCombatPlayerProtoWithCaptains(true, []*entity.TroopPos{nil, nil, captain.NewTroopPos(0)}, m.guildSnapshot.GetSnapshot)
			if attacker == nil {
				logrus.Debug("合璧抢位，GenCombatPlayerProtoWithCaptains 失败")
				result.Add(hebi.ERR_ROB_POS_FAIL_SERVER_ERR)
				return
			}

			atkCaptainProto = captain.EncodeHebiCaptain()

			hero.Hebi().RobPosCdEndTime = ctime.Add(miscData.RobPosCdDuration)
			result.Changed()
			result.Ok()
		}) {
			return
		}
		response, ok := m.fight(attacker, hc.Id(), room.hostId, room.guestId, atkCaptainProto, room.hostCaptain, false, room.hostGoodsData.HebiType, 0)
		if !ok {
			hc.Send(hebi.ERR_ROB_POS_FAIL_SERVER_ERR)
			return
		}
		if !response.AttackerWin {
			hc.Send(hebi.NewS2cRobPosMsg(proto.RoomId, false, response.Link))
			return
		}

		// 抢位成功
		oldHeroId := room.hostId
		goodsData, succ := room.RobPos(hc.Id(), ctime)
		if !succ {
			logrus.Debugf("合璧抢位，room.RobPos 失败")
			hc.Send(hebi.ERR_ROB_POS_FAIL_SERVER_ERR)
			return
		}

		m.hebiManager.UpdateMsg(ctime)
		hc.Send(hebi.NewS2cRobPosMsg(proto.RoomId, true, response.Link))

		// 物品补回去
		if goodsData != nil {
			m.dep.HeroData().FuncWithSend(oldHeroId, func(hero *entity.Hero, result herolock.LockResult) {
				heromodule.AddGoods(hctx, hero, result, goodsData, 1)
				result.Changed()
				result.Ok()
			})
		}

		defenserId = hostHero.Id
		if data := m.dep.Datas().MailHelp().HebiRoomBeenRobbed; data != nil {
			defenserMail = data.NewTextMail(shared_proto.MailType_MailNormal)

			attackerName := m.dep.Datas().MiscConfig().FlagHeroName.FormatIgnoreEmpty(attacker.Hero.GuildFlagName, attacker.Hero.Name)
			defenserMail.Text = data.NewTextFields().WithAttacker(attackerName).JsonString()
		}

	}, func() {
		hc.Send(hebi.ERR_ROB_POS_FAIL_SERVER_ERR)
	})

	if defenserId != 0 {
		m.dep.World().Send(defenserId, hebi.NewS2cSomeoneRobbedMyPosMsg(proto.RoomId))

		if defenserMail != nil {
			// 给target发邮件
			ctime := m.dep.Time().CurrentTime()
			m.mail.SendProtoMail(defenserId, defenserMail, ctime)
		}
	}
}

//gogen:iface
func (m *HebiModule) ProcessRob(proto *hebi.C2SRobProto, hc iface.HeroController) {
	roomId := u64.FromInt32(proto.RoomId)

	if proto.RoomId < 0 || roomId >= m.dep.Datas().HebiMiscData().RoomsMaxSize {
		hc.Send(hebi.ERR_ROB_FAIL_INVALID_ROOM_ID)
		return
	}

	m.hebiManager.Func("ProcessRob", func() {
		// 验证房间
		room := m.hebiManager.GetRoom(roomId)
		if room == nil || room.state != shared_proto.HebiRoomState_HebiRoomHebiRunning {
			hc.Send(hebi.ERR_ROB_FAIL_ROOM_NOT_IN_HEBI)
			return
		}

		if m.hebiManager.InAnyRoom(hc.Id()) {
			hc.Send(hebi.ERR_ROB_FAIL_ALREADY_IN_ROOM)
			return
		}

		miscData := m.dep.Datas().HebiMiscData()
		ctime := m.dep.Time().CurrentTime()

		var hostGuildId, guestGuildId int64
		hostHero := m.dep.HeroSnapshot().Get(room.hostId)
		if hostHero != nil {
			hostGuildId = hostHero.GuildId
		}
		guestHero := m.dep.HeroSnapshot().Get(room.guestId)
		if guestHero != nil {
			guestGuildId = guestHero.GuildId
		}

		// 战斗
		var atkCaptainProto *shared_proto.HebiCaptainProto
		var attacker *shared_proto.CombatPlayerProto
		if hc.FuncWithSend(func(hero *entity.Hero, result herolock.LockResult) {
			if hero.Hebi().DailyRobCount >= miscData.DailyRobCount {
				result.Add(hebi.ERR_ROB_FAIL_NO_TIMES)
				return
			}

			if ctime.Before(hero.Hebi().RobCdEndTime) {
				result.Add(hebi.ERR_ROB_FAIL_IN_CD)
				return
			}

			if hero.GuildId() > 0 {
				if hostGuildId == hero.GuildId() || guestGuildId == hero.GuildId() {
					result.Add(hebi.ERR_ROB_FAIL_IN_SAME_GUILD)
					return
				}
			}

			captain := hero.Military().Captain(hero.Hebi().CaptainId)
			if captain == nil {
				result.Add(hebi.ERR_ROB_FAIL_CAPTAIN_ERR)
				return
			}

			attacker = hero.GenCombatPlayerProtoWithCaptains(true, []*entity.TroopPos{nil, nil, captain.NewTroopPos(0)}, m.guildSnapshot.GetSnapshot)
			if attacker == nil {
				logrus.Debug("合璧抢夺，GenCombatPlayerProtoWithCaptains 失败")
				result.Add(hebi.ERR_ROB_FAIL_SERVER_ERR)
				return
			}

			atkCaptainProto = captain.EncodeHebiCaptain()

			result.Changed()
			result.Ok()
		}) {
			return
		}

		// 战斗
		firstDefenderId := room.hostId
		secondDefenderId := room.guestId
		firstCaptain := room.hostCaptain
		secondCaptain := room.guestCaptain
		defenderFightAmount := room.guestCaptain.FightAmount
		if room.hostCaptain.FightAmount > room.guestCaptain.FightAmount {
			firstDefenderId = room.guestId
			secondDefenderId = room.hostId
			firstCaptain = room.guestCaptain
			secondCaptain = room.hostCaptain
			defenderFightAmount = room.hostCaptain.FightAmount
		}

		// 先打第一个
		response1, ok := m.fight(attacker, hc.Id(), firstDefenderId, secondDefenderId, atkCaptainProto, firstCaptain, true, room.hostGoodsData.HebiType, 1)
		if !ok {
			hc.Send(hebi.ERR_ROB_FAIL_SERVER_ERR)
			return
		}
		if !response1.AttackerWin {
			hc.FuncWithSend(func(hero *entity.Hero, result herolock.LockResult) {
				// 输了才加cd
				hero.Hebi().RobCdEndTime = ctime.Add(miscData.RobCdDuration)
				result.Ok()
			})

			hc.Send(hebi.NewS2cRobMsg(proto.RoomId, false, response1.Link, false, "", nil))
			return
		}

		// 赢了再打第二个
		response2, ok := m.fight(attacker, hc.Id(), secondDefenderId, firstDefenderId, atkCaptainProto, secondCaptain, true, room.hostGoodsData.HebiType, 2)
		if !ok {
			hc.Send(hebi.ERR_ROB_FAIL_SERVER_ERR)
			return
		}
		if !response2.AttackerWin {
			hc.FuncWithSend(func(hero *entity.Hero, result herolock.LockResult) {
				// 输了才加cd
				hero.Hebi().RobCdEndTime = ctime.Add(miscData.RobCdDuration)
				result.Ok()
			})

			hc.Send(hebi.NewS2cRobMsg(proto.RoomId, true, response1.Link, false, response2.Link, nil))
			return
		}

		// 都赢了
		hostId := room.hostId
		guestId := room.guestId
		hostGoods, guestGoods := room.GetHostGoods(), room.GetGuestGoods()
		copySelf := room.copySelf
		succ, prizeId := room.Rob()
		if !succ {
			logrus.Debugf("合璧抢夺，rob.succ == false")
			hc.Send(hebi.ERR_ROB_FAIL_SERVER_ERR)
			return
		}

		robPrize := m.hebiManager.GetRobPrize(prizeId, ctime, u64.FromInt32(attacker.TotalFightAmount), u64.FromInt32(defenderFightAmount))
		if robPrize != nil {
			hc.Send(hebi.NewS2cRobMsg(proto.RoomId, true, response1.Link, true, response2.Link, robPrize.Encode()))

			quality := shared_proto.Quality_InvalidQuality
			if g := room.GetHostGoods(); g != nil {
				quality = g.Quality
			}
			if g := room.GetGuestGoods(); g != nil {
				if quality < g.Quality {
					quality = g.Quality
				}
			}

			// 奖励改成走邮件
			hc.FuncWithSend(func(hero *entity.Hero, result herolock.LockResult) {
				hero.Hebi().DailyRobCount++

				//heromodule.AddPrize(hctx, hero, result, robPrize, ctime)

				// 完成任务
				if quality > 0 {
					hero.HistoryAmount().IncreaseOneWithSubType(server_proto.HistoryAmountType_RobHebi, uint64(quality))
					hero.HistoryAmount().IncreaseOne(server_proto.HistoryAmountType_RobHebi)

					heromodule.UpdateTaskProgress(hero, result, shared_proto.TaskTargetType_TASK_TARGET_HEBI_ROB)
				}

				result.Changed()
				result.Ok()
			})
		} else {
			hc.Send(hebi.NewS2cRobMsg(proto.RoomId, true, response1.Link, true, response2.Link, nil))
		}

		attackerName := m.dep.Datas().MiscConfig().FlagHeroName.FormatIgnoreEmpty(attacker.Hero.GuildFlagName, attacker.Hero.Name)

		m.hebiManager.GiveBeRobbedPrize(hostId, guestId, hostGoods, guestGoods,
			prizeId, attackerName, u64.FromInt32(attacker.TotalFightAmount), u64.FromInt32(defenderFightAmount), ctime, copySelf)

		// 抢壁（干扰）的那个家伙给自己联盟增加进度
		m.hebiManager.updateGuildTaskProgress(hc.Id(), 0, m.dep.Datas().GuildConfig().HebiRobbedSuccessTaskProgress)

		noticeMsg := hebi.NewS2cSomeoneRobbedMyPrizeMsg(proto.RoomId)
		m.dep.World().Send(hostId, noticeMsg)
		m.dep.World().Send(guestId, noticeMsg)

		m.hebiManager.UpdateMsg(ctime)

		// 您[color=#e2e2e2]成功干扰[/color]了[color=#ff304e]{{name}}[/color]与[color=#ff304e]{{hero_name}}[/color]合璧，得到丰富奖励。
		if mailData := m.dep.Datas().MailHelp().HebiRobPrize; mailData != nil {
			mailProto := mailData.NewTextMail(shared_proto.MailType_MailNormal)
			mailProto.Prize = robPrize.Encode()
			mailProto.Text = mailData.NewTextFields().WithName(hostHero.GuildFlagName()).WithHeroName(guestHero.GuildFlagName()).JsonString()
			m.mail.SendProtoMail(hc.Id(), mailProto, ctime)
		}
	}, func() {
		hc.Send(hebi.ERR_ROB_FAIL_SERVER_ERR)
	})
}

func (m *HebiModule) fight(attacker *shared_proto.CombatPlayerProto, attackerId, defenderId, defenderPartnerId int64, atkCaptain, defCaptain *shared_proto.HebiCaptainProto, isRob bool, hebiType shared_proto.HebiType, fightNum uint64) (response *server_proto.CombatResponseServerProto, ok bool) {
	miscData := m.dep.Datas().HebiMiscData()

	var defCaptainProto *shared_proto.HebiCaptainProto
	var defender *shared_proto.CombatPlayerProto
	m.dep.HeroData().Func(defenderId, func(hero *entity.Hero, err error) (heroChanged bool) {
		captain := hero.Military().Captain(u64.FromInt32(defCaptain.Id))
		if captain == nil {
			logrus.Debugf("合璧抢夺，GenCombatPlayerProtoWithCaptains, host 设置的 captain id 不存在 id:%v", defCaptain.Id)
			return
		}
		defender = hero.GenCombatPlayerProtoWithCaptains(true, []*entity.TroopPos{nil, nil, captain.NewTroopPos(0)}, m.guildSnapshot.GetSnapshot)
		if defender == nil {
			logrus.Debug("合璧抢夺，host.GenCombatPlayerProtoWithCaptains 失败")
			return
		}

		defCaptainProto = captain.EncodeHebiCaptain()

		return
	})
	if defender == nil {
		return
	}

	tfctx := entity.NewTlogFightContext(operate_type.BattleHebi, 0, 0, 0)
	response = m.fightService.SendFightRequest(tfctx, miscData.CombatScene, attackerId, defenderId, attacker, defender)
	if response == nil {
		logrus.Errorf("合璧抢夺，战斗 response==nil")
		return
	}

	if response.ReturnCode != 0 {
		logrus.Errorf("合璧抢夺，战斗计算发生错误，%s", response.ReturnMsg)
		return
	}

	ctime := m.dep.Time().CurrentTime()
	// atk record
	m.hebiManager.AddHeroRecord(attackerId, 0, attacker.Hero, defender.Hero, atkCaptain, defCaptainProto, isRob, true, hebiType, fightNum, response, ctime)
	// def record
	m.hebiManager.AddHeroRecord(defenderId, defenderPartnerId, defender.Hero, attacker.Hero, defCaptainProto, atkCaptain, isRob, false, hebiType, fightNum, response, ctime)

	ok = true
	return
}

func (m *HebiModule) UpdateGuildInfo(heroId, guildId int64) {
	m.hebiManager.FuncNoWait("UpdateGuildInfo", func() {
		var g *guildsnapshotdata.GuildSnapshot
		if guildId > 0 {
			g = m.dep.GuildSnapshot().GetSnapshot(guildId)
		}

		m.hebiManager.updateGuildInfo(heroId, g)
	}, func() {
		logrus.Debugf("合璧更新联盟信息异常 guildId:%v heroId:%v", guildId, heroId)
	})
}

func (m *HebiModule) UpdateGuildInfoBatch(heroIds []int64, guildId int64) {
	m.hebiManager.FuncNoWait("UpdateGuildInfoBatch", func() {
		var g *guildsnapshotdata.GuildSnapshot
		if guildId > 0 {
			g = m.dep.GuildSnapshot().GetSnapshot(guildId)
		}

		for _, id := range heroIds {
			m.hebiManager.updateGuildInfo(id, g)
		}
	}, func() {
		logrus.Debugf("合璧更新联盟信息异常 guildId:%v heroIds:%v", guildId, heroIds)
	})
}

//gogen:iface c2s_hero_record_list
func (m *HebiModule) ProcessHebiHeroRecordList(hc iface.HeroController) {
	hc.Send(m.hebiManager.GetHeroRecordMsg(hc.Id()))
}

//gogen:iface c2s_view_show_prize
func (m *HebiModule) ProcessViewHebiShowPrize(proto *hebi.C2SViewShowPrizeProto, hc iface.HeroController) {
	heroLevel := u64.FromInt32(proto.HeroLevel)
	if heroLevel <= 0 {
		hc.Send(hebi.ERR_VIEW_SHOW_PRIZE_FAIL_INVALID_HERO_LEVEL)
		return
	}

	goodsId := u64.FromInt32(proto.Goods)
	g := m.dep.Datas().GetGoodsData(goodsId)
	if g == nil {
		hc.Send(hebi.ERR_VIEW_SHOW_PRIZE_FAIL_INVALID_GOODS)
		return
	}

	prize := m.dep.Datas().GetHebiPrizeData(hebiConf.GenHebiPrizeId(heroLevel, g.HebiType, g.GoodsQuality.Level))
	if prize == nil {
		hc.Send(hebi.ERR_VIEW_SHOW_PRIZE_FAIL_NO_PRIZE)
		return
	}

	hc.Send(prize.ShowPrizeMsg)
}
