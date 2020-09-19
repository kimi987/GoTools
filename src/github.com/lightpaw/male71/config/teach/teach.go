package teach

import (
	"github.com/lightpaw/male7/config/dungeon"
	"github.com/lightpaw/male7/config/monsterdata"
	"github.com/lightpaw/male7/config/resdata"
	"github.com/lightpaw/male7/config/taskdata"
	"github.com/lightpaw/male7/util/check"
)

//gogen:config
type TeachChapterData struct {
	_ struct{} `file:"教学/关卡.txt"`
	_ struct{} `proto:"shared_proto.TeachChapterDataProto"`
	_ struct{} `protoconfig:"teach_data"`

	Id uint64

	MinHeroLevel      uint64
	PassBaYeTaskStage uint64 `validator:"uint"`
	PassDungeonId     uint64 `validator:"uint"`

	Title string
	Desc  string
	Image string

	Prize *resdata.Prize

	AtkStartMonster *monsterdata.MonsterMasterData
	AtkEndMonster   *monsterdata.MonsterMasterData
	DefMonster      *monsterdata.MonsterMasterData

	PrevData *TeachChapterData `protofield:"-" head:"-"`
}

func (data *TeachChapterData) InitAll(filename string, conf interface {
	GetTeachChapterDataArray() []*TeachChapterData
	GetBaYeStageData(uint64) *taskdata.BaYeStageData
	GetDungeonData(uint64) *dungeon.DungeonData
}) {

	var prevData *TeachChapterData
	for _, d := range conf.GetTeachChapterDataArray() {
		if d.PassBaYeTaskStage > 0 {
			check.PanicNotTrue(conf.GetBaYeStageData(d.PassBaYeTaskStage) != nil, "%v,id:%v PassBaYeTaskStage:%v 不存在", filename, d.Id, d.PassBaYeTaskStage)
		}

		if d.PassDungeonId > 0 {
			check.PanicNotTrue(conf.GetDungeonData(d.PassDungeonId) != nil, "%v,id:%v PassDungeonId:%v 不存在", filename, d.Id, d.PassDungeonId)
		}

		if prevData != nil {
			check.PanicNotTrue(prevData.Id+1 == d.Id, "%v,id:%v 必须依次递增。", filename, d.Id)
			check.PanicNotTrue(prevData.MinHeroLevel <= d.MinHeroLevel, "%v, id:%v, min_hero_level:%v 最小君主等级必须>=前一关的君主等级", filename, d.Id, d.MinHeroLevel)
			check.PanicNotTrue(prevData.PassBaYeTaskStage <= d.PassBaYeTaskStage, "%v, id:%v pass_ba_ye_task_stage:%v 通过霸业阶段必须>=前一关的通关霸业阶段", filename, d.Id, d.PassBaYeTaskStage)
		}

		d.PrevData = prevData
		prevData = d
	}
}
