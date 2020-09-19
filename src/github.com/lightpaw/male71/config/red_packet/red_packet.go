package red_packet

import (
	"github.com/lightpaw/male7/config/icon"
	"github.com/lightpaw/male7/config/resdata"
	"github.com/lightpaw/male7/pb/shared_proto"
	"github.com/lightpaw/male7/util/check"
	"time"
)

//gogen:config
type RedPacketData struct {
	_ struct{} `file:"杂项/红包.txt"`
	_ struct{} `protogen:"true"`

	Id              uint64
	Icon            *icon.Icon              `desc:"红包图标"`
	Name            string                  `desc:"红包名字"`
	Desc            string                  `desc:"红包描述"`
	DefaultText     string                  `desc:"红包默认文字"`
	Cost            *resdata.Cost           `desc:"红包价格，现在一定是元宝"`
	AmountType      shared_proto.AmountType `desc:"红包里包的数值类型"`
	Money           uint64                  `desc:"总钱数"`
	AllGarbbedPrize *resdata.Prize          `desc:"红包都抢完送的奖励"`
	ExpiredDuration time.Duration           `desc:"有效期" default:"24h"`
	MinPartMoney    uint64                  `desc:"抢红包保底金额" default:"1" protofield:"-"`
	MinCount        uint64                  `desc:"最小份数" default:"1"`
	MaxCount        uint64                  `desc:"最大份数" default:"100"`
	MaxTextLen      uint64                  `desc:"红包吉利话最大长度" default:"10"`
}

func (d *RedPacketData) Init(filename string) {
	check.PanicNotTrue(d.AmountType == shared_proto.AmountType_PT_YUANBAO || d.AmountType == shared_proto.AmountType_PT_DIANQUAN, "%v, amount_type:%v 必须配银两或点券。", filename, d.AmountType)
	check.PanicNotTrue(d.Money > 0, "%v, money: %v 必须>0", filename, d.Money)
	check.PanicNotTrue(d.MinPartMoney >0, "%v, min_part_money:%v 必须>0", filename, d.MinPartMoney)
	check.PanicNotTrue(d.MinCount > 0,  "%v, min_count:%v 必须>0", filename, d.MinCount)
	check.PanicNotTrue(d.MaxCount > 0,  "%v, max_count:%v 必须>0", filename, d.MaxCount)
	check.PanicNotTrue(d.MinCount <= d.MaxCount,  "%v, 必须 min_count:%v <= max_count:%v", filename, d.MinCount, d.MaxCount)
	check.PanicNotTrue(d.MaxTextLen > 0,  "%v, 必须 max_text_len:%v >0", filename, d.MaxTextLen)
}
