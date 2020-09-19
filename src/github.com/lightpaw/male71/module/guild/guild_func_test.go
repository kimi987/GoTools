package guild

import (
	"context"
	"github.com/lightpaw/male7/config"
	"github.com/lightpaw/male7/config/confpath"
	"github.com/lightpaw/male7/config/domestic_data"
	"github.com/lightpaw/male7/entity"
	"github.com/lightpaw/male7/gen/ifacemock"
	"github.com/lightpaw/male7/gen/pb/guild"
	"github.com/lightpaw/male7/mock"
	"github.com/lightpaw/male7/pb/shared_proto"
	"github.com/lightpaw/male7/service/conflict/guilddataservice/sharedguilddata"
	"github.com/lightpaw/male7/service/heromodule"
	"github.com/lightpaw/male7/util/concurrent"
	"github.com/lightpaw/male7/util/u64"
	"github.com/lightpaw/pbutil"
	. "github.com/onsi/gomega"
	"github.com/pkg/errors"
	"testing"
	"time"
)

func TestNewGuildFunc(t *testing.T) {
	RegisterTestingT(t)

	guildFunc := newFunc()
	Ω(guildFunc).ShouldNot(BeNil())
}

func TestGuildModule_ProcessCreateGuild(t *testing.T) {
	RegisterTestingT(t)

	fn := newFunc()

	ctime := fn.time.CurrentTime()

	guilds := sharedguilddata.NewGuilds(fn.datas)
	g1 := sharedguilddata.NewGuild(fn.newGuildId(), "name", "flag", fn.datas, ctime)
	g2 := sharedguilddata.NewGuild(fn.newGuildId(), "name2", "flag2", fn.datas, ctime)
	guilds.Add(g1)
	guilds.Add(g2)

	// 名字冲突
	errMsg := fn.createGuild(ifacemock.HeroController, "name", "flag1", guilds)
	Ω(errMsg).Should(Equal(guild.ErrCreateGuildFailNameDuplicate))

	// 旗号冲突
	errMsg = fn.createGuild(ifacemock.HeroController, "name1", "flag2", guilds)
	Ω(errMsg).Should(Equal(guild.ErrCreateGuildFailFlagNameDuplicate))

	// 设置英雄
	hero := entity.NewHero(1, "hero1", fn.datas.HeroInitData(), ctime)
	hero.SetCountryId(fn.datas.GuildConfig().DefaultGuildCountry.Id)
	mock.DefaultHero(hero)

	mock.LockResult.Reset()

	// 玩家存在联盟
	hero.SetGuild(1)
	errMsg = fn.createGuild(ifacemock.HeroController, "name1", "flag1", guilds)
	Ω(errMsg).Should(Equal(guild.ErrCreateGuildFailInTheGuild))

	Ω(mock.LockResult.IsChanged()).Should(BeFalse())
	Ω(mock.LockResult.IsOk()).Should(BeFalse())
	Ω(mock.LockResult.PopMsg()).Should(BeNil())

	hero.SetGuild(0)

	if !fn.datas.GuildConfig().CreateGuildCost.IsEmpty() {
		// 钱不够
		errMsg = fn.createGuild(ifacemock.HeroController, "name1", "flag1", guilds)
		Ω(errMsg).Should(Equal(guild.ErrCreateGuildFailCostNotEnough))

		Ω(mock.LockResult.IsChanged()).Should(BeFalse())
		Ω(mock.LockResult.IsOk()).Should(BeFalse())
		Ω(mock.LockResult.PopMsg()).Should(BeNil())
	}

	// 钱也够了
	mock.SetHeroEnoughCost(hero, fn.datas.GuildConfig().CreateGuildCost)
	Ω(heromodule.HasEnoughCost(hero, fn.datas.GuildConfig().CreateGuildCost)).Should(BeTrue())

	ifacemock.DbService.Mock(ifacemock.DbService.CreateGuild, func(ctx context.Context, id int64, data []byte) error {
		Ω(id).Should(Equal(int64(3)))
		Ω(data).ShouldNot(BeEmpty())
		return errors.Errorf("错误")
	})
	errMsg = fn.createGuild(ifacemock.HeroController, "name1", "flag1", guilds)
	Ω(errMsg).Should(Equal(guild.ErrCreateGuildFailServerError))

	Ω(mock.LockResult.IsChanged()).Should(BeFalse())
	Ω(mock.LockResult.IsOk()).Should(BeFalse())
	Ω(mock.LockResult.PopMsg()).Should(BeNil())

	// 钱还在
	Ω(heromodule.HasEnoughCost(hero, fn.datas.GuildConfig().CreateGuildCost)).Should(BeTrue())

	// 成功
	ifacemock.DbService.Mock(ifacemock.DbService.CreateGuild, func(ctx context.Context, id int64, data []byte) error {
		Ω(id).Should(Equal(int64(4)))
		Ω(data).ShouldNot(BeEmpty())
		return nil
	})

	errMsg = fn.createGuild(ifacemock.HeroController, "name1", "flag1", guilds)
	Ω(errMsg).Should(BeNil())
	Ω(mock.LockResult.IsChanged()).Should(BeTrue())
	Ω(mock.LockResult.IsOk()).Should(BeTrue())

	if !fn.datas.GuildConfig().CreateGuildCost.IsEmpty() {
		// 钱扣掉了
		Ω(heromodule.HasEnoughCost(hero, fn.datas.GuildConfig().CreateGuildCost)).Should(BeFalse())
	}

	// 帮派有了
	Ω(hero.GuildId()).Should(Equal(int64(4)))

	g := guilds.Get(4)
	Ω(g).ShouldNot(BeNil())

	member := g.GetMember(hero.Id())
	Ω(member).ShouldNot(BeNil())

	Ω(member.ClassLevelData()).Should(Equal(fn.datas.GuildConfig().GetLeaderClassLevel()))
}

