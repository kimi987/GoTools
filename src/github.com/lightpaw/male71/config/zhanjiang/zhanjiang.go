package zhanjiang

import (
	"github.com/lightpaw/male7/config/monsterdata"
	"github.com/lightpaw/male7/config/resdata"
	"github.com/lightpaw/male7/config/scene"
	"github.com/lightpaw/male7/gen/pb/zhanjiang"
	"github.com/lightpaw/male7/util/check"
	"github.com/lightpaw/male7/util/u64"
	"github.com/lightpaw/pbutil"
)

// 过关斩将章节
//gogen:config
type ZhanJiangChapterData struct {
	_ struct{} `file:"过关斩将/章节.txt"`
	_ struct{} `proto:"shared_proto.ZhanJiangChapterProto"`
	_ struct{} `protoconfig:"ZhanJiangChapter"`

	ChapterId   uint64 `validator:"int>0" key:"true"` // 章节Id，章节id从小到大就是挑战顺序
	ChapterName string `validator:"string>0"`         // 章节名
	ChapterDesc string `validator:"string>0"`         // 章节描述

	BgImg          string                  `validator:"string>0"`                 // 背景图
	ZhanJiangDatas []*ZhanJiangGuanQiaData `head:"guan_qia" protofield:"GuanQia"` // 关卡

	PreChapter *ZhanJiangChapterData `head:"-" protofield:"-"`
}

func (*ZhanJiangChapterData) InitAll(filename string, configs interface {
	GetZhanJiangChapterDataArray() []*ZhanJiangChapterData
}) {
	var prevGuanQia *ZhanJiangGuanQiaData
	var preChapter *ZhanJiangChapterData
	for idx, data := range configs.GetZhanJiangChapterDataArray() {
		check.PanicNotTrue(data.ChapterId == uint64(idx+1), "过关斩将章节必须从1开始逐行加1")

		if preChapter != nil {
			data.PreChapter = preChapter
		}

		preChapter = data

		for _, guanQia := range data.ZhanJiangDatas {
			check.PanicNotTrue(guanQia.Prev == nil, "%s 关卡 %d-%s 被配置在多个章节中!", filename, guanQia.Id, guanQia.Name)
			guanQia.Prev = prevGuanQia
			if prevGuanQia != nil {
				prevGuanQia.Next = guanQia
			}
			prevGuanQia = guanQia // 设置新的前置关卡
			guanQia.ChapterData = data
		}
	}
}

// 过关斩将关卡
//gogen:config
type ZhanJiangGuanQiaData struct {
	_ struct{} `file:"过关斩将/关卡.txt"`
	_ struct{} `proto:"shared_proto.ZhanJiangGuanQiaProto"`

	Id             uint64                `validator:"int>0"`                                 // 关卡id
	Name           string                `validator:"string>0"`                              // 关卡名字
	PositionDesc   string                `validator:"string>0"`                              // 关卡地点描述
	Desc           string                `validator:"string>0"`                              // 关卡描述
	BgImg          string                `validator:"string>0"`                              // 关卡背景图
	ZhanJiangDatas []*ZhanJiangData      `head:"guan" protofield:"Guan"`                     // 关卡中的所有的小关卡
	AbilityExp     uint64                `validator:"uint"`                                  // 成长值
	Prev           *ZhanJiangGuanQiaData `head:"-" protofield:"Prev,config.U64ToI32(%s.Id)"` // 上一关卡，可能为空
	Next           *ZhanJiangGuanQiaData `head:"-" protofield:"Next,config.U64ToI32(%s.Id)"` // 下一关卡，可能为空
	ShowGongXun    uint64                `validator:"uint"`                                  // 功勋，展示用的
	ShowPrize      *resdata.Prize        `default:"nullable"`                                // 通关奖励，展示用的

	ChapterData *ZhanJiangChapterData `head:"-" protofield:"-"` // 对应的章节

	PassMsg pbutil.Buffer `head:"-" protofield:"-"`
}

func (data *ZhanJiangGuanQiaData) Init() {
	data.PassMsg = zhanjiang.NewS2cPassMsg(u64.Int32(data.Id)).Static()
}

//gogen:config
type ZhanJiangData struct {
	_ struct{} `file:"过关斩将/小关卡.txt"`
	_ struct{} `proto:"shared_proto.ZhanJiangDataProto"`

	Id   uint64 `validator:"int>0"`    // id
	Name string `validator:"string>0"` // 名字
	Desc string `validator:"string>0"` // 描述
	Icon string `validator:"string>0"` // 图标

	PassPrize   *resdata.Prize     `default:"nullable" protofield:"-"` // 通关奖励
	Plunder     *resdata.Plunder   `default:"nullable" protofield:"-"` // 掉落奖励
	ShowPrize   *resdata.Prize     `default:"nullable"`                // 掉落展示奖励
	Monster     *monsterdata.MonsterMasterData                         // 副本怪物
	CombatScene *scene.CombatScene `protofield:"-"`                    // 战斗场景
	GongXun     uint64             `validator:"uint"`                  // 功勋
}

// 过关斩将其他数据
//gogen:config
type ZhanJiangMiscData struct {
	_ struct{} `singleton:"true"`
	_ struct{} `file:"过关斩将/其他.txt"`
	_ struct{} `proto:"shared_proto.ZhanJiangMiscDataProto"`
	_ struct{} `protoconfig:"ZhanJiangMisc"`

	DefaultTimes uint64 `validator:"int>0" default:"3" protofield:"MaxTimes"` // 默认次数
}
