package mingc

import (
	"github.com/lightpaw/logrus"
	"github.com/lightpaw/male7/entity"
	"github.com/lightpaw/male7/gen/iface"
	"github.com/lightpaw/male7/gen/pb/mingc"
	"github.com/lightpaw/male7/pb/shared_proto"
	"github.com/lightpaw/male7/service/conflict/heroservice/herolock"
	"github.com/lightpaw/male7/util/timeutil"
	"github.com/lightpaw/male7/util/u64"
	"time"
)

func NewMingcModule(dep iface.ServiceDep, baiZhan iface.BaiZhanService) *MingcModule {
	m := &MingcModule{}
	m.dep = dep
	m.baiZhan = baiZhan
	m.mingcSrv = dep.Mingc()

	return m
}

//gogen:iface
type MingcModule struct {
	dep      iface.ServiceDep
	baiZhan  iface.BaiZhanService
	mingcSrv iface.MingcService
}

//gogen:iface
func (m *MingcModule) ProcessMingcList(proto *mingc.C2SMingcListProto, hc iface.HeroController) {
	hc.Send(m.mingcSrv.MingcsMsg(u64.FromInt32(proto.Ver)))
}

//gogen:iface
func (m *MingcModule) ProcessViewMingc(proto *mingc.C2SViewMingcProto, hc iface.HeroController) {
	if mc := m.mingcSrv.Mingc(u64.FromInt32(proto.Id)); mc == nil {
		hc.Send(mingc.ERR_VIEW_MINGC_FAIL_INVALID_ID)
	} else {
		hc.Send(mingc.NewS2cViewMingcMsg(mc.Encode(m.dep.Datas().GetMingcBaseData(mc.Id()), m.dep.GuildSnapshot().GetGuildBasicProto)))
	}
}

//gogen:iface
func (m *MingcModule) ProcessMcBuildLog(proto *mingc.C2SMcBuildLogProto, hc iface.HeroController) {
	mcId := u64.FromInt32(proto.McId)

	mc := m.mingcSrv.Mingc(mcId)
	if mc == nil {
		hc.Send(mingc.ERR_MC_BUILD_LOG_FAIL_INVALID_MC_ID)
		return
	}

	msg := m.mingcSrv.McBuildLogMsg(mc)
	if msg == nil {
		logrus.Errorf("名城营建记录消息创建失败.mc:%v", mcId)
		hc.Send(mingc.NewS2cMcBuildLogMsg(nil))
		return
	}

	hc.Send(msg)
}

//gogen:iface
func (m *MingcModule) ProcessMcBuild(proto *mingc.C2SMcBuildProto, hc iface.HeroController) {
	mcId := u64.FromInt32(proto.McId)

	mc := m.mingcSrv.Mingc(mcId)
	if mc == nil {
		hc.Send(mingc.ERR_MC_BUILD_FAIL_INVALID_MC_ID)
		return
	}

	ctime := m.dep.Time().CurrentTime()

	var gid int64
	var nextTime time.Time
	if hc.FuncWithSend(func(hero *entity.Hero, result herolock.LockResult) {
		if hero.Level() < m.dep.Datas().McBuildMiscData().BuildMinHeroLevel {
			result.Add(mingc.ERR_MC_BUILD_FAIL_HERO_LEVEL_LIMIT)
			return
		}

		mcBuild := hero.McBuild()

		if mcBuild.BuildCount >= m.dep.Datas().McBuildMiscData().DailyBuildMaxCount {
			hc.Send(mingc.ERR_MC_BUILD_FAIL_NO_COUNT)
			return
		}

		if ctime.Before(mcBuild.NextTime) {
			result.Add(mingc.ERR_MC_BUILD_FAIL_IN_CD)
			return
		}

		if hero.GuildId() <= 0 {
			result.Add(mingc.ERR_MC_BUILD_FAIL_NO_GUILD)
			return
		}

		nextTime = ctime.Add(m.dep.Datas().McBuildMiscData().BuildCd)
		mcBuild.Build(nextTime)

		gid = hero.GuildId()
		result.Changed()
		result.Ok()
	}) {
		return
	}

	m.mingcSrv.Build(hc.Id(), gid, m.baiZhan.GetJunXianLevel(hc.Id()), mcId)
	hc.Send(mingc.NewS2cMcBuildMsg(proto.McId, u64.Int32(mc.Level()), u64.Int32(mc.Support()), u64.Int32(mc.DailyAddedSupport()), timeutil.Marshal32(nextTime)))
}

//gogen:iface
func (m *MingcModule) ProcessMingcHostGuild(proto *mingc.C2SMingcHostGuildProto, hc iface.HeroController) {
	mcId := u64.FromInt32(proto.McId)

	mc := m.mingcSrv.Mingc(mcId)
	if mc == nil {
		hc.Send(mingc.ERR_MC_BUILD_FAIL_INVALID_MC_ID)
		return
	}

	gid := mc.HostGuildId()

	var p *shared_proto.GuildSnapshotProto
	if g := m.dep.GuildSnapshot().GetSnapshot(gid); g != nil {
		p = g.Encode(m.dep.HeroSnapshot().GetBasicSnapshotProto)
	} else {
		p = &shared_proto.GuildSnapshotProto{}
	}

	hc.Send(mingc.NewS2cMingcHostGuildMsg(p))
}
