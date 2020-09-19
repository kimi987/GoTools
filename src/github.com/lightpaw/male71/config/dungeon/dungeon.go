package dungeon

import (
	"github.com/lightpaw/male7/config/monsterdata"
	"github.com/lightpaw/male7/config/resdata"
	"github.com/lightpaw/male7/config/scene"
	"github.com/lightpaw/male7/pb/shared_proto"
	"github.com/lightpaw/male7/util/check"
	"sort"
	"time"
	"github.com/lightpaw/male7/config/pvetroop"
	"github.com/lightpaw/male7/util/u64"
	"github.com/lightpaw/male7/config/captain"
)

const (
	StarConditionIsWin               = 1
	StarConditionGeHpPercent         = 2
	StarConditionLeTimeLimit         = 3
	StarConditionLeCaptainDeathCount = 4
)

// 副本章节配置
//gogen:config
type DungeonChapterData struct {
	_ struct{} `file:"推图副本/推图副本章节.txt"`
	_ struct{} `proto:"shared_proto.DungeonChapterProto"`
	_ struct{} `protoconfig:"DungeonChapter"`

	Id uint64 `validator:"int>0" key:"true"` // 章节编号（主键）
	//Sequence    uint64                   `validator:"int>0"`                                // 章节序列
	ChapterName string               `validator:"string>0"`                            // 章节名
	ChapterDesc string               `validator:"string>0"`                            // 章节描述
	Captain     *captain.CaptainData `protofield:"CaptainSoul,config.U64ToI32(%s.Id)"` // 将魂

	Type shared_proto.DifficultType `validator:"int"` // 难度类型

	BgImg        string         `validator:"string>0"` // 背景图
	DungeonDatas []*DungeonData `head:"-"`             // 副本
	PassPrize    *resdata.Prize                        // 通关奖励

	StarPrize []*resdata.Prize // 星数奖励
	Star      []uint64         // 星数要求

	FirstDungeon *DungeonData `head:"-" protofield:"-"` // 该章节该难度第一个副本
	LastDungeon  *DungeonData `head:"-" protofield:"-"` // 该章节该难度最后一个副本

	starDungeonIds []uint64
}

func (d *DungeonChapterData) Init(filename string, datas interface {
	GetDungeonDataArray() []*DungeonData
}) {
	//d.Type = shared_proto.DifficultType_ORDINARY

	for _, data := range datas.GetDungeonDataArray() {
		if data.ChapterId == d.Id && data.Type == d.Type {
			d.DungeonDatas = append(d.DungeonDatas, data)

			if data.Star > 0 {
				d.starDungeonIds = append(d.starDungeonIds, data.Id)
			}
		}
	}

	check.PanicNotTrue(len(d.DungeonDatas) > 0, "%s 中配置的 章节id: %d, 难度: %v, 副本数量为0!", filename, d.Id, d.Type)

	sort.Sort(dungeonDataSlice(d.DungeonDatas))
	d.FirstDungeon = d.DungeonDatas[0]
	d.LastDungeon = d.DungeonDatas[len(d.DungeonDatas)-1]

	for _, v := range d.DungeonDatas {
		v.chapterStarDungeonIds = d.starDungeonIds
	}

	//checkPuzzleArray := make([]bool, len(d.DungeonDatas))
	//
	//for idx, data := range d.DungeonDatas {
	//
	//	// 检查拼图
	//	check.PanicNotTrue(data.UnlockPuzzleIndex < uint64(len(d.DungeonDatas)),
	//		"%s 中配置的章节id: %d，难度: %v，解锁的拼图的index[%d]必须小于本章副本的数量[%d]!", filename, d.ChapterId, d.Type, data.UnlockPuzzleIndex, len(d.DungeonDatas))
	//	check.PanicNotTrue(!checkPuzzleArray[data.UnlockPuzzleIndex],
	//		"%s 中配置的章节id: %d，难度: %v，解锁的拼图的index[%d]出现了重复!", filename, d.ChapterId, d.Type, data.UnlockPuzzleIndex)
	//	checkPuzzleArray[data.UnlockPuzzleIndex] = true
	//
	//	if idx != 0 {
	//		check.PanicNotTrue(data.UnlockPassDungeon == d.DungeonDatas[idx-1],
	//			"%s 中配置的章节id: %d，难度: %v，前置副本必须是相同章节，相同难度的副本，且不为空!", filename, d.ChapterId, d.Type)
	//	}
	//}
}

func (d *DungeonChapterData) GetStarDungeonIds() []uint64 {
	return d.starDungeonIds
}

