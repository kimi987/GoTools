package data

import (
	"github.com/lightpaw/male7/config/i18n"
	"github.com/lightpaw/male7/pb/shared_proto"
	"time"
	"github.com/lightpaw/male7/util/u64"
	"github.com/lightpaw/male7/util/check"
	"github.com/lightpaw/male7/util/timeutil"
	"github.com/lightpaw/pbutil"
	"github.com/lightpaw/logrus"
	"github.com/lightpaw/male7/gen/pb/misc"
)

const (
	KeySelf        = "self"
	KeyNum         = "num"
	KeyEquip       = "equip"
	KeyNpc         = "npc"
	KeyCaptain     = "captain"
	KeyQuality     = "quality"
	KeyText        = "text"
	KeyGuild       = "guild"
	KeyName        = "name"
	KeyCountry     = "country"
	KeyScore       = "score"
	KeyMingc       = "mingc"
)

//gogen:config
type BroadcastData struct {
	_ struct{} `file:"文字/广播.txt"`

	Id                 string
	Sequence           uint64
	Text               *i18n.I18nRef
	Condition          []uint64        `validator:"uint"`
	BcType             shared_proto.BCType
	SendHourMinute     []uint64        `validator:"uint" default:"1200,1600,2000"`
	SendDuration       []time.Duration `head:"-"`
	SendChat           bool            `default:"false"`
	OnlySendOnce       bool            `default:"false"`
	TimingBroadcastMsg pbutil.Buffer   `head:"-"`
}

func (c *BroadcastData) InitAll(filename string, configs interface {
	GetBroadcastDataArray() []*BroadcastData
}) {
	seqMap := make(map[uint64]struct{})
	for _, d := range configs.GetBroadcastDataArray() {
		 if _, ok := seqMap[d.Sequence]; ok {
			 logrus.Panicf("%v sequence 不能重复 id:%v sequence:%v", filename, d.Id, d.Sequence)
		 } else {
		 	seqMap[d.Sequence] = struct{}{}
		 }
	}
}

func (c *BroadcastData) Init(filename string) {
	u64.Sort(c.SendHourMinute)
	for _, hhmm := range c.SendHourMinute {
		hh := hhmm / 100
		mm := hhmm % 100

		check.PanicNotTrue(hh < 24, "%s 系统定时广播小时数必须[0 <= hour < 24], hour: %v", filename, hh)
		check.PanicNotTrue(mm < 60, "%s 系统定时广播分钟数必须[0 <= minute < 60], minute: %v", filename, mm)

		dd := time.Duration(hh)*time.Hour + time.Duration(mm)*time.Minute
		c.SendDuration = append(c.SendDuration, dd)
	}

	if len(c.SendDuration) > 0 {
		check.PanicNotTrue(c.BcType == shared_proto.BCType_BC_TIMING_SYS, "%v 系统定时广播，bc_type 必须为 %v id:%v", filename, shared_proto.BCType_BC_TIMING_SYS, c.Id)
	}

	if c.BcType == shared_proto.BCType_BC_TIMING_SYS {
		c.TimingBroadcastMsg = misc.NewS2cSysTimingBroadcastMsg(c.Text.New().JsonString()).Static()
	}
}

func (d *BroadcastData) NewTextFields() *i18n.Fields {
	return d.Text.New()
}

func (c *BroadcastData) GetNextTime(ctime time.Time) time.Time {
	// 多个里面找一个时间最小的
	t := ctime.Add(24 * time.Hour)
	for _, d := range c.SendDuration {
		nextTime := getNextResetDailyTime(ctime, d)
		if nextTime.Before(t) {
			t = nextTime
		}
	}

	return t
}

// 每日0点重置，返回下一天的0点（传入今日0点，也返回下一天0点）
func getNextResetDailyTime(ctime time.Time, resetDuration time.Duration) time.Time {

	todayZeroTime := timeutil.DailyTime.PrevTime(ctime)
	resetTime := todayZeroTime.Add(resetDuration)

	if ctime.Before(resetTime) {
		return resetTime
	}

	return resetTime.Add(24 * time.Hour)
}

func (d *BroadcastData) CanBroadcast(num uint64) (cond uint64, succ bool) {
	// 完全相等才发广播，以后出现其他情况再特殊处理
	i := u64.GetIndex(d.Condition, num)
	if i >= 0 {
		cond = d.Condition[i]
		succ = true
		return
	}

	return
}
