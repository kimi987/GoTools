package mingc_war

import (
	"github.com/lightpaw/male7/config/data"
	"github.com/lightpaw/male7/config/mingcdata"
	"github.com/lightpaw/male7/config/singleton"
	"github.com/lightpaw/male7/entity/cb"
	"github.com/lightpaw/male7/gen/iface"
	"github.com/lightpaw/male7/gen/pb/chat"
	"github.com/lightpaw/male7/gen/pb/mingc_war"
	"github.com/lightpaw/male7/gen/pb/misc"
	"github.com/lightpaw/male7/pb/server_proto"
	"github.com/lightpaw/male7/pb/shared_proto"
	"github.com/lightpaw/male7/util/call"
	"github.com/lightpaw/male7/util/i32"
	"github.com/lightpaw/male7/util/idbytes"
	"github.com/lightpaw/male7/util/timeutil"
	"github.com/lightpaw/male7/util/u64"
	"time"
)

type McWarBuilding struct {
	atk  bool
	pos  cb.Cube
	data *mingcdata.MingcWarBuildingData
	name string

	prosperity      uint64
	lastDestroyTime time.Time

	troops map[int64]*McWarTroop

	// 投石机
	touShiTargets        []cb.Cube
	touShiTurnEndTime    time.Time
	touShiPrepareEndTime time.Time
	touShiTargetIndex    int
	lastFireHeroName     string
	lastFireHeroCountry  int32

	// 炮弹爆炸队列
	explodeBombs []*TouShiBomb
}

func newMcWarBuilding(atk bool, pos cb.Cube, name string, data *mingcdata.MingcWarBuildingData) *McWarBuilding {
	b := &McWarBuilding{}
	b.atk = atk
	b.name = name
	b.pos = pos
	b.data = data
	b.prosperity = data.Prosperity
	b.troops = make(map[int64]*McWarTroop)

	return b
}

func (b *McWarBuilding) unmarshal(p *server_proto.McWarBuildingServerProto, mcId uint64, datas iface.ConfigDatas) {
	b.atk = p.Atk
	b.name = p.Name
	b.pos = cb.Cube(p.Pos)
	b.prosperity = p.Prosperity
	b.lastDestroyTime = timeutil.Unix64(p.LastDestroyTime)
	b.data = datas.GetMingcWarBuildingData(uint64(p.Type))

	b.troops = make(map[int64]*McWarTroop)
	for _, tp := range p.Troop {
		t := &McWarTroop{datas: datas}
		if t.unmarshal(tp) {
			heroId, _ := idbytes.ToId(tp.Hero.Id)
			b.troops[heroId] = t
		}
	}

	for _, t := range datas.GetMingcWarSceneData(mcId).GetTouShiTarget(b.pos) {
		b.touShiTargets = append(b.touShiTargets, cb.Cube(t))
	}
	b.touShiPrepareEndTime = timeutil.Unix64(p.TouShiPrepareEndTime)
	b.touShiTurnEndTime = timeutil.Unix64(p.TouShiTurnEndTime)
	b.touShiTargetIndex = int(p.TouShiTargetIndex)
	b.lastFireHeroCountry = p.LastFireHeroCountry
	b.lastFireHeroName = p.LastFireHeroName
	// 防止策划改配置
	if int(b.touShiTargetIndex) >= len(b.touShiTargets) {
		b.touShiTargetIndex = int(i32.Max(0, int32(len(b.touShiTargets)-1)))
	}
}

func (b *McWarBuilding) encodeServer() *server_proto.McWarBuildingServerProto {
	p := &server_proto.McWarBuildingServerProto{}
	p.Atk = b.atk
	p.Name = b.name
	p.Pos = uint64(b.pos)
	p.Type = b.data.Type
	p.Prosperity = b.prosperity
	p.LastDestroyTime = timeutil.Marshal64(b.lastDestroyTime)
	for _, t := range b.troops {
		p.Troop = append(p.Troop, t.encodeServer())
	}

	p.TouShiPrepareEndTime = timeutil.Marshal64(b.touShiPrepareEndTime)
	p.TouShiTurnEndTime = timeutil.Marshal64(b.touShiTurnEndTime)
	p.TouShiTargetIndex = int32(b.touShiTargetIndex)
	p.LastFireHeroName = b.lastFireHeroName
	p.LastFireHeroCountry = b.lastFireHeroCountry

	return p
}

