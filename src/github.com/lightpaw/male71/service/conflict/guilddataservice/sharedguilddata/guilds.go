package sharedguilddata

import (
	"github.com/lightpaw/male7/util/i64"
	"sort"
	"github.com/lightpaw/male7/util/u64"
	"github.com/lightpaw/pbutil"
	"github.com/lightpaw/male7/pb/shared_proto"
	"github.com/lightpaw/male7/gen/pb/guild"
	"sync"
	"time"
	"github.com/lightpaw/male7/config"
	"github.com/lightpaw/logrus"
)

type Guilds interface {
	Get(int64) *Guild
	Read(int64) *Guild
	IdRankArray() []*Guild
	Add(*Guild)
	Remove(int64) *Guild
	// 线程安全
	GetIdByName(name string) int64
	ChangeName(oldName, newName string, guildId int64)
	GetIdByFlagName(flagName string) int64
	ChangeFlagName(oldName, newName string, guildId int64)
	Walk(f Func)
	RefreshPrestigeRank() []int64
	PrestigeRankMsg(uint64, time.Time, GenerateRankMsgFunc) pbutil.Buffer
	ListGuild(heroLevel, junXianLevel, towerMaxFloor, countryId uint64) []*Guild
}

type Func func(g *Guild)

type Funcs func(guilds Guilds)

type GenerateRankMsgFunc func(uint64, int) pbutil.Buffer

func NewGuilds(datas *config.ConfigDatas) Guilds {
	m := &guilds{}
	m.dataMap = make(map[int64]*Guild)
	m.nameMap = i64.NewStringI64Map()
	m.flagNameMap = i64.NewStringI64Map()
	m.prestigeRankLimit = len(datas.GetGuildRankPrizeDataArray())
	m.prestigeRankMsgMap = make(map[uint64]*prestigeRankMsg)
	for _, data := range datas.GetCountryDataArray() {
		m.prestigeRankMsgMap[data.Id] = &prestigeRankMsg{}
	}
	return m
}

type prestigeRankMsg struct {
	mux         sync.RWMutex
	msg         pbutil.Buffer
	timeoutTime time.Time
}

func (p *prestigeRankMsg) get() (time.Time, pbutil.Buffer) {
	p.mux.RLock()
	defer p.mux.RUnlock()

	return p.timeoutTime, p.msg
}

func (p *prestigeRankMsg) Set(time time.Time, msg pbutil.Buffer) {
	p.mux.Lock()
	defer p.mux.Unlock()

	p.timeoutTime = time
	p.msg = msg
}

const timeout = 5 * time.Second

func (p *prestigeRankMsg) GetRankMsg(cid uint64, limit int, ctime time.Time, f GenerateRankMsgFunc) pbutil.Buffer {

	timeoutTime, msg := p.get()
	if ctime.Before(timeoutTime) {
		return msg
	}

	msg = f(cid, limit)
	p.Set(ctime.Add(timeout), msg)
	return msg
}

type guilds struct {
	// 数据以此为准
	dataMap map[int64]*Guild

	nameMap     *i64.StringI64Map
	flagNameMap *i64.StringI64Map

	// 按id排好序，主要是查询帮派时候用
	idRankArray []*Guild
	// 当日声望排名的静态消息
	prestigeRankMsgMap map[uint64]*prestigeRankMsg
	prestigeRankLimit  int // 排行数限制
	rankFunc           GenerateRankMsgFunc
}

func (m *guilds) Get(id int64) *Guild {
	g := m.dataMap[id]
	if g != nil {
		g.SetChanged()
	}

	return g
}

func (m *guilds) Read(id int64) *Guild {
	return m.dataMap[id]
}

func (m *guilds) Walk(f Func) {
	for _, g := range m.dataMap {
		f(g)
	}
}

func (m *guilds) IdRankArray() []*Guild {
	return m.idRankArray
}

// --- 添加

func (m *guilds) Add(toAdd *Guild) {
	if m.dataMap[toAdd.id] != nil {
		// 防御性
		return
	}

	m.dataMap[toAdd.id] = toAdd
	m.nameMap.Set(toAdd.Name(), toAdd.Id())
	m.flagNameMap.Set(toAdd.FlagName(), toAdd.Id())

	if toAdd.GetNpcTemplate() == nil || !toAdd.GetNpcTemplate().RejectUserJoin {
		// 非a类联盟才可以加进来
		lastIdGuild := GetLast(m.idRankArray)
		if lastIdGuild == nil || lastIdGuild.id < toAdd.id {
			// 这是正常情况
			m.idRankArray = append(m.idRankArray, toAdd)
		} else {
			// 异常了，排下序吧
			m.idRankArray = append(m.idRankArray, toAdd)
			sort.Sort(idRankSlice(m.idRankArray))
		}
	}
}

// 添加 ---

func (m *guilds) Remove(toRemoveId int64) *Guild {
	toRemove := m.dataMap[toRemoveId]
	if toRemove != nil {
		delete(m.dataMap, toRemoveId)

		m.nameMap.RemoveIfSame(toRemove.Name(), toRemoveId)
		m.flagNameMap.RemoveIfSame(toRemove.FlagName(), toRemoveId)

		m.idRankArray = RemoveAndLeftShift(m.idRankArray, toRemoveId)
	}

	return toRemove
}

