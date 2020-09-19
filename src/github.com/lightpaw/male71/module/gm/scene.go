package gm

import (
	"context"
	"github.com/lightpaw/logrus"
	"github.com/lightpaw/male7/entity"
	"github.com/lightpaw/male7/gen/iface"
	"github.com/lightpaw/male7/gen/pb/military"
	"github.com/lightpaw/male7/gen/pb/misc"
	"github.com/lightpaw/male7/gen/pb/region"
	"github.com/lightpaw/male7/service/conflict/guilddataservice/sharedguilddata"
	"github.com/lightpaw/male7/service/conflict/heroservice/herolock"
	"github.com/lightpaw/male7/util/u64"
	"sort"
	"github.com/lightpaw/male7/util/ctxfunc"
	"github.com/lightpaw/male7/pb/shared_proto"
)

func (m *GmModule) newSceneGmGroup() *gm_group {
	return &gm_group{
		tab: "场景",
		handler: []*gm_handler{
			newStringHandler("自己流亡", "", m.removeBase),
			newIntHandler("加/减繁荣度(不会改等级)", "1000", m.addProsperity),
			newIntHandler("群嘲", "0", m.chaoFeng),
			newIntHandler("护驾", "0", m.huJia),
			//newHeroStringHandler("创建帮派", "帮派名字 旗号", func(amount string, hc iface.HeroController) {
			//	strArray := strings.Split(strings.TrimSpace(amount), " ")
			//
			//	modules.GuildModule().(*guild2.GuildModule).ProcessCreateGuild(&guild.C2SCreateGuildProto{
			//		Name:     strArray[0],
			//		FlagName: strArray[1],
			//	}, hc)
			//}),
			//newHeroStringHandler("打印帮派", "", func(amount string, hc iface.HeroController) {
			//	ids, _ := db.ListGuild(0, 10)
			//	for _, id := range ids {
			//		g, _ := db.LoadGuild(id)
			//		fmt.Println(g.EncodeClient(false))
			//		fmt.Println(g.EncodeClient(true))
			//	}
			//}),
		},
	}
}

func (m *GmModule) removeBase(amount string, hc iface.HeroController) {

	realm := m.realmService.GetBigMap()
	if realm != nil {
		realm.RemoveBase(hc.Id(), false, m.datas.TextHelp().MRDRBroken4a.Text, m.datas.TextHelp().MRDRBroken4d.Text)
	}

	hc.FuncWithSend(func(hero *entity.Hero, result herolock.LockResult) {
		hero.SetProsperity(0)
		hero.SetBaseLevel(0)

		result.Changed()
		result.Ok()
	})

	hc.Disconnect(misc.ErrDisconectReasonFailGm)
}

func (m *GmModule) addProsperity(amount int64, hc iface.HeroController) {
	var maxCanReduceHomeProsperity uint64
	hc.FuncWithSend(func(hero *entity.Hero, result herolock.LockResult) {

		maxCanReduceHomeProsperity = u64.Sub(hero.Prosperity(), 1) // 还要剩余1点
	})

	realm := m.realmService.GetBigMap()
	if realm != nil {
		if amount > 0 {
			realm.AddProsperity(hc.Id(), uint64(amount))
		} else {
			realm.GmReduceProsperity(hc.Id(), u64.Min(maxCanReduceHomeProsperity, uint64(-amount)))
		}
	}
}

func (m *GmModule) huJia(amount int64, hc iface.HeroController) {
	var selfGuild, selfRegion int64
	if hc.FuncNotError(func(hero *entity.Hero) (heroChanged bool) {
		selfGuild = hero.GuildId()
		selfRegion = hero.BaseRegion()
		return
	}) {
		return
	}

	if selfGuild == 0 {
		logrus.Debugf("没有帮派，护驾失败")
		return
	}

	var heroIds []int64
	if !m.sharedGuildService.Func(func(guilds sharedguilddata.Guilds) {
		g := guilds.Read(selfGuild)
		if g != nil {
			heroIds = g.AllUserMemberIds()
		}
	}) {
		return
	}

	c := amount
	for _, targetId := range heroIds {
		if targetId == hc.Id() {
			continue
		}

		if amount > 0 && c <= 0 {
			break
		}

		var proto *region.C2SInvasionProto
		m.heroDataService.FuncWithSend(targetId, func(hero *entity.Hero, result herolock.LockResult) {

			if hero.BaseRegion() != selfRegion {
				return
			}

			if hero.BaseLevel() <= 0 {
				return
			}

			setHeroTroopCaptain(hero, result, m.datas)

			var troopIndex int = -1
			for i, t := range hero.Troops() {
				if t.IsOutside() {
					continue
				}

				// 如果武将个数不足，给他搞几个武将
				for _, pos := range t.Pos() {
					captain := pos.Captain()
					if captain != nil {
						// 补满兵
						captain.SetSoldier(captain.SoldierCapcity())
						result.Add(military.NewS2cCaptainChangeDataMsg(u64.Int32(captain.Id()), u64.Int32(captain.Soldier()), u64.Int32(captain.SoldierCapcity()), u64.Int32(captain.FightAmount())))
					}
				}

				troopIndex = i
				break
			}

			if troopIndex >= 0 {
				// 搞点武将
				proto = &region.C2SInvasionProto{}
				proto.Operate = int32(shared_proto.TroopOperate_ToAssist)
				proto.Target = hc.IdBytes()
				proto.TroopIndex = int32(troopIndex)
			}

			result.Changed()
			result.Ok()
		})

		// 出发
		if proto != nil {
			m.modules.RegionModule().(interface {
				ProcessInvasion(*region.C2SInvasionProto, iface.HeroController)
			}).ProcessInvasion(proto, m.getOrCreateFakeHeroControler(targetId))

			c--
		}
	}
}

