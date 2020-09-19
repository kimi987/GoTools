package random_event

import (
	"github.com/lightpaw/male7/gen/iface"
	"github.com/lightpaw/male7/gen/pb/random_event"
	"github.com/lightpaw/male7/entity"
	"github.com/lightpaw/male7/service/conflict/heroservice/herolock"
	"github.com/lightpaw/logrus"
	"github.com/lightpaw/male7/service/heromodule"
	"github.com/lightpaw/male7/service/operate_type"
	"github.com/lightpaw/male7/pb/shared_proto"
	"github.com/lightpaw/male7/util/must"
	"math/rand"
	"github.com/lightpaw/male7/util/u64"
)

func NewRandomEventModule(dep iface.ServiceDep, seasonService iface.SeasonService) *RandomEventModule {
	return &RandomEventModule{
		dep:    dep,
		datas:  dep.Datas(),
		season: seasonService,
	}
}

//gogen:iface
type RandomEventModule struct {
	dep    iface.ServiceDep
	datas  iface.ConfigDatas
	season iface.SeasonService
}

//gogen:iface
func (m *RandomEventModule) ProcessChooseOption(proto *random_event.C2SChooseOptionProto, hc iface.HeroController) {

	ctime := m.dep.Time().CurrentTime()
	posX := proto.PosX
	posY := proto.PosY
	option := proto.Option
	if hc.FuncWithSend(func(hero *entity.Hero, result herolock.LockResult) {
		event := hero.RandomEvent().GetEvent(posX, posY)
		if event == nil {
			logrus.Debug("随机事件选项，失效事件")
			result.Add(random_event.ERR_CHOOSE_OPTION_FAIL_INVALID_EVENT)
			return
		}

		eventId := event.Id()
		if eventId == 0 {
			logrus.Debug("随机事件选项，事件还未生成")
			result.Add(random_event.ERR_CHOOSE_OPTION_FAIL_NO_CATCH)
			return
		}
		eventData := m.datas.RandomEventData().Get(eventId)
		optionSize := len(eventData.OptionDatas)
		if option < 0 || option >= int32(optionSize) {
			logrus.Debug("随机事件选项，选项下标越界")
			result.Add(random_event.ERR_CHOOSE_OPTION_FAIL_INVALID_OPTION)
			return
		}
		var hctx *heromodule.HeroContext
		optionData := eventData.OptionDatas[option]
		if cost := optionData.Cost; cost != nil { // 摳消耗
			hctx = heromodule.NewContext(m.dep, operate_type.RandomEventChooseOption)
			if !heromodule.TryReduceCost(hctx, hero, result,cost) {
				logrus.Debugf("随机事件选项，選項消耗不足")
				result.Add(random_event.ERR_CHOOSE_OPTION_FAIL_COST_NOT_ENOUGH)
				return
			}
		}
		optionPrize := optionData.SuccessPrize
		bFailed := false
		if optionData.FailedRate > 0 {
			bFailed = rand.Uint64() % 10000 < optionData.FailedRate
		}
		if bFailed { // 失敗
			optionPrize = optionData.FailedPrize
		}
		if optionPrize != nil { // 發奬勵
			prize := optionPrize.CatchPrize()
			if optionPrize.GfAdd > 0 {
				if gf := hero.Domestic().GetBuilding(shared_proto.BuildingType_GUAN_FU); gf != nil {
					cp := *prize
					prize = &cp
					addValue := optionPrize.GfAdd * gf.Level
					if prize.HeroExp > 0 {
						prize.HeroExp += addValue
					}
					if prize.CaptainExp > 0 {
						prize.CaptainExp += addValue
					}
					if prize.Stone > 0 {
						if prize.UnsafeStone > 0 {
							prize.UnsafeStone += addValue
						}
						if prize.SafeStone > 0 {
							prize.SafeStone += addValue
						}
					}
					if prize.Gold > 0 {
						if prize.UnsafeGold > 0 {
							prize.UnsafeGold += addValue
						}

						if prize.SafeGold > 0 {
							prize.SafeGold += addValue
						}
					}
				}
			}
			if hctx == nil {
				hctx = heromodule.NewContext(m.dep, operate_type.RandomEventChooseOption)
			}
			heromodule.AddPrize(hctx, hero, result, prize, ctime)
			result.Add(random_event.NewS2cChooseOptionMsg(posX, posY, option, !bFailed, must.Marshal(prize.Encode())))
		} else {
			result.Add(random_event.NewS2cChooseOptionMsg(posX, posY, option, !bFailed, []byte{}))
		}

		hero.RandomEvent().RemoveEvent(posX, posY)

		result.Ok()
	}) {
		return
	}
}

//gogen:iface
func (m *RandomEventModule) ProcessOpenEvent(proto *random_event.C2SOpenEventProto, hc iface.HeroController) {

	posX := proto.PosX
	posY := proto.PosY
	if hc.FuncWithSend(func(hero *entity.Hero, result herolock.LockResult) {
		event := hero.RandomEvent().GetEvent(posX, posY)
		if event == nil {
			logrus.Debug("打开随机事件，失效事件")
			result.Add(random_event.ERR_OPEN_EVENT_FAIL_INVALID_EVENT)
			return
		}
		if eventId := event.Id(); eventId == 0 {
			typeArea := m.datas.RandomEventPositionDictionary().GetTypeArea(event.Cube())
			if typeArea == 0 {
				logrus.Debug("打开随机事件，配置表中数据有变更，删除")
				result.Add(random_event.ERR_OPEN_EVENT_FAIL_INVALID_EVENT)

				hero.RandomEvent().RemoveEvent(posX, posY)
				return
			}
			eventData := m.datas.RandomEventDataDictionary().CatchEventData(m.season.Season().Season, typeArea) // 必有数据
			// 记录，下次客户端申请直接可以获取
			event.SetId(eventData.Id)
			result.Add(random_event.NewS2cOpenEventMsg(posX, posY, u64.Int32(eventData.Id), eventData.OptionsProto4Send))
			// 是否激活新的图鉴
			if hero.RandomEvent().TrySetHandbooks(eventData.Id) {
				result.Add(random_event.NewS2cAddEventHandbookMsg(u64.Int32(eventData.Id)))

				heromodule.UpdateTaskProgress(hero, result, shared_proto.TaskTargetType_TASK_TARGET_RNDEVENT_HANDBOOKS)
			}

		} else {
			eventData := m.datas.RandomEventData().Get(eventId)
			if eventData == nil {
				logrus.Debug("打开随机事件，配置表中无数据，删除")
				result.Add(random_event.ERR_OPEN_EVENT_FAIL_INVALID_EVENT)

				hero.RandomEvent().RemoveEvent(posX, posY)
				return
			}
			result.Add(random_event.NewS2cOpenEventMsg(posX, posY, u64.Int32(eventId), eventData.OptionsProto4Send))
		}

		result.Ok()
	}) {
		return
	}
}
