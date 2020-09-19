package secret_tower

import (
	"github.com/lightpaw/logrus"
	"github.com/lightpaw/male7/config/towerdata"
	"github.com/lightpaw/male7/entity"
	"github.com/lightpaw/male7/gen/iface"
	"github.com/lightpaw/male7/gen/pb/secret_tower"
	"github.com/lightpaw/male7/pb/shared_proto"
	"github.com/lightpaw/male7/service/herosnapshot/snapshotdata"
	"github.com/lightpaw/male7/util/i64"
	"github.com/lightpaw/male7/util/idbytes"
	"github.com/lightpaw/male7/util/must"
	"github.com/lightpaw/male7/util/timeutil"
	"github.com/lightpaw/male7/util/u64"
	"github.com/lightpaw/pbutil"
	"time"
	"github.com/lightpaw/male7/util/collection"
)

func newSecretTowerTeam(manager *team_manager, secretTowerData *towerdata.SecretTowerData, leader *secret_tower_team_member, guildId int64, ctime time.Time) *secret_tower_team {
	team := &secret_tower_team{
		towerData:       secretTowerData,
		teamId:          manager.newTeamId(),
		leader:          leader,
		members:         make([]*secret_tower_team_member, 0, secretTowerData.MaxAttackerCount),
		removeTeam:      manager.removeTeam,
		onHeroLeaveTeam: manager.onHeroLeaveTeam,
		createTime:      ctime,
		expireTime:      ctime.Add(secretTowerData.TeamExpireDuration),
		guildId:         guildId,
		chatRecord:      collection.NewRingList(10),
	}

	team.add(leader, ctime)

	return team
}

type r_secret_tower_team interface {
	// 序列化用于队伍展示，可以多线程访问
	Encode4Show() []byte
}

// 密室队伍
type secret_tower_team struct {
	// 密室数据
	towerData *towerdata.SecretTowerData

	// 队伍id
	teamId int64

	// 队长
	leader *secret_tower_team_member

	// 所有成员
	members []*secret_tower_team_member

	// 移除队伍
	removeTeam func(team *secret_tower_team)

	// 处理成员离开队伍的回调
	onHeroLeaveTeam func(heroId int64, isKick bool)

	// 队伍创建时间
	createTime time.Time

	// 队伍过期时间
	expireTime time.Time

	// 开始挑战的保护结束时间
	protectEndTime time.Time

	// 帮派id，非0表示需要改帮派的玩家
	guildId int64

	// 队伍展示缓存
	team4ShowCache []byte
	// 队伍详情缓存
	teamDetailCache []byte
	teamDetailMsg   pbutil.Buffer
	// 聊天记录
	chatRecord *collection.RingList
}

func (t *secret_tower_team) SendChatRecord(member *secret_tower_team_member) {
	if recLen := t.chatRecord.Length(); recLen > 0 {
		recs := make([]*shared_proto.SecretTowerChatRecordProto, 0, recLen)
		t.chatRecord.ReverseRange(func(v interface{}) (toContinue bool) {
			recs = append(recs, v.(*shared_proto.SecretTowerChatRecordProto))
			return true
		})
		member.SendMsg(secret_tower.NewS2cTeamHistoryTalkMsg(recs))
	}
}

func (t *secret_tower_team) IsLeader(id int64) bool {
	return t.leader.Id() == id
}

func (t *secret_tower_team) MemberCount() int {
	return len(t.members)
}

func (t *secret_tower_team) Members() []*secret_tower_team_member {
	return t.members
}

func (t *secret_tower_team) GuildId() int64 {
	return t.guildId
}

func (t *secret_tower_team) IsFull() bool {
	return uint64(t.MemberCount()) >= t.towerData.MaxAttackerCount
}

func (t *secret_tower_team) ProtectEndTime() time.Time {
	return t.protectEndTime
}

func (t *secret_tower_team) IsMemberNotEnough() bool {
	return uint64(len(t.members)) < t.towerData.MinAttackerCount
}

