package fight

import (
	"github.com/lightpaw/male7/config/kv"
	"github.com/lightpaw/male7/config/race"
	"github.com/lightpaw/male7/config/scene"
	"github.com/lightpaw/male7/gen/iface"
	"github.com/lightpaw/male7/pb/server_proto"
	"github.com/lightpaw/male7/pb/shared_proto"
	"github.com/lightpaw/male7/util/u64"
	"math/rand"
	"github.com/lightpaw/logrus"
	"github.com/lightpaw/jsoniter"
	"time"
	"github.com/google/uuid"
	"github.com/lightpaw/rpc7"
	"golang.org/x/net/context"
	"github.com/pkg/errors"
	"github.com/lightpaw/male7/entity"
	"github.com/lightpaw/male7combatx"
	"github.com/lightpaw/male7/service/operate_type"
	"runtime/debug"
	"github.com/lightpaw/male7/service/monitor/metrics"
	"github.com/lightpaw/male7/util/must"
	"github.com/lightpaw/male7/util"
)

const (
	CombatxGrpcTarget = "/m7/grpc/combatx"
)

func NewFightXService(datas iface.ConfigDatas, heroData iface.HeroDataService,
	serverConfig *kv.IndividualServerConfig, cluster iface.ClusterService, tlog iface.TlogService) *FightXService {

	configProto := datas.CombatConfig().Encode()

	configProtoBytes := must.Marshal(configProto)
	configProtoMd5 := util.Md5String(configProtoBytes)

	combatConfig, err := combatx.NewConfig(configProto)
	if err != nil {
		logrus.WithError(err).Panic("初始化战斗服配置失败")
	}

	m := &FightXService{}
	m.datas = datas
	m.individualServerConfig = serverConfig
	m.tlog = tlog
	m.heroData = heroData
	m.combatProtoMd5 = configProtoMd5
	m.combatProtoBytes = configProtoBytes
	m.combatConfig = combatConfig

	if serverConfig.EnableCosUploader {
		data, err := cluster.GetConfig(cosConfigPath)
		if err != nil {
			logrus.WithError(err).Panic("获取Cos配置失败")
		}

		appid := jsoniter.Get(data, "appid").ToString()
		secretid := jsoniter.Get(data, "secretid").ToString()
		secretKey := jsoniter.Get(data, "secretKey").ToString()
		region := jsoniter.Get(data, "region").ToString()
		bucketName := jsoniter.Get(data, "bucketName").ToString()

		m.uploader = combatx.NewCosUploader(appid, secretid, secretKey, region, bucketName,
			2*time.Second, combatx.CosPrefix, "temp", combatx.LocalPrefix)
	} else {
		m.uploader = combatx.NewLocalUploader("temp") // TODO
	}

	if cluster.EtcdClient() != nil {
		m.combatClient, err = rpc7.NewEtcdClient(cluster.EtcdClient(), CombatxGrpcTarget)
		if err != nil {
			logrus.WithError(err).Panic("初始化战斗服RPC（etcd）失败")
		}
	} else {
		if len(serverConfig.CombatClusterAddr) > 0 {
			m.combatClient, err = rpc7.NewClient(serverConfig.CombatClusterAddr)
			if err != nil {
				logrus.WithError(err).Panic("初始化战斗服RPC失败")
			}
		}
	}

	return m
}

//gogen:iface
type FightXService struct {
	datas                  iface.ConfigDatas
	individualServerConfig *kv.IndividualServerConfig
	tlog                   iface.TlogService
	heroData               iface.HeroDataService

	combatProtoMd5   string
	combatProtoBytes []byte
	combatConfig     *combatx.Config

	combatClient *rpc7.Client
	uploader     combatx.Uploader
}

func (m *FightXService) Close() {
	if m.combatClient != nil {
		m.combatClient.Close()
	}
}

func (m *FightXService) config() *race.RaceConfig {
	return m.datas.RaceConfig()
}

// 新版战斗
func (m *FightXService) SendFightRequest(tfctx *entity.TlogFightContext, combatScene *scene.CombatScene, attackerId, defenserId int64, attacker, defenser *shared_proto.CombatPlayerProto) *server_proto.CombatXResponseServerProto {
	return m.SendFightRequestReturnResult(tfctx, combatScene, attackerId, defenserId, attacker, defenser, false)
}

func (m *FightXService) SendFightRequestReturnResult(tfctx *entity.TlogFightContext, combatScene *scene.CombatScene, attackerId, defenserId int64, attacker, defenser *shared_proto.CombatPlayerProto, returnResult bool) (resp *server_proto.CombatXResponseServerProto) {

	config := m.datas.CombatConfig()

	proto := &server_proto.CombatXRequestServerProto{}
	proto.Debug = m.individualServerConfig.IsDebugFight

	proto.Seed = rand.Int63()
	proto.UploadFilePath = uuid.New().String()
	proto.ReturnResult = returnResult

	if defenser.WallStat != nil {
		proto.MapRes = combatScene.WallMapRes
	} else {
		proto.MapRes = combatScene.MapRes
	}
	proto.MapXLen = u64.Int32(config.CombatXLen)
	proto.MapYLen = u64.Int32(config.CombatYLen)

	proto.AttackerId = attackerId
	proto.DefenserId = defenserId

	proto.Attacker = attacker
	proto.Defenser = defenser

	// 设置race data
	raceDataMap := make(map[shared_proto.Race]*shared_proto.RaceDataProto)
	for _, t := range proto.Attacker.Troops {
		raceDataMap[t.Captain.Race] = m.config().GetProto(t.Captain.Race)
	}
	for _, t := range proto.Defenser.Troops {
		raceDataMap[t.Captain.Race] = m.config().GetProto(t.Captain.Race)
	}

	if m.combatClient != nil {
		proto.ConfigSum = m.combatProtoMd5

		//var resp *server_proto.CombatResponseServerProto
	retry:
	// TODO
	//ctxfunc.Timeout2s(func(ctx context.Context) (err error) {
	//	resp, err = remoteFightX(m.combatClient, ctx, proto)
	//	if err != nil {
	//		logrus.WithError(err).Warn("combat服务器返回错误")
	//	}
	//	return nil
	//})

		if resp != nil {
			if resp.ReturnCode == 0 {
				m.tlogFight(tfctx, proto, resp)
			}

			if resp.ReturnCode == ReturnCodeNeedConfig && len(proto.Config) <= 0 {
				// 没有配置
				if len(m.combatProtoBytes) > 0 {
					proto.Config = m.combatProtoBytes
					goto retry
				}
			}

			return resp
		}
	}

	resp = combatx.Handle(m.combatConfig, m.uploader, proto)
	if resp != nil && resp.ReturnCode == 0 {
		m.tlogFight(tfctx, proto, resp)
	}
	return resp
}

