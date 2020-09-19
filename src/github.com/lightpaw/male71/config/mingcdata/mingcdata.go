package mingcdata

import (
	"github.com/lightpaw/male7/config/icon"
	"github.com/lightpaw/male7/pb/shared_proto"
)

//gogen:config
type MingcBaseData struct {
	_ struct{} `file:"名城战/名城.txt"`
	_ struct{} `protogen:"true"`
	_ struct{} `protoimport:"mingc.proto"`
	_ struct{} `protoimport:"base.proto"`

	Id uint64 `desc:"名城id"`

	// 名字
	Name string `desc:"名城名称"`

	// 模型
	Model string `desc:"名城模型"`

	// 野外坐标
	BaseX uint64 `desc:"野外坐标"`
	BaseY uint64

	// 半径
	Radius uint64 `desc:"野外半径"`

	// 名城类型，都城，郡城，县城
	Type shared_proto.MincType

	// 周朝都城
	ZhouCaptain bool

	DefaultYinliang uint64 `desc:"初始银两"`

	DailyAddYinliang uint64 `desc:"仓库每日新增银两"`

	MaxYinliang uint64 `desc:"仓库银两上限"`

	HostDailyAddYinliang uint64 `desc:"占领盟每日收益"`

	// 名称初始所属国家
	Country uint64 `desc:"名称初始所属国家，洛阳配0"  validator:"uint"`

	WarIcon *icon.Icon `protofield:",%s.Id,string" desc:"名城图标"`

	AtkMinHufu uint64 `desc:"申请攻打名城最低虎符出价"`

	AtkMinGuildLevel uint64 `desc:"申请攻打名城最低联盟等级"`

	AstMaxGuild uint64 `desc:"申请协助的最大联盟数"`

	BaseMinDistance uint64 `desc:"建城最小距离"`
}
