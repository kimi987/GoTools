package mingc_war

import (
	"github.com/lightpaw/male7/gen/iface"
	"github.com/lightpaw/male7/util/u64"
	"github.com/lightpaw/male7/service/conflict/guilddataservice/sharedguilddata"
	"github.com/lightpaw/logrus"
	"github.com/lightpaw/pbutil"
	"github.com/lightpaw/male7/gen/pb/mingc_war"
	"github.com/lightpaw/male7/util/logp"
	"github.com/lightpaw/male7/entity/cb"
	"github.com/lightpaw/male7/pb/shared_proto"
	"github.com/lightpaw/male7/entity"
	"time"
	"github.com/lightpaw/male7/service/conflict/heroservice/herolock"
	"github.com/lightpaw/male7/util/timeutil"
	"github.com/lightpaw/male7/util/ctxfunc"
	"context"
	"github.com/lightpaw/male7/service/heromodule"
	"github.com/lightpaw/male7/service/operate_type"
	"github.com/lightpaw/male7/gen/pb/guild"
	"github.com/lightpaw/male7/util/i32"
)

func NewMingcWarModule(dep iface.ServiceDep, mingcSrv iface.MingcWarService, realm iface.RealmService) *MingcWarModule {
	m := &MingcWarModule{dep: dep}
	m.mcWarSrv = mingcSrv
	m.realm = realm

	return m
}

//gogen:iface
type MingcWarModule struct {
	dep      iface.ServiceDep
	mcWarSrv iface.MingcWarService
	realm    iface.RealmService
}

//gogen:iface
func (m *MingcWarModule) ProcessViewMcWar(proto *mingc_war.C2SViewMcWarProto, hc iface.HeroController) {
	hc.Send(m.mcWarSrv.ViewMsg(u64.FromInt32(proto.Ver)))
}

//gogen:iface c2s_view_mc_war_self_guild
func (m *MingcWarModule) ProcessViewMcWarSelfGuild(hc iface.HeroController) {
	if gid, succ := hc.LockGetGuildId(); succ {
		if g := m.dep.GuildSnapshot().GetSnapshot(gid); g != nil {
			hc.Send(mingc_war.NewS2cViewMcWarSelfGuildMsg(m.mcWarSrv.ViewSelfGuildProto(gid)))
			return
		} else {
			hc.Send(mingc_war.ERR_VIEW_MC_WAR_SELF_GUILD_FAIL_NO_GUILD)
		}
	} else {
		hc.Send(mingc_war.ERR_VIEW_MC_WAR_SELF_GUILD_FAIL_SERVER_ERR)
	}
}

//gogen:iface
func (m *MingcWarModule) ProcessViewMcWarScene(proto *mingc_war.C2SViewMcWarSceneProto, hc iface.HeroController) {
	succMsg, errMsg := m.mcWarSrv.ViewMcWarSceneMsg(u64.FromInt32(proto.McId))
	if errMsg != nil {
		hc.Send(errMsg)
		return
	}

	hc.Send(succMsg)
}

