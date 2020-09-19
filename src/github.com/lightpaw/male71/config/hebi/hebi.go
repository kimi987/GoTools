package hebi

import (
	"github.com/lightpaw/male7/config/resdata"
	"github.com/lightpaw/male7/pb/shared_proto"
	"time"
	"github.com/lightpaw/male7/config/goods"
	"github.com/lightpaw/male7/config/scene"
	"github.com/lightpaw/male7/util/check"
	"github.com/lightpaw/pbutil"
	"github.com/lightpaw/male7/gen/pb/hebi"
)

//gogen:config
type HebiMiscData struct {
	_ struct{} `singleton:"true"`
	_ struct{} `file:"合璧/合璧杂项.txt"`
	_ struct{} `proto:"shared_proto.HebiMiscProto"`
	_ struct{} `protoconfig:"hebi_misc"`

	RoomsMaxSize            uint64        `default:"100"`
	DailyRobCount           uint64        `default:"3"`
	RobCdDuration           time.Duration
	RobPosCdDuration        time.Duration `default:"10s"`
	RobProtectDuration      time.Duration
	HebiDuration            time.Duration
	HeShiBiType             shared_proto.HebiType
	RoomWaitExpiredDuration time.Duration `default:"3h" protofield:"-"`
	HebiHeroRecordMaxSize   uint64        `default:"50" protofield:"-"`

	// 战斗场景
	CombatScene *scene.CombatScene `protofield:",%s.Id"`

	CopySelfGoods *goods.GoodsData `protofield:",config.U64ToI32(%s.Id)"`

	HebiGoods []*goods.GoodsData `protofield:"-" head:"-"`
}

func (d *HebiMiscData) Init(filename string, configs interface {
	GetGoodsDataArray() []*goods.GoodsData
}) {
	d.HebiGoods = make([]*goods.GoodsData, 0)
	for _, g := range configs.GetGoodsDataArray() {
		if g.SpecType == shared_proto.GoodsSpecType_GAT_HEBI {
			d.HebiGoods = append(d.HebiGoods, g)
		}
	}
	check.PanicNotTrue(d.RoomWaitExpiredDuration > 0, "%v 合璧房间等待状态必须>0 %v", filename, d.RoomWaitExpiredDuration)
}

func GenHebiPrizeId(heroLevel uint64, hebiType shared_proto.HebiType, quality uint64) uint64 {
	return heroLevel<<16 | uint64(hebiType)<<8 | quality
}

//gogen:config
type HebiPrizeData struct {
	_ struct{} `file:"合璧/合璧奖励.txt"`
	_ struct{} `proto:"shared_proto.HebiPrizeProto"`
	_ struct{} `protoconfig:"hebi_prize"`

	Id           uint64                `head:"-,GenHebiPrizeId(%s.HeroLevel%c %s.HebiType%c %s.Quality%c)"`
	HeroLevel    uint64
	HebiType     shared_proto.HebiType
	Quality      uint64
	Prize        *resdata.Prize                         // 展示奖励
	AmountPrize  *resdata.Prize        `protofield:"-"` // 数值奖励
	PlunderPrize *resdata.PlunderPrize `protofield:"-"` // 物品奖励
	ShowPrizeMsg pbutil.Buffer `head:"-" protofield:"-"`
}

func (d *HebiPrizeData) Init(filename string) {
	if d.Prize == nil {
		d.ShowPrizeMsg = hebi.NewS2cViewShowPrizeMsg(&shared_proto.PrizeProto{}).Static()
	} else {
		d.ShowPrizeMsg = hebi.NewS2cViewShowPrizeMsg(d.Prize.PrizeProto()).Static()
	}
}