func (m *GmModule) chaoFeng(amount int64, hc iface.HeroController) {
	var selfGuild, selfRegion int64
	var selfBaseX, selfBaseY int
	if hc.FuncNotError(func(hero *entity.Hero) (heroChanged bool) {
		selfGuild = hero.GuildId()
		selfRegion = hero.BaseRegion()
		selfBaseX, selfBaseY = hero.BaseX(), hero.BaseY()
		return
	}) {
		return
	}

	// 查询本地图的所有玩家，这些家伙每个人搞一个部队来干我
	var heros []*entity.Hero
	err := ctxfunc.Timeout2s(func(ctx context.Context) (err error) {
		heros, err = m.db.LoadAllHeroData(ctx)
		return
	})
	if err != nil {
		logrus.WithError(err).Debugf("群嘲失败")
		return
	}

	// 根据距离我远近来排序近的排前面
	a := &hero_near_slice{
		baseX: selfBaseX,
		baseY: selfBaseY,
		a:     heros,
	}
	sort.Sort(a)

	var heroIds []int64
	for _, hero := range heros {
		if hero.Id() == hc.Id() {
			continue
		}

		heroIds = append(heroIds, hero.Id())
	}

	c := amount
	for _, targetId := range heroIds {
		if amount > 0 && c <= 0 {
			break
		}

		var proto *region.C2SInvasionProto
		m.heroDataService.FuncWithSend(targetId, func(hero *entity.Hero, result herolock.LockResult) {

			if hero.BaseRegion() != selfRegion {
				return
			}

			if hero.BaseLevel() <= 0 {
				return
			}

			if selfGuild != 0 && hero.GuildId() == selfGuild {
				// 自己人
				return
			}

			setHeroTroopCaptain(hero, result, m.datas)

			var troopIndex int = -1
			for i, t := range hero.Troops() {
				if t.IsOutside() {
					continue
				}

				// 如果武将个数不足，给他搞几个武将
				for _, pos := range t.Pos() {
					captain := pos.Captain()
					if captain != nil {
						// 补满兵
						captain.SetSoldier(captain.SoldierCapcity())
						result.Add(military.NewS2cCaptainChangeDataMsg(u64.Int32(captain.Id()), u64.Int32(captain.Soldier()), u64.Int32(captain.SoldierCapcity()), u64.Int32(captain.FightAmount())))
					}
				}

				troopIndex = i
				break
			}

			if troopIndex >= 0 {
				proto = &region.C2SInvasionProto{}
				proto.Operate = int32(shared_proto.TroopOperate_ToInvasion)
				proto.Target = hc.IdBytes()
				proto.TroopIndex = int32(troopIndex)
			}

			result.Changed()
			result.Ok()
		})

		// 出发
		if proto != nil {
			m.modules.RegionModule().(interface {
				ProcessInvasion(*region.C2SInvasionProto, iface.HeroController)
			}).ProcessInvasion(proto, m.getOrCreateFakeHeroControler(targetId))

			c--
		}
	}
}

func setHeroTroopCaptain(hero *entity.Hero, result herolock.LockResult, datas iface.ConfigDatas) {

	//ctime := time.Now()
	//
	//// 解锁所有的武将
	//for _, data := range datas.GetCaptainDataArray() {
	//	heromodule.TryAddCaptain(hero, result, data, ctime)
	//}
	//
	//for _, t := range hero.Troops() {
	//	if t.IsOutside() {
	//		continue
	//	}
	//
	//	for i, pos := range t.Pos() {
	//		captain := pos.Captain()
	//		if captain != nil {
	//			continue
	//		}
	//
	//		for _, captain := range hero.Military().Captains() {
	//			if captain.IsOutSide() {
	//				continue
	//			}
	//
	//			captain.SetSoldier(captain.SoldierCapcity())
	//
	//			t.SetCaptainIfAbsent(uint64(i), captain, 0)
	//
	//			sendIndex := t.Sequence()*uint64(len(t.Pos())) + uint64(i)
	//			result.Add(military.NewS2cRecruitCaptainV2MarshalMsg(captain.EncodeClient(), u64.Int32(sendIndex+1)))
	//
	//			heromodule.UpdateAllTroopFightAmount(hero, result)
	//
	//			break
	//		}
	//	}
	//}
	//
	//// 更新任务进度
	//heromodule.UpdateTaskProgress(hero, result, shared_proto.TaskTargetType_TASK_TARGET_CAPTAIN_COUNT)
	//heromodule.UpdateTaskProgress(hero, result, shared_proto.TaskTargetType_TASK_TARGET_CAPTAIN_LEVEL_COUNT)

}
