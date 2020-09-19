package guild

import (
	"github.com/lightpaw/male7/gen/ifacemock"
	"github.com/lightpaw/male7/gen/pb/guild"
	"github.com/lightpaw/male7/mock"
	"github.com/lightpaw/male7/service/conflict/guilddataservice/sharedguilddata"
	"github.com/lightpaw/male7/service/herosnapshot/snapshotdata"
	. "github.com/onsi/gomega"
	"testing"
	"time"
)

func TestImpeachNpcLeader(t *testing.T) {
	// 弹劾Npc盟主

	RegisterTestingT(t)

	fn := newFunc()

	testImpeachNpcLeaderFail(fn)
	testImpeachNpcLeaderTick(fn)
	testImpeachNpcLeaderVote(fn)
}

func testImpeachNpcLeaderFail(fn *guild_func) {
	// 弹劾Npc盟主

	ctime := fn.time.CurrentTime()

	g := sharedguilddata.NewGuild(1, "name", "flag", fn.datas, ctime)
	self := sharedguilddata.NewMember(1, fn.datas.GuildClassLevelData().MinKeyData, ctime)
	self.AddContribution(100)
	g.AddMember(self)

	leader := sharedguilddata.NewMember(-1, fn.datas.GuildClassLevelData().MaxKeyData, ctime)
	g.SetLeader(leader.Id())
	g.AddMember(leader)

	// 现在Npc就是不能弹劾
	if !self.ClassLevelData().Permission.ImpeachNpcLeader {
		// 这个职位不能弹劾
		successMsg, errMsg, broadcastChanged := fn.impeachLeader(g, self)
		Ω(successMsg).Should(BeNil())
		Ω(errMsg).Should(Equal(guild.ErrImpeachLeaderFailDeny))
		Ω(broadcastChanged).Should(Equal(false))

		array := fn.datas.GetGuildClassLevelDataArray()
		self.SetClassLevelData(array[len(array)-2])
	}

	// 现在NPC就是不能弹劾，不管有没有权限
	successMsg, errMsg, broadcastChanged := fn.impeachLeader(g, self)
	Ω(successMsg).Should(BeNil())
	Ω(errMsg).Should(Equal(guild.ErrImpeachLeaderFailDeny))
	Ω(broadcastChanged).Should(Equal(false))

	//// 超过今日弹劾时间
	//endTime := fn.datas.GuildConfig().GetImpeachLeaderTime(ctime, true)
	//mock.SetTime(endTime.Add(1))
	//
	//successMsg, errMsg, broadcastChanged := fn.impeachLeader(g, self)
	//Ω(successMsg).Should(BeNil())
	//Ω(errMsg).Should(Equal(guild.ErrImpeachLeaderFailInvalidTime))
	//Ω(broadcastChanged).Should(Equal(false))

}

