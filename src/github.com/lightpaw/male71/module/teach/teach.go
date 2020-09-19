package teach

import (
	"github.com/lightpaw/male7/entity"
	"github.com/lightpaw/male7/gen/iface"
	"github.com/lightpaw/male7/gen/pb/teach"
	"github.com/lightpaw/male7/service/conflict/heroservice/herolock"
	"github.com/lightpaw/male7/service/heromodule"
	"github.com/lightpaw/male7/service/operate_type"
	"github.com/lightpaw/male7/util/u64"
)

func NewTeachModule(dep iface.ServiceDep) *TeachModule {
	m := &TeachModule{}
	m.dep = dep
	return m
}

//gogen:iface
type TeachModule struct {
	dep iface.ServiceDep
}

//gogen:iface
func (m *TeachModule) ProcessFight(proto *teach.C2SFightProto, hc iface.HeroController) {
	dataId := u64.FromInt32(proto.Id)
	data := m.dep.Datas().GetTeachChapterData(dataId)
	if data == nil {
		hc.Send(teach.ERR_FIGHT_FAIL_INVALID_ID)
		return
	}

	hc.FuncWithSend(func(hero *entity.Hero, result herolock.LockResult) {
		if data.PrevData != nil {
			if _, ok := hero.Teach().PassedChapterIds[data.PrevData.Id]; !ok {
				result.Add(teach.ERR_FIGHT_FAIL_PREV_CHAPTER_NOT_PASS)
				return
			}
		}

		if hero.Level() < data.MinHeroLevel {
			result.Add(teach.ERR_FIGHT_FAIL_HERO_LEVEL_LIMIT)
			return
		}

		if hero.TaskList().GetCompletedBaYeStage() < data.PassBaYeTaskStage {
			result.Add(teach.ERR_FIGHT_FAIL_BA_YE_TASK_LIMIT)
			return
		}

		if data.PassDungeonId > 0 {
			if !hero.Dungeon().IsPass(m.dep.Datas().GetDungeonData(data.PassDungeonId)) {
				result.Add(teach.ERR_FIGHT_FAIL_DUNGEON_LIMIT)
				return
			}
		}

		hero.Teach().PassedChapterIds[data.Id] = struct{}{}

		result.Add(teach.NewS2cFightMsg(proto.Id))

		result.Changed()
		result.Ok()
	})

}

//gogen:iface
func (m *TeachModule) ProcessCollectPrize(proto *teach.C2SCollectPrizeProto, hc iface.HeroController) {
	dataId := u64.FromInt32(proto.Id)
	data := m.dep.Datas().GetTeachChapterData(dataId)
	if data == nil {
		hc.Send(teach.ERR_COLLECT_PRIZE_FAIL_INVALID_ID)
		return
	}

	hc.FuncWithSend(func(hero *entity.Hero, result herolock.LockResult) {
		if _, ok := hero.Teach().PassedChapterIds[dataId]; !ok {
			hc.Send(teach.ERR_COLLECT_PRIZE_FAIL_NOT_PASS)
			return
		}
		if _, ok := hero.Teach().CollectedChapterIds[dataId]; ok {
			hc.Send(teach.ERR_COLLECT_PRIZE_FAIL_ALREADY_COLLECTED)
			return
		}

		hero.Teach().CollectedChapterIds[dataId] = struct{}{}
		ctime := m.dep.Time().CurrentTime()
		hctx := heromodule.NewContext(m.dep, operate_type.TeachCollectPrize)
		heromodule.AddPrize(hctx, hero, result, data.Prize, ctime)

		result.Add(teach.NewS2cCollectPrizeMsg(proto.Id))

		result.Changed()
		result.Ok()
	})
}
