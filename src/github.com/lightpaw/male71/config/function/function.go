package function

import (
	"github.com/lightpaw/male7/config/domestic_data"
	"github.com/lightpaw/male7/config/herodata"
	"github.com/lightpaw/male7/config/taskdata"
	"github.com/lightpaw/male7/config/towerdata"
	"github.com/lightpaw/male7/gen/pb/misc"
	"github.com/lightpaw/male7/pb/shared_proto"
	"github.com/lightpaw/male7/util/check"
	"github.com/lightpaw/pbutil"
	"github.com/lightpaw/male7/config/dungeon"
)

// 功能开启数据

//gogen:config
type FunctionOpenData struct {
	_ struct{} `file:"功能开启/功能开启.txt"`
	_ struct{} `proto:"shared_proto.FunctionOpenDataProto"`
	_ struct{} `protoconfig:"FunctionOpenDatas"`

	//Id           uint64 `head:"-,uint64(%s.FunctionType)" protofield:"-"`
	//FunctionType shared_proto.FunctionType            // 功能类型
	FunctionType uint64 `key:"true"`
	Desc         string `default:" "`                  // 描述
	Icon         string `default:" "`                  // 图标
	NotifyOrder  uint64 `validator:"uint" default:"0"` // 提示排序

	GuanFuLevel *domestic_data.BuildingData `default:"nullable" protofield:",config.U64ToI32(%s.Level)"`           // 官府数据
	Building    *domestic_data.BuildingData `default:"nullable" protofield:"BuildingId,config.U64ToI32(%s.Level)"` // 官府数据
	HeroLevel   *herodata.HeroLevelData     `default:"nullable" protofield:",config.U64ToI32(%s.Level)"`           // 君主等级数据
	MainTask    *taskdata.MainTaskData      `default:"nullable" protofield:",config.U64ToI32(%s.Sequence)"`        // 主线任务数据
	BaYeStage   *taskdata.BaYeStageData     `default:"nullable" protofield:",config.U64ToI32(%s.Stage)"`           // 霸业阶段数据
	TowerFloor  *towerdata.TowerData        `default:"nullable" protofield:",config.U64ToI32(%s.Floor)"`           // 千重楼层级
	Dungeon     *dungeon.DungeonData        `default:"nullable" protofield:",config.U64ToI32(%s.Id)"`              // 幻境

	OpenMsg pbutil.Buffer `head:"-" protofield:"-"`
}

func (data *FunctionOpenData) Init(filepath string) {
	check.PanicNotTrue(data.FunctionType <= 255, "%s 配置的function_type不能超过255，如果要超过找后端接触封印")

	if data.GuanFuLevel != nil {
		check.PanicNotTrue(data.GuanFuLevel.Type == shared_proto.BuildingType_GUAN_FU, "%s 功能开启中的建筑类型必须是官府!%+v", filepath, data.GuanFuLevel.Type)
	}
	data.OpenMsg = misc.NewS2cOpenFunctionMsg(int32(data.FunctionType)).Static()
}
