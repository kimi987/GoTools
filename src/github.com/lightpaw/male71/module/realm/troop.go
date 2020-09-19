package realm

import (
	"sync"

	"github.com/lightpaw/male7/config/data"
	"github.com/lightpaw/male7/config/monsterdata"
	"github.com/lightpaw/male7/config/resdata"
	"github.com/lightpaw/male7/entity"
	"github.com/lightpaw/male7/entity/npcid"
	"github.com/lightpaw/male7/gen/iface"
	"github.com/lightpaw/male7/module/realm/realmevent"
	"github.com/lightpaw/male7/module/realm/realmface"
	"github.com/lightpaw/male7/pb/shared_proto"
	"github.com/lightpaw/male7/util/i64"
	"github.com/lightpaw/male7/util/idbytes"
	"github.com/lightpaw/male7/util/u64"
	"github.com/lightpaw/pbutil"
	"time"
	"github.com/lightpaw/male7/pb/server_proto"
	"github.com/lightpaw/male7/config/regdata"
	"github.com/lightpaw/male7/config/singleton"
	"github.com/lightpaw/male7/entity/hexagon"
	"github.com/lightpaw/male7/service/conflict/heroservice/herolock"
	"github.com/lightpaw/male7/util/must"
	"github.com/lightpaw/male7/util/imath"
	"github.com/pkg/errors"
	"github.com/lightpaw/male7/util/timeutil"
	"github.com/lightpaw/male7/constants"
	"github.com/lightpaw/male7/gen/pb/region"
	"github.com/lightpaw/male7/util/i32"
)

var (
	_ realmface.Captain = (*captain)(nil)
	_ realmface.Troop   = (*troop)(nil)
)

type troop struct {
	idbytes.IdHolder // 队伍的id. 并不保存到db中. 开服时重新生成, 运行时作为唯一标识

	originTargetId               int64
	originTargetX, originTargetY int // 初始设置目标位置

	startingBase      *baseWithData // db里只存id
	startingBaseLevel uint64        // 对于多等级怪物有效
	targetBase        *baseWithData // db里只存id
	targetBaseLevel   uint64        // 对于多等级怪物有效

	// 某些怪物属于独立的君主（匈奴入侵怪物）
	mmd   *monsterdata.MonsterMasterData
	owner *shared_proto.HeroBasicSnapshotProto // 拥有者，可能为空

	ownerCanSeeTarget bool // 英雄自己可以看到

	createTime time.Time // 部队创建时间

	moveSpeedRate float64 // 初始移动速度加成

	moveStartTime            time.Time
	moveArriveTime           time.Time
	robbingEndTime           time.Time // 抢劫结束时间
	nextAddPrizeTime         time.Time // 下次给奖励时间
	nextReduceProsperityTime time.Time // 下次扣繁荣度时间
	nextAddHateTime          time.Time // 下次加仇恨的时间
	nextRobBaowuTime         time.Time // 下次加宝物的时间

	state realmface.TroopState

	//gold, food, wood, stone uint64
	//
	//jadeOre uint64

	accumRobPrize      *resdata.PrizeBuilder
	accumRobPrizeProto *shared_proto.PrizeProto

	accumReduceProsperity uint64 // TODO 保存这个值

	accumRobBaowuCount uint64

	captains []realmface.Captain // 这里面没有nil

	// 击退列表
	killEnemy []string

	event realmevent.Event

	// 城墙属性，城墙固定伤兵数，城墙属性为空表示没有城墙
	wallLevel     uint64
	wallStat      *data.SpriteStat // 可能为空
	wallFixDamage uint64

	assembly *assembly

	npcTimes uint64 // 本次讨伐野怪次数

	// 头顶对话
	dialogue uint64

	// --- msg ---
	//proto        *shared_proto.MilitaryInfoProto
	//protoBytes   []byte
	protoBytesCache []byte
	protoChanged    bool

	careMeHeroIds map[int64]struct{}
	//kimi 道路轨迹点
	realmTracks []*block_index
	//观察者列表
	watchList *sync.Map

	assistDefendStartTime time.Time // 援助驻扎开始时间

	updateMsgMa pbutil.Buffer
	removeMsgMa pbutil.Buffer

	addSeeMeMsg    pbutil.Buffer
	removeSeeMeMsg pbutil.Buffer
}

func (r *Realm) newTempTroop(id int64, startingBase *baseWithData, startingBaseLevel uint64, encode CaptainEncoder, captains []*entity.TroopPos) *troop {

	troopCaptains := make([]realmface.Captain, 0, len(captains))
	for i, pos := range captains {
		c := pos.Captain()
		if c == nil {
			continue
		}

		captain := &captain{
			idAtHero: c.Id(),
			index:    i + 1,
			proto:    encode(pos, false),
		}

		troopCaptains = append(troopCaptains, captain)
	}

	return r.newBasicTroop(id, startingBase, startingBaseLevel, startingBase, startingBaseLevel, troopCaptains, realmface.Temp, 0, 0, 0)
}

func (r *Realm) newTempTroopWithInfo(id int64, startingBase *baseWithData, startingBaseLevel uint64, captains []*shared_proto.CaptainInfoProto) *troop {

	troopCaptains := make([]realmface.Captain, 0, len(captains))
	for i, c := range captains {
		if c == nil {
			continue
		}

		captain := &captain{
			idAtHero: u64.FromInt32(c.Id),
			index:    i + 1,
			proto:    c,
		}

		troopCaptains = append(troopCaptains, captain)
	}

	return r.newBasicTroop(id, startingBase, startingBaseLevel, startingBase, startingBaseLevel, troopCaptains, realmface.Temp, 0, 0, 0)
}

func (r *Realm) newNpcTroop(id int64, startingBase *baseWithData, startingBaseLevel uint64,
	targetBase *baseWithData, targetBaseLevel uint64, mmd *monsterdata.MonsterMasterData, state realmface.TroopState) *troop {

	captains := newNpcTroopCaptains(mmd)
	t := r.newBasicTroop(id, startingBase, startingBaseLevel, targetBase, targetBaseLevel, captains, state, 0, 0, 0)
	t.mmd = mmd
	t.owner = newNpcOwner(mmd, startingBase, startingBaseLevel)

	return t
}

func newNpcOwner(mmd *monsterdata.MonsterMasterData, startingBase *baseWithData, startingBaseLevel uint64) *shared_proto.HeroBasicSnapshotProto {
	owner := mmd.EncodeSnapshot(startingBase.Id())
	owner.BaseX = int32(startingBase.BaseX())
	owner.BaseY = int32(startingBase.BaseY())
	owner.BaseRegion = i64.Int32(startingBase.RegionID())
	owner.Prosperity = u64.Int32(startingBase.Prosperity())

	if info := startingBase.internalBase.getBaseInfoByLevel(startingBaseLevel); info != nil {
		owner.BaseLevel = u64.Int32(info.GetBaseLevel())
	} else {
		owner.BaseLevel = u64.Int32(startingBase.BaseLevel())
	}

	return owner
}

func (r *Realm) newMoveNpcTroop(id int64, startingBase *baseWithData, startingBaseLevel uint64,
	targetBase *baseWithData, targetBaseLevel uint64, mmd *monsterdata.MonsterMasterData, isInvade bool,
	arriveOffset time.Duration) *troop {

	var state realmface.TroopState
	if isInvade {
		state = realmface.MovingToInvade
	} else {
		state = realmface.MovingToAssist
	}

	t := r.newNpcTroop(id, startingBase, startingBaseLevel, targetBase, targetBaseLevel, mmd, state)
	r.initMoveTroop(t, arriveOffset)
	//kimi 添加行军路径
	t.AddRealmPoint(r)

	return t
}

func (r *Realm) newDefendingNpcTroop(id int64, startingBase *baseWithData, startingBaseLevel uint64,
	targetBase *baseWithData, targetBaseLevel uint64, mmd *monsterdata.MonsterMasterData) *troop {

	t := r.newNpcTroop(id, startingBase, startingBaseLevel, targetBase, targetBaseLevel, mmd, realmface.Defending)

	now := r.services.timeService.CurrentTime()
	createTime := now.Add(-time.Second)

	t.createTime = createTime
	t.moveStartTime = createTime
	t.moveArriveTime = now

	return t
}

func IsOnlyOwnerCanSeeBase(targetBase *baseWithData) bool {
	return npcid.IsHomeNpcId(targetBase.Id())
}

