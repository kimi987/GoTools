package maildata

import (
	"github.com/lightpaw/male7/pb/shared_proto"
	"github.com/lightpaw/male7/util/u64"
	"github.com/lightpaw/male7/config/i18n"
	"github.com/lightpaw/male7/config/resdata"
)

const(
	TagPvp = 0 // pvp战报
	TagPve = 1 // pve战报
	TagAct = 2 // 活动战报
)

//gogen:config
type MailData struct {
	_        struct{} `file:"文字/邮件.txt"`
	Id       string
	Icon     uint64   `validator:"uint"`
	Title    *i18n.I18nRef
	SubTitle *i18n.I18nRef
	Text     *i18n.I18nRef
	Keep     bool
	Tag      uint64 `validator:"int" default:"0" desc:"战报邮件分类"`

	Desc      *i18n.I18nRef
	Image     string
	ImageWord uint64 `validator:"uint" default:"0"`

	Prize      *resdata.Prize `default:"nullable"`
	prizeProto *shared_proto.PrizeProto
}

func (d *MailData) Init(filename string) {
	if d.Prize != nil {
		d.prizeProto = d.Prize.Encode4Init()
	}
}

func (d *MailData) NewTextMail(t shared_proto.MailType) *shared_proto.MailProto {
	proto := &shared_proto.MailProto{}
	proto.Icon = u64.Int32(d.Icon)
	proto.Title = d.Title.KeysOnlyJson()
	proto.SubTitle = d.SubTitle.KeysOnlyJson()
	proto.Text = d.Text.KeysOnlyJson()
	proto.Keep = d.Keep
	proto.ReportTag = u64.Int32(d.Tag)
	proto.Image = d.Image
	proto.ImageWord = u64.Int32(d.ImageWord)
	proto.MailType = t

	if t == shared_proto.MailType_MailNormal {
		// 普通邮件才有奖励
		if d.prizeProto != nil {
			proto.Prize = d.prizeProto
		}
	}
	return proto
}

func (d *MailData) NewTitleFields() *i18n.Fields {
	return d.Title.New()
}

func (d *MailData) NewSubTitleFields() *i18n.Fields {
	return d.SubTitle.New()
}

func (d *MailData) NewTextFields() *i18n.Fields {
	return d.Text.New()
}
