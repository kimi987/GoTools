package guild

import (
	"github.com/lightpaw/male7/config/guild_data"
	"github.com/lightpaw/male7/entity"
	"github.com/lightpaw/male7/gen/ifacemock"
	"github.com/lightpaw/male7/gen/pb/guild"
	"github.com/lightpaw/male7/service/conflict/guilddataservice/sharedguilddata"
	. "github.com/onsi/gomega"
	"testing"
	"github.com/lightpaw/logrus"
	"github.com/lightpaw/male7/util/idbytes"
)

func TestInvate(t *testing.T) {
	RegisterTestingT(t)

	logrus.SetLevel(logrus.DebugLevel)

	fn := newFunc()

	testAllInvate(fn, false)
	testAllInvate(fn, true)
}

func testAllInvate(fn *guild_func, hasOriginGuild bool) {

	gctx := &sharedguilddata.GuildContext{
		OperType:     sharedguilddata.ReplyJoinGuild,
		OperatorId:   101,
		OperatorName: idbytes.PlayerName(101),
	}

	testInvateOther(fn, hasOriginGuild)
	testCancelInvateOther(fn, hasOriginGuild)

	testInvateOtherExpired(fn, hasOriginGuild)

	testReplyInvateSelfFail(gctx, fn, hasOriginGuild)
	testReplyInvateSelfNo(fn, hasOriginGuild)
	testReplyInvateSelfYes(gctx, fn, hasOriginGuild)
}

func testInvateOtherExpired(fn *guild_func, hasOriginGuild bool) {

	guilds, target := testInvateOther(fn, hasOriginGuild)
	g := guilds.Get(1)

	ctime := fn.time.CurrentTime()
	fn.tickRemoveExpiredInvateRequest(g, ctime)

	Ω(target.GetBeenInvateGuildIds()).Should(Equal([]int64{g.Id()}))
	Ω(g.GetInvateHeroIds()).Should(Equal([]int64{target.Id()}))

	expiredTime := ctime.Add(fn.datas.GuildConfig().InvateDuration)
	fn.tickRemoveExpiredInvateRequest(g, expiredTime)

	Ω(target.GetBeenInvateGuildIds()).Should(BeEmpty())
	Ω(g.GetInvateHeroIds()).Should(BeEmpty())
}

func testInvateOther(fn *guild_func, hasOriginGuild bool) (sharedguilddata.Guilds, *entity.Hero) {

	guilds, target := newHeroGuild(fn, hasOriginGuild)
	g := guilds.Get(1)
	//origin := guilds.Get(2)

	ctime := fn.time.CurrentTime()
	self := sharedguilddata.NewMember(101, fn.datas.GuildClassLevelData().MaxKeyData, ctime)
	g.SetLeader(self.Id())
	g.AddMember(self)

	successMsg, errMsg, broadcastChanged := fn.guildInvateOther(guilds, g, self, target.Id())
	Ω(errMsg).Should(BeNil())
	Ω(successMsg).Should(BeEquivalentTo(guild.NewS2cGuildInvateOtherMsg(target.IdBytes())))
	Ω(broadcastChanged).Should(BeTrue())

	Ω(target.GetBeenInvateGuildIds()).Should(Equal([]int64{g.Id()}))
	Ω(g.GetInvateHeroIds()).Should(Equal([]int64{target.Id()}))

	return guilds, target
}

func testCancelInvateOther(fn *guild_func, hasOriginGuild bool) {

	guilds, hero := testInvateOther(fn, hasOriginGuild)
	g := guilds.Get(1)
	leader := g.GetMember(101)

	successMsg, errMsg, broadcastChanged := fn.guildCancelInvateOther(g, leader, hero.Id())

	Ω(successMsg.Buffer()).Should(BeEquivalentTo(guild.NewS2cGuildInvateOtherMsg(hero.IdBytes()).Buffer()))
	Ω(errMsg).Should(BeNil())
	Ω(broadcastChanged).Should(BeTrue())

	Ω(hero.GetBeenInvateGuildIds()).Should(BeEmpty())
	Ω(g.GetInvateHeroIds()).Should(BeEmpty())

}

