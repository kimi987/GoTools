package secret_tower

import (
	"github.com/lightpaw/logrus"
	"github.com/lightpaw/male7/config/towerdata"
	"github.com/lightpaw/male7/gen/iface"
	"github.com/lightpaw/male7/gen/pb/secret_tower"
	"github.com/lightpaw/male7/util/u64"
	"github.com/lightpaw/pbutil"
	"sync"
	"time"
)

func newTeamManager(datas iface.ConfigDatas, worldService iface.WorldService, heroDataService iface.HeroDataService) *team_manager {
	m := &team_manager{}

	teamIdGen := int64(0)
	m.newTeamId = func() int64 {
		teamIdGen++
		return teamIdGen
	}

	m.allTeams = make(map[int64]*secret_tower_team)
	m.heroIdAndTeamIdMap = Newhero_join_team_map()
	m.levelTeams = make([]*tower_teams, len(datas.GetSecretTowerDataArray()))
	m.invites = newTeamInvites(worldService)
	m.heroDataService = heroDataService

	for _, data := range datas.GetSecretTowerDataArray() {
		m.levelTeams[data.Id-1] = newTowerTeams(data)
	}

	return m
}

type team_manager struct {
	// 队伍id生成
	newTeamId func() int64

	heroDataService iface.HeroDataService

	// 所有队伍
	allTeams     map[int64]*secret_tower_team
	sync.RWMutex // 所有队伍的读写锁

	// 所以队伍根据密室id分组
	levelTeams []*tower_teams

	// 玩家id跟队伍的map，用来判断玩家是不是有队伍，玩家是不是在自己的队伍
	heroIdAndTeamIdMap *hero_join_team_map

	// 队伍数量的缓存
	teamCountMsgCache pbutil.Buffer

	invites *team_invites
}

func (m *team_manager) addTeam(team *secret_tower_team) {
	func() {
		m.Lock()
		defer m.Unlock()
		m.allTeams[team.teamId] = team
	}()
	m.mustTowerTeams(team.towerData).addTeam(team)
	m.teamCountMsgCache = nil
}

func (m *team_manager) removeTeam(team *secret_tower_team) {
	func() {
		m.Lock()
		defer m.Unlock()
		delete(m.allTeams, team.teamId)
	}()
	m.mustTowerTeams(team.towerData).removeTeam(team)
	m.invites.OnDestroyTeam(team.teamId)
	m.teamCountMsgCache = nil
}

func (m *team_manager) getRTeam(teamId int64) r_secret_tower_team {
	team := m.getTeam(teamId)
	if team == nil {
		return nil
	}
	return team
}

func (m *team_manager) getTeam(teamId int64) *secret_tower_team {
	m.RLock()
	defer m.RUnlock()
	return m.allTeams[teamId]
}

func (m *team_manager) mustTowerTeams(data *towerdata.SecretTowerData) *tower_teams {
	return m.levelTeams[data.Id-1]
}

func (m *team_manager) onHeroJoinTeam(heroId int64, team *secret_tower_team) {
	m.heroIdAndTeamIdMap.Set(heroId, team)
	// 策划确认了，所有邀请不清除
	//m.invites.OnJoinTeam(heroId)
}

func (m *team_manager) onHeroLeaveTeam(heroId int64, isKick bool) {
	if isKick {
		team := m.getHeroJoinTeam(heroId)
		if team != nil {
			m.heroIdAndTeamIdMap.Remove(heroId)
			m.invites.OnKickOutTeam(heroId, team.teamId)
		} else {
			logrus.Errorln("玩家离开队伍的时候竟然没找到玩家队伍信息")
		}
	} else {
		m.heroIdAndTeamIdMap.Remove(heroId)
	}
}

func (m *team_manager) isHeroJoinTeam(heroId int64) bool {
	return m.getHeroJoinTeam(heroId) != nil
}

func (m *team_manager) getHeroJoinTeam(heroId int64) *secret_tower_team {
	team, _ := m.heroIdAndTeamIdMap.Get(heroId)
	return team
}

func (m *team_manager) getTeamCountMsgCache() pbutil.Buffer {
	if m.teamCountMsgCache == nil {
		towerIdArray := make([]int32, 0, len(m.levelTeams))
		towerTeamCountArray := make([]int32, 0, len(m.levelTeams))
		for _, levelTeams := range m.levelTeams {
			teamCount := len(levelTeams.teams)
			if teamCount <= 0 {
				continue
			}

			towerIdArray = append(towerIdArray, u64.Int32(levelTeams.data.Id))
			towerTeamCountArray = append(towerTeamCountArray, int32(teamCount))
		}

		m.teamCountMsgCache = secret_tower.NewS2cRequestTeamCountMsg(towerIdArray, towerTeamCountArray)
	}
	return m.teamCountMsgCache
}

func (m *team_manager) update(ctime time.Time) {
	for _, teams := range m.levelTeams {
		teams.update(m.heroDataService, ctime)
	}
}

func newTowerTeams(data *towerdata.SecretTowerData) *tower_teams {
	return &tower_teams{
		data:  data,
		teams: make(map[int64]*secret_tower_team, 32),
	}
}

// 层级数据
type tower_teams struct {
	data  *towerdata.SecretTowerData
	teams map[int64]*secret_tower_team

	// 队伍列表缓存跟缓存过期时间
	teamListCache pbutil.Buffer
}

func (t *tower_teams) addTeam(team *secret_tower_team) {
	t.teams[team.teamId] = team
	t.invalidTeamListCache()
}

func (t *tower_teams) removeTeam(team *secret_tower_team) {
	delete(t.teams, team.teamId)
	t.invalidTeamListCache()
}

func (t *tower_teams) update(heroDataService iface.HeroDataService, ctime time.Time) {
	for _, team := range t.teams {
		if ctime.After(team.expireTime) {
			// 过期了
			team.broadcast(secret_tower.TEAM_EXPIRED_S2C)
			team.destroy()
			continue
		}
		team.update(heroDataService)
	}
	t.invalidTeamListCache()
}

func (t *tower_teams) invalidTeamListCache() {
	t.teamListCache = nil
}

func (t *tower_teams) getTeamListCache() pbutil.Buffer {
	if t.teamListCache == nil {
		teamList := make([][]byte, 0, len(t.teams))
		for _, team := range t.teams {
			teamList = append(teamList, team.Encode4Show())
		}
		t.teamListCache = secret_tower.NewS2cRequestTeamListMsg(u64.Int32(t.data.Id), teamList)
	}

	return t.teamListCache
}