//gogen:iface
func (m *MingcWarModule) ProcessApplyAtk(proto *mingc_war.C2SApplyAtkProto, hc iface.HeroController) {
	mingcId := u64.FromInt32(proto.Mcid)
	mc := m.dep.Datas().GetMingcBaseData(mingcId)
	if mc == nil {
		hc.Send(mingc_war.ERR_APPLY_ATK_FAIL_INVALID_MCID)
		return
	}

	if m.dep.Datas().GetMingcWarSceneData(mc.Id) == nil {
		logrus.Warnf("MingcWarModule.ProcessApplyAtk mc:%v 没有开启名城战", mc.Id)
		hc.Send(mingc_war.ERR_APPLY_ATK_FAIL_INVALID_MCID)
		return
	}

	gid, ok := hc.LockGetGuildId()
	if !ok || gid <= 0 {
		logrus.Debugf("MingcWarModule hc.LockGetGuildId error gid:%v", gid)
		hc.Send(mingc_war.ERR_APPLY_ATK_FAIL_NOT_LEADER)
		return
	}

	if mc.Type == shared_proto.MincType_MC_DU {
		ctime := m.dep.Time().CurrentTime()
		serverStartTime := m.dep.SvrConf().GetServerStartTime()

		countryId := hc.LockHeroCountry()
		if mc.Country == countryId {
			if ctime.Before(serverStartTime.Add(m.dep.Datas().MingcMiscData().StartSelfCapitalAfterServerOpen)) {
				hc.Send(mingc_war.ERR_APPLY_ATK_FAIL_DU_CHENG_NOT_OPEN)
				return
			}
		} else {
			if ctime.Before(serverStartTime.Add(m.dep.Datas().MingcMiscData().StartOtherCapitalAfterServerOpen)) {
				hc.Send(mingc_war.ERR_APPLY_ATK_FAIL_OTHER_DU_CHENG_NOT_OPEN)
				return
			}

			if !m.dep.Mingc().IsHoldCountryCapital(gid, countryId) {
				hc.Send(mingc_war.ERR_APPLY_ATK_FAIL_NOT_HOLD_CAPITAL)
				return
			}

			if len(m.dep.Mingc().CountryHoldInitMcs(mc.Country)) != 1 {
				hc.Send(mingc_war.ERR_APPLY_ATK_FAIL_OTHER_COUNTRY_HAS_OTHER_MC)
				return
			}
		}
	}

	cost := u64.FromInt32(proto.Cost)

	var errMsg pbutil.Buffer
	m.dep.Guild().FuncGuild(gid, func(g *sharedguilddata.Guild) {
		if g == nil || g.LeaderId() != hc.Id() {
			errMsg = mingc_war.ERR_APPLY_ATK_FAIL_NOT_LEADER
			return
		}

		if g.LevelData().Level < mc.AtkMinGuildLevel {
			errMsg = mingc_war.ERR_APPLY_ATK_FAIL_GUILD_LEVEL_LIMIT
			return
		}

		_, succ := g.ReduceHufu(cost)
		if !succ {
			errMsg = mingc_war.ERR_APPLY_ATK_FAIL_HUFU_NOT_ENOUGH
			return
		}
	})

	if errMsg != nil {
		hc.Send(errMsg)
		return
	}

	succMsg, errMsg := m.mcWarSrv.ApplyAtk(gid, mc, cost)
	if succMsg != nil {
		hc.Send(succMsg)
		// 虎符广播刷新
		m.dep.Guild().ClearSelfGuildMsgCache(gid)
		m.dep.Guild().FuncGuild(gid, func(g *sharedguilddata.Guild) {
			if g == nil {
				return
			}
			if ids := g.AllUserMemberIds(); len(ids) > 0 {
				m.dep.World().MultiSend(ids, guild.NewS2cUpdateHufuMsg(u64.Int32(g.GetHufu())))
			}
		})
		return
	}

	hc.Send(errMsg)

	// 失败，还虎符
	m.dep.Guild().FuncGuild(gid, func(g *sharedguilddata.Guild) {
		if g == nil {
			return
		}

		g.AddHufu(cost)
	})
}

//gogen:iface
func (m *MingcWarModule) ProcessApplyAst(proto *mingc_war.C2SApplyAstProto, hc iface.HeroController) {
	mingcId := u64.FromInt32(proto.Mcid)
	mc := m.dep.Datas().GetMingcBaseData(mingcId)
	if mc == nil {
		hc.Send(mingc_war.ERR_APPLY_AST_FAIL_INVALID_MCID)
		return
	}

	gid, ok := hc.LockGetGuildId()
	if !ok || gid <= 0 {
		logrus.Debugf("MingcWarModule hc.LockGetGuildId error gid:%v", gid)
		hc.Send(mingc_war.ERR_APPLY_AST_FAIL_NOT_LEADER)
		return
	}

	g := m.dep.GuildSnapshot().GetSnapshot(gid)
	if g == nil || g.LeaderId != hc.Id() {
		hc.Send(mingc_war.ERR_APPLY_AST_FAIL_NOT_LEADER)
		return
	}

	succMsg, errMsg := m.mcWarSrv.ApplyAst(gid, proto.Atk, mc)
	if errMsg != nil {
		hc.Send(errMsg)
		return
	}

	hc.Send(succMsg)
}

