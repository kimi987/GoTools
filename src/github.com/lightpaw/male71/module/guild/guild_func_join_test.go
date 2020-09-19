package guild

import (
	"github.com/lightpaw/male7/config/guild_data"
	"github.com/lightpaw/male7/entity"
	"github.com/lightpaw/male7/gen/ifacemock"
	"github.com/lightpaw/male7/gen/pb/guild"
	"github.com/lightpaw/male7/mock"
	"github.com/lightpaw/male7/service/conflict/guilddataservice/sharedguilddata"
	"github.com/lightpaw/male7/util/u64"
	. "github.com/onsi/gomega"
	"testing"
	"time"
	"github.com/lightpaw/male7/util/idbytes"
)

// 入盟测试
func TestJoin(t *testing.T) {
	RegisterTestingT(t)

	fn := newFunc()
	//ctime := fn.time.CurrentTime()

	gctx := &sharedguilddata.GuildContext{
		OperType:     sharedguilddata.ReplyJoinGuild,
		OperatorId:   101,
		OperatorName: idbytes.PlayerName(101),
	}

	testJoinFail(fn)

	testJoinAutoJoin(fn)
	testExistJoinAutoJoin(fn)

	testJoinRejectAutoJoin(fn, false)
	testCancelJoinGuildRequest(fn, false)
	testReplyJoinGuildFail(gctx, fn, false)
	testReplyJoinGuildNo(gctx, fn, false)
	testReplyJoinGuildYes(gctx, fn, false)

	testJoinRejectAutoJoin(fn, true)
	testCancelJoinGuildRequest(fn, true)
	testReplyJoinGuildFail(gctx, fn, true)
	testReplyJoinGuildNo(gctx, fn, true)
	testReplyJoinGuildYes(gctx, fn, true)

	testJoinRequestExpired(fn)
}

func testJoinRequestExpired(fn *guild_func) {
	// 成功加入
	guilds, hero := testJoinRejectAutoJoin(fn, false)
	g := guilds.Get(1)

	Ω(hero.GetJoinGuildIds()).Should(Equal([]int64{g.Id()}))
	Ω(g.GetRequestJoinHeroIds()).Should(Equal([]int64{hero.Id()}))

	// 没到时间
	ctime := fn.time.CurrentTime()
	fn.tickRemoveExpiredJoinRequest(g, ctime)

	Ω(hero.GetJoinGuildIds()).Should(Equal([]int64{g.Id()}))
	Ω(g.GetRequestJoinHeroIds()).Should(Equal([]int64{hero.Id()}))

	// 到时间
	expiredTime := ctime.Add(fn.datas.GuildConfig().JoinRequestDuration)

	fn.tickRemoveExpiredJoinRequest(g, expiredTime)

	Ω(hero.GetJoinGuildIds()).Should(BeEmpty())
	Ω(g.GetRequestJoinHeroIds()).Should(BeEmpty())
}

func testJoinRejectAutoJoin(fn *guild_func, hasOriginGuild bool) (sharedguilddata.Guilds, *entity.Hero) {
	// 没有帮派，加入新帮派

	guilds, hero := newHeroGuild(fn, hasOriginGuild)
	toJoin := guilds.Get(1)
	origin := guilds.Get(2)

	toJoin.SetJoinCondition(true, 0, 0, 0)

	errMsg := fn.userRequestJoin(ifacemock.HeroController, guilds, toJoin.Id())
	Ω(errMsg).Should(BeNil())

	checkHeroJoinGuild(fn, hero, origin, toJoin, false)

	return guilds, hero
}

