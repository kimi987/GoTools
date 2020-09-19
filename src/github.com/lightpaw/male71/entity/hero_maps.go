package entity

import (
	"github.com/lightpaw/male7/pb/server_proto"
	"github.com/lightpaw/male7/util/u64"
)

func newHeroMaps() *hero_maps {
	return &hero_maps{
		heroMaps:      make(map[server_proto.HeroMapCategory]*heromap),
		keyMaps:       make(map[server_proto.HeroMapCategory]*herokeys),
	}
}

type hero_maps struct {
	heroMaps map[server_proto.HeroMapCategory]*heromap
	keyMaps  map[server_proto.HeroMapCategory]*herokeys
}

func (maps *hero_maps) unmarshal(p *server_proto.HeroServerProto) {
	for _, proto := range p.HeroMaps {
		m := maps.heroMaps[proto.Category]
		if m == nil {
			continue
		}
		u64.CopyMapTo(m.internalMap, proto.DataMap)
	}
	for _, proto := range p.HeroKeys {
		m := maps.keyMaps[proto.Category]
		if m == nil {
			continue
		}
		for _, k := range proto.Keys {
			m.Add(k)
		}
	}
}

func (maps *hero_maps) encode() (mps []*server_proto.HeroMapServerProto, kps []*server_proto.HeroKeysServerProto) {
	for _, m := range maps.heroMaps {
		if len(m.internalMap) > 0 {
			mps = append(mps, m.encode())
		}
	}
	for _, m := range maps.keyMaps {
		if len(m.internalMap) > 0 {
			kps = append(kps, m.encode())
		}
	}
	return
}

func (maps *hero_maps) getOrCreateMap(category server_proto.HeroMapCategory, isDaily bool) *heromap {
	m := maps.heroMaps[category]
	if m == nil {
		m = &heromap{
			category:    category,
			internalMap: make(map[uint64]uint64),
			isDaily:     isDaily,
		}
		maps.heroMaps[category] = m
	}
	return m
}

func (maps *hero_maps) getOrCreateKeys(category server_proto.HeroMapCategory, isDaily bool) *herokeys {
	m := maps.keyMaps[category]
	if m == nil {
		m = &herokeys{
			category:    category,
			internalMap: make(map[uint64]struct{}),
			isDaily:     isDaily,
		}
		maps.keyMaps[category] = m
	}
	return m
}

func (m *hero_maps) resetDaily() {
	for _, v := range m.heroMaps {
		if !v.isDaily {
			continue
		}
		v.internalMap = make(map[uint64]uint64)
	}
	for _, v := range m.keyMaps {
		if !v.isDaily {
			continue
		}
		v.internalMap = make(map[uint64]struct{})
	}
}

type heromap struct {
	category    server_proto.HeroMapCategory
	internalMap map[uint64]uint64
	isDaily     bool
}

func (m *heromap) Size() int {
	return len(m.internalMap)
}

func (m *heromap) Clear() {
	m.internalMap = make(map[uint64]uint64)
}

func (m *heromap) Exist(key uint64) bool {
	_, exist := m.internalMap[key]
	return exist
}

func (m *heromap) Get(key uint64) uint64 {
	return m.internalMap[key]
}

func (m *heromap) Set(key, value uint64) {
	m.internalMap[key] = value
}

func (m *heromap) Increse(key uint64) uint64 {
	value := m.internalMap[key]
	value++
	m.internalMap[key] = value

	return value
}

func (m *heromap) encode() *server_proto.HeroMapServerProto {
	proto := &server_proto.HeroMapServerProto{}
	proto.Category = m.category
	proto.DataMap = m.internalMap

	return proto
}

type herokeys struct {
	category    server_proto.HeroMapCategory
	internalMap map[uint64]struct{}
	isDaily     bool
}

func (m *herokeys) Exist(key uint64) bool {
	_, exist := m.internalMap[key]
	return exist
}

func (m *herokeys) Add(key uint64) {
	m.internalMap[key] = struct{}{}
}

func (m *herokeys) Remove(key uint64) {
	delete(m.internalMap, key)
}

func (m *herokeys) Size() int {
	return len(m.internalMap)
}

func (m *herokeys) encode() *server_proto.HeroKeysServerProto {
	proto := &server_proto.HeroKeysServerProto{}
	proto.Category = m.category
	for k := range m.internalMap {
		proto.Keys = append(proto.Keys, k)
	}

	return proto
}