//gogen:iface
func (m *MingcWarModule) ProcessCancelApplyAst(proto *mingc_war.C2SCancelApplyAstProto, hc iface.HeroController) {
	mingcId := u64.FromInt32(proto.Mcid)
	mc := m.dep.Datas().GetMingcBaseData(mingcId)
	if mc == nil {
		hc.Send(mingc_war.ERR_CANCEL_APPLY_AST_FAIL_INVALID_MCID)
		return
	}

	gid, ok := hc.LockGetGuildId()
	if !ok || gid <= 0 {
		logrus.Debugf("MingcWarModule hc.LockGetGuildId error gid:%v", gid)
		hc.Send(mingc_war.ERR_CANCEL_APPLY_AST_FAIL_NOT_LEADER)
		return
	}

	g := m.dep.GuildSnapshot().GetSnapshot(gid)
	if g == nil || g.LeaderId != hc.Id() {
		hc.Send(mingc_war.ERR_APPLY_AST_FAIL_NOT_LEADER)
		return
	}

	succMsg, errMsg := m.mcWarSrv.CancelApplyAst(gid, mc)
	if errMsg != nil {
		hc.Send(errMsg)
		return
	}

	hc.Send(succMsg)
}

//gogen:iface
func (m *MingcWarModule) ProcessReplyApplyAst(proto *mingc_war.C2SReplyApplyAstProto, hc iface.HeroController) {
	mingcId := u64.FromInt32(proto.Mcid)
	mc := m.dep.Datas().GetMingcBaseData(mingcId)
	if mc == nil {
		hc.Send(mingc_war.ERR_REPLY_APPLY_AST_FAIL_INVALID_MCID)
		return
	}

	gid, ok := hc.LockGetGuildId()
	if !ok || gid <= 0 {
		logp.Debugf("gid:%v", gid)
		hc.Send(mingc_war.ERR_REPLY_APPLY_AST_FAIL_NOT_LEADER)
		return
	}

	g := m.dep.GuildSnapshot().GetSnapshot(gid)
	if g == nil || g.LeaderId != hc.Id() {
		hc.Send(mingc_war.ERR_REPLY_APPLY_AST_FAIL_NOT_LEADER)
		return
	}

	succMsg, errMsg := m.mcWarSrv.ReplyApplyAst(gid, int64(proto.Gid), mc, proto.Agree)
	if errMsg != nil {
		hc.Send(errMsg)
		return
	}

	hc.Send(succMsg)
}

//gogen:iface
func (m *MingcWarModule) ProcessWatch(proto *mingc_war.C2SWatchProto, hc iface.HeroController) {
	succMsg, errMsg := m.mcWarSrv.Watch(hc, u64.FromInt32(proto.McId))
	if errMsg != nil {
		hc.Send(errMsg)
		return
	}
	hc.Send(succMsg)
}

//gogen:iface
func (m *MingcWarModule) ProcessQuitWatch(proto *mingc_war.C2SQuitWatchProto, hc iface.HeroController) {
	succMsg, errMsg := m.mcWarSrv.QuitWatch(hc, u64.FromInt32(proto.McId))
	if errMsg != nil {
		hc.Send(errMsg)
		return
	}
	hc.Send(succMsg)
}

//gogen:iface
func (m *MingcWarModule) ProcessJoinFight(proto *mingc_war.C2SJoinFightProto, hc iface.HeroController) {
	if hc.FuncWithSend(func(hero *entity.Hero, result herolock.LockResult) {
		if hero.Level() < m.dep.Datas().MingcMiscData().JoinFightHeroMinLevel {
			result.Add(mingc_war.ERR_JOIN_FIGHT_FAIL_HERO_LEVEL_LIMIT)
			return
		}
		result.Ok()
	}) {
		return
	}

	succMsg, errMsg := m.mcWarSrv.JoinFight(hc, u64.FromInt32(proto.Mcid), u64.FromInt32Array(proto.CaptainId), proto.XIndex)
	if errMsg != nil {
		hc.Send(errMsg)
		return
	}

	hc.Send(succMsg)

	// 免战
	m.setMian(hc)
}

func (m *MingcWarModule) setMian(hc iface.HeroController) {
	ctime := m.dep.Time().CurrentTime()
	start, end := m.mcWarSrv.McWarStartEndTime()
	if ctime.Before(start) || ctime.After(end) {
		return
	}

	var setMian bool
	var mianEndTime time.Time
	hc.FuncWithSend(func(hero *entity.Hero, result herolock.LockResult) {
		if end.Before(hero.GetMianDisappearTime()) {
			return
		}

		mianEndTime = end
		hero.SetMianRebackTime(hero.GetMianDisappearTime())
		result.Changed()
		result.Ok()
		setMian = true
	})

	if !setMian {
		return
	}

	realm := m.realm.GetBigMap()
	if realm == nil {
		logrus.Warnf("名城战参战的人要开免战，但没有主城")
		return
	}

	realm.Mian(hc.Id(), mianEndTime, true)
}

