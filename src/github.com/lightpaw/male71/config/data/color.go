package data

import (
	"github.com/lightpaw/male7/pb/shared_proto"
	"github.com/lightpaw/male7/util/u64"
	"github.com/lightpaw/male7/util/check"
)

//gogen:config
type ColorData struct {
	_ struct{} `file:"文字/品质颜色.txt"`

	QualityKey             uint64               `key:"true" validator:"uint"`
	Quality                shared_proto.Quality `head:"-"`
	ColorCode              string
	ColorName              string
	QualityColorText       string               `head:"-"`
	CaptainSoulQualityName string               `default:"S"`
}

func (c *ColorData) Init(filename string) {
	_, ok := shared_proto.Quality_name[u64.Int32(c.QualityKey)]
	check.PanicNotTrue(ok, "%v 品质不存在 quality_key: %v", filename, c.QualityKey)
	c.Quality = shared_proto.Quality(u64.Int32(c.QualityKey))
	c.QualityColorText = c.GetColorText(c.ColorName)
}

func (c *ColorData) GetColorText(text string) string {
	return "[color=" + c.ColorCode + "]" + text + "[/color]"
}
