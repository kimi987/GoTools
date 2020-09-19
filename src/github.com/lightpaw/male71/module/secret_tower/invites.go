package secret_tower

import (
	"github.com/lightpaw/male7/gen/iface"
	"github.com/lightpaw/male7/gen/pb/secret_tower"
	"github.com/lightpaw/male7/util/i64"
	"github.com/lightpaw/male7/util/timeutil"
	"sync"
	"time"
)

// 邀请
func newTeamInvites(worldService iface.WorldService) *team_invites {
	return &team_invites{
		worldService:                   worldService,
		heroBeenInviteMap:              make(map[int64]invite_me_teams),
		teamInviteHeroMap:              make(map[int64]invite_heroids),
		heroLastGetInviteMeTeamTimeMap: make(map[int64]time.Time),
	}
}

type team_invites struct {
	worldService iface.WorldService

	// 玩家邀请的列表
	heroBeenInviteMap map[int64]invite_me_teams

	// 队伍邀请人的列表
	teamInviteHeroMap map[int64]invite_heroids

	// 每个玩家最后一次获得邀请我的队伍列表的时间
	heroLastGetInviteMeTeamTimeMap map[int64]time.Time

	sync.RWMutex
}

type invite_me_team struct {
	teamId     int64     // 队伍id
	inviteTime time.Time // 该队伍邀请我的时间
}

type invite_me_teams []*invite_me_team // 邀请我的队伍
func (teams invite_me_teams) RemoveTeam(teamId int64) invite_me_teams {
	for idx, team := range teams {
		if team.teamId == teamId {
			if idx != len(teams)-1 {
				// 不是最后一个
				copy(teams[idx:], teams[idx+1:])
			}

			return teams[:len(teams)-1]
		}
	}

	return teams
}

type invite_heroids []int64 // 队伍邀请的玩家id
func (i invite_heroids) Contains(heroId int64) bool {
	return i64.Contains(i, heroId)
}

// 是有有邀请我加入这支队伍
func (t *team_invites) RefreshLastGetInviteMeTeamTime(heroId int64, ctime time.Time) bool {
	t.Lock()
	defer t.Unlock()

	inviteMeTeams := t.heroBeenInviteMap[heroId]
	if len(inviteMeTeams) <= 0 {
		return false
	}

	t.heroLastGetInviteMeTeamTimeMap[heroId] = ctime

	return true
}

func (t *team_invites) HasInviteJoinTeam(heroId int64, teamId int64) bool {
	t.RLock()
	defer t.RUnlock()

	inviteHeroIds := t.teamInviteHeroMap[teamId]
	if len(inviteHeroIds) <= 0 {
		return false
	}

	return inviteHeroIds.Contains(heroId)
}

// 邀请我加入了某支队伍
func (t *team_invites) OnInviteJoinTeam(heroId int64, teamId int64, ctime time.Time) {
	var teamCount int
	func() {
		t.Lock()
		defer t.Unlock()

		inviteHeroIds := t.teamInviteHeroMap[teamId]
		if inviteHeroIds == nil {
			inviteHeroIds = []int64{heroId}
		} else {
			if i64.Contains(inviteHeroIds, heroId) {
				return
			}

			inviteHeroIds = append(inviteHeroIds, heroId)
		}

		t.teamInviteHeroMap[teamId] = inviteHeroIds

		inviteMeTeamIds := t.heroBeenInviteMap[heroId]
		if inviteMeTeamIds == nil {
			inviteMeTeamIds = make([]*invite_me_team, 0, 1)
		}
		inviteMeTeamIds = append(inviteMeTeamIds, &invite_me_team{teamId: teamId, inviteTime: ctime})
		t.heroBeenInviteMap[heroId] = inviteMeTeamIds

		teamCount = len(inviteMeTeamIds)
	}()

	if teamCount > 0 {
		// 新的邀请，发消息
		t.worldService.Send(heroId, secret_tower.NewS2cReceiveInviteMsg(int32(teamCount), true))
	}
}

// 邀请我们加入了某支队伍
func (t *team_invites) OnInviteAllJoinTeam(heroIds []int64, teamId int64, ctime time.Time) {
	teamCounts := make([]int, len(heroIds))
	func() {
		t.Lock()
		defer t.Unlock()

		inviteHeroIds := t.teamInviteHeroMap[teamId]
		for i := 0; i < len(heroIds); i++ {
			heroId := heroIds[i]

			if inviteHeroIds == nil {
				inviteHeroIds = []int64{heroId}
			} else {
				if i64.Contains(inviteHeroIds, heroId) {
					continue
				}

				inviteHeroIds = append(inviteHeroIds, heroId)
			}

			inviteMeTeamIds := t.heroBeenInviteMap[heroId]
			if inviteMeTeamIds == nil {
				inviteMeTeamIds = make([]*invite_me_team, 0, 1)
			}
			inviteMeTeamIds = append(inviteMeTeamIds, &invite_me_team{teamId: teamId, inviteTime: ctime})
			t.heroBeenInviteMap[heroId] = inviteMeTeamIds

			teamCounts[i] = len(inviteMeTeamIds)
		}

		t.teamInviteHeroMap[teamId] = inviteHeroIds
	}()

	for i := 0; i < len(heroIds); i++ {
		if teamCount := teamCounts[i]; teamCount > 0 {
			heroId := heroIds[i]
			// 新的邀请，发消息
			t.worldService.Send(heroId, secret_tower.NewS2cReceiveInviteMsg(int32(teamCount), true))
		}
	}
}

