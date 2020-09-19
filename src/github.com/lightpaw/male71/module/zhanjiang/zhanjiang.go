package zhanjiang

import (
	"github.com/lightpaw/logrus"
	"github.com/lightpaw/male7/entity"
	"github.com/lightpaw/male7/gen/iface"
	"github.com/lightpaw/male7/gen/pb/zhanjiang"
	"github.com/lightpaw/male7/service/conflict/heroservice/herolock"
	"github.com/lightpaw/male7/service/heromodule"
	"github.com/lightpaw/male7/util/must"
	"github.com/lightpaw/male7/util/u64"
	"github.com/lightpaw/male7/pb/server_proto"
	"github.com/lightpaw/male7/config/i18n"
	"github.com/lightpaw/male7/config/data"
	"github.com/lightpaw/male7/pb/shared_proto"
	"github.com/lightpaw/male7/service/operate_type"
)

func NewZhanJiangModule(dep iface.ServiceDep, fightService iface.FightService, militaryModule iface.MilitaryModule, guildSnapshotService iface.GuildSnapshotService) *ZhanJiangModule {
	m := &ZhanJiangModule{
		dep:                  dep,
		datas:                dep.Datas(),
		timeService:          dep.Time(),
		fightService:         fightService,
		militaryModule:       militaryModule,
		guildSnapshotService: guildSnapshotService,
		broadcast:            dep.Broadcast(),
	}
	return m
}

//gogen:iface
type ZhanJiangModule struct {
	dep                  iface.ServiceDep
	datas                iface.ConfigDatas
	timeService          iface.TimeService
	fightService         iface.FightService
	militaryModule       iface.MilitaryModule
	guildSnapshotService iface.GuildSnapshotService
	broadcast            iface.BroadcastService
}

//gogen:iface
func (m *ZhanJiangModule) ProcessOpen(proto *zhanjiang.C2SOpenProto, hc iface.HeroController) {
	data := m.datas.GetZhanJiangGuanQiaData(u64.FromInt32(proto.Id))
	if data == nil {
		logrus.Debugf("斩将数据没找到")
		hc.Send(zhanjiang.ERR_OPEN_FAIL_NOT_FOUND)
		return
	}

	hc.FuncWithSend(func(hero *entity.Hero, result herolock.LockResult) {
		heroZhanJiang := hero.ZhanJiang()

		if data.Prev != nil && !heroZhanJiang.IsPass(data.Prev) {
			logrus.Debugf("斩将前置关卡没有通关")
			result.Add(zhanjiang.ERR_OPEN_FAIL_PREV_NOT_PASS)
			return
		}

		if heroZhanJiang.CurChallenge() != nil {
			logrus.Debugf("斩将当前已经开启了，请挑战完")
			result.Add(zhanjiang.ERR_OPEN_FAIL_IS_OPEN)
			return
		}

		if heroZhanJiang.OpenTimes() >= m.datas.ZhanJiangMiscData().DefaultTimes {
			logrus.Debugf("斩将没次数了")
			result.Add(zhanjiang.ERR_OPEN_FAIL_NO_OPEN_TIMES)
			return
		}

		result.Ok()
		result.Changed()

		// 开启
		captainId := heroZhanJiang.LastCaptainId()
		if captainId == 0 || hero.Military().Captain(captainId) == nil || !heroZhanJiang.IsPass(data) {
			// 首次通关、没有设置过、武将没找到
			c := hero.Military().MaxFightingAmountCaptain()
			if c != nil {
				captainId = c.Id()
			}
		}
		heroZhanJiang.StartChallenge(entity.NewZhanJiangChallenge(data, captainId))

		result.Add(zhanjiang.NewS2cOpenMsg(u64.Int32(data.Id), u64.Int32(captainId)))
	})
}