func (t *secret_tower_team) Add(member *secret_tower_team_member, ctime time.Time) (isExist, teamFull bool) {
	_, existMember := t.GetMember(member.Id())
	if existMember != nil {
		isExist = true
		logrus.Debugln("玩家已经在队伍中了")
		return
	}

	if t.IsFull() {
		// 超出队伍人数了
		teamFull = true
		logrus.Debugln("队伍已满")
		return
	}

	t.add(member, ctime)
	return
}

func (t *secret_tower_team) add(member *secret_tower_team_member, ctime time.Time) {
	t.members = append(t.members, member)
	t.changeProtectEndTime(ctime)
	t.invalidCache()

	return
}

func (t *secret_tower_team) changeProtectEndTime(ctime time.Time) {
	t.protectEndTime = ctime.Add(t.towerData.StartProtectDuration)
}

// 成员离开
func (t *secret_tower_team) Leave(id int64) (destroy, notFound bool) {
	idx, member := t.GetMember(id)
	if member == nil {
		notFound = true
		return
	}

	t.removeMember(idx, false)

	if len(t.members) <= 0 {
		destroy = true
		t.destroy()
		return
	}

	if !t.IsLeader(id) {
		return
	}

	// 找一个战斗力最高的
	suc := t.resetLeaderId()
	if !suc {
		// 队长离开导致没有非挑战模式的玩家接任队伍导致队伍解散
		t.broadcast(secret_tower.TEAM_DESTROYED_BECAUSE_OF_LEADER_LEAVE_S2C)

		for idx := len(t.members) - 1; idx >= 0; idx-- {
			t.removeMember(idx, false)
		}

		destroy = true
		t.destroy()
	}

	return
}

// 踢出成员
func (t *secret_tower_team) Kick(id int64) (beenKickMember *secret_tower_team_member) {
	idx, beenKickMember := t.GetMember(id)
	if beenKickMember == nil {
		return
	}

	t.removeMember(idx, true)

	return
}

// 移除队员
func (t *secret_tower_team) removeMember(idx int, isKick bool) {
	member := t.members[idx]

	if idx != len(t.members)-1 {
		copy(t.members[idx:], t.members[idx+1:])
	}
	t.members = t.members[:len(t.members)-1]
	t.onHeroLeaveTeam(member.Id(), isKick)
	t.invalidCache()
}

// 摧毁这个队伍
func (t *secret_tower_team) destroy() {
	for _, member := range t.Members() {
		t.onHeroLeaveTeam(member.Id(), false)
	}

	t.removeTeam(t)
}

// 移动成员
func (t *secret_tower_team) Move(id int64, up bool) (opSuccess, failAndIsFirst, notFound bool) {
	oldIdx, member := t.GetMember(id)
	if member == nil {
		notFound = true
		return
	}

	newIdx := 0

	if up {
		if oldIdx == 0 {
			failAndIsFirst = true
			return
		}
		newIdx = oldIdx - 1
	} else {
		if oldIdx == len(t.members)-1 {
			return
		}
		newIdx = oldIdx + 1
	}

	t.members[oldIdx], t.members[newIdx] = t.members[newIdx], t.members[oldIdx]

	t.invalidCache()

	opSuccess = true

	return
}

func (t *secret_tower_team) UpdateMemberPos(ids []int64) bool {

	if len(t.members) != len(ids) {
		// 长度不一致
		return false
	}

	// 互相在对方的列表中
	newMembers := make([]*secret_tower_team_member, len(t.members))
	for _, m := range t.members {
		index := i64.GetIndex(ids, m.Id())
		if index < 0 || newMembers[index] != nil {
			return false
		}

		newMembers[index] = m
	}

	t.members = newMembers
	t.invalidCache()
	return true
}

// 获得队伍成员，如果member为空，表示没找到
func (t *secret_tower_team) GetMember(id int64) (idx int, member *secret_tower_team_member) {
	for idx, member := range t.members {
		if member.Id() == id {
			return idx, member
		}
	}

	return -1, nil
}

