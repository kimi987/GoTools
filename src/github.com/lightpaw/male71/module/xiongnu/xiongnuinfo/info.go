package xiongnuinfo

import (
	"github.com/lightpaw/logrus"
	"github.com/lightpaw/male7/config/xiongnu"
	xnmsg "github.com/lightpaw/male7/gen/pb/xiongnu"
	"github.com/lightpaw/male7/pb/server_proto"
	"github.com/lightpaw/male7/pb/shared_proto"
	"github.com/lightpaw/male7/service/herosnapshot/snapshotdata"
	"github.com/lightpaw/male7/util/atomic"
	"github.com/lightpaw/male7/util/i64"
	"github.com/lightpaw/male7/util/idbytes"
	"github.com/lightpaw/male7/util/timeutil"
	"github.com/lightpaw/male7/util/u64"
	"time"
	"github.com/lightpaw/male7/config"
	"github.com/lightpaw/male7/module/xiongnu/xiongnuface"
	"sync"
	"github.com/lightpaw/male7/util/imath"
)

func NewResistXiongNuInfo(
	configDatas *config.ConfigDatas, // 配置
	guildId int64,                   // 联盟id
	data *xiongnu.ResistXiongNuData, // 抗击匈奴数据
	defenders []int64,               // 防守者
	startTime time.Time,
) *resistXiongNuInfo {
	return &resistXiongNuInfo{
		misc:                configDatas.ResistXiongNuMisc(),
		configDatas:         configDatas,
		data:                data,
		guildId:             guildId,
		defenders:           defenders,
		startTime:           startTime,
		wipeOutMonsterCount: atomic.NewUint64(0),
		accumReduceMorale:   atomic.NewUint64(0),
		syncChange:          atomic.NewBool(false),
	}
}

// 活动信息
type resistXiongNuInfo struct {
	misc        *xiongnu.ResistXiongNuMisc // 杂项
	data        *xiongnu.ResistXiongNuData // 抗击匈奴数据
	configDatas *config.ConfigDatas        // 配置

	guildId      int64 // 联盟id
	baseId       int64 // 主城在场景中的id
	baseX, baseY int32

	startTime time.Time // 开启时间

	wipeOutMonsterCount *atomic.Uint64 // 消灭的怪物数量
	//morale              *atomic.Uint64 // 当前士气
	accumReduceMorale *atomic.Uint64 // 累计减少的士气值

	givePrizeMembers []int64 // 给奖励的成员列表

	defenders []int64 // 防守者

	// 防守战力，多线程访问
	defenserFightAmountRef atomic.Value

	wave                   uint64    // 第几波
	addMonsterCount        uint64    // 在场景中添加的怪物数量
	nextRefreshMonsterTime time.Time // 下一次刷新怪物的时间

	defeated bool // 是不是打败了匈奴
	resist   bool // 是否开始反击了

	syncChange *atomic.Bool // 是否需要同步变化

	// 玩家击杀数据
	heroFightMap sync.Map
}

func (info *resistXiongNuInfo) UpdateDefenserFightAmount(defenserId []int64, newFightAmounts, newEnemyCounts []uint64,
	newExpireTime time.Time, old *xiongnuface.DefenserFightAmount) *xiongnuface.DefenserFightAmount {

	newData := newRoDefenserFightAmount(defenserId, newFightAmounts, newEnemyCounts, newExpireTime, old)
	info.defenserFightAmountRef.Store(newData)
	return newData
}

func (info *resistXiongNuInfo) GetDefenserFightAmount() *xiongnuface.DefenserFightAmount {
	ref := info.defenserFightAmountRef.Load()
	if ref != nil {
		return ref.(*xiongnuface.DefenserFightAmount)
	}
	return nil
}