func (b *McWarBuilding) encode() *shared_proto.McWarSceneBuildingProto {
	p := &shared_proto.McWarSceneBuildingProto{}
	p.Atk = b.atk
	p.PosX, p.PosY = b.pos.XYI32()
	p.Type = b.data.Type
	p.Prosperity = u64.Int32(b.prosperity)
	p.LastDestroyProsperityTime = timeutil.Marshal32(b.lastDestroyTime)

	for _, t := range b.touShiTargets {
		x, y := t.XYI32()
		p.TouShiTargetPosX = append(p.TouShiTargetPosX, x)
		p.TouShiTargetPosY = append(p.TouShiTargetPosY, y)
	}
	p.TouShiPrepareEndTime = timeutil.Marshal32(b.touShiPrepareEndTime)
	p.TouShiTurnEndTime = timeutil.Marshal32(b.touShiTurnEndTime)
	p.TouShiTargetIndex = int32(b.touShiTargetIndex)
	p.LastFireHeroName = b.lastFireHeroName
	p.LastFireHeroCountry = b.lastFireHeroCountry

	return p
}

func (b *McWarBuilding) onUpdate(scene *McWarScene, ctime time.Time) (changed bool) {
	if !b.data.CanBeAtked || b.prosperity <= 0 {
		return
	}

	miscData := scene.dep.Datas().MingcMiscData()
	if ctime.Before(b.lastDestroyTime.Add(miscData.DestroyProsperityDuration)) {
		return
	}

	// 部队打掉血
	b.lastDestroyTime = ctime
	changed = b.troopDestroyBuilding(scene)

	// 炮弹打掉血
	changed = b.touShiBombExplode(scene, ctime) || changed

	// 名城战系统聊天
	b.sendBuildingDestroyChat(scene)

	return
}

func (b *McWarBuilding) sendBuildingDestroyChat(scene *McWarScene) {
	if b.prosperity > 0 {
		return
	}

	var ownerText, enemyText string
	if b.atk {
		ownerText = scene.dep.Datas().TextHelp().McWarChatOurAtkBuildingDestroy.New().WithBuilding(b.name).JsonString()
		enemyText = scene.dep.Datas().TextHelp().McWarChatEnemyAtkBuildingDestroy.New().WithBuilding(b.name).JsonString()
	} else {
		ownerText = scene.dep.Datas().TextHelp().McWarChatOurDefBuildingDestroy.New().WithBuilding(b.name).JsonString()
		enemyText = scene.dep.Datas().TextHelp().McWarChatEnemyDefBuildingDestroy.New().WithBuilding(b.name).JsonString()
	}

	go call.CatchPanic(func() {
		scene.dep.Chat().SysChatSendFunc(0, 0, shared_proto.ChatType_ChatMcWar, ownerText, shared_proto.ChatMsgType_ChatMsgMcWarSys, true, true, true, false, func(proto *shared_proto.ChatMsgProto) {
			if b.atk {
				scene.recordAtkChat(proto)
			} else {
				scene.recordDefChat(proto)
			}
			scene.broadcastCampExclude(chat.NewS2cOtherSendChatMarshalMsg(proto), b.atk, true, 0)
		})

		scene.dep.Chat().SysChatSendFunc(0, 0, shared_proto.ChatType_ChatMcWar, enemyText, shared_proto.ChatMsgType_ChatMsgMcWarSys, true, true, true, false, func(proto *shared_proto.ChatMsgProto) {
			if !b.atk {
				scene.recordAtkChat(proto)
			} else {
				scene.recordDefChat(proto)
			}
			scene.broadcastCampExclude(chat.NewS2cOtherSendChatMarshalMsg(proto), !b.atk, true, 0)
		})
	}, "名城战据点摧毁系统聊天")

}

func (b *McWarBuilding) troopDestroyBuilding(scene *McWarScene) (destroyed bool) {
	miscData := scene.dep.Datas().MingcMiscData()
	troops := b.sortAtkTroopsByMode(b.troops)

	var atkCount uint64
	var allDestroyAmount uint64
	for _, t := range troops {
		if atkCount >= miscData.DestroyProsperityMaxTroop {
			break
		}
		atkCount++

		destroy := t.getDestroyProsperity()
		allDestroyAmount += destroy
		b.troopDestroyBuildingRecord(t, scene, destroy)
	}

	if atkCount <= 0 {
		return
	}

	destroyed = true
	b.reduceBuildingPropsperity(allDestroyAmount, scene)
	return
}