func (m *MingcWarModule) rebackMian(hc iface.HeroController) {
	var setMian bool
	var mianEndTime time.Time
	hc.FuncWithSend(func(hero *entity.Hero, result herolock.LockResult) {
		rebackTime := hero.GetMianRebackTime()
		if timeutil.IsZero(rebackTime) {
			return
		}
		if rebackTime.After(hero.GetMianDisappearTime()) {
			return
		}

		mianEndTime = rebackTime
		hero.SetMianRebackTime(time.Time{})
		result.Changed()
		result.Ok()
		setMian = true
	})

	if !setMian {
		return
	}

	realm := m.realm.GetBigMap()
	if realm == nil {
		logrus.Warnf("名城战参战的人要还原免战，但没有主城")
		return
	}

	realm.Mian(hc.Id(), mianEndTime, true)
}

//gogen:iface c2s_quit_fight
func (m *MingcWarModule) ProcessQuitFight(hc iface.HeroController) {
	succMsg, errMsg := m.mcWarSrv.QuitFight(hc.Id())
	if errMsg != nil {
		hc.Send(errMsg)
		return
	}

	hc.Send(succMsg)

	// 还原免战
	m.rebackMian(hc)
}

//gogen:iface
func (m *MingcWarModule) ProcessViewMingcWarMc(proto *mingc_war.C2SViewMingcWarMcProto, hc iface.HeroController) {
	mcId := u64.FromInt32(proto.Mcid)
	succMsg, errMsg := m.mcWarSrv.ViewMcWarMcMsg(mcId)
	if errMsg != nil {
		hc.Send(errMsg)
		return
	}

	hc.Send(succMsg)
}

//gogen:iface
func (m *MingcWarModule) ProcessSceneMove(proto *mingc_war.C2SSceneMoveProto, hc iface.HeroController) {
	dest := cb.XYCubeI32(proto.DestPosX, proto.DestPosY)
	succMsg, errMsg := m.mcWarSrv.SceneMove(hc.Id(), dest)
	if errMsg != nil {
		hc.Send(errMsg)
		return
	}

	hc.Send(succMsg)
}

//gogen:iface c2s_scene_back
func (m *MingcWarModule) ProcessSceneBack(hc iface.HeroController) {
	succMsg, errMsg := m.mcWarSrv.SceneBack(hc.Id())
	if errMsg != nil {
		hc.Send(errMsg)
		return
	}

	hc.Send(succMsg)
}

//gogen:iface c2s_scene_troop_relive
func (m *MingcWarModule) ProcessSceneTroopRelive(hc iface.HeroController) {
	succMsg, errMsg := m.mcWarSrv.SceneTroopRelive(hc.Id())
	if errMsg != nil {
		hc.Send(errMsg)
		return
	}

	hc.Send(succMsg)
}

//gogen:iface
func (m *MingcWarModule) ProcessSceneChangeMode(proto *mingc_war.C2SSceneChangeModeProto, hc iface.HeroController) {
	if proto.Mode <= 0 {
		hc.Send(mingc_war.ERR_SCENE_CHANGE_MODE_FAIL_INVALID_MODE)
		return
	}
	if _, ok := shared_proto.MingcWarModeType_name[proto.Mode]; !ok {
		hc.Send(mingc_war.ERR_SCENE_CHANGE_MODE_FAIL_INVALID_MODE)
		return
	}

	newMode := shared_proto.MingcWarModeType(proto.Mode)
	succMsg, errMsg := m.mcWarSrv.SceneChangeMode(hc.Id(), newMode)
	if errMsg != nil {
		hc.Send(errMsg)
		return
	}

	hc.Send(succMsg)
}

//gogen:iface
func (m *MingcWarModule) ProcessSceneTouShiBuildingTurnTo(proto *mingc_war.C2SSceneTouShiBuildingTurnToProto, hc iface.HeroController) {
	succMsg, errMsg := m.mcWarSrv.SceneTouShiBuildingTurnTo(hc.Id(), cb.XYCube(int(proto.PosX), int(proto.PosY)), proto.Left)
	if errMsg != nil {
		hc.Send(errMsg)
		return
	}

	hc.Send(succMsg)
}

//gogen:iface
func (m *MingcWarModule) ProcessSceneTouShiBuildingFire(proto *mingc_war.C2SSceneTouShiBuildingFireProto, hc iface.HeroController) {
	succMsg, errMsg := m.mcWarSrv.SceneTouShiBuildingFire(hc.Id(), cb.XYCube(int(proto.PosX), int(proto.PosY)))
	if errMsg != nil {
		hc.Send(errMsg)
		return
	}

	hc.Send(succMsg)
}