func TestGuildModule_ProcessLeaveGuild(t *testing.T) {
	RegisterTestingT(t)

	fn := newFunc()

	ctime := fn.time.CurrentTime()

	// 设置英雄
	hero := entity.NewHero(1, "hero1", fn.datas.HeroInitData(), ctime)
	mock.DefaultHero(hero)

	// 玩家存在联盟
	hero.SetGuild(1)

	guilds := sharedguilddata.NewGuilds(fn.datas)
	g := sharedguilddata.NewGuild(fn.newGuildId(), "name", "flag", fn.datas, ctime)
	guilds.Add(g)

	self := sharedguilddata.NewMember(1, fn.datas.GuildClassLevelData().MaxKeyData, ctime)
	g.AddMember(self)
	g.SetLeader(1)

	self2 := sharedguilddata.NewMember(2, fn.datas.GuildClassLevelData().MinKeyData, ctime)
	g.AddMember(self2)

	// 联盟有人，盟主不能退
	successMsg, errMsg, broadcastChanged := fn.leaveGuild(ifacemock.HeroController, guilds, g, self)
	Ω(successMsg).Should(BeNil())
	Ω(errMsg).Should(Equal(guild.ErrLeaveGuildFailLeader))
	Ω(broadcastChanged).Should(Equal(false))

	// 2号离开联盟
	hero2 := entity.NewHero(2, "hero2", fn.datas.HeroInitData(), ctime)
	mock.DefaultHero(hero2)
	hero2.SetGuild(1)

	successMsg, errMsg, broadcastChanged = fn.leaveGuild(ifacemock.HeroController, guilds, g, self2)
	Ω(successMsg).Should(Equal(guild.LEAVE_GUILD_S2C))
	Ω(errMsg).Should(BeNil())
	Ω(broadcastChanged).Should(Equal(true))

	Ω(g.MemberCount()).Should(Equal(1))
	Ω(g.GetMember(hero2.Id())).Should(BeNil())
	Ω(hero2.GuildId()).Should(Equal(int64(0)))

	// 联盟还存在
	Ω(guilds.Get(1)).Should(Equal(g))

	// 1号离开联盟
	mock.DefaultHero(hero)
	successMsg, errMsg, broadcastChanged = fn.leaveGuild(ifacemock.HeroController, guilds, g, self)
	Ω(successMsg).Should(Equal(guild.LEAVE_GUILD_S2C))
	Ω(errMsg).Should(BeNil())
	Ω(broadcastChanged).Should(Equal(true))

	Ω(g.MemberCount()).Should(Equal(0))
	Ω(g.GetMember(hero.Id())).Should(BeNil())
	Ω(hero.GuildId()).Should(Equal(int64(0)))

	// 联盟也不存在了
	Ω(guilds.Get(1)).Should(BeNil())
}

func TestGuildModule_ProcessKickOther(t *testing.T) {
	RegisterTestingT(t)

	fn := newFunc()

	ctime := fn.time.CurrentTime()

	// 设置英雄
	hero := entity.NewHero(1, "hero1", fn.datas.HeroInitData(), ctime)
	hero.SetGuild(1)
	mock.DefaultHero(hero)

	guilds := sharedguilddata.NewGuilds(fn.datas)
	g := sharedguilddata.NewGuild(fn.newGuildId(), "name", "flag", fn.datas, ctime)
	guilds.Add(g)

	self := sharedguilddata.NewMember(1, fn.datas.GuildClassLevelData().MaxKeyData, ctime)
	g.AddMember(self)

	g.SetLeader(-1)
	// Npc联盟不能踢人
	successMsg, errMsg, broadcastChanged := fn.kickOther(g, self, 2)
	Ω(successMsg).Should(BeNil())
	Ω(errMsg).Should(Equal(guild.ErrKickOtherFailNpc))
	Ω(broadcastChanged).Should(Equal(false))

	// 正常联盟
	g.SetLeader(1)

	hero2 := entity.NewHero(2, "hero2", fn.datas.HeroInitData(), ctime)
	hero2.SetGuild(1)
	mock.DefaultHero(hero2)
	other := sharedguilddata.NewMember(2, fn.datas.GuildClassLevelData().MinKeyData, ctime)
	g.AddMember(other)

	// 渣渣都有踢人权限????让策划改表
	Ω(other.ClassLevelData().Permission.KickLowerMember).Should(BeFalse())

	// 没有权限
	successMsg, errMsg, broadcastChanged = fn.kickOther(g, other, 2)
	Ω(successMsg).Should(BeNil())
	Ω(errMsg).Should(Equal(guild.ErrKickOtherFailDeny))
	Ω(broadcastChanged).Should(Equal(false))

	// 目标id无效
	successMsg, errMsg, broadcastChanged = fn.kickOther(g, self, 3)
	Ω(successMsg).Should(BeNil())
	Ω(errMsg).Should(Equal(guild.ErrKickOtherFailNotInGuild))
	Ω(broadcastChanged).Should(Equal(false))

	// 权限不足
	successMsg, errMsg, broadcastChanged = fn.kickOther(g, self, 1)
	Ω(successMsg).Should(BeNil())
	Ω(errMsg).Should(Equal(guild.ErrKickOtherFailDeny))
	Ω(broadcastChanged).Should(Equal(false))

	// TODO 弹劾

	for i := uint64(0); i < fn.datas.GuildConfig().DailyMaxKickCount; i++ {
		// 成功
		Ω(hero2.GuildId()).Should(Equal(int64(1)))

		successMsg, errMsg, broadcastChanged = fn.kickOther(g, self, 2)
		//Ω(successMsg).Should(BeEquivalentTo(guild.NewS2cKickOtherMsg(idbytes.ToBytes(2), idbytes.PlayerName(2))))
		Ω(errMsg).Should(BeNil())
		Ω(broadcastChanged).Should(Equal(true))

		Ω(g.MemberCount()).Should(Equal(1))
		Ω(g.GetMember(hero2.Id())).Should(BeNil())
		Ω(hero2.GuildId()).Should(Equal(int64(0)))

		Ω(g.GetKickMemberCount()).Should(Equal(i + 1))

		// 再次加入，再踢一次
		hero2.SetGuild(1)
		g.AddMember(other)
	}

	// 每日踢人上限
	successMsg, errMsg, broadcastChanged = fn.kickOther(g, self, 2)
	Ω(successMsg).Should(BeNil())
	Ω(errMsg).Should(Equal(guild.ErrKickOtherFailLimit))
	Ω(broadcastChanged).Should(Equal(false))
}