func (r *Realm) newBasicTroop(id int64, startingBase *baseWithData, startingBaseLevel uint64,
	targetBase *baseWithData, targetBaseLevel uint64, captains []realmface.Captain,
	state realmface.TroopState, assemblyCount, dialogue uint64, moveSpeedRate float64) *troop {

	result := &troop{
		IdHolder:       idbytes.NewIdHolder(id),
		originTargetId: targetBase.Id(),
		originTargetX:  targetBase.BaseX(),
		originTargetY:  targetBase.BaseY(),

		startingBase:      startingBase,
		startingBaseLevel: startingBaseLevel,
		targetBase:        targetBase,
		targetBaseLevel:   targetBaseLevel,

		ownerCanSeeTarget: IsOnlyOwnerCanSeeBase(targetBase),
		moveSpeedRate:     moveSpeedRate,

		state:         state,
		dialogue:      dialogue,
		accumRobPrize: resdata.NewPrizeBuilder(),
		captains:      captains,
		careMeHeroIds: make(map[int64]struct{}),
	}

	if assemblyCount > 0 {
		result.assembly = newAssembly(result, assemblyCount)
	}

	return result
}

type CaptainEncoder func(t *entity.TroopPos, fullSoldier bool) *shared_proto.CaptainInfoProto

func (r *Realm) newTroop(id int64, startingBase *baseWithData, startingBaseLevel uint64,
	targetBase *baseWithData, targetBaseLevel uint64, state realmface.TroopState,
	encode CaptainEncoder, captains []*entity.TroopPos, assemblyCount, dialogue uint64, moveSpeedRate float64) *troop {
	var array []realmface.Captain
	for i, pos := range captains {
		c := pos.Captain()
		if c == nil {
			continue
		}

		captain := &captain{
			idAtHero: c.Id(),
			index:    i + 1,
			proto:    encode(pos, false),
		}

		array = append(array, captain)
	}

	return r.newTroopWithCaptains(id, startingBase, startingBaseLevel, targetBase, targetBaseLevel, state, array, assemblyCount, dialogue, moveSpeedRate)
}

func (r *Realm) newTroopWithCaptains(id int64, startingBase *baseWithData, startingBaseLevel uint64,
	targetBase *baseWithData, targetBaseLevel uint64, state realmface.TroopState, captains []realmface.Captain, assemblyCount, dialogue uint64, moveSpeedRate float64) *troop {

	result := r.newBasicTroop(id, startingBase, startingBaseLevel, targetBase, targetBaseLevel, captains, state, assemblyCount, dialogue, moveSpeedRate)
	r.initMoveTroop(result, 0)
	//kimi 添加行军路径
	result.AddRealmPoint(r)

	return result
}

func (r *Realm) initTroopMoveTime(result *troop, arriveOffset time.Duration) {
	now := r.services.timeService.CurrentTime()

	if arriveOffset > 0 {
		arriveTime := now.Add(arriveOffset)
		createTime := arriveTime.Add(-result.getMoveDuration(r))

		result.createTime = createTime
		result.moveStartTime = createTime
		result.moveArriveTime = arriveTime
	} else {
		result.createTime = now
		result.moveStartTime = now
		result.moveArriveTime = now.Add(result.getMoveDuration(r))
	}
}

func (r *Realm) initMoveTroop(result *troop, arriveOffset time.Duration) {

	switch result.State() {
	case realmface.MovingToInvade:
		r.initTroopMoveTime(result, arriveOffset)
		result.event = r.newEvent(result.moveArriveTime, result.invasionArrivedEvent)
	case realmface.MovingToAssist:
		r.initTroopMoveTime(result, arriveOffset)
		result.event = r.newEvent(result.moveArriveTime, result.assistArrivedEvent)
	case realmface.MovingToAssembly:
		r.initTroopMoveTime(result, arriveOffset)
		result.event = r.newEvent(result.moveArriveTime, result.assemblyArrivedEvent)
	case realmface.MovingToWorkshopBuild:
		r.initTroopMoveTime(result, arriveOffset)
		result.event = r.newEvent(result.moveArriveTime, result.workshopArrivedEvent)
	case realmface.MovingToWorkshopProd:
		r.initTroopMoveTime(result, arriveOffset)
		result.event = r.newEvent(result.moveArriveTime, result.workshopArrivedEvent)
	case realmface.MovingToWorkshopPrize:
		r.initTroopMoveTime(result, arriveOffset)
		result.event = r.newEvent(result.moveArriveTime, result.workshopPrizeArrivedEvent)
	case realmface.MovingToInvesigate:
		r.initTroopMoveTime(result, arriveOffset)
		result.event = r.newEvent(result.moveArriveTime, result.investigateArrivedEvent)
	default:
		return
	}
}

func (r *Realm) rebuildTroop(proto *server_proto.TroopServerProto, ctime time.Time) error {
	if proto.GetTroopState() == realmface.Temp {
		return errors.New("rebuildTroop，无效的部队状态")
	}

	moveArrvieTime := timeutil.Unix64(proto.MoveArriveTime)
	// 容错处理，如果到达时间距离当前时间超过限制，删掉这个马
	if ctime.Before(moveArrvieTime) && moveArrvieTime.Sub(ctime) > constants.RebuildMaxTroopMoveDuration {
		return errors.New("rebuildTroop发现到达时间超过限制")
	}

	robbingEndTime := timeutil.Unix64(proto.RobbingEndTime)
	if ctime.Before(robbingEndTime) && robbingEndTime.Sub(ctime) > constants.RebuildMaxRecurringRobDuration {
		return errors.New("rebuildTroop发现持续掠夺时间超过限制")
	}

	startingBase := r.getBase(proto.StartingBaseId)
	if startingBase == nil {
		return errors.New("rebuildTroop没有找到hero的base")
	}

	var targetBase *baseWithData
	if targetBaseID := proto.TargetBaseId; targetBaseID != 0 {
		targetBase = r.getBase(targetBaseID)
		if targetBase == nil {
			return errors.New("rebuildTroop没有找到troop的targetBase. 可能没有清掉?")
		}
	}

	if realmface.TroopState(proto.TroopState) == realmface.Robbing && targetBase == nil {
		return errors.New("rebuildTroop没有找到troop的targetBase")
	}

	var mmd *monsterdata.MonsterMasterData
	var owner *shared_proto.HeroBasicSnapshotProto
	if proto.MonsterMasterId != 0 {
		mmd = r.services.datas.GetMonsterMasterData(proto.MonsterMasterId)
		if mmd == nil {
			return errors.New("rebuildTroop没有找到troop的MonsterMasterData")
		}

		owner = newNpcOwner(mmd, startingBase, proto.StartingBaseLevel)
	}

	targetBaseLevel := proto.TargetBaseLevel

	result := &troop{
		IdHolder:       idbytes.NewIdHolder(proto.Id),
		originTargetId: proto.OriginTargetId,

		startingBase:      startingBase,
		startingBaseLevel: proto.StartingBaseLevel,
		targetBase:        targetBase,
		targetBaseLevel:   targetBaseLevel,

		mmd:   mmd,
		owner: owner,

		originTargetX: int(proto.BackHomeTargetX),
		originTargetY: int(proto.BackHomeTargetY),

		ownerCanSeeTarget: proto.TargetIsOwnerCanSee,

		createTime:               timeutil.Unix64(proto.CreateTime),
		moveStartTime:            timeutil.Unix64(proto.MoveStartTime),
		moveArriveTime:           moveArrvieTime,
		robbingEndTime:           robbingEndTime,
		nextReduceProsperityTime: timeutil.Unix64(proto.NextReduceProsperityTime),
		nextAddHateTime:          timeutil.Unix64(proto.NextAddHateTime),

		state:         realmface.TroopState(proto.TroopState),
		accumRobPrize: resdata.NewPrizeBuilder(),
		captains:      make([]realmface.Captain, 0, len(proto.Captains)),
		careMeHeroIds: make(map[int64]struct{}),

		assistDefendStartTime: timeutil.Unix64(proto.AssistDefendStartTime),

		npcTimes: proto.MultiLevelMonsterCount,
	}

	if prizeProto := proto.AccumRobPrize; prizeProto != nil {
		result.accumRobPrize.Add(resdata.UnmarshalPrize(prizeProto, r.services.datas))
	}
	result.accumReduceProsperity = proto.AccumReduceProsperity

	for index, c := range proto.Captains {
		result.captains = append(result.captains, &captain{
			idAtHero: u64.FromInt32(c.Id),
			index:    int(proto.CaptainIndex[index]),
			proto:    c,
		})
	}

	// 初始化
	result.onChanged()

	startingBase.selfTroops[result.Id()] = result
	if targetBase != nil {
		targetBase.targetingTroops[result.Id()] = result

		targetBase.remindAttackOrRobCountChanged(r)
	}

	r.rebuildTroopEvent(result, targetBase, targetBaseLevel, ctime)

	return nil
}