//gogen:iface c2s_scene_drum
func (m *MingcWarModule) ProcessSceneDrum(hc iface.HeroController) {
	if hero := m.dep.HeroSnapshot().Get(hc.Id()); hero == nil {
		logrus.Errorf("MingcWarModule.ProcessSceneDrum, 找不到 heroSnapshot。heroid:%v", hc.Id())
		hc.Send(mingc_war.ERR_SCENE_DRUM_FAIL_BAI_ZHAN_LEVEL_LIMIT)
		return
	} else if hero.BaiZhanJunXianLevel < m.dep.Datas().MingcMiscData().DrumMinBaiZhanLevel {
		hc.Send(mingc_war.ERR_SCENE_DRUM_FAIL_BAI_ZHAN_LEVEL_LIMIT)
		return
	}

	succMsg, errMsg := m.mcWarSrv.SceneDrum(hc.Id())
	if errMsg != nil {
		hc.Send(errMsg)
		return
	}

	hc.Send(succMsg)
}

//gogen:iface
func (m *MingcWarModule) ProcessViewMcWarRecord(proto *mingc_war.C2SViewMcWarRecordProto, hc iface.HeroController) {
	warId := u64.FromInt32(proto.WarId)
	mcId := u64.FromInt32(proto.McId)
	if mc := m.dep.Mingc().Mingc(mcId); mc == nil {
		hc.Send(mingc_war.ERR_VIEW_MC_WAR_RECORD_FAIL_NO_RECORD)
		return
	}

	var record *shared_proto.McWarFightRecordProto
	if err := ctxfunc.Timeout3s(func(ctx context.Context) (err error) {
		record, err = m.dep.Db().LoadMcWarRecord(ctx, warId, mcId)
		return
	}); err != nil || record == nil {
		logrus.WithError(err).Errorf("MingcWarModule.Db().LoadMcWarRecord 异常")
		hc.Send(mingc_war.ERR_VIEW_MC_WAR_RECORD_FAIL_NO_RECORD)
		return
	}

	hc.Send(mingc_war.NewS2cViewMcWarRecordMsg(proto.WarId, record))
}

//gogen:iface
func (m *MingcWarModule) ProcessViewMcWarTroopRecord(proto *mingc_war.C2SViewMcWarTroopRecordProto, hc iface.HeroController) {
	if mc := m.dep.Mingc().Mingc(u64.FromInt32(proto.McId)); mc == nil {
		hc.Send(mingc_war.ERR_VIEW_MC_WAR_TROOP_RECORD_FAIL_NO_RECORD)
		return
	}

	var record *shared_proto.McWarTroopAllRecordProto
	if err := ctxfunc.Timeout3s(func(ctx context.Context) (err error) {
		record, err = m.dep.Db().LoadMcWarHeroRecord(ctx, u64.FromInt32(proto.WarId), u64.FromInt32(proto.McId), hc.Id())
		return
	}); err != nil || record == nil {
		logrus.WithError(err).Errorf("MingcWarModule.Db().LoadMcWarHeroRecord 异常")
		hc.Send(mingc_war.ERR_VIEW_MC_WAR_TROOP_RECORD_FAIL_NO_RECORD)
		return
	}

	hc.Send(mingc_war.NewS2cViewMcWarTroopRecordMsg(record))
}

