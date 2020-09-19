package realm

import (
	"fmt"
	"github.com/lightpaw/male7/config"
	"github.com/lightpaw/male7/config/blockdata"
	"github.com/lightpaw/male7/config/confpath"
	"github.com/lightpaw/male7/entity"
	"github.com/lightpaw/male7/entity/cb"
	npcid2 "github.com/lightpaw/male7/entity/npcid"
	"github.com/lightpaw/male7/gen/ifacemock"
	"github.com/lightpaw/male7/mock"
	"github.com/lightpaw/male7/module/realm/realmerr"
	"github.com/lightpaw/male7/module/realm/realmface"
	"github.com/lightpaw/male7/util/event"
	"github.com/lightpaw/male7/util/timeutil"
	"github.com/lightpaw/male7/util/u64"
	. "github.com/onsi/gomega"
	"testing"
	"time"
	"github.com/lightpaw/male7/pb/shared_proto"
	"github.com/lightpaw/male7/pb/server_proto"
	"github.com/lightpaw/male7/config/season"
	"github.com/lightpaw/male7/config/scene"
	"github.com/lightpaw/male7/service/conflict/guilddataservice/guildsnapshotdata"
	"github.com/lightpaw/male7/util/i64"
)

func TestRealm_RandomAroundBase(t *testing.T) {
	RegisterTestingT(t)
	r, _ := newMockRealm()

	for i := 0; i < 100; i++ {
		ok, x, y := r.randomBasePos()
		Ω(ok).Should(BeTrue())
		Ω(r.mapData.IsValidHomePosition(x, y)).Should(BeTrue())

		x, y, ok = r.RandomAroundBase(x, y)
		Ω(ok).Should(BeTrue())
		Ω(r.mapData.IsValidHomePosition(x, y)).Should(BeTrue())

		cube := cb.XYCube(x, y)

		Ω(r.conflict.baseConflictCount[cube]).Should(Equal(1))
		neighbors := NeighborsInBaseConflictRange(x, y)
		for _, n := range neighbors {
			Ω(r.conflict.baseConflictCount[n] > 0).Should(BeTrue())
		}
	}
}

func TestRealm_RandomHomePos(t *testing.T) {
	RegisterTestingT(t)
	r, _ := newMockRealm()

	baseMap := map[cb.Cube]struct{}{}
	cubeMap := map[cb.Cube]struct{}{}
	for {
		ok, x, y := r.randomBasePos()
		Ω(ok).Should(BeTrue())
		Ω(r.mapData.IsValidHomePosition(x, y)).Should(BeTrue())

		cube := cb.XYCube(x, y)
		baseMap[cube] = struct{}{}

		_, exist := cubeMap[cube]
		Ω(exist).Should(BeFalse())
		cubeMap[cb.XYCube(x, y)] = struct{}{}

		Ω(r.conflict.baseConflictCount[cube]).Should(Equal(1))
		neighbors := NeighborsInBaseConflictRange(x, y)
		for _, n := range neighbors {
			Ω(r.conflict.baseConflictCount[n] > 0).Should(BeTrue())

			cubeMap[n] = struct{}{}

			_, exist := baseMap[n]
			Ω(exist).Should(Equal(cube == n))
		}

		if len(baseMap) > 100 {
			break
		}
	}

	r.mapData.RangeBlock(u64.Min(r.GetRadius(), 5), blockdata.BlockRangeTypeCenterSpiral, func(blockX, blockY uint64) (toContinue bool) {
		r.mapData.RangeBlockHomeCubes(blockX, blockY, func(cube cb.Cube) (toContinue bool) {

			cc := r.baseConflictCount[cube]
			//Ω(cc > 0).Should(BeTrue())

			if _, exist := baseMap[cube]; exist {
				Ω(cc).Should(Equal(1))
			}

			baseCount := 0
			neighbors := NeighborsInBaseConflictRange(cube.XY())
			for _, n := range neighbors {
				if _, exist := baseMap[n]; exist {
					baseCount++
				}
			}

			Ω(baseCount).Should(Equal(cc))

			return true
		})
		return true
	})
}

func TestRealm_AddBase(t *testing.T) {
	RegisterTestingT(t)

	r, _ := newMockRealm()

	ok, x, y := r.randomBasePos()
	Ω(ok).Should(BeTrue())
	Ω(r.mapData.IsValidHomePosition(x, y)).Should(BeTrue())

	ctime := r.services.timeService.CurrentTime()

	// 设置英雄
	hero := entity.NewHero(1, "hero1", r.services.datas.HeroInitData(), ctime)
	mock.DefaultHero(hero)

	// 存在 baseRegion
	hero.SetBaseXY(1, 0, 0)
	processed, err := r.AddBase(ifacemock.HeroController.Id(), x, y, realmface.AddBaseHomeNewHero)
	Ω(processed).Should(BeTrue())
	Ω(err).Should(Equal(realmerr.ErrAddBaseAlreadyHasRealm))

	hero.SetBaseXY(0, 0, 0)

	// 城还活着
	hero.SetBaseLevel(1)
	hero.SetProsperity(1)

	ok, x, y = r.randomBasePos()
	Ω(ok).Should(BeTrue())
	Ω(r.mapData.IsValidHomePosition(x, y)).Should(BeTrue())

	processed, err = r.AddBase(ifacemock.HeroController.Id(), x, y, realmface.AddBaseHomeReborn)
	Ω(processed).Should(BeTrue())
	Ω(err).Should(Equal(realmerr.ErrAddBaseHomeAlive))

	hero.SetBaseLevel(0)
	hero.SetProsperity(0)

	// 第一次添加，成功
	processed, err = r.AddBase(ifacemock.HeroController.Id(), x, y, realmface.AddBaseHomeNewHero)
	Ω(processed).Should(BeTrue())
	Ω(err).Should(Succeed())

	Ω(hero.BaseRegion()).Should(Equal(r.id))
	Ω(hero.BaseX()).Should(Equal(x))
	Ω(hero.BaseY()).Should(Equal(y))
	Ω(hero.BaseLevel()).Should(Equal(uint64(1)))
	Ω(hero.Prosperity()).Should(Equal(hero.ProsperityCapcity()))

	base := r.getBase(hero.Id())
	Ω(base).ShouldNot(BeNil())
	Ω(base.Id()).Should(Equal(hero.Id()))
	Ω(base.Prosperity()).Should(Equal(hero.Prosperity()))
	Ω(base.ProsperityCapcity()).Should(Equal(hero.ProsperityCapcity()))

	ok, x, y = r.randomBasePos()
	Ω(ok).Should(BeTrue())
	Ω(r.mapData.IsValidHomePosition(x, y)).Should(BeTrue())

	// 第二次添加，失败
	processed, err = r.AddBase(ifacemock.HeroController.Id(), x, y, realmface.AddBaseHomeNewHero)
	Ω(processed).Should(BeTrue())
	Ω(err).Should(Equal(realmerr.ErrAddBaseAlreadyHasRealm))
}