func (r *Realm) rebuildHeroTroop(t *entity.TroopInvaseInfo, hero *entity.Hero, ctime time.Time) error {

	if t.State() == realmface.Temp {
		return errors.New("rebuildHeroTroop，无效的部队状态")
	}

	// 容错处理，如果到达时间距离当前时间超过限制，删掉这个马
	if ctime.Before(t.MoveArriveTime()) && t.MoveArriveTime().Sub(ctime) > constants.RebuildMaxTroopMoveDuration {
		return errors.New("rebuildTroop发现到达时间超过限制")
	}

	if ctime.Before(t.RobbingEndTime()) && t.RobbingEndTime().Sub(ctime) > constants.RebuildMaxRecurringRobDuration {
		return errors.New("rebuildTroop发现持续掠夺时间超过限制")
	}

	startingBase := r.getBase(hero.Id())
	if startingBase == nil {
		return errors.New("rebuildTroop没有找到hero的base")
	}

	var targetBase *baseWithData
	if targetBaseID := t.TargetBaseID(); targetBaseID != 0 {
		targetBase = r.getBase(targetBaseID)
		if targetBase == nil {
			return errors.New("rebuildTroop没有找到troop的targetBase. 可能没有清掉?")
		}
	}

	if t.State() == realmface.Robbing && targetBase == nil {
		return errors.New("rebuildTroop没有找到troop的targetBase")
	}

	var assemblyMemberCount uint64
	var joinAssembly *assembly
	if t.AssemblyId() != 0 {
		if t.AssemblyId() == t.Id() {
			assemblyData := r.services.datas.GetAssemblyData(regdata.GetAssemblyTypeByTarget(t.OriginTargetID()))
			if assemblyData == nil {
				return errors.New("rebuildTroop，没有找到集结目标的集结配置，assemblyData == nil")
			}

			assemblyMemberCount = assemblyData.MemberCount
		} else {
			// 加入别人的集结
			assemblyBaseId := entity.GetTroopHeroId(t.AssemblyId())
			assemblyBase := r.getBase(assemblyBaseId)
			if assemblyBase == nil {
				return errors.New("rebuildTroop，没有找到加入的集结目标")
			}

			assemblyTroop := assemblyBase.selfTroops[t.AssemblyId()]
			if assemblyTroop == nil {
				return errors.New("rebuildTroop，没有找到加入的集结目标的部队")
			}

			if assemblyTroop.assembly == nil {
				return errors.New("rebuildTroop，没有找到加入的集结目标的部队的集结")
			}

			if assemblyTroop.State() == realmface.Assembly {
				switch t.State() {
				case realmface.MovingToAssembly, realmface.AssemblyArrived:
				default:
					return errors.Errorf("rebuildTroop，集结者是等待集结状态，自己是 %v", t.State())
				}
			} else {
				if t.State() != realmface.AssemblyArrived {
					return errors.Errorf("rebuildTroop，集结者不是等待集结状态，自己是 %v", t.State())
				}
			}

			joinAssembly = assemblyTroop.assembly
		}

	}

	targetBaseLevel := t.TargetBaseLevel()

	result := &troop{
		IdHolder:       idbytes.NewIdHolder(t.Id()),
		originTargetId: t.OriginTargetID(),

		startingBase:      startingBase,
		startingBaseLevel: startingBase.BaseLevel(),
		targetBase:        targetBase,
		targetBaseLevel:   targetBaseLevel,

		originTargetX: t.BackHomeTargetX(),
		originTargetY: t.BackHomeTargetY(),

		ownerCanSeeTarget: t.TargetIsOwnerCanSee(),
		moveSpeedRate:     hero.GetMoveSpeedRate(),

		createTime:               t.CreateTime(),
		moveStartTime:            t.MoveStartTime(),
		moveArriveTime:           t.MoveArriveTime(),
		robbingEndTime:           t.RobbingEndTime(),
		nextReduceProsperityTime: t.NextReduceProsperityTime(),
		nextAddHateTime:          t.NextAddHateTime(),
		nextRobBaowuTime:         t.NextRobBaowuTime(),

		state:         t.State(),
		dialogue:      t.Dialogue(),
		accumRobPrize: resdata.NewPrizeBuilder(),
		captains:      make([]realmface.Captain, 0, len(t.Captains())),
		careMeHeroIds: make(map[int64]struct{}),
	}

	if t.AssemblyId() != 0 {
		if t.AssemblyId() == t.Id() {
			// 自己创建的集结
			result.assembly = newAssembly(result, assemblyMemberCount)
		} else {
			result.assembly = joinAssembly
			joinAssembly.addTroop(result)
		}
	}

	if prizeProto := t.AccumRobPrize(); prizeProto != nil {
		result.accumRobPrize.Add(resdata.UnmarshalPrize(prizeProto, r.services.datas))
	}
	result.accumReduceProsperity = t.AccumReduceProsperity()

	for index, captainId := range t.Captains() {
		if captainId <= 0 {
			continue
		}

		if c := hero.Military().Captain(captainId); c != nil {
			var xIndex int32
			if index < len(t.CaptainXIndex()) {
				xIndex = t.CaptainXIndex()[index]
			}

			result.captains = append(result.captains, &captain{
				idAtHero: captainId,
				index:    index + 1,
				proto:    c.EncodeCaptainInfo(false, xIndex),
			})
		}
	}

	// 初始化
	result.onChanged()

	startingBase.selfTroops[result.Id()] = result
	if targetBase != nil {
		targetBase.targetingTroops[result.Id()] = result

		targetBase.remindAttackOrRobCountChanged(r)
	}

	r.rebuildTroopEvent(result, targetBase, targetBaseLevel, ctime)

	return nil
}

func (r *Realm) rebuildTroopEvent(result *troop, targetBase *baseWithData, targetBaseLevel uint64, ctime time.Time) {
	switch result.state {
	case realmface.InvadeMovingBack, realmface.AssistMovingBack, realmface.AssemblyMovingBack,
		realmface.WorkshopBuildMovingBack, realmface.WorkshopProdMovingBack, realmface.WorkshopPrizeMovingBack, realmface.InvesigateMovingBack:
		result.event = r.newEvent(result.moveArriveTime, result.returnedToBaseEvent)
	case realmface.MovingToInvesigate:
		//侦察事件
		result.event = r.newEvent(result.moveArriveTime, result.investigateArrivedEvent)
	case realmface.MovingToAssist:
		result.event = r.newEvent(result.moveArriveTime, result.assistArrivedEvent)
	case realmface.MovingToInvade:
		result.event = r.newEvent(result.moveArriveTime, result.invasionArrivedEvent)
	case realmface.MovingToAssembly:
		result.event = r.newEvent(result.moveArriveTime, result.assemblyArrivedEvent)
	case realmface.MovingToWorkshopBuild:
		result.event = r.newEvent(result.moveArriveTime, result.workshopArrivedEvent)
	case realmface.MovingToWorkshopProd:
		result.event = r.newEvent(result.moveArriveTime, result.workshopArrivedEvent)
	case realmface.MovingToWorkshopPrize:
		result.event = r.newEvent(result.moveArriveTime, result.workshopPrizeArrivedEvent)
	case realmface.Assembly:
		result.event = r.newEvent(result.moveArriveTime, result.assemblyTimesUpEvent)
	case realmface.Robbing:
		nextTickTime := result.InitRobbing(result.robbingEndTime, r, targetBase.internalBase.getBaseInfoByLevel(targetBaseLevel), ctime)
		result.event = r.newEvent(nextTickTime, result.doRobEvent)
	}
}

func (t *troop) getMonsterSpecType() monsterdata.SpecType {
	if t.mmd != nil {
		return t.mmd.GetSpec()
	}
	return monsterdata.None
}

func (t *troop) getStartingBaseInfo() baseInfo {
	return t.startingBase.internalBase.getBaseInfoByLevel(t.startingBaseLevel)
}

func (t *troop) getTargetBaseInfo() baseInfo {
	if t.targetBase != nil {
		return t.targetBase.internalBase.getBaseInfoByLevel(t.targetBaseLevel)
	}
	return nil
}

func (t *troop) getStartingBaseFlagHeroName(r *Realm) string {
	if t.owner != nil {
		//return r.toFlagHeroName(t.owner.Basic.GuildFlagName, t.owner.Basic.Name)
		return t.owner.Basic.Name
	}
	return r.toBaseFlagHeroName(t.startingBase, t.startingBaseLevel)
}

func (t *troop) getStartingBaseSnapshot(r *Realm) *shared_proto.HeroBasicSnapshotProto {
	if t.owner != nil {
		return t.owner
	}
	return t.getStartingBaseInfo().EncodeAsHeroBasicSnapshot(r.services.heroSnapshotService.Get)
}