type dungeonDataSlice []*DungeonData

func (p dungeonDataSlice) Len() int           { return len(p) }
func (p dungeonDataSlice) Less(i, j int) bool { return p[i].Id < p[j].Id }
func (p dungeonDataSlice) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }

//gogen:config
type DungeonData struct {
	_ struct{} `file:"推图副本/推图副本.txt"`
	_ struct{} `proto:"shared_proto.DungeonDataProto"`

	Id uint64 `validator:"int>0"` // 副本id
	//Sequence            uint64       `validator:"int>0"`                                       // 副本序列
	Name                string                     `validator:"string>0"`                                       // 副本名
	Desc                string                     `validator:"string>0"`                                       // 副本描述
	UnlockHeroLevel     uint64                     `validator:"uint"`                                           // 解锁需要的君主等级，0表示跟君主等级无关
	UnlockPassDungeon   []*DungeonData             `protofield:",config.U64a2I32a(GetDungeonDataKeyArray(%s))"` // 解锁需要通关的副本(可能为空)，只有该副本完成了，才可以打
	UnlockBayeStage     uint64                     `validator:"uint" default:"0"`
	ChapterId           uint64                     `validator:"int>0"` // 章节Id
	Type                shared_proto.DifficultType `validator:"int"`   // 难度类型
	Star                uint64                     `validator:"uint" default:"0"`
	StarCondition       []uint64
	StarConditionValue  []uint64                   `validator:"uint,duplicate"`
	PassLimit           uint64                     `validator:"uint" default:"0"` // 每日通关次数限制
	Sp                  uint64                     `validator:"uint" default:"0"` // 通关所需体力值消耗
	UnlockPuzzleIndex   uint64                     `validator:"uint"`             // 解锁的拼图的位置，从0开始
	StoryId             uint64                     `validator:"uint"`             // 故事id
	DialogId            uint64                     `validator:"uint" default:"0"` // 剧情id
	PreBattleDialogId   uint64                     `validator:"uint" default:"0"` // 战前剧情id
	AfterBattleDialogId uint64                     `validator:"uint" default:"0"` // 战后剧情id
	BallonToolTip       string                     `default:"ballon_tool_tip"`    // 气泡提示
	NpcName             string                     `default:"npc_name"`           // 关卡npc名称
	NpcIcon             string                     `default:"npc_icon"`           // 关卡npc的icon
	FirstPassPrize      *resdata.Prize             `default:"nullable"`           // 首次通关奖励
	PassPrize           *resdata.Prize             `default:"nullable"`           // 通关奖励
	Plunder             *resdata.Plunder           `protofield:"-"`               // 掉落奖励
	ShowPrize           *resdata.Prize                                            // 掉落展示奖励
	Monster             *monsterdata.MonsterMasterData                            // 副本怪物
	CombatScene         *scene.CombatScene         `protofield:"-"`               // 战斗场景
	CombatSceneRes      string                     `head:"-"`                     // 战斗场景资源

	Prev *DungeonData `head:"-" protofield:"Prev,config.U64ToI32(%s.Id)"` // 上一关卡，可能为空
	Next *DungeonData `head:"-" protofield:"Next,config.U64ToI32(%s.Id)"` // 下一关卡，可能为空

	YuanJunData []*monsterdata.MonsterCaptainData `head:"yuan_jun_id"`

	// 剧情数据
	PlotIdx []uint64 `default:"nullable"`
	PlotId  []uint64 `default:"nullable"`

	GuideTroop *DungeonGuideTroopData `head:"-"` // 引导布阵

	chapterStarDungeonIds []uint64
}

func (d *DungeonData) Init(filename string, configs interface {
	GetMonsterCaptainData(key uint64) *monsterdata.MonsterCaptainData
	GetPveTroopData(key uint64) *pvetroop.PveTroopData
}) {
	troopCap := configs.GetPveTroopData(uint64(shared_proto.PveTroopType_DUNGEON)).Capacity
	check.PanicNotTrue(u64.FromInt(len(d.YuanJunData)) <= troopCap, "%v %v 的援军不能大于部队容量%v", filename, d.Name, troopCap)

	for i := 0; i < len(d.YuanJunData); i++ {
		idxi := d.YuanJunData[i].Index
		check.PanicNotTrue(idxi >= 1 && idxi <= troopCap, "%v %v 的援军index不能大于部队容量%v index:%v", filename, d.Name, troopCap, idxi)
		for j := 0; j < len(d.YuanJunData); j++ {
			idxj := d.YuanJunData[j].Index
			if i == j {
				continue
			}
			check.PanicNotTrue(d.YuanJunData[i] != d.YuanJunData[j] && idxi != idxj, "%v %v 的援军 id 必须和援军 index 都不能重复。 id:%v index:%v", filename, d.Name, d.YuanJunData[i].Id, idxi)
		}
	}

	if d.Monster.WallStat == nil {
		d.CombatSceneRes = d.CombatScene.MapRes
	} else {
		d.CombatSceneRes = d.CombatScene.WallMapRes
	}
}

