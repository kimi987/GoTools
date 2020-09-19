package guild_data

import (
	"fmt"
	"github.com/lightpaw/male7/config/monsterdata"
	"github.com/lightpaw/male7/pb/shared_proto"
	"github.com/lightpaw/male7/util/check"
	"sort"
	"strconv"
	"strings"
	"unicode/utf8"
	"github.com/lightpaw/male7/config/country"
)

//gogen:config
type NpcGuildSuffixName struct {
	_ struct{} `singleton:"true"`
	_ struct{} `file:"联盟/联盟Npc后缀名.txt"`

	Name []string
}

func (c *NpcGuildSuffixName) Init(filename string) {

	check.PanicNotTrue(len(c.Name) > 0, "%s 配置的Npc联盟后缀配置为空", filename)

	for _, s := range c.Name {
		check.PanicNotTrue(utf8.RuneCountInString(s) == 1, "%s 配置的Npc联盟后缀必须为1个字符", filename)
	}
}

func (c *NpcGuildSuffixName) Must(idx uint64) string {
	if len(c.Name) > 0 {
		idx = idx % uint64(len(c.Name))
		return c.Name[idx]
	}

	return ""
}

//gogen:config
type NpcMemberData struct {
	_ struct{} `file:"联盟/联盟Npc成员.txt"`

	Id uint64

	// 君主信息
	Master *monsterdata.MonsterMasterData

	// 每日贡献度
	ContributionAmount uint64

	// 7日贡献度
	ContributionAmount7 uint64

	// 总贡献度
	TotalContributionAmount uint64
}

func (m *NpcMemberData) EncodeSnapshot(id int64) *shared_proto.HeroBasicSnapshotProto {
	return m.Master.EncodeSnapshot(id)
}

//gogen:config
type NpcGuildTemplate struct {
	_  struct{} `file:"联盟/联盟Npc模板.txt"`
	Id uint64

	Name string // 联盟名

	FlagName string // 旗号

	Text         string
	InternalText string
	Labels       []string

	Level *GuildLevelData

	Country *country.CountryData

	RejectUserJoin bool // true表示A类联盟(拒绝玩家加入)，否则表示B类联盟（允许玩家加入）

	NpcLeaderVote uint64 // B类联盟NPC 盟主基础票数

	// 盟主
	Leader *NpcMemberData

	// 初始成员，按照7日贡献度排序
	Members []*NpcMemberData

	contribution7Members []*NpcMemberData

	// 缓存组合出来的数据
	guildNames []string
	flagNames  []string
}

func (t *NpcGuildTemplate) Init(filename string, configs interface {
	NpcGuildSuffixName() *NpcGuildSuffixName
}) {

	t.Name = strings.TrimSpace(t.Name)
	t.FlagName = strings.TrimSpace(t.FlagName)

	//check.PanicNotTrue(len(t.Members) > 0, "%s 联盟模板%v-%s 配置的联盟成员个数必须>0", filename, t.Id, t.Name)

	maxMemberCount := 512
	check.PanicNotTrue(len(t.Members)+1 <= maxMemberCount, "%s 联盟模板%v-%s 配置的联盟成员(包含盟主)个数必须<=%v", filename, t.Id, t.Name, maxMemberCount)

	for _, m := range t.Members {
		check.PanicNotTrue(m.Id != t.Leader.Id, "%s 联盟模板%v-%s 配置的联盟成员列表中包含盟主的id", filename, t.Id, t.Name)
	}

	t.contribution7Members = make([]*NpcMemberData, len(t.Members))
	copy(t.contribution7Members, t.Members)
	sort.Sort(c7slice(t.contribution7Members))

	check.PanicNotTrue(len(t.Name) > 0, "%s 联盟模板%v-%s 配置的名字为空", filename, t.Id, t.Name)
	check.PanicNotTrue(len(t.FlagName) > 0, "%s 联盟模板%v-%s 配置的旗号为空", filename, t.Id, t.Name)

	if t.RejectUserJoin {

		check.PanicNotTrue(!strings.Contains(t.Name, "%s"), "%s 联盟模板%v-%s (A类联盟)配置的名字不能包含 %s", filename, t.Id, t.Name, "%s")
		check.PanicNotTrue(!strings.Contains(t.FlagName, "%s"), "%s 联盟模板%v-%s (A类联盟)配置的名字不能包含 %s", filename, t.Id, t.Name, "%s")

	} else {
		check.PanicNotTrue(strings.Contains(t.Name, "%s"), "%s 联盟模板%v-%s 配置的名字必须包含 %s", filename, t.Id, t.Name, "%s")
		check.PanicNotTrue(strings.Contains(t.FlagName, "%s"), "%s 联盟模板%v-%s 配置的名字必须包含 %s", filename, t.Id, t.Name, "%s")

		// 帮派名字 XXX1盟 XXX2盟
		// 旗号 甲A 甲B 甲C
		for i, name := range configs.NpcGuildSuffixName().Name {
			t.guildNames = append(t.guildNames, fmt.Sprintf(t.Name, strconv.Itoa(i+1)))
			t.flagNames = append(t.flagNames, fmt.Sprintf(t.FlagName, name))
		}
	}

}

func (t *NpcGuildTemplate) GetCombineNames() []string {
	return t.guildNames
}

func (t *NpcGuildTemplate) GetCombineFlagNames() []string {
	return t.flagNames
}

func (t *NpcGuildTemplate) Contribution7Members() []*NpcMemberData {
	return t.contribution7Members
}

func (t *NpcGuildTemplate) GetNpc(seq uint64) *NpcMemberData {
	if seq > 0 && seq <= uint64(len(t.contribution7Members)) {
		return t.contribution7Members[seq-1]
	}

	if seq == 0 {
		return t.Leader
	}

	if len(t.contribution7Members) <= 0 {
		return nil
	}

	return t.contribution7Members[len(t.contribution7Members)-1]
}

// guild id

//const seqbit = 9 // 512
//const seqmask = 1<<9 - 1
//const gidmask = math.MaxInt64 >> seqbit
//
//func InvalidNpcGuildId(guildId int64) bool {
//	return guildId <= 0 || guildId > gidmask
//}
//
//func InvalidNpcMemberSequence(sequence int) bool {
//	return sequence < 0 || sequence > seqmask
//}
//
//func EncodeNpcMemberId(guildId int64, sequence int) int64 {
//
//	if InvalidNpcGuildId(guildId) {
//		logrus.Errorf("EncodeNpcMemberId invalid npc guild id")
//		return -1
//	}
//
//	if InvalidNpcMemberSequence(sequence) {
//		logrus.Errorf("EncodeNpcMemberId invalid npc member sequene")
//		return -1
//	}
//
//	return -int64(uint64(guildId)<<seqbit | uint64(sequence))
//}
//
//func NpcMemberSequence(id int64) uint64 {
//	return uint64(-id) & seqmask
//}
//
//func NpcMemberGuildId(id int64) int64 {
//	return int64(uint64(-id) >> seqbit)
//}

type c7slice []*NpcMemberData

func (a c7slice) Len() int      { return len(a) }
func (a c7slice) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a c7slice) Less(i, j int) bool {
	ai := a[i]
	aj := a[j]
	if ai.ContributionAmount7 == aj.ContributionAmount7 {
		// 相同，id小的在前面
		return ai.Id < aj.Id
	}

	return ai.ContributionAmount7 > aj.ContributionAmount7
}
