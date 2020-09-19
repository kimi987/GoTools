package secret_tower

import (
	"github.com/lightpaw/male7/gen/iface"
	"github.com/lightpaw/male7/gen/pb/secret_tower"
	"github.com/lightpaw/male7/pb/shared_proto"
	"github.com/lightpaw/male7/entity"
	"github.com/lightpaw/male7/service/conflict/heroservice/herolock"
)

//gogen:iface c2s_request_team_count
func (m *SecretTowerModule) ProcessRequestTeamCount(hc iface.HeroController) {
	m.queueWaitFunc(func() {
		hc.Send(m.manager.getTeamCountMsgCache())
	})
}

//gogen:iface c2s_request_team_list
func (m *SecretTowerModule) ProcessRequestTeamList(proto *secret_tower.C2SRequestTeamListProto, hc iface.HeroController) {
	_, err := m.requestTeamList(uint64(proto.GetSecretTowerId()), hc)
	if err != nil {
		hc.Send(err.ErrMsg())
		return
	}
}

//gogen:iface c2s_create_team
func (m *SecretTowerModule) ProcessCreateTeam(proto *secret_tower.C2SCreateTeamProto, hc iface.HeroController) {
	_, _, err := m.createTeam(uint64(proto.GetSecretTowerId()), proto.GetIsGuild(), hc)
	if err != nil {
		hc.Send(err.ErrMsg())
		return
	}
}

//gogen:iface c2s_join_team
func (m *SecretTowerModule) ProcessJoinTeam(proto *secret_tower.C2SJoinTeamProto, hc iface.HeroController) {
	if proto.GetTeamId() != 0 {
		_, _, err := m.joinTeam(int64(proto.GetTeamId()), hc)
		if err != nil {
			hc.Send(err.ErrMsg())
			return
		}
	} else {
		_, _, err := m.autoJoinTeam(uint64(proto.GetSecretTowerId()), hc)
		if err != nil {
			hc.Send(err.ErrMsg())
			return
		}
	}

}

//gogen:iface c2s_leave_team
func (m *SecretTowerModule) ProcessLeaveTeam(hc iface.HeroController) {
	_, _, err := m.leaveTeam(hc)
	if err != nil {
		hc.Send(err.ErrMsg())
		return
	}
}

//gogen:iface c2s_kick_member
func (m *SecretTowerModule) ProcessKickMember(proto *secret_tower.C2SKickMemberProto, hc iface.HeroController) {
	_, _, err := m.kickMember(proto.GetId(), hc)
	if err != nil {
		hc.Send(err.ErrMsg())
		return
	}
}

//gogen:iface c2s_move_member
func (m *SecretTowerModule) ProcessMoveMember(proto *secret_tower.C2SMoveMemberProto, hc iface.HeroController) {
	_, _, err := m.moveMember(proto.GetId(), proto.GetUp(), hc)
	if err != nil {
		hc.Send(err.ErrMsg())
		return
	}
}

//gogen:iface c2s_update_member_pos
func (m *SecretTowerModule) ProcessUpdateMemberPos(proto *secret_tower.C2SUpdateMemberPosProto, hc iface.HeroController) {
	_, _, err := m.updateMemberPos(proto.GetId(), hc)
	if err != nil {
		hc.Send(err.ErrMsg())
		return
	}
}

//gogen:iface c2s_change_mode
func (m *SecretTowerModule) ProcessChangeMode(proto *secret_tower.C2SChangeModeProto, hc iface.HeroController) {
	_, _, err := m.changeMode(shared_proto.TowerTeamMode(proto.GetMode()), hc)
	if err != nil {
		hc.Send(err.ErrMsg())
		return
	}
}

//gogen:iface c2s_invite
func (m *SecretTowerModule) ProcessInvite(proto *secret_tower.C2SInviteProto, hc iface.HeroController) {
	_, _, err := m.invite(proto.GetId(), hc)
	if err != nil {
		hc.Send(err.ErrMsg())
		return
	}
}

//gogen:iface c2s_invite_all
func (m *SecretTowerModule) ProcessInviteAll(proto *secret_tower.C2SInviteAllProto, hc iface.HeroController) {
	_, _, err := m.inviteAll(proto.Id, hc)
	if err != nil {
		hc.Send(err.ErrMsg())
		return
	}
}

//gogen:iface c2s_request_invite_list
func (m *SecretTowerModule) ProcessRequestInviteList(hc iface.HeroController) {
	m.requestInviteList(hc)
}

//gogen:iface c2s_request_team_detail
func (m *SecretTowerModule) ProcessRequestTeamDetail(hc iface.HeroController) {
	_, _, err := m.requestTeamDetail(hc)
	if err != nil {
		hc.Send(err.ErrMsg())
		return
	}
}

//gogen:iface c2s_start_challenge
func (m *SecretTowerModule) ProcessStartChallenge(hc iface.HeroController) {
	_, _, err := m.startChallenge(hc)
	if err != nil {
		hc.Send(err.ErrMsg())
		return
	}
}

//gogen:iface c2s_quick_query_team_basic
func (m *SecretTowerModule) ProcessQuickQueryTeamBasic(proto *secret_tower.C2SQuickQueryTeamBasicProto, hc iface.HeroController) {
	m.quickQueryTeamBasic(proto, hc)
}

//gogen:iface c2s_change_guild_mode
func (m *SecretTowerModule) ProcessChangeGuildMode(hc iface.HeroController) {
	_, _, err := m.changeGuildMode(hc)
	if err != nil {
		hc.Send(err.ErrMsg())
		return
	}
}

//gogen:iface c2s_list_record
func (m *SecretTowerModule) ProcessListRecord(hc iface.HeroController) {
	hc.FuncWithSend(func(hero *entity.Hero, result herolock.LockResult) {
		result.Add(secret_tower.NewS2cListRecordMsg(hero.SecretTower().GetRecords()))
		result.Ok()
	})
}