func (t *troop) getStartingBaseBasicProto(r *Realm) *shared_proto.HeroBasicProto {
	if t.owner != nil {
		return t.owner.Basic
	}
	return t.getStartingBaseInfo().GetHeroBasicProto(r.services.heroSnapshotService.Get)
}

func (t *troop) getTargetBaseFlagHeroName(r *Realm) string {
	if t.targetBase != nil {
		return r.toBaseFlagHeroName(t.targetBase, t.targetBaseLevel)
	}
	return idbytes.PlayerName(t.originTargetId)
}

func (t *troop) setCareMeHeroId(heroId int64) {
	t.careMeHeroIds[heroId] = struct{}{}
}

func (t *troop) tryRemoveCareMeHeroId(heroId int64) bool {
	if _, exist := t.careMeHeroIds[heroId]; exist {
		delete(t.careMeHeroIds, heroId)
		return true
	}
	return false
}

func (t *troop) CanSee(area *realmface.ViewArea) bool {
	if area == nil {
		return false
	}
	if t.DontShowMa() {
		return false
	}
	return area.CanSeeLine(t.startingBase.BaseX(), t.startingBase.BaseY(), t.originTargetX, t.originTargetY)
}
//AddWatchListOnCreate 在创建的时候添加观察者
func (t *troop) AddWatchListOnCreate(r *Realm) {
	poses:= t.GetPos()

	for _,v := range poses {
		r.rangeCareHeros(v[0],v[1], func(hc iface.HeroController){
			t.AddWatcher(hc)
		})
	}
}

//GetPos 获取坐标
func (t *troop) GetPos() [][2]int {
	startX, startY, targetX, targetY := t.startingBase.BaseX(), t.startingBase.BaseY(), t.originTargetX, t.originTargetY
	distantX, distantY := constants.RealmIndexBlockSize, constants.RealmIndexBlockSize

	if startX/constants.RealmIndexBlockSize > targetX/constants.RealmIndexBlockSize {
		distantX = -distantX
	}

	if startY/constants.RealmIndexBlockSize > targetY/constants.RealmIndexBlockSize {
		distantY = -distantY
	}
	var points [][2]int
	for {
		points = append(points, [2]int{startX, startY})
		// logrus.Debugf("添加军队索引点[%d] [%d]", pointX, pointY)
		if startX/constants.RealmIndexBlockSize != targetX/constants.RealmIndexBlockSize {
			startX += distantX
		} else if startY/constants.RealmIndexBlockSize == targetY/constants.RealmIndexBlockSize {
			break
		}
		if startY/constants.RealmIndexBlockSize != targetY/constants.RealmIndexBlockSize {
			startY += distantY
		}
	}

	return points
}

//AddWatcher 添加观察者(这里是为了玩家在野外视野移动中,给观察对象加一个观察者索引)
func (t *troop) AddWatcher(hc iface.HeroController) {
	if t.watchList == nil {
		t.watchList = &sync.Map{}
	}
	t.watchList.Store(hc.Id(), hc)
}

//RemoveWatcher 移除观察者
func (t *troop) RemoveWatcher(hc iface.HeroController) {
	if t.watchList == nil {
		return
	}
	t.watchList.Delete(hc.Id())
}

//RangeWatchList 遍历观察者
func (t *troop) RangeWatchList(f func(hc iface.HeroController)) {
	if t.watchList == nil {
		return
	}
	t.watchList.Range(func(k, v interface{}) bool {
		hc, ok := v.(iface.HeroController)
		if ok {
			f(hc)
		}
		return true
	})
}

func (t *troop) DontShowMa() bool {
	// 军团怪和匈奴怪的野外防守马不显示出来
	if t.State() == realmface.Defending &&
		(npcid.IsJunTuanNpcId(t.startingBase.Id()) || npcid.IsXiongNuNpcId(t.startingBase.Id())) {
		return true
	}

	if t.State() == realmface.Assembly {
		// 集结等待阶段，不展示
		return true
	}

	return false
}

func (t *troop) isMatchCondition(r *Realm, cond *server_proto.MilitaryConditionProto) bool {
	if t.assembly != nil && t.assembly.self != t {
		// 加入集结的，不显示出来
		return false
	}

	return t.isMatchCondition0(r, cond, 3) // 最多3层
}

func (t *troop) isMatchCondition0(r *Realm, cond *server_proto.MilitaryConditionProto, layer int) bool {

	// 层数太多，或者条件太多，返回false(主要防止攻击，也可能是客户端Bug)
	if layer <= 0 || len(cond.Attributes)+len(cond.Conditions) > 10 {
		return false
	}

	if cond.IsOr {
		// 条件中有一个为true，返回true
		for _, a := range cond.Attributes {
			if t.isMatchAttribute(r, a) {
				return true
			}
		}

		for _, c := range cond.Conditions {
			if t.isMatchCondition0(r, c, layer-1) {
				return true
			}
		}

		// 没有任何条件匹配，返回false
		return false
	} else {
		// 条件中有一个false，返回false
		for _, a := range cond.Attributes {
			if !t.isMatchAttribute(r, a) {
				return false
			}
		}

		for _, c := range cond.Conditions {
			if !t.isMatchCondition0(r, c, layer-1) {
				return false
			}
		}

		return true
	}
}

func (t *troop) isMatchAttribute(r *Realm, a *server_proto.MilitaryAttributeProto) bool {

	// attribute中只要有一个符合条件，返回true，如果全部跳过，则返回false

	// 部队状态
	if a.TroopState != 0 && a.TroopState == int32(t.state) {
		return true
	}

	// 部队出发城池（部队拥有者）
	if id := a.StartBaseId; id != 0 && id == t.startingBase.Id() {
		return true
	}

	// 部队目的地城池
	if id := a.TargetBaseId; id != 0 && id == t.originTargetId {
		return true
	}

	// 部队出发城池联盟
	if a.StartBaseGuildId != 0 && a.StartBaseGuildId == t.startingBase.GuildId() {
		return true
	}

	// 部队目的地城池联盟
	if a.TargetBaseGuildId != 0 {

		tb := t.targetBase
		if tb == nil {
			tb = r.getBase(t.originTargetId)
		}
		if tb != nil {
			if a.TargetBaseGuildId == tb.GuildId() {
				return true
			}
		} else {
			// 可能已经回城了(targetBase == nil)
			if !npcid.IsNpcId(t.originTargetId) {
				if th := r.services.heroSnapshotService.Get(t.originTargetId); th != nil {
					if a.TargetBaseGuildId == th.GuildId {
						return true
					}
				}
			}
		}
	}

	if a.JoinAssemblyHeroId != 0 {
		if t.assembly != nil && t.assembly.isBaseJoined(a.JoinAssemblyHeroId) {
			return true
		}
	}

	return false
}

// stat 可能为空，stat为空表示没有城墙
func (t *troop) setWall(level uint64, stat *data.SpriteStat, fixDamage uint64) {
	t.wallLevel = u64.Max(level, 1)
	t.wallStat = stat
	t.wallFixDamage = fixDamage
}

func (t *troop) getMoveDuration(r *Realm) time.Duration {
	distance := t.getMoveDistance(r.config().Edge)

	if distance > 0 {
		if npcid.IsXiongNuNpcId(t.startingBase.Id()) {
			return r.services.datas.ResistXiongNuMisc().MoveDuration(distance)
		} else {
			var offset time.Duration
			if t.startingBase.isHeroHomeBase() && !npcid.IsNpcId(t.originTargetId) {
				offset = r.config().TroopMoveOffsetDuration
			}

			speed := r.CalcMoveSpeed(t.originTargetId, t.moveSpeedRate)
			return singleton.MoveDuration(distance, speed) + offset
		}
	}
	return 0
}

func (t *troop) getMoveDistance(edge float64) float64 {
	return hexagon.OffsetDistance(t.originTargetX, t.originTargetY, t.startingBase.BaseX(), t.startingBase.BaseY(), edge)
}

