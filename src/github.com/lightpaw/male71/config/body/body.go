package body

import (
	"github.com/lightpaw/logrus"
	"github.com/lightpaw/male7/config/captain"
	"github.com/lightpaw/male7/config/herodata"
	"github.com/lightpaw/male7/config/icon"
	"github.com/lightpaw/male7/gen/pb/domestic"
	"github.com/lightpaw/male7/pb/shared_proto"
	"github.com/lightpaw/male7/util/check"
	"github.com/lightpaw/male7/util/u64"
	"github.com/lightpaw/pbutil"
)

// 头像数据
//gogen:config
type BodyData struct {
	_ struct{} `file:"杂项/形象.txt"`
	_ struct{} `proto:"shared_proto.BodyProto"`
	_ struct{} `protoconfig:"Bodys"`

	Id                  uint64                           `key:"1"`
	Icon                *icon.Icon                       `protofield:"IconId,%s.Id"`
	UnlockNeedCaptain   *captain.CaptainData             `default:"nil" protofield:"UnlockNeedCaptainSoul,config.U64ToI32(%s.Id)"` // 解锁需要的将魂id
	UnlockNeedHeroLevel *herodata.HeroLevelData          `default:"nil" protofield:",config.U64ToI32(%s.Level)"`                   // 解锁需要的英雄等级
	ChangeBodyMsg       pbutil.Buffer                    `head:"-" protofield:"-"`                                                 // 变更形象的协议
	DefaultBody         bool                             `default:"true" protofield:"-"`                                           // 是否是默认形象
	CountryOfficial     uint64                           `default:"0" validator:"uint" protofield:"-"`
	CountryOfficialType shared_proto.CountryOfficialType `head:"-" protofield:"-"`
}

func (data *BodyData) Init(filename string) {
	if data.DefaultBody {
		check.PanicNotTrue(data.UnlockNeedCaptain == nil, "%s 配置的默认头像 %d 里面必须不需要解锁将魂!", filename, data.Id)
		check.PanicNotTrue(data.UnlockNeedHeroLevel == nil, "%s 配置的默认头像 %d 里面必须不需要等级限制!", filename, data.Id)
	}

	data.ChangeBodyMsg = domestic.NewS2cChangeBodyMsg(u64.Int32(data.Id)).Static()

	if _, ok := shared_proto.CountryOfficialType_name[int32(data.CountryOfficial)]; !ok {
		logrus.Panicf("%v, id:%v 找不到官职类型:%v", filename, data.Id, data.CountryOfficial)
	}
	data.CountryOfficialType = shared_proto.CountryOfficialType(int32(data.CountryOfficial))
}