func (b *McWarBuilding) reduceBuildingPropsperity(amount uint64, scene *McWarScene) {
	b.prosperity = u64.Sub(b.prosperity, amount)

	x, y := b.pos.XYI32()
	scene.broadcast(mingc_war.NewS2cSceneBuildingDestroyProsperityMsg(x, y, u64.Int32(b.prosperity)))
}

func (b *McWarBuilding) troopDestroyBuildingRecord(t *McWarTroop, scene *McWarScene, destroy uint64) {
	// 每个联盟摧毁的繁荣度记录
	if r, ok := scene.record.guilds[t.gid]; ok {
		r.destroyed += destroy
	}
	t.destroyBuilding += destroy

	if t.rankObj != nil {
		t.rankObj.addDestroy(destroy)
		scene.troopsRank.needSort = true
	}
	return
}

// 建筑中的敌方队伍，攻城车排前面
func (b *McWarBuilding) sortAtkTroopsByMode(ts map[int64]*McWarTroop) (troops []*McWarTroop) {
	var normalTroops, freeTankTroops []*McWarTroop

	for _, t := range ts {
		if t.atk == b.atk || t.action.getState() != shared_proto.MingcWarTroopState_MC_TP_STATION {
			continue
		}

		switch t.mode {
		case shared_proto.MingcWarModeType_MC_MT_NORMAL:
			normalTroops = append(normalTroops, t)
		case shared_proto.MingcWarModeType_MC_MT_FREE_TANK:
			freeTankTroops = append(freeTankTroops, t)
		}
	}

	troops = append(troops, freeTankTroops...)
	troops = append(troops, normalTroops...)

	return
}

// atk: 移动到据点的一方，可能是攻方或守方
func (b *McWarBuilding) arriveOrFight(moveTroop *McWarTroop, scene *McWarScene, ctime time.Time) (arrived bool) {
	// 己方复活点
	if b.atk == moveTroop.atk && b.data.Type == shared_proto.MingcWarBuildingType_MC_B_RELIVE {
		return true
	}

	// 是否入侵
	var invade bool
	if moveTroop.atk != b.atk {
		invade = true
	}

	var fight bool

	defer func() {
		if fight {
			x, y := b.pos.XYI32()
			scene.broadcast(mingc_war.NewS2cSceneBuildingFightMsg(idbytes.ToBytes(moveTroop.heroId), int32(moveTroop.gid), arrived, x, y))
		}
	}()

	for _, stationTroop := range b.troops {
		if stationTroop.atk == moveTroop.atk {
			continue
		}
		if stationTroop.action.getState() != shared_proto.MingcWarTroopState_MC_TP_STATION || stationTroop.action.getPos() != b.pos {
			continue
		}

		fight = true

		moveTroopWin := doFight(b, moveTroop, stationTroop, invade, scene, ctime)
		// 进攻失败，据点战斗结束
		if !moveTroopWin {
			if d := scene.dep.Datas().TextHelp().McWarAtkFail; d != nil {
				scene.dep.World().Send(moveTroop.heroId, misc.NewS2cScreenShowWordsMsg(d.Text.New().WithBuilding(b.name).WithHero(stationTroop.hero.Name).JsonString()))
			}
			return
		}

		// 一场战斗成功，对应守方回复活点
		stationTroop.failToRelive(scene, b, ctime, 0)

		if d := scene.dep.Datas().TextHelp().McWarDefFail; d != nil {
			scene.dep.World().Send(stationTroop.heroId, misc.NewS2cScreenShowWordsMsg(d.Text.New().WithBuilding(b.name).WithHero(moveTroop.hero.Name).JsonString()))
		}
	}

	arrived = true

	return
}