const (
	ReturnCodeNeedConfig = 5
)

func remoteFightX(c *rpc7.Client, ctx context.Context, proto *server_proto.CombatXRequestServerProto) (*server_proto.CombatXResponseServerProto, error) {
	if data, err := proto.Marshal(); err != nil {
		return nil, errors.Wrapf(err, "fightx proto marshal fail")
	} else {
		result, err := c.HandleBytes(ctx, "fightx", "", 0, data)
		if err != nil {
			return nil, errors.Wrapf(err, "fightx fail")
		}

		s2c := &server_proto.CombatXResponseServerProto{}
		if err := s2c.Unmarshal(result); err != nil {
			return nil, errors.Wrapf(err, "fightx s2c.Unmarshal() fail")
		}
		return s2c, nil
	}
}

func (m *FightXService) tlogFight(tfctx *entity.TlogFightContext, req *server_proto.CombatXRequestServerProto, resp *server_proto.CombatXResponseServerProto) {
	if m.tlog.DontGenTlog() {
		return
	}

	defer func() {
		if err := recover(); err != nil {
			logrus.WithField("stack", string(debug.Stack())).WithField("err", err).Errorf("recovered from FightService.%v panic. SEVERE!!!", "tlogFight")
			metrics.IncPanic()
		}
	}()

	var atkResult, defResult uint64
	if resp.AttackerWin {
		atkResult, defResult = 1, 2
	} else {
		atkResult, defResult = 2, 1
	}

	if heroId := req.AttackerId; heroId > 0 {
		var talen uint64
		var taids, taraces, tafight, tasoldiers [tlogTroopLen]uint64
		for i, t := range req.Attacker.Troops {
			if t != nil {
				idx := i
				if t.FightIndex > 0 {
					idx = int(t.FightIndex - 1)
				}

				if idx < len(taids) {
					taids[idx] = u64.FromInt32(t.Captain.Id)
					taraces[idx] = uint64(t.Captain.Race)
					tafight[idx] = u64.FromInt32(t.Captain.FightAmount)
					if len(resp.AttackerAliveSoldier) > 0 {
						tasoldiers[idx] = u64.FromInt32(resp.AttackerAliveSoldier[t.Captain.Id])
					}
				}

				talen++
			}
		}

		m.tlog.TlogRoundFlowById(heroId, tfctx.BattleType, tfctx.BattleID, operate_type.BattleTypeAtk, 0, atkResult, 0, u64.FromInt(len(req.Attacker.Troops)), talen, taids[0], taids[1], taids[2], taids[3], taids[4], taraces[0], tafight[0], taraces[1], tafight[1], taraces[2], tafight[2], taraces[3], tafight[3], taraces[4], tafight[4], tasoldiers[0], tasoldiers[1], tasoldiers[2], tasoldiers[3], tasoldiers[4], u64.FromInt32(resp.Score))
	}

	if heroId := req.DefenserId; heroId > 0 {
		var tdlen uint64
		var tdids, tdraces, tdfight, tdsoldiers [tlogTroopLen]uint64
		for i, t := range req.Defenser.Troops {
			if t != nil {
				idx := i
				if t.FightIndex > 0 {
					idx = int(t.FightIndex - 1)
				}

				if idx < len(tdids) {
					tdids[idx] = u64.FromInt32(t.Captain.Id)
					tdraces[idx] = uint64(t.Captain.Race)
					tdfight[idx] = u64.FromInt32(t.Captain.FightAmount)

					if len(resp.DefenserAliveSoldier) > 0 {
						tdsoldiers[idx] = u64.FromInt32(resp.DefenserAliveSoldier[t.Captain.Id])
					}
				}
				tdlen++
			}
		}

		m.tlog.TlogRoundFlowById(heroId, tfctx.BattleType, tfctx.BattleID, operate_type.BattleTypeDef, 0, defResult, 0, u64.FromInt(len(req.Defenser.Troops)), tdlen, tdids[0], tdids[1], tdids[2], tdids[3], tdids[4], tdraces[0], tdfight[0], tdraces[1], tdfight[1], tdraces[2], tdfight[2], tdraces[3], tdfight[3], tdraces[4], tdfight[4], tdsoldiers[0], tdsoldiers[1], tdsoldiers[2], tdsoldiers[3], tdsoldiers[4], u64.FromInt32(resp.Score))
	}
}