var noInviteMsg = secret_tower.NewS2cReceiveInviteMsg(0, false).Static()

// 加入了队伍
func (t *team_invites) OnJoinTeam(heroId int64) {
	t.Lock()
	defer t.Unlock()

	inviteMeTeams := t.heroBeenInviteMap[heroId]
	if len(inviteMeTeams) <= 0 {
		return
	}

	delete(t.heroBeenInviteMap, heroId)
	delete(t.heroLastGetInviteMeTeamTimeMap, heroId)
	// 发消息
	t.worldService.Send(heroId, noInviteMsg)

	for _, team := range inviteMeTeams {
		inviteHeroIds := t.teamInviteHeroMap[team.teamId]
		if len(inviteHeroIds) > 0 {
			oldInviteCount := len(inviteHeroIds)
			inviteHeroIds = i64.RemoveIfPresent(inviteHeroIds, heroId)
			leftInviteCount := len(inviteHeroIds)
			if oldInviteCount != leftInviteCount {
				if leftInviteCount <= 0 {
					delete(t.teamInviteHeroMap, team.teamId)
				} else {
					t.teamInviteHeroMap[team.teamId] = inviteHeroIds
				}
			}
		}
	}
}

// 踢出了队伍
func (t *team_invites) OnKickOutTeam(heroId int64, teamId int64) {
	t.Lock()
	defer t.Unlock()

	inviteHeroIds := t.teamInviteHeroMap[teamId]
	if len(inviteHeroIds) <= 0 {
		return
	}

	oldInviteCount := len(inviteHeroIds)
	inviteHeroIds = i64.RemoveIfPresent(inviteHeroIds, heroId)
	leftInviteCount := len(inviteHeroIds)
	if oldInviteCount != leftInviteCount {
		if leftInviteCount <= 0 {
			delete(t.teamInviteHeroMap, teamId)
		} else {
			t.teamInviteHeroMap[teamId] = inviteHeroIds
		}

		inviteMeTeams := t.heroBeenInviteMap[heroId]
		oldInviteMeTeamCount := len(inviteMeTeams)
		inviteMeTeams = inviteMeTeams.RemoveTeam(teamId)
		leftInviteMeTeamCount := len(inviteMeTeams)
		if oldInviteMeTeamCount != leftInviteMeTeamCount {
			if leftInviteMeTeamCount <= 0 {
				delete(t.heroBeenInviteMap, heroId)
				delete(t.heroLastGetInviteMeTeamTimeMap, heroId)

				t.worldService.Send(heroId, noInviteMsg)
			} else {
				t.heroBeenInviteMap[heroId] = inviteMeTeams
				t.worldService.Send(heroId, secret_tower.NewS2cReceiveInviteMsg(t.inviteMeTeamAndHaveNewUnderLock(heroId)))
			}
		}
	}
}

// 移除队伍
func (t *team_invites) OnDestroyTeam(teamId int64) {
	t.Lock()
	defer t.Unlock()

	inviteHeroIds := t.teamInviteHeroMap[teamId]
	if len(inviteHeroIds) <= 0 {
		return
	}

	delete(t.teamInviteHeroMap, teamId)

	for _, heroId := range inviteHeroIds {
		inviteMeTeams := t.heroBeenInviteMap[heroId]
		if len(inviteMeTeams) > 0 {
			oldInviteCount := len(inviteMeTeams)
			inviteMeTeams = inviteMeTeams.RemoveTeam(teamId)
			leftInviteCount := len(inviteMeTeams)
			if oldInviteCount != leftInviteCount {
				if leftInviteCount <= 0 {
					delete(t.heroBeenInviteMap, heroId)
					delete(t.heroLastGetInviteMeTeamTimeMap, heroId)

					t.worldService.Send(heroId, noInviteMsg)
				} else {
					t.heroBeenInviteMap[heroId] = inviteMeTeams
					t.worldService.Send(heroId, secret_tower.NewS2cReceiveInviteMsg(t.inviteMeTeamAndHaveNewUnderLock(heroId)))
				}
			}
		}
	}
}

// 遍历邀请我的队伍
func (t *team_invites) WalkInviteMeTeam(heroId int64, f func(teamId int64) (toContinue bool)) {
	t.RLock()
	defer t.RUnlock()

	inviteMeTeams := t.heroBeenInviteMap[heroId]
	if len(inviteMeTeams) > 0 {
		for _, team := range inviteMeTeams {
			if !f(team.teamId) {
				break
			}
		}
	}
}

// 邀请我的队伍数量
func (t *team_invites) InviteMeTeamAndHaveNew(heroId int64) (count int32, haveNew bool) {
	t.RLock()
	defer t.RUnlock()

	return t.inviteMeTeamAndHaveNewUnderLock(heroId)
}

// 邀请我的队伍数量
func (t *team_invites) inviteMeTeamAndHaveNewUnderLock(heroId int64) (count int32, haveNew bool) {
	inviteMeTeams := t.heroBeenInviteMap[heroId]
	if len(inviteMeTeams) <= 0 {
		return 0, false
	}

	if lastTime := t.heroLastGetInviteMeTeamTimeMap[heroId]; !timeutil.IsZero(lastTime) {
		for _, team := range inviteMeTeams {
			if team.inviteTime.After(lastTime) {
				return int32(len(inviteMeTeams)), true
			}
		}
	}

	return int32(len(inviteMeTeams)), false
}
