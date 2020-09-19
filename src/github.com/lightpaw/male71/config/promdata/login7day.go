package promdata

import "github.com/lightpaw/male7/config/resdata"

// 累积登陆奖励

//gogen:config
type LoginDayData struct {
	_ struct{} `file:"福利/7日登陆奖励.txt"`
	_ struct{} `protogen:"true"`

	Day uint64 `key:"true"`

	Prize *resdata.Prize
}