func (b *McWarBuilding) touShiTurnTo(left bool, ctime time.Time, turnDuration time.Duration) (succ bool) {
	if b.data.Type != shared_proto.MingcWarBuildingType_MC_B_TOU_SHI {
		return
	}

	var step int
	if left {
		step = -1
	} else {
		step = 1
	}

	succ = true
	newTurn := b.touShiTargetIndex + step
	if int(newTurn) >= len(b.touShiTargets) {
		b.touShiTargetIndex = 0
	} else if int(newTurn) < 0 {
		b.touShiTargetIndex = int(i32.Max(0, int32(len(b.touShiTargets)-1)))
	} else {
		b.touShiTargetIndex = newTurn
	}
	b.touShiTurnEndTime = ctime.Add(turnDuration)

	return
}

func (b *McWarBuilding) updateLastFireHero(t *McWarTroop, miscConfig *singleton.MiscConfig) {
	b.lastFireHeroName = miscConfig.FlagHeroName.FormatIgnoreEmpty(t.hero.GuildFlagName, t.hero.Name)
	b.lastFireHeroCountry = t.hero.CountryId
}

func (b *McWarBuilding) touShiBombExplode(scene *McWarScene, ctime time.Time) (destroyed bool) {
	var bombCount int
	for _, bomb := range b.explodeBombs {
		if bomb == nil {
			bombCount++
			continue
		}

		if ctime.Before(bomb.explodeTime.Add(-mingcdata.McWarLoopDuration)) {
			break
		}

		bomb.explode(b, scene, ctime)
		bombCount++
	}

	if bombCount <= 0 {
		return
	}

	destroyed = true

	if bombCount >= len(b.explodeBombs) {
		b.explodeBombs = []*TouShiBomb{}
	} else {
		b.explodeBombs = b.explodeBombs[bombCount:]
	}

	return
}

type TouShiBomb struct {
	fireTime    time.Time
	explodeTime time.Time

	fireTroop *McWarTroop

	baseHurt     uint64
	hurtPercent  *data.Amount
	hurtMaxTroop uint64

	destroyPerpsperity uint64
}

func newTouShiBomb(fireTroop *McWarTroop, ctime time.Time, miscData *mingcdata.MingcMiscData) *TouShiBomb {
	b := &TouShiBomb{}
	b.fireTroop = fireTroop
	b.fireTime = ctime
	b.explodeTime = ctime.Add(miscData.TouShiBuildingBombFlyDuration)
	b.baseHurt = miscData.TouShiBuildingBaseHurt
	b.hurtPercent = miscData.TouShiBuildingHurtPercent
	b.hurtMaxTroop = miscData.TouShiBuildingBaseHurtMaxTroop
	b.destroyPerpsperity = miscData.TouShiBuildingDestroyProsperity
	return b
}

func (bomb *TouShiBomb) explode(b *McWarBuilding, scene *McWarScene, ctime time.Time) {
	// 部队伤害
	var hurtCount uint64
	for _, t := range b.troops {
		if t.action.getState() != shared_proto.MingcWarTroopState_MC_TP_STATION {
			continue
		}
		if hurtCount >= bomb.hurtMaxTroop {
			break
		}
		hurtCount++

		hurt := t.hurtByPercentToAlive(bomb.hurtPercent)

		troopAlive, thisHurt := t.hurtByAmount(bomb.baseHurt)
		hurt += thisHurt
		scene.broadcast(mingc_war.NewS2cSceneTroopUpdateMsg(idbytes.ToBytes(t.heroId), t.encode()).Static())

		if !troopAlive {
			t.failToRelive(scene, b, ctime, 0)

			resp := &server_proto.CombatXResponseServerProto{AttackerWin: true}
			scene.addFightRecord(bomb.fireTroop, t, resp, hurt, 0, 0, 0, ctime, b.pos, true)
			scene.dep.World().Send(t.heroId, mingc_war.SCENE_TROOP_RECORD_ADD_NOTICE_S2C)
			scene.dep.World().Send(bomb.fireTroop.heroId, mingc_war.SCENE_TROOP_RECORD_ADD_NOTICE_S2C)
		}
	}

	// 建筑伤害
	b.reduceBuildingPropsperity(bomb.destroyPerpsperity, scene)
	b.troopDestroyBuildingRecord(bomb.fireTroop, scene, bomb.destroyPerpsperity)
	scene.broadcast(mingc_war.NewS2cSceneTouShiBombExplodeNoticeMsg(idbytes.ToBytes(bomb.fireTroop.heroId)))
}
