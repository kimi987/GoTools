package entity

import (
	"github.com/lightpaw/male7/pb/server_proto"
	"github.com/lightpaw/male7/pb/shared_proto"
	"math"
)

func newBools() *hero_bools {
	return &hero_bools{}
}

type hero_bools struct {
	serverBools []uint64
}

func (hero *Hero) Bools() *hero_bools {
	return hero.bools
}

func (b *hero_bools) encodeClient() int32 {
	if len(b.serverBools) > 0 {
		return int32(b.serverBools[0] & math.MaxInt32)
	}
	return 0
}

func (b *hero_bools) encodeArrayClient() []uint32 {
	if len(b.serverBools) > 0 {

		is := make([]uint32, len(b.serverBools)*2)

		for i, v := range b.serverBools {
			if v != 0 {
				lowIdx := i * 2
				highIdx := lowIdx + 1
				is[lowIdx] = uint32(v & math.MaxUint32)
				is[highIdx] = uint32((v >> 32) & math.MaxUint32)
			}
		}
		return is
	}
	return nil
}

func (b *hero_bools) unmarshal(p *server_proto.HeroServerProto) {
	b.serverBools = p.HeroBools
}

func (b *hero_bools) ensureLen(n int) {
	if len(b.serverBools) < n {
		newBools := make([]uint64, n)
		copy(newBools, b.serverBools)
		b.serverBools = newBools
	}
}

var boolLen = len(shared_proto.HeroBoolType_name)

func (b *hero_bools) Get(t shared_proto.HeroBoolType) bool {
	index := int(t)
	if index < boolLen {
		subIndex := index / 64
		bitOffset := index % 64

		if subIndex < len(b.serverBools) {
			bools := b.serverBools[subIndex]
			return bools&(1<<uint64(bitOffset)) != 0
		}
	}

	return false
}

func (b *hero_bools) TrySet(t shared_proto.HeroBoolType) bool {
	if !b.Get(t) {
		b.SetTrue(t)
		return true
	}
	return false
}

func (b *hero_bools) SetTrue(t shared_proto.HeroBoolType) {
	index := int(t)
	if index < boolLen {
		subIndex := index / 64
		bitOffset := index % 64

		b.ensureLen(subIndex + 1)
		bools := b.serverBools[subIndex]
		b.serverBools[subIndex] = bools | 1<<uint64(bitOffset)
	}
}

func (b *hero_bools) SetFalse(t shared_proto.HeroBoolType) {
	index := int(t)
	if index < boolLen {
		subIndex := index / 64
		bitOffset := index % 64

		b.ensureLen(subIndex + 1)
		bools := b.serverBools[subIndex]
		b.serverBools[subIndex] = bools & ^(1 << uint64(bitOffset))
	}
}
