package gm

import (
	"fmt"
	"github.com/lightpaw/logrus"
	"github.com/lightpaw/male7/config"
	"github.com/lightpaw/male7/config/kv"
	"github.com/lightpaw/male7/config/regdata"
	"github.com/lightpaw/male7/entity"
	"github.com/lightpaw/male7/entity/npcid"
	"github.com/lightpaw/male7/gen/iface"
	"github.com/lightpaw/male7/gen/pb/gm"
	"github.com/lightpaw/male7/pb/shared_proto"
	"github.com/lightpaw/male7/service"
	"github.com/lightpaw/male7/service/cluster"
	"github.com/lightpaw/male7/service/conflict/heroservice/herolock"
	"github.com/lightpaw/male7/service/heromodule"
	"github.com/lightpaw/male7/service/monitor/metrics"
	"github.com/lightpaw/male7/service/operate_type"
	"github.com/lightpaw/male7/util/idbytes"
	"github.com/lightpaw/male7/util/imath"
	"github.com/lightpaw/male7/util/msg"
	"github.com/lightpaw/male7/util/must"
	"github.com/lightpaw/male7/util/u64"
	"github.com/lightpaw/pbutil"
	"runtime/debug"
	"strconv"
	"strings"
)

func NewGmModule(dep iface.ServiceDep, db iface.DbService, config *kv.IndividualServerConfig, datas *config.ConfigDatas,
	modules iface.Modules, realmService iface.RealmService, reminderService iface.ReminderService, buffService iface.BuffService,
	pushService iface.PushService, farmService iface.FarmService, mingcWarService iface.MingcWarService, mingcService iface.MingcService,
	clusterService *cluster.ClusterService, seasonService iface.SeasonService, gameExporter iface.GameExporter, country iface.CountryService,
	tick iface.TickerService) *GmModule {

	m := &GmModule{
		dep:                 dep,
		time:                dep.Time(),
		db:                  db,
		tick:                tick,
		config:              config,
		datas:               datas,
		modules:             modules,
		heroDataService:     dep.HeroData(),
		world:               dep.World(),
		reminderService:     reminderService,
		realmService:        realmService,
		heroSnapshotService: dep.HeroSnapshot(),
		sharedGuildService:  dep.Guild(),
		pushService:         pushService,
		farmService:         farmService,
		mingcWarService:     mingcWarService,
		mingcService:        mingcService,
		clusterService:      clusterService,
		seasonService:       seasonService,
		buffService:         buffService,
		country:             country,
		gameExporter:        gameExporter,
	}

	if m.config.IsDebug {
		if m.config.IsDebugYuanbao {
			m.groups = []*gm_group{
				{
					tab: "常用",
					handler: []*gm_handler{
						newCmdIntHandler("加元宝(负数表示减)_10", "加元宝(负数表示减)", "100000", func(amount int64, hc iface.HeroController) {
							hc.FuncWithSend(func(hero *entity.Hero, result herolock.LockResult) {
								m.addYuanbao(amount, hero, result, hc)
								result.Changed()
								return
							})
						}),
					},
				},
			}
		} else {
			m.groups = []*gm_group{
				m.newCommonGmGroup(),
				m.newDomesticGmGroup(),
				m.newGoodsGmGroup(),
				m.newSingleGoodsGmGroup(),
				m.newSingleEquipmentGmGroup(),
				m.newSingleGemGmGroup(),
				m.newResetGmGroup(),
				m.newTaskGmGroup(),
				m.newDungeonGmGroup(),
				m.newMailGmGroup(),
				m.newLevelGmGroup(),
				m.newSceneGmGroup(),
				m.newZhanJiangGmGroup(),
				m.newMiscGmGroup(),
				m.newPrintEquipsGmGroup(),
				m.newMingcWarGmGroup(),
				m.newMingcGmGroup(),
				m.newRedPacketGmGroup(),
				m.newCountryGmGroup(),
			}
		}

	}

	// 处理模块消息
	var i interface{}
	i = m
	if modules, ok := i.(interface {
		initModuleHandler() []*gm_group
	}); ok {
		m.groups = append(m.groups, modules.initModuleHandler()...)
	}

	var protoBytes [][]byte
	for _, g := range m.groups {
		proto := &shared_proto.GmCmdListProto{}
		proto.Tab = g.tab

		for _, h := range g.handler {
			hp := &shared_proto.GmCmdProto{}
			hp.Cmd = h.cmd
			hp.Desc = h.desc
			hp.HasInput = len(h.defaultInput) > 0
			hp.DefaultInput = h.defaultInput

			proto.Cmd = append(proto.Cmd, hp)
		}

		protoBytes = append(protoBytes, must.Marshal(proto))
	}
	m.listCmdMsg = gm.NewS2cListCmdMarshalMsg(protoBytes).Static()

	m.hctx = heromodule.NewContext(m.dep, operate_type.GMCmd)

	return m
}