func testCancelJoinGuildRequest(fn *guild_func, hasOriginGuild bool) {

	// 取消加入帮派申请

	// 帮派不存在
	mock.DefaultHero(entity.NewHero(101, "hero101", fn.datas.HeroInitData(), time.Now()))
	guilds := sharedguilddata.NewGuilds(fn.datas)
	errMsg := fn.userCancelRequestJoin(ifacemock.HeroController, guilds, 1)
	Ω(errMsg).Should(BeEquivalentTo(guild.ErrUserCancelJoinRequestFailInvalidId))

	// 成功加入
	guilds, hero := testJoinRejectAutoJoin(fn, hasOriginGuild)
	g := guilds.Get(1)

	Ω(hero.GetJoinGuildIds()).Should(Equal([]int64{g.Id()}))
	Ω(g.GetRequestJoinHeroIds()).Should(Equal([]int64{hero.Id()}))

	// 成功取消
	errMsg = fn.userCancelRequestJoin(ifacemock.HeroController, guilds, 1)
	Ω(errMsg).Should(BeNil())

	Ω(hero.GetJoinGuildIds()).Should(BeEmpty())
	Ω(g.GetInvateHeroIds()).Should(BeEmpty())
}

func testReplyJoinGuildFail(gctx *sharedguilddata.GuildContext, fn *guild_func, hasOriginGuild bool) {
	// 成功加入
	guilds, target := testJoinRejectAutoJoin(fn, hasOriginGuild)
	g := guilds.Get(1)
	origin := guilds.Get(2)

	ctime := fn.time.CurrentTime()
	self := sharedguilddata.NewMember(101, fn.datas.GuildClassLevelData().MinKeyData, ctime)
	g.AddMember(self)

	// 权限不足
	successMsg, errMsg, broadcastChanged := fn.guildReplyJoinRequest(gctx, guilds, g, self, target.Id(), false)
	Ω(successMsg).Should(BeNil())
	Ω(errMsg).Should(BeEquivalentTo(guild.ErrGuildReplyJoinRequestFailDeny))
	Ω(broadcastChanged).Should(BeFalse())

	self.SetClassLevelData(fn.datas.GuildClassLevelData().MaxKeyData)

	// 帮派成员
	successMsg, errMsg, broadcastChanged = fn.guildReplyJoinRequest(gctx, guilds, g, self, self.Id(), false)
	Ω(successMsg).Should(BeNil())
	Ω(errMsg).Should(BeEquivalentTo(guild.ErrGuildReplyJoinRequestFailInvalidRequest))
	Ω(broadcastChanged).Should(BeFalse())

	// 不在申请列表
	successMsg, errMsg, broadcastChanged = fn.guildReplyJoinRequest(gctx, guilds, g, self, -1, false)
	Ω(successMsg).Should(BeNil())
	Ω(errMsg).Should(BeEquivalentTo(guild.ErrGuildReplyJoinRequestFailInvalidRequest))
	Ω(broadcastChanged).Should(BeFalse())

	// 对方是盟主
	if hasOriginGuild {
		origin.SetLeader(target.Id())
		m := origin.GetMember(target.Id())
		Ω(m).ShouldNot(BeNil())
		m.SetClassLevelData(fn.datas.GuildClassLevelData().MaxKeyData)

		successMsg, errMsg, broadcastChanged = fn.guildReplyJoinRequest(gctx, guilds, g, self, target.Id(), true)
		Ω(successMsg).Should(BeNil())
		Ω(errMsg).Should(BeEquivalentTo(guild.ErrGuildReplyJoinRequestFailInvalidRequest))
		Ω(broadcastChanged).Should(BeFalse())

		Ω(target.GetJoinGuildIds()).Should(BeEmpty())
		Ω(g.GetRequestJoinHeroIds()).Should(BeEmpty())
	}
}

func testReplyJoinGuildNo(gctx *sharedguilddata.GuildContext, fn *guild_func, hasOriginGuild bool) {

	// 拒绝加入帮派

	// 成功加入
	guilds, target := testJoinRejectAutoJoin(fn, hasOriginGuild)
	g := guilds.Get(1)
	origin := guilds.Get(2)

	Ω(target.GetJoinGuildIds()).Should(Equal([]int64{g.Id()}))
	Ω(g.GetRequestJoinHeroIds()).Should(Equal([]int64{target.Id()}))

	ctime := fn.time.CurrentTime()
	self := sharedguilddata.NewMember(101, fn.datas.GuildClassLevelData().MaxKeyData, ctime)
	g.AddMember(self)

	// 成功拒绝
	successMsg, errMsg, broadcastChanged := fn.guildReplyJoinRequest(gctx, guilds, g, self, target.Id(), false)
	Ω(successMsg.Buffer()).Should(BeEquivalentTo(guild.NewS2cGuildReplyJoinRequestMsg(target.IdBytes(), false).Buffer()))
	Ω(errMsg).Should(BeNil())
	Ω(broadcastChanged).Should(BeTrue())

	Ω(target.GetJoinGuildIds()).Should(BeEmpty())
	Ω(g.GetInvateHeroIds()).Should(BeEmpty())

	checkHeroInGuild(target, origin, true)
	checkHeroInGuild(target, g, false)
}

func testReplyJoinGuildYes(gctx *sharedguilddata.GuildContext, fn *guild_func, hasOriginGuild bool) {

	// 同意加入帮派

	// 成功加入
	guilds, target := testJoinRejectAutoJoin(fn, hasOriginGuild)
	g := guilds.Get(1)
	origin := guilds.Get(2)

	Ω(target.GetJoinGuildIds()).Should(Equal([]int64{g.Id()}))
	Ω(g.GetRequestJoinHeroIds()).Should(Equal([]int64{target.Id()}))

	ctime := fn.time.CurrentTime()
	self := sharedguilddata.NewMember(101, fn.datas.GuildClassLevelData().MaxKeyData, ctime)
	g.AddMember(self)

	// 同意加入
	successMsg, errMsg, broadcastChanged := fn.guildReplyJoinRequest(gctx, guilds, g, self, target.Id(), true)
	if !hasOriginGuild {
		Ω(errMsg).Should(BeNil())
		Ω(successMsg).Should(BeEquivalentTo(guild.NewS2cGuildReplyJoinRequestMsg(target.IdBytes(), true)))
		Ω(broadcastChanged).Should(BeTrue())

		checkHeroJoinGuild(fn, target, origin, g, true)
	} else {
		Ω(successMsg).Should(BeNil())
		Ω(errMsg).Should(BeEquivalentTo(guild.ErrGuildReplyJoinRequestFailInvalidRequest))
	}

}

func testJoinAutoJoin(fn *guild_func) {
	// 没有帮派，加入新帮派
	testJoinAutoJoin000(fn, false)
}

func testExistJoinAutoJoin(fn *guild_func) {
	// 有帮派，加入新帮派
	testJoinAutoJoin000(fn, true)
}

func testJoinAutoJoin000(fn *guild_func, hasOriginGuild bool) sharedguilddata.Guilds {
	// 从已有的帮派跳到另外的帮派

	guilds, hero := newHeroGuild(fn, hasOriginGuild)
	toJoin := guilds.Get(1)
	origin := guilds.Get(2)

	errMsg := fn.userRequestJoin(ifacemock.HeroController, guilds, toJoin.Id())
	Ω(errMsg).Should(BeNil())

	checkHeroJoinGuild(fn, hero, origin, toJoin, true)

	return guilds
}

func newHeroGuild(fn *guild_func, hasOriginGuild bool) (sharedguilddata.Guilds, *entity.Hero) {
	ctime := fn.time.CurrentTime()

	hero := entity.NewHero(1, "hero1", fn.datas.HeroInitData(), ctime)

	var origin *sharedguilddata.Guild
	if hasOriginGuild {
		origin = sharedguilddata.NewGuild(2, "name", "flag", fn.datas, ctime)
		self := sharedguilddata.NewMember(1, fn.datas.GuildClassLevelData().MinKeyData, ctime)
		origin.AddMember(self)
		hero.SetGuild(2)
	}

	g := sharedguilddata.NewGuild(1, "name", "flag", fn.datas, ctime)

	guilds := sharedguilddata.NewGuilds(fn.datas)
	if origin != nil {
		guilds.Add(origin)
	}
	guilds.Add(g)

	hero.SetCountryId(g.CountryId())
	mock.DefaultHero(hero)

	return guilds, hero
}

func checkHeroJoinGuild(fn *guild_func, hero *entity.Hero, origin, toJoin *sharedguilddata.Guild, assentJoined bool) {

	if assentJoined {
		checkHeroInGuild(hero, toJoin, true)
		checkHeroInGuild(hero, origin, false)

		Ω(hero.GetJoinGuildIds()).Should(BeEmpty())
	} else {
		checkHeroInGuild(hero, toJoin, false)
		checkHeroInGuild(hero, origin, true)

		Ω(toJoin.GetMember(hero.Id())).Should(BeNil())
		Ω(toJoin.GetRequestJoinHeroIds()).Should(Equal([]int64{hero.Id()}))
		Ω(hero.GetJoinGuildIds()).Should(Equal([]int64{toJoin.Id()}))
	}
}

func checkHeroInGuild(hero *entity.Hero, g *sharedguilddata.Guild, assertInGuild bool) {

	if assertInGuild {
		if g == nil {
			Ω(hero.GuildId()).Should(Equal(int64(0)))
		} else {
			Ω(hero.GuildId()).Should(Equal(g.Id()))
			Ω(g.GetMember(hero.Id())).ShouldNot(BeNil())
		}
	} else {
		if g != nil {
			Ω(hero.GuildId()).ShouldNot(Equal(g.Id()))
			Ω(g.GetMember(hero.Id())).Should(BeNil())
		}
	}
}

func testJoinFail(fn *guild_func) {

	ctime := fn.time.CurrentTime()

	hero := entity.NewHero(1, "hero1", fn.datas.HeroInitData(), ctime)
	mock.DefaultHero(hero)

	guilds := sharedguilddata.NewGuilds(fn.datas)

	// 帮派不存在
	errMsg := fn.userRequestJoin(ifacemock.HeroController, guilds, 1)
	Ω(errMsg).Should(BeEquivalentTo(guild.ErrUserRequestJoinFailInvalidId))

	g := sharedguilddata.NewGuild(1, "name", "flag", fn.datas, ctime)
	guilds.Add(g)

	// A类帮派
	g.SetNpcTemplate(&guild_data.NpcGuildTemplate{
		RejectUserJoin: true,
	})
	errMsg = fn.userRequestJoin(ifacemock.HeroController, guilds, 1)
	Ω(errMsg).Should(BeEquivalentTo(guild.ErrUserRequestJoinFailNpc))
	g.SetNpcTemplate(nil)

	// 帮派已满
	for i := uint64(0); i < g.LevelData().MemberCount; i++ {
		g.AddMember(sharedguilddata.NewMember(int64(i+1), fn.datas.GuildClassLevelData().MinKeyData, ctime))
	}
	errMsg = fn.userRequestJoin(ifacemock.HeroController, guilds, 1)
	Ω(errMsg).Should(BeEquivalentTo(guild.ErrUserRequestJoinFailFull))
	for i := uint64(0); i < g.LevelData().MemberCount; i++ {
		g.RemoveMember(int64(i+1), 0)
	}

	heroLevel := uint64(10)
	junxianLevel := uint64(0) // TODO
	towerFloor := uint64(30)
	g.SetJoinCondition(false, heroLevel, junxianLevel, towerFloor)
	// TODO 军衔

	// 英雄等级不满足
	errMsg = fn.userRequestJoin(ifacemock.HeroController, guilds, 1)
	Ω(errMsg).Should(BeEquivalentTo(guild.ErrUserRequestJoinFailCondition))

	for i := uint64(0); i < heroLevel; i++ {
		hero.AddExp(hero.LevelData().Sub.UpgradeExp)
	}

	// 千重楼最高层数不满足
	errMsg = fn.userRequestJoin(ifacemock.HeroController, guilds, 1)
	Ω(errMsg).Should(BeEquivalentTo(guild.ErrUserRequestJoinFailCondition))
	for i := uint64(0); i < towerFloor; i++ {
		hero.Tower().IncreseCurrentFloor(ctime, fn.datas.MiscConfig().TowerAutoKeepFloor)
	}

	// 允许自动加入

	// 自己是别的帮派的盟主
	g2 := sharedguilddata.NewGuild(2, "name2", "flag2", fn.datas, ctime)
	guilds.Add(g2)
	self := sharedguilddata.NewMember(hero.Id(), fn.datas.GuildClassLevelData().MaxKeyData, ctime)
	g2.AddMember(self)
	g2.SetLeader(self.Id())
	hero.SetGuild(2)
	hero.SetCountryId(g.CountryId())

	errMsg = fn.userRequestJoin(ifacemock.HeroController, guilds, 1)
	Ω(errMsg).Should(BeEquivalentTo(guild.ErrUserRequestJoinFailLeader))

	self.SetClassLevelData(fn.datas.GuildClassLevelData().MinKeyData)
	g2.SetLeader(0)
	//guilds.Remove(g2.Id())
	//g2.RemoveMember(hero.Id())
	//hero.SetGuild(0)

	// 不允许自动加入
	g2.SetJoinCondition(true, 0, 0, 0)

	// 不能申请自己的联盟
	errMsg = fn.userRequestJoin(ifacemock.HeroController, guilds, 2)
	Ω(errMsg).Should(BeEquivalentTo(guild.ErrUserRequestJoinFailSelfGuild))

	// 自己的申请入帮申请已达上限
	g.SetJoinCondition(true, heroLevel, junxianLevel, towerFloor)
	for i := uint64(0); i < fn.datas.GuildConfig().UserMaxJoinRequestCount; i++ {
		hero.AddJoinGuildIds(int64(100 + i))
	}
	errMsg = fn.userRequestJoin(ifacemock.HeroController, guilds, 1)
	Ω(errMsg).Should(BeEquivalentTo(guild.ErrUserRequestJoinFailSelfFull))

	for i := uint64(0); i < fn.datas.GuildConfig().UserMaxJoinRequestCount; i++ {
		hero.RemoveJoinGuildIds(int64(100 + i))
	}
}