func TestGuildModule_ProcessUpdateText(t *testing.T) {
	RegisterTestingT(t)

	fn := newFunc()

	ctime := fn.time.CurrentTime()
	g := sharedguilddata.NewGuild(1, "name", "flag", fn.datas, ctime)
	self := sharedguilddata.NewMember(1, fn.datas.GuildClassLevelData().MinKeyData, ctime)

	toUpdate := "new text"

	// 没有权限
	successMsg, errMsg, broadcastChanged := fn.updateText(g, self, toUpdate)
	Ω(successMsg).Should(BeNil())
	Ω(errMsg).Should(Equal(guild.ErrUpdateTextFailDeny))
	Ω(broadcastChanged).Should(Equal(false))

	// npc帮派
	g.SetLeader(-1)
	successMsg, errMsg, broadcastChanged = fn.updateText(g, self, toUpdate)
	Ω(successMsg).Should(BeNil())
	Ω(errMsg).Should(Equal(guild.ErrUpdateTextFailNpc))
	Ω(broadcastChanged).Should(Equal(false))

	Ω(g.GetText()).Should(BeEmpty())
	// 帮主权限
	g.SetLeader(self.Id())
	self.SetClassLevelData(fn.datas.GuildClassLevelData().MaxKeyData)

	successMsg, errMsg, broadcastChanged = fn.updateText(g, self, toUpdate)
	Ω(successMsg).ShouldNot(BeNil())
	Ω(successMsg.Buffer()).Should(Equal(guild.NewS2cUpdateTextMsg(toUpdate).Buffer()))
	Ω(errMsg).Should(BeNil())
	Ω(broadcastChanged).Should(Equal(false))

	Ω(g.GetText()).Should(Equal(toUpdate))
}

func TestGuildModule_ProcessUpdateInternalText(t *testing.T) {
	RegisterTestingT(t)

	fn := newFunc()

	ctime := fn.time.CurrentTime()
	g := sharedguilddata.NewGuild(1, "name", "flag", fn.datas, ctime)
	self := sharedguilddata.NewMember(1, fn.datas.GuildClassLevelData().MinKeyData, ctime)

	toUpdate := "new text"

	// 没有权限
	successMsg, errMsg, broadcastChanged := fn.updateInternalText(g, self, toUpdate)
	Ω(successMsg).Should(BeNil())
	Ω(errMsg).Should(Equal(guild.ErrUpdateInternalTextFailDeny))
	Ω(broadcastChanged).Should(Equal(false))

	// npc帮派
	g.SetLeader(-1)
	successMsg, errMsg, broadcastChanged = fn.updateInternalText(g, self, toUpdate)
	Ω(successMsg).Should(BeNil())
	Ω(errMsg).Should(Equal(guild.ErrUpdateInternalTextFailNpc))
	Ω(broadcastChanged).Should(Equal(false))

	Ω(g.GetInternalText()).Should(BeEmpty())
	// 帮主权限
	g.SetLeader(self.Id())
	self.SetClassLevelData(fn.datas.GuildClassLevelData().MaxKeyData)

	successMsg, errMsg, broadcastChanged = fn.updateInternalText(g, self, toUpdate)
	Ω(successMsg).ShouldNot(BeNil())
	Ω(successMsg.Buffer()).Should(Equal(guild.NewS2cUpdateInternalTextMsg(toUpdate).Buffer()))
	Ω(errMsg).Should(BeNil())
	Ω(broadcastChanged).Should(Equal(false))

	Ω(g.GetInternalText()).Should(Equal(toUpdate))
}

func TestGuildModule_ProcessUpdateLabels(t *testing.T) {
	RegisterTestingT(t)

	fn := newFunc()

	ctime := fn.time.CurrentTime()
	g := sharedguilddata.NewGuild(1, "name", "flag", fn.datas, ctime)
	self := sharedguilddata.NewMember(1, fn.datas.GuildClassLevelData().MinKeyData, ctime)

	labels := []string{"a", "b", "c", "d"}

	// 没有权限
	successMsg, errMsg, broadcastChanged := fn.updateLabel(g, self, labels)
	Ω(successMsg).Should(BeNil())
	Ω(errMsg).Should(Equal(guild.ErrUpdateGuildLabelFailDeny))
	Ω(broadcastChanged).Should(Equal(false))

	// npc帮派
	g.SetLeader(-1)
	successMsg, errMsg, broadcastChanged = fn.updateLabel(g, self, labels)
	Ω(successMsg).Should(BeNil())
	Ω(errMsg).Should(Equal(guild.ErrUpdateGuildLabelFailNpc))
	Ω(broadcastChanged).Should(Equal(false))

	Ω(g.GetLabels()).Should(BeEmpty())
	// 帮主权限
	g.SetLeader(self.Id())
	self.SetClassLevelData(fn.datas.GuildClassLevelData().MaxKeyData)

	successMsg, errMsg, broadcastChanged = fn.updateLabel(g, self, labels)
	Ω(successMsg).ShouldNot(BeNil())
	Ω(successMsg.Buffer()).Should(Equal(guild.NewS2cUpdateGuildLabelMsg(labels).Buffer()))
	Ω(errMsg).Should(BeNil())
	Ω(broadcastChanged).Should(Equal(false))

	Ω(g.GetLabels()).Should(Equal(labels))
}