func newRoDefenserFightAmount(defenserId []int64, newFightAmounts, newEnemyCounts []uint64, newExpireTime time.Time, old *xiongnuface.DefenserFightAmount) *xiongnuface.DefenserFightAmount {

	newRo := &xiongnuface.DefenserFightAmount{}
	newRo.ExpireTime = newExpireTime

	if old != nil && u64.Equals(old.FightAmounts, newFightAmounts) &&
		u64.Equals(old.EnemyCounts, newEnemyCounts) {
		newRo.Version = old.Version
		newRo.FightAmounts = old.FightAmounts
		newRo.EnemyCounts = old.EnemyCounts
		newRo.DiffVersionMsg = old.DiffVersionMsg
		newRo.SameVersionMsg = old.SameVersionMsg
	} else {
		newRo.Version = timeutil.Marshal32(newExpireTime)
		if old != nil && newRo.Version == old.Version {
			// 这里只能保证跟原来的version不一样
			// 可能同时有多个线程来创建新的，创建的新对象拥有相同的version
			// 但是不管，反正数据有效期很短，马上就会产生新的
			newRo.Version += 1
		}

		newRo.FightAmounts = newFightAmounts
		newRo.EnemyCounts = newEnemyCounts
		newRo.DiffVersionMsg = xnmsg.NewS2cGetDefenserFightAmountMsg(newRo.Version,
			idbytes.ToBytesArray(defenserId), u64.Int32Array(newFightAmounts), u64.Int32Array(newEnemyCounts)).Static()
		newRo.SameVersionMsg = xnmsg.NewS2cGetDefenserFightAmountMsg(newRo.Version,
			nil, nil, nil).Static()
	}

	return newRo
}

func (info *resistXiongNuInfo) Data() *xiongnu.ResistXiongNuData {
	return info.data
}

func (info *resistXiongNuInfo) GuildId() int64 {
	return info.guildId
}

func (info *resistXiongNuInfo) BaseId() int64 {
	return info.baseId
}

func (info *resistXiongNuInfo) BaseX() int32 {
	return info.baseX
}

func (info *resistXiongNuInfo) BaseY() int32 {
	return info.baseY
}

func (info *resistXiongNuInfo) SetBase(baseId int64, baseX, baseY int32) {
	info.baseId = baseId
	info.baseX = baseX
	info.baseY = baseY
	info.syncChange.Store(true)
}

func (info *resistXiongNuInfo) NextRefreshMonsterTime() time.Time {
	return info.nextRefreshMonsterTime
}

func (info *resistXiongNuInfo) SetNextRefreshMonsterTime(toSet time.Time) {
	info.nextRefreshMonsterTime = toSet
}

func (info *resistXiongNuInfo) AddMonsterCount() uint64 {
	return info.addMonsterCount
}

func (info *resistXiongNuInfo) IncAddMonsterCount(toAdd uint64) {
	info.addMonsterCount = u64.Plus(info.addMonsterCount, toAdd)
}

func (info *resistXiongNuInfo) SetAddMonsterCount(toSet uint64) {
	info.addMonsterCount = toSet
}

func (info *resistXiongNuInfo) SetDefeated() {
	info.defeated = true
}

func (info *resistXiongNuInfo) IsDefeated() bool {
	return info.defeated
}

func (info *resistXiongNuInfo) SetResist() {
	info.resist = true
}

func (info *resistXiongNuInfo) IsResist() bool {
	return info.resist
}

func (info *resistXiongNuInfo) Defenders() []int64 {
	return info.defenders
}

func (info *resistXiongNuInfo) WipeOutMonsterCount() uint64 {
	return info.wipeOutMonsterCount.Load()
}

// 消灭怪物
func (info *resistXiongNuInfo) WipeOutMonster(ctime time.Time, toReduceMorale uint64) {
	info.wipeOutMonsterCount.Inc()

	if toReduceMorale > 0 {
		info.accumReduceMorale.Add(toReduceMorale)
	}
	info.syncChange.Store(true)
}

func (info *resistXiongNuInfo) Morale() (morale uint64) {
	return u64.Sub(info.configDatas.ResistXiongNuMisc().MaxMorale, info.accumReduceMorale.Load())
}

func (info *resistXiongNuInfo) GivePrizeMembers() []int64 {
	return info.givePrizeMembers
}

func (info *resistXiongNuInfo) AddGivePrizeMember(id int64) {
	info.givePrizeMembers = append(info.givePrizeMembers, id)
	info.syncChange.Store(true)
}

func (info *resistXiongNuInfo) StartTime() time.Time {
	return info.startTime
}

