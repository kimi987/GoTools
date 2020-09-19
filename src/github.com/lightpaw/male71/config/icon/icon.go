package icon

import (
	"github.com/lightpaw/male7/util/check"
)

//gogen:config
type Icon struct {
	_ struct{} `file:"杂项/图标.txt"`
	_ struct{} `proto:"shared_proto.IconProto"`
	_ struct{} `protoconfig:"Icons"`

	Id           string                                                      // icon的id
	DefaultIcon  string `validator:"string>0" head:"icon" protofield:"Icon"` // 默认图像
	MiddleIcon   string `default:" "`                                        // 中图像
	BigIcon      string `default:" "`                                        // 大图像
	HeadIcon     string `default:" "`
	TailIcon     string `default:" "`
	SuperBigIcon string `default:" "`                                        // 大图像
	CaptainHead  bool   `default:"false" protofield:"-"`                     // 是否武将头像

	Tab  uint64 `validator:"int" default:"0"`   // 标签
	Text string `validator:"string" default:" "` // 图标文字
}

func (c *Icon) InitAll(filename string, configs interface {
	GetIconArray() []*Icon
}) {
	var headCount uint64
	for _, ic := range configs.GetIconArray() {
		if ic.CaptainHead {
			headCount++
		}
	}
	check.PanicNotTrue(headCount > 0, "%v 至少要配置%s个武将头像", filename, 1)
}
