package resdata

import (
	"github.com/lightpaw/config"
	"github.com/lightpaw/male7/pb/shared_proto"
	"github.com/lightpaw/male7/util/check"
	"strings"
)

//gogen:config
type AmountShowSortData struct {
	_ struct{} `file:"杂项/展示排序.txt"`
	_ struct{} `proto:"shared_proto.AmountShowSortProto"`
	_ struct{} `protoconfig:"AmountShowSortDatas"`

	Id       uint64                    // id
	Name     string                    `protofield:"-"` // 名字
	TypeList []shared_proto.AmountType `head:"-"`       // 数值类型
}

func (data *AmountShowSortData) Init(parser *config.ObjectParser, filename string) {
	array := parser.OriginStringArray("amount_types")
	for _, value := range array {
		if len(value) <= 0 {
			continue
		}

		i, ok := shared_proto.AmountType_value[strings.ToUpper(value)]
		check.PanicNotTrue(ok, "%s 配置了不存在的数值类型: %d-%s, %s", filename, data.Id, data.Name, value)

		if i == 0 {
			continue
		}

		data.TypeList = append(data.TypeList, shared_proto.AmountType(i))
	}
}
