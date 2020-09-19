package gm

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"github.com/lightpaw/logrus"
	"github.com/lightpaw/male7/config/towerdata"
	"github.com/lightpaw/male7/entity"
	"github.com/lightpaw/male7/gen/iface"
	"github.com/lightpaw/male7/gen/pb/gm"
	"github.com/lightpaw/male7/gen/pb/misc"
	"github.com/lightpaw/male7/gen/pb/region"
	"github.com/lightpaw/male7/gen/pb/tag"
	"github.com/lightpaw/male7/module/rank/ranklist"
	"github.com/lightpaw/male7/pb/shared_proto"
	"github.com/lightpaw/male7/service/conflict/heroservice/herolock"
	"github.com/lightpaw/male7/service/heromodule"
	"github.com/lightpaw/male7/util/i64"
	"github.com/lightpaw/male7/util/must"
	"github.com/lightpaw/male7/util/u64"
	rand2 "math/rand"
	"github.com/lightpaw/male7/config/data"
	"github.com/lightpaw/male7/gen/pb/military"
)

func (m *GmModule) newMiscGmGroup() *gm_group {
	return &gm_group{
		tab: "其他",
		handler: []*gm_handler{
			newHeroIntHandler("千重楼爬塔", "1", m.towerUp),
			newIntHandler("设置军衔", "1", m.junXian),
			newHeroIntHandler("标签", "4", m.tag),
			newIntHandler("个人紧急军情", "4", m.selfReminder),
			newIntHandler("联盟紧急军情", "4", m.guildReminder),
			newHeroIntHandler("武将加功勋", "10000", m.addCaptainGongxun),
			newStringHandler("完成问卷调查", "1", m.completeSurvey),
			newStringHandler("开启抗击匈奴", "", m.openResistXiongNu),
			newStringHandler("系统广播", "", m.sendSysBroadcast),
			newStringHandler("系统定时广播", "", m.sendSysTimingBroadcast),
			newStringHandler("农场一键成熟", "", m.farmRipe),
			newStringHandler("农场一键可偷", "", m.farmCanSteal),
		},
	}
}

func (m *GmModule) junXian(amount int64, hc iface.HeroController) {
	m.modules.BaiZhanModule().GmSetJunXian(amount, hc)
}

func (m *GmModule) selfReminder(amount int64, hc iface.HeroController) {
	hc.Send(region.NewS2cSelfBeenAttackRobChangedMsg(i64.Int32(amount), i64.Int32(0)))
}

func (m *GmModule) guildReminder(amount int64, hc iface.HeroController) {
	hc.Send(region.NewS2cGuildBeenAttackRobChangedMsg(i64.Int32(amount)))
}

// 完成问卷调查
func (m *GmModule) completeSurvey(amount string, hc iface.HeroController) {
	data := m.datas.GetSurveyData(amount)
	if data == nil {
		hc.Send(gm.NewS2cGmMsg("问卷调查数据没找到"))
		return
	}

	if heromodule.IsLocked(hc.Func, m.sharedGuildService.GetSnapshot, data.Condition) {
		hc.Send(gm.NewS2cGmMsg(fmt.Sprintf("问卷调查未解锁: %v-%s", data.Id, data.Name)))
		return
	}

	hc.FuncWithSend(func(hero *entity.Hero, result herolock.LockResult) {
		if hero.Survey().IsCompleted(data) {
			hc.Send(gm.NewS2cGmMsg(fmt.Sprintf("问卷调查奖励已经领取: %v-%s", data.Id, data.Name)))
			return
		}

		result.Changed()
		result.Ok()
		hero.Survey().Complete(data)

		proto := m.datas.MailHelp().SurveyMail.NewTextMail(shared_proto.MailType_MailNormal)
		proto.Prize = data.PrizeProto

		ctime := m.time.CurrentTime()
		m.modules.MailModule().SendProtoMail(hc.Id(), proto, ctime)

		result.Add(data.CompleteMsg)

		hc.Send(gm.NewS2cGmMsg("问卷调查奖励已经领取"))
	})
}

func (m *GmModule) tag(amount int64, hero *entity.Hero, result herolock.LockResult, hc iface.HeroController) {
	var guildFlagName string
	if snapshot := m.sharedGuildService.GetSnapshot(hero.GuildId()); snapshot != nil {
		guildFlagName = snapshot.FlagName
	}

	for i := int64(0); i < amount; i++ {
		if uint64(hero.Tag().TagCount()) >= m.datas.TagMiscData().MaxCount {
			return
		}

		bytes := make([]byte, rand2.Intn(8)+3)
		rand.Read(bytes)

		addTag := base64.StdEncoding.EncodeToString(bytes)
		t, record := hero.Tag().AddTag(hero.IdBytes(), hero.Name(), guildFlagName, addTag, m.time.CurrentTime())

		result.Add(tag.NewS2cAddOrUpdateTagMsg(hero.IdBytes(), must.Marshal(record), must.Marshal(t)))
	}
}