func TestGuildModule_ProcessUpdateClassNames(t *testing.T) {
	RegisterTestingT(t)

	fn := newFunc()

	ctime := fn.time.CurrentTime()
	g := sharedguilddata.NewGuild(1, "name", "flag", fn.datas, ctime)
	self := sharedguilddata.NewMember(1, fn.datas.GuildClassLevelData().MinKeyData, ctime)

	names := []string{"a", "b", "c", "d"}

	// 没有权限
	successMsg, errMsg, broadcastChanged := fn.updateClassNames(g, self, names)
	Ω(successMsg).Should(BeNil())
	Ω(errMsg).Should(Equal(guild.ErrUpdateClassNamesFailDeny))
	Ω(broadcastChanged).Should(Equal(false))

	// npc帮派
	g.SetLeader(-1)
	successMsg, errMsg, broadcastChanged = fn.updateClassNames(g, self, names)
	Ω(successMsg).Should(BeNil())
	Ω(errMsg).Should(Equal(guild.ErrUpdateClassNamesFailNpc))
	Ω(broadcastChanged).Should(Equal(false))
	g.SetLeader(1)

	Ω(g.GetClassNames()).Should(BeEmpty())
	// 帮主权限
	g.SetLeader(self.Id())
	self.SetClassLevelData(fn.datas.GuildClassLevelData().MaxKeyData)

	// 名字使用中
	duplicateName := append(names, fn.datas.GuildClassTitleData().MinKeyData.Name)
	successMsg, errMsg, broadcastChanged = fn.updateClassNames(g, self, duplicateName)
	Ω(successMsg).Should(BeNil())
	Ω(errMsg).Should(Equal(guild.ErrUpdateClassNamesFailInvalidDuplicate))
	Ω(broadcastChanged).Should(Equal(false))

	g.SetClassTitle(&shared_proto.GuildClassTitleProto{
		CustomClassTitleName: []string{"e"},
	})
	duplicateName = append(names, "e")
	successMsg, errMsg, broadcastChanged = fn.updateClassNames(g, self, duplicateName)
	Ω(successMsg).Should(BeNil())
	Ω(errMsg).Should(Equal(guild.ErrUpdateClassNamesFailInvalidDuplicate))
	Ω(broadcastChanged).Should(Equal(false))

	// 成功
	successMsg, errMsg, broadcastChanged = fn.updateClassNames(g, self, names)
	Ω(successMsg).ShouldNot(BeNil())
	Ω(successMsg.Buffer()).Should(Equal(guild.NewS2cUpdateClassNamesMsg(names).Buffer()))
	Ω(errMsg).Should(BeNil())
	Ω(broadcastChanged).Should(Equal(false))

	Ω(g.GetClassNames()).Should(Equal(names))
}

