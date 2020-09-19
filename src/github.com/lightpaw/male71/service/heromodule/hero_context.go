package heromodule

import (
	"github.com/lightpaw/logrus"
	"github.com/lightpaw/male7/config/data"
	"github.com/lightpaw/male7/config/i18n"
	"github.com/lightpaw/male7/entity"
	"github.com/lightpaw/male7/gen/iface"
	"github.com/lightpaw/male7/pb/shared_proto"
	"github.com/lightpaw/male7/service/conflict/guilddataservice/sharedguilddata"
	"github.com/lightpaw/male7/service/conflict/heroservice/herolock"
	"github.com/lightpaw/male7/service/herosnapshot/snapshotdata"
	"github.com/lightpaw/male7/service/monitor/metrics"
	"github.com/lightpaw/male7/service/operate_type"
	"github.com/lightpaw/male7/util/u64"
	"runtime/debug"
)

func NewContext(dep iface.ServiceDep, operateType *operate_type.OperateType) *HeroContext {
	s := &HeroContext{
		datas:       dep.Datas(),
		guild:       dep.GuildSnapshot(),
		broadcast:   dep.Broadcast(),
		world:       dep.World(),
		chat:        dep.Chat(),
		tlog:        dep.Tlog(),
		operateType: operateType,
	}
	return s
}

func NewContext2(
	datas iface.ConfigDatas,
	broadcast iface.BroadcastService,
	guild iface.GuildSnapshotService,
	world iface.WorldService,
	chat iface.ChatService,
	tlog iface.TlogService,
	operateType *operate_type.OperateType) *HeroContext {
	s := &HeroContext{
		datas:       datas,
		guild:       guild,
		broadcast:   broadcast,
		world:       world,
		chat:        chat,
		tlog:        tlog,
		operateType: operateType,
	}
	return s
}

func (hctx *HeroContext) Copy(t *operate_type.OperateType) *HeroContext {
	s := &HeroContext{
		datas:       hctx.datas,
		guild:       hctx.guild,
		broadcast:   hctx.broadcast,
		world:       hctx.world,
		chat:        hctx.chat,
		tlog:        hctx.tlog,
		operateType: t,
	}
	return s
}

type HeroContext struct {
	datas       iface.ConfigDatas
	broadcast   iface.BroadcastService
	guild       iface.GuildSnapshotService
	world       iface.WorldService
	chat        iface.ChatService
	tlog        iface.TlogService
	operateType *operate_type.OperateType

	// 提供给宝物日志使用
	baowuOp        shared_proto.BaowuOpType
	baowuOtherName string
	baowuOtherX    int32
	baowuOtherY    int32
	rareBaowuIds   []uint64
}

func (ctx *HeroContext) SetBaowuInfo(baowuOp shared_proto.BaowuOpType, baowuOtherName string,
	baowuOtherX int32, baowuOtherY int32, rareBaowuIds []uint64) {
	ctx.baowuOp = baowuOp
	ctx.baowuOtherName = baowuOtherName
	ctx.baowuOtherX = baowuOtherX
	ctx.baowuOtherY = baowuOtherY
	ctx.rareBaowuIds = rareBaowuIds
}

func (ctx *HeroContext) GetBaowuOp() shared_proto.BaowuOpType {
	return ctx.baowuOp
}

func (ctx *HeroContext) BaowuOtherInfo() (string, int32, int32) {
	return ctx.baowuOtherName, ctx.baowuOtherX, ctx.baowuOtherY
}

func (ctx *HeroContext) IsRareBaowu(baowuId uint64) bool {
	return u64.Contains(ctx.rareBaowuIds, baowuId)
}

func (ctx *HeroContext) Datas() iface.ConfigDatas {
	return ctx.datas
}

func (ctx *HeroContext) Broadcast() iface.BroadcastService {
	return ctx.broadcast
}

func (ctx *HeroContext) BroadcastHelp() *data.BroadcastHelp {
	return ctx.datas.BroadcastHelp()
}

func (ctx *HeroContext) Tlog() iface.TlogService {
	return ctx.tlog
}

func (ctx *HeroContext) World() iface.WorldService {
	return ctx.world
}

func (ctx *HeroContext) OperType() *operate_type.OperateType {
	return ctx.operateType
}

func (ctx *HeroContext) OperId() uint64 {
	return ctx.operateType.Id()
}

// **************** 系统广播 *********************

func (ctx *HeroContext) AddBroadcast(d *data.BroadcastData, hero *entity.Hero, result herolock.LockResult, subType, num uint64, f func() *i18n.Fields) {
	defer func() {
		if r := recover(); r != nil {
			logrus.WithField("stack", string(debug.Stack())).WithField("err", r).Errorf("系统广播异常！")
			metrics.IncPanic()
		}
	}()

	if d.OnlySendOnce {
		if hero.Misc().BroadcastSendTypeExist(d.Sequence, subType, num) {
			return
		}
	}
	_, succ := d.CanBroadcast(num)
	if !succ {
		return
	}

	text := f()
	args := text.JsonString()
	logrus.Debugf("========== 系统广播: %v, %v", d.Text, args)

	ctx.broadcast.Broadcast(args, d.SendChat)

	if d.OnlySendOnce {
		hero.Misc().AddBroadcastSendType(d.Sequence, subType, num)
		result.Changed()
	}

	result.Ok()
}

func (ctx *HeroContext) AddGuildBroadcast(d *data.BroadcastData, gid int64, subType, num uint64, f func() *i18n.Fields) {
	defer func() {
		if r := recover(); r != nil {
			logrus.WithField("stack", string(debug.Stack())).WithField("err", r).Errorf("系统广播异常！")
			metrics.IncPanic()
		}
	}()

	_, succ := d.CanBroadcast(num)
	if !succ {
		return
	}

	text := f()

	if text == nil {
		logrus.Debugf("系统广播 %v，text == nil。", d.Id)
		return
	}

	// 直接发
	ctx.broadcast.Broadcast(text.JsonString(), d.SendChat)
}

func (ctx *HeroContext) GetFlagHeroName(hero *entity.Hero) (self string) {
	var guildFlag string
	g := ctx.guild.GetSnapshot(hero.GuildId())
	if g != nil {
		guildFlag = g.FlagName
	}

	return ctx.datas.MiscConfig().FlagHeroName.FormatIgnoreEmpty(guildFlag, hero.Name())
}

func (ctx *HeroContext) GetFlagGuildNameFromSnapshot(gid int64) (self string, ok bool) {
	g := ctx.guild.GetSnapshot(gid)
	if g == nil {
		return
	}

	self = ctx.datas.MiscConfig().FlagHeroName.FormatIgnoreEmpty(g.FlagName, g.Name)
	ok = true
	return
}

func (ctx *HeroContext) GetFlagGuildName(g *sharedguilddata.Guild) (self string) {
	return ctx.datas.MiscConfig().FlagHeroName.FormatIgnoreEmpty(g.FlagName(), g.Name())
}

func (ctx *HeroContext) GetFlagName(flag string, name string) (self string) {
	return ctx.datas.MiscConfig().FlagHeroName.FormatIgnoreEmpty(flag, name)
}

func GetFlagHeroName(datas iface.ConfigDatas, hero *snapshotdata.HeroSnapshot) string {
	return datas.MiscConfig().FlagHeroName.FormatIgnoreEmpty(hero.GuildFlagName(), hero.Name)
}