//gogen:iface
type GmModule struct {
	dep  iface.ServiceDep
	time iface.TimeService
	db   iface.DbService
	tick iface.TickerService

	config *kv.IndividualServerConfig
	datas  *config.ConfigDatas

	modules         iface.Modules
	heroDataService iface.HeroDataService
	//sharedGuildService iface.SharedGuildService
	world           iface.WorldService
	country         iface.CountryService
	reminderService iface.ReminderService

	realmService iface.RealmService

	heroSnapshotService iface.HeroSnapshotService

	sharedGuildService iface.GuildService

	pushService iface.PushService

	farmService iface.FarmService

	mingcWarService iface.MingcWarService

	mingcService iface.MingcService

	clusterService *cluster.ClusterService

	seasonService iface.SeasonService

	buffService iface.BuffService

	gameExporter iface.GameExporter

	groups []*gm_group

	listCmdMsg pbutil.Buffer

	hctx *heromodule.HeroContext
}

type gm_group struct {
	tab string

	handler []*gm_handler
}

func newIntHandler(desc, defaultInput string, f func(amount int64, hc iface.HeroController)) *gm_handler {
	return newCmdIntHandler(newCmd(desc), desc, defaultInput, f)
}

func newCmdIntHandler(cmd, desc, defaultInput string, f func(amount int64, hc iface.HeroController)) *gm_handler {
	return newCmdStringHandler(cmd, desc, defaultInput, func(input string, hc iface.HeroController) {
		i, err := strconv.ParseInt(input, 10, 64)
		if err != nil {
			logrus.WithError(err).Warnf("GM命令收到的input不是数字，cmd:%s input: %s", cmd, input)
		}

		f(i, hc)
	})
}

var cmdMap = map[string]struct{}{}

func newCmd(desc string) string {
	cmd := desc
	if _, exist := cmdMap[cmd]; exist {
		for i := 0; i < 1000; i++ {
			cmd = fmt.Sprintf("%s_%v", desc, i)
			if _, exist := cmdMap[cmd]; !exist {
				break
			}
		}
	}
	cmdMap[cmd] = struct{}{}
	return cmd
}

func newStringHandler(desc, defaultInput string, f func(input string, hc iface.HeroController)) *gm_handler {
	return newCmdStringHandler(newCmd(desc), desc, defaultInput, f)
}

func newHeroIntHandler(desc, defaultInput string, f func(amount int64, hero *entity.Hero, result herolock.LockResult, hc iface.HeroController)) *gm_handler {
	return newCmdIntHandler(newCmd(desc), desc, defaultInput, func(amount int64, hc iface.HeroController) {
		hc.FuncWithSend(func(hero *entity.Hero, result herolock.LockResult) {
			f(amount, hero, result, hc)
			result.Changed()
			return
		})
	})
}

func newHeroStringHandler(desc, defaultInput string, f func(input string, hero *entity.Hero, result herolock.LockResult, hc iface.HeroController)) *gm_handler {
	return newCmdStringHandler(newCmd(desc), desc, defaultInput, func(input string, hc iface.HeroController) {
		hc.FuncWithSend(func(hero *entity.Hero, result herolock.LockResult) {
			f(input, hero, result, hc)
			result.Changed()
			return
		})
	})
}

func newCmdStringHandler(cmd, desc, defaultInput string, f func(input string, hc iface.HeroController)) *gm_handler {
	h := &gm_handler{}
	h.cmd = cmd
	h.cmdSpace = cmd + " "
	h.desc = desc
	h.defaultInput = defaultInput
	h.handle = f

	return h
}

type gm_handler struct {
	cmd      string
	cmdSpace string

	desc string

	defaultInput string

	handle func(input string, hc iface.HeroController)
}

//gogen:iface c2s_list_cmd
func (m *GmModule) ProcessListCmdMsg(hc iface.HeroController) {
	hc.Send(m.listCmdMsg)
}

