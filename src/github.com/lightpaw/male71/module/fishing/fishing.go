package fishing

import (
	"github.com/lightpaw/logrus"
	"github.com/lightpaw/male7/config/resdata"
	"github.com/lightpaw/male7/entity"
	"github.com/lightpaw/male7/gen/iface"
	"github.com/lightpaw/male7/gen/pb/fishing"
	"github.com/lightpaw/male7/pb/server_proto"
	"github.com/lightpaw/male7/pb/shared_proto"
	"github.com/lightpaw/male7/service/conflict/heroservice/herolock"
	"github.com/lightpaw/male7/service/heromodule"
	"github.com/lightpaw/male7/util/u64"
	"github.com/lightpaw/male7/config/fishing_data"
	"github.com/lightpaw/male7/util/i64"
	"github.com/lightpaw/male7/gen/pb/misc"
	"fmt"
	"github.com/lightpaw/male7/service/operate_type"
	"github.com/lightpaw/male7/config/goods"
	"math/rand"
)

// 钓鱼
func NewFishingModule(dep iface.ServiceDep, guildSnapshotService iface.GuildSnapshotService) *FishingModule {
	m := &FishingModule{}
	m.dep = dep
	m.datas = dep.Datas()
	m.timeService = dep.Time()
	m.chatService = dep.Chat()
	m.guildSnapshotService = guildSnapshotService
	return m
}

//gogen:iface
type FishingModule struct {
	dep                  iface.ServiceDep
	datas                iface.ConfigDatas
	timeService          iface.TimeService
	chatService          iface.ChatService
	guildSnapshotService iface.GuildSnapshotService
}

