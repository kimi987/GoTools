package bai_zhan

import (
	"github.com/lightpaw/logrus"
	"github.com/lightpaw/male7/entity"
	"github.com/lightpaw/male7/gen/pb/bai_zhan"
	"github.com/lightpaw/male7/module/bai_zhan/bai_zhan_objs"
	"github.com/lightpaw/male7/module/rank/rankface"
	"github.com/lightpaw/male7/pb/server_proto"
	"github.com/lightpaw/male7/pb/shared_proto"
	"github.com/lightpaw/male7/service/conflict/heroservice/herolock"
	"github.com/lightpaw/male7/service/heromodule"
	"time"
	"github.com/lightpaw/male7/config/maildata"
)

func (m *BaiZhanModule) GmResetDaily() {
	m.tryResetDaily(m.lastResetTime)
}

func (m *BaiZhanModule) GmResetChallengeTimes(heroId int64) {
	m.baiZhanService.Func(func(objs *bai_zhan_objs.BaiZhanObjs) {
		if obj := objs.GetBaiZhanObj(heroId); obj != nil {
			obj.ResetChallengeTimes()

			m.worldService.Send(heroId, bai_zhan.RESET_S2C)
		}
	})
}

// 每日重置
func (m *BaiZhanModule) tryResetDaily(newResetTime time.Time) {
	m.baiZhanService.Func(func(objs *bai_zhan_objs.BaiZhanObjs) {
		// 更新排名，上升，下降，如果军衔变更发送军衔变更，重置玩家数据(今日积分，俸禄)
		validCount := uint64(0)

		for _, data := range m.baiZhanLevelDatas {
			// 重置军衔等级
			validCount += m.resetJunXianLevel(data, newResetTime)
			// 重置镜像
			data.clearMirrors()
		}

		var rankObjs = make([]rankface.RankObj, 0, validCount)

		objs.Walk(func(obj *bai_zhan_objs.HeroBaiZhanObj) {
			if obj.Point() <= 0 {
				// 0分的，且军衔没被清除的，清除军衔
				// 所有军衔0分者，从军衔内镜像匹配池里移除
				oldLevelData := obj.LevelData()
				obj.RemoveJunXian()
				if oldLevelData != obj.LevelData() {
					m.onJunXianLevelChanged(obj, false)
					m.notifyLevelChanged(obj, shared_proto.LevelChangeType_LEVEL_DOWN)
				} else {
					m.notifyLevelChanged(obj, shared_proto.LevelChangeType_LEVEL_KEEP)
				}
			} else {
				// 有积分，但是没上榜的人都是保级
				rankObjs = append(rankObjs, m.newRankObj(obj))

				obj.ClearPoint()

				// 重新生成镜像，0分的肯定没镜像
				if obj.CombatMirror() != nil {
					m.mustBaiZhanLevelData(obj.LevelData()).addMirror(obj)
				}
			}

			// 每日重置
			obj.ResetDaily()
		})

		// 设置最近一次重置时间
		m.lastResetTime = newResetTime

		// 新的百战排行榜
		m.rankModule.UpdateBaiZhanRankList(rankObjs)

		// 重新生成积分排行榜
		m.resetPointRankList()

		// 广播数据重置
		m.worldService.Broadcast(bai_zhan.RESET_S2C)
	})
}

