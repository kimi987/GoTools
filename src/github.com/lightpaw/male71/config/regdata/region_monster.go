package regdata

import "github.com/lightpaw/male7/config/basedata"

//gogen:config
type RegionMonsterData struct {
	_ struct{} `file:"地图/地区定点野怪.txt"`

	Id uint64

	// Npc列表
	Base *basedata.NpcBaseData

	// 坐标
	BaseX int
	BaseY int

	// 所属的Region Data Id
	RegionId uint64
}