// 钓鱼
//gogen:iface
func (m *FishingModule) ProcessFishing(proto *fishing.C2SFishingProto, hc iface.HeroController) {

	fishType := uint64(proto.FishType)
	if fishType != fishing_data.FishTypeYuanbao && fishType != fishing_data.FishTypeFree {
		logrus.WithField("fish_type", proto.FishType).Debugf("钓鱼，无效的钓鱼类型")
		hc.Send(fishing.ERR_FISHING_FAIL_INVALID_FISH_TYPE)
		return
	}

	var fishGoods *goods.GoodsData
	if proto.UseGoods {
		if fishType != fishing_data.FishTypeYuanbao {
			logrus.WithField("fish_type", proto.FishType).Debugf("钓鱼，鱼饵只能代替元宝钓鱼")
			hc.Send(fishing.ERR_FISHING_FAIL_INVALID_FISH_TYPE)
			return
		}

		fishGoods = m.datas.GoodsConfig().FishGoods
		if fishGoods == nil {
			logrus.WithField("fish_type", proto.FishType).Debugf("钓鱼，鱼饵物品不存在")
			hc.Send(fishing.ERR_FISHING_FAIL_INVALID_FISH_TYPE)
			return
		}
	}

	// 获得消耗
	fishId := fishing_data.FishId(fishType, u64.FromInt32(proto.GetTimes()))
	fishingCostData := m.datas.GetFishingCostData(fishId)
	if fishingCostData == nil {
		logrus.WithFields(logrus.Fields{"times": proto.GetTimes()}).Debugf("消耗没找到[%s]\n", proto.GetTimes())
		hc.Send(fishing.ERR_FISHING_FAIL_COST_NOUT_FOUND)
		return
	}

	hctx := heromodule.NewContext(m.dep, operate_type.FishingFishing)

	//var broadcastMsgs []pbutil.Buffer
	if hc.FuncWithSend(func(hero *entity.Hero, result herolock.LockResult) {
		// 检查消耗够不够
		heroFishing := hero.Fishing()
		fishingTimes := heroFishing.FishingTimes(fishId) // 获得该种类型钓鱼的次数
		if fishingCostData.DailyTimes > 0 && fishingTimes >= fishingCostData.DailyTimes {
			logrus.WithField("fish_type", proto.FishType).Debugf("钓鱼，钓鱼次数达到当日上限")
			hc.Send(fishing.ERR_FISHING_FAIL_DAILY_TIMES_LIMIT)
			return
		}

		newFishingTimes := fishingTimes

		ctime := m.timeService.CurrentTime()
		var cost *resdata.Cost
		if proto.UseGoods {
			// 使用物品
			if !heromodule.TryReduceGoods(hctx, hero, result, fishGoods, fishingCostData.Times) {
				logrus.Debugf("钓鱼[%v]次的鱼饵不够啊\n", proto.GetTimes())
				result.Add(fishing.ERR_FISHING_FAIL_RES_NOT_ENOUGH)
				return
			}
		} else {
			if fishingCostData.FreeCountdown > 0 {
				nextTime := heroFishing.FishingCountdown(fishId)
				if nextTime > 0 && ctime.Unix() < nextTime {
					logrus.Debugf("钓鱼，倒计时未结束")
					hc.Send(fishing.ERR_FISHING_FAIL_COUNTDOWN)
					return
				}
			}

			cost = fishingCostData.GetCost(fishingTimes) // 获得该种类型的消耗，可能为空，因为免费
			if cost != nil && !heromodule.TryReduceCost(hctx, hero, result, cost) {
				logrus.WithFields(logrus.Fields{"fishingCostData": fishingCostData}).Debugf("钓鱼[%v]次的钱不够啊\n", proto.GetTimes())
				result.Add(fishing.ERR_FISHING_FAIL_RES_NOT_ENOUGH)
				return
			}

			// 使用花钱时候才自增
			newFishingTimes = heroFishing.IncreFishingTimes(fishId)
		}

		//guildFlagName := ""
		//if g := m.guildSnapshotService.GetSnapshot(hero.GuildId()); g != nil {
		//	guildFlagName = g.FlagName
		//}

		var priorityIndex int32 = -1
		if proto.GetTimes() > 1 {
			if proto.FishType == fishing_data.FishTypeYuanbao {
				// 策划新需求（比如单抽了2次，此时10连会在第8个道具位置显示出保底道具）
				priorityIndex = fishing_data.FishCombo - u64.Int32(hero.MiscData().GetFishCombo()) - 1
			} else {
				priorityIndex = rand.Int31n(proto.GetTimes())
			}
		} else if proto.FishType == fishing_data.FishTypeYuanbao && hero.MiscData().IncFishCombo() >= fishing_data.FishCombo {
			hero.MiscData().SetFishCombo(0)
			priorityIndex = 0
		}

		// 随机奖励
		prizeBuilder := resdata.NewPrizeBuilder()

		prizeBytes := make([][]byte, 0, proto.GetTimes())
		haveSoulToGoodsArray := make([]bool, 0, proto.GetTimes())
		showIndexes := make([]int32, 0, 1)

		fishTimes := hero.HistoryAmount().AmountWithSubType(server_proto.HistoryAmountType_Fishing, fishType)
		var fishDatas []*fishing_data.FishData
		for index := int32(0); index < proto.GetTimes(); index++ {
			fishTimes++

			var data *fishing_data.FishData
			if index == priorityIndex {
				data = m.datas.FishRandomer().PriorityFishing(fishType, heroFishing.CaptainSet())
				if data == nil {
					data = m.datas.FishRandomer().Fishing(fishType, fishTimes, heroFishing.CaptainSet())
				}
			} else {
				data = m.datas.FishRandomer().Fishing(fishType, fishTimes, heroFishing.CaptainSet())
			}

			// 有判断的，奖励不可能为空
			prizeBuilder.Add(data.Prize)

			haveSoulToGoods := false
		out:
			for _, captain := range data.Prize.Captain {
				if hero.Military().Captain(captain.Id) != nil {
					haveSoulToGoods = true
					break
				}

				for _, fd := range fishDatas {
					for _, so := range fd.Prize.Captain {
						if so.Id == captain.Id {
							haveSoulToGoods = true
							break out
						}
					}
				}
			}
			fishDatas = append(fishDatas, data)

			prizeBytes = append(prizeBytes, data.GetPrizeBytes())
			haveSoulToGoodsArray = append(haveSoulToGoodsArray, haveSoulToGoods)

			if data.IsShow {
				showIndexes = append(showIndexes, index+1)
			}

			if len(data.Prize.Gem) > 0 {
				for _, v := range data.Prize.Gem {
					// 系统广播
					heromodule.AddGemBroadcast(hctx, hero, v, result)
				}
			}
		}

		var newFishingCountdown int64
		if !proto.UseGoods {
			if fishingCostData.FreeCountdown > 0 && cost == nil {
				newFishingCountdown = ctime.Add(fishingCostData.FreeCountdown).Unix()
				heroFishing.SetFishingCountdown(fishId, newFishingCountdown)
			}
		}

		if fishingCostData.FishType == fishing_data.FishTypeYuanbao {
			newAmount := hero.MiscData().GetFishPoint() + fishingCostData.Times
			hero.MiscData().SetFishPoint(newAmount)
			result.Add(fishing.NewS2cUpdateFishPointMsg(u64.Int32(newAmount)))
		}

		// 发消息
		result.Add(fishing.NewS2cFishingMsg(prizeBytes, haveSoulToGoodsArray, showIndexes, proto.GetTimes(), proto.FishType, u64.Int32(newFishingTimes), i64.Int32(newFishingCountdown)))

		if fishType == fishing_data.FishTypeYuanbao {
			if hero.Bools().TrySet(shared_proto.HeroBoolType_BOOL_YUANBAO_FISH) {
				result.Add(misc.NewS2cSetHeroBoolMsg(int32(shared_proto.HeroBoolType_BOOL_YUANBAO_FISH)))
			}
		} else if fishType == fishing_data.FishTypeFree {
			if hero.Bools().TrySet(shared_proto.HeroBoolType_BOOL_FREE_FISH) {
				result.Add(misc.NewS2cSetHeroBoolMsg(int32(shared_proto.HeroBoolType_BOOL_FREE_FISH)))
			}
		}

		// 给奖励
		heromodule.AddPrize(hctx, hero, result, prizeBuilder.Build(), m.timeService.CurrentTime())

		hero.HistoryAmount().IncreaseWithSubType(server_proto.HistoryAmountType_Fishing, fishType, fishingCostData.Times)
		heromodule.UpdateTaskProgress(hero, result, shared_proto.TaskTargetType_TASK_TARGET_ACCUM_FISHING)

		heromodule.IncreTaskProgress(hero, result, shared_proto.TaskTargetType_TASK_TARGET_FISHING, fishingCostData.Times)

		result.Changed()
		result.Ok()

		var tlogFishType uint64
		if fishType == fishing_data.FishTypeYuanbao {
			if proto.UseGoods {
				tlogFishType = operate_type.FishUseGoods
			} else {
				tlogFishType = operate_type.FishYuanbao
				if fishingCostData.Times > 1 {
					tlogFishType = operate_type.FishYuanbaoTen
				}
			}
		} else if fishType == fishing_data.FishTypeFree {
			tlogFishType = operate_type.FishNormal
		}

		// tlog
		hctx.Tlog().TlogFishFlow(hero, tlogFishType, u64.FromInt(fishingCostData.Cost.Id))

	}) {
		return
	}
}