//gogen:iface c2s_give_up
func (m *ZhanJiangModule) ProcessGiveUp(hc iface.HeroController) {
	hc.FuncWithSend(func(hero *entity.Hero, result herolock.LockResult) {
		heroZhanJiang := hero.ZhanJiang()

		if heroZhanJiang.CurChallenge() == nil {
			logrus.Debugf("斩将未开启，无法放弃")
			result.Add(zhanjiang.ERR_GIVE_UP_FAIL_NOT_OPEN)
			return
		}

		result.Ok()
		result.Changed()

		// 开启
		heroZhanJiang.EndChallenge()

		result.Add(zhanjiang.GIVE_UP_S2C)
	})
}

//gogen:iface c2s_update_captain
func (m *ZhanJiangModule) ProcessUpdateCaptain(proto *zhanjiang.C2SUpdateCaptainProto, hc iface.HeroController) {
	hc.FuncWithSend(func(hero *entity.Hero, result herolock.LockResult) {
		captain := hero.Military().Captain(u64.FromInt32(proto.Id))
		if captain == nil {
			logrus.Debugf("武将没找到")
			result.Add(zhanjiang.ERR_UPDATE_CAPTAIN_FAIL_NOT_FOUND)
			return
		}

		heroZhanJiang := hero.ZhanJiang()

		if heroZhanJiang.CurChallenge() == nil {
			logrus.Debugf("斩将未开启，无法放弃")
			result.Add(zhanjiang.ERR_UPDATE_CAPTAIN_FAIL_NOT_OPEN)
			return
		}

		result.Ok()
		result.Changed()

		heroZhanJiang.CurChallenge().SetCaptainId(captain.Id())

		result.Add(zhanjiang.NewS2cUpdateCaptainMsg(proto.Id))
	})
}