func (m *guilds) GetIdByName(name string) int64 {
	key, _ := m.nameMap.Get(name)
	return key
}

func (m *guilds) ChangeName(oldName, newName string, guildId int64) {
	if m.nameMap.RemoveIfSame(oldName, guildId) {
		m.nameMap.Set(newName, guildId)
	}
}

func (m *guilds) GetIdByFlagName(flagName string) int64 {
	key, _ := m.flagNameMap.Get(flagName)
	return key
}

func (m *guilds) ChangeFlagName(oldName, newName string, guildId int64) {
	if m.flagNameMap.RemoveIfSame(oldName, guildId) {
		m.flagNameMap.Set(newName, guildId)
	}
}

// prestigeRankArrMap
func (m *guilds) RefreshPrestigeRank() (refreshedGuilds []int64) {
	logrus.Debug("刷新联盟国家排行榜")
	// 根据国家筛选出符合排序条件的联盟
	prestigeRankArrMap := make(map[uint64][]*Guild)
	for _, guild := range m.dataMap {
		logrus.Debug(guild.name)
		// 上次声望奖励排名清零
		if guild.GetLastPrestigeRank() > 0 {
			guild.SetLastPrestigeRank(0)
			refreshedGuilds = append(refreshedGuilds, guild.id)
		}

		// 过滤掉没有任何声望值的
		if guild.YesterdayPrestige() == 0 {
			logrus.Debug(guild.name, "昨日声望为0，跳过")
			continue
		}
		country := guild.CountryId()
		prestigeRankArrMap[country] = append(prestigeRankArrMap[country], guild)
	}

	// 对每个国家的联盟根据声望降序
	for _, guilds := range prestigeRankArrMap {
		sort.Sort(prestigeRankSlice(guilds))
		// 只对排行前N名（配置）的联盟进行排名标记
		gLen := len(guilds)
		for i := 0; i < gLen && i < m.prestigeRankLimit; i++ {
			guild := guilds[i]
			rank := u64.FromInt(i + 1)
			guild.SetLastPrestigeRank(rank)
			if !i64.Contains(refreshedGuilds, guild.id) {
				refreshedGuilds = append(refreshedGuilds, guild.id)
			}
			logrus.Debug("联盟国家排行榜名次", rank, guild.name)
		}
	}
	return
}

var emptyPrestigeRankMsg = guild.NewS2cViewDailyGuildRankMsg(&shared_proto.RankProto{}).Static()

func (m *guilds) PrestigeRankMsg(cid uint64, ctime time.Time, f GenerateRankMsgFunc) pbutil.Buffer {
	if msg := m.prestigeRankMsgMap[cid]; msg != nil {
		return msg.GetRankMsg(cid, m.prestigeRankLimit, ctime, f)
	}
	return emptyPrestigeRankMsg
}

// 无联盟玩家打开界面，收到联盟列表的特殊函数，只取10个
func (m *guilds) ListGuild(heroLevel, junXianLevel, towerMaxFloor, countryId uint64) (list []*Guild) {
	logrus.Debug("刷新无联盟玩家打开界面联盟列表")
	// 筛选出符合基本条件的联盟
	var guildList1 []*Guild
	var guildList2 []*Guild
	for _, guild := range m.dataMap {
		if guild.MemberCount() >= u64.Int(guild.levelData.MemberCount) {
			continue
		}
		if guild.requiredHeroLevel > heroLevel {
			continue
		}
		if guild.requiredJunXianLevel > junXianLevel {
			continue
		}
		if guild.requiredTowerMaxFloor > towerMaxFloor {
			continue
		}
		if guild.rejectAutoJoin {
			guildList2 = append(guildList2, guild)
		} else {
			guildList1 = append(guildList1, guild)
		}
	}
	size := len(guildList1)
	if size > 0 {
		sort.Sort(prestigeCoreRankSlice(guildList1))
		if size > 5 {
			size = 5
		}
		for i := 0; i < size; i++ {
			list = append(list, guildList1[i])
		}
	}
	size = 10 - size
	if size2 := len(guildList2); size2 > 0 {
		sort.Sort(prestigeCoreRankSlice(guildList2))
		if size2 > size {
			size2 = size
		}
		for i := 0; i < size2; i++ {
			list = append(list, guildList2[i])
		}
		size -= size2
	}
	if size > 0 { // 还有剩余空间没填满
		guildMap := make(map[int64]*Guild, len(list))
		for _, g := range list {
			guildMap[g.id] = g
		}
		var guildList3 []*Guild
		for _, guild := range m.dataMap {
			if guild.CountryId() == countryId {
				guildList3 = append(guildList3, guild)
			}
		}
		sort.Sort(prestigeCoreRankSlice(guildList3))
		for _, g := range guildList3 {
			if guildMap[g.id] != nil {
				continue
			}
			list = append(list, g)
			if size--; size <= 0 {
				break
			}
		}
	}
	return
}