func (t *troop) speedUp(r *Realm, ctime time.Time, speedUpRate float64) bool {

	st := t.moveStartTime.Unix()
	et := t.moveArriveTime.Unix()
	ct := ctime.Unix()

	d0 := et - ct // 剩余行军时间
	if d0 <= 1 {
		// 最少加速1秒，已经到达目的地
		return false
	}

	d := et - st
	if d < d0 {
		// imposible
		return false
	}

	// 计算新的速度，原来速度*倍率
	distance := t.getMoveDistance(r.config().Edge)
	oldVolocityPerSecond := distance / float64(d)
	//newVelocityPerSecond := oldVolocityPerSecond * speedRate
	newVelocityPerSecond := oldVolocityPerSecond * (1 + speedUpRate)

	// 计算出当前所在的位置，rate
	rate := float64(d0) / float64(d)

	// 计算剩余距离，加速之后到达终点所需时间 d1
	fd1 := singleton.MoveDuration(distance*rate, newVelocityPerSecond).Seconds()
	d1 := int64(fd1)
	if d1 >= d0 {
		// 加速之后没变化
		return false
	}

	newEndTime := ct + d1
	t.moveArriveTime = timeutil.Unix64(newEndTime)                  // 新的到达时间
	t.moveStartTime = timeutil.Unix64(newEndTime - int64(fd1/rate)) // 新的开始时间
	t.onChanged()

	if t.event != nil {
		t.event.UpdateTime(t.moveArriveTime)
		r.updateNextEventTime(t.event)
	}

	return true
}

func (t *troop) backHome(r *Realm, ctime time.Time, destroyAssembly, updateTroop bool) {

	rate := timeutil.Rate(t.moveStartTime, t.moveArriveTime, ctime)

	if assembly := t.assembly; assembly != nil {

		if destroyAssembly && assembly.self == t && rate >= 1 {
			// 在目标的城解散
			assembly.destroyAndTroopBackHome(r, t.originTargetId, t.originTargetX, t.originTargetY, ctime, updateTroop)
			return
		}

		if assembly.self != t {
			// 玩家部队自己要回家，不涉及到全部人
			for i, member := range assembly.member {
				if member == t {
					assembly.removeTroopByIndex(i)
					break
				}
			}
			t.assembly = nil

			if t.State() == realmface.AssemblyArrived {
				t.updateCaptainStat(assembly.addedStat, nil)
				assembly.updateAddedStat(r, 0)
			}
		}

		assembly.broadcastChanged(r)
	}

	t.backHomeFrom(r, t.originTargetId, t.originTargetX, t.originTargetY, rate, ctime)

	r.broadcastMaToCared(t, addTroopTypeUpdate, 0)

	if updateTroop {
		r.heroBaseFuncNotError(t.startingBase.Base(), func(hero *entity.Hero) (heroChanged bool) {
			hero.UpdateTroop(t, false)
			return true
		})
	}
}

func (t *troop) backHomeFrom(r *Realm, fromBaseId int64, fromX, fromY int, rate float64, ctime time.Time) {
	t.onAddHateWhenBackHome(r, ctime)
	r.updateHeroSpecMonsterTask(t)

	switch t.state {
	case realmface.MovingToAssist, realmface.Defending:
		t.state = realmface.AssistMovingBack
	case realmface.MovingToAssembly, realmface.AssemblyArrived:
		if t.originTargetId == fromBaseId {
			t.state = realmface.AssemblyMovingBack
		} else {
			t.state = realmface.InvadeMovingBack
		}
	case realmface.MovingToWorkshopBuild:
		t.state = realmface.WorkshopBuildMovingBack
	case realmface.MovingToWorkshopProd:
		t.state = realmface.WorkshopProdMovingBack
	case realmface.MovingToWorkshopPrize:
		t.state = realmface.WorkshopPrizeMovingBack
	case realmface.MovingToInvesigate:
		t.state = realmface.InvesigateMovingBack
	default:
		t.state = realmface.InvadeMovingBack
	}

	if t.targetBase != nil {
		t.addAstDefendLog(r, ctime)

		delete(t.targetBase.targetingTroops, t.Id())
		t.targetBase.remindAttackOrRobCountChanged(r)
		t.targetBase = nil

		if npcid.IsXiongNuNpcId(t.startingBase.Id()) {
			t.startingBase.updateXiongNuTarget()
		}
	}

	t.assistDefendStartTime = time.Time{}

	t.originTargetId = fromBaseId
	t.originTargetX = fromX
	t.originTargetY = fromY

	moveDuration := t.getMoveDuration(r)

	arriveDuration := moveDuration
	if rate > 0 && rate < 1 {
		arriveDuration = time.Duration(float64(moveDuration) * rate)
	}

	t.moveArriveTime = ctime.Add(arriveDuration)
	t.moveStartTime = t.moveArriveTime.Add(-moveDuration)

	if t.event != nil {
		t.event.RemoveFromQueue()
	}
	t.event = r.newEvent(t.moveArriveTime, t.returnedToBaseEvent)

	t.onChanged()
}

// 援助驻扎加繁荣度日志
func (t *troop) addAstDefendLog(r *Realm, ctime time.Time) {
	if t.state != realmface.AssistMovingBack {
		return
	}
	if t.targetBase == nil || t.targetBase.Prosperity() >= t.targetBase.ProsperityCapcity() {
		return
	}
	if timeutil.IsZero(t.assistDefendStartTime) || ctime.Before(t.assistDefendStartTime) {
		return
	}

	toAdd := uint64(ctime.Sub(t.assistDefendStartTime).Minutes()) * r.dep.Datas().RegionConfig().AstDefendRestoreHomeProsperityAmount.Calc(t.targetBase.ProsperityCapcity())
	if toAdd <= 0 {
		return
	}

	heroId := t.startingBase.Id()
	hero := r.services.heroSnapshotService.Get(heroId)
	if hero == nil {
		return
	}

	var flagName string
	if g := r.dep.GuildSnapshot().GetSnapshot(t.startingBase.GuildId()); g != nil {
		flagName = g.FlagName
	}

	fullName := r.dep.Datas().MiscConfig().FlagHeroName.FormatIgnoreEmpty(flagName, hero.Name)
	r.AddAstDefendLog(t.targetBase.Id(), ctime, t.assistDefendStartTime, fullName, toAdd)
}

func (t *troop) onAddHateWhenBackHome(r *Realm, ctime time.Time) {
	if info := t.getTargetBaseInfo(); info != nil {
		if hateData := info.getHateData(); hateData != nil {
			if hateTypeData := hateData.TypeData(); hateTypeData != nil {
				r.heroBaseFuncWithSend(t.startingBase.Base(), func(hero *entity.Hero, result herolock.LockResult) {
					t.doAddHate(hero, result, hateTypeData, ctime)
				})
			}
		}
	}
}

// 处理从目标主城离开
// 1、到达目标被打死
// 2、持续掠夺/援助中，被打死
// 3、自己主城爆了
// 4、回家（主动或被动）
func (t *troop) leaveTarget(hero *entity.Hero, result herolock.LockResult, ctime time.Time) {
	if info := t.getTargetBaseInfo(); info != nil {
		if hateData := info.getHateData(); hateData != nil {
			if hateTypeData := hateData.TypeData(); hateTypeData != nil {
				t.doAddHate(hero, result, hateTypeData, ctime)
			}
		}
	}
}

func (t *troop) doAddHate(hero *entity.Hero, result herolock.LockResult, hateTypeData *regdata.RegionMultiLevelNpcTypeData, ctime time.Time) {

	// 有仇恨
	heroTroop := hero.GetTroopByIndex(entity.GetTroopIndex(t.Id()))
	if heroTroop == nil {
		return
	}

	ii := heroTroop.GetInvateInfo()
	if ii == nil || ii.AccumAddHate() <= 0 {
		return
	}

	info := hero.GetOrCreateNpcTypeInfo(hateTypeData.Type)
	if info == nil {
		return
	}

	if info.GetHate() < hateTypeData.MaxHate {
		toAdd := ii.ClearAccumAddHate()
		if info.SetHate(u64.Min(hateTypeData.MaxHate, info.GetHate()+toAdd)) {
			result.Add(region.NewS2cUpdateMultiLevelNpcHateMsg(int32(hateTypeData.Type), u64.Int32(info.GetHate())))
		}
	}

}

//func (t *troop) backHome(r *Realm, currentX, currentY int, startMovingTime time.Time) {
//	moveDuration := r.moveDuration(currentX, currentY, t.startingBase.BaseX(), t.startingBase.BaseY())
//
//	t.state = realmface.InvadeMovingBack
//	if t.targetBase != nil {
//		t.backHomeTargetX, t.backHomeTargetY = t.targetBase.BaseX(), t.targetBase.BaseY()
//		delete(t.targetBase.targetingTroops, t.Id())
//		t.targetBase = nil
//	}
//
//	t.moveStartTime = startMovingTime
//	t.moveArriveTime = startMovingTime.Add(moveDuration)
//	if t.event != nil {
//		t.event.RemoveFromQueue()
//	}
//	t.event = r.newEvent(t.moveArriveTime, t.returnedToBaseEvent)
//	t.proto = t.doEncode()
//	t.refreshMsg()
//}

