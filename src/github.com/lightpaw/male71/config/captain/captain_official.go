package captain

import (
	"github.com/lightpaw/male7/util/check"
	"github.com/lightpaw/male7/config/data"
)

var (
	EmptyOfficialData = emptyOfficialData()
)

//gogen:config
type CaptainOfficialData struct {
	_ struct{} `file:"武将/武将官职.txt"`
	_ struct{} `protoconfig:"captain_official"`
	_ struct{} `proto:"shared_proto.CaptainOfficialProto"`

	Id           uint64 `validator:"int>0"`    // 官职id
	OfficialName string `validator:"string>0"` // 官职名称
	SpriteStat   *data.SpriteStat              // 加成
	NeedGongxun  uint64 `validator:"int>0"`    // 需要的功勋
	Icon         string
	Desc         string
}

func (d *CaptainOfficialData) InitAll(filename string, array []*CaptainOfficialData) {
	for i, d := range array {
		if i > 0 {
			prev := array[i-1]
			check.PanicNotTrue(prev.NeedGongxun <= d.NeedGongxun, "%s 每一级官职需要的功勋不能小于上一级的。id:%v need_gongxun:%v", filename, d.Id, d.NeedGongxun)
		}
	}
}

// 构造一个默认的初始官职
func emptyOfficialData() *CaptainOfficialData {
	stat := data.EmptyStat()
	return &CaptainOfficialData{Id: 0, OfficialName: "", SpriteStat: stat, NeedGongxun: 0}
}
