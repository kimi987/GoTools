package guild

import (
	"github.com/lightpaw/male7/entity"
	"github.com/lightpaw/male7/gen/ifacemock"
	"github.com/lightpaw/male7/mock"
	"github.com/lightpaw/male7/service/conflict/guilddataservice/sharedguilddata"
	"github.com/lightpaw/male7/service/herosnapshot/snapshotdata"
	"github.com/lightpaw/male7/util/u64"
	. "github.com/onsi/gomega"
	"testing"
	"time"
)

// 每日重置
func TestResetDaily(t *testing.T) {
	RegisterTestingT(t)

	fn := newFunc()

	ctime := fn.time.CurrentTime()
	g := sharedguilddata.NewGuild(1, "g1", "f1", fn.datas, ctime)
	self := sharedguilddata.NewMember(1, fn.datas.GuildClassLevelData().MinKeyData, ctime)
	npc := sharedguilddata.NewMember(-1, fn.datas.GuildClassLevelData().MinKeyData, ctime)

	g.AddMember(self)
	g.AddMember(npc)

	total := uint64(0)
	for i := 0; i < fn.datas.GuildConfig().ContributionDay-1; i++ {
		toAdd := uint64(i + 1)
		total += toAdd
		self.AddContribution(toAdd)
		self.AddDonation(toAdd)

		Ω(self.ContributionAmount()).Should(Equal(toAdd))
		Ω(self.ContributionAmount7()).Should(Equal(total))
		Ω(self.DonationAmount()).Should(Equal(toAdd))
		Ω(self.DonationAmount7()).Should(Equal(total))

		ctime = ctime.Add(24 * time.Hour)
		fn.resetDaily(g, ctime)

		Ω(self.ContributionAmount()).Should(Equal(uint64(0)))
		Ω(self.ContributionAmount7()).Should(Equal(total))
		Ω(self.DonationAmount()).Should(Equal(uint64(0)))
		Ω(self.DonationAmount7()).Should(Equal(total))
	}

	// 第7天
	toAdd := uint64(10)
	total += toAdd - 1 // 减去第一天的
	self.AddContribution(toAdd)
	self.AddDonation(toAdd)

	ctime = ctime.Add(24 * time.Hour)
	fn.resetDaily(g, ctime)

	Ω(self.ContributionAmount()).Should(Equal(uint64(0)))
	Ω(self.ContributionAmount7()).Should(Equal(total))
	Ω(self.DonationAmount()).Should(Equal(uint64(0)))
	Ω(self.DonationAmount7()).Should(Equal(total))

	Ω(npc.ContributionAmount()).Should(Equal(uint64(0)))
	Ω(npc.ContributionAmount7()).Should(Equal(uint64(0)))
	Ω(npc.DonationAmount()).Should(Equal(uint64(0)))
	Ω(npc.DonationAmount7()).Should(Equal(uint64(0)))
}

// npc设置职位
func TestNpcSetClass(t *testing.T) {
	RegisterTestingT(t)

	//fn := newFunc()
	//
	//guilds := sharedguilddata.NewGuilds()
	//Ω(guilds.IdRankArray()).Should(BeEmpty())
	//
	//fn.tryKeepFreeNpcGuild(guilds)
	//Ω(len(guilds.IdRankArray())).Should(Equal(u64.Int(fn.datas.GuildConfig().FreeNpcGuildKeepCount)))
	//
	//g := guilds.IdRankArray()[0]
	//
	//ctime := fn.time.CurrentTime()
	//self := sharedguilddata.NewMember(1, fn.datas.GuildClassLevelData().MinKeyData, ctime.Add(1))
	//g.AddMember(self)
	//
	//// 我没有贡献，没有职位
	//totalMemberCount := uint64(0)
	//for _, d := range g.LevelData().ClassMemberCount {
	//	totalMemberCount += d
	//}
	//
	//if uint64(g.MemberCount()) > totalMemberCount {
	//	fn.npcSetMemberClass(g)
	//	Ω(self.ClassLevelData()).Should(Equal(fn.datas.GuildClassLevelData().MinKeyData))
	//}
	//
	//for i, d := range fn.datas.GuildClassLevelData().Array {
	//
	//	if i == 0 || i+1 == len(fn.datas.GuildClassLevelData().Array) {
	//		continue
	//	}
	//
	//	amount := uint64(0)
	//	g.WalkMember(func(member *sharedguilddata.GuildMember) {
	//		if member.ClassLevelData() == d {
	//			if amount == 0 {
	//				amount = member.ContributionAmount7()
	//			} else {
	//				amount = u64.Min(amount, member.ContributionAmount7())
	//			}
	//		}
	//	})
	//
	//	levelData := d
	//	g.WalkMember(func(member *sharedguilddata.GuildMember) {
	//		// 这个值比上一级领导的贡献都大，应该升到上一级去
	//		if member.ClassLevelData() != fn.datas.GuildClassLevelData().MaxKeyData &&
	//			levelData.Level < member.ClassLevelData().Level {
	//			if amount >= member.ContributionAmount7() {
	//				levelData = member.ClassLevelData()
	//			}
	//		}
	//	})
	//
	//	self.AddContribution(u64.Sub(amount, self.ContributionAmount()))
	//
	//	fn.npcSetMemberClass(g)
	//	Ω(self.ClassLevelData()).Should(Equal(levelData))
	//}
	//
	//// 你的贡献度就是超过帮主，也不会把你变成帮主
	//leader := g.GetMember(g.LeaderId())
	//self.AddContribution(leader.ContributionAmount7() + 1)
	//
	//fn.npcSetMemberClass(g)
	//Ω(self.ClassLevelData()).Should(Equal(fn.datas.GuildClassLevelData().Array[len(fn.datas.GuildClassLevelData().Array)-2]))
}

// npc踢人
func TestNpcKick(t *testing.T) {
	RegisterTestingT(t)

	fn := newFunc()

	ctime := fn.time.CurrentTime()

	guilds := sharedguilddata.NewGuilds(fn.datas)
	Ω(guilds.IdRankArray()).Should(BeEmpty())

	fn.tryKeepFreeNpcGuild(guilds)
	Ω(len(guilds.IdRankArray())).Should(Equal(u64.Int(fn.datas.GuildConfig().FreeNpcGuildKeepCount)))

	g := guilds.IdRankArray()[0]

	originMemberCount := g.MemberCount()

	// 帮派不满，不踢人
	fn.npcTryKickMember(g, ctime)
	Ω(g.MemberCount()).Should(Equal(originMemberCount))

	// 满人
	heroMap := make(map[int64]*snapshotdata.HeroSnapshot)
	n := g.EmptyMemberCount()
	for i := uint64(0); i < n; i++ {
		heroId := int64(i + 1)
		toAdd := sharedguilddata.NewMember(heroId, fn.datas.GuildClassLevelData().MinKeyData, ctime)
		g.AddMember(toAdd)

		// 贡献从大到小
		toAdd.AddContribution(n - i)

		// 离线时间从小到大
		heroMap[int64(heroId)] = &snapshotdata.HeroSnapshot{
			LastOfflineTime: ctime.Add(time.Duration(i)),
		}

		hero := entity.NewHero(heroId, "name", fn.datas.HeroInitData(), ctime)
		hero.SetGuild(g.Id())
		mock.DefaultHero(hero)
	}

	ifacemock.HeroSnapshotService.Mock(ifacemock.HeroSnapshotService.Get, func(id int64) *snapshotdata.HeroSnapshot {
		return heroMap[id]
	})

	// 离线时间未达标，不踢人
	fn.npcTryKickMember(g, ctime)
	Ω(g.MemberCount()).Should(Equal(int(g.LevelData().MemberCount)))

	// 只有一个人超过离线最大时间

	removed := g.GetMember(1)

	ctime = ctime.Add(fn.datas.GuildConfig().NpcKickOfflineDuration)
	fn.npcTryKickMember(g, ctime)
	Ω(g.MemberCount()).Should(Equal(int(g.LevelData().MemberCount) - 1))
	Ω(g.GetMember(removed.Id())).Should(BeNil())

	// 加回去，有3个人超过最大离线时间，第三个被踢掉，贡献值最低
	g.AddMember(removed)

	removed = g.GetMember(3)

	ctime = ctime.Add(2)
	fn.npcTryKickMember(g, ctime)
	Ω(g.MemberCount()).Should(Equal(int(g.LevelData().MemberCount) - 1))
	Ω(g.GetMember(removed.Id())).Should(BeNil())
}

// 保持npc帮派个数
func TestKeepNpc(t *testing.T) {
	RegisterTestingT(t)

	fn := newFunc()

	// 当前不需要新增NPC帮派
	guilds := sharedguilddata.NewGuilds(fn.datas)
	Ω(guilds.IdRankArray()).Should(BeEmpty())

	fn.tryKeepFreeNpcGuild(guilds)
	Ω(len(guilds.IdRankArray())).Should(Equal(u64.Int(fn.datas.GuildConfig().FreeNpcGuildKeepCount)))
	Ω(fn.freeNpcGuildChanged).Should(BeFalse())

	// 再调用一次，还是这么多个
	fn.doUpdateFreeNpcGuildChanged()
	fn.tryKeepFreeNpcGuild(guilds)

	Ω(len(guilds.IdRankArray())).Should(Equal(u64.Int(fn.datas.GuildConfig().FreeNpcGuildKeepCount)))
	Ω(fn.freeNpcGuildChanged).Should(BeFalse())

	ctime := fn.time.CurrentTime()

	// 其中一个加人，变成不是空虚的联盟
	g := guilds.IdRankArray()[0]
	n := u64.Sub(g.EmptyMemberCount(), fn.datas.GuildConfig().FreeNpcGuildEmptyCount)
	for i := uint64(0); i < n; i++ {
		toAdd := sharedguilddata.NewMember(int64(i+1), fn.datas.GuildClassLevelData().MinKeyData, ctime)
		g.AddMember(toAdd)
	}

	// 新增了一个联盟
	fn.doUpdateFreeNpcGuildChanged()
	fn.tryKeepFreeNpcGuild(guilds)

	Ω(len(guilds.IdRankArray())).Should(Equal(u64.Int(fn.datas.GuildConfig().FreeNpcGuildKeepCount + 1)))
	Ω(fn.freeNpcGuildChanged).Should(BeFalse())
}