// 重置军衔等级
func (m *BaiZhanModule) resetJunXianLevel(data *bai_zhan_level_data, newResetTime time.Time) (rankCount uint64) {
	//结算时，记某一军衔内的玩家数量为S
	//	除诸侯以外的军衔，累计得分排名排在前INT(S*x%)+1且得分≥X的玩家，成为升级者
	//	除小卒以外的军衔，累计得分排名排在后INT(S*y%)+1且得分≤Y的玩家，成为降级者
	//	所有军衔0分者，从军衔内镜像匹配池里移除
	//	上述以外的玩家，成为保级者

	levelData := data.levelData

	rankList := m.mustPointRankList(levelData)

	rankList.Lock()
	defer rankList.Unlock()

	rankCount = rankList.RankCount()

	levelUpMaxRank, _, levelDownMinRank, _ := rankList.LevelUpAndDownRankAndPoint()

	// 榜单上面的都是有分的
	rankList.walk(func(obj *bai_zhan_objs.HeroBaiZhanObj) {
		switch obj.LevelChangeType(levelUpMaxRank, levelDownMinRank) {
		case shared_proto.LevelChangeType_LEVEL_UP:
			// 升级
			if levelData.NextLevel != nil {
				levelChange, historyLevelChanged := obj.ResetLevelData(levelData.NextLevel, newResetTime)
				if levelChange {
					m.onJunXianLevelChanged(obj, historyLevelChanged)
					m.notifyLevelChanged(obj, shared_proto.LevelChangeType_LEVEL_UP)

					m.dep.Tlog().TlogBaiZhanFlowById(obj.Id(), levelData.Level, levelData.NextLevel.Level)

				} else {
					m.notifyLevelChanged(obj, shared_proto.LevelChangeType_LEVEL_KEEP)
				}
				return
			}

			logrus.Errorln("能够升级的竟然没有下一级")
		case shared_proto.LevelChangeType_LEVEL_DOWN:
			// 降级
			if levelData.PrevLevel != nil {
				levelChange, historyLevelChanged := obj.ResetLevelData(levelData.PrevLevel, newResetTime)
				if levelChange {
					m.onJunXianLevelChanged(obj, historyLevelChanged)
					m.notifyLevelChanged(obj, shared_proto.LevelChangeType_LEVEL_DOWN)
					
					m.dep.Tlog().TlogBaiZhanFlowById(obj.Id(), levelData.Level, levelData.NextLevel.Level)
				} else {
					m.notifyLevelChanged(obj, shared_proto.LevelChangeType_LEVEL_KEEP)
				}
				return
			}

			logrus.Errorln("能够降级的竟然没有上一级")
		}

		// 保级
		m.notifyLevelChanged(obj, shared_proto.LevelChangeType_LEVEL_KEEP)
	})

	return
}

func (m *BaiZhanModule) onJunXianLevelChanged(obj *bai_zhan_objs.HeroBaiZhanObj, historyLevelChanged bool) {

	m.heroDataService.FuncWithSend(obj.Id(), func(hero *entity.Hero, result herolock.LockResult) {
		hero.HistoryAmount().Set(server_proto.HistoryAmountType_MaxJunXianLevel, obj.HistoryMaxJunXianLevelData().Level)
		heromodule.UpdateTaskProgress(hero, result, shared_proto.TaskTargetType_TASK_TARGET_BAI_ZHAN_JUN_XIAN)

		// 历史最高记录改变
		if historyLevelChanged {
			result.Add(obj.HistoryMaxJunXianLevelData().MaxJunXianLevelChangedMsg)
		}

		result.Changed()
		result.Ok()
	})

}

func (m *BaiZhanModule) notifyLevelChanged(obj *bai_zhan_objs.HeroBaiZhanObj, changeType shared_proto.LevelChangeType) {
	if obj.CombatMirrorFightAmount() <= 0 {
		return
	}

	var mailData *maildata.MailData

	switch changeType {
	case shared_proto.LevelChangeType_LEVEL_UP:
		mailData = m.configDatas.MailHelp().BaiZhanJunXianLevelUp
	case shared_proto.LevelChangeType_LEVEL_KEEP:
		if obj.LevelData().NextLevel == nil {
			mailData = m.configDatas.MailHelp().BaiZhanJunXianLevelMaxKeep
		} else {
			mailData = m.configDatas.MailHelp().BaiZhanJunXianLevelKeep
		}
	default:
		mailData = m.configDatas.MailHelp().BaiZhanJunXianLevelDown
	}

	// 发邮件
	proto := mailData.NewTextMail(shared_proto.MailType_MailNormal)
	proto.Text = mailData.Text.New().WithKey("jun_xian", obj.LevelData().Name).JsonString()
	m.mailModule.SendProtoMail(obj.Id(), proto, m.timeService.CurrentTime())
}
