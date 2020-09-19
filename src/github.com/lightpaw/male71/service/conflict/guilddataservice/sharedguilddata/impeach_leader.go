package sharedguilddata

import (
	"github.com/lightpaw/male7/entity/npcid"
	"github.com/lightpaw/male7/pb/server_proto"
	"github.com/lightpaw/male7/pb/shared_proto"
	"github.com/lightpaw/male7/util/i64"
	"github.com/lightpaw/male7/util/idbytes"
	"github.com/lightpaw/male7/util/timeutil"
	"github.com/lightpaw/male7/util/u64"
	"sort"
	"time"
)

func createImpeachLeader(leaderId, impeachMemberId int64, impeachStartTime, impeachEndTime time.Time, memberMap map[int64]*GuildMember, extraCandidateCount, npcLeaderVote uint64) *impeach_leader {

	isNpcLeader := npcid.IsNpcId(leaderId)

	npcLeader := leaderId
	var candidates []int64
	count := int(extraCandidateCount)
	if !isNpcLeader {
		npcLeader = 0
		// User leader，候选人为原来盟主，弹劾者，加另外1个7日贡献度最高的玩家
		candidates = i64.AddIfAbsent(candidates, leaderId)
		candidates = i64.AddIfAbsent(candidates, impeachMemberId)
		if count > 1 {
			for _, c := range getUserCandidate(memberMap, count-1, impeachMemberId, leaderId) {
				candidates = i64.AddIfAbsent(candidates, c)
			}
		}
	} else {
		// NPC leader，候选人为NPC盟主，加另外2个7日贡献度最高的玩家
		candidates = i64.AddIfAbsent(candidates, leaderId)
		for _, c := range getUserCandidate(memberMap, count, npcLeader) {
			candidates = i64.AddIfAbsent(candidates, c)
		}
	}

	result := newImpeachLeader(npcLeader, impeachStartTime, impeachEndTime, candidates, impeachMemberId, npcLeaderVote)

	if isNpcLeader {
		// 所有的NPC都给NPC帮主投票
		for _, m := range memberMap {
			if npcid.IsNpcId(m.Id()) {
				result.vote(m.Id(), npcLeader)
			}
		}
	}

	return result
}

type contribution7Slice []*GuildMember

func (a contribution7Slice) Len() int      { return len(a) }
func (a contribution7Slice) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a contribution7Slice) Less(i, j int) bool {
	ia := a[i].ContributionAmount7()
	ja := a[j].ContributionAmount7()
	if ia == ja {
		// 后入盟者优先，入盟时间大的在前面
		return a[i].createTime.After(a[j].createTime)
	}

	return ia > ja // 大的在前
}

func getUserCandidate(memberMap map[int64]*GuildMember, count int, ignore ...int64) (candidates []int64) {

	if count <= 0 {
		return
	}

	var array []*GuildMember
	for _, m := range memberMap {
		if npcid.IsNpcId(m.Id()) || i64.Contains(ignore, m.Id()) {
			continue
		}

		array = append(array, m)
	}

	sort.Sort(contribution7Slice(array))
	for _, m := range array {
		candidates = append(candidates, m.Id())

		if len(candidates) >= count {
			break
		}
	}

	return
}

func newImpeachLeader(npcLeader int64, impeachStartTime, impeachEndTime time.Time, candidates []int64, impeachMemberId int64, npcLeaderVote uint64) *impeach_leader {
	m := &impeach_leader{}
	m.impeachStartTime = impeachStartTime
	m.impeachEndTime = impeachEndTime
	m.candidates = candidates
	m.impeachMemberId = impeachMemberId
	m.voteMap = make(map[int64]int64)

	m.npcLeader = npcLeader
	m.npcLeaderVote = npcLeaderVote

	return m
}

type impeach_leader struct {
	// 弹劾开始时间
	impeachStartTime time.Time

	// 弹劾结束时间
	impeachEndTime time.Time

	// 弹劾候选人
	candidates []int64

	// 弹劾发起者
	impeachMemberId int64

	// 选票
	voteMap map[int64]int64

	// npc 盟主票数，0表示不是npc盟主
	npcLeader int64

	npcLeaderVote uint64
}

func (m *impeach_leader) encodeClient(memberMap map[int64]*GuildMember) *shared_proto.GuildImpeachProto {
	proto := &shared_proto.GuildImpeachProto{}

	proto.ImpeachEndTime = timeutil.Marshal32(m.impeachEndTime)

	proto.ImpeachMemberId = idbytes.ToBytes(m.impeachMemberId)

	for _, id := range m.candidates {
		proto.Candidates = append(proto.Candidates, idbytes.ToBytes(id))
		proto.Points = append(proto.Points, u64.Int32(m.getScore(id, memberMap)))
	}

	for voteHero, voteTarget := range m.voteMap {
		proto.VoteHeros = append(proto.VoteHeros, idbytes.ToBytes(voteHero))
		proto.VoteTarget = append(proto.VoteTarget, idbytes.ToBytes(voteTarget))
	}

	return proto
}

