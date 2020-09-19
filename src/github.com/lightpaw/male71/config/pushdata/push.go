package pushdata

import (
	"github.com/lightpaw/male7/pb/shared_proto"
	"strings"
)

//gogen:config
type PushData struct {
	_ struct{} `file:"杂项/推送.txt"`
	_ struct{} `protogen:"true"`
	_ struct{} `protoimport:"setting.proto"`

	Id uint64 `head:"-,uint64(%s.Type)"`

	Type shared_proto.SettingType `desc:"推送类型"`

	Title string `desc:"标题"`

	Content string `desc:"正文"`

	TickTime string `validator:"string" desc:"触发时间"`
}

type PushFunc func(d *PushData) (title, content string)

func (d *PushData) ReplaceContent(replacer ...string) string {
	return replaceContent(d.Content, replacer...)
}

func replaceContent(content string, replacer ...string) string {

	n := len(replacer) / 2
	if n <= 0 {
		return content
	}

	for i := 0; i < n; i++ {
		idx := i * 2
		content = strings.Replace(content, replacer[idx], replacer[idx+1], -1)
	}
	return content
}