func TestRealm_AddBaseReborn(t *testing.T) {
	RegisterTestingT(t)

	r, _ := newMockRealm()

	ok, x, y := r.randomBasePos()
	Ω(ok).Should(BeTrue())
	Ω(r.mapData.IsValidHomePosition(x, y)).Should(BeTrue())

	ctime := r.services.timeService.CurrentTime()

	// 设置英雄
	hero := entity.NewHero(1, "hero1", r.services.datas.HeroInitData(), ctime)
	mock.DefaultHero(hero)

	// 第一次添加，成功
	processed, err := r.AddBase(ifacemock.HeroController.Id(), x, y, realmface.AddBaseHomeReborn)
	Ω(processed).Should(BeTrue())
	Ω(err).Should(Succeed())

	Ω(hero.BaseRegion()).Should(Equal(r.id))
	Ω(hero.BaseX()).Should(Equal(x))
	Ω(hero.BaseY()).Should(Equal(y))
	Ω(hero.BaseLevel()).Should(Equal(uint64(1)))
	Ω(hero.Prosperity()).Should(Equal(hero.ProsperityCapcity() / 3))

	base := r.getBase(hero.Id())
	Ω(base).ShouldNot(BeNil())
	Ω(base.Id()).Should(Equal(hero.Id()))
	Ω(base.Prosperity()).Should(Equal(hero.ProsperityCapcity() / 3))
	Ω(base.ProsperityCapcity()).Should(Equal(hero.ProsperityCapcity()))
}

func TestRealm_AddBaseTransfer(t *testing.T) {
	RegisterTestingT(t)

	r, _ := newMockRealm()

	ok, x, y := r.randomBasePos()
	Ω(ok).Should(BeTrue())
	Ω(r.mapData.IsValidHomePosition(x, y)).Should(BeTrue())

	ctime := r.services.timeService.CurrentTime()

	// 设置英雄
	hero := entity.NewHero(1, "hero1", r.services.datas.HeroInitData(), ctime)
	mock.DefaultHero(hero)

	processed, err := r.AddBase(ifacemock.HeroController.Id(), x, y, realmface.AddBaseHomeTransfer)
	Ω(processed).Should(BeTrue())
	Ω(err).Should(Equal(realmerr.ErrAddBaseHomeNotAlive))

	hero.SetProsperity(1)
	hero.SetBaseLevel(2)

	processed, err = r.AddBase(ifacemock.HeroController.Id(), x, y, realmface.AddBaseHomeTransfer)
	Ω(processed).Should(BeTrue())
	Ω(err).Should(Succeed())

	Ω(hero.BaseRegion()).Should(Equal(r.id))
	Ω(hero.BaseX()).Should(Equal(x))
	Ω(hero.BaseY()).Should(Equal(y))
	Ω(hero.BaseLevel()).Should(Equal(uint64(2)))
	Ω(hero.Prosperity()).Should(Equal(uint64(1)))

	base := r.getBase(hero.Id())
	Ω(base).ShouldNot(BeNil())
	Ω(base.Id()).Should(Equal(hero.Id()))
	Ω(base.Prosperity()).Should(Equal(uint64(1)))
	Ω(base.ProsperityCapcity()).Should(Equal(hero.ProsperityCapcity()))
}

func TestRealm_RemoveBase(t *testing.T) {
	RegisterTestingT(t)

	r, _ := newMockRealm()

	ctime := r.services.timeService.CurrentTime()

	// 设置英雄
	hero := entity.NewHero(1, "hero1", r.services.datas.HeroInitData(), ctime)
	mock.DefaultHero(hero)

	// 移除不存在的主城
	processed, err, originX, originY := r.RemoveBase(hero.Id(), true,
		r.services.datas.TextHelp().MRDRMoveBase4a.Text, r.services.datas.TextHelp().MRDRMoveBase4d.Text)
	Ω(processed).Should(BeTrue())
	Ω(err).Should(Equal(realmerr.ErrRemoveBaseSelfNoBase))
	Ω(cb.XYCube(originX, originY)).Should(Equal(cb.XYCube(0, 0)))

	// 新建主城
	ok, x, y := r.randomBasePos()
	Ω(ok).Should(BeTrue())
	Ω(r.mapData.IsValidHomePosition(x, y)).Should(BeTrue())

	processed, err = r.AddBase(ifacemock.HeroController.Id(), x, y, realmface.AddBaseHomeNewHero)
	Ω(processed).Should(BeTrue())
	Ω(err).Should(Succeed())

	Ω(hero.BaseRegion()).Should(Equal(r.id))

	// 存在部队的测试，在部队测试中补充

	// 移除主城
	processed, err, originX, originY = r.RemoveBase(hero.Id(), true,
		r.services.datas.TextHelp().MRDRMoveBase4a.Text, r.services.datas.TextHelp().MRDRMoveBase4d.Text)
	Ω(processed).Should(BeTrue())
	Ω(err).Should(Succeed())
	Ω(cb.XYCube(originX, originY)).Should(Equal(cb.XYCube(x, y)))

	Ω(hero.BaseRegion()).Should(Equal(int64(0)))
}

