package entity

import (
	"github.com/lightpaw/male7/pb/server_proto"
	"github.com/lightpaw/male7/pb/shared_proto"
	"github.com/lightpaw/male7/util/idbytes"
	"github.com/lightpaw/male7/util/u64"
	"time"
	"github.com/lightpaw/male7/util/timeutil"
)

type RelationType uint8

const (
	Unkown RelationType = 0
	Friend RelationType = 1
	Black  RelationType = 2
)

type relation struct {
	createTime   time.Time
	relationType RelationType
}

func newRelation() *hero_relation {
	return &hero_relation{
		relationMap:  make(map[int64]*relation),
		enemyMap:     make(map[int64]time.Time),
		importantMap: make(map[int64]time.Time),
	}
}

type hero_relation struct {
	relationMap  map[int64]*relation
	enemyMap     map[int64]time.Time
	importantMap map[int64]time.Time

	friendCount uint64
	blackCount  uint64
}

func (r *hero_relation) unmarshal(proto *server_proto.HeroRelationServerProto, ctime time.Time) {
	if proto == nil {
		return
	}

	if len(proto.Friend) == len(proto.FriendCreateTime) {
		for k, t := range proto.Friend {
			r.AddFriend(t, timeutil.Unix64(proto.FriendCreateTime[k]))
		}
	} else {
		for _, t := range proto.Friend {
			r.AddFriend(t, ctime)
		}
	}

	if len(proto.Black) == len(proto.BlackCreateTime) {
		for k, t := range proto.Black {
			r.AddBlack(t, timeutil.Unix64(proto.BlackCreateTime[k]))
		}
	} else {
		for _, t := range proto.Black {
			r.AddFriend(t, ctime)
		}
	}

	for id, t := range proto.Enemy {
		r.AddEnemy(id, timeutil.Unix64(t))
	}

	for id, t := range proto.Important {
		r.SetImportantFriend(id, timeutil.Unix64(t))
	}
}

func (r *hero_relation) encode() *server_proto.HeroRelationServerProto {
	proto := &server_proto.HeroRelationServerProto{}

	for k, v := range r.relationMap {
		switch v.relationType {
		case Friend:
			proto.Friend = append(proto.Friend, k)
			proto.FriendCreateTime = append(proto.FriendCreateTime, timeutil.Marshal64(v.createTime))
		case Black:
			proto.Black = append(proto.Black, k)
			proto.BlackCreateTime = append(proto.BlackCreateTime, timeutil.Marshal64(v.createTime))
		}
	}

	proto.Enemy = make(map[int64]int64)
	for k, v := range r.enemyMap {
		proto.Enemy[k] = timeutil.Marshal64(v)
	}

	proto.Important = make(map[int64]int64)
	for k, v := range r.importantMap {
		proto.Important[k] = timeutil.Marshal64(v)
	}

	return proto
}

func (r *hero_relation) encodeClient() *shared_proto.HeroRelationProto {
	proto := &shared_proto.HeroRelationProto{}

	for k, v := range r.relationMap {
		switch v.relationType {
		case Friend:
			proto.Friend = append(proto.Friend, idbytes.ToBytes(k))
			proto.FriendCreateTime = append(proto.FriendCreateTime, timeutil.Marshal32(v.createTime))
		case Black:
			proto.Black = append(proto.Black, idbytes.ToBytes(k))
			proto.BlackCreateTime = append(proto.BlackCreateTime, timeutil.Marshal32(v.createTime))
		}
	}

	for k, v := range r.enemyMap {
		p := &shared_proto.HeroEnemyProto{EnemyId: idbytes.ToBytes(k), AddTime: timeutil.Marshal32(v)}
		proto.Enemy = append(proto.Enemy, p)
	}

	for k, v := range r.importantMap {
		p := &shared_proto.HeroImportantFriendProto{FriendId: idbytes.ToBytes(k), SetTime: timeutil.Marshal32(v)}
		proto.Important = append(proto.Important, p)
	}

	return proto
}

func (r *hero_relation) GetRelation(heroId int64) RelationType {
	if target := r.relationMap[heroId]; target != nil {
		return target.relationType
	}
	return 0
}

