package rank

import (
	"github.com/lightpaw/male7/module/rank/rankface"
	"github.com/lightpaw/male7/module/rank/ranklist"
	"github.com/lightpaw/male7/pb/shared_proto"
	"github.com/lightpaw/male7/util/u64"
	"strings"
)

func (m *RankModule) CountryOfficial(size int, name string, countryId uint64, official shared_proto.CountryOfficialType) (heros []*shared_proto.HeroBasicSnapshotProto) {
	holder := m.rankHolder(shared_proto.RankType_Tower)
	if holder == nil {
		return
	}
	holder.LockFunc(func(h ranklist.LockedRankHolder) {
		h.Walk(func(list rankface.RankList) {
			list.Walk(func(obj rankface.RankObj) {
				hero := obj.EncodeHeroSnapshotProto()

				if hero.Basic.CountryId != u64.Int32(countryId) {
					return
				}
				if hero.Basic.Official != official {
					return
				}

				name = strings.TrimSpace(name)
				if name != "" && strings.Index(hero.Basic.Name, name) < 0 {
					return
				}

				heros = append(heros, hero)
				if len(heros) >= size {
					return
				}
			})
			if len(heros) >= size {
				return
			}
		})
	})

	return
}