func (t *secret_tower_team) resetLeaderId() (suc bool) {
	var maxFightAmountMember *secret_tower_team_member
	for _, member := range t.members {
		if member.mode == shared_proto.TowerTeamMode_HELP {
			continue
		}

		if maxFightAmountMember == nil {
			maxFightAmountMember = member
		} else if member.FightAmount() > maxFightAmountMember.FightAmount() {
			maxFightAmountMember = member
		}
	}

	if maxFightAmountMember == nil {
		return
	}

	t.leader = maxFightAmountMember
	t.invalidCache()

	return true
}

func (t *secret_tower_team) broadcast(buffer pbutil.Buffer) {
	for _, member := range t.members {
		member.SendMsg(buffer)
	}
}

func (t *secret_tower_team) broadcastIgnore(buffer pbutil.Buffer, ignore int64) {
	for _, member := range t.members {
		if member.Id() != ignore {
			member.SendMsg(buffer)
		}
	}
}

func (t *secret_tower_team) update(heroDataService iface.HeroDataService) {
	for _, m := range t.members {
		m.Update(heroDataService)
	}
	t.invalidCache()
}

func (t *secret_tower_team) invalidCache() {
	t.team4ShowCache = nil
	t.teamDetailCache = nil
	t.teamDetailMsg = nil
}

// 序列化用于队伍展示，可以多线程访问
func (t *secret_tower_team) Encode4Show() []byte {
	if t.team4ShowCache == nil {
		t.team4ShowCache = must.Marshal(t.encode4Show())
	}

	return t.team4ShowCache
}

func (t *secret_tower_team) encode4Show() *shared_proto.SecretTeamShowProto {
	proto := &shared_proto.SecretTeamShowProto{}

	proto.TeamId = i64.Int32(t.teamId)
	proto.Leader = t.leader.EncodeClient()
	proto.SecretTowerId = u64.Int32(t.towerData.Id)
	proto.CurMemberCount = int32(len(t.members))
	proto.MaxMemberCount = u64.Int32(t.towerData.MaxAttackerCount)
	proto.GuildId = i64.Int32(t.guildId)
	proto.CreateTime = timeutil.Marshal32(t.createTime)

	for _, member := range t.Members() {
		snapshot := member.heroSnapshotGetter()
		if t.IsLeader(member.Id()) {
			continue
		}
		proto.MemberTowerFloor = append(proto.MemberTowerFloor, u64.Int32(snapshot.TowerMaxFloor))
	}

	return proto
}

func (t *secret_tower_team) TeamDetailMsg() pbutil.Buffer {
	if t.teamDetailMsg == nil {
		t.teamDetailMsg = secret_tower.NewS2cRequestTeamDetailMsg(t.EncodeDetail()).Static()
	}

	return t.teamDetailMsg
}

// 序列化用于队内展示
func (t *secret_tower_team) EncodeDetail() []byte {
	if t.teamDetailCache == nil {
		t.teamDetailCache = must.Marshal(t.encodeDetail())
	}

	return t.teamDetailCache
}

func (t *secret_tower_team) encodeDetail() *shared_proto.SecretTeamDetailProto {
	proto := &shared_proto.SecretTeamDetailProto{}

	proto.TeamId = i64.Int32(t.teamId)
	proto.LeaderId = t.leader.IdBytes()

	proto.Members = make([][]byte, 0, len(t.members))
	for _, member := range t.members {
		proto.Members = append(proto.Members, member.EncodeClient())
	}

	proto.SecretTowerId = u64.Int32(t.towerData.Id)
	proto.GuildId = i64.Int32(t.guildId)

	proto.ProtectEndTime = timeutil.Marshal32(t.ProtectEndTime())

	return proto
}

