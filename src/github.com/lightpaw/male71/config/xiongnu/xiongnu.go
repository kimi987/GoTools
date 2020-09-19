package xiongnu

import (
	"github.com/lightpaw/male7/config/basedata"
	"github.com/lightpaw/male7/config/monsterdata"
	"github.com/lightpaw/male7/config/resdata"
	"github.com/lightpaw/male7/config/singleton"
	"github.com/lightpaw/male7/util/check"
	"github.com/lightpaw/male7/util/timeutil"
	"github.com/lightpaw/male7/util/u64"
	"time"
)

// 抗击匈奴其他数据
//gogen:config
type ResistXiongNuMisc struct {
	_ struct{} `singleton:"true"`
	_ struct{} `file:"抗击匈奴/杂项.txt"`
	_ struct{} `proto:"shared_proto.ResistXiongNuMiscProto"`
	_ struct{} `protoconfig:"ResistXiongNuMisc"`
	_ timeutil.CycleTime

	DefenseMemberCount         uint64                                                                  // 防守人数
	InvadeDuration             time.Duration                                                           // 入侵时长
	ResistDuration             time.Duration                                                           // 反击时长
	InvadeWaveDuration         []time.Duration `protofield:",timeutil.DurationArrayToSecondArray(%s)"` // 入侵时间，开始后多久开始入侵
	TroopMoveVelocityPerSecond float64         `default:"0.25" protofield:"-"`                         // 移动速度
	MinMoveDuration            time.Duration   `protofield:"-"`                                        // 最小行军时长
	MaxMoveDuration            time.Duration   `protofield:"-"`                                        // 最大行军时长
	RobbingDuration            time.Duration   `protofield:"-"`                                        // 持续掠夺时长
	MaxMorale                  uint64                                                                  // 最大的士气
	WipeOutReduceMorale        uint64                                                                  // 消灭匈奴部队减少的士气
	OneMoraleReduceSoldierPer  uint64                                                                  // 一点士气减少士兵的万分比
	AttackDuration             time.Duration   `default:"3s" protofield:"-"`                           // 出兵间隔
	OpenNeedGuildLevel         uint64                                                                  // 开启需要联盟等级
	MaxCanOpenTimesPerDay      uint64          `default:"1"`                                           // 每天可以开启的次数
	MinBaseLevel               uint64          `default:"4"`                                           // 最小主城等级
	MaxDistance                uint64          `default:"500"`                                         // 与设置者的最大距离（里）(客户端用)
	StartAfterServerOpen       time.Duration   `default:"24h"`                                         // 开服多久才能开启匈奴

	BaseMinRange uint64 `default:"20" protofield:"-"`
	BaseMaxRange uint64 `default:"100" protofield:"-"`
}

func (data *ResistXiongNuMisc) GetMinRange() int {
	return u64.Int(u64.Min(data.BaseMinRange, data.BaseMaxRange))
}

func (data *ResistXiongNuMisc) GetMaxRange() int {
	return u64.Int(u64.Max(data.BaseMinRange, data.BaseMaxRange))
}

func (data *ResistXiongNuMisc) MoveDuration(distance float64) time.Duration {
	duration := singleton.MoveDuration(distance, data.TroopMoveVelocityPerSecond)
	if duration < data.MinMoveDuration {
		return data.MinMoveDuration
	}
	if duration > data.MaxMoveDuration {
		return data.MaxMoveDuration
	}
	return duration
}

func (data *ResistXiongNuMisc) GetInvadeWaveDuration(wave uint64) time.Duration {
	if len(data.InvadeWaveDuration) > 0 {
		idx := u64.Min(wave, uint64(len(data.InvadeWaveDuration)-1))
		return data.InvadeWaveDuration[idx]
	}
	return 0
}

