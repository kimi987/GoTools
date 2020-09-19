package fight

import (
	"github.com/lightpaw/male7/config/kv"
	"github.com/lightpaw/male7/config/race"
	"github.com/lightpaw/male7/config/scene"
	"github.com/lightpaw/male7/gen/iface"
	"github.com/lightpaw/male7/pb/server_proto"
	"github.com/lightpaw/male7/pb/shared_proto"
	"github.com/lightpaw/male7/util/i32"
	"github.com/lightpaw/male7/util/u64"
	"github.com/lightpaw/male7combat/combat"
	"math/rand"
	"github.com/lightpaw/logrus"
	"github.com/lightpaw/jsoniter"
	"time"
	"github.com/google/uuid"
	"github.com/lightpaw/rpc7"
	"golang.org/x/net/context"
	"github.com/pkg/errors"
	"github.com/lightpaw/male7/util/ctxfunc"
	"github.com/lightpaw/male7/entity"
	"runtime/debug"
	"github.com/lightpaw/male7/service/monitor/metrics"
	"github.com/lightpaw/male7/service/operate_type"
	"github.com/lightpaw/male7/util/imath"
)

const (
	CombatGrpcTarget = "/m7/grpc/combat"
)

func NewFightService(datas iface.ConfigDatas, individualServerConfig iface.IndividualServerConfig, heroData iface.HeroDataService,
	serverConfig *kv.IndividualServerConfig, cluster iface.ClusterService, tlog iface.TlogService) *FightService {

	m := &FightService{}
	m.datas = datas
	m.individualServerConfig = individualServerConfig
	m.tlog = tlog
	m.heroData = heroData

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

		m.uploader = NewCosUploader(appid, secretid, secretKey, region, bucketName, 2*time.Second, "{{cos}}")
	}

	var err error
	if cluster.EtcdClient() != nil {
		m.combatClient, err = rpc7.NewEtcdClient(cluster.EtcdClient(), CombatGrpcTarget)
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
type FightService struct {
	datas                  iface.ConfigDatas
	individualServerConfig iface.IndividualServerConfig
	tlog                   iface.TlogService
	heroData               iface.HeroDataService

	combatClient *rpc7.Client
	uploader     combat.Uploader
}

func (m *FightService) Close() {
	if m.combatClient != nil {
		m.combatClient.Close()
	}
}

func (m *FightService) config() *race.RaceConfig {
	return m.datas.RaceConfig()
}

func (m *FightService) SendFightRequest(tfctx *entity.TlogFightContext, combatScene *scene.CombatScene, attackerId, defenserId int64, attacker, defenser *shared_proto.CombatPlayerProto) *server_proto.CombatResponseServerProto {
	return m.SendFightRequestReturnResult(tfctx, combatScene, attackerId, defenserId, attacker, defenser, false)
}

func (m *FightService) SendFightRequestReturnResult(tfctx *entity.TlogFightContext, combatScene *scene.CombatScene, attackerId, defenserId int64, attacker, defenser *shared_proto.CombatPlayerProto, returnResult bool) (resp *server_proto.CombatResponseServerProto) {

	config := m.datas.MilitaryConfig()

	proto := &server_proto.CombatRequestServerProto{}

	proto.Seed = rand.Int63()
	proto.UploadFilePath = uuid.New().String()
	proto.ReturnResult = returnResult

	if defenser.WallStat != nil {
		proto.MapRes = combatScene.WallMapRes
		proto.MinWallAttackRound = u64.Int32(m.datas.MilitaryConfig().MinWallAttackRound)
		proto.MaxWallAttachFixDamageRound = u64.Int32(m.datas.MilitaryConfig().MaxWallAttachFixDamageRound)
		proto.MaxWallBeenHurtPercent = i32.MultiF64(10000, config.MaxWallBeenHurtPercent)
	} else {
		proto.MapRes = combatScene.MapRes
	}
	proto.MapXLen = u64.Int32(config.CombatXLen)
	proto.MapYLen = u64.Int32(config.CombatYLen)
	proto.MaxRound = u64.Int32(config.CombatMaxRound)

	proto.Coef = i32.MultiF64(10000, config.CombatCoef)
	proto.CritRate = i32.MultiF64(10000, config.CombatCritRate)
	proto.RestraintRate = i32.MultiF64(10000, config.CombatRestraintRate)

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
	proto.Races = make([]*shared_proto.RaceDataProto, 0, len(raceDataMap))
	for _, r := range raceDataMap {
		proto.Races = append(proto.Races, r)
	}

	proto.ScorePercent = config.CombatScorePercent

	if m.combatClient != nil {
		//var resp *server_proto.CombatResponseServerProto
		ctxfunc.Timeout2s(func(ctx context.Context) (err error) {
			resp, err = remoteFight(m.combatClient, ctx, proto)
			if err != nil {
				logrus.WithError(err).Warn("combat服务器返回错误")
			}
			return nil
		})

		if resp != nil {
			if resp != nil && resp.ReturnCode == 0 {
				m.tlogFight(tfctx, proto, resp)
			}

			return resp
		}
	}

	if m.uploader == nil {
		resp = combat.LocalHandle(proto)
	} else {
		resp = combat.Handle(m.uploader, proto)
	}

	if resp != nil && resp.ReturnCode == 0 {
		m.tlogFight(tfctx, proto, resp)
	}
	return resp
}

func (m *FightService) SendMultiFightRequest(tfctx *entity.TlogFightContext, combatScene *scene.CombatScene, attackerId, defenserId []int64, attacker, defenser []*shared_proto.CombatPlayerProto,
	concurrentFightCount, attackerContinueWinCount, defenserContinueWinCount int32) *server_proto.MultiCombatResponseServerProto {
	return m.SendMultiFightRequestReturnResult(tfctx, combatScene, attackerId, defenserId, attacker, defenser, concurrentFightCount, attackerContinueWinCount, defenserContinueWinCount, false)
}

func (m *FightService) SendMultiFightRequestReturnResult(tfctx *entity.TlogFightContext, combatScene *scene.CombatScene, attackerId, defenserId []int64, attacker, defenser []*shared_proto.CombatPlayerProto,
	concurrentFightCount, attackerContinueWinCount, defenserContinueWinCount int32, returnResult bool) (resp *server_proto.MultiCombatResponseServerProto) {

	config := m.datas.MilitaryConfig()

	proto := &server_proto.MultiCombatRequestServerProto{}

	proto.Seed = rand.Int63()
	proto.UploadFilePath = uuid.New().String()
	proto.ReturnResult = returnResult

	proto.MapRes = combatScene.MapRes
	for _, def := range defenser {
		if def.WallStat != nil {
			proto.MapRes = combatScene.WallMapRes
			proto.MinWallAttackRound = u64.Int32(m.datas.MilitaryConfig().MinWallAttackRound)
			proto.MaxWallAttachFixDamageRound = u64.Int32(m.datas.MilitaryConfig().MaxWallAttachFixDamageRound)
			proto.MaxWallBeenHurtPercent = i32.MultiF64(10000, config.MaxWallBeenHurtPercent)
			break
		}
	}
	proto.MapXLen = u64.Int32(config.CombatXLen)
	proto.MapYLen = u64.Int32(config.CombatYLen)
	proto.MaxRound = u64.Int32(config.CombatMaxRound)

	proto.Coef = i32.MultiF64(10000, config.CombatCoef)
	proto.CritRate = i32.MultiF64(10000, config.CombatCritRate)
	proto.RestraintRate = i32.MultiF64(10000, config.CombatRestraintRate)

	proto.AttackerId = attackerId
	proto.DefenserId = defenserId

	proto.Attacker = attacker
	proto.Defenser = defenser

	proto.ConcurrentFightCount = concurrentFightCount
	proto.AttackerContinueWinCount = attackerContinueWinCount
	proto.DefenserContinueWinCount = defenserContinueWinCount

	// 设置race data
	raceDataMap := make(map[shared_proto.Race]*shared_proto.RaceDataProto)
	for _, attacker := range proto.Attacker {
		for _, t := range attacker.Troops {
			raceDataMap[t.Captain.Race] = m.config().GetProto(t.Captain.Race)
		}
	}
	for _, defenser := range proto.Defenser {
		for _, t := range defenser.Troops {
			raceDataMap[t.Captain.Race] = m.config().GetProto(t.Captain.Race)
		}
	}
	proto.Races = make([]*shared_proto.RaceDataProto, 0, len(raceDataMap))
	for _, r := range raceDataMap {
		proto.Races = append(proto.Races, r)
	}

	proto.ScorePercent = config.CombatScorePercent

	if m.combatClient != nil {
		//var resp *server_proto.MultiCombatResponseServerProto
		ctxfunc.Timeout2s(func(ctx context.Context) (err error) {
			resp, err = remoteMultiFight(m.combatClient, ctx, proto)
			if err != nil {
				logrus.WithError(err).Warn("combat服务器返回错误")
			}
			return nil
		})

		if resp != nil {
			if resp.ReturnCode == 0 {
				m.tlogMultiFight(tfctx, proto, resp)
			}
			return resp
		}
	}

	if m.uploader == nil {
		resp = combat.LocalHandleMulti(proto)
	} else {
		resp = combat.HandleMulti(m.uploader, proto)
	}

	if resp != nil && resp.ReturnCode == 0 {
		m.tlogMultiFight(tfctx, proto, resp)
	}

	return resp
}

func remoteFight(c *rpc7.Client, ctx context.Context, proto *server_proto.CombatRequestServerProto) (*server_proto.CombatResponseServerProto, error) {
	if data, err := proto.Marshal(); err != nil {
		return nil, errors.Wrapf(err, "fight proto marshal fail")
	} else {
		result, err := c.HandleBytes(ctx, "fight", "", 0, data)
		if err != nil {
			return nil, errors.Wrapf(err, "fight fail")
		}

		s2c := &server_proto.CombatResponseServerProto{}
		if err := s2c.Unmarshal(result); err != nil {
			return nil, errors.Wrapf(err, "fight s2c.Unmarshal() fail")
		}
		return s2c, nil
	}
}

func remoteMultiFight(c *rpc7.Client, ctx context.Context, proto *server_proto.MultiCombatRequestServerProto) (*server_proto.MultiCombatResponseServerProto, error) {
	if data, err := proto.Marshal(); err != nil {
		return nil, errors.Wrapf(err, "multi_fight proto marshal fail")
	} else {
		result, err := c.HandleBytes(ctx, "multi_fight", "", 0, data)
		if err != nil {
			return nil, errors.Wrapf(err, "multi_fight fail")
		}

		s2c := &server_proto.MultiCombatResponseServerProto{}
		if err := s2c.Unmarshal(result); err != nil {
			return nil, errors.Wrapf(err, "multi_fight s2c.Unmarshal() fail")
		}
		return s2c, nil
	}
}

const (
	tlogTroopLen = 5
)

func (m *FightService) tlogFight(tfctx *entity.TlogFightContext, req *server_proto.CombatRequestServerProto, resp *server_proto.CombatResponseServerProto) {
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

func (m *FightService) tlogMultiFight(tfctx *entity.TlogFightContext, req *server_proto.MultiCombatRequestServerProto, resp *server_proto.MultiCombatResponseServerProto) {
	if m.tlog.DontGenTlog() {
		return
	}

	defer func() {
		if err := recover(); err != nil {
			logrus.WithField("stack", string(debug.Stack())).WithField("err", err).Errorf("recovered from FightService.%v panic. SEVERE!!!", "TlogMultiFight")
			metrics.IncPanic()
		}
	}()

	var atkResult, defResult uint64
	if resp.AttackerWin {
		atkResult, defResult = 1, 2
	} else {
		atkResult, defResult = 2, 1
	}

	aliveSoldierMap := make(map[int64]*server_proto.AliveSoldierProto, len(resp.AliveSoldiers))
	for _, as := range resp.AliveSoldiers {
		aliveSoldierMap[as.Id] = as
	}

	n := imath.Min(len(req.Attacker), len(req.AttackerId))
	for i := 0; i < n; i++ {
		atker := req.Attacker[i]
		atkerId := req.AttackerId[i]

		if heroId := atkerId; heroId > 0 {
			as := aliveSoldierMap[atkerId]

			var talen uint64
			var taids, taraces, tafight, tasoldiers [tlogTroopLen]uint64
			for i, t := range atker.Troops {
				if t != nil {
					idx := i
					if t.FightIndex > 0 {
						idx = int(t.FightIndex - 1)
					}

					if idx < len(taids) {
						taids[idx] = u64.FromInt32(t.Captain.Id)
						taraces[idx] = uint64(t.Captain.Race)
						tafight[idx] = u64.FromInt32(t.Captain.FightAmount)
						if as != nil && len(as.AliveSoldier) > 0 {
							tasoldiers[idx] = u64.FromInt32(as.AliveSoldier[t.Captain.Id])
						}
					}

					talen++
				}
			}

			var teamType uint64
			if heroId == tfctx.LeaderId {
				teamType = operate_type.BattleTypeLeader
			} else {
				teamType = operate_type.BattleTypeMember
			}
			m.tlog.TlogRoundFlowById(heroId, tfctx.BattleType, tfctx.BattleID, operate_type.BattleTypeAtk, teamType, atkResult, 0, u64.FromInt(len(atker.Troops)), talen, taids[0], taids[1], taids[2], taids[3], taids[4], taraces[0], tafight[0], taraces[1], tafight[1], taraces[2], tafight[2], taraces[3], tafight[3], taraces[4], tafight[4], tasoldiers[0], tasoldiers[1], tasoldiers[2], tasoldiers[3], tasoldiers[4], u64.FromInt32(resp.Score))
		}
	}

	n = imath.Min(len(req.Defenser), len(req.DefenserId))
	for i := 0; i < n; i++ {
		defenser := req.Defenser[i]
		defenserId := req.DefenserId[i]

		if heroId := defenserId; heroId > 0 {
			as := aliveSoldierMap[defenserId]

			var tdlen uint64
			var tdids, tdraces, tdfight, tdsoldiers [tlogTroopLen]uint64
			for i, t := range defenser.Troops {
				if t != nil {
					idx := i
					if t.FightIndex > 0 {
						idx = int(t.FightIndex - 1)
					}

					if idx < len(tdids) {
						tdids[idx] = u64.FromInt32(t.Captain.Id)
						tdraces[idx] = uint64(t.Captain.Race)
						tdfight[idx] = u64.FromInt32(t.Captain.FightAmount)

						if as != nil && len(as.AliveSoldier) > 0 {
							tdsoldiers[idx] = u64.FromInt32(as.AliveSoldier[t.Captain.Id])
						}
					}

					tdlen++
				}
			}

			var teamType uint64
			if heroId == tfctx.LeaderId {
				teamType = operate_type.BattleTypeLeader
			} else {
				teamType = operate_type.BattleTypeMember
			}
			m.tlog.TlogRoundFlowById(heroId, tfctx.BattleType, tfctx.BattleID, operate_type.BattleTypeDef, teamType, defResult, 0, u64.FromInt(len(defenser.Troops)), tdlen, tdids[0], tdids[1], tdids[2], tdids[3], tdids[4], tdraces[0], tdfight[0], tdraces[1], tdfight[1], tdraces[2], tdfight[2], tdraces[3], tdfight[3], tdraces[4], tdfight[4], tdsoldiers[0], tdsoldiers[1], tdsoldiers[2], tdsoldiers[3], tdsoldiers[4], u64.FromInt32(resp.Score))
		}
	}

}
