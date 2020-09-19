package basedata

import (
	"github.com/lightpaw/male7/util/check"
	"github.com/lightpaw/male7/entity/npcid"
)

//gogen:config
type HomeNpcBaseData struct {
	_ struct{} `file:"地图/玩家主城野怪.txt"`

	Id uint64

	Data *NpcBaseData

	// 出现的偏移位置
	EvenOffsetX int `validator:"int"`
	EvenOffsetY int `validator:"int"`

	// 刷新时机，首次升级到X级
	HomeBaseLevel uint64 `validator:"int"`

	BaYeStage uint64 `validator:"int"`
}

func (b *HomeNpcBaseData) Init(filename string) {
	check.PanicNotTrue(b.Id <= npcid.NpcDataMask, "%s npc城池的配置数据的id最大不能超过 %d, id: %d", filename, npcid.NpcDataMask, b.Id)

	b.Data.DestroyWhenLose = true // TODO
	check.PanicNotTrue(b.Data.DestroyWhenLose, "%s npc城池[%v]的怪物必须是击破流亡的，DestroyWhenLose必须设置为1", filename, b.Id)

	baseLevelRefresh := b.HomeBaseLevel > 0
	baYeRefresh := b.BaYeStage > 0
	check.PanicNotTrue(baseLevelRefresh != baYeRefresh, "%s 配置的刷新时机，刷新主城等级和霸业阶段只能配置1个", filename)
}
