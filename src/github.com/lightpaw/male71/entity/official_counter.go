package entity

import (
	"github.com/lightpaw/male7/util/u64"
	"github.com/lightpaw/male7/pb/shared_proto"
)

func newCaptainOfficialCounter() *CaptainOfficialCounter {
	return &CaptainOfficialCounter {
		officials: make(map[uint64]*captainOfficial),
	}
}

type CaptainOfficialCounter struct {
	// 各种现有官职
	officials map[uint64]*captainOfficial
}

func (c *CaptainOfficialCounter) encodeServer() map[uint64]*shared_proto.Int32ArrayProto {
	m := make(map[uint64]*shared_proto.Int32ArrayProto)
	for id, o := range c.officials {
		p := &shared_proto.Int32ArrayProto{}
		for idx, _ := range o.noViewed {
			p.V = append(p.V, idx)
		}
		m[id] = p
	}
	return m
}

func (c *CaptainOfficialCounter) unmarshal(configMap map[uint64]uint64, m map[uint64]*shared_proto.Int32ArrayProto) {
	for id, count := range configMap {
		if count == 0 {
			continue
		}
		c.officials[id] = &captainOfficial {
			officialPos: make([]uint64, count),
			noViewed: make(map[int32]struct{}),
		}
	}
	for id, view := range m {
		if o := c.officials[id]; o != nil {
			for _, idx := range view.V {
				o.noViewed[idx] = struct{}{}
			}
		}
	}
}

func (c *CaptainOfficialCounter) encode() []*shared_proto.Int32PairInt32ArrayProto {
	p := []*shared_proto.Int32PairInt32ArrayProto{}
	for id, official := range c.officials {
		proto := &shared_proto.Int32PairInt32ArrayProto {
			K: u64.Int32(id),
		}
		for idx, _ := range official.noViewed {
			proto.V = append(proto.V, idx)
		}
		p = append(p, proto)
	}
	return p
}

func (c *CaptainOfficialCounter) Get(id uint64) *captainOfficial {
	return c.officials[id]
}

func (c *CaptainOfficialCounter) Set(id, idx, captainId uint64) {
	if official := c.officials[id]; official != nil {
		official.set(idx, captainId)
	}
}

func (c *CaptainOfficialCounter) OfficialPositionCount(id uint64) uint64 {
	official := c.officials[id]
	if official == nil {
		return 0
	}
	return official.positionCount()
}

func (c *CaptainOfficialCounter) change(configMap map[uint64]uint64) {
	for id, count := range configMap {
		if count == 0 {
			continue
		}
		if official := c.officials[id]; official == nil {
			official := &captainOfficial {
				officialPos: make([]uint64, count),
				noViewed: make(map[int32]struct{}),
			}
			for i, max := int32(0), u64.Int32(count); i < max; i++ {
				official.noViewed[i] = struct{}{}
			}
			c.officials[id] = official
		} else {
			for i, max := u64.Int32(official.capacity()), u64.Int32(count); i < max; i++ {
				official.noViewed[i] = struct{}{}
			}
			newArr := make([]uint64, count)
			copy(newArr, official.officialPos)
			official.officialPos = newArr
		}
	}
}

// 每一种官职
type captainOfficial struct {
	officialPos []uint64 // 每种官职占的武将下标坑,存的是武将id
	noViewed map[int32]struct{} // 哪些下标未预览过
}

func (c *captainOfficial) capacity() uint64 {
	return u64.FromInt(len(c.officialPos))
}

func (c *captainOfficial) checkEmpty(idx uint64) bool {
	return idx < c.capacity() && c.officialPos[idx] == 0
}

func (c *captainOfficial) get(idx uint64) uint64 {
	if idx < c.capacity() {
		return c.officialPos[idx]
	}
	return 0
}

func (c *captainOfficial) set(idx, captainId uint64) {
	if idx < c.capacity() {
		c.officialPos[idx] = captainId
	}
}

func (c *captainOfficial) positionCount() uint64 {
	count := uint64(0)
	for _, cid := range c.officialPos {
		if cid != 0 {
			count++
		}
	}
	return count
}
