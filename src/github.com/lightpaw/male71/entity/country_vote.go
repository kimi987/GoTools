package entity

import (
	"github.com/lightpaw/male7/pb/server_proto"
	"github.com/lightpaw/male7/pb/shared_proto"
	"github.com/lightpaw/male7/util/timeutil"
	"github.com/lightpaw/male7/util/u64"
	"time"
)

func NewCountryChangeNameVote() *CountryChangeNameVote {
	c := &CountryChangeNameVote{}
	return c
}

type CountryChangeNameVote struct {
	id            int32
	newName       string
	endTime       time.Time
	agreeCount    uint64
	disagreeCount uint64
	completed     bool
}

func (c *CountryChangeNameVote) encode() (p *shared_proto.CountryChangeNameVoteProto) {
	p = &shared_proto.CountryChangeNameVoteProto{}
	p.Id = c.id
	p.NewName = c.newName
	p.EndTime = timeutil.Marshal32(c.endTime)
	p.AgreeCount = u64.Int32(c.agreeCount)
	p.DisagreeCount = u64.Int32(c.disagreeCount)

	return p
}

func (c *CountryChangeNameVote) encodeServer() (p *server_proto.CountryChangeNameVoteServerProto) {
	p = &server_proto.CountryChangeNameVoteServerProto{}
	p.Id = c.id
	p.NewName = c.newName
	p.EndTime = timeutil.Marshal64(c.endTime)
	p.AgreeCount = c.agreeCount
	p.DisagreeCount = c.disagreeCount
	p.Completed = c.completed

	return p
}

func (c *CountryChangeNameVote) unmarshal(p *server_proto.CountryChangeNameVoteServerProto) {
	if p == nil {
		return
	}

	c.id = p.Id
	c.newName = p.NewName
	c.endTime = timeutil.Unix64(p.EndTime)
	c.agreeCount = p.AgreeCount
	c.disagreeCount = p.DisagreeCount
	c.completed = p.Completed
}

func (c *Country) StartChangeNameVote(newName string, endTime time.Time) {
	c.lock.Lock()
	defer c.lock.Unlock()

	c.vote = NewCountryChangeNameVote()
	c.vote.id = (timeutil.Marshal32(endTime) << 8) | u64.Int32(c.id)
	c.vote.endTime = endTime
	c.vote.newName = newName
}

func (c *Country) OnChangeNameVoteEnd() (ended bool, changed bool) {
	c.lock.Lock()
	defer c.lock.Unlock()

	if c.vote.completed {
		return
	}

	c.vote.completed = true
	ended = true

	if c.vote.agreeCount <= c.vote.disagreeCount {
		return
	}

	c.name = c.vote.newName
	changed = true
	return
}

func (c *Country) VoteEndTime() time.Time {
	c.lock.RLock()
	defer c.lock.RUnlock()

	return c.vote.endTime
}

func (c *Country) VoteId() int32 {
	c.lock.RLock()
	defer c.lock.RUnlock()

	return c.vote.id
}

func (c *Country) VoteNewName() string {
	c.lock.RLock()
	defer c.lock.RUnlock()

	return c.vote.newName
}

func (c *Country) UpdateChangeNameVoteCount(id int32, toUpdate int, agree bool) (newCount uint64) {
	c.lock.Lock()
	defer c.lock.Unlock()

	if toUpdate == 0 {
		return
	}

	if c.vote.id != id {
		return
	}

	if agree {
		c.vote.agreeCount = u64.AddInt(c.vote.agreeCount, toUpdate)
		newCount = c.vote.agreeCount
	} else {
		c.vote.disagreeCount = u64.AddInt(c.vote.disagreeCount, toUpdate)
		newCount = c.vote.disagreeCount
	}
	return
}

func (c *Country) InChangeNameVoteDuration(ctime time.Time) bool {
	c.lock.RLock()
	defer c.lock.RUnlock()

	return !c.vote.completed && ctime.Before(c.vote.endTime)
}

func (c *Country) IsChangeNameVoteCompleted() bool {
	c.lock.RLock()
	defer c.lock.RUnlock()

	return c.vote.completed || timeutil.IsZero(c.vote.endTime)
}