func (*DungeonData) InitAll(filename string, configs interface {
	GetDungeonDataArray() []*DungeonData
}) {
	var prev *DungeonData
	for _, d := range configs.GetDungeonDataArray() {
		if d.Star > 0 {
			check.PanicNotTrue(u64.Int(d.Star) == len(d.StarCondition) && u64.Int(d.Star) == len(d.StarConditionValue), "%s 关卡 %d-%s 中获星条件数组长度与星星数量不匹配！星星数：%v 条件数组长度：%v 数值数组长度：%s", filename, d.Id, d.Name, d.Star, len(d.StarCondition), len(d.StarConditionValue))
		}
		check.PanicNotTrue(d.Prev == nil, "%s 关卡 %d-%s 被配置在多个章节中!", filename, d.Id, d.Name)
		d.Prev = prev
		if prev != nil {
			prev.Next = d
		}
		prev = d
	}
}

func (d *DungeonData) GetChapterStarDungeonIds() []uint64 {
	return d.chapterStarDungeonIds
}

// 根据传递的战斗信息来返回获取星星的激活列表
func (d *DungeonData) CalculateEnabledStars(resultHpPercent, passCostSeconds, captainDeathCount uint64) ([]bool, uint64) {
	if d.Star == 0 {
		return nil, 0
	}
	enabledStars := make([]bool, d.Star)
	var star uint64
	for i, condition := range d.StarCondition {
		switch condition {
		case StarConditionIsWin:
			enabledStars[i] = true
		case StarConditionGeHpPercent:
			enabledStars[i] = resultHpPercent >= d.StarConditionValue[i]
		case StarConditionLeTimeLimit:
			enabledStars[i] = passCostSeconds <= d.StarConditionValue[i]
		case StarConditionLeCaptainDeathCount:
			enabledStars[i] = captainDeathCount <= d.StarConditionValue[i]
		}

		if enabledStars[i] {
			star++
		}
	}
	return enabledStars, star
}

// 副本其他数据
//gogen:config
type DungeonMiscData struct {
	_ struct{} `singleton:"true"`
	_ struct{} `file:"推图副本/推图副本杂项.txt"`
	_ struct{} `proto:"shared_proto.DungeonMiscProto"`
	_ struct{} `protoconfig:"DungeonMisc"`

	MaxAutoTimes        uint64        `validator:"int>0" default:"100"`               // 最大扫荡次数
	RecoverAutoDuration time.Duration `default:"10m"`                                 // 恢复扫荡间隔
	DefaultAutoTimes    uint64        `validator:"int>0" default:"50" protofield:"-"` // 默认给的扫荡次数
	AutoPerTimes        uint64        `default:"5"`                                   // 客户端配置，每次扫荡多少次
}

// 副本引导布阵配置
//gogen:config
type DungeonGuideTroopData struct {
	_ struct{} `file:"推图副本/引导布阵.txt"`
	_ struct{} `proto:"shared_proto.DungeonGuideTroopDataProto"`

	Id       uint64   `validator:"int>0" protofield:"-"`      // 副本id
	NotFirst bool                                             // 是否第2次才教学引导
	Captain  []uint64                                         // 引导武将id
	SrcPos   []uint64 `validator:"int>0,duplicate,notAllNil"` // 引导初始阵位
	SrcPosX  []uint64 `validator:"int>0,duplicate,notAllNil"` // 引导初始阵位X
	DstPos   []uint64 `validator:"int>0,duplicate,notAllNil"` // 引导最终阵位
	DstPosX  []uint64 `validator:"int>0,duplicate,notAllNil"` // 引导最终阵位X
}

func (d *DungeonGuideTroopData) Init(filename string, configs interface {
	GetDungeonData(key uint64) *DungeonData
}) {
	dungeonData := configs.GetDungeonData(d.Id)
	check.PanicNotTrue(dungeonData != nil, "%v 引导布阵副本id %v 在副本表中没有配置", filename, d.Id)
	dungeonData.GuideTroop = d
}
