package regdata

import (
	"github.com/lightpaw/male7/config/basedata"
	"github.com/lightpaw/male7/util/check"
	"github.com/lightpaw/male7/entity/npcid"
)

//gogen:config
type BaozNpcData struct {
	_ struct{} `file:"地图/宝藏怪物.txt"`
	_ struct{} `protogen:"true"`

	Id uint64

	Npc *basedata.NpcBaseData `protofield:",config.U64ToI32(%s.Id),int32" desc:"Npc野怪数据 NpcBaseDataProto"`

	// 刷怪数量
	KeepCount uint64 `protofield:"-"`

	// 出征所需君主等级
	RequiredHeroLevel uint64 `desc:"出征所需君主等级"`

	// 稀有宝物id
	RareBaowuIds []uint64 `default:"nullable" protofield:"-"`
}

func (d *BaozNpcData) Init(filename string) {

	check.PanicNotTrue(d.Id <= npcid.NpcDataMask, "%s npc城池的配置数据的id最大不能超过 %d, id: %d", filename, npcid.NpcDataMask, d.Id)

	check.PanicNotTrue(d.KeepCount < 256, "%s 配置的刷怪数必须 < 256", filename)

}