func (m *GmModule) towerUp(floor int64, hero *entity.Hero, result herolock.LockResult, hc iface.HeroController) {

	heroTower := hero.Tower()

	var nextFloorData *towerdata.TowerData
	if heroTower.CurrentFloor() <= 0 {
		nextFloorData = m.datas.TowerData().MinKeyData
	} else {
		curTowerData := m.datas.GetTowerData(heroTower.CurrentFloor())
		if curTowerData == nil {
			logrus.Debugf("我去，当前层竟然数据都找不到![%v]", heroTower.CurrentFloor())
			return
		}

		nextFloorData = curTowerData.NextFloor()
	}

	if nextFloorData == nil {
		logrus.Debugln("已经爬到最后一层了")
		return
	}

	var historyMaxFloorChanged bool

	for i := int64(0); i < floor; i++ {
		// 爬一层
		if nextFloorData == nil {
			break
		}

		if nextFloorData.UnlockSecretTower != nil {
			// 激活密室
			heroSecretTower := hero.SecretTower()

			openNewSecretTower := heroSecretTower.TryOpenAndGiveDefaultTimes(nextFloorData.UnlockSecretTower)
			if openNewSecretTower {
				result.Add(nextFloorData.UnlockSecretTower.UnlockMsg)
			}
		}

		if heroTower.IncreseCurrentFloor(m.time.CurrentTime(), m.datas.MiscConfig().TowerAutoKeepFloor) {
			historyMaxFloorChanged = true
		}

		nextFloorData = nextFloorData.NextFloor()
	}

	if historyMaxFloorChanged {
		// 历史最高层数变更了
		m.modules.RankModule().AddOrUpdateRankObj(ranklist.NewTowerRankObj(m.heroSnapshotService.Get, hc.Id(), heroTower.HistoryMaxFloor(), heroTower.HistoryMaxFloorTime()))
		result.Changed()
	}

	// 更新任务进度
	heromodule.UpdateTaskProgress(hero, result, shared_proto.TaskTargetType_TASK_TARGET_TOWER_FLOOR)
	hc.Disconnect(misc.ErrDisconectReasonFailGm)
}

func (m *GmModule) towerUpTo(floor int64, hero *entity.Hero, result herolock.LockResult, hc iface.HeroController) {

	heroTower := hero.Tower()

	if int64(heroTower.CurrentFloor()) < floor {
		m.towerUp(floor-int64(heroTower.CurrentFloor()), hero, result, hc)
	}
}

func (m *GmModule) addCaptainGongxun(toAdd int64, hero *entity.Hero, result herolock.LockResult, hc iface.HeroController) {
	for _, c := range hero.Military().Captains() {
		newGongXun := c.AddGongXun(u64.FromInt64(toAdd))
		result.Add(military.NewS2cAddGongxunMsg(u64.Int32(c.Id()), u64.Int32(newGongXun)))
		result.Changed()
	}
}

func (m *GmModule) openResistXiongNu(input string, hc iface.HeroController) {
	guildId, _ := hc.LockGetGuildId()
	if guildId == 0 {
		hc.Send(gm.NewS2cGmMsg("没有联盟"))
		return
	}

	suc := m.modules.XiongNuModule().GmStart(hc, guildId)
	if suc {
		hc.Send(gm.NewS2cGmMsg("开启成功"))
	} else {
		hc.Send(gm.NewS2cGmMsg("开启失败"))
	}
}

func (m *GmModule) farmRipe(input string, hc iface.HeroController) {
	m.farmService.GMRipe(hc.Id())
}

func (m *GmModule) farmCanSteal(input string, hc iface.HeroController) {
	m.farmService.GMCanSteal(hc.Id())
}

func (m *GmModule) sendSysTimingBroadcast(input string, hc iface.HeroController) {
	hc.FuncWithSend(func(hero *entity.Hero, result herolock.LockResult) {
		funcList := []func(){
			func() {
				d := m.datas.BroadcastHelp().BaiZhanTouFang
				result.AddBroadcast(d.TimingBroadcastMsg)
			},
		}

		for _, f := range funcList {
			f()
		}
	})
}

