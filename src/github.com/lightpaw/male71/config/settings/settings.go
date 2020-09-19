package settings

import (
	"github.com/lightpaw/male7/pb/shared_proto"
	"github.com/lightpaw/male7/util/u64"
	"github.com/lightpaw/logrus"
)

var (
	DefaultPrivacySettings = []shared_proto.PrivacySettingType{}
)

// 推送设置
//gogen:config
type SettingMiscData struct {
	_ struct{} `singleton:"true"`
	_ struct{} `file:"设置/默认设置.txt"`

	DefaultSettings []shared_proto.SettingType
}

// 隐私设置
//gogen:config
type PrivacySettingData struct {
	_ struct{} `file:"设置/隐私设置.txt"`
	_ struct{} `protogen:"true"`

	Id           uint64  `protofield:"-"`
	SettingType  shared_proto.PrivacySettingType `head:"-"`
	Name         string
	NameType     string
	DefaultOpen  bool    `default:"false" protofield:"-"`
	RuleTitle    string
	RuleDesc     string
}

func (d *PrivacySettingData) Init(filename string) {
	d.SettingType = shared_proto.PrivacySettingType(u64.Int32(d.Id))
	switch d.SettingType {
	case shared_proto.PrivacySettingType_PST_SHOW_VIP_LEVEL,
	shared_proto.PrivacySettingType_PST_SHOW_VIP_HEAD_FRAME,
	shared_proto.PrivacySettingType_PST_SHOW_VIP_CHAT_BUBBLE,
	shared_proto.PrivacySettingType_PST_SHARE_GUILD_GIFT_GIVER:
		if d.DefaultOpen {
			DefaultPrivacySettings = append(DefaultPrivacySettings, d.SettingType)
		}
	default:
		logrus.Panicf("%s 新增的隐私设置id：%d 请通知程序员开发功能完毕后再配表", filename, d.Id)
	}
}
