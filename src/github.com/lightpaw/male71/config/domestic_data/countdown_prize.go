package domestic_data

import (
	"github.com/lightpaw/male7/config/resdata"
	"time"
	"math/rand"
	"github.com/lightpaw/male7/util/check"
)

//gogen:config
type CountdownPrizeDescData struct {
	_  struct{} `file:"内政/倒计时礼包描述.txt"`
	_  struct{} `proto:"shared_proto.CountdownPrizeDescDataProto"`
	_  struct{} `protoconfig:"countdown_prize_desc"`
	Id uint64

	Desc string
}

// 倒计时礼包
//gogen:config
type CountdownPrizeData struct {
	_  struct{} `file:"内政/倒计时礼包.txt"`
	Id uint64

	Plunder *resdata.Plunder

	// 等待时间
	WaitDuration time.Duration

	Descs []*CountdownPrizeDescData

	firstData *CountdownPrizeData
	nextData  *CountdownPrizeData
}

func (d *CountdownPrizeData) Init(filename string, dataMap map[uint64]*CountdownPrizeData) {

	d.firstData = dataMap[1]
	check.PanicNotTrue(d.firstData != nil, "%s 配置的倒计时礼包id必须从1开始连续配置，ID[%d]不存在", filename, 1)

	if d.Id > 1 {
		prev := dataMap[d.Id-1]
		check.PanicNotTrue(prev != nil, "%s 配置的倒计时礼包id必须从1开始连续配置，ID[%d]不存在", filename, d.Id-1)

		prev.nextData = d
	}

}

func (d *CountdownPrizeData) FirstData() *CountdownPrizeData {
	return d.firstData
}

func (d *CountdownPrizeData) NextData() *CountdownPrizeData {
	return d.nextData
}

func (d *CountdownPrizeData) RandomDesc() *CountdownPrizeDescData {
	return d.Descs[rand.Intn(len(d.Descs))]
}