//gogen:iface
func (m *MingcWarModule) ProcessSceneSpeedUp(proto *mingc_war.C2SSceneSpeedUpProto, hc iface.HeroController) {
	goodsData := m.dep.Datas().GetGoodsData(u64.FromInt32(proto.GoodsId))
	if goodsData == nil {
		logrus.WithField("goods", proto.GoodsId).Debug("名城战行军加速，物品不存在")
		hc.Send(mingc_war.ERR_SCENE_SPEED_UP_FAIL_INVALID_GOODS)
		return
	}

	if goodsData.EffectType != shared_proto.GoodsEffectType_EFFECT_SPEED_UP {
		logrus.WithField("goods", goodsData.Name).Debug("名城战行军加速，物品不是行军加速物品")
		hc.Send(mingc_war.ERR_SCENE_SPEED_UP_FAIL_INVALID_GOODS)
		return
	}

	hctx := heromodule.NewContext(m.dep, operate_type.MingcWarSpeedUp)

	var reserveResult *entity.ReserveResult
	if hc.FuncWithSend(func(hero *entity.Hero, result herolock.LockResult) {
		ctime := m.dep.Time().CurrentTime()
		var hasEnoughGoods bool
		if proto.Money {
			if goodsData.YuanbaoPrice > 0 {
				if hasEnoughGoods, reserveResult = heromodule.ReserveYuanbao(hctx, hero, result, goodsData.YuanbaoPrice, ctime); !hasEnoughGoods {
					logrus.Debug("名城战行军加速，购买，元宝不足")
					result.Add(mingc_war.ERR_SCENE_SPEED_UP_FAIL_COST_NOT_ENOUGH)
					return
				}
			} else if goodsData.DianquanPrice > 0 {
				if hasEnoughGoods, reserveResult = heromodule.ReserveDianquan(hctx, hero, result, goodsData.DianquanPrice, ctime); !hasEnoughGoods {
					logrus.Debug("名城战行军加速，购买，点券不足")
					result.Add(mingc_war.ERR_SCENE_SPEED_UP_FAIL_COST_NOT_ENOUGH)
					return
				}
			} else {
				logrus.Debugf("名城战行军加速，不支持元宝或点券购买")
				result.Add(mingc_war.ERR_SCENE_SPEED_UP_FAIL_COST_NOT_SUPPORT)
				return
			}
		} else {
			// 预约扣物品
			if hasEnoughGoods, reserveResult = heromodule.ReserveGoods(hctx, hero, result, goodsData, 1, ctime); !hasEnoughGoods {
				logrus.Debug("名城战行军加速，预约扣除物品，个数不足")
				result.Add(mingc_war.ERR_SCENE_SPEED_UP_FAIL_GOODS_NOT_ENOUGH)
				return
			}
		}

		result.Changed()
		result.Ok()
	}) {
		return
	}

	heromodule.ConfirmReserveResult(hctx, hc, reserveResult, func() (success bool) {
		succMsg, errMsg := m.mcWarSrv.SceneSpeedUp(hc.Id(), goodsData.GoodsEffect.TroopSpeedUpRate)
		if errMsg != nil {
			hc.Send(errMsg)
			return
		}
		success = true
		hc.Send(succMsg)
		return
	})

}

//gogen:iface c2s_view_scene_troop_record
func (m *MingcWarModule) ProcessViewSceneTroopRecord(hc iface.HeroController) {
	succMsg, errMsg := m.mcWarSrv.ViewSceneTroopRecord(hc.Id())
	if errMsg != nil {
		hc.Send(errMsg)
		return
	}

	hc.Send(succMsg)
}

//gogen:iface
func (m *MingcWarModule) ProcessApplyRefreshRank(proto *mingc_war.C2SApplyRefreshRankProto, hc iface.HeroController)  {
	resultMsg, myRankMsg := m.mcWarSrv.CatchTroopsRank(hc.Id(), i32.Uint64(proto.Version))
	hc.Send(resultMsg)
	if myRankMsg != nil {
		hc.Send(myRankMsg)
	}
}

//gogen:iface
func (m *MingcWarModule) ProcessViewMyGuildMemberRank(proto *mingc_war.C2SViewMyGuildMemberRankProto, hc iface.HeroController)  {
	if gid, ok := hc.LockGetGuildId(); ok && gid > 0 {
		//if mc := m.dep.Mingc().Mingc(u64.FromInt32(proto.McId)); mc == nil {
		//	hc.Send(mingc_war.ERR_VIEW_MY_GUILD_MEMBER_RANK_FAIL_NO_RECORD)
		//	return
		//}

		var t *shared_proto.McWarTroopsInfoProto
		if err := ctxfunc.Timeout3s(func(ctx context.Context) (err error) {
			t, err = m.dep.Db().LoadMcWarGuildRecord(ctx, u64.FromInt32(proto.WarId), u64.FromInt32(proto.McId), gid)
			return
		}); err != nil || t == nil {
			logrus.WithError(err).Errorf("MingcWarModule.Db().LoadMcWarTroop 异常")
			hc.Send(mingc_war.ERR_VIEW_MY_GUILD_MEMBER_RANK_FAIL_NO_RECORD)
			return
		}

		hc.Send(mingc_war.NewS2cViewMyGuildMemberRankMsg(t))

	} else {
		logrus.Debugf("MingcWarModule hc.LockGetGuildId error gid:%v", gid)
		hc.Send(mingc_war.ERR_VIEW_MY_GUILD_MEMBER_RANK_FAIL_NO_RECORD)
	}
}
