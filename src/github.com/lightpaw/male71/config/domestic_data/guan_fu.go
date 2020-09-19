package domestic_data

import (
	"time"
	"github.com/lightpaw/male7/config/resdata"
)

//gogen:config
type GuanFuLevelData struct {
	_     struct{} `file:"内政/官府等级.txt"`
	_     struct{} `protoconfig:"guan_fu_level_data"`
	_     struct{} `proto:"shared_proto.GuanFuLevelProto"`
	Level uint64   `validator:"int>0"`

	// 繁荣度
	RestoreProsperity                        uint64        `default:"10"`  // 繁荣度恢复
	MoveBaseRestoreHomeProsperityDuration    time.Duration `default:"4h"`  // 迁城恢复繁荣度的buf的间隔
	MoveBaseRestoreHomeProsperity            uint64        `default:"100"` // 迁城buf期间每次恢复的繁荣度
	BuyProsperityRestoreDurationWith1Yuanbao time.Duration `default:"1m"`  // 1点券恢复多长时间
	BuyProsperityRestoreDurationWith1Cost    time.Duration `default:"1m"`  // 1点券恢复多长时间

	// 联盟工坊给的奖励
	WorkshopPrize *resdata.Prize
}