func stateToActionMoveType(state realmface.TroopState) (actionType shared_proto.TroopOperate, moveType realmface.MoveType) {
	return state.Operate(), state.MoveType()
}

func (t *troop) TargetIsOwnerCanSee() bool {
	return t.ownerCanSeeTarget
}

// 从各个基地和event中删除这个部队.
// 并没有把部队从英雄中删除, 也没有把武将状态和士兵状态还原. 也没有发消息通知其他人
func (t *troop) removeWithoutReturnCaptain(r *Realm) {
	delete(t.startingBase.selfTroops, t.Id())
	t.RemoveRealmPoints()
	if t.targetBase != nil {
		delete(t.targetBase.targetingTroops, t.Id())
		t.targetBase.remindAttackOrRobCountChanged(r)
		t.assistDefendStartTime = time.Time{}
	}

	if t.event != nil {
		t.event.RemoveFromQueue()
		t.event = nil
	}
}

//func (t *troop) initProto(r *Realm) {
//	t.getProtoBytes(r) = must.Marshal(t.doEncode(r))
//}
//
//func (t *troop) refreshMsg(r *Realm) {
//	t.getProtoBytes(r) = must.Marshal(t.doEncode(r))
//
//	freeMsg(&t.updateMsgBoth)
//	freeMsg(&t.updateMsgMa)
//	freeMsg(&t.updateMsgRealm)
//}

func (t *troop) clearMsgs() {
	freeMsg(&t.updateMsgMa)
	freeMsg(&t.removeMsgMa)

	freeMsg(&t.addSeeMeMsg)
	freeMsg(&t.removeSeeMeMsg)
}

func freeMsg(msg *pbutil.Buffer) {
	if *msg != nil {
		//(*msg).DoFreeEvenItsStaticAndFuckMeIfItExplodes()
		(*msg) = nil
	}
}

const (
	addTroopTypeCanSee = 0
	addTroopTypeInvate = 1
	addTroopTypeUpdate = 2

	removeTroopTypeCantSee = 0
	removeTroopTypeDestroy = 1
)

//kimi 修改了函数名,与接口一致
func (t *troop) GetCantSeeMeMsg() pbutil.Buffer {
	if t.removeSeeMeMsg == nil {
		t.removeSeeMeMsg = t.newRemoveSeeMeMsg(removeTroopTypeCantSee)
	}

	return t.removeSeeMeMsg
}

func (t *troop) newRemoveSeeMeMsg(removeType int32) pbutil.Buffer {
	return region.NewS2cRemoveTroopUnitMsg(removeType, t.IdBytes())
}

func (t *troop) getRemoveMsgMa() pbutil.Buffer {
	if t.removeMsgMa == nil {
		t.removeMsgMa = region.NewS2cRemoveMilitaryInfoMsg(t.IdBytes(), false, true).Static()
	}

	return t.removeMsgMa
}

func (t *troop) getRemoveMsgToSelf() pbutil.Buffer {
	return region.NewS2cRemoveSelfMilitaryInfoMsg(t.IdBytes())
}

func (t *troop) onChanged() {
	t.protoChanged = true
}

func (t *troop) getProtoBytes(r *Realm) []byte {
	if t.protoChanged {
		t.protoChanged = false

		t.protoBytesCache = must.Marshal(t.doEncode(r))

		freeMsg(&t.updateMsgMa)
		freeMsg(&t.addSeeMeMsg)
	}

	return t.protoBytesCache
}

func (t *troop) getAddSeeMeMsg(r *Realm) pbutil.Buffer {
	if t.protoChanged || t.addSeeMeMsg == nil {
		t.getProtoBytes(r)
		t.addSeeMeMsg = t.newAddSeeMeMsg(r, addTroopTypeCanSee)
	}

	return t.addSeeMeMsg
}

func (t *troop) newAddSeeMeMsg(r *Realm, addType int32) pbutil.Buffer {
	return region.NewS2cAddTroopUnitMarshalMsg(addType, t.doEncodeUnit(r))
}

func (t *troop) getUpdateMsgMa(r *Realm) pbutil.Buffer {
	if t.protoChanged || t.updateMsgMa == nil {
		t.updateMsgMa = region.NewS2cUpdateMilitaryInfoMsg(t.getProtoBytes(r), false, true).Static()
	}

	return t.updateMsgMa
}

func (t *troop) getUpdateMsgToSelf(r *Realm) pbutil.Buffer {
	return region.NewS2cUpdateSelfMilitaryInfoMsg(u64.Int32(entity.GetTroopIndex(t.Id())+1), t.IdBytes(), t.getProtoBytes(r))
}

func (t *troop) getStartingBaseReportProto(r *Realm) *shared_proto.ReportHeroProto {
	if t.owner != nil {
		return t.startingBase.newReportHeroProto(r, t.owner.Basic)
	}
	return t.startingBase.toReportHeroProto(r, t.startingBaseLevel)
}

func (t *troop) getTargetBaseReportProto(r *Realm) *shared_proto.ReportHeroProto {
	if t.targetBase != nil {
		return t.targetBase.toReportHeroProto(r, t.targetBaseLevel)
	}
	return nil
}

func (t *troop) toReportHeroProto(r *Realm, aliveMap, killMap map[int32]int32) *shared_proto.ReportHeroProto {
	proto := t.getStartingBaseReportProto(r)

	tfa := data.NewTroopFightAmount()
	for _, captain := range t.captains {
		//proto.TotalSoldier += captain.Proto().TotalSoldier
		proto.TotalSoldier += captain.Proto().Soldier // 不使用总兵力，使用开打前的兵力

		cp := &shared_proto.ReportCaptainProto{
			Index:         imath.Int32(captain.Index()),
			Captain:       captain.Proto(),
			CombatSoldier: captain.Proto().Soldier,
		}

		if cp.CombatSoldier > 0 {
			cp.FightAmount = data.ProtoFightAmount(captain.Proto().TotalStat, cp.CombatSoldier, captain.Proto().SpellFightAmountCoef)
			tfa.AddInt32(cp.FightAmount)

			if len(aliveMap) > 0 {
				alive := aliveMap[u64.Int32(captain.Id())]
				cp.AliveSoldier = alive
				proto.AliveSoldier += alive
			}

			if len(killMap) > 0 {
				cp.KillSoldier = killMap[u64.Int32(captain.Id())]
			}
		}

		// 武将详情
		proto.Captains = append(proto.Captains, cp)
	}
	proto.TotalFightAmount = tfa.ToI32()

	proto.AliveSoldier = i32.Min(proto.AliveSoldier, proto.TotalSoldier)

	// 城墙
	if t.wallStat != nil && t.wallStat.Strength > 0 {
		life := u64.Max(t.wallStat.Strength, 1)

		if life > 0 {
			proto.WallLevel = u64.Int32(t.wallLevel)
			proto.WallCombatSoldier = u64.Int32(life)

			if len(aliveMap) > 0 {
				proto.WallAliveSoldier = i32.Min(aliveMap[0], proto.WallCombatSoldier)
			}

			if len(killMap) > 0 {
				proto.WallKillSoldier = i32.Min(killMap[0], proto.WallCombatSoldier)
			}
		}
	}

	return proto
}

func (t *troop) toCombatPlayerProto(r *Realm) *shared_proto.CombatPlayerProto {

	player := &shared_proto.CombatPlayerProto{}
	player.Hero = t.getStartingBaseBasicProto(r)

	player.Troops = make([]*shared_proto.CombatTroopsProto, 0, len(t.captains))
	tfa := data.NewTroopFightAmount()
	for _, captain := range t.captains {
		if captain.Proto().Soldier <= 0 {
			continue
		}

		tps := &shared_proto.CombatTroopsProto{}
		tps.FightIndex = int32(captain.Index())
		tps.Captain = captain.Proto()

		player.Troops = append(player.Troops, tps)

		tfa.AddInt32(tps.Captain.FightAmount)
	}
	player.TotalFightAmount = tfa.ToI32()

	if t.wallStat != nil {
		player.WallStat = t.wallStat.Encode()
		player.WallFixDamage = u64.Int32(t.wallFixDamage)
		player.TotalWallLife = i32.Max(player.WallStat.Strength, 1)
		player.WallLevel = u64.Int32(t.wallLevel)
	}

	return player
}