// 抗击匈奴难度
//gogen:config
type ResistXiongNuData struct {
	_ struct{} `file:"抗击匈奴/难度.txt"`
	_ struct{} `proto:"shared_proto.ResistXiongNuDataProto"`
	_ struct{} `protoconfig:"ResistXiongNuData"`

	Level              uint64                                                                                                          // 难度等级
	Name               string                                                                                                          // 名字
	NpcBaseData        *basedata.NpcBaseData            `protofield:"Npc,config.U64ToI32(%s.Id)"`                                      // 主城数据
	AssistMonsters     []*monsterdata.MonsterMasterData `validator:"int,duplicate"`                                                    // 协助怪
	ResistWaves        []*ResistXiongNuWaveData         `protofield:"-"`                                                               // 攻城波次
	FirstResistWave    *ResistXiongNuWaveData           `head:"-" protofield:"-"`                                                      // 首次攻城波次
	ScorePrizes        []*resdata.Prize                 `validator:"int,duplicate" head:"-"`                                           // 给客户端的评分奖励，兼容旧版本
	ScorePlunderPrizes []*resdata.PlunderPrize          `validator:"int,duplicate" protofield:"-"`                                     // 评分奖励
	ScorePrestiges     []uint64                         `validator:"int>0,duplicate" protofield:"ScorePrestiges,config.U64a2I32a(%s)"` // 评分奖励联盟声望
	ResistSucPrize     *resdata.Prize                                                                                                  // 反击奖励
	ResistSucPrestige  uint64                                                                                                          // 反击奖励联盟声望
	NextLevel          *ResistXiongNuData               `head:"-" protofield:"-"`                                                      // 下级难度，可能为空
	TotalMonsterCount  uint64                           `head:"-"`                                                                     // 总的怪物数量
	MaxFightAmount     uint64                           `head:"-"`                                                                     // 最大战斗力

	ShowPrizes []*resdata.Prize `validator:"int,duplicate"`

	// 防守总兵力
	totalSoldier uint64

	maxWave uint64 // 最大波数
}

func (*ResistXiongNuData) InitAll(filename string, configs interface {
	GetResistXiongNuDataArray() []*ResistXiongNuData
	GetGuildLevelPrizeArray() []*resdata.GuildLevelPrize
}) {
	var prevLevel *ResistXiongNuData

	for idx, data := range configs.GetResistXiongNuDataArray() {
		check.PanicNotTrue(data.Level == uint64(idx+1), "%s 抗击匈奴难度只能从1开始配置，每个难度逐渐加1", filename)

		if prevLevel != nil {
			prevLevel.NextLevel = data
		}

		prevLevel = data
	}
}

func (data *ResistXiongNuData) Init(filename string, configs interface {
	GetResistXiongNuScoreDataArray() []*ResistXiongNuScoreData
}) {
	check.PanicNotTrue(len(configs.GetResistXiongNuScoreDataArray()) == len(data.ScorePlunderPrizes), "%s 配置的匈奴评分数量跟评分奖励数量不匹配!%d", filename, data.Level)
	check.PanicNotTrue(len(configs.GetResistXiongNuScoreDataArray()) == len(data.ScorePrestiges), "%s 配置的匈奴评分数量跟评分奖励声望数量不匹配!%d", filename, data.Level)
	if len(data.ShowPrizes) > 0 {
		check.PanicNotTrue(len(configs.GetResistXiongNuScoreDataArray()) == len(data.ShowPrizes), "%s 配置的匈奴评分数量跟展示奖励数量不匹配!%d", filename, data.Level)
	}
	check.PanicNotTrue(data.NpcBaseData.DestroyWhenLose, "%s 配置的匈奴npc 必须死亡摧毁主城!%d-%s", filename, data.Level, data.Name)

	check.PanicNotTrue(len(data.ResistWaves) > 0, "%s 配置的匈奴npc 阶段数据没有配置!%d-%s", filename, data.Level, data.Name)

	data.FirstResistWave = data.ResistWaves[0]
	data.maxWave = uint64(len(data.ResistWaves))

	data.ScorePrizes = make([]*resdata.Prize, 0)
	for _, pp := range data.ScorePlunderPrizes {
		data.ScorePrizes = append(data.ScorePrizes, pp.Prize)
	}

	var prev *ResistXiongNuWaveData
	for _, waveData := range data.ResistWaves {
		check.PanicNotTrue(waveData.Next == nil, "%s %d-%s 配置的波次中，竟然同一个波次被配置在多场中", filename, data.Level, data.Name)

		if prev != nil {
			prev.Next = waveData
		}

		prev = waveData

		data.TotalMonsterCount += uint64(len(waveData.Monsters))

		for _, m := range waveData.Monsters {
			data.MaxFightAmount = u64.Max(data.MaxFightAmount, m.CalculateFightAmount())
		}
	}

	var totalSoldier uint64
	for _, captain := range data.NpcBaseData.Npc.Captains {
		totalSoldier += captain.Soldier
	}
	for _, monster := range data.AssistMonsters {
		for _, captain := range monster.Captains {
			totalSoldier += captain.Soldier
		}
	}
	data.totalSoldier = totalSoldier
}