func TestRealm_InvasionSameRealm(t *testing.T) {
	RegisterTestingT(t)

	// 测试过程，创建2个英雄，在r1中新建主城，然后，A对B进行出征

	// 验证马的状态是否正确

	// 将时间设置到到达目的地，马的状态将转为持续掠夺

	// 将时间设置到扣繁荣度tick，验证马抢到的资源，B扣的繁荣度

	// 将时间设置到持续掠夺结束时间，验证马到时间回家

	// 将时间设置到马已经回到家，验证马到家改状态

	r, _ := newMockRealm()

	ifacemock.FightXService.Mock(ifacemock.FightXService.SendFightRequest, func(ctx *entity.TlogFightContext, combatScene *scene.CombatScene, attackerId, defenserId int64, attacker, defenser *shared_proto.CombatPlayerProto) *server_proto.CombatXResponseServerProto {
		return &server_proto.CombatXResponseServerProto{AttackerWin: true}
	})

	ifacemock.SeasonService.Mock(ifacemock.SeasonService.Season, func() *season.SeasonData {
		return r.services.datas.SeasonData().MinKeyData
	})

	ctime := r.services.timeService.CurrentTime()

	// 第一个英雄
	defenserHero := entity.NewHero(1, "hero1", r.services.datas.HeroInitData(), ctime)
	defenserHero.SetGuild(1)
	defenserHero.GetUnsafeResource().AddResource(
		u64.Sub(defenserHero.ProtectedCapcity()+10, defenserHero.GetUnsafeResource().Gold()),
		u64.Sub(defenserHero.ProtectedCapcity()+20, defenserHero.GetUnsafeResource().Food()),
		u64.Sub(defenserHero.ProtectedCapcity()+30, defenserHero.GetUnsafeResource().Wood()),
		u64.Sub(defenserHero.ProtectedCapcity()+40, defenserHero.GetUnsafeResource().Stone()),
	)
	mock.DefaultHero(defenserHero)

	ok, x, y := r.randomBasePos()
	Ω(ok).Should(BeTrue())
	Ω(r.mapData.IsValidHomePosition(x, y)).Should(BeTrue())

	processed, err := r.AddBase(ifacemock.HeroController.Id(), x, y, realmface.AddBaseHomeNewHero)
	Ω(processed).Should(BeTrue())
	Ω(err).Should(Succeed())

	defenserBase := r.getBase(defenserHero.Id())
	Ω(defenserBase).ShouldNot(BeNil())
	Ω(defenserBase.RegionID()).Should(Equal(r.id))
	Ω(defenserBase.MianDisappearTime()).ShouldNot(Equal(int32(0)))

	// 第二个英雄
	attackerHero := entity.NewHero(2, "hero2", r.services.datas.HeroInitData(), ctime)
	attackerHero.SetGuild(2)
	attackerHero.SetWhiteFlag(1, 1, ctime.Add(24*time.Hour))
	mock.DefaultHero(attackerHero)

	ifacemock.GuildService.Mock(ifacemock.GuildService.GetSnapshot, func(id int64) *guildsnapshotdata.GuildSnapshot {
		if id == 2 {
			return &guildsnapshotdata.GuildSnapshot{
				Id:       2,
				Name:     "guild2",
				FlagName: "flag2",
			}
		}
		return nil
	})

	ok, x, y = r.randomBasePos()
	Ω(ok).Should(BeTrue())
	Ω(r.mapData.IsValidHomePosition(x, y)).Should(BeTrue())

	processed, err = r.AddBase(ifacemock.HeroController.Id(), x, y, realmface.AddBaseHomeNewHero)
	Ω(processed).Should(BeTrue())
	Ω(err).Should(Succeed())

	attackerBase := r.getBase(attackerHero.Id())
	Ω(attackerBase).ShouldNot(BeNil())
	Ω(attackerBase.RegionID()).Should(Equal(r.id))
	Ω(attackerBase.MianDisappearTime()).ShouldNot(Equal(int32(0)))

	// 添加部队
	mock.SetHeroTroopCaptain(attackerHero, r.services.datas)
	heroTroop := attackerHero.GetTroopByIndex(0)
	for _, pos := range heroTroop.Pos() {
		captain := pos.Captain()
		Ω(captain).ShouldNot(BeNil())
		Ω(captain.IsOutSide()).Should(BeFalse())
	}
	Ω(heroTroop.IsOutside()).Should(BeFalse())

	// 出征
	processed, err = r.Invasion(ifacemock.HeroController, true, 1, 1, 0)
	Ω(processed).Should(BeTrue())
	Ω(err).Should(Equal(realmerr.ErrInvasionMian))

	// 移除免战，再打
	defenserBase.SetMianDisappearTime(0)

	processed, err = r.Invasion(ifacemock.HeroController, true, 1, 1, 0)
	Ω(processed).Should(BeTrue())
	Ω(err).Should(Succeed())

	Ω(attackerBase.MianDisappearTime()).Should(Equal(int32(0)))

	troop := attackerBase.selfTroops[heroTroop.Id()]
	Ω(troop).ShouldNot(BeNil())

	//speed := r.CalcMoveSpeed(defenserBase.Id())
	Ω(troop.State()).Should(Equal(realmface.MovingToInvade)) // 马出发中
	Ω(troop.MoveStartTime()).Should(Equal(ctime))
	Ω(troop.MoveArriveTime()).Should(Equal(ctime.Add(
		troop.getMoveDuration(r),
	)))

	baseEquips(attackerHero, attackerBase)
	baseEquips(defenserHero, defenserBase)
	troopEquips(heroTroop, troop, attackerBase, defenserBase)

	// 到达时间

	// 出征到达
	event := r.eventQueue.Peek()
	Ω(event).ShouldNot(BeNil())
	Ω(r.queueFunc(true, func() {
		mock.SetTime(troop.moveArriveTime.Add(1))
		event.Data().(func(*Realm))(r) // 更新一下
	})).Should(BeTrue())

	troop = attackerBase.selfTroops[heroTroop.Id()] // 马还在
	Ω(troop).ShouldNot(BeNil())

	Ω(troop.State()).Should(Equal(realmface.Robbing)) // 持续掠夺
	//Ω(troop.MoveStartTime()).Should(Equal(time.Time{}))
	//Ω(troop.MoveArriveTime()).Should(Equal(time.Time{}))

	// 第一次，抢的东西不加到累积奖励中
	gold, food, wood, stone := troop.Carrying()
	Ω(gold == 0).Should(BeTrue())
	Ω(food == 0).Should(BeTrue())
	Ω(wood == 0).Should(BeTrue())
	Ω(stone == 0).Should(BeTrue())

	originGold, originFood, originWood, originStone := attackerHero.GetUnsafeResource().Gold(), attackerHero.GetUnsafeResource().Food(), attackerHero.GetUnsafeResource().Wood(), attackerHero.GetUnsafeResource().Stone()

	baseEquips(attackerHero, attackerBase)
	baseEquips(defenserHero, defenserBase)
	troopEquips(heroTroop, troop, attackerBase, defenserBase)

	// 拔掉白旗
	Ω(attackerHero.GetWhiteFlagHeroId()).Should(Equal(int64(0)))
	Ω(attackerHero.GetWhiteFlagGuildId()).Should(Equal(int64(0)))
	Ω(attackerHero.GetWhiteFlagDisappearTime()).Should(Equal(time.Time{}))

	event = r.eventQueue.Peek()
	Ω(event).ShouldNot(BeNil())
	Ω(r.queueFunc(true, func() {
		// 到达结束时间
		if troop.robbingEndTime.Before(troop.nextReduceProsperityTime) {
			mock.SetTime(troop.nextReduceProsperityTime.Add(1))
		} else {
			mock.SetTime(troop.robbingEndTime.Add(1))
		}

		event.Data().(func(*Realm))(r) // 更新一下
	})).Should(BeTrue())

	ctime = r.services.timeService.CurrentTime()

	Ω(troop.State()).Should(Equal(realmface.InvadeMovingBack)) // 抢完回家
	Ω(troop.MoveStartTime()).Should(Equal(ctime))

	//speed = r.CalcMoveSpeed(defenserBase.Id())

	Ω(troop.MoveArriveTime()).Should(Equal(ctime.Add(troop.getMoveDuration(r))))

	baseEquips(attackerHero, attackerBase)
	baseEquips(defenserHero, defenserBase)
	troopEquips(heroTroop, troop, attackerBase, nil)

	Ω(defenserBase.targetingTroops[troop.Id()]).Should(BeNil())

	// 插白旗
	Ω(defenserHero.GetWhiteFlagHeroId()).Should(Equal(int64(2)))
	Ω(defenserHero.GetWhiteFlagGuildId()).Should(Equal(int64(2)))
	Ω(defenserHero.GetWhiteFlagDisappearTime()).Should(Equal(ctime.Add(r.services.datas.RegionConfig().WhiteFlagDuration)))

	// 到家
	gold, food, wood, stone = troop.Carrying()
	//Ω(gold > 0).Should(BeTrue())
	//Ω(food > 0).Should(BeTrue())
	//Ω(wood > 0).Should(BeTrue())
	//Ω(stone > 0).Should(BeTrue())

	event = r.eventQueue.Peek()
	Ω(event).ShouldNot(BeNil())
	Ω(r.queueFunc(true, func() {
		mock.SetTime(troop.moveArriveTime)
		event.Data().(func(*Realm))(r) // 更新一下
	})).Should(BeTrue())

	troop = attackerBase.selfTroops[heroTroop.Id()] // 马没了
	Ω(troop).Should(BeNil())

	for _, pos := range heroTroop.Pos() {
		captain := pos.Captain()
		Ω(captain).ShouldNot(BeNil())
		Ω(captain.IsOutSide()).Should(BeFalse())
	}
	Ω(heroTroop.IsOutside()).Should(BeFalse())
	Ω(heroTroop.GetInvateInfo()).Should(BeNil())

	// 加资源
	Ω(attackerHero.GetUnsafeResource().Gold()).Should(Equal(originGold + gold))
	Ω(attackerHero.GetUnsafeResource().Food()).Should(Equal(originFood + food))
	Ω(attackerHero.GetUnsafeResource().Wood()).Should(Equal(originWood + wood))
	Ω(attackerHero.GetUnsafeResource().Stone()).Should(Equal(originStone + stone))

	baseEquips(attackerHero, attackerBase)
	baseEquips(defenserHero, defenserBase)
}

func TestRealm_InvasionAssembly(t *testing.T) {
	RegisterTestingT(t)

	// 测试过程，创建3个英雄，在r1中新建主城，A和C属于同一个联盟，然后，A对B进行集结出征

	// 验证马的状态是否正确

	// 将时间设置到到达目的地，马的状态将转为持续掠夺

	// 将时间设置到扣繁荣度tick，验证马抢到的资源，B扣的繁荣度

	// 将时间设置到持续掠夺结束时间，验证马到时间回家

	// 将时间设置到马已经回到家，验证马到家改状态

	r, _ := newMockRealm()

	ifacemock.FightXService.Mock(ifacemock.FightXService.SendFightRequest, func(ctx *entity.TlogFightContext, combatScene *scene.CombatScene, attackerId, defenserId int64, attacker, defenser *shared_proto.CombatPlayerProto) *server_proto.CombatXResponseServerProto {
		return &server_proto.CombatXResponseServerProto{AttackerWin: true}
	})

	ifacemock.SeasonService.Mock(ifacemock.SeasonService.Season, func() *season.SeasonData {
		return r.services.datas.SeasonData().MinKeyData
	})

	ctime := r.services.timeService.CurrentTime()

	// 第一个英雄
	defenserHero := entity.NewHero(1, "hero1", r.services.datas.HeroInitData(), ctime)
	defenserHero.SetGuild(1)
	defenserHero.GetUnsafeResource().AddResource(
		u64.Sub(defenserHero.ProtectedCapcity()+10, defenserHero.GetUnsafeResource().Gold()),
		u64.Sub(defenserHero.ProtectedCapcity()+20, defenserHero.GetUnsafeResource().Food()),
		u64.Sub(defenserHero.ProtectedCapcity()+30, defenserHero.GetUnsafeResource().Wood()),
		u64.Sub(defenserHero.ProtectedCapcity()+40, defenserHero.GetUnsafeResource().Stone()),
	)
	mock.DefaultHero(defenserHero)

	ok, x, y := r.randomBasePos()
	Ω(ok).Should(BeTrue())
	Ω(r.mapData.IsValidHomePosition(x, y)).Should(BeTrue())

	processed, err := r.AddBase(ifacemock.HeroController.Id(), x, y, realmface.AddBaseHomeNewHero)
	Ω(processed).Should(BeTrue())
	Ω(err).Should(Succeed())

	defenserBase := r.getBase(defenserHero.Id())
	Ω(defenserBase).ShouldNot(BeNil())
	Ω(defenserBase.RegionID()).Should(Equal(r.id))
	Ω(defenserBase.MianDisappearTime()).ShouldNot(Equal(int32(0)))

	ifacemock.GuildService.Mock(ifacemock.GuildService.GetSnapshot, func(id int64) *guildsnapshotdata.GuildSnapshot {
		if id == 2 {
			return &guildsnapshotdata.GuildSnapshot{
				Id:       2,
				Name:     "guild2",
				FlagName: "flag2",
			}
		}
		return nil
	})

	// 第三个英雄，集结出征
	attackerHero := entity.NewHero(2, "hero2", r.services.datas.HeroInitData(), ctime)
	attackerHero.SetGuild(2)
	attackerHero.SetWhiteFlag(1, 1, ctime.Add(24*time.Hour))
	mock.DefaultHero(attackerHero)

	ok, x, y = r.randomBasePos()
	Ω(ok).Should(BeTrue())
	Ω(r.mapData.IsValidHomePosition(x, y)).Should(BeTrue())

	processed, err = r.AddBase(ifacemock.HeroController.Id(), x, y, realmface.AddBaseHomeNewHero)
	Ω(processed).Should(BeTrue())
	Ω(err).Should(Succeed())

	attackerBase := r.getBase(attackerHero.Id())
	Ω(attackerBase).ShouldNot(BeNil())
	Ω(attackerBase.RegionID()).Should(Equal(r.id))
	Ω(attackerBase.MianDisappearTime()).ShouldNot(Equal(int32(0)))

	// 添加部队
	mock.SetHeroTroopCaptain(attackerHero, r.services.datas)
	attackerHeroTroop := attackerHero.GetTroopByIndex(0)
	for _, pos := range attackerHeroTroop.Pos() {
		captain := pos.Captain()
		Ω(captain).ShouldNot(BeNil())
		Ω(captain.IsOutSide()).Should(BeFalse())
	}
	Ω(attackerHeroTroop.IsOutside()).Should(BeFalse())

	// 出征
	processed, err = r.CreateAssembly(ifacemock.HeroController, 1, 1, 0, 5, 1*time.Second)
	Ω(processed).Should(BeTrue())
	Ω(err).Should(Equal(realmerr.ErrCreateAssemblyMian))

	// 移除免战，再打
	defenserBase.SetMianDisappearTime(0)

	processed, err = r.CreateAssembly(ifacemock.HeroController, 1, 1, 0, 5, 1*time.Second)
	Ω(processed).Should(BeTrue())
	Ω(err).Should(Succeed())

	Ω(attackerBase.MianDisappearTime()).Should(Equal(int32(0)))

	troop := attackerBase.selfTroops[attackerHeroTroop.Id()]
	Ω(troop).ShouldNot(BeNil())

	// 集结过期
	event := r.eventQueue.Peek()
	Ω(event).ShouldNot(BeNil())
	Ω(r.queueFunc(true, func() {
		mock.SetTime(troop.moveArriveTime.Add(1))
		event.Data().(func(*Realm))(r) // 更新一下
	})).Should(BeTrue())

	Ω(attackerHeroTroop.GetInvateInfo()).Should(BeNil())
	troop = attackerBase.selfTroops[attackerHeroTroop.Id()]
	Ω(troop).Should(BeNil())

	// 创建长时间的集结
	processed, err = r.CreateAssembly(ifacemock.HeroController, 1, 1, 0, 5, 48*time.Hour)
	Ω(processed).Should(BeTrue())
	Ω(err).Should(Succeed())

	troop = attackerBase.selfTroops[attackerHeroTroop.Id()]
	Ω(troop).ShouldNot(BeNil())
	Ω(troop.GetAssembly()).ShouldNot(BeNil())

	// 加入集结
	// 第3个英雄，参加集结
	joinHero := entity.NewHero(3, "hero3", r.services.datas.HeroInitData(), ctime)
	joinHero.SetGuild(2)
	joinHero.SetWhiteFlag(1, 1, ctime.Add(24*time.Hour))
	mock.DefaultHero(joinHero)

	ok, x, y = r.randomBasePos()
	Ω(ok).Should(BeTrue())
	Ω(r.mapData.IsValidHomePosition(x, y)).Should(BeTrue())

	processed, err = r.AddBase(ifacemock.HeroController.Id(), x, y, realmface.AddBaseHomeNewHero)
	Ω(processed).Should(BeTrue())
	Ω(err).Should(Succeed())

	joinBase := r.getBase(joinHero.Id())
	Ω(joinBase).ShouldNot(BeNil())
	Ω(joinBase.RegionID()).Should(Equal(r.id))
	Ω(joinBase.MianDisappearTime()).ShouldNot(Equal(int32(0)))

	// 添加部队
	mock.SetHeroTroopCaptain(joinHero, r.services.datas)
	joinHeroTroop := joinHero.GetTroopByIndex(0)
	for _, pos := range joinHeroTroop.Pos() {
		captain := pos.Captain()
		Ω(captain).ShouldNot(BeNil())
		Ω(captain.IsOutSide()).Should(BeFalse())
	}
	Ω(joinHeroTroop.IsOutside()).Should(BeFalse())

	// 加入集结
	processed, err = r.JoinAssembly(ifacemock.HeroController, attackerHero.Id(), attackerHeroTroop.Id(), 0)
	Ω(processed).Should(BeTrue())
	Ω(err).Should(Succeed())

	Ω(joinBase.MianDisappearTime()).Should(Equal(int32(0)))

	joinTroop := joinBase.selfTroops[joinHeroTroop.Id()]
	Ω(joinTroop).ShouldNot(BeNil())
	Ω(joinTroop).Should(Equal(attackerBase.targetingTroops[joinHeroTroop.Id()]))
	Ω(joinTroop.GetAssembly()).Should(Equal(troop.GetAssembly()))

	if joinTroop.moveArriveTime.After(troop.moveArriveTime) {
		speedUpRate := i64.DivideTimes(joinTroop.moveArriveTime.Unix(), troop.moveArriveTime.Unix())

		processed, err = r.SpeedUp(ifacemock.HeroController, joinHeroTroop.Id(), 0, float64(speedUpRate), 0)
		Ω(processed).Should(BeTrue())
		Ω(err).Should(Succeed())
	}

	Ω(joinTroop.moveArriveTime.Before(troop.moveArriveTime)).Should(BeTrue())

	// 到达集结地
	event = r.eventQueue.Peek()
	Ω(event).ShouldNot(BeNil())
	Ω(r.queueFunc(true, func() {
		mock.SetTime(joinTroop.moveArriveTime.Add(1))
		event.Data().(func(*Realm))(r) // 更新一下
	})).Should(BeTrue())

	Ω(joinTroop.State()).Should(Equal(realmface.AssemblyArrived))

	// 到达集结出发时间
	event = r.eventQueue.Peek()
	Ω(event).ShouldNot(BeNil())
	Ω(r.queueFunc(true, func() {
		mock.SetTime(troop.moveArriveTime.Add(1))
		event.Data().(func(*Realm))(r) // 更新一下
	})).Should(BeTrue())

	// 加入集结的人，几乎没有变化
	Ω(joinTroop.State()).Should(Equal(realmface.AssemblyArrived))

	// 创建集结的人，状态变成出征
	Ω(troop.State()).Should(Equal(realmface.MovingToInvade))

	// 到达集结地，开始战斗
	event = r.eventQueue.Peek()
	Ω(event).ShouldNot(BeNil())
	Ω(r.queueFunc(true, func() {
		mock.SetTime(troop.moveArriveTime.Add(1))
		event.Data().(func(*Realm))(r) // 更新一下
	})).Should(BeTrue())

	Ω(troop.State()).Should(Equal(realmface.Robbing))
	Ω(troop.assembly).ShouldNot(BeNil())
	Ω(attackerHeroTroop.GetInvateInfo()).ShouldNot(BeNil())

	// 持续掠夺结束
	event = r.eventQueue.Peek()
	Ω(event).ShouldNot(BeNil())
	Ω(r.queueFunc(true, func() {
		mock.SetTime(troop.robbingEndTime.Add(1))
		event.Data().(func(*Realm))(r) // 更新一下
	})).Should(BeTrue())

	Ω(troop.State()).Should(Equal(realmface.InvadeMovingBack))
	Ω(troop.GetAssembly()).Should(BeNil())

	Ω(joinTroop.State()).Should(Equal(realmface.InvadeMovingBack))
	Ω(joinTroop.GetAssembly()).Should(BeNil())

	// 到家
	event = r.eventQueue.Peek()
	Ω(event).ShouldNot(BeNil())
	Ω(r.queueFunc(true, func() {
		mock.SetTime(troop.moveArriveTime)
		event.Data().(func(*Realm))(r) // 更新一下
	})).Should(BeTrue())

	event = r.eventQueue.Peek()
	Ω(event).ShouldNot(BeNil())
	Ω(r.queueFunc(true, func() {
		mock.SetTime(troop.moveArriveTime)
		event.Data().(func(*Realm))(r) // 更新一下
	})).Should(BeTrue())

	troop = attackerBase.selfTroops[attackerHeroTroop.Id()] // 马没了
	Ω(troop).Should(BeNil())

	for _, pos := range attackerHeroTroop.Pos() {
		captain := pos.Captain()
		Ω(captain).ShouldNot(BeNil())
		Ω(captain.IsOutSide()).Should(BeFalse())
	}
	Ω(attackerHeroTroop.IsOutside()).Should(BeFalse())
	Ω(attackerHeroTroop.GetInvateInfo()).Should(BeNil())

	joinTroop = joinBase.selfTroops[joinHeroTroop.Id()] // 马没了
	Ω(joinTroop).Should(BeNil())

	for _, pos := range joinHeroTroop.Pos() {
		captain := pos.Captain()
		Ω(captain).ShouldNot(BeNil())
		Ω(captain.IsOutSide()).Should(BeFalse())
	}
	Ω(joinHeroTroop.IsOutside()).Should(BeFalse())
	Ω(joinHeroTroop.GetInvateInfo()).Should(BeNil())

	baseEquips(attackerHero, attackerBase)
	baseEquips(defenserHero, defenserBase)
	baseEquips(joinHero, joinBase)
}

func TestMian(t *testing.T) {

	RegisterTestingT(t)

	r, _ := newMockRealm()
	ctime := r.services.timeService.CurrentTime()

	// 第一个英雄
	hero := entity.NewHero(1, "hero1", r.services.datas.HeroInitData(), ctime)
	mock.DefaultHero(hero)

	processed, err := r.Mian(hero.Id(), ctime, false)
	Ω(processed).Should(BeTrue())
	Ω(err).Should(Succeed())

	Ω(hero.GetMianDisappearTime()).Should(Equal(time.Time{}))

	disappearTime := ctime.Add(time.Second)
	processed, err = r.Mian(hero.Id(), disappearTime, false)
	Ω(processed).Should(BeTrue())
	Ω(err).Should(Equal(realmerr.ErrMianSelfNoBase))

	// 添加主城
	ok, x, y := r.randomBasePos()
	Ω(ok).Should(BeTrue())
	Ω(r.mapData.IsValidHomePosition(x, y)).Should(BeTrue())

	processed, err = r.AddBase(ifacemock.HeroController.Id(), x, y, realmface.AddBaseHomeNewHero)
	Ω(processed).Should(BeTrue())
	Ω(err).Should(Succeed())

	base := r.getBase(hero.Id())
	Ω(base.MianDisappearTime()).ShouldNot(Equal(int32(0)))

	// 移除新手免战
	base.SetMianDisappearTime(0)

	// 设置成功
	processed, err = r.Mian(hero.Id(), disappearTime, false)
	Ω(processed).Should(BeTrue())
	Ω(err).Should(Succeed())

	Ω(base.MianDisappearTime()).Should(Equal(timeutil.Marshal32(disappearTime)))
	Ω(hero.GetMianDisappearTime()).Should(Equal(disappearTime))

	// 再次设置
	processed, err = r.Mian(hero.Id(), disappearTime, false)
	Ω(processed).Should(BeTrue())
	Ω(err).Should(Equal(realmerr.ErrMianExist))

	processed, err = r.Mian(hero.Id(), disappearTime, true)
	Ω(processed).Should(BeTrue())
	Ω(err).Should(Equal(realmerr.ErrMianCantOverwrite))

	// 加大时间，再次设置
	disappearTime = disappearTime.Add(time.Second)

	processed, err = r.Mian(hero.Id(), disappearTime, false)
	Ω(processed).Should(BeTrue())
	Ω(err).Should(Equal(realmerr.ErrMianExist))

	processed, err = r.Mian(hero.Id(), disappearTime, true)
	Ω(processed).Should(BeTrue())
	Ω(err).Should(Succeed())

	Ω(base.MianDisappearTime()).Should(Equal(timeutil.Marshal32(disappearTime)))
	Ω(hero.GetMianDisappearTime()).Should(Equal(disappearTime))

	// 移除免战
	Ω(r.queueFunc(true, func() {
		r.doRemoveBaseMian(base, true)
	})).Should(BeTrue())

	Ω(base.MianDisappearTime()).Should(Equal(int32(0)))
	Ω(hero.GetMianDisappearTime()).Should(Equal(time.Time{}))
}

func TestUpgradeBase(t *testing.T) {
	RegisterTestingT(t)

	r, _ := newMockRealm()
	ctime := r.services.timeService.CurrentTime()

	// 第一个英雄
	hero := entity.NewHero(1, "hero1", r.services.datas.HeroInitData(), ctime)
	mock.DefaultHero(hero)

	// 添加主城
	ok, x, y := r.randomBasePos()
	Ω(ok).Should(BeTrue())
	Ω(r.mapData.IsValidHomePosition(x, y)).Should(BeTrue())

	processed, err := r.AddBase(ifacemock.HeroController.Id(), x, y, realmface.AddBaseHomeNewHero)
	Ω(processed).Should(BeTrue())
	Ω(err).Should(Succeed())

	base := r.getBase(hero.Id())

	// 加繁荣度
	prosperity := hero.Prosperity()
	prosperityCapcity := hero.ProsperityCapcity()
	for _, v := range r.services.datas.GetBaseLevelDataArray() {
		nextLevel := r.services.datas.GetBaseLevelData(v.Level + 1)
		if nextLevel == nil {
			continue
		}

		hero.AddProsperityCapcity(v.UpgradeProsperity)

		// 繁荣度上限已经加了
		prosperityCapcity += v.UpgradeProsperity
		Ω(hero.ProsperityCapcity()).Should(Equal(prosperityCapcity))

		// 还没加
		Ω(hero.Prosperity()).Should(Equal(prosperity))
		addProsperity := v.UpgradeProsperity - prosperity
		prosperity = v.UpgradeProsperity

		processed, err = r.AddProsperity(hero.Id(), addProsperity)
		Ω(err).Should(Succeed())
		Ω(processed).Should(BeTrue())

		Ω(r.queueFunc(true, func() {
			Ω(base.Prosperity()).Should(Equal(prosperity))
			Ω(base.ProsperityCapcity()).Should(Equal(prosperityCapcity))
		})).Should(BeTrue())

		// 升级
		processed, err = r.UpgradeBase(ifacemock.HeroController)
		Ω(processed).Should(BeTrue())
		Ω(err).Should(Succeed())

		Ω(base.BaseLevel()).Should(Equal(nextLevel.Level))

		// Npc资源点占用
		checkNpcResourceConflict(r, base, hero)
	}

	// 原来的地图迁移
	ok, newX, newY := r.randomBasePos()
	Ω(ok).Should(BeTrue())
	Ω(r.mapData.IsValidHomePosition(newX, newY)).Should(BeTrue())
	Ω(cb.XYCube(newX, newY)).ShouldNot(Equal(cb.XYCube(x, y)))

	processed, err = r.MoveBase(ifacemock.HeroController, x, y, newX, newY, false)
	Ω(processed).Should(BeTrue())
	Ω(err).Should(Succeed())

	checkNpcResourceConflict(r, base, hero)

	// 移除基地
	processed, err, originX, originY := r.RemoveBase(hero.Id(), true, r.getTextHelp().MRDRMoveBase4a.Text, r.getTextHelp().MRDRMoveBase4d.Text)
	Ω(processed).Should(BeTrue())
	Ω(err).Should(Succeed())
	Ω(cb.XYCube(newX, newY)).Should(Equal(cb.XYCube(originX, originY)))

	base = r.getBase(hero.Id())
	Ω(base).Should(BeNil())

	Ω(baseLen(r)).Should(BeZero())
	Ω(r.resourceConflictHeroMap).Should(BeEmpty())
	Ω(r.conflict.baseConflictCount).Should(BeEmpty())
}

func baseEquips(hero *entity.Hero, base *baseWithData) {

	Ω(base.Id()).Should(Equal(hero.Id()))
	Ω(base.IdBytes()).Should(Equal(hero.IdBytes()))
	Ω(base.internalBase.HeroName()).Should(Equal(hero.Name()))
	//Ω(base.level()).Should(Equal(hero.Level()))

	Ω(base.GuildId()).Should(Equal(hero.GuildId()))

	if base.isHeroHomeBase() {
		Ω(base.RegionID()).Should(Equal(hero.BaseRegion()))
		Ω(base.BaseX()).Should(Equal(hero.BaseX()))
		Ω(base.BaseY()).Should(Equal(hero.BaseY()))
		Ω(base.Prosperity()).Should(Equal(hero.Prosperity()))
		Ω(base.ProsperityCapcity()).Should(Equal(hero.ProsperityCapcity()))
		Ω(base.BaseLevel()).Should(Equal(hero.BaseLevel()))

		if b := GetHomeBase(base); b != nil {
			Ω(b.WhiteFlagGuildId()).Should(Equal(hero.GetWhiteFlagGuildId()))
			Ω(b.WhiteFlagDisappearTime()).Should(Equal(timeutil.Marshal32(hero.GetWhiteFlagDisappearTime())))
		}
	}

	//if base.isHeroTentBase() {
	//	Ω(base.RegionID()).Should(Equal(hero.TentBaseRegion()))
	//	Ω(base.BaseX()).Should(Equal(hero.TentBaseX()))
	//	Ω(base.BaseY()).Should(Equal(hero.TentBaseY()))
	//	Ω(base.Prosperity()).Should(Equal(hero.TentProsperity()))
	//	Ω(base.ProsperityCapcity()).Should(Equal(hero.TentOutsideProsperityCapcity()))
	//	Ω(base.BaseLevel()).Should(Equal(hero.TentBaseLevel()))
	//
	//	Ω(base.WhiteFlagGuildId()).Should(Equal(int64(0)))
	//	Ω(base.WhiteFlagDisappearTime()).Should(Equal(int32(0)))
	//}

}

func troopEquips(heroTroop *entity.Troop, troop *troop, selfBase, targetBase *baseWithData) {

	Ω(troop.startingBase).Should(Equal(selfBase))
	Ω(troop.startingBase.selfTroops[troop.Id()]).Should(Equal(troop))

	if targetBase != nil {
		Ω(troop.targetBase == targetBase).Should(BeTrue())
		Ω(troop.targetBase.targetingTroops[troop.Id()] == troop).Should(BeTrue())

		Ω(targetBase.Id()).Should(Equal(heroTroop.GetInvateInfo().TargetBaseID()))
	} else {
		Ω(troop.targetBase).Should(BeNil())
	}

	Ω(heroTroop.IsOutside()).Should(BeTrue())
	Ω(heroTroop.GetInvateInfo()).ShouldNot(BeNil())
	Ω(selfBase.RegionID()).Should(Equal(heroTroop.GetInvateInfo().RegionID()))

	Ω(troop.Id()).Should(Equal(heroTroop.Id()))
	Ω(troop.state).Should(Equal(heroTroop.GetInvateInfo().State()))

	Ω(troop.MoveStartTime()).Should(Equal(heroTroop.GetInvateInfo().MoveStartTime()))
	Ω(troop.moveArriveTime).Should(Equal(heroTroop.GetInvateInfo().MoveArriveTime()))
	Ω(troop.NextReduceProsperityTime()).Should(Equal(heroTroop.GetInvateInfo().NextReduceProsperityTime()))
	Ω(troop.RobbingEndTime()).Should(Equal(heroTroop.GetInvateInfo().RobbingEndTime()))
	Ω(troop.BackHomeTargetX()).Should(Equal(heroTroop.GetInvateInfo().BackHomeTargetX()))
	Ω(troop.BackHomeTargetY()).Should(Equal(heroTroop.GetInvateInfo().BackHomeTargetY()))

	Ω(troop.AccumRobPrize()).Should(Equal(heroTroop.GetInvateInfo().AccumRobPrize()))

	captainMap := make(map[uint64]*entity.Captain)
	for _, pos := range heroTroop.Pos() {
		c := pos.Captain()
		if c != nil {
			captainMap[c.Id()] = c
		}
	}

	Ω(len(captainMap)).Should(Equal(len(troop.captains)))

	for _, c := range troop.captains {
		Ω(c).ShouldNot(BeNil())
		Ω(c.Index() <= len(heroTroop.Pos())).Should(BeTrue())

		heroCaptain := captainMap[c.Id()]
		Ω(heroCaptain).ShouldNot(BeNil())

		Ω(c.Proto().Soldier).Should(Equal(u64.Int32(heroCaptain.Soldier())))
	}
}

func newMockRealm() (home, tent *Realm) {
	datas, err := config.LoadConfigDatas(confpath.GetConfigPath())
	Ω(err).Should(Succeed())

	realmArray := make([]*Realm, 2)

	otherRealmEventQueue := event.NewFuncQueue(1024, "RealmService.OtherRealmEventQueue")
	service := &services{
		dep:                  mock.MockDep(),
		datas:                datas,
		dbService:            ifacemock.DbService,
		timeService:          ifacemock.TimeService,
		world:                ifacemock.WorldService,
		heroDataService:      ifacemock.HeroDataService,
		fightModule:          ifacemock.FightXService,
		mail:                 ifacemock.MailModule,
		heroSnapshotService:  ifacemock.HeroSnapshotService,
		guildService:         ifacemock.GuildService,
		reminderService:      ifacemock.ReminderService,
		otherRealmEventQueue: otherRealmEventQueue,
		farmService:          ifacemock.FarmService,
		seasonService:        ifacemock.SeasonService,
		extraTimesService:    ifacemock.ExtraTimesService,
		pushService:          ifacemock.PushService,
		baiZhanService:       ifacemock.BaiZhanService,
		//getRealm: func(id int64) *Realm {
		//	for _, v := range realmArray {
		//		if v != nil && v.id == id {
		//			return v
		//		}
		//	}
		//	return nil
		//},
		homeNpcBaseIdGen: npcid2.NewNpcIdGen(0, npcid2.NpcType_HomeNpc),
	}

	//datas.RegionData().MinKeyData.InitRadius = 2

	home = newRealm(realmface.GetRealmId(1, 0, datas.RegionData().MinKeyData.RegionType), service, datas.RegionData().MinKeyData)
	tent = newRealm(realmface.GetRealmId(1, 1, datas.RegionData().MinKeyData.RegionType), service, datas.RegionData().MinKeyData)
	realmArray[0], realmArray[1] = home, tent

	home.directExecFunc = true
	tent.directExecFunc = true
	//go home.loop()
	//go tent.loop()

	return
}

func TestFmt(t *testing.T) {
	a0()
	a0(1)
	a0(1, 2)
	a0(1, 2, 3)
}

func a0(aaa ...int) {
	a(aaa)
}

func a(arr []int) {
	fmt.Println(arr)
}

func TestBaseRestoreProsperity(t *testing.T) {
	RegisterTestingT(t)

	b := &heroHome{}

	now := time.Now()
	d := time.Second

	// 第一次返回false
	Ω(b.tryRestoreProsperity(now, d)).Should(BeFalse())
	Ω(b.nextRestoreProsperityTime).Should(Equal(now.Add(d)))

	// 时间没到，返回false
	Ω(b.tryRestoreProsperity(now.Add(d-1), d)).Should(BeFalse())
	Ω(b.nextRestoreProsperityTime).Should(Equal(now.Add(d)))

	// 时间到了，返回true
	Ω(b.tryRestoreProsperity(now.Add(d), d)).Should(BeTrue())
	Ω(b.nextRestoreProsperityTime).Should(Equal(now.Add(2 * d)))

	// 时间一下跳多次，返回true，以ctime为准
	Ω(b.tryRestoreProsperity(now.Add(9*d+1), d)).Should(BeTrue())
	Ω(b.nextRestoreProsperityTime).Should(Equal(now.Add(10*d + 1)))
}

func getBaseForTest(r *Realm, heroId int64) (base *baseWithData) {
	r.queueFunc(true, func() {
		base = r.getBase(heroId)
	})
	return
}

func setBaseForTest(r *Realm, base *baseWithData) {
	r.queueFunc(true, func() {
		r.addBase(base)
	})
	return
}

//func TestRealm_SpeedUp(t *testing.T) {
//	RegisterTestingT(t)
//
//	for i := 0; i <= 10; i++ {
//		Ω(calCurStartRate(0, 10, int64(i), 0)).Should(Equal(float64(i) / 10))
//	}
//
//	Ω(calCurStartRate(50, 100, 75, 0.5)).Should(Equal(0.75))
//
//	Ω(calculateClientStartTime(50, 100, 0.5)).Should(Equal(int32(0)))
//
//	Ω(calculateClientStartTime(110, 200, 0.1)).Should(Equal(int32(100)))
//}
