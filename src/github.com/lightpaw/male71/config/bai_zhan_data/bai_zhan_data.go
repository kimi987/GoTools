package bai_zhan_data

import (
	"github.com/lightpaw/male7/config/data"
	"github.com/lightpaw/male7/config/i18n"
	"github.com/lightpaw/male7/config/icon"
	"github.com/lightpaw/male7/config/monsterdata"
	"github.com/lightpaw/male7/config/resdata"
	"github.com/lightpaw/male7/config/scene"
	"github.com/lightpaw/male7/gen/pb/bai_zhan"
	"github.com/lightpaw/male7/util/check"
	"github.com/lightpaw/male7/util/timeutil"
	"github.com/lightpaw/male7/util/u64"
	"github.com/lightpaw/pbutil"
	"time"
)

// 军衔等级数据
//gogen:config
type JunXianLevelData struct {
	_ struct{} `file:"百战千军/等级.txt"`
	_ struct{} `proto:"shared_proto.JunXianLevelDataProto"`
	_ struct{} `protoconfig:"jun_xian_level"`

	Level uint64 `validator:"int>0"` // 军衔等级

	Name *i18n.I18nRef `validator:"string>0"` // 军衔名字

	Icon *icon.Icon `protofield:"IconId,%s.Id"` // 图标

	StrongMatchNpcGuardCaptain []*monsterdata.MonsterMasterData `default:"nullable" protofield:"-"` // 强制匹配的NPC守卫，长度可能为空
	NpcGuardCaptain            []*monsterdata.MonsterMasterData `protofield:"-"`                    // NPC守卫

	PrevLevel *JunXianLevelData `head:"-" protofield:"-"` // 上一级军衔
	NextLevel *JunXianLevelData `head:"-" protofield:"-"` // 下一级军衔

	ReaddJunXianLevelData *JunXianLevelData `protofield:"-"` // 重新加回来的军衔等级数据

	LevelUpPercent uint64 `validator:"uint"` // 等级上升的百分比，最后一级没有
	LevelUpPoint   uint64 `validator:"uint"` // 等级上升的分数

	LevelDownPercent uint64 `validator:"uint"` // 等级下降的百分比，第一级没有
	LevelDownPoint   uint64 `validator:"uint"` // 等级下降的分数

	DailySalary *resdata.Prize // 每日俸禄
	DailyHufu   uint64 `validator:"uint"`

	AssemblyStat *data.SpriteStat // 集结属性

	MinKeepLevelCount uint64 `protofield:"-"` // 最小的保级人数

	CombatScene *scene.CombatScene `protofield:"-"` // 战斗场景

	MaxJunXianLevelChangedMsg pbutil.Buffer `head:"-" protofield:"-"`
}

func (data *JunXianLevelData) Init() {
	data.MaxJunXianLevelChangedMsg = bai_zhan.NewS2cMaxJunXianLevelChangedMsg(u64.Int32(data.Level)).Static()
}

func (*JunXianLevelData) InitAll(filename string, configDatas interface {
	GetJunXianLevelDataArray() []*JunXianLevelData
}) {
	array := configDatas.GetJunXianLevelDataArray()
	for idx, data := range array {
		check.PanicNotTrue(uint64(idx+1) == data.Level, "军衔等级配置 %s level=%v 等级必须从1开始逐级递增!", filename, data.Level)

		check.PanicNotTrue(len(data.NpcGuardCaptain) > 0, "军衔等级配置 %s level=%v NPC守卫必须配置!", filename, data.Level)

		check.PanicNotTrue(data.ReaddJunXianLevelData.Level <= data.Level, "军衔等级配置 %s 配置的移除军衔后加回来的军衔等级 [%d] 只能<=当前军衔等级 [%d] 还要低啊!", filename, data.ReaddJunXianLevelData.Level, data.Level)

		if idx != 0 {
			data.PrevLevel = array[idx-1]
			check.PanicNotTrue(data.LevelDownPercent > 0 && data.LevelDownPercent < 100, "军衔等级配置 level=%d, level_down_percent [%d] 必须>0且<100!", data.Level, data.LevelDownPercent)
			check.PanicNotTrue(data.LevelDownPoint > 0, "军衔等级配置 level=%d, level_down_point [%d] 必须>0!", data.Level, data.LevelDownPoint)
			if idx != len(array)-1 {
				check.PanicNotTrue(data.LevelUpPoint >= data.LevelDownPoint, "军衔等级配置 %d 每一级的等级上升分数 [%d] 必须>每一级的等级下降分数 [%d]!", data.Level, data.LevelUpPoint, data.LevelDownPoint)
			}
		}

		if idx != len(array)-1 {
			data.NextLevel = array[idx+1]
			check.PanicNotTrue(data.LevelUpPercent > 0 && data.LevelUpPercent < 100, "军衔等级配置 level=%d, level_up_percent [%d] 必须>0且<100!", data.Level, data.LevelUpPercent)
			check.PanicNotTrue(data.LevelUpPoint > 0, "军衔等级配置 level=%d, level_up_point [%d] 必须>0!", data.Level, data.LevelUpPoint)
		}
	}
}

