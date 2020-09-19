package worker

import (
	"github.com/lightpaw/logrus"
	"github.com/lightpaw/male7/entity"
	"github.com/lightpaw/male7/gen/pb/country"
	"github.com/lightpaw/male7/gen/pb/domestic"
	"github.com/lightpaw/male7/gen/pb/login"
	"github.com/lightpaw/male7/gen/pb/mingc_war"
	"github.com/lightpaw/male7/gen/pb/misc"
	"github.com/lightpaw/male7/gen/service"
	"github.com/lightpaw/male7/pb/shared_proto"
	"github.com/lightpaw/male7/service/conflict/heroservice/herolock"
	"github.com/lightpaw/male7/service/heromodule"
	"github.com/lightpaw/male7/service/herosnapshot/snapshotdata"
	"github.com/lightpaw/male7/service/operate_type"
	"github.com/lightpaw/male7/util/u64"
)

func (m *MessageWorker) processLoadedMsg(data *service.MsgData) bool {
	if m.user == nil {
		logrus.Errorf("worker received loaded msg, but not login, received: %d-%d",
			data.ModuleID, data.SequenceID)
		m.Send(login.ERR_LOADED_FAIL_NO_LOGIN)
		m.Close()
		return false
	}

	hc := m.user.GetHeroController()
	if hc == nil {
		logrus.Errorf("worker received loaded msg, but not create hero, received: %d-%d",
			data.ModuleID, data.SequenceID)
		m.Send(login.ERR_LOADED_FAIL_NO_CREATED)
		m.Close()
		return false
	}

	hctx := heromodule.NewContext(service.ServiceDep, operate_type.HeroLoginLoad)

	ctime := service.TimeService.CurrentTime()

	heroId := hc.Id()
	logrus.Debugf("worker.processLoadedMsg() hero: %v", heroId)

	// weekly reset
	resetWeeklyTime := service.TickerService.GetWeeklyTickTime().GetPrevTickTime()
	// daily reset
	resetDailyTime := service.TickerService.GetDailyTickTime().GetPrevTickTime()
	resetDailyZeroTime := service.TickerService.GetDailyZeroTickTime().GetPrevTickTime()
	resetDailyMcTime := service.TickerService.GetDailyMcTickTime().GetPrevTickTime()

	// 轩辕会武每日重置
	xuanyResetTime := service.XuanyuanModule.GetResetTickTime().GetPrevTickTime()

	// season reset
	seasonResetTime := service.SeasonService.GetSeasonTickTime().GetPrevTickTime()
	curSeason := service.SeasonService.SeasonByTime(seasonResetTime)

	var heroName string
	var guildId, baseRegion int64
	var countryId uint64

	var snapshot *snapshotdata.HeroSnapshot
	if hc.FuncWithSend(func(hero *entity.Hero, result herolock.LockResult) {

		heroName = hero.Name()
		guildId = hero.GuildId()
		baseRegion = hero.BaseRegion()
		countryId = hero.CountryId()

		// 每日22点重置
		if hero.GetDailyMcResetTime().Before(resetDailyMcTime) {
			hero.ResetDailyMc(resetDailyMcTime, service.ConfigDatas)
		}

		// 每周重置
		if hero.GetWeeklyResetTime().Before(resetWeeklyTime) {
			hero.ResetWeekly(resetWeeklyTime)
		}
		// 每日0点重置
		if hero.GetDailyZeroResetTime().Before(resetDailyZeroTime) {
			hero.ResetDailyZero(resetDailyZeroTime, service.ConfigDatas)
			// vip登录
			heromodule.VipResetDailyZero(hctx, hero, result, service.ConfigDatas, ctime)
		}

		// 每日重置
		isResetDaily := hero.GetDailyResetTime().Before(resetDailyTime)
		if isResetDaily {
			hero.ResetDaily(resetDailyTime, service.ConfigDatas)
			// tlog
			service.TlogService.TlogResourceStockFlow(hero, hero.GetYuanbao(), hero.GetDianquan(), hero.Military().FreeSoldier(ctime), hero.GetAllRes(shared_proto.ResType_GOLD), hero.GetAllRes(shared_proto.ResType_STONE), hero.GetGuildContributionCoin())
		}

		// 重置轩辕会武
		hero.Xuanyuan().TryResetDaily(xuanyResetTime)

		// season reset
		if hero.GetSeasonResetTime().Before(seasonResetTime) {
			heromodule.HeroResetSeason(hero, result, ctime, seasonResetTime, curSeason)
		}

		// 收一次税收
		hero.UpdateTax(ctime, service.ConfigDatas.MiscGenConfig().TaxDuration, service.ConfigDatas.GetBuffEffectData)

		if heroProto, err := hero.EncodeClient(ctime, service.GuildService.GetSnapshot, service.ConfigDatas, service.IndividualServerConfig.GetServerInfo()).Marshal(); err != nil {
			logrus.WithError(err).Errorf("worker.processLoginModuleMsg hero.Marshal error, %d-%v", heroId, heroName)
			return
		} else {
			result.Add(login.NewS2cLoadedMsg(heroProto))
		}

		if isResetDaily {
			heromodule.UpdateTaskProgress(hero, result, shared_proto.TaskTargetType_TASK_TARGET_ACCUM_LOGIN_DAY)
		}

		if !heromodule.TryRefreshWorkshop(hctx, hero, result, service.ConfigDatas, ctime) {
			heromodule.SendWorkshopMsg(hero, result, service.ConfigDatas)
		}

		if !heromodule.TryRefreshBlackMarketGoods(hero, result, service.ConfigDatas, ctime, false) {
			heromodule.SendBlackMarketGoodsMsg(hero, result, false)
		}

		// 推送主城增益数量
		heromodule.UpdateAdvantageCount(hero, result, service.ConfigDatas, ctime)
		result.Add(domestic.NewS2cUpdateAdvantageCountMsg(int32(hero.CurrentAdvantageCount())))

		// 上线了
		hero.Online(ctime)

		snapshot = service.HeroSnapshotService.NewSnapshot(hero)

		result.Ok()
	}) {
		m.Send(login.ERR_LOGIN_FAIL_SERVER_ERROR)

		m.Close()
		return false
	}

	logrus.Debugf("worker.processLoadedMsg() end hero: %v-%s", heroId, heroName)

	service.HeroSnapshotService.Online(snapshot)

	service.GuildModule.OnHeroOnline(hc, guildId)

	service.CountryService.OnHeroOnline(hc, countryId)

	service.LocationHeroCache.UpdateHero(snapshot, ctime.Unix())

	// 名城战在战斗阶段，推送展示图标
	if msg, isFightState := service.MingcWarService.BuildFightStartMsg(ctime); isFightState {
		m.Send(msg)
	}

	// 在名城战参战，需要被传入副本场景中
	if mcId, ok := service.MingcWarService.JoiningFightMingc(hc.Id()); ok {
		m.Send(mingc_war.NewS2cIsJoiningFightOnLoginMsg(u64.Int32(mcId)))
	}

	if gid, ok := hc.LockGetGuildId(); ok {
		// 在名城战申请攻打和协助阶段，推送红点
		if service.MingcWarService.ApplyAtkNotice(hc.Id(), gid) || service.MingcWarService.ApplyAstNotice(hc.Id(), gid) {
			m.Send(mingc_war.RED_POINT_NOTICE_S2C)
		}
	}

	if countryId := hc.LockHeroCountry(); countryId > 0 {
		if c := service.CountryService.Country(countryId); c != nil && c.InChangeNameVoteDuration(ctime) {
			m.Send(country.CHANGE_NAME_START_NOTICE_S2C)
		}
	}

	r := service.RealmService.GetBigMap()
	hc.Send(r.GetMapData().GetRadiusBlock(r.GetRadius()).GetUpdateMapRadiusMsg())

	if baseRegion != 0 {
		r := service.RealmService.GetBigMap()
		if r != nil {
			r.OnHeroLogin(hc)
		} else {
			// 主城对应的场景不存在，流亡
			hc.FuncWithSend(func(hero *entity.Hero, result herolock.LockResult) {

				if hero.BaseRegion() == baseRegion {
					logrus.Errorf("英雄登陆 %d %s，主城对应的场景不存在，删除主城，region: %d", hero.Id(), hero.Name(), hero.BaseRegion())

					// 这个场景不存在了，当没有主城了，下次登陆重新创建
					hero.ClearBase()

					// 踢下线
					hc.Disconnect(misc.ErrDisconectReasonFailKick)
				}
			})
		}
	}

	heromodule.OnHeroOnlineEvent(hc)

	m.user.SetLoaded()

	return true
}
