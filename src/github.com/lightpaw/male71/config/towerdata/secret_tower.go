package towerdata

import (
	"github.com/lightpaw/male7/config/monsterdata"
	"github.com/lightpaw/male7/config/resdata"
	"github.com/lightpaw/male7/config/scene"
	"github.com/lightpaw/male7/gen/pb/secret_tower"
	"github.com/lightpaw/male7/pb/shared_proto"
	"github.com/lightpaw/male7/util/check"
	"github.com/lightpaw/male7/util/u64"
	"github.com/lightpaw/pbutil"
	"math/rand"
	"time"
)

const denominator int64 = 10000

//gogen:config
type SecretTowerData struct {
	_ struct{} `file:"系统模块/重楼密室.txt"`
	_ struct{} `proto:"shared_proto.SecretTowerDataProto"`
	_ struct{} `protoconfig:"SecretTower"`

	// 楼层，从1开始
	Id uint64 `key:"1"`

	Name string `validator:"string>0"`

	// 首胜奖励
	FirstPassPrize *resdata.Prize

	// 通关奖励
	Plunder   *resdata.Plunder `protofield:"-"`
	ShowPrize *resdata.Prize   `protofield:"Prize"`

	Image string `default:""`

	// 超级大奖
	SuperPlunder     *resdata.Plunder `default:"nullable" protofield:"-"`
	SuperPlunderRate uint64           `validator:"int>0" protofield:"-"`
	SuperShowPrize   *resdata.Prize   `default:"nullable" protofield:"SuperPrize"`

	// 联盟协助贡献
	GuildHelpContribution uint64 `validator:"int>0"`

	// 盟友数量、连胜数、盟友连胜数有关
	//GuildMemberCountContribution            []uint64 `validator:"int>0" protofield:"-"` // 盟友数量
	//ContinueWinTimesContribution            []uint64 `validator:"int>0" protofield:"-"` // 连胜数
	//GuildMemberContinueWinTimesContribution []uint64 `validator:"int>0" protofield:"-"` // 盟友连胜数

	// 在有玩家加入队伍后，队伍的开启保护时间间隔，就是有人进来了，几秒内不可以开始游戏
	StartProtectDuration time.Duration `protofield:"-"`

	// 最大的进攻方数量
	MaxAttackerCount uint64 `validator:"int>0"`

	// 最少的进攻方数量
	MinAttackerCount uint64 `validator:"uint"`

	// 同时能够进行的战斗数量
	ConcurrentFightCount uint64 `validator:"int>0" protofield:"-"`

	// 进攻方最大的连胜次数
	MaxAttackerContinuewWinTimes uint64 `validator:"uint"`

	MaxDefenserContinuewWinTimes uint64 `validator:"uint" protofield:"-" default:"0"`

	// 描述
	Desc string

	MonsterLeaderId uint64 // 防守武将队长id

	// 防守武将
	Monster []*monsterdata.MonsterMasterData

	// 战斗场景
	CombatScene *scene.CombatScene `protofield:"-"`

	// 解锁需要的千重楼
	UnlockTowerData *TowerData    `protofield:",config.U64ToI32(%s.Floor)"`
	UnlockMsg       pbutil.Buffer `head:"-" protofield:"-"`

	TeamExpireDuration time.Duration `default:"15m" protofield:"-"` // 队伍过期时间
}

func (d *SecretTowerData) Init(filename string) {
	d.UnlockMsg = secret_tower.NewS2cUnlockSecretTowerMsg(u64.Int32(d.Id))

	check.PanicNotTrue(d.UnlockTowerData.UnlockSecretTower == nil, "%s 配置的密室 %d 被多个千重楼解锁!%v", filename, d.Id, d.UnlockTowerData)

	d.UnlockTowerData.UnlockSecretTower = d

	check.PanicNotTrue(d.MinAttackerCount <= d.MaxAttackerCount, "%s 配置的密室 %d， 配置的最小进攻方人数必须[%d]<=最大进攻方人数[%d]", filename, d.Id, d.MinAttackerCount, d.MaxAttackerCount)

	check.PanicNotTrue(len(d.Monster) > 0, "%s 配置的密室 %d， 配置的怪物npc数量[%d]必须>=1", filename, d.Id, len(d.Monster))

	//check.PanicNotTrue(len(d.GuildMemberCountContribution) == int(d.MaxAttackerCount),
	//	"%s 配置的密室 %d 配置的盟友数量联盟贡献奖励配置的条数必须跟最大进攻方人数一样多!%d, %d", filename, d.Id, len(d.GuildMemberCountContribution), int(d.MaxAttackerCount))
	//
	//maxContinueWinTimes := int(d.MaxAttackerContinuewWinTimes)
	//if d.MaxAttackerContinuewWinTimes == 0 {
	//	maxContinueWinTimes = len(d.Monster)
	//}
	//check.PanicNotTrue(len(d.ContinueWinTimesContribution) == maxContinueWinTimes,
	//	"%s 配置的密室 %d 配置的连胜联盟贡献奖励配置的条数必须跟最大连胜次数一样多!%d, %d", filename, d.Id, len(d.ContinueWinTimesContribution), maxContinueWinTimes)
	//check.PanicNotTrue(len(d.GuildMemberContinueWinTimesContribution) == maxContinueWinTimes,
	//	"%s 配置的密室 %d 配置的盟友连胜联盟贡献奖励配置的条数必须跟最大连胜次数一样多!%d, %d", filename, d.Id, len(d.GuildMemberContinueWinTimesContribution), maxContinueWinTimes)

	if d.MaxAttackerContinuewWinTimes > 0 {
		check.PanicNotTrue(d.MinAttackerCount*d.MaxAttackerContinuewWinTimes >= uint64(len(d.Monster)), "%s 配置的密室 %d 配置的最小进攻方数量[%d]*最大的进攻方连胜数必须[%d]>=敌方NPC人数[%d]", filename, d.Id, d.MinAttackerCount, d.MaxAttackerContinuewWinTimes, len(d.Monster))
	}

	leaderExist := false
	for _, m := range d.Monster {
		if m.Id == d.MonsterLeaderId {
			leaderExist = true
			break
		}
	}
	check.PanicNotTrue(leaderExist, "%s 配置的密室 %d 配置的武将首领id[%d]不存在!", filename, d.Id, d.MonsterLeaderId)
}