// 军衔等级数据
//gogen:config
type JunXianPrizeData struct {
	_ struct{} `file:"百战千军/等级奖励.txt"`
	_ struct{} `proto:"shared_proto.JunXianLevelPrizeProto"`
	_ struct{} `protoconfig:"jun_xian_level_prize"`

	Id uint64 `validator:"int>0"` // id

	LevelData *JunXianLevelData `head:"level" protofield:"Level,config.U64ToI32(%s.Level)"` // 需要的军衔等级

	Point uint64 `validator:"uint"` // 需要的积分

	Prize        *resdata.Prize           // 达成的奖励
	CollectedMsg pbutil.Buffer `head:"-"` // 领取奖励的消息
}

func (d *JunXianPrizeData) Init() {
	d.CollectedMsg = bai_zhan.NewS2cCollectJunXianPrizeMsg(u64.Int32(d.Id)).Static()
}

func (d *JunXianPrizeData) InitAll(filename string, array []*JunXianPrizeData) {
	for idx, prizeData := range array {
		if idx == 0 {
			continue
		}

		check.PanicNotTrue(prizeData.Id == uint64(idx+1), "%s 中配置的军衔奖励数据必须从1开始，逐级加1!%+v", filename, prizeData)

		prevPrizeData := array[idx-1]
		check.PanicNotTrue(prevPrizeData.LevelData.Level <= prizeData.LevelData.Level,
			"%s 中配置的数据，id越大，在上一层奖励需要的军衔等级 [%d] 必须<=下一层奖励需要的军衔等级 [%d] ，同时id必须从1开始，逐级加1!%v",
			filename, prevPrizeData.LevelData.Level, prizeData.LevelData.Level, prizeData)

		if prevPrizeData.LevelData.Level == prizeData.LevelData.Level {
			check.PanicNotTrue(prevPrizeData.Point < prizeData.Point, "%s 中配置的数据，id越大，在等级相同的情况下积分必须越配越大，同时id必须从1开始，逐级加1!%v", filename, prizeData)
		}
	}
}

// 百战其他数据
//gogen:config
type BaiZhanMiscData struct {
	_ struct{} `singleton:"true"`
	_ struct{} `file:"百战千军/杂项.txt"`
	_ struct{} `proto:"shared_proto.BaiZhanMiscProto"`
	_ struct{} `protoconfig:"BaiZhanMisc"`
	_ timeutil.CycleTime

	WinPoint  uint64 `validator:"int>0"` // 胜利加的积分
	FailPoint uint64 `validator:"int>0"` // 失败加的积分

	InitTimes        uint64          `validator:"uint" default:"0" protofield:"-"`
	RecoverTimesTime []time.Duration `protofield:",timeutil.DurationArrayToSecondArray(%s)"` // 恢复次数的时间
	RecoverTimes     []uint64        `validator:"int>0,duplicate"`                           // 加的次数

	MaxRecord uint64 `validator:"int>0" default:"20"` // 最大挑战记录保存条数

	ShowRankCount uint64 `validator:"int>0" default:"6"` // 排行榜上面展示的数量
}

func (d *BaiZhanMiscData) Init(filename string) {
	check.PanicNotTrue(d.WinPoint > d.FailPoint, "%s 配置的胜利加的积分必须比失败加的积分多!%d, %d", filename, d.WinPoint, d.FailPoint)
	check.PanicNotTrue(len(d.RecoverTimesTime) == len(d.RecoverTimes),
		"%s 配置的恢复次数的时间的长度跟恢复次数的长度必须一致!%d, %d", filename, len(d.RecoverTimesTime), len(d.RecoverTimes))

	for idx, duration := range d.RecoverTimesTime {
		check.PanicNotTrue(duration >= 0 && duration < 24*time.Hour,
			"恢复次数的时间间隔 [%d] 必须是>=重置间隔 [%d] 且<明日重置间隔 [%d]!", duration, 0, 24*time.Hour)

		if idx != 0 {
			check.PanicNotTrue(d.RecoverTimesTime[idx-1] < d.RecoverTimesTime[idx],
				"恢复次数的时间间隔配置必须越配置越大 [%d] 必须< [%d]!", d.RecoverTimesTime[idx-1], d.RecoverTimesTime[idx])
		}
	}
}

func (d *BaiZhanMiscData) GetAddPoint(isWin bool) (point uint64) {
	if isWin {
		return d.WinPoint
	} else {
		return d.FailPoint
	}
}

func (d *BaiZhanMiscData) GetCanChallengeTimes(since, offset time.Duration) (challengeTimes uint64) {
	// 算出当前时间跟重置时间点之后的时间间隔

	challengeTimes += d.InitTimes
	for idx, duration := range d.RecoverTimesTime {
		od := duration - offset
		if od < 0 {
			od += 24 * time.Hour
		}

		if since >= od {
			challengeTimes += d.RecoverTimes[idx]
		}
	}

	return
}
