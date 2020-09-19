package herosnapshot

import (
	"context"
	"github.com/lightpaw/male7/entity"
	"github.com/lightpaw/male7/gen/iface"
	"github.com/lightpaw/male7/gen/ifacemock"
	"github.com/lightpaw/male7/service/herosnapshot/snapshotdata"
	. "github.com/onsi/gomega"
	"github.com/pkg/errors"
	"sync"
	"testing"
	"github.com/lightpaw/male7/gen/pb/util"
	"github.com/lightpaw/male7/pb/shared_proto"
	"github.com/lightpaw/male7/util/idbytes"
	"github.com/lightpaw/male7/util/must"
)

type panicDb struct {
	iface.DbService
	t *testing.T
}

func (p *panicDb) LoadHero(ctx context.Context, id int64) (*entity.Hero, error) {
	p.t.Fatal("should not call db.LoadHero")
	return nil, errors.New("should not be here")
}

func TestHeroSnapshotService_Encode(t *testing.T) {
	RegisterTestingT(t)

	snapshot := &snapshotdata.HeroSnapshot{
		Id:               1,
		Name:             "hero",
		INTERNAL_VERSION: 2,
	}

	heroProto := snapshot.EncodeClient()
	Ω(heroProto).ShouldNot(BeNil())
	Ω(snapshot.EncodeClientBytes()).Should(Equal(util.SafeMarshal(heroProto)))

	heroBasicProtoBytes := snapshot.EncodeBasic4ClientBytes()
	Ω(heroBasicProtoBytes).ShouldNot(BeEmpty())
	Ω(util.SafeMarshal(snapshot.EncodeBasic4Client())).Should(Equal(heroBasicProtoBytes))
}

func TestHeroSnapshotService_Cache(t *testing.T) {
	RegisterTestingT(t)

	service := newHeroSnapshotService(&panicDb{t: t}, ifacemock.ConfigDatas, ifacemock.GuildSnapshotService, ifacemock.BaiZhanService, 2)

	snapshot := &snapshotdata.HeroSnapshot{
		Id:               1,
		INTERNAL_VERSION: 2,
	}

	service.Update(snapshot)

	Ω(service.Get(1)).Should(BeIdenticalTo(snapshot))

	oldVersion := &snapshotdata.HeroSnapshot{
		Id:               1,
		INTERNAL_VERSION: 1,
	}
	service.Update(oldVersion)

	Ω(service.Get(1)).Should(BeIdenticalTo(snapshot))

	s2 := &snapshotdata.HeroSnapshot{
		Id:               2,
		INTERNAL_VERSION: 2,
	}
	service.Update(s2)

	s3 := &snapshotdata.HeroSnapshot{
		Id:               3,
		INTERNAL_VERSION: 2,
	}
	service.Update(s3)

	Ω(service.getFromCache(1)).Should(BeNil())
}

func TestHeroSnapshotService_Online(t *testing.T) {
	RegisterTestingT(t)

	service := newHeroSnapshotService(&panicDb{t: t}, ifacemock.ConfigDatas, ifacemock.GuildSnapshotService, ifacemock.BaiZhanService, 2)

	snapshot := &snapshotdata.HeroSnapshot{
		Id:               1,
		INTERNAL_VERSION: 2,
	}

	service.Online(snapshot)

	Ω(service.getFromCache(1)).Should(BeIdenticalTo(snapshot))
	for i := 2; i < 5; i++ {
		snapshot := &snapshotdata.HeroSnapshot{Id: int64(i), INTERNAL_VERSION: 1}
		service.Update(snapshot)
		Ω(service.getFromCache(int64(i))).Should(BeIdenticalTo(snapshot))
	}

	Ω(service.getFromCache(1)).Should(BeIdenticalTo(snapshot))
	Ω(service.getFromCache(2)).Should(BeNil())
	Ω(service.getFromCache(3)).ShouldNot(BeNil())

	service.Offline(1)

	for i := 5; i < 8; i++ {
		snapshot := &snapshotdata.HeroSnapshot{Id: int64(i), INTERNAL_VERSION: 1}
		service.Update(snapshot)
		Ω(service.getFromCache(int64(i))).Should(BeIdenticalTo(snapshot))
	}
	Ω(service.getFromCache(1)).Should(BeNil())
}

type recordCallback struct {
	sync.Mutex
	callback []*snapshotdata.HeroSnapshot
}

func (s *recordCallback) OnHeroSnapshotUpdate(snapshot *snapshotdata.HeroSnapshot) {
	s.Lock()
	s.callback = append(s.callback, snapshot)
	s.Unlock()
}