func testImpeachNpcLeader(fn *guild_func) *sharedguilddata.Guild {
	// 弹劾Npc盟主

	ctime := fn.time.CurrentTime()

	g := sharedguilddata.NewGuild(1, "name", "flag", fn.datas, ctime)
	self := sharedguilddata.NewMember(1, fn.datas.GuildClassLevelData().MinKeyData, ctime)
	self.AddContribution(100)
	g.AddMember(self)

	// 加2个人，当候选人
	candidate := sharedguilddata.NewMember(101, fn.datas.GuildClassLevelData().MinKeyData, ctime)
	candidate.AddContribution(50)
	g.AddMember(candidate)

	noCandidate := sharedguilddata.NewMember(102, fn.datas.GuildClassLevelData().MinKeyData, ctime)
	noCandidate.AddContribution(10)
	g.AddMember(noCandidate)

	leader := sharedguilddata.NewMember(-1, fn.datas.GuildClassLevelData().MaxKeyData, ctime)
	g.SetLeader(leader.Id())
	g.AddMember(leader)

	// 加几个NPC
	for i := 0; i < 3; i++ {
		npc := sharedguilddata.NewMember(int64(-2-i), fn.datas.GuildClassLevelData().MinKeyData, ctime)
		g.AddMember(npc)
	}

	if !self.ClassLevelData().Permission.ImpeachNpcLeader {
		array := fn.datas.GetGuildClassLevelDataArray()
		self.SetClassLevelData(array[len(array)-2])
	}

	endTime := fn.datas.GuildConfig().GetImpeachLeaderTime(ctime, true)
	mock.SetTime(endTime.Add(-1))

	// 现在NPC就是不能弹劾，不管有没有权限
	successMsg, errMsg, broadcastChanged := fn.impeachLeader(g, self)
	Ω(successMsg).Should(BeNil())
	Ω(errMsg).Should(Equal(guild.ErrImpeachLeaderFailDeny))
	Ω(broadcastChanged).Should(Equal(false))

	//successMsg, errMsg, broadcastChanged := fn.impeachLeader(g, self)
	//Ω(successMsg).Should(BeEquivalentTo(guild.IMPEACH_LEADER_S2C))
	//Ω(errMsg).Should(BeNil())
	//Ω(broadcastChanged).Should(Equal(true))
	//
	//Ω(g.GetImpeachLeader()).ShouldNot(BeNil())
	//Ω(g.GetImpeachLeader().IsValidCandidate(leader.Id())).Should(BeTrue())
	//Ω(g.GetImpeachLeader().IsValidCandidate(self.Id())).Should(BeTrue())
	//if fn.datas.GuildConfig().ImpeachExtraCandidateCount == 2 {
	//	Ω(g.GetImpeachLeader().IsValidCandidate(candidate.Id())).Should(BeTrue())
	//	Ω(g.GetImpeachLeader().IsValidCandidate(noCandidate.Id())).Should(BeFalse())
	//}
	//Ω(g.GetImpeachLeader().IsValidCandidate(candidate.Id())).Should(Equal(
	//	fn.datas.GuildConfig().ImpeachExtraCandidateCount > 1,
	//))
	//Ω(g.GetImpeachLeader().IsValidCandidate(noCandidate.Id())).Should(Equal(
	//	fn.datas.GuildConfig().ImpeachExtraCandidateCount > 2,
	//))

	return g
}

func testImpeachNpcLeaderTick(fn *guild_func) {

	g := testImpeachNpcLeader(fn)

	//self := g.GetMember(1)
	leader := g.GetMember(-1)

	ctime := fn.time.CurrentTime()
	endTime := fn.datas.GuildConfig().GetImpeachLeaderTime(ctime, true)

	// 到点tick
	fn.tickUpdateImpeachLeader(g, endTime)
	Ω(g.GetImpeachLeader()).Should(BeNil())

	testCheckLeader(fn, g, leader)
}

func testImpeachNpcLeaderVote(fn *guild_func) {

	g := testImpeachNpcLeader(fn)

	ctime := fn.time.CurrentTime()
	endTime := ctime.Add(time.Hour)
	g.StartImpeachLeader(0, ctime, endTime, fn.datas.GuildConfig().ImpeachExtraCandidateCount)

	self := g.GetMember(1)
	//leader := g.GetMember(-1)

	// 计算npc的票数
	npcVoteScore := uint64(0)
	g.WalkMember(func(member *sharedguilddata.GuildMember) {
		if member.IsNpc() {
			npcVoteScore += member.ClassLevelData().VoteScore
		}
	})

	n := npcVoteScore/fn.datas.GuildClassLevelData().MinKeyData.VoteScore + 1
	for i := uint64(0); i < n; i++ {
		member := sharedguilddata.NewMember(int64(2+i), fn.datas.GuildClassLevelData().MinKeyData, ctime)
		g.AddMember(member)

		vote(fn, g, member, self.Id())
	}

	// 投完票，不会提前结束
	Ω(g.GetImpeachLeader()).ShouldNot(BeNil())

	// 开始tick时间
	fn.tickUpdateImpeachLeader(g, endTime.Add(1))

	Ω(g.GetImpeachLeader()).Should(BeNil())

	testCheckLeader(fn, g, self)

	g.WalkMember(func(member *sharedguilddata.GuildMember) {
		Ω(member.IsNpc()).Should(BeFalse())
	})
}

func TestImpeachUserLeader(t *testing.T) {
	// 弹劾Npc盟主

	RegisterTestingT(t)

	fn := newFunc()

	testImpeachUserLeaderFail(fn)
	testImpeachUserLeaderTickSame(fn)
	testImpeachUserLeaderTickDiffVote(fn)

	testImpeachUserLeaderJoinVote(fn)
}

