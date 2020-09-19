package secret_tower

import (
	"github.com/lightpaw/male7/gen/ifacemock"
	. "github.com/onsi/gomega"
	"testing"
	"time"
)

// 邀请的测试

func NewTeamInvites() *team_invites {
	return newTeamInvites(ifacemock.WorldService)
}

func TestTeam_invites(t *testing.T) {
	RegisterTestingT(t)

	invites := NewTeamInvites()

	heroId1 := int64(1)
	heroId2 := int64(2)

	invites.OnDestroyTeam(1)
	invites.OnJoinTeam(heroId1)

	Ω(invites.HasInviteJoinTeam(heroId1, 1)).Should(BeEquivalentTo(false))
	Ω(invites.InviteMeTeamAndHaveNew(heroId1)).Should(BeEquivalentTo(0))

	invites.OnInviteJoinTeam(heroId1, 3, time.Time{})
	Ω(invites.HasInviteJoinTeam(heroId1, 1)).Should(BeEquivalentTo(false))
	Ω(invites.HasInviteJoinTeam(heroId1, 3)).Should(BeEquivalentTo(true))
	Ω(invites.InviteMeTeamAndHaveNew(heroId1)).Should(BeEquivalentTo(1))

	invites.OnDestroyTeam(3)
	Ω(invites.InviteMeTeamAndHaveNew(heroId1)).Should(BeEquivalentTo(0))
	Ω(invites.HasInviteJoinTeam(heroId1, 1)).Should(BeEquivalentTo(false))
	Ω(invites.HasInviteJoinTeam(heroId1, 3)).Should(BeEquivalentTo(false))

	invites.OnInviteJoinTeam(heroId1, 3, time.Time{})
	Ω(invites.InviteMeTeamAndHaveNew(heroId1)).Should(BeEquivalentTo(1))
	invites.OnInviteJoinTeam(heroId1, 4, time.Time{})
	Ω(invites.InviteMeTeamAndHaveNew(heroId1)).Should(BeEquivalentTo(2))
	Ω(invites.HasInviteJoinTeam(heroId1, 1)).Should(BeEquivalentTo(false))
	Ω(invites.HasInviteJoinTeam(heroId1, 3)).Should(BeEquivalentTo(true))
	Ω(invites.HasInviteJoinTeam(heroId1, 4)).Should(BeEquivalentTo(true))

	invites.OnJoinTeam(heroId1)
	Ω(invites.InviteMeTeamAndHaveNew(heroId1)).Should(BeEquivalentTo(0))
	Ω(invites.HasInviteJoinTeam(heroId1, 1)).Should(BeEquivalentTo(false))
	Ω(invites.HasInviteJoinTeam(heroId1, 3)).Should(BeEquivalentTo(false))
	Ω(invites.HasInviteJoinTeam(heroId1, 4)).Should(BeEquivalentTo(false))

	invites.OnInviteJoinTeam(heroId1, 3, time.Time{})
	Ω(invites.InviteMeTeamAndHaveNew(heroId1)).Should(BeEquivalentTo(1))
	invites.OnInviteJoinTeam(heroId1, 4, time.Time{})
	Ω(invites.InviteMeTeamAndHaveNew(heroId1)).Should(BeEquivalentTo(2))
	invites.OnInviteJoinTeam(heroId2, 3, time.Time{})
	Ω(invites.InviteMeTeamAndHaveNew(heroId2)).Should(BeEquivalentTo(1))
	invites.OnInviteJoinTeam(heroId2, 5, time.Time{})
	Ω(invites.InviteMeTeamAndHaveNew(heroId2)).Should(BeEquivalentTo(2))

	walkTeamIds := []int64{3, 4}
	invites.WalkInviteMeTeam(heroId1, func(teamId int64) (toContinue bool) {
		Ω(teamId).Should(BeEquivalentTo(walkTeamIds[0]))

		walkTeamIds = walkTeamIds[1:]
		return true
	})

	walkTeamIds = []int64{3, 5}
	invites.WalkInviteMeTeam(heroId2, func(teamId int64)(toContinue bool) {
		Ω(teamId).Should(BeEquivalentTo(walkTeamIds[0]))

		walkTeamIds = walkTeamIds[1:]
		return true
	})

	Ω(invites.HasInviteJoinTeam(heroId1, 1)).Should(BeEquivalentTo(false))
	Ω(invites.HasInviteJoinTeam(heroId1, 3)).Should(BeEquivalentTo(true))
	Ω(invites.HasInviteJoinTeam(heroId1, 4)).Should(BeEquivalentTo(true))
	Ω(invites.HasInviteJoinTeam(heroId1, 5)).Should(BeEquivalentTo(false))

	Ω(invites.HasInviteJoinTeam(heroId2, 1)).Should(BeEquivalentTo(false))
	Ω(invites.HasInviteJoinTeam(heroId2, 3)).Should(BeEquivalentTo(true))
	Ω(invites.HasInviteJoinTeam(heroId2, 4)).Should(BeEquivalentTo(false))
	Ω(invites.HasInviteJoinTeam(heroId2, 5)).Should(BeEquivalentTo(true))

	invites.OnJoinTeam(heroId1)
	Ω(invites.HasInviteJoinTeam(heroId1, 1)).Should(BeEquivalentTo(false))
	Ω(invites.HasInviteJoinTeam(heroId1, 3)).Should(BeEquivalentTo(false))
	Ω(invites.HasInviteJoinTeam(heroId1, 4)).Should(BeEquivalentTo(false))
	Ω(invites.HasInviteJoinTeam(heroId1, 5)).Should(BeEquivalentTo(false))

	Ω(invites.HasInviteJoinTeam(heroId2, 1)).Should(BeEquivalentTo(false))
	Ω(invites.HasInviteJoinTeam(heroId2, 3)).Should(BeEquivalentTo(true))
	Ω(invites.HasInviteJoinTeam(heroId2, 4)).Should(BeEquivalentTo(false))
	Ω(invites.HasInviteJoinTeam(heroId2, 5)).Should(BeEquivalentTo(true))

	invites.OnInviteJoinTeam(heroId1, 3, time.Time{})
	Ω(invites.InviteMeTeamAndHaveNew(heroId1)).Should(BeEquivalentTo(1))
	invites.OnInviteJoinTeam(heroId1, 4, time.Time{})
	Ω(invites.InviteMeTeamAndHaveNew(heroId1)).Should(BeEquivalentTo(2))

	invites.OnDestroyTeam(3)
	Ω(invites.HasInviteJoinTeam(heroId1, 1)).Should(BeEquivalentTo(false))
	Ω(invites.HasInviteJoinTeam(heroId1, 3)).Should(BeEquivalentTo(false))
	Ω(invites.HasInviteJoinTeam(heroId1, 4)).Should(BeEquivalentTo(true))
	Ω(invites.HasInviteJoinTeam(heroId1, 5)).Should(BeEquivalentTo(false))

	Ω(invites.HasInviteJoinTeam(heroId2, 1)).Should(BeEquivalentTo(false))
	Ω(invites.HasInviteJoinTeam(heroId2, 3)).Should(BeEquivalentTo(false))
	Ω(invites.HasInviteJoinTeam(heroId2, 4)).Should(BeEquivalentTo(false))
	Ω(invites.HasInviteJoinTeam(heroId2, 5)).Should(BeEquivalentTo(true))
}
