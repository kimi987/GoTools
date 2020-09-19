package entity

import (
	"github.com/lightpaw/male7/pb/server_proto"
	"github.com/lightpaw/male7/pb/shared_proto"
	"github.com/lightpaw/male7/util/u64"
)

func NewHeroTeach() *HeroTeach {
	h := &HeroTeach{}
	h.PassedChapterIds = make(map[uint64]struct{})
	h.CollectedChapterIds = make(map[uint64]struct{})

	return h
}

type HeroTeach struct {
	PassedChapterIds    map[uint64]struct{}
	CollectedChapterIds map[uint64]struct{}
}

func (h *HeroTeach) encode() *shared_proto.HeroTeachProto {
	p := &shared_proto.HeroTeachProto{}
	p.PassedChapterIds = u64.MapKey2Int32Arrary(h.PassedChapterIds)
	p.CollectedChapterIds = u64.MapKey2Int32Arrary(h.CollectedChapterIds)

	return p
}

func (h *HeroTeach) encodeServer() *server_proto.HeroTeachServerProto {
	p := &server_proto.HeroTeachServerProto{}
	p.PassedChapterIds = u64.MapKey2Uint64Array(h.PassedChapterIds)
	p.CollectedChapterIds = u64.MapKey2Uint64Array(h.CollectedChapterIds)
	return p
}

func (h *HeroTeach) unmarshal(p *server_proto.HeroTeachServerProto) {
	if p == nil {
		return
	}
	h.PassedChapterIds = u64.Uint64ArrayToMapKey(p.PassedChapterIds)
	h.CollectedChapterIds = u64.Uint64ArrayToMapKey(p.CollectedChapterIds)
}