func testImpeachUserLeaderFail(fn *guild_func) {

	// 弹劾Npc盟主
	ctime := fn.time.CurrentTime()

	g := sharedguilddata.NewGuild(1, "name", "flag", fn.datas, ctime)
	self := sharedguilddata.NewMember(1, fn.datas.GuildClassLevelData().MinKeyData, ctime)
	self.AddContribution(100)
	g.AddMember(self)

	leader := sharedguilddata.NewMember(100, fn.datas.GuildClassLevelData().MaxKeyData, ctime)
	g.AddMember(leader)
	g.SetLeader(leader.Id())

	leaderSnapshot := &snapshotdata.HeroSnapshot{}
	leaderSnapshot.LastOfflineTime = ctime
	ifacemock.HeroSnapshotService.Mock(ifacemock.HeroSnapshotService.Get, func(id int64) *snapshotdata.HeroSnapshot {
		if id == 100 {
			return leaderSnapshot
		}
		return nil
	})

	// 盟主弹劾自己
	successMsg, errMsg, broadcastChanged := fn.impeachLeader(g, leader)
	Ω(successMsg).Should(BeNil())
	Ω(errMsg).Should(Equal(guild.ErrImpeachLeaderFailConditionNotReach))
	Ω(broadcastChanged).Should(Equal(false))

	if g.MemberCount() < fn.datas.GuildConfig().ImpeachUserLeaderMemberCount {
		// 人数不足
		successMsg, errMsg, broadcastChanged = fn.impeachLeader(g, self)
		Ω(successMsg).Should(BeNil())
		Ω(errMsg).Should(Equal(guild.ErrImpeachLeaderFailConditionNotReach))
		Ω(broadcastChanged).Should(Equal(false))
	}

	// 达到人数
	for i := 2; i < fn.datas.GuildConfig().ImpeachUserLeaderMemberCount; i++ {
		noCandidate := sharedguilddata.NewMember(int64(i), fn.datas.GuildClassLevelData().MinKeyData, ctime)
		g.AddMember(noCandidate)
	}

	// 离线时间不足
	successMsg, errMsg, broadcastChanged = fn.impeachLeader(g, self)
	Ω(successMsg).Should(BeNil())
	Ω(errMsg).Should(Equal(guild.ErrImpeachLeaderFailConditionNotReach))
	Ω(broadcastChanged).Should(Equal(false))

	successMsg, errMsg, broadcastChanged = fn.impeachLeader(g, self)
	Ω(successMsg).Should(BeNil())
	Ω(errMsg).Should(Equal(guild.ErrImpeachLeaderFailConditionNotReach))
	Ω(broadcastChanged).Should(Equal(false))
}

func testImpeachUserLeader(fn *guild_func) *sharedguilddata.Guild {

	ctime := fn.time.CurrentTime()

	g := sharedguilddata.NewGuild(1, "name", "flag", fn.datas, ctime)
	self := sharedguilddata.NewMember(1, fn.datas.GuildClassLevelData().MinKeyData, ctime)
	self.AddContribution(100)
	g.AddMember(self)

	leader := sharedguilddata.NewMember(100, fn.datas.GuildClassLevelData().MaxKeyData, ctime)
	g.AddMember(leader)
	g.SetLeader(leader.Id())

	leaderSnapshot := &snapshotdata.HeroSnapshot{}
	leaderSnapshot.LastOfflineTime = ctime.Add(-fn.datas.GuildConfig().ImpeachUserLeaderOffline - 1)

	ifacemock.HeroSnapshotService.Mock(ifacemock.HeroSnapshotService.Get, func(id int64) *snapshotdata.HeroSnapshot {
		if id == 100 {
			return leaderSnapshot
		}
		return nil
	})

	// 加2个人进来，当候选人
	candidate1 := sharedguilddata.NewMember(101, fn.datas.GuildClassLevelData().MinKeyData, ctime.Add(-1))
	candidate1.AddContribution(50)
	g.AddMember(candidate1)

	candidate2 := sharedguilddata.NewMember(102, fn.datas.GuildClassLevelData().MinKeyData, ctime.Add(1))
	candidate2.AddContribution(10)
	g.AddMember(candidate2)

	// 达到弹劾人数
	for i := 4; i < fn.datas.GuildConfig().ImpeachUserLeaderMemberCount; i++ {
		noCandidate := sharedguilddata.NewMember(int64(i), fn.datas.GuildClassLevelData().MinKeyData, ctime)
		g.AddMember(noCandidate)
	}

	// 不会超过今日弹劾时间
	endTime := fn.datas.GuildConfig().GetImpeachLeaderTime(ctime, false)
	Ω(ctime.Before(endTime)).Should(BeTrue())

	// 弹劾成功
	successMsg, errMsg, broadcastChanged := fn.impeachLeader(g, self)
	Ω(successMsg).Should(BeEquivalentTo(guild.IMPEACH_LEADER_S2C))
	Ω(errMsg).Should(BeNil())
	Ω(broadcastChanged).Should(Equal(true))

	Ω(g.GetImpeachLeader()).ShouldNot(BeNil())
	Ω(g.GetImpeachLeader().IsValidCandidate(leader.Id())).Should(BeTrue()) // 原盟主一定在
	Ω(g.GetImpeachLeader().IsValidCandidate(self.Id())).Should(BeTrue())   // 弹劾的人一定在
	if fn.datas.GuildConfig().ImpeachExtraCandidateCount == 3 {
		Ω(g.GetImpeachLeader().IsValidCandidate(candidate1.Id())).Should(BeTrue())
		Ω(g.GetImpeachLeader().IsValidCandidate(candidate2.Id())).Should(BeTrue())
	}
	Ω(g.GetImpeachLeader().IsValidCandidate(candidate1.Id())).Should(Equal(
		fn.datas.GuildConfig().ImpeachExtraCandidateCount > 1,
	))
	Ω(g.GetImpeachLeader().IsValidCandidate(candidate2.Id())).Should(Equal(
		fn.datas.GuildConfig().ImpeachExtraCandidateCount > 2,
	))

	// leader 不足候选人之列
	//Ω(g.GetImpeachLeader().IsValidCandidate(leader.Id())).Should(BeFalse())

	return g
}