func (data *ResistXiongNuData) MaxWave() uint64 {
	return data.maxWave
}

func (data *ResistXiongNuData) GetTotalSoldier() uint64 {
	return data.totalSoldier
}

func (data *ResistXiongNuData) WaveData(wave uint64) *ResistXiongNuWaveData {
	if wave <= 0 || wave > uint64(len(data.ResistWaves)) {
		return nil
	}
	return data.ResistWaves[wave-1]
}

// 抗击匈奴波次
//gogen:config
type ResistXiongNuWaveData struct {
	_ struct{} `file:"抗击匈奴/攻城波次.txt"`

	Id       uint64                           // id
	Name     string                           // 名字
	Monsters []*monsterdata.MonsterMasterData // 攻城怪物

	Next *ResistXiongNuWaveData `head:"-"`
}

// 抗击匈奴评分
//gogen:config
type ResistXiongNuScoreData struct {
	_ struct{} `file:"抗击匈奴/评分.txt"`
	_ struct{} `proto:"shared_proto.ResistXiongNuScoreProto"`
	_ struct{} `protoconfig:"ResistXiongNuScore"`

	Level                     uint64                    // 评分等级
	Name                      string                    // 评分的名字，此处要支持多语言
	WipeOutInvadeMonsterCount uint64 `validator:"uint"` // 消灭入侵队伍数量
	UnlockNextLevel           bool                      // 该评分是否解锁下一难度
}

func (*ResistXiongNuScoreData) InitAll(filename string, configs interface {
	GetResistXiongNuScoreDataArray() []*ResistXiongNuScoreData
}) {
	var prev *ResistXiongNuScoreData
	var hasUnlockNextLevel bool

	for idx, data := range configs.GetResistXiongNuScoreDataArray() {
		check.PanicNotTrue(data.Level == uint64(idx+1), "%s 抗击匈奴评分 level 只能从1开始配置，每个评分逐渐加1", filename)

		if hasUnlockNextLevel {
			check.PanicNotTrue(data.UnlockNextLevel, "%s 抗击匈奴评分里面，上一个评分解锁了下一个难度，但是更高的评分却不解锁下一难度!", filename)
		}

		if data.UnlockNextLevel {
			hasUnlockNextLevel = true
		}

		if prev != nil {
			check.PanicNotTrue(prev.WipeOutInvadeMonsterCount < data.WipeOutInvadeMonsterCount, "%s 抗击匈奴评分里面，上个难度必须比下一个难度消灭的怪物数量少!%v,%v", filename, prev, data)
		}

		prev = data
	}

	check.PanicNotTrue(hasUnlockNextLevel, "%s 抗击匈奴评分里面，没有解锁下一难度的评分!", filename)
}
