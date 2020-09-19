package gardendata

import (
	"time"
)

//gogen:config
type GardenConfig struct {
	_ struct{} `singleton:"true"`
	_ struct{} `file:"花园/杂项.txt"`
	_ struct{} `proto:"shared_proto.GradonConfigProto"`
	_ struct{} `protoconfig:"GradonConfig"`

	TreasuryTreeFullTimes       uint64                                       // 浇水浇满次数
	TreasuryTreeCollectDuration time.Duration `default:"24h" protofield:"-"` // 浇水浇满后多久可以领奖
	TreasuryTreeHelpMeLogCount  uint64        `default:"5"`                  // 帮助日志最大个数
}