func (info *resistXiongNuInfo) InvadeEndTime() time.Time {
	return info.startTime.Add(info.misc.InvadeDuration)
}

func (info *resistXiongNuInfo) ResistTime() time.Time {
	return info.startTime.Add(info.misc.InvadeDuration)
}

func (info *resistXiongNuInfo) EndTime() time.Time {
	return info.startTime.Add(info.misc.InvadeDuration).Add(info.misc.ResistDuration).Add(-time.Second) // 减少一秒，客户端怕倒计时完了，还没给他发结束，导致他界面出问题
}

func (info *resistXiongNuInfo) NextWaveTime() time.Time {
	if info.wave >= uint64(len(info.data.ResistWaves)) {
		return time.Time{}
	}

	return info.startTime.Add(info.misc.GetInvadeWaveDuration(info.wave))
}

func (info *resistXiongNuInfo) NeedSyncChange() bool {
	return info.syncChange.Swap(false)
}

func (info *resistXiongNuInfo) IncWave() {
	info.wave++
	info.addMonsterCount = 0
	info.syncChange.Store(true)
}

func (info *resistXiongNuInfo) Wave() uint64 {
	return info.wave
}

func (info *resistXiongNuInfo) AddHeroFightSoldier(heroId int64, killSoldier, beenKilledSoldier uint64) {

	obj := info.getOrCreateHeroFightObject(heroId)
	obj.killSoldier.Add(killSoldier)
	obj.beenKilledSoldier.Add(beenKilledSoldier)

}

func (info *resistXiongNuInfo) GetHeroFightObject(heroId int64) *HeroFightObject {
	if o, _ := info.heroFightMap.Load(heroId); o != nil {
		return o.(*HeroFightObject)
	}
	return nil
}

func (info *resistXiongNuInfo) getOrCreateHeroFightObject(heroId int64) *HeroFightObject {
	if o, _ := info.heroFightMap.Load(heroId); o != nil {
		return o.(*HeroFightObject)
	}

	hero := newHeroFightObject(heroId)
	if actual, loaded := info.heroFightMap.LoadOrStore(heroId, hero); loaded {
		return actual.(*HeroFightObject)
	}
	return hero
}

func (info *resistXiongNuInfo) RangeHeroFightObject(f func(obj *HeroFightObject) bool) {
	info.heroFightMap.Range(func(key, value interface{}) bool {
		obj, ok := value.(*HeroFightObject)
		if ok {
			return f(obj)
		}
		return true
	})
}

func newHeroFightObject(heroId int64) *HeroFightObject {
	return &HeroFightObject{
		heroId:            heroId,
		killSoldier:       atomic.NewUint64(0),
		beenKilledSoldier: atomic.NewUint64(0),
	}
}

type HeroFightObject struct {
	heroId int64

	// 击杀士兵数
	killSoldier *atomic.Uint64

	// 士兵被击杀数
	beenKilledSoldier *atomic.Uint64
}

func (info *resistXiongNuInfo) Encode() *server_proto.XiongNuServerProto {
	proto := &server_proto.XiongNuServerProto{}

	proto.Level = info.data.Level
	proto.GuildId = info.guildId
	proto.BaseId = info.baseId
	proto.BaseX = info.baseX
	proto.BaseY = info.baseY
	proto.StartTime = timeutil.Marshal64(info.startTime)
	proto.WipeOutMonsterCount = info.wipeOutMonsterCount.Load()
	proto.GivePrizeMembers = info.givePrizeMembers
	proto.Defenders = info.defenders
	proto.Wave = info.wave
	proto.AddMonsterCount = info.addMonsterCount
	proto.NextRefreshMonsterTime = timeutil.Marshal64(info.nextRefreshMonsterTime)
	proto.AccumReduceMorale = info.accumReduceMorale.Load()
	proto.Defeated = info.defeated
	proto.Resist = info.resist

	info.RangeHeroFightObject(func(obj *HeroFightObject) bool {
		proto.HeroId = append(proto.HeroId, obj.heroId)
		proto.HeroKillSoldier = append(proto.HeroKillSoldier, obj.killSoldier.Load())
		proto.HeroBeenKillSoldier = append(proto.HeroBeenKillSoldier, obj.beenKilledSoldier.Load())
		return true
	})

	return proto
}

