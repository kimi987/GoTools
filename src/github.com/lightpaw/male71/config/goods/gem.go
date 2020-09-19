// 宝石
package goods

import (
	"github.com/lightpaw/male7/config/data"
	"github.com/lightpaw/male7/config/icon"
	"github.com/lightpaw/male7/config/race"
	"github.com/lightpaw/male7/pb/shared_proto"
	"github.com/lightpaw/male7/util/check"
	"sort"
)

//gogen:config
type GemDatas struct {
	_ struct{} `singleton:"true"`
	_ struct{} `file:"物品/宝石.txt"`

	levelDescArrayByGemTypeMap map[uint64][]*GemData `head:"-"` // 宝石根据宝石类型的二维数组，一维数组中的数据等级从大到小排序
	levelAscArrayByGemTypeMap  map[uint64][]*GemData `head:"-"` // 宝石根据宝石类型的二维数组，一维数组中的数据等级从大到小排序
}

func (d *GemDatas) Init(filename string, configDatas interface {
	GetGemDataArray() []*GemData
	GetRaceDataArray() []*race.RaceData
}) {
	d.levelAscArrayByGemTypeMap = make(map[uint64][]*GemData)
	d.levelDescArrayByGemTypeMap = make(map[uint64][]*GemData)

	for _, data := range configDatas.GetGemDataArray() {
		descArray := d.GetLevelDescArrayByGemType(data.GemType)
		descArray = append(descArray, data)
		d.levelDescArrayByGemTypeMap[data.GemType] = descArray
	}

	for gemType, descArray := range d.levelDescArrayByGemTypeMap {
		check.PanicNotTrue(len(descArray) != 0, "没有配置宝石类型为[%d]的宝石!", gemType)
		sort.Sort(gemLevelDescSlice(descArray))

		ascArray := make([]*GemData, len(descArray))
		for i := 0; i < len(descArray); i++ {
			ascArray[len(descArray)-1-i] = descArray[i]
		}
		d.levelAscArrayByGemTypeMap[gemType] = ascArray
	}
}

func (d *GemDatas) GetLevelDescArrayByGemType(gemType uint64) []*GemData {
	return d.levelDescArrayByGemTypeMap[gemType]
}

func (d *GemDatas) GetLevelAscArrayByGemType(gemType uint64) []*GemData {
	return d.levelAscArrayByGemTypeMap[gemType]
}

type gemLevelDescSlice []*GemData

func (a gemLevelDescSlice) Len() int           { return len(a) }
func (a gemLevelDescSlice) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a gemLevelDescSlice) Less(i, j int) bool { return a[i].Level > a[j].Level }

const startGemId = 11000000
const endGemId = 12000000

func isGemId(goodsId uint64) bool {
	return goodsId > startGemId && goodsId < endGemId
}

func GetGemId(gemType, level uint64) uint64 {
	return startGemId + gemType*1000 + level
}

func GetGemTypeLevelByGoodsId(goodsId uint64) (gemType, gemLevel uint64) {
	if isGemId(goodsId) {
		dataId := goodsId - startGemId
		gemType = dataId / 1000
		gemLevel = dataId % 1000
	}

	return
}

//func GetGemType(gemId uint64) uint64 {
//	return (gemId - startGemId) / 1000
//}
//
//func GetGemLevel(gemId uint64) uint64 {
//	return (gemId - startGemId) % 1000
//}

// 宝石
//gogen:config
type GemData struct {
	_ struct{} `file:"物品/宝石.txt"`
	_ struct{} `proto:"shared_proto.GemDataProto"`
	_ struct{} `protoconfig:"gem"`

	Id uint64 `validator:"int>0"`

	Name string // 名字

	Desc string // 描述

	Icon *icon.Icon `protofield:"IconId,%s.Id"` // 图标

	Quality shared_proto.Quality // 品质

	BaseStat *data.SpriteStat // 加的属性

	GemType uint64 `validator:"int>0"` // 宝石类型

	Level uint64 `validator:"int>0"` // 宝石等级，从1级开始

	YuanbaoPrice uint64 `validator:"uint" default:"0"` // 元宝价格

	UpgradeNeedCount uint64 `validator:"int>0"` // 升级到下一级需要的数量

	UpgradeToThisNeedFirstLevelCount uint64 `head:"-" protofield:"-"` // 升级到该级需要的宝石数量

	PrevLevel *GemData `head:"-" protofield:",config.U64ToI32(%s.Id)"` // 前置宝石，可能为空，为空的话，Level必须=1，否则的话前置等级宝石的 Level+1 必须等于 Level
	NextLevel *GemData `head:"-" protofield:",config.U64ToI32(%s.Id)"` // 下一级宝石，可能为空，为空的话，表示满级了
}

func (d *GemData) Init(filename string, dataMap map[uint64]*GemData) {

	gemId := GetGemId(d.GemType, d.Level)
	check.PanicNotTrue(d.Id == gemId, "%s 配置了无效的Id[%d]，应该等于 11000000 + gemType*1000 + level", filename, d.Id)
	check.PanicNotTrue(isGemId(d.Id), "%s 配置了无效的Id[%d]，应该等于 11000000 + gemType*1000 + level", filename, d.Id)

	if d.Level > 1 {
		prevLevel := dataMap[GetGemId(d.GemType, d.Level-1)]
		d.PrevLevel = prevLevel

		check.PanicNotTrue(d.PrevLevel != nil, "%s 配置的宝石 %d，等级为 %d，竟然没有前置宝石!", filename, d.Id, d.Level)
		check.PanicNotTrue(d.PrevLevel.GemType == d.GemType, "%s 配置的宝石 %d，等级为 %d，前置宝石的宝石类型跟自己不匹配!", filename, d.Id, d.Level)

		prevLevel.NextLevel = d
	}

	for _, g := range dataMap {
		if g != d {
			check.PanicNotTrue(d.GemType != g.GemType || d.Level != g.Level, "%s 配置了多个相同宝石类型[%d]，宝石等级[%d]的宝石!%d, %d", filename, d.GemType, d.Level, d.Id, g.Id)
		}
	}
}

func (d *GemData) InitAll(filename string, configDatas interface {
	GetGemDataArray() []*GemData
}) {
	for _, data := range configDatas.GetGemDataArray() {
		if data.Level != 1 {
			continue
		}

		data.UpgradeToThisNeedFirstLevelCount = 1

		for data.NextLevel != nil {
			data.NextLevel.UpgradeToThisNeedFirstLevelCount = data.UpgradeNeedCount * data.UpgradeToThisNeedFirstLevelCount

			check.PanicNotTrue(data.YuanbaoPrice*data.UpgradeNeedCount == data.NextLevel.YuanbaoPrice, "%s 配置的宝石[%d]元宝价格错误，必须是等于上一级价格*升级所需个数", filename, data.NextLevel.Id)

			data = data.NextLevel
		}
	}
}