//gogen:iface
func (m *GmModule) ProcessGmMsg(proto *gm.C2SGmProto, hc iface.HeroController) {

	if !m.config.GetIsDebug() {
		logrus.Errorf("不是debug模式，但是收到debug消息")
		//hc.Disconnect()
		return
	}

	defer func() {
		if r := recover(); r != nil {
			// 严重错误. 英雄线程这里不能panic
			logrus.WithField("err", r).WithField("stack", string(debug.Stack())).Warn("GmMsg recovered from panic!!! SERIOUS PROBLEM")
			metrics.IncPanic()
		}
	}()

	logrus.Debugf("收到GM命令：%s", proto.Cmd)

	cmd := strings.TrimSpace(proto.Cmd)

	for _, g := range m.groups {
		for _, h := range g.handler {
			if strings.HasPrefix(cmd, h.cmdSpace) || cmd == h.cmd {
				input := ""
				if len(cmd) > len(h.cmdSpace) {
					input = cmd[len(h.cmdSpace):]
				}

				h.handle(input, hc)
				return
			}
		}
	}
	hc.Send(gm.NewS2cGmMsg("GM无效的命令: " + proto.Cmd))
}

//gogen:iface
func (m *GmModule) ProcessInvaseTargetIdMsg(proto *gm.C2SInvaseTargetIdProto, hc iface.HeroController) {

	var heroBaseX, heroBaseY int
	hc.Func(func(hero *entity.Hero, err error) (heroChanged bool) {
		heroBaseX, heroBaseY = hero.BaseX(), hero.BaseY()
		return false
	})

	mapData := m.realmService.GetBigMap().GetMapData()
	ux, uy := mapData.GetBlockByPos(heroBaseX, heroBaseY)

	startX := ux * mapData.BlockData().XLen
	startY := uy * mapData.BlockData().YLen

	sequence := regdata.BlockSequence(ux, uy)

	var data *regdata.RegionMultiLevelNpcData
	for _, data = range m.datas.GetRegionMultiLevelNpcDataArray() {
		if int32(data.TypeData.Type) == proto.NpcType {
			break
		}
	}

	id := npcid.GetNpcId(sequence, data.Id, npcid.NpcType_MultiLevelMonster)
	baseX := startX + data.OffsetBaseX
	baseY := startY + data.OffsetBaseY

	hc.Send(gm.NewS2cInvaseTargetIdMsg(idbytes.ToBytes(id), u64.Int32(baseX), u64.Int32(baseY)))
}

type hero_near_slice struct {
	baseX, baseY int
	a            []*entity.Hero
}

func (a *hero_near_slice) score(hero *entity.Hero) int {
	return imath.Abs(hero.BaseX()-a.baseX) + imath.Abs(hero.BaseY()-a.baseY)
}

func (a *hero_near_slice) Len() int           { return len(a.a) }
func (a *hero_near_slice) Swap(i, j int)      { a.a[i], a.a[j] = a.a[j], a.a[i] }
func (a *hero_near_slice) Less(i, j int) bool { return a.score(a.a[i]) < a.score(a.a[j]) }

func (m *GmModule) getOrCreateFakeHeroControler(id int64) iface.HeroController {
	sender := m.world.GetUserCloseSender(id)
	if sender != nil {
		u, ok := sender.(iface.ConnectedUser)
		if ok {
			return u.GetHeroController()
		}
	} else {
		sender = fakeSender
	}

	return service.NewHeroController(id, sender, "127.0.0.1", 0x100007f, 0, m.heroDataService.NewHeroLocker(id))
}

var fakeSender = &fake_sender{}

type fake_sender struct{}

func (m *fake_sender) Id() int64                        { return 0 }
func (m *fake_sender) SendAll(msgs []pbutil.Buffer)     {}
func (m *fake_sender) Send(msg pbutil.Buffer)           {}
func (m *fake_sender) SendIfFree(msg pbutil.Buffer)     {}
func (m *fake_sender) Disconnect(err msg.ErrMsg)        {}
func (m *fake_sender) DisconnectAndWait(err msg.ErrMsg) {}
func (m *fake_sender) IsClosed() bool                   { return false }

//func (module *GmModule) processGoodsCmd(args []string, hc iface.HeroController) bool {
//
//	module.goodsCmd("", hc)
//
//	return true
//}
//
//func (module *GmModule) goodsCmd(args string, hc iface.HeroController) {
//	for _, data := range module.datas.GoodsData().Array {
//		newCount := hero.Depot().AddGoods(data.Id, 100)
//		result.Add(depot.NewS2cUpdateGoodsMsg(u64.Int32(data.Id), u64.Int32(newCount)))
//	}
//}
//
//func (module *GmModule) processEquipmentCmd(args []string, hc iface.HeroController) bool {
//
//	module.equipmentCmd("", hc)
//
//	return true
//}
//
//func (module *GmModule) equipmentCmd(args string, hc iface.HeroController) {
//	for _, data := range module.datas.EquipmentData().Array {
//		e := entity.NewEquipment(hero.Depot().NewEquipmentId(), data)
//		hero.Depot().AddEquipment(e)
//		result.Add(equipment.NewS2cAddEquipmentMsg(must.Marshal(e.EncodeClient())))
//	}
//}
