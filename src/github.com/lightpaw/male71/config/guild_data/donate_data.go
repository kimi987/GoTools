package guild_data

import "github.com/lightpaw/male7/config/resdata"

func DonateId(seq, times uint64) uint64 {
	return seq*10000 + times
}

//gogen:config
type GuildDonateData struct {
	_  struct{} `file:"联盟/联盟捐献.txt"`
	_  struct{} `proto:"shared_proto.GuildDonateProto"`
	_  struct{} `protoconfig:"GuildDonate"`
	Id uint64   `head:"-,DonateId(%s.Sequence%c %s.Times)" protofield:"-"`

	Sequence uint64

	Times uint64

	// 捐献消耗
	Cost *resdata.Cost

	// 给联盟加多少建设值
	GuildBuildingAmount uint64

	// 给自己加多少贡献值
	ContributionAmount uint64

	// 给自己加多少捐献值
	DonationAmount uint64

	// 给自己加多少贡献币
	ContributionCoin uint64

	// 推荐官府等级，< 这个等级弹出警告
	RecommandGuanfuLevel uint64 `default:"1"`
}
