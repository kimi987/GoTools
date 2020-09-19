package maildata

//import (
//	"github.com/lightpaw/male7/config/i18n"
//	"github.com/lightpaw/male7/pb/shared_proto"
//	"github.com/lightpaw/male7/util/check"
//	"github.com/lightpaw/male7/util/u64"
//	"unicode"
//)
//
////gogen:config
//type I18nMailData struct {
//	_    struct{} `file:"文字/国际化邮件.txt"`
//	Id   string
//	Icon uint64 `validator:"uint"`
//	Keep bool
//
//	titleKey    *i18n.Key
//	subTitleKey *i18n.Key
//	textKey     *i18n.Key
//}
//
//func (d *I18nMailData) Init() {
//
//	prefixRunes := []rune("mail")
//	for _, r := range []rune(d.Id) {
//		if unicode.IsUpper(r) {
//			prefixRunes = append(prefixRunes, '_', unicode.ToLower(r))
//		} else {
//			prefixRunes = append(prefixRunes, r)
//		}
//	}
//	prefix := string(prefixRunes)
//
//	d.titleKey = i18n.MailKey.Key(prefix + "_title")
//	d.subTitleKey = i18n.MailKey.Key(prefix + "_sub_title")
//	d.textKey = i18n.MailKey.Key(prefix + "_text")
//
//	check.PanicNotTrue(d.titleKey != nil, "i18n邮件[%s]的title没找到，key: %s（程序说可以加才能加）", d.Id, prefix+"_title")
//	check.PanicNotTrue(d.subTitleKey != nil, "i18n邮件[%s]的sub_title没找到，key: %s（程序说可以加才能加）", d.Id, prefix+"_sub_title")
//	check.PanicNotTrue(d.textKey != nil, "i18n邮件[%s]的text没找到，key: %s（程序说可以加才能加）", d.Id, prefix+"_text")
//
//}
//
//func (d *I18nMailData) NewTextMail() *shared_proto.MailProto {
//	proto := &shared_proto.MailProto{}
//	proto.Icon = u64.Int32(d.Icon)
//	proto.Title = d.titleKey.KeysOnlyJson()
//	proto.SubTitle = d.subTitleKey.KeysOnlyJson()
//	proto.Text = d.textKey.KeysOnlyJson()
//	proto.Keep = d.Keep
//
//	return proto
//}
//
//func (d *I18nMailData) NewTitleFields() i18n.Fields {
//	return i18n.NewKey(d.titleKey)
//}
//
//func (d *I18nMailData) NewSubTitleFields() i18n.Fields {
//	return i18n.NewKey(d.subTitleKey)
//}
//
//func (d *I18nMailData) NewTextFields() i18n.Fields {
//	return i18n.NewKey(d.textKey)
//}