func (info *resistXiongNuInfo) EncodeClient() *shared_proto.XiongNuProto {
	proto := &shared_proto.XiongNuProto{}

	proto.Level = u64.Int32(info.data.Level)
	proto.BaseId = i64.Int32(info.baseId)
	proto.BaseX = info.baseX
	proto.BaseY = info.baseY
	proto.StartTime = timeutil.Marshal32(info.startTime)
	proto.WipeOutMonsterCount = u64.Int32(info.wipeOutMonsterCount.Load())
	proto.Morale = u64.Int32(info.Morale())
	proto.Wave = u64.Int32(info.wave)
	proto.NextWaveTime = timeutil.Marshal32(info.NextWaveTime())
	proto.CanResistIds = make([][]byte, 0, len(info.givePrizeMembers))
	for _, defenderId := range info.givePrizeMembers {
		bytes := idbytes.ToBytes(defenderId)
		proto.CanResistIds = append(proto.CanResistIds, bytes)
	}

	return proto
}

func (info *resistXiongNuInfo) Unmarshal(proto *server_proto.XiongNuServerProto) {
	info.givePrizeMembers = proto.GivePrizeMembers
	info.baseId = proto.BaseId
	info.baseX = proto.BaseX
	info.baseY = proto.BaseY
	info.wipeOutMonsterCount.Store(proto.WipeOutMonsterCount)
	info.wave = proto.Wave
	info.addMonsterCount = proto.AddMonsterCount
	info.nextRefreshMonsterTime = timeutil.Unix64(proto.NextRefreshMonsterTime)
	info.accumReduceMorale.Store(proto.AccumReduceMorale)
	info.defeated = proto.Defeated
	info.resist = proto.Resist
	info.syncChange.Store(true)

	n := imath.Minx(len(proto.HeroId), len(proto.HeroKillSoldier), len(proto.HeroBeenKillSoldier))
	if n > 0 {
		for i := 0; i < n; i++ {
			obj := info.getOrCreateHeroFightObject(proto.HeroId[i])
			obj.killSoldier.Store(proto.HeroKillSoldier[i])
			obj.beenKilledSoldier.Store(proto.HeroBeenKillSoldier[i])
		}
	}
}

func (info *resistXiongNuInfo) EncodeLast(scoreLevel uint64, heroSnapshotGetter func(id int64) *snapshotdata.HeroSnapshot) *shared_proto.LastResistXiongNuProto {
	proto := &shared_proto.LastResistXiongNuProto{}

	proto.Level = u64.Int32(info.data.Level)
	proto.WipeOutMonsterCount = u64.Int32(info.wipeOutMonsterCount.Load())
	proto.DestroyBase = info.defeated
	proto.Defenders = make([]*shared_proto.HeroBasicProto, 0, len(info.defenders))
	proto.Grade = u64.Int32(scoreLevel)

	for _, defenderId := range info.defenders {
		defenderSnapshot := heroSnapshotGetter(defenderId)
		if defenderSnapshot != nil {
			proto.Defenders = append(proto.Defenders, defenderSnapshot.EncodeBasic4Client())
		} else {
			logrus.Debugf("战斗结束，没找到防守成员镜像数据: %d", defenderId)
		}
	}

	return proto
}

func (info *resistXiongNuInfo) EncodeFight(heroSnapshotGetter func(id int64) *snapshotdata.HeroSnapshot) *shared_proto.ResistXiongNuFightProto {
	proto := &shared_proto.ResistXiongNuFightProto{}

	info.RangeHeroFightObject(func(obj *HeroFightObject) bool {
		if hero := heroSnapshotGetter(obj.heroId); hero != nil {
			proto.Hero = append(proto.Hero, hero.EncodeBasic4Client())
			proto.KillSoldier = append(proto.KillSoldier, u64.Int32(obj.killSoldier.Load()))
			proto.BeenKillSoldier = append(proto.BeenKillSoldier, u64.Int32(obj.beenKilledSoldier.Load()))
		}
		return true
	})

	for _, defenserId := range info.defenders {
		proto.Defenser = append(proto.Defenser, idbytes.ToBytes(defenserId))
	}

	return proto
}