// 钓鱼兑换
//gogen:iface c2s_fish_point_exchange
func (m *FishingModule) ProcessFishPointExchange(hc iface.HeroController) {

	hc.FuncWithSend(func(hero *entity.Hero, result herolock.LockResult) {

		fishPointExchangeIndex := hero.MiscData().GetFishPointExchangeIndex()
		configMaxIndex := len(m.datas.MiscGenConfig().FishMaxPoint) - 1
		if fishPointExchangeIndex > int32(configMaxIndex) {
			fishPointExchangeIndex = int32(configMaxIndex)
		}
		if hero.MiscData().GetFishPoint() < m.datas.MiscGenConfig().FishMaxPoint[fishPointExchangeIndex] {
			logrus.Debug("钓鱼兑换，积分不足")
			result.Add(fishing.ERR_FISH_POINT_EXCHANGE_FAIL_POINT_NOT_ENOUGH)
			return
		}

		// 可以兑换的将魂
		//var array []*resdata.ResCaptainData
		//for _, cs := range m.datas.MiscGenConfig().FishPointCaptain {
		//	if cs != nil {
		//		if hero.Military().Captain(cs.Id) == nil {
		//			array = append(array, cs)
		//		}
		//	}
		//}
		//
		//if len(array) <= 0 {
		//	logrus.Debug("钓鱼兑换，已拥有所有的兑换物")
		//	result.Add(fishing.ERR_FISH_POINT_EXCHANGE_FAIL_OWNER_ALL)
		//	return
		//}

		// 设置新的fishPointExchangeIndex下标
		//if int(fishPointExchangeIndex) + 1 < len(m.datas.MiscGenConfig().FishMaxPoint) {
		hero.MiscData().IncFishPointExchangeIndex()
		//}

		// 扣积分
		toSet := u64.Sub(hero.MiscData().GetFishPoint(), m.datas.MiscGenConfig().FishMaxPoint[fishPointExchangeIndex])
		hero.MiscData().SetFishPoint(toSet)
		result.Add(fishing.NewS2cUpdateFishPointMsg(u64.Int32(toSet)))

		// 获得的奖励
		prize := m.datas.MiscGenConfig().FishPointPlunder.Try()

		existCaptain := false
		for _, c := range prize.Captain {
			if hero.Military().Captain(c.Id) != nil {
				existCaptain = true
			}
		}

		ctime := m.timeService.CurrentTime()
		hctx := heromodule.NewContext(m.dep, operate_type.FishingFishPointExchange)
		heromodule.AddPrize(hctx, hero, result, prize, ctime)

		if fishPointExchangeIndex < int32(configMaxIndex) {
			fishPointExchangeIndex++
		}
		result.Add(fishing.NewS2cFishPointExchangeMarshalMsg(prize.Encode(), fishPointExchangeIndex, existCaptain))

		result.Ok()
	})

}

