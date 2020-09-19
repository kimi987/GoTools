package gm

import (
	"github.com/lightpaw/logrus"
	"github.com/lightpaw/male7/config/dungeon"
	"github.com/lightpaw/male7/entity"
	"github.com/lightpaw/male7/gen/iface"
	dm "github.com/lightpaw/male7/gen/pb/dungeon"
	"github.com/lightpaw/male7/gen/pb/misc"
	"github.com/lightpaw/male7/pb/shared_proto"
	"github.com/lightpaw/male7/service/conflict/heroservice/herolock"
	"github.com/lightpaw/male7/service/heromodule"
	"github.com/lightpaw/male7/util/must"
	"github.com/lightpaw/male7/util/u64"
)

func (m *GmModule) newDungeonGmGroup() *gm_group {
	return &gm_group{
		tab: "副本",
		handler: []*gm_handler{
			newStringHandler("重置推图", "", m.resetDungeon),
			newStringHandler("重置推图扫荡次数", "", m.resetDungeonAutoTimes),
			newHeroStringHandler("通关1章", "", m.pass1Chapter),
			newHeroStringHandler("通关1关", "", m.pass1Times),
			newHeroStringHandler("通关1章（精英）", "", m.pass1EliteChapter),
			newHeroStringHandler("通关1关（精英）", "", m.pass1EliteTimes),
		},
	}
}

func (m *GmModule) resetDungeon(input string, hc iface.HeroController) {
	hc.FuncWithSend(func(hero *entity.Hero, result herolock.LockResult) {
		heroDungeon := hero.Dungeon()
		heroDungeon.GMReset()

		result.Changed()
		result.Ok()
	})

	hc.Disconnect(misc.ErrDisconectReasonFailGm)
}

func (m *GmModule) resetDungeonAutoTimes(input string, hc iface.HeroController) {
	hc.FuncWithSend(func(hero *entity.Hero, result herolock.LockResult) {
		ctime := m.time.CurrentTime()

		heroDungeon := hero.Dungeon()
		t := heroDungeon.ChallengeTimes()
		t.SetTimes(t.MaxTimes(), ctime)

		result.Add(dm.NewS2cUpdateChallengeTimesMsg(t.StartTimeUnix32()))

		result.Ok()
	})
}

func (m *GmModule) pass1Chapter(input string, hero *entity.Hero, result herolock.LockResult, hc iface.HeroController) {
	m.passDungeon(hero, result, hc, 10, true, shared_proto.DifficultType_ORDINARY)
}

func (m *GmModule) pass1Times(input string, hero *entity.Hero, result herolock.LockResult, hc iface.HeroController) {
	m.passDungeon(hero, result, hc, 1, false, shared_proto.DifficultType_ORDINARY)
}

func (m *GmModule) pass1EliteChapter(input string, hero *entity.Hero, result herolock.LockResult, hc iface.HeroController) {
	m.passDungeon(hero, result, hc, 10, true, shared_proto.DifficultType_ELITE)
}

func (m *GmModule) pass1EliteTimes(input string, hero *entity.Hero, result herolock.LockResult, hc iface.HeroController) {
	m.passDungeon(hero, result, hc, 1, false, shared_proto.DifficultType_ELITE)
}

func (m *GmModule) passDungeon(hero *entity.Hero, result herolock.LockResult, hc iface.HeroController, totalPassTimes int, breakEndChapter bool, difficultType shared_proto.DifficultType) {
	passTimes := 0

	heroDungeon := hero.Dungeon()
	var firstChapterData *dungeon.DungeonChapterData
	for _, chapterData := range m.datas.GetDungeonChapterDataArray() {
		if chapterData.Type ==  difficultType {
			// 普通副本
			for _, dungeonData := range chapterData.DungeonDatas {
				if heroDungeon.IsPass(dungeonData) {
					continue
				}

				if hero.Level() < dungeonData.UnlockHeroLevel {
					continue
				}

				if firstChapterData == nil {
					firstChapterData = chapterData
				} else {
					if firstChapterData != chapterData {
						break
					}
				}

				if heroDungeon.IsUnlockPreDungeon(dungeonData) {
					logrus.Debugf("通关副本: %s", dungeonData.Name)

					// 解锁了

					enabledStars := []bool{}
					for i := uint64(0) ; i < dungeonData.Star; i++ {
						enabledStars = append(enabledStars, true)
					}
					heroDungeon.Pass(dungeonData, enabledStars) // 给你满星通关，开心不

					// 通关奖励
					prize := dungeonData.Plunder.Try()
					heromodule.AddPrize(m.hctx, hero, result, prize, m.time.CurrentTime())

					result.Add(dm.NewS2cChallengeMsg(u64.Int32(dungeonData.Id), "", []byte{}, must.Marshal(prize.Encode()),
					true, true, enabledStars, 10,
					u64.Int32(heroDungeon.GetChapterStar(dungeonData.GetChapterStarDungeonIds())),
					u64.Int32(heroDungeon.GetLimitDungeonPassTimes(dungeonData.Id)), true))

					passTimes++

					if passTimes >= totalPassTimes {
						hc.Disconnect(misc.ErrDisconectReasonFailGm)
						return
					}
				}
			}
		}
	}
}