func testReplyInvateSelfFail(gctx *sharedguilddata.GuildContext, fn *guild_func, hasOriginGuild bool) {

	guilds, hero := testInvateOther(fn, hasOriginGuild)
	g := guilds.Get(1)
	origin := guilds.Get(2)

	// 拒绝不存在的帮派id
	errMsg := fn.userRejectInvateRequest(ifacemock.HeroController, guilds, 100)
	Ω(errMsg).Should(BeEquivalentTo(guild.ErrUserReplyInvateRequestFailInvalidId))

	// 同意没邀请你的帮派
	if origin != nil {
		errMsg = fn.userAgreeInvateRequest(gctx, ifacemock.HeroController, guilds, origin.Id())
		Ω(errMsg).Should(BeEquivalentTo(guild.ErrUserReplyInvateRequestFailInvalidId))
	}

	// 帮派不存在
	guilds.Remove(g.Id())
	errMsg = fn.userAgreeInvateRequest(gctx, ifacemock.HeroController, guilds, g.Id())
	Ω(errMsg).Should(BeEquivalentTo(guild.ErrUserReplyInvateRequestFailInvalidId))

	Ω(hero.GetBeenInvateGuildIds()).Should(BeEmpty())

	testReplyInvateSelfFailNpc(gctx, fn, hasOriginGuild)
	testReplyInvateSelfFailLeader(gctx, fn, hasOriginGuild)
}

func testReplyInvateSelfFailNpc(gctx *sharedguilddata.GuildContext, fn *guild_func, hasOriginGuild bool) {

	guilds, hero := testInvateOther(fn, hasOriginGuild)
	g := guilds.Get(1)

	// 纯npc帮派
	g.SetNpcTemplate(&guild_data.NpcGuildTemplate{
		RejectUserJoin: true,
	})
	errMsg := fn.userAgreeInvateRequest(gctx, ifacemock.HeroController, guilds, g.Id())
	Ω(errMsg).Should(BeEquivalentTo(guild.ErrUserReplyInvateRequestFailInvalidId))

	Ω(hero.GetBeenInvateGuildIds()).Should(BeEmpty())
	Ω(g.GetInvateHeroIds()).Should(BeEmpty())
}

func testReplyInvateSelfFailLeader(gctx *sharedguilddata.GuildContext, fn *guild_func, hasOriginGuild bool) {

	if !hasOriginGuild {
		return
	}

	guilds, hero := testInvateOther(fn, hasOriginGuild)
	g := guilds.Get(1)
	origin := guilds.Get(2)
	origin.GetMember(hero.Id()).SetClassLevelData(fn.datas.GuildClassLevelData().MaxKeyData)
	origin.SetLeader(hero.Id())

	errMsg := fn.userAgreeInvateRequest(gctx, ifacemock.HeroController, guilds, g.Id())
	//Ω(errMsg).Should(BeEquivalentTo(guild.ErrUserReplyInvateRequestFailLeader))
	// 有帮派不能加入
	Ω(errMsg).Should(BeEquivalentTo(guild.ErrUserReplyInvateRequestFailInvalidId))

	// 如果是帮主，那么邀请状态不会取消
	Ω(hero.GetBeenInvateGuildIds()).Should(Equal([]int64{g.Id()}))
	Ω(g.GetInvateHeroIds()).Should(Equal([]int64{hero.Id()}))
}

func testReplyInvateSelfNo(fn *guild_func, hasOriginGuild bool) {

	guilds, hero := testInvateOther(fn, hasOriginGuild)
	g := guilds.Get(1)

	errMsg := fn.userRejectInvateRequest(ifacemock.HeroController, guilds, 1)
	Ω(errMsg).Should(BeNil())

	Ω(hero.GetBeenInvateGuildIds()).Should(BeEmpty())
	Ω(g.GetInvateHeroIds()).Should(BeEmpty())
}

func testReplyInvateSelfYes(gctx *sharedguilddata.GuildContext, fn *guild_func, hasOriginGuild bool) {
	guilds, hero := testInvateOther(fn, hasOriginGuild)
	g := guilds.Get(1)
	origin := guilds.Get(2)

	errMsg := fn.userAgreeInvateRequest(gctx, ifacemock.HeroController, guilds, 1)
	if hasOriginGuild {
		// 有帮派不能加入
		Ω(errMsg).Should(BeEquivalentTo(guild.ErrUserReplyInvateRequestFailInvalidId))

		checkHeroInGuild(hero, origin, true)
		checkHeroInGuild(hero, g, false)
	} else {
		Ω(errMsg).Should(BeNil())

		checkHeroInGuild(hero, origin, false)
		checkHeroInGuild(hero, g, true)

		Ω(hero.GetBeenInvateGuildIds()).Should(BeEmpty())
		Ω(g.GetInvateHeroIds()).Should(BeEmpty())
	}

}