//gogen:iface c2s_challenge
func (m *ZhanJiangModule) ProcessChallenge(proto *zhanjiang.C2SChallengeProto, hc iface.HeroController) {
	hc.FuncWithSend(func(hero *entity.Hero, result herolock.LockResult) {
		heroZhanJiang := hero.ZhanJiang()
		curChallenge := heroZhanJiang.CurChallenge()
		if curChallenge == nil {
			logrus.Debugf("未开启")
			result.Add(zhanjiang.ERR_CHALLENGE_FAIL_NOT_OPEN)
			return
		}

		if curChallenge.PassCount() != u64.FromInt32(proto.PassCount) {
			logrus.Debugf("请求过于频繁")
			result.Add(zhanjiang.ERR_CHALLENGE_FAIL_NO_DUPLICATE)
			return
		}

		captain := hero.Military().Captain(curChallenge.CaptainId())
		if captain == nil {
			logrus.Debugf("未设置出战武将")
			result.Add(zhanjiang.ERR_CHALLENGE_FAIL_NOT_SET_CAPTAIN)
			return
		}

		// 写死在第三个
		attacker := hero.GenCombatPlayerProtoWithCaptains(true, []*entity.TroopPos{nil, nil, captain.NewTroopPos(0)}, m.guildSnapshotService.GetSnapshot)
		if attacker == nil {
			logrus.Errorf("竟然一个武将生成proto失败了")
			result.Add(zhanjiang.ERR_CHALLENGE_FAIL_SERVER_ERROR)
			return
		}

		zhanJiangData := curChallenge.GuanQia().ZhanJiangDatas[curChallenge.PassCount()]

		tfctx := entity.NewTlogFightContext(operate_type.BattleZhanJiang, zhanJiangData.Id, operate_type.CountryChange, 0)
		response := m.fightService.SendFightRequest(tfctx, zhanJiangData.CombatScene, hero.Id(), zhanJiangData.Monster.GetNpcId(), attacker, zhanJiangData.Monster.GetPlayer())
		if response == nil {
			logrus.Errorf("过关斩将，response==nil")
			result.Add(zhanjiang.ERR_CHALLENGE_FAIL_SERVER_ERROR)
			return
		}

		if response.ReturnCode != 0 {
			logrus.Errorf("过关斩将，战斗计算发生错误，%s", response.ReturnMsg)
			result.Add(zhanjiang.ERR_CHALLENGE_FAIL_SERVER_ERROR)
			return
		}

		result.Changed()
		result.Ok()

		hero.HistoryAmount().IncreaseOne(server_proto.HistoryAmountType_AccumZhanJiangGuanqia)
		heromodule.UpdateTaskProgress(hero, result, shared_proto.TaskTargetType_TASK_TARGET_ZHANJIANG_START_TIMES)

		if !response.AttackerWin {
			result.Add(zhanjiang.NewS2cChallengeMsg(false, response.Link, must.Marshal(response.AttackerShare), u64.Int32(zhanJiangData.Id)))
			return
		}

		curChallenge.IncPassCount()

		prize := zhanJiangData.PassPrize
		if zhanJiangData.Plunder != nil {
			prize = zhanJiangData.Plunder.Try()
		}

		oldExp := captain.AbilityExp()
		hctx := heromodule.NewContext(m.dep, operate_type.ZhanJiangChallenge)
		// 给奖励
		ctime := m.timeService.CurrentTime()
		if prize != nil {
			heromodule.AddCaptainPrize(hctx, hero, result, prize, []uint64{captain.Id()}, ctime)
		}
		heromodule.AddGongXun(hero, result, captain, zhanJiangData.GongXun)

		m.dep.Tlog().TlogGUOGUANFlow(hero, curChallenge.CaptainId(), captain.Level(), oldExp, captain.AbilityExp(), u64.Sub(captain.AbilityExp(), oldExp), zhanJiangData.GongXun)

		if curChallenge.IsAllPass() {
			curGuanQia := curChallenge.GuanQia()
			heroZhanJiang.EndChallenge()
			heroZhanJiang.Pass(curChallenge.GuanQia())
			heroZhanJiang.SetLastCaptainId(curChallenge.CaptainId())

			//toAdd := curChallenge.GuanQia().AbilityExp
			//if captain.Ability() >= captain.RebirthLevelData().AbilityLimit {
			//	// 可以添加的经验值，只能加到升级前1点
			//	maxCanAdd := u64.Sub(captain.AbilityData().UpgradeExp, captain.AbilityExp()+1)
			//	toAdd = u64.Min(toAdd, maxCanAdd)
			//}
			//
			//if toAdd > 0 {
			//	// 给武将加成长值
			//	hctx := heromodule.NewContext(m.dep, operate_type.MilitaryCaptainUpgradeAbilityExp)
			//	heromodule.AddCaptainAbilityExp(hctx, hero, result, captain, toAdd)
			//}

			// 成长值达到上限的话不加
			if captain.Ability() < captain.RebirthLevelData().AbilityLimit {
				toAdd := curChallenge.GuanQia().AbilityExp
				if toAdd > 0 {
					hctx := heromodule.NewContext(m.dep, operate_type.MilitaryCaptainUpgradeAbilityExp)
					heromodule.AddCaptainAbilityExp(hctx, hero, result, captain, toAdd, ctime)
				}
			}

			heromodule.UpdateTaskProgress(hero, result, shared_proto.TaskTargetType_TASK_TARGET_ZHANJIANG_GUANQIA_COMPLETE)

			// 系统广播
			if d := hctx.BroadcastHelp().ZhanJiangZhangJieComplete; d != nil {
				hctx.AddBroadcast(d, hero, result, captain.Id(), zhanJiangData.Id, func() *i18n.Fields {
					text := d.NewTextFields()
					text.WithClickHeroFields(data.KeySelf, hctx.GetFlagHeroName(hero), hero.Id())
					text.WithClickCaptainFields(data.KeyCaptain, hctx.Broadcast().GetCaptainText(captain), hero.Id(), captain.Id())
					text.WithFields(data.KeyText, curGuanQia.ChapterData.ChapterName)
					return text
				})
			}
		}

		result.Add(zhanjiang.NewS2cChallengeMsg(true, response.Link, must.Marshal(response.AttackerShare), u64.Int32(zhanJiangData.Id)))

		if curChallenge.IsAllPass() {
			result.Add(curChallenge.GuanQia().PassMsg)
		}
	})
}