func (s *recordCallback) copy() []*snapshotdata.HeroSnapshot {
	s.Lock()
	array := make([]*snapshotdata.HeroSnapshot, len(s.callback))
	copy(array, s.callback)
	s.Unlock()

	return array
}

func TestHeroSnapshotService_Callback(t *testing.T) {
	RegisterTestingT(t)

	service := newHeroSnapshotService(&panicDb{t: t}, ifacemock.ConfigDatas, ifacemock.GuildSnapshotService, ifacemock.BaiZhanService, 2)
	callback := &recordCallback{}
	service.RegisterCallback(callback)

	snapshot := &snapshotdata.HeroSnapshot{
		Id:               1,
		INTERNAL_VERSION: 2,
	}

	service.Online(snapshot)

	Ω(callback.copy()).Should(HaveLen(0))
	service.Get(1)
	service.Offline(1)
	service.Online(snapshot)
	Ω(callback.copy()).Should(HaveLen(0))

	service.Update(snapshot)
	Eventually(func() []*snapshotdata.HeroSnapshot {
		return callback.copy()
	}).Should(HaveLen(1))
	Ω(callback.copy()[0]).Should(BeIdenticalTo(snapshot))

	s2 := &snapshotdata.HeroSnapshot{
		Id:               1,
		INTERNAL_VERSION: 3,
	}
	service.Update(s2)
	Eventually(func() []*snapshotdata.HeroSnapshot {
		return callback.copy()
	}).Should(HaveLen(2))
	Ω(callback.copy()[1]).Should(BeIdenticalTo(s2))

	// if update old snapshot, should not trigger callback
	s3 := &snapshotdata.HeroSnapshot{
		Id:               1,
		INTERNAL_VERSION: 2,
	}
	service.Update(s3)
	Eventually(func() []*snapshotdata.HeroSnapshot {
		return callback.copy()
	}).Should(HaveLen(2))

	service.Offline(1)
	service.Update(s3)
	Eventually(func() []*snapshotdata.HeroSnapshot {
		return callback.copy()
	}).Should(HaveLen(2))

}

type notExistDB struct {
	iface.DbService
	getCount int
}

func (d *notExistDB) LoadHero(ctx context.Context, id int64) (*entity.Hero, error) {
	d.getCount++
	return nil, nil
}

func TestHeroSnapshotService_GetNotExist(t *testing.T) {
	RegisterTestingT(t)

	db := &notExistDB{}
	service := newHeroSnapshotService(db, ifacemock.ConfigDatas, ifacemock.GuildSnapshotService, ifacemock.BaiZhanService, 2)
	callback := &recordCallback{}
	service.RegisterCallback(callback)

	Ω(service.Get(1)).Should(BeNil())
	Ω(db.getCount).Should(BeNumerically("==", 1))

	Ω(service.getFromCache(1)).Should(BeNil())
	Ω(service.Get(1)).Should(BeNil())
	Ω(db.getCount).Should(BeNumerically("==", 1)) // nil result should be cached

	snapshot := &snapshotdata.HeroSnapshot{
		Id:               1,
		INTERNAL_VERSION: 1,
	}

	service.Update(snapshot)
	Ω(service.getFromCache(1)).Should(BeIdenticalTo(snapshot))
	Ω(service.Get(1)).Should(BeIdenticalTo(snapshot))
}

func TestEncode(t *testing.T) {
	RegisterTestingT(t)

	basicProto := &shared_proto.HeroBasicProto{
		Id:   idbytes.ToBytes(1),
		Name: "hero",
	}

	proto := &shared_proto.HeroBasicSnapshotProto{
		Basic: basicProto,
	}

	snapshot := &snapshotdata.HeroSnapshot{
		Id:               1,
		IdBytes:          idbytes.ToBytes(1),
		Name:             "hero",
		INTERNAL_VERSION: 2,
		GuildId:          1,
	}

	Ω(snapshot.EncodeClientBytes()).Should(BeEquivalentTo(must.Marshal(proto)))
	Ω(snapshot.EncodeClient()).Should(BeEquivalentTo(proto))
	Ω(snapshot.EncodeBasic4ClientBytes()).Should(BeEquivalentTo(must.Marshal(basicProto)))
	Ω(snapshot.EncodeBasic4Client()).Should(BeEquivalentTo(basicProto))

	snapshot.ClearProto()
	Ω(snapshot.EncodeClientBytes()).Should(BeEquivalentTo(must.Marshal(proto)))
	Ω(snapshot.EncodeClient()).Should(BeEquivalentTo(proto))
	Ω(snapshot.EncodeBasic4ClientBytes()).Should(BeEquivalentTo(must.Marshal(basicProto)))
	Ω(snapshot.EncodeBasic4Client()).Should(BeEquivalentTo(basicProto))
}