func TestUpdateJoinCondition(t *testing.T) {
	RegisterTestingT(t)

	fn := newFunc()

	ctime := fn.time.CurrentTime()

	g := sharedguilddata.NewGuild(1, "name", "flag", fn.datas, ctime)
	self := sharedguilddata.NewMember(1, fn.datas.GuildClassLevelData().MinKeyData, ctime)
	g.AddMember(self)

	// npc盟主
	g.SetLeader(-1)
	successMsg, errMsg, broadcastChanged := fn.updateJoinCondition(g, self, true, 2, 3, 4)
	Ω(successMsg).Should(BeNil())
	Ω(errMsg).Should(Equal(guild.ErrUpdateJoinConditionFailNpc))
	Ω(broadcastChanged).Should(Equal(false))

	// 没权限
	g.SetLeader(0)
	successMsg, errMsg, broadcastChanged = fn.updateJoinCondition(g, self, true, 2, 3, 4)
	Ω(successMsg).Should(BeNil())
	Ω(errMsg).Should(Equal(guild.ErrUpdateJoinConditionFailDeny))
	Ω(broadcastChanged).Should(Equal(false))

	// 盟主权限
	self.SetClassLevelData(fn.datas.GuildClassLevelData().MaxKeyData)
	g.SetLeader(self.Id())

	// 成功
	testUpdateJoinCondition(fn, g, self, true, 2, 3, 4)

}

func testUpdateJoinCondition(fn *guild_func, g *sharedguilddata.Guild, self *sharedguilddata.GuildMember,
	rejectAutoJoin bool, requiredHeroLevel, requiredJunXianLevel, requiredTowerMaxFloor uint64) {
	successMsg, errMsg, broadcastChanged := fn.updateJoinCondition(g, self, rejectAutoJoin, requiredHeroLevel, requiredJunXianLevel, requiredTowerMaxFloor)
	Ω(successMsg.Buffer()).Should(BeEquivalentTo(guild.NewS2cUpdateJoinConditionMsg(rejectAutoJoin, u64.Int32(requiredHeroLevel), u64.Int32(requiredJunXianLevel), u64.Int32(requiredTowerMaxFloor)).Buffer()))
	Ω(errMsg).Should(BeNil())
	Ω(broadcastChanged).Should(Equal(true))

	Ω(g.IsRejectAutoJoin()).Should(BeTrue())
	Ω(g.GetRequiredHeroLevel()).Should(Equal(uint64(2)))
	Ω(g.GetRequiredJunXianLevel()).Should(Equal(uint64(3)))
	Ω(g.GetRequiredTowerMaxFloor()).Should(Equal(uint64(4)))
}
