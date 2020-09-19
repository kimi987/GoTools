package guild_data

import (
	"github.com/lightpaw/male7/config/i18n"
	"github.com/lightpaw/male7/pb/shared_proto"
	"github.com/lightpaw/male7/util/timeutil"
	"time"
)

//gogen:config
type GuildLogHelp struct {
	_ struct{} `singleton:"true"`
	_ struct{} `file:"联盟/联盟日志.txt"`

	JoinGuild          *GuildLogData
	LeaveGuild         *GuildLogData
	ReplyJoinGuild     *GuildLogData
	KickLeaveGuild     *GuildLogData
	UpdateMemberClass  *GuildLogData
	UpgradeTechnology  *GuildLogData
	PrestigePrize      *GuildLogData
	InvaseMonster      *GuildLogData
	CollectSalary      *GuildLogData
	CreateGuild        *GuildLogData
	NewLeaderImpeach   *GuildLogData
	NewLeaderDemise    *GuildLogData
	StartImpeach       *GuildLogData
	TerminateImpeach   *GuildLogData
	UpdateName         *GuildLogData
	UpdateFlagName     *GuildLogData
	UpdateCountry      *GuildLogData
	UpgradeLevel       *GuildLogData
	UpdateInternalText *GuildLogData

	FAttSucc    *GuildLogData // 进攻胜利
	FDefFail    *GuildLogData // 防守失败
	FAttDestroy *GuildLogData // 破坏流亡
	FDefDestroy *GuildLogData // 被破坏流亡

	StartXiongNu             *GuildLogData // 开启抗击匈奴
	WipeOutXiongNuTroop      *GuildLogData // 消灭匈奴部队
	UnlockXiongNu            *GuildLogData // 解锁匈奴难度
	ResistXiongNuAddPrestige *GuildLogData // 抗击匈奴获得联盟声望
	XiongNuBaseDestroy       *GuildLogData // 打爆匈奴主城

	YinliangMingcHost       *GuildLogData // 名城占领
	YinliangMcWarAtkWin     *GuildLogData // 名城进攻成功
	YinliangMcWarAtkFail    *GuildLogData // 名城进攻失败
	YinliangGuildReceive    *GuildLogData // 其他联盟送钱
	YinliangGuildSend       *GuildLogData // 送其他联盟钱
	YinliangGuildSendMember *GuildLogData // 送联盟成员
	YinliangGuildPaySalary  *GuildLogData // 发工资

	McWarApplyAtkFail *GuildLogData // 名城战申请攻打失败
	McWarApplyAtkSucc *GuildLogData // 名城战申请攻打成功
	McWarAtkWin       *GuildLogData // 名城战攻打成功
	McWarDefWin       *GuildLogData // 名城战防守成功
	HufuAdd *GuildLogData // 成员捐虎符
	BaowuAddPrestige *GuildLogData // 成员开启宝物增加联盟声望
	UpdateMark *GuildLogData // 联盟标记
	AssemblyWin *GuildLogData // 集结野战胜利
}

//gogen:config
type GuildLogData struct {
	_  struct{} `file:"联盟/联盟日志.txt"`
	Id string

	LogType shared_proto.GuildLogType

	Icon string `validator:"string"`

	Text *i18n.I18nRef

	SendChat bool `default:"false"`
}

func (d *GuildLogData) NewLogProto(t time.Time) *shared_proto.GuildLogProto {

	proto := &shared_proto.GuildLogProto{}
	proto.Icon = d.Icon
	proto.Text = d.Text.KeysOnlyJson()
	proto.Time = timeutil.Marshal32(t)
	proto.Type = d.LogType
	proto.DataId = d.Id

	return proto
}

func (d *GuildLogData) NewHeroLogProto(t time.Time, heroIdBytes []byte, head string) *shared_proto.GuildLogProto {

	proto := d.NewLogProto(t)
	if len(proto.Icon) <= 0 {
		proto.HeroId = heroIdBytes
		proto.Icon = head
	}

	return proto
}