func (m *impeach_leader) encodeServer() *server_proto.GuildImpeachServerProto {
	proto := &server_proto.GuildImpeachServerProto{}

	proto.ImpeachEndTime = timeutil.Marshal64(m.impeachEndTime)
	proto.Candidates = m.candidates
	proto.ImpeachMemberId = m.impeachMemberId

	for voteHero, voteTarget := range m.voteMap {
		proto.VoteHeros = append(proto.VoteHeros, voteHero)
		proto.VoteTarget = append(proto.VoteTarget, voteTarget)
	}

	return proto
}

func (m *impeach_leader) IsNpcLeader() bool {
	return npcid.IsNpcId(m.npcLeader)
}

func (m *impeach_leader) IsValidCandidate(heroId int64) bool {
	return i64.Contains(m.candidates, heroId)
}

func (m *impeach_leader) vote(selfId, voteTargetId int64) {
	m.voteMap[selfId] = voteTargetId
}

func (m *impeach_leader) removeMember(heroId int64) {

	// 删掉这个人的投票
	delete(m.voteMap, heroId)

	// 如果是候选人，删掉候选人
	if m.IsValidCandidate(heroId) {
		m.candidates = i64.LeftShiftRemoveIfPresent(m.candidates, heroId)

		// 删掉所有投给这个人的票
		for k, v := range m.voteMap {
			if v == heroId {
				delete(m.voteMap, k)
			}
		}
	}
}

func (m *impeach_leader) tryImpeach(memberMap map[int64]*GuildMember) (newLeader *GuildMember) {
	if npcid.IsNpcId(m.npcLeader) {
		return m.tryImpeachNpcLeader(memberMap)
	} else {
		return m.tryImpeachUserLeader(memberMap)
	}
}

// 弹劾玩家盟主
func (m *impeach_leader) getMaxScoreCandidate(memberMap map[int64]*GuildMember) (newLeader *GuildMember) {
	// 若投票持续时间结束，则得票数最多的候选人获胜（若并列，后入盟者获胜）

	var maxScore int
	for _, candidate := range m.candidates {
		cm := memberMap[candidate]
		if cm == nil {
			continue
		}

		score := m.getVoteMember(candidate)
		if maxScore < score {
			newLeader = cm
			maxScore = score
		} else if maxScore == score {
			if newLeader == nil {
				newLeader = cm
			} else {
				// 比较2个人的入盟时间，后入盟的优先
				if newLeader.createTime.Before(cm.createTime) {
					newLeader = cm
				}
			}
		}
	}

	return
}

// 弹劾玩家盟主
func (m *impeach_leader) tryImpeachUserLeader(memberMap map[int64]*GuildMember) (newLeader *GuildMember) {
	// 一旦某个候选人得票数超过联盟当前人数的一半，则投票结束，该候选人获胜
	halfMemberCount := len(memberMap) / 2

	for _, candidate := range m.candidates {
		cm := memberMap[candidate]
		if cm == nil {
			continue
		}

		c := m.getVoteMember(candidate)
		if c > halfMemberCount {
			return cm
		}
	}

	return nil
}

// 弹劾npc盟主
func (m *impeach_leader) tryImpeachNpcLeader(memberMap map[int64]*GuildMember) (newLeader *GuildMember) {
	// 一旦某个玩家得票数超过原NPC盟主，则投票结束，弹劾成功

	npcLeaderScore := m.getVoteClassScore(m.npcLeader, memberMap)

	// 遍历所有的玩家候选人，获取他的得票数，如果得票数超过Npc盟主，则弹劾成功
	for _, candidate := range m.candidates {
		if m.npcLeader == candidate {
			continue
		}

		cm := memberMap[candidate]
		if cm == nil {
			continue
		}

		score := m.getVoteClassScore(candidate, memberMap)
		if score > npcLeaderScore {
			return cm
		}
	}

	return nil
}

func (m *impeach_leader) getScore(candidate int64, memberMap map[int64]*GuildMember) uint64 {
	if npcid.IsNpcId(m.npcLeader) {
		return m.getVoteClassScore(candidate, memberMap)
	} else {
		return uint64(m.getVoteMember(candidate))
	}
}

func (m *impeach_leader) getVoteClassScore(candidate int64, memberMap map[int64]*GuildMember) uint64 {
	var score uint64
	if npcid.IsNpcId(m.npcLeader) && m.npcLeader == candidate {
		score = m.npcLeaderVote
	}

	for heroId, voteTargetId := range m.voteMap {
		if voteTargetId == candidate {
			member := memberMap[heroId]
			if member != nil {
				score += member.classLevelData.VoteScore
			}
		}
	}

	return score
}

func (m *impeach_leader) getVoteMember(candidate int64) int {
	var score int
	for _, voteTargetId := range m.voteMap {
		if voteTargetId == candidate {
			// 一人一票
			score++
		}
	}

	return score
}