func newSecretTowerTeamMember(hero *entity.Hero, sendMsgFunc func(int64, pbutil.Buffer), mode shared_proto.TowerTeamMode, heroSnapshotService iface.HeroSnapshotService) *secret_tower_team_member {

	id := hero.Id()

	heroSnapshot := heroSnapshotService.Get(id)
	if heroSnapshot == nil {
		heroSnapshot = heroSnapshotService.NewSnapshot(hero)
	}

	heroSnapshotGetter := func() *snapshotdata.HeroSnapshot {
		snap := heroSnapshotService.GetFromCache(id)

		if snap != nil && snap != heroSnapshot {
			heroSnapshot = snap
		}

		return heroSnapshot
	}

	member := &secret_tower_team_member{
		IdHolder:           hero.IdHolder,
		heroSnapshotGetter: heroSnapshotGetter,
		sendMsgFunc:        sendMsgFunc,
		mode:               mode,
	}

	member.syncCaptains(hero)

	return member
}

// 密室成员
type secret_tower_team_member struct {
	idbytes.IdHolder

	heroSnapshotGetter func() *snapshotdata.HeroSnapshot // 这个方法不会返回空的

	sendMsgFunc func(int64, pbutil.Buffer)

	// 模式
	mode shared_proto.TowerTeamMode

	// 出战的武将
	captainRaces            []shared_proto.Race // 出战武将的职业
	captainTotalFightAmount uint64              // 出战武将的战斗力

	saidWords uint64    // 气泡说话
	saidTime  time.Time // 最后一条气泡的时间

	// SecretTowerTeamMemberProto
	cacheClientBytes []byte
}

func (m *secret_tower_team_member) HeroNameAndGuildInfo() (heroName string, guildId int64, guildName, guildFlagName string) {
	snapshot := m.heroSnapshotGetter()
	if g := snapshot.Guild(); g != nil {
		return snapshot.Name, g.Id, g.Name, g.FlagName
	}
	return snapshot.Name, 0, "", ""
}

func (m *secret_tower_team_member) ChangeMode(newMode shared_proto.TowerTeamMode) (noChange bool) {
	if m.mode == newMode {
		return false
	}

	m.mode = newMode

	return true
}

func (m *secret_tower_team_member) FightAmount() uint64 {
	return m.captainTotalFightAmount
}

func (m *secret_tower_team_member) SendMsg(buffer pbutil.Buffer) {
	m.sendMsgFunc(m.Id(), buffer)
}

func (m *secret_tower_team_member) Update(heroDataService iface.HeroDataService) {
	heroDataService.FuncNotError(m.Id(), func(hero *entity.Hero) (heroChanged bool) {
		m.syncCaptains(hero)
		return
	})
	m.cacheClientBytes = nil
}

func (m *secret_tower_team_member) syncCaptains(hero *entity.Hero) {
	pveTroop := hero.PveTroop(shared_proto.PveTroopType_DUNGEON)
	if pveTroop == nil {
		return
	}

	m.captainRaces = pveTroop.CaptainRaces()
	m.captainTotalFightAmount = pveTroop.TotalFullFightAmount()
}

func (m *secret_tower_team_member) EncodeClient() []byte {
	if m.cacheClientBytes == nil {
		m.cacheClientBytes = must.Marshal(m.encodeClient())
	}
	return m.cacheClientBytes
}

func (m *secret_tower_team_member) encodeClient() *shared_proto.SecretTowerTeamMemberProto {
	proto := &shared_proto.SecretTowerTeamMemberProto{}

	proto.Hero = m.heroSnapshotGetter().EncodeClient()
	proto.Mode = m.mode
	proto.CaptainRace = m.captainRaces
	proto.FightAmount = u64.Int32(m.captainTotalFightAmount)
	proto.SaidTime = timeutil.Marshal32(m.saidTime)
	proto.SaidWords = u64.Int32(m.saidWords)

	return proto
}

func (m *secret_tower_team_member) Encode4Result() *shared_proto.SecretMemberResultProto {
	proto := &shared_proto.SecretMemberResultProto{}

	proto.Id = m.IdBytes()

	snapshot := m.heroSnapshotGetter()
	proto.Name = snapshot.Name
	proto.Head = snapshot.Head
	proto.Level = u64.Int32(snapshot.Level)
	if g := snapshot.Guild(); g != nil {
		proto.GuildId = i64.Int32(g.Id)
		proto.GuildName = g.Name
		proto.GuildFlagName = g.FlagName
	}

	return proto
}