func (m *GmModule) sendSysBroadcast(input string, hc iface.HeroController) {
	hc.FuncWithSend(func(hero *entity.Hero, result herolock.LockResult) {
		hctx := m.hctx
		funcList := []func() string{
			func() string {
				d := m.datas.BroadcastHelp().BaseLevel
				text := d.NewTextFields()
				text.WithClickHeroFields(data.KeySelf, hctx.GetFlagHeroName(hero), hero.Id())
				text.WithFields(data.KeyNum, hero.BaseLevel())

				result.AddBroadcast(misc.NewS2cSysBroadcastMsg(text.JsonString()).Static())
				return text.JsonString()
			},
			func() string {
				if len(hero.Military().Captains()) <= 0 {
					return ""
				}
				var captain *entity.Captain
				for _, v := range hero.Military().Captains() {
					captain = v
					break
				}
				d := m.datas.BroadcastHelp().CaptainAbilityExp
				text := d.NewTextFields().WithClickHeroFields(data.KeySelf, hctx.GetFlagHeroName(hero), hero.Id())
				text.WithFields(data.KeyQuality, hctx.Datas().GetColorData(uint64(captain.Quality())).QualityColorText)

				result.AddBroadcast(misc.NewS2cSysBroadcastMsg(text.JsonString()).Static())
				return text.JsonString()
			},
			func() string {
				if len(hero.Military().Captains()) <= 0 {
					return ""
				}
				var captain *entity.Captain
				for _, v := range hero.Military().Captains() {
					captain = v
					break
				}
				d := m.datas.BroadcastHelp().CaptainReBrith
				text := d.NewTextFields().WithClickHeroFields(data.KeySelf, hctx.GetFlagHeroName(hero), hero.Id())
				text.WithFields(data.KeyNum, captain.RebirthLevel())

				result.AddBroadcast(misc.NewS2cSysBroadcastMsg(text.JsonString()).Static())
				return text.JsonString()
			},
			func() string {
				d := m.datas.BroadcastHelp().HeroLevel
				text := d.NewTextFields()
				text.WithClickHeroFields(data.KeySelf, hctx.GetFlagHeroName(hero), hero.Id())
				text.WithFields(data.KeyNum, hero.Level())

				result.AddBroadcast(misc.NewS2cSysBroadcastMsg(text.JsonString()).Static())
				return text.JsonString()
			},
			func() string {
				d := m.datas.BroadcastHelp().JiuGuanBaoJi
				text := d.NewTextFields().WithClickHeroFields(data.KeySelf, hctx.GetFlagHeroName(hero), hero.Id())
				text.WithFields(data.KeyNum, 10)

				result.AddBroadcast(misc.NewS2cSysBroadcastMsg(text.JsonString()).Static())
				return text.JsonString()
			},
			func() string {
				tutor := m.datas.TutorData().MaxKeyData
				d := m.datas.BroadcastHelp().JiuGuanQingJiao
				text := d.NewTextFields().WithClickHeroFields(data.KeySelf, hctx.GetFlagHeroName(hero), hero.Id())
				text.WithFields(data.KeyNpc, tutor.Name)

				result.AddBroadcast(misc.NewS2cSysBroadcastMsg(text.JsonString()).Static())
				return text.JsonString()
			},
			func() string {
				d := m.datas.BroadcastHelp().TowerFloor
				text := d.NewTextFields()
				text.WithClickHeroFields(data.KeySelf, hctx.GetFlagHeroName(hero), hero.Id())
				text.WithFields(data.KeyNum, hero.Tower().HistoryMaxFloor())

				result.AddBroadcast(misc.NewS2cSysBroadcastMsg(text.JsonString()).Static())
				return text.JsonString()
			},
			func() string {
				d := m.datas.BroadcastHelp().GemGet
				gemData := m.datas.GemData().Must(1)
				text := d.NewTextFields()
				text.WithClickHeroFields(data.KeySelf, hctx.GetFlagHeroName(hero), hero.Id())
				text.WithFields(data.KeyNum, gemData.Level)

				result.AddBroadcast(misc.NewS2cSysBroadcastMsg(text.JsonString()).Static())
				return text.JsonString()
			},
			func() string {
				d := m.datas.BroadcastHelp().GuildCreate
				g := m.dep.GuildSnapshot().GetSnapshot(hero.GuildId())
				if g == nil {
					return ""
				}
				text := d.NewTextFields().WithClickHeroFields(data.KeySelf, hctx.GetFlagName(g.FlagName, g.Name), g.Id)
				text.WithClickHeroFields(data.KeyName, hero.Name(), hero.Id())

				result.AddBroadcast(misc.NewS2cSysBroadcastMsg(text.JsonString()).Static())
				return text.JsonString()
			},
			func() string {
				d := m.datas.BroadcastHelp().GuildLevel
				g := m.dep.GuildSnapshot().GetSnapshot(hero.GuildId())
				if g == nil {
					return ""
				}
				text := d.NewTextFields()
				text.WithClickGuildFields(data.KeyGuild, hctx.GetFlagName(g.FlagName, g.Name), g.Id)
				text.WithFields(data.KeyNum, g.GuildLevel.Level)

				result.AddBroadcast(misc.NewS2cSysBroadcastMsg(text.JsonString()).Static())
				return text.JsonString()
			},
			func() string {
				d := m.datas.BroadcastHelp().FishCaptainSoul
				text := d.NewTextFields().WithClickHeroFields(data.KeySelf, hctx.GetFlagHeroName(hero), hero.Id())

				csData := m.datas.CaptainData().Must(1)
				colorData := m.datas.ColorData().Must(uint64(csData.Rarity.Color))
				text.WithFields(data.KeyText, csData.Rarity.Name)
				text.WithClickCaptainFields(data.KeyCaptain, colorData.GetColorText(csData.Name), hero.Id(), csData.Id)

				result.AddBroadcast(misc.NewS2cSysBroadcastMsg(text.JsonString()).Static())
				return text.JsonString()
			},
		}

		for _, f := range funcList {
			content := f()
			logrus.Debugf("========== gm命令 系统广播：%v", content)
		}
	})
}