func (m *FishingModule) GmFishingRate(heroId int64, fishType, times uint64) {

	// GM 测试钓鱼概率
	var setId uint64
	m.dep.HeroData().Func(heroId, func(hero *entity.Hero, err error) (heroChanged bool) {
		setId = hero.Fishing().CaptainSet()
		return
	})

	// 随机奖励
	prizeBuilder := resdata.NewPrizeBuilder()

	fishDataMap := make(map[uint64]uint64)
	for index := uint64(0); index < times; index++ {

		// 里面不会返回空
		data := m.datas.FishRandomer().Fishing(fishType, 1000, setId)

		// 有判断的，奖励不可能为空
		prizeBuilder.Add(data.Prize)

		fishDataMap[data.Id]++
	}

	text := fmt.Sprintf("钓鱼 %+v", prizeBuilder.Build().Encode())
	m.chatService.SysChat(0, heroId, shared_proto.ChatType_ChatWorld, text, shared_proto.ChatMsgType_ChatMsgSys, true, false, true, false)

}

// 金杆钓设置
//gogen:iface
func (m *FishingModule) ProcessSetFishingCaptain(proto *fishing.C2SSetFishingCaptainProto, hc iface.HeroController) {
	data := m.datas.GetFishingCaptainProbabilityData(u64.FromInt32(proto.GetCaptainId()))
	if data == nil {
		logrus.Debug("金杆钓设置，无效的配置id")
		hc.Send(fishing.ERR_SET_FISHING_CAPTAIN_FAIL_INVALID_ID)
		return
	}
	hc.FuncWithSend(func(hero *entity.Hero, result herolock.LockResult) {
		vipData := m.datas.GetVipLevelData(hero.VipLevel())
		if vipData != nil && !vipData.FishingCaptainProbability {
			logrus.Debug("金杆钓设置，vip等级不足")
			result.Add(fishing.ERR_SET_FISHING_CAPTAIN_FAIL_VIP_LEVEL_LIMIT)
			return
		}
		if hero.Fishing().CaptainSet() == data.CaptainId {
			logrus.Debug("金杆钓设置，重复设置")
			result.Add(fishing.ERR_SET_FISHING_CAPTAIN_FAIL_DUPLICATE)
			return
		}
		hero.Fishing().ChangeCaptainSet(data.CaptainId)
		result.Add(fishing.NewS2cSetFishingCaptainMsg(proto.CaptainId))
		result.Ok()
	})
}