func (t *troop) doEncodeUnit(r *Realm) *shared_proto.TroopUnitProto {
	proto := &shared_proto.TroopUnitProto{}

	proto.TroopId = t.IdBytes()

	actionType, moveType := stateToActionMoveType(t.state)
	proto.Action = int32(actionType)
	proto.MoveType = moveType.Int32()

	proto.MoveStartTime = timeutil.Marshal32(t.moveStartTime)
	proto.MoveArrivedTime = timeutil.Marshal32(t.moveArriveTime)

	proto.StartBaseId = t.startingBase.IdBytes()
	proto.StartBaseX = imath.Int32(t.startingBase.BaseX())
	proto.StartBaseY = imath.Int32(t.startingBase.BaseY())
	proto.StartBaseGuildId = i64.Int32(t.startingBase.GuildId())

	proto.TargetBaseId = idbytes.ToBytes(t.originTargetId)
	proto.TargetBaseX = imath.Int32(t.originTargetX)
	proto.TargetBaseY = imath.Int32(t.originTargetY)

	proto.Type = t.startingBase.TroopType()
	if t.assembly != nil && t.assembly.self == t {
		proto.Type = shared_proto.TroopType_TT_ASSEMBLY
	}

	if !npcid.IsNpcId(t.originTargetId) {
		target := r.services.heroSnapshotService.Get(t.originTargetId)
		if target != nil {
			proto.TargetBaseGuildId = i64.Int32(target.GuildId)
		}
	}

	proto.Dialogue = u64.Int32(t.dialogue)

	return proto
}

func (t *troop) doEncodeToServer(r *Realm) *server_proto.TroopServerProto {
	proto := &server_proto.TroopServerProto{}

	proto.StartingBaseId = t.startingBase.Id()
	proto.StartingBaseLevel = t.startingBaseLevel

	proto.RealmId = r.Id()

	proto.Id = t.Id()
	if t.targetBase != nil {
		proto.TargetBaseId = t.targetBase.Id()
	}
	proto.TargetBaseLevel = t.targetBaseLevel
	proto.TroopState = uint32(t.state)

	proto.OriginTargetId = t.originTargetId
	proto.BackHomeTargetX = int32(t.BackHomeTargetX())
	proto.BackHomeTargetY = int32(t.BackHomeTargetY())

	proto.CreateTime = timeutil.Marshal64(t.createTime)
	proto.MoveStartTime = timeutil.Marshal64(t.moveStartTime)
	proto.MoveArriveTime = timeutil.Marshal64(t.moveArriveTime)
	proto.RobbingEndTime = timeutil.Marshal64(t.robbingEndTime)
	proto.NextReduceProsperityTime = timeutil.Marshal64(t.nextReduceProsperityTime)
	proto.NextAddHateTime = timeutil.Marshal64(t.nextAddHateTime)

	if t.accumRobPrizeProto != nil {
		proto.AccumRobPrize = t.accumRobPrizeProto
	}
	proto.AccumReduceProsperity = t.accumReduceProsperity

	proto.Captains = make([]*shared_proto.CaptainInfoProto, 0, len(t.captains))
	proto.CaptainIndex = make([]int32, 0, len(t.captains))
	for _, c := range t.captains {
		proto.Captains = append(proto.Captains, c.Proto())
		proto.CaptainIndex = append(proto.CaptainIndex, int32(c.Index()))
	}

	if t.mmd != nil {
		proto.MonsterMasterId = t.mmd.Id
	}

	proto.AssistDefendStartTime = timeutil.Marshal64(t.assistDefendStartTime)

	proto.MultiLevelMonsterCount = t.npcTimes

	return proto
}

func (t *troop) unmarshalCaptainSoldier(captainIndex, captainSoldier []int32) {
	n := imath.Min(len(captainIndex), len(captainSoldier))

	for _, c := range t.Captains() {
		var soldier int32
		for i := 0; i < n; i++ {
			if int(captainIndex[i]) == c.Index() {
				soldier = captainSoldier[i]
				break
			}
		}

		c.Proto().Soldier = i32.Min(soldier, c.Proto().TotalSoldier)
		c.Proto().FightAmount = data.ProtoFightAmount(c.Proto().TotalStat, c.Proto().Soldier, c.Proto().SpellFightAmountCoef)
	}
}

func (t *troop) doEncode(r *Realm) *shared_proto.MilitaryInfoProto {
	proto := &shared_proto.MilitaryInfoProto{}

	proto.RegionId = i64.Int32(t.startingBase.RegionID())
	proto.CombineId = t.IdBytes()

	actionType, moveType := stateToActionMoveType(t.state)
	proto.Action = int32(actionType)
	proto.MoveType = moveType.Int32()

	proto.CreateTime = timeutil.Marshal32(t.createTime)
	proto.MoveStartTime = timeutil.Marshal32(t.moveStartTime)
	proto.MoveArrivedTime = timeutil.Marshal32(t.moveArriveTime)
	proto.RobbingEndTime = timeutil.Marshal32(t.robbingEndTime)

	// ignore priority_action_id，改成从最傻逼的开打

	//proto.Self = t.startingBase.internalBase.getBaseInfoByLevel(t.startingBaseLevel).EncodeAsHeroBasicSnapshot(r.services.heroSnapshotService.Get)
	proto.Self = t.getStartingBaseSnapshot(r)
	proto.SelfIsTent = !t.startingBase.isHeroHomeBase()

	selfNpcDataId, selfNpcType := t.startingBase.internalBase.getNpcConfig()
	proto.SelfNpcDataId, proto.SelfNpcType = u64.Int32(selfNpcDataId), selfNpcType

	for _, captain := range t.captains {
		proto.Captains = append(proto.Captains, captain.Proto())
		proto.CaptainIndex = append(proto.CaptainIndex, int32(captain.Index()))
	}

	// 目标
	tb := t.targetBase
	if tb == nil {
		tb = r.getBase(t.originTargetId)
	}
	if tb != nil {
		proto.Target = tb.internalBase.getBaseInfoByLevel(t.targetBaseLevel).EncodeAsHeroBasicSnapshot(r.services.heroSnapshotService.Get)

		npcDataId, npcType := tb.internalBase.getNpcConfig()
		proto.TargetNpcDataId, proto.TargetNpcType = u64.Int32(npcDataId), npcType
	} else {

		if npcid.IsNpcId(t.originTargetId) {
			npcDataId := npcid.GetNpcDataId(t.originTargetId)
			npcType := npcid.GetNpcIdType(t.originTargetId)
			proto.TargetNpcDataId, proto.TargetNpcType = u64.Int32(npcDataId), npcType

			switch npcType {
			case npcid.NpcType_MultiLevelMonster:
				data := r.services.datas.GetRegionMultiLevelNpcData(npcDataId)
				if data != nil {
					proto.Target = data.GetLevelBaseData(t.targetBaseLevel).Npc.EncodeSnapshot(t.originTargetId, r.id, t.originTargetX, t.originTargetY)
				}
			default:
				data := r.services.datas.GetNpcBaseData(npcDataId)
				if data != nil {
					proto.Target = data.EncodeSnapshot(t.originTargetId, r.id, t.originTargetX, t.originTargetY)
				}
			}

		} else {
			targetSnapshot := r.services.heroSnapshotService.Get(t.originTargetId)
			if targetSnapshot != nil {
				proto.Target = targetSnapshot.EncodeClient()
			}
		}
	}

	proto.TargetBaseX = imath.Int32(t.originTargetX)
	proto.TargetBaseY = imath.Int32(t.originTargetY)

	proto.KillEnemy = t.killEnemy

	proto.AccumRobPrize = t.AccumRobPrize()

	proto.MultiLevelMonsterCount = u64.Int32(t.npcTimes)

	if t.assembly != nil {
		proto.AssemblyId = t.assembly.self.IdBytes()
		proto.AssemblyCount = u64.Int32(t.assembly.Count())
		proto.AssemblyTotalCount = u64.Int32(t.assembly.TotalCount())

		if npcid.IsNpcId(t.assembly.self.originTargetId) {
			proto.AssemblyTargetNpcType = npcid.GetNpcIdType(t.assembly.self.originTargetId)
		}
	}

	return proto
}

func (t *troop) doEncodeToXiongNu(r *Realm) *shared_proto.XiongNuTroopProto {
	proto := &shared_proto.XiongNuTroopProto{}

	//proto.Self = t.startingBase.internalBase.getBaseInfoByLevel(t.startingBaseLevel).EncodeAsHeroBasicSnapshot(r.services.heroSnapshotService.Get)
	proto.Self = t.getStartingBaseSnapshot(r)

	for _, captain := range t.captains {
		proto.Captains = append(proto.Captains, captain.Proto())
		proto.CaptainIndex = append(proto.CaptainIndex, int32(captain.Index()))
	}

	return proto
}

func (t *troop) FightAmount() uint64 {
	tfa := data.NewTroopFightAmount()
	for _, c := range t.captains {
		tfa.AddInt32(c.Proto().FightAmount)
	}

	return tfa.ToU64()
}