func (d *SecretTowerData) InitAll(filename string, array []*SecretTowerData) {
	for idx, data := range array {
		if idx > 0 {
			prev := array[idx-1]
			check.PanicNotTrue(prev.UnlockTowerData.Floor < data.UnlockTowerData.Floor, "%s 配置的密室 %d 密室id越大，解锁该密室的千重楼需要的层级要比id更小的要大!%v, %v", filename, data.Id, prev, data)
		}
	}
}

func (d *SecretTowerData) RandomGiveSuperPrize() (giveSuperPrize bool) {
	return rand.Int63n(denominator) <= int64(d.SuperPlunderRate)
}

//// 获得在联盟组队模式下的帮派贡献
//func (d *SecretTowerData) CalcGuildModeContribution(guildMemberCount int, continueWinTimes, maxMemberContinueWinTimes int64) (contribution uint64) {
//	if guildMemberCount > 0 {
//		if guildMemberCount > len(d.GuildMemberCountContribution) {
//			contribution += d.GuildMemberCountContribution[len(d.GuildMemberCountContribution)-1]
//		} else {
//			contribution += d.GuildMemberCountContribution[guildMemberCount-1]
//		}
//	}
//
//	if continueWinTimes > 0 {
//		if int(continueWinTimes) > len(d.ContinueWinTimesContribution) {
//			contribution += d.ContinueWinTimesContribution[len(d.ContinueWinTimesContribution)-1]
//		} else {
//			contribution += d.ContinueWinTimesContribution[continueWinTimes-1]
//		}
//	}
//
//	if maxMemberContinueWinTimes > 0 {
//		if int(maxMemberContinueWinTimes) > len(d.ContinueWinTimesContribution) {
//			contribution += d.ContinueWinTimesContribution[len(d.ContinueWinTimesContribution)-1]
//		} else {
//			contribution += d.ContinueWinTimesContribution[maxMemberContinueWinTimes-1]
//		}
//	}
//
//	return
//}

func (d *SecretTowerData) GenMonsterCombatPlayer() (defenserIds []int64, defenserProtos []*shared_proto.CombatPlayerProto) {
	defenserIds = make([]int64, 0, len(d.Monster))
	defenserProtos = make([]*shared_proto.CombatPlayerProto, 0, len(d.Monster))

	for _, m := range d.Monster {
		defenserIds = append(defenserIds, m.GetNpcId())
		defenserProtos = append(defenserProtos, m.GetPlayer())
	}

	return
}

//gogen:config
type SecretTowerMiscData struct {
	_ struct{} `singleton:"true"`
	_ struct{} `file:"系统模块/重楼密室杂项.txt"`
	_ struct{} `proto:"shared_proto.SecretTowerMiscProto"`
	_ struct{} `protoconfig:"SecretTowerMisc"`

	MaxTimes             uint64 `validator:"int>0"`                              // 最大的次数
	MaxHelpTimes         uint64 `validator:"int>0"`                              // 最大的协助次数
	MaxGuildContribution uint64 `validator:"int>0" protofield:"-" default:"100"` // 最大能够获得的联盟贡献
	MaxRecord            uint64 `default:"10"`
}

//gogen:config
type SecretTowerWordsData struct {
	_ struct{} `file:"系统模块/重楼密室聊天.txt"`
	_ struct{} `protogen:"true"`

	Id    uint64
	Words string `desc:"气泡说话文字"`
}