func (r *hero_relation) IsEnemy(heroId int64) bool {
	if _, ok := r.enemyMap[heroId]; ok {
		return true
	}
	return false
}

func (r *hero_relation) RemoveRelation(heroId int64) RelationType {
	rt := r.GetRelation(heroId)
	switch rt {
	case Friend:
		r.friendCount = u64.Sub(r.friendCount, 1)
	case Black:
		r.blackCount = u64.Sub(r.blackCount, 1)
	default:
		return rt
	}

	delete(r.relationMap, heroId)
	return rt
}

func (r *hero_relation) RemoveEnemy(heroId int64) {
	delete(r.enemyMap, heroId)
}

func (r *hero_relation) CancelImportantFriend(heroId int64) bool {
	if _, ok := r.importantMap[heroId]; !ok {
		return false
	}
	delete(r.importantMap, heroId)
	return true
}

func (r *hero_relation) AddFriend(heroId int64, ctime time.Time) {
	rt := r.GetRelation(heroId)
	switch rt {
	case Friend:
		return
	case Black:
		r.blackCount = u64.Sub(r.blackCount, 1)
	}

	r.relationMap[heroId] = &relation{ctime, Friend}
	r.friendCount++
}

func (r *hero_relation) AddBlack(heroId int64, ctime time.Time) {
	rt := r.GetRelation(heroId)
	switch rt {
	case Friend:
		r.friendCount = u64.Sub(r.friendCount, 1)
	case Black:
		return
	}

	r.relationMap[heroId] = &relation{ctime, Black}
	r.blackCount++
}

func (r *hero_relation) AddEnemy(heroId int64, t time.Time) {
	r.enemyMap[heroId] = t
}

func (r *hero_relation) SetImportantFriend(heroId int64, t time.Time) bool {
	if _, ok := r.importantMap[heroId]; ok {
		return false
	}
	r.importantMap[heroId] = t
	return true
}

func (r *hero_relation) RelationIds() []int64 {
	ids := make([]int64, 0, len(r.relationMap))
	for k := range r.relationMap {
		ids = append(ids, k)
	}
	return ids
}

func (r *hero_relation) RelationAndEnemyIds() []int64 {
	ids := make([]int64, 0, len(r.relationMap))
	for k := range r.relationMap {
		ids = append(ids, k)
	}
	for k := range r.enemyMap {
		if _, ok := r.relationMap[k]; ok {
			continue
		}
		ids = append(ids, k)
	}
	return ids
}

func (r *hero_relation) FriendIds() []int64 {
	return r.idsByType(Friend)
}

func (r *hero_relation) BlackIds() []int64 {
	return r.idsByType(Black)
}

func (r *hero_relation) EnemyIds() []int64 {
	return r.idsByEnemy()
}

func (r *hero_relation) FriendAndEnemyIds() (ids []int64) {
	if len(r.enemyMap) <= 0 {
		return r.idsByType(Friend)
	}
	if r.friendCount <= 0 {
		return r.idsByEnemy()
	}

	ids = r.idsByEnemy()
	for k, v := range r.relationMap {
		if v.relationType != Friend {
			continue
		}
		if r.IsEnemy(k) {
			continue
		}
		ids = append(ids, k)
	}
	return
}

func (r *hero_relation) idsByEnemy() (ids []int64) {
	if len(r.enemyMap) <= 0 {
		return
	}

	ids = make([]int64, 0, len(r.enemyMap))
	for k := range r.enemyMap {
		ids = append(ids, k)
	}
	return
}

func (r *hero_relation) idsByType(t RelationType) (ids []int64) {
	var count uint64
	switch t {
	case Friend:
		count = r.friendCount
	case Black:
		count = r.blackCount
	default:
		return
	}

	if count > 0 {
		ids = make([]int64, 0, count)
		for k, v := range r.relationMap {
			if t == v.relationType {
				ids = append(ids, k)
			}
		}
	}
	return ids
}

func (r *hero_relation) RelationCount(t RelationType) uint64 {
	switch t {
	case Friend:
		return r.friendCount
	case Black:
		return r.blackCount
	}
	return 0
}