func (t *troop) StartingBase() realmface.Base {
	// 接口不能直接return
	if t.startingBase != nil {
		return t.startingBase.Base()
	}
	return nil
}

func (t *troop) TargetBase() realmface.Base {
	// 接口不能直接return
	if t.targetBase != nil {
		return t.targetBase.Base()
	}
	return nil
}

func (t *troop) State() realmface.TroopState {
	return t.state
}

func (t *troop) BackHomeTargetX() int {
	return t.originTargetX
}

func (t *troop) BackHomeTargetY() int {
	return t.originTargetY
}

func (t *troop) CreateTime() time.Time {
	return t.createTime
}

func (t *troop) MoveStartTime() time.Time {
	return t.moveStartTime
}

func (t *troop) MoveArriveTime() time.Time {
	return t.moveArriveTime
}

func (t *troop) RobbingEndTime() time.Time {
	return t.robbingEndTime
}

func (t *troop) getRobDuration(r *Realm, targetBaseInfo baseInfo) time.Duration {
	if npcid.IsXiongNuNpcId(t.startingBase.Id()) {
		return r.services.datas.ResistXiongNuMisc().RobbingDuration
	} else {
		return targetBaseInfo.BeenRobMaxDuration(r)
	}
}

func (t *troop) InitRobbing(robbingEndTime time.Time, r *Realm, targetBaseInfo baseInfo, ctime time.Time) (nextTickTime time.Time) {
	t.robbingEndTime = robbingEndTime

	robbingStartTime := t.robbingEndTime.Add(-t.getRobDuration(r, targetBaseInfo))
	nextTickTime = t.robbingEndTime

	if d := targetBaseInfo.BeenRobTickDuration(r); d > 0 {
		if t.nextAddPrizeTime.Before(ctime) {
			t.nextAddPrizeTime = timeutil.NextTickTime(robbingStartTime, ctime, d)
		}
		nextTickTime = timeutil.Min(nextTickTime, t.nextAddPrizeTime)
	}

	if d := targetBaseInfo.BeenRobLostProsperityDuration(r); d > 0 {
		if t.nextReduceProsperityTime.Before(ctime) {
			// 扣对方繁荣度
			t.nextReduceProsperityTime = timeutil.NextTickTime(robbingStartTime, ctime, d)
		}
		nextTickTime = timeutil.Min(nextTickTime, t.nextReduceProsperityTime)
	}

	if t.startingBase.isHeroHomeBase() && t.targetBase != nil && t.targetBase.isHeroHomeBase() {
		if d := r.services.datas.RegionConfig().RobBaowuTickDuration; d > 0 {
			// 人打人才会抢宝物
			if t.nextRobBaowuTime.Before(ctime) {
				t.nextRobBaowuTime = timeutil.NextTickTime(robbingStartTime, ctime, d)
			}
			nextTickTime = timeutil.Min(nextTickTime, t.nextRobBaowuTime)
		}
	}

	if hateData := targetBaseInfo.getHateData(); hateData != nil && hateData.HateTickDuration > 0 {
		if t.nextAddHateTime.Before(ctime) {
			// 加仇恨
			t.nextAddHateTime = timeutil.NextTickTime(robbingStartTime, ctime, hateData.HateTickDuration)
		}
		nextTickTime = timeutil.Min(nextTickTime, t.nextAddHateTime)
	}

	return
}

func (t *troop) NextReduceProsperityTime() time.Time {
	return t.nextReduceProsperityTime
}

func (t *troop) NextAddHateTime() time.Time {
	return t.nextAddHateTime
}

func (t *troop) NextRobBaowuTime() time.Time {
	return t.nextRobBaowuTime
}

func (t *troop) Carrying() (gold uint64, food uint64, wood uint64, stone uint64) {
	return t.accumRobPrize.GetUnsafeResource()
}

func (t *troop) UpdateAssistDefendStartTime(startTime time.Time) {
	t.assistDefendStartTime = startTime
}

//func (t *troop) ClearCarrying() (gold uint64, food uint64, wood uint64, stone uint64) {
//	gold = t.gold
//	t.gold = 0
//
//	food = t.food
//	t.food = 0
//
//	wood = t.wood
//	t.wood = 0
//
//	stone = t.stone
//	t.stone = 0
//
//	return
//}
//
//func (t *troop) JadeOre() uint64 {
//	return t.jadeOre
//}

func (t *troop) clearAccumRobPrizeProto() {
	t.accumRobPrizeProto = nil
}

func (t *troop) AccumRobPrize() *shared_proto.PrizeProto {
	if t.accumRobPrizeProto == nil {
		t.accumRobPrizeProto = t.accumRobPrize.Build().Encode()
	}
	return t.accumRobPrizeProto
}

func (t *troop) BuildAccumRobPrize() *resdata.Prize {
	return t.accumRobPrize.Build()
}

func (t *troop) AccumReduceProsperity() uint64 {
	return t.accumReduceProsperity
}

func (t *troop) Captains() []realmface.Captain {
	return t.captains
}

func (t *troop) AliveSoldier() uint64 {

	var s int32
	for _, c := range t.captains {
		s += c.Proto().Soldier
	}

	return u64.FromInt32(s)
}

func (t *troop) TotalAliveSoldier() uint64 {
	totalSoldier := t.AliveSoldier()

	if t.assembly != nil {
		for _, st := range t.assembly.member {
			if st != nil {
				totalSoldier += st.AliveSoldier()
			}
		}
	}

	return totalSoldier
}

func (t *troop) HasAliveSoldier() bool {

	for _, c := range t.captains {
		if c.Proto().Soldier > 0 {
			return true
		}
	}

	return false
}

const enemyCount = 3

func (t *troop) AddKillEnemy(flagHeroName string) {
	t.onChanged()

	if len(t.killEnemy) >= enemyCount {
		copy(t.killEnemy, t.killEnemy[1:])
		t.killEnemy[len(t.killEnemy)-1] = flagHeroName
		return
	}

	t.killEnemy = append(t.killEnemy, flagHeroName)
}

func (t *troop) AssemblyId() int64 {
	if t.assembly != nil {
		return t.assembly.self.Id()
	}
	return 0
}

func (t *troop) AssemblyTargetId() int64 {
	if t.assembly != nil {
		return t.assembly.self.originTargetId
	}
	return 0
}

func (t *troop) Dialogue() uint64 {
	return t.dialogue
}

func (t *troop) NpcTimes() uint64 {
	return t.npcTimes
}

//AddRealmPoint 添加地图索引点
func (t *troop) AddRealmPoint(r *Realm) {
	distantX, distantY := constants.RealmIndexBlockSize, constants.RealmIndexBlockSize
	pointX, pointY := t.startingBase.BaseX(), t.startingBase.BaseY()
	if pointX/constants.RealmIndexBlockSize > t.originTargetX/constants.RealmIndexBlockSize {
		distantX = -distantX
	}

	if pointY/constants.RealmIndexBlockSize > t.originTargetY/constants.RealmIndexBlockSize {
		distantY = -distantY
	}

	for {
		//添加索引点
		r.blockManager.AddTroopIndex(pointX, pointY, t)
		// logrus.Debugf("添加军队索引点[%d] [%d]", pointX, pointY)
		if pointX/constants.RealmIndexBlockSize != t.originTargetX/constants.RealmIndexBlockSize {
			pointX += distantX
		} else if pointY/constants.RealmIndexBlockSize == t.originTargetY/constants.RealmIndexBlockSize {
			break
		}
		if pointY/constants.RealmIndexBlockSize != t.originTargetY/constants.RealmIndexBlockSize {
			pointY += distantY
		}
	}
}

//RemoveRealmPoints 移除地图索引点
func (t *troop) RemoveRealmPoints() {
	for _, v := range t.realmTracks {
		if v != nil {
			v.RemoveTroopIndex(t)
			// logrus.Debugf("移除军队索引点")
		}
	}
	t.realmTracks = nil

}

// --- captain ---

type captain struct {
	idAtHero uint64

	// 打架时的index
	index int
	proto *shared_proto.CaptainInfoProto
}

func (c *captain) Index() int {
	return c.index
}

func (c *captain) Id() uint64 {
	return c.idAtHero
}

func (c *captain) Proto() *shared_proto.CaptainInfoProto {
	return c.proto
}

func (c *captain) GetAndSetSoldierCount(newCount uint64) (oldCount uint64) {
	oldCount = u64.FromInt32(c.proto.Soldier)

	c.proto.Soldier = u64.Int32(newCount)
	c.proto.FightAmount = data.ProtoFightAmount(c.proto.TotalStat, c.proto.Soldier, c.proto.SpellFightAmountCoef)
	return
}