func TestGuildModule_ProcessUpdateClassTitle(t *testing.T) {
	RegisterTestingT(t)

	fn := newFunc()

	ctime := fn.time.CurrentTime()
	g := sharedguilddata.NewGuild(1, "name", "flag", fn.datas, ctime)
	self := sharedguilddata.NewMember(1, fn.datas.GuildClassLevelData().MinKeyData, ctime)
	g.AddMember(self)

	proto := &shared_proto.GuildClassTitleProto{}
	var heroIds []int64
	toSetTitleData := fn.datas.GetGuildClassTitleDataArray()
	for i := range toSetTitleData {
		target := sharedguilddata.NewMember(int64(i+2), fn.datas.GuildClassLevelData().MinKeyData, ctime)
		g.AddMember(target)

		heroIds = append(heroIds, target.Id())
	}

	// 没有权限
	successMsg, errMsg, broadcastChanged := fn.updateClassTitle(g, self, proto,
		heroIds, toSetTitleData)
	Ω(successMsg).Should(BeNil())
	Ω(errMsg).Should(Equal(guild.ErrUpdateClassTitleFailDeny))
	Ω(broadcastChanged).Should(Equal(false))

	// npc帮派
	g.SetLeader(-1)
	successMsg, errMsg, broadcastChanged = fn.updateClassTitle(g, self, proto,
		heroIds, toSetTitleData)
	Ω(successMsg).Should(BeNil())
	Ω(errMsg).Should(Equal(guild.ErrUpdateClassTitleFailNpc))
	Ω(broadcastChanged).Should(Equal(false))

	Ω(g.GetCustomClassTitle()).Should(BeEmpty())
	// 帮主权限
	g.SetLeader(self.Id())
	self.SetClassLevelData(fn.datas.GuildClassLevelData().MaxKeyData)

	// 名字重复 系统职位
	for _, d := range fn.datas.GetGuildClassLevelDataArray() {
		proto.CustomClassTitleName = append(proto.CustomClassTitleName, d.Name)
		break
	}
	successMsg, errMsg, broadcastChanged = fn.updateClassTitle(g, self, proto,
		heroIds, toSetTitleData)
	Ω(successMsg).Should(BeNil())
	Ω(errMsg).Should(Equal(guild.ErrUpdateClassTitleFailNameExist))
	Ω(broadcastChanged).Should(Equal(false))

	// 系统职称
	proto.CustomClassTitleName = nil
	for _, d := range fn.datas.GetGuildClassTitleDataArray() {
		proto.CustomClassTitleName = append(proto.CustomClassTitleName, d.Name)
		break
	}
	successMsg, errMsg, broadcastChanged = fn.updateClassTitle(g, self, proto,
		heroIds, toSetTitleData)
	Ω(successMsg).Should(BeNil())
	Ω(errMsg).Should(Equal(guild.ErrUpdateClassTitleFailNameExist))
	Ω(broadcastChanged).Should(Equal(false))

	// 自定义职位
	g.SetClassNames([]string{"1团", "2团"})
	proto.CustomClassTitleName = []string{"1团", "2团"}
	successMsg, errMsg, broadcastChanged = fn.updateClassTitle(g, self, proto,
		heroIds, toSetTitleData)
	Ω(successMsg).Should(BeNil())
	Ω(errMsg).Should(Equal(guild.ErrUpdateClassTitleFailNameExist))
	Ω(broadcastChanged).Should(Equal(false))

	g.SetClassNames(nil)

	// 英雄不存在
	successMsg, errMsg, broadcastChanged = fn.updateClassTitle(g, self, proto,
		[]int64{int64(len(toSetTitleData) + 2)}, toSetTitleData)
	Ω(successMsg).Should(BeNil())
	Ω(errMsg).Should(Equal(guild.ErrUpdateClassTitleFailInvalidMemberId))
	Ω(broadcastChanged).Should(Equal(false))

	// 成功
	successMsg, errMsg, broadcastChanged = fn.updateClassTitle(g, self, proto,
		heroIds, toSetTitleData)
	Ω(successMsg).ShouldNot(BeNil())
	Ω(successMsg).Should(Equal(guild.UPDATE_CLASS_TITLE_S2C))
	Ω(errMsg).Should(BeNil())
	Ω(broadcastChanged).Should(Equal(true))

	Ω(g.GetCustomClassTitle()).Should(Equal([]string{"1团", "2团"}))

	for i, v := range toSetTitleData {
		member := g.GetMember(int64(i + 2))

		Ω(member.ClassTitleData()).Should(Equal(v))
	}

	// 再成功一次，原来有的人就没有了
	successMsg, errMsg, broadcastChanged = fn.updateClassTitle(g, self, proto,
		[]int64{1}, toSetTitleData)
	Ω(successMsg).ShouldNot(BeNil())
	Ω(successMsg).Should(Equal(guild.UPDATE_CLASS_TITLE_S2C))
	Ω(errMsg).Should(BeNil())
	Ω(broadcastChanged).Should(Equal(true))

	Ω(g.GetCustomClassTitle()).Should(Equal([]string{"1团", "2团"}))

	for i := range toSetTitleData {
		member := g.GetMember(int64(i + 2))
		Ω(member.ClassTitleData()).Should(BeNil())
	}
}

func TestGuildModule_ProcessUpdateFlagType(t *testing.T) {
	RegisterTestingT(t)

	fn := newFunc()

	ctime := fn.time.CurrentTime()
	g := sharedguilddata.NewGuild(1, "name", "flag", fn.datas, ctime)
	self := sharedguilddata.NewMember(1, fn.datas.GuildClassLevelData().MinKeyData, ctime)

	// 没有权限
	successMsg, errMsg, broadcastChanged := fn.updateFlagType(g, self, 1)
	Ω(successMsg).Should(BeNil())
	Ω(errMsg).Should(Equal(guild.ErrUpdateFlagTypeFailDeny))
	Ω(broadcastChanged).Should(Equal(false))

	// npc帮派
	g.SetLeader(-1)
	successMsg, errMsg, broadcastChanged = fn.updateFlagType(g, self, 1)
	Ω(successMsg).Should(BeNil())
	Ω(errMsg).Should(Equal(guild.ErrUpdateFlagTypeFailNpc))
	Ω(broadcastChanged).Should(Equal(false))

	Ω(g.FlagType()).Should(Equal(uint64(0)))
	// 帮主权限
	g.SetLeader(self.Id())
	self.SetClassLevelData(fn.datas.GuildClassLevelData().MaxKeyData)

	successMsg, errMsg, broadcastChanged = fn.updateFlagType(g, self, 1)
	Ω(successMsg).ShouldNot(BeNil())
	Ω(successMsg.Buffer()).Should(Equal(guild.NewS2cUpdateFlagTypeMsg(1).Buffer()))
	Ω(errMsg).Should(BeNil())
	Ω(broadcastChanged).Should(Equal(false))

	Ω(g.FlagType()).Should(Equal(uint64(1)))
}

// 更新帮派成员职位
func TestGuildModule_ProcessUpdateMemberClassLevel(t *testing.T) {
	RegisterTestingT(t)

	fn := newFunc()

	ctime := fn.time.CurrentTime()

	g := sharedguilddata.NewGuild(1, "name", "flag", fn.datas, ctime)
	self := sharedguilddata.NewMember(1, fn.datas.GuildClassLevelData().MinKeyData, ctime)
	g.AddMember(self)

	// 没有权限
	successMsg, errMsg, broadcastChanged := fn.updateMemberClassLevel(g, self, 2, fn.datas.GuildClassLevelData().MinKeyData)
	Ω(successMsg).Should(BeNil())
	Ω(errMsg).Should(Equal(guild.ErrUpdateMemberClassLevelFailDeny))
	Ω(broadcastChanged).Should(Equal(false))

	// npc帮派
	g.SetLeader(-1)
	successMsg, errMsg, broadcastChanged = fn.updateMemberClassLevel(g, self, 2, fn.datas.GuildClassLevelData().MinKeyData)
	Ω(successMsg).Should(BeNil())
	Ω(errMsg).Should(Equal(guild.ErrUpdateMemberClassLevelFailNpc))
	Ω(broadcastChanged).Should(Equal(false))

	// 帮主权限
	g.SetLeader(self.Id())
	self.SetClassLevelData(fn.datas.GuildClassLevelData().MaxKeyData)

	// target not exist
	successMsg, errMsg, broadcastChanged = fn.updateMemberClassLevel(g, self, 2, fn.datas.GuildClassLevelData().MinKeyData)
	Ω(successMsg).Should(BeNil())
	Ω(errMsg).Should(Equal(guild.ErrUpdateMemberClassLevelFailTargetNotInGuild))
	Ω(broadcastChanged).Should(Equal(false))

	other := sharedguilddata.NewMember(2, fn.datas.GuildClassLevelData().MinKeyData, ctime)
	g.AddMember(other)

	for i := 1; i < len(fn.datas.GetGuildClassLevelDataArray())-1; i++ {
		toSet := fn.datas.GetGuildClassLevelDataArray()[i]

		successMsg, errMsg, broadcastChanged = fn.updateMemberClassLevel(g, self, 2, toSet)
		Ω(successMsg).Should(BeNil())
		Ω(errMsg).Should(BeNil())
		Ω(broadcastChanged).Should(Equal(true))

		Ω(other.ClassLevelData()).Should(Equal(toSet))
	}

	// 转移帮主
	successMsg, errMsg, broadcastChanged = fn.updateMemberClassLevel(g, self, 2, fn.datas.GuildClassLevelData().MaxKeyData)
	Ω(successMsg).Should(BeNil())
	Ω(errMsg).Should(BeNil())
	Ω(broadcastChanged).Should(Equal(true))

	Ω(self.ClassLevelData()).Should(Equal(fn.datas.GuildClassLevelData().MinKeyData))
	testCheckLeader(fn, g, other)

	// 玩家大于X个时，需要倒计时
	for i := 2; i < u64.Int(fn.datas.GuildConfig().ChangeLeaderCountdownMemberCount); i++ {
		other := sharedguilddata.NewMember(int64(i+1), fn.datas.GuildClassLevelData().MinKeyData, ctime)
		g.AddMember(other)
	}

	// 倒计时
	successMsg, errMsg, broadcastChanged = fn.updateMemberClassLevel(g, other, 1, fn.datas.GuildClassLevelData().MaxKeyData)
	Ω(successMsg).Should(BeNil())
	Ω(errMsg).Should(BeNil())
	Ω(broadcastChanged).Should(Equal(true))

	Ω(self.ClassLevelData()).Should(Equal(fn.datas.GuildClassLevelData().MinKeyData))
	testCheckLeader(fn, g, other)

	Ω(g.GetChangeLeaderId()).Should(Equal(self.Id()))
	Ω(g.HasChangeLeaderCountDown()).Should(BeTrue())

	countdownTime := ctime.Add(fn.datas.GuildConfig().ChangeLeaderCountdownDuration)
	result := g.TryTickChangeLeader(countdownTime.Add(-1), fn.datas.GuildClassLevelData().MinKeyData, fn.datas.GuildClassLevelData().MaxKeyData)
	Ω(result).Should(BeFalse())

	result = g.TryTickChangeLeader(countdownTime, fn.datas.GuildClassLevelData().MinKeyData, fn.datas.GuildClassLevelData().MaxKeyData)
	Ω(result).Should(BeTrue())

	testCheckLeader(fn, g, self)
	Ω(other.ClassLevelData()).Should(Equal(fn.datas.GuildClassLevelData().MinKeyData))

	// 取消
	successMsg, errMsg, broadcastChanged = fn.updateMemberClassLevel(g, self, 2, fn.datas.GuildClassLevelData().MaxKeyData)
	Ω(successMsg).Should(BeNil())
	Ω(errMsg).Should(BeNil())
	Ω(broadcastChanged).Should(Equal(true))

	Ω(g.GetChangeLeaderId()).Should(Equal(other.Id()))
	Ω(g.HasChangeLeaderCountDown()).Should(BeTrue())

	g.CancelChangeLeader()
	Ω(g.GetChangeLeaderId()).Should(Equal(int64(0)))
	Ω(g.HasChangeLeaderCountDown()).Should(BeFalse())
}

func TestDonate(t *testing.T) {
	RegisterTestingT(t)

	fn := newFunc()

	ctime := fn.time.CurrentTime()

	g := sharedguilddata.NewGuild(1, "name", "flag", fn.datas, ctime)
	self := sharedguilddata.NewMember(1, fn.datas.GuildClassLevelData().MinKeyData, ctime)
	g.AddMember(self)

	hero := entity.NewHero(1, "hero1", fn.datas.HeroInitData(), ctime)
	mock.DefaultHero(hero)

	var maxLevelDonateBuilding *domestic_data.BuildingData
	for _, v := range fn.datas.GetBuildingDataArray() {
		if v.Type == domestic_data.DonateBuilding {
			if maxLevelDonateBuilding == nil || maxLevelDonateBuilding.Level < v.Level {
				maxLevelDonateBuilding = v
			}
		}
	}
	hero.Domestic().SetBuilding(maxLevelDonateBuilding)

	for hero.Level() < fn.datas.GuildConfig().GuildDonateNeedHeroLevel {
		hero.AddExp(hero.LevelData().Sub.UpgradeExp * 10) // 升级
	}

	var contributionAmount, donationAmount, donationYuanbao, buildingAmount, coin uint64
	for i, v := range fn.datas.GetGuildDonateDataArray() {
		if uint64(i) >= maxLevelDonateBuilding.Effect.GuildDonateTimes {
			break
		}

		mock.SetHeroEnoughCost(hero, v.Cost)

		successMsg, errMsg, broadcastChanged := fn.donation(ifacemock.HeroController, g, self, v.Sequence)
		Ω(successMsg).Should(BeEquivalentTo(guild.NewS2cDonateMsg(u64.Int32(v.Sequence), u64.Int32(v.Times), u64.Int32(v.Id), u64.Int32(g.GetBuildingAmount()),
			u64.Int32(self.ContributionAmount()), u64.Int32(self.ContributionTotalAmount()),
			u64.Int32(self.ContributionAmount7()), u64.Int32(self.DonationAmount()),
			u64.Int32(self.DonationTotalAmount()), u64.Int32(self.DonationAmount7()),
			u64.Int32(self.DonateTotalYuanbao()))))
		Ω(errMsg).Should(BeNil())
		Ω(broadcastChanged).Should(Equal(true))

		// 贡献点
		contributionAmount += v.ContributionAmount
		Ω(self.ContributionAmount()).Should(Equal(contributionAmount))

		// 捐献值
		donationAmount += v.DonationAmount
		Ω(self.DonationAmount()).Should(Equal(donationAmount))

		// 捐献元宝
		donationYuanbao += v.Cost.Yuanbao
		Ω(self.DonateTotalYuanbao()).Should(Equal(donationYuanbao))

		// 超过一天后，都加建设值
		buildingAmount += v.GuildBuildingAmount

		// 建设值
		Ω(g.GetBuildingAmount()).Should(Equal(buildingAmount))

		coin += v.ContributionCoin
		Ω(hero.GetGuildContributionCoin()).Should(Equal(coin))

		t, _ := hero.GetDonateTimes(v.Sequence)
		Ω(t).Should(Equal(v.Times))
	}

}

func TestGuildModule_ProcessSendYinliangToOtherGuild(t *testing.T) {
	RegisterTestingT(t)

	fn := newFunc()

	guilds := sharedguilddata.NewGuilds(fn.datas)

	ctime := fn.time.CurrentTime()
	g := sharedguilddata.NewGuild(1, "name", "flag", fn.datas, ctime)
	guilds.Add(g)
	self := sharedguilddata.NewMember(1, fn.datas.GuildClassLevelData().MinKeyData, ctime)
	g.SetLeader(1)
	g.AddMember(self)
	self.SetClassLevelData(fn.datas.GuildClassLevelData().MaxKeyData)
	g.SetYinliang(10000)

	g2 := sharedguilddata.NewGuild(2, "name", "flag", fn.datas, ctime)
	self2 := sharedguilddata.NewMember(2, fn.datas.GuildClassLevelData().MinKeyData, ctime)
	guilds.Add(g2)
	g2.SetLeader(2)
	g2.AddMember(self2)

	_, err, _ := fn.sendYinliangToGuild(guilds, g, self, 2, 1000)
	Ω(err).Should(BeNil())
	Ω(g.GetYinliang()).Should(Equal(uint64(10000 - 1000)))
	Ω(g2.GetYinliang()).Should(Equal(uint64(1000)))
}

func TestGuildModule_ProcessUpgradeLevel(t *testing.T) {
	RegisterTestingT(t)

	fn := newFunc()

	ctime := fn.time.CurrentTime()

	g := sharedguilddata.NewGuild(1, "name", "flag", fn.datas, ctime)
	self := sharedguilddata.NewMember(1, fn.datas.GuildClassLevelData().MinKeyData, ctime)
	g.AddMember(self)

	successMsg, errMsg, broadcastChanged := fn.upgradeLevel(g, self)
	Ω(successMsg).Should(BeNil())
	Ω(errMsg).Should(Equal(guild.ErrUpgradeLevelFailDeny))
	Ω(broadcastChanged).Should(Equal(false))

	// 帮主权限
	g.SetLeader(self.Id())
	self.SetClassLevelData(fn.datas.GuildClassLevelData().MaxKeyData)

	// 建设值不足
	successMsg, errMsg, broadcastChanged = fn.upgradeLevel(g, self)
	Ω(successMsg).Should(BeNil())
	Ω(errMsg).Should(Equal(guild.ErrUpgradeLevelFailCostNotEnough))
	Ω(broadcastChanged).Should(Equal(false))

	// 加建设值
	g.AddBuildingAmount(g.LevelData().UpgradeBuilding, ctime)

	originLevel := g.LevelData()

	// 成功
	successMsg, errMsg, broadcastChanged = fn.upgradeLevel(g, self)
	Ω(successMsg).Should(BeEquivalentTo(guild.UPGRADE_LEVEL_S2C))
	Ω(errMsg).Should(BeNil())
	Ω(broadcastChanged).Should(Equal(true))

	Ω(g.GetBuildingAmount()).Should(Equal(uint64(0)))

	if g.LevelData().UpgradeDuration > 0 {
		// 有倒计时
		Ω(g.GetUpgradeEndTime()).Should(Equal(ctime.Add(g.LevelData().UpgradeDuration)))
		Ω(g.LevelData()).Should(BeEquivalentTo(originLevel))

		fn.tryUpgradeGuildLevel(g, g.GetUpgradeEndTime().Add(-1))
		Ω(g.GetUpgradeEndTime()).Should(Equal(ctime.Add(g.LevelData().UpgradeDuration)))
		Ω(g.LevelData()).Should(BeEquivalentTo(originLevel))

		fn.tryUpgradeGuildLevel(g, g.GetUpgradeEndTime())
	}

	Ω(g.GetUpgradeEndTime()).Should(Equal(time.Time{}))
	Ω(g.LevelData()).Should(BeEquivalentTo(originLevel.NextLevel()))
}

func testCheckLeader(fn *guild_func, g *sharedguilddata.Guild, leader *sharedguilddata.GuildMember) {
	Ω(g.LeaderId()).Should(Equal(leader.Id()))
	Ω(leader.ClassLevelData()).Should(BeEquivalentTo(fn.datas.GuildClassLevelData().MaxKeyData))
}

func TestGuildTarget(t *testing.T) {
	RegisterTestingT(t)

	fn := newFunc()

	ctime := fn.time.CurrentTime()

	g := sharedguilddata.NewGuild(1, "name", "flag", fn.datas, ctime)

	for _, t := range fn.datas.GetGuildTargetArray() {
		isDoing, st, et := g.IsDoingTarget(t, fn.datas.GuildConfig(), ctime)

		switch t.TargetType {
		case shared_proto.GuildTargetType_GuildLevelUp:
			Ω(isDoing).Should(Equal(g.LevelData().Level < t.Target))
			Ω(st).Should(Equal(time.Time{}))
			Ω(et).Should(Equal(time.Time{}))

			if isDoing {
				g.SetLevelData(fn.datas.GetGuildLevelData(t.Target))
				isDoing, st, et = g.IsDoingTarget(t, fn.datas.GuildConfig(), ctime)
				Ω(isDoing).Should(BeFalse())
				Ω(st).Should(Equal(time.Time{}))
				Ω(et).Should(Equal(time.Time{}))
			}

		case shared_proto.GuildTargetType_PrestigeUp:
			Ω(isDoing).Should(Equal(g.GetPrestige() < t.Target))
			Ω(st).Should(Equal(time.Time{}))
			Ω(et).Should(Equal(time.Time{}))

			if isDoing {
				g.RefreshHistoryMaxPrestige(t.Target)
				isDoing, st, et = g.IsDoingTarget(t, fn.datas.GuildConfig(), ctime)
				Ω(isDoing).Should(BeFalse())
				Ω(st).Should(Equal(time.Time{}))
				Ω(et).Should(Equal(time.Time{}))
			}

		case shared_proto.GuildTargetType_ImpeachNpcLeader:
			// TODO
		case shared_proto.GuildTargetType_ImpeachUserLeader:
			// TODO
		case shared_proto.GuildTargetType_UpdateMemverClass:
			// TODO
		case shared_proto.GuildTargetType_UserLeaderUseless:
			Ω(isDoing).Should(BeFalse())

			g.UpdateLeaderOfflineTime(ctime.Add(-fn.datas.GuildConfig().ImpeachUserLeaderOffline - 1))

			isDoing, _, _ = g.IsDoingTarget(t, fn.datas.GuildConfig(), ctime)
			Ω(isDoing).Should(BeTrue())

			// 弹劾开始
			g.StartImpeachLeader(0, ctime, ctime.Add(time.Hour), 1)
			isDoing, _, _ = g.IsDoingTarget(t, fn.datas.GuildConfig(), ctime)
			Ω(isDoing).Should(BeFalse())

			// 弹劾结束
			g.SetImpeachLeader(nil)
			isDoing, _, _ = g.IsDoingTarget(t, fn.datas.GuildConfig(), ctime)
			Ω(isDoing).Should(BeTrue())
		case shared_proto.GuildTargetType_GuildChangeCountry:
			Ω(isDoing).Should(BeFalse())
			Ω(st).Should(Equal(time.Time{}))
			Ω(et).Should(Equal(time.Time{}))

			array := fn.datas.CountryData().Array
			last := array[len(array)-1]

			waitTime := ctime.Add(fn.datas.GuildGenConfig().GuildChangeCountryWaitDuration)
			nextTime := ctime.Add(fn.datas.GuildGenConfig().GuildChangeCountryCooldown)
			g.SetChangeCountry(last, waitTime.Unix(), nextTime.Unix())

			isDoing, st, et = g.IsDoingTarget(t, fn.datas.GuildConfig(), ctime)
			Ω(isDoing).Should(BeTrue())
			Ω(st).Should(Equal(ctime))
			Ω(et).Should(Equal(waitTime))

			// 弹劾结束
			g.CancelChangeCountry()
			isDoing, _, _ = g.IsDoingTarget(t, fn.datas.GuildConfig(), ctime)
			Ω(isDoing).Should(BeTrue())
		}
	}

}

func newFunc() *guild_func {
	datas, err := config.LoadConfigDatas(confpath.GetConfigPath())
	Ω(err).Should(Succeed())

	selfGuildMsgCache := concurrent.NewI64VersionBufferCacheMap(func(int64, uint64) (pbutil.Buffer, error) {
		return pbutil.Empty, nil
	})

	return newGuildFunc(mock.MockDep(), datas, ifacemock.HebiModule, ifacemock.DbService, ifacemock.MailModule,
		ifacemock.RealmService, ifacemock.XiongNuModule, ifacemock.XiongNuService, ifacemock.CountryService,
		ifacemock.BaiZhanService, ifacemock.RankModule, ifacemock.PushService, ifacemock.MingcWarService, 0, selfGuildMsgCache.Clear,
		ifacemock.GuildService.UpdateSnapshot, ifacemock.GuildService.GetSnapshot, ifacemock.GuildService.RemoveSnapshot,
		ifacemock.ChatService, ifacemock.GuildService.AddLog,
		newGuildTemplateArray(datas.GetNpcGuildTemplateArray(), nil))
}