func testImpeachUserLeaderTickSame(fn *guild_func) {

	g := testImpeachUserLeader(fn)
	//self := g.GetMember(1)
	leader := g.GetMember(100) // ctime
	//candidate1 := g.GetMember(101) // ctime - 1
	//candidate2 := g.GetMember(102) // ctime + 1

	ctime := fn.time.CurrentTime()
	endTime := fn.datas.GuildConfig().GetImpeachLeaderTime(ctime, false)

	// 到点tick
	fn.tickUpdateImpeachLeader(g, endTime)
	Ω(g.GetImpeachLeader()).Should(BeNil())

	// 票数相同，没有人投票，入盟晚的当选
	testCheckLeader(fn, g, leader)
}

func testImpeachUserLeaderTickDiffVote(fn *guild_func) {

	g := testImpeachUserLeader(fn)
	self := g.GetMember(1) // ctime
	//leader := g.GetMember(100) // ctime
	candidate1 := g.GetMember(101) // ctime - 1
	//candidate2 := g.GetMember(102) // ctime + 1

	vote(fn, g, self, self.Id())
	vote(fn, g, candidate1, candidate1.Id())

	ctime := fn.time.CurrentTime()
	endTime := fn.datas.GuildConfig().GetImpeachLeaderTime(ctime, false)

	// 到点tick
	fn.tickUpdateImpeachLeader(g, endTime)
	Ω(g.GetImpeachLeader()).Should(BeNil())

	// 票数相同，没有人投票，入盟晚的当选
	testCheckLeader(fn, g, self)
}

func testImpeachUserLeaderJoinVote(fn *guild_func) {

	g := testImpeachUserLeader(fn)
	self := g.GetMember(1)

	n := g.MemberCount()/2 + 1

	// 投票投出来
	g.WalkMember(func(member *sharedguilddata.GuildMember) {
		if n > 0 {
			vote(fn, g, member, self.Id())
		}

		n--
	})

	// 投完票，不会提前结束
	Ω(g.GetImpeachLeader()).ShouldNot(BeNil())

	// 开始tick时间
	ctime := fn.time.CurrentTime()
	endTime := fn.datas.GuildConfig().GetImpeachLeaderTime(ctime, false)
	fn.tickUpdateImpeachLeader(g, endTime.Add(1))

	Ω(g.GetImpeachLeader()).Should(BeNil())

	Ω(g.GetImpeachLeader()).Should(BeNil())
	testCheckLeader(fn, g, self)

}

func vote(fn *guild_func, g *sharedguilddata.Guild, self *sharedguilddata.GuildMember, targetId int64) {
	successMsg, errMsg, broadcastChanged := fn.impeachLeaderVote(g, self, targetId)
	Ω(successMsg).ShouldNot(BeNil())
	Ω(errMsg).Should(BeNil())
	Ω(broadcastChanged).Should(Equal(true))
}
