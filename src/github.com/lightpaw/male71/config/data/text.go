package data

import (
	"github.com/lightpaw/male7/util/check"
	"github.com/lightpaw/male7/config/i18n"
)

//gogen:config
type Text struct {
	_  struct{} `file:"文字/文本.txt"`
	Id string

	Text *i18n.I18nRef
}

func (t *Text) Init(filename string) {
	check.PanicNotTrue(len(t.Text.Value) > 0, "%s 中key=%s 的Text字段没有配置", filename, t.Id)
}

func (t *Text) New() *i18n.Fields {
	return t.Text.New()
}