package hebi

import (
	"testing"
	. "github.com/onsi/gomega"
	"github.com/lightpaw/male7/mock"
	"github.com/lightpaw/male7/gen/ifacemock"
	"github.com/lightpaw/male7/entity"
	"github.com/lightpaw/male7/gen/pb/hebi"
	"github.com/lightpaw/male7/pb/shared_proto"
	"github.com/lightpaw/male7/util/u64"
)

func TestHebiModule_ProcessChangeCaptain(t *testing.T) {
	RegisterTestingT(t)

	dep := mock.NewMockDep2()

	ctime := dep.Time().CurrentTime()
	var heroId int64 = 1
	hero := entity.NewHero(heroId, "hero1", dep.Datas().HeroInitData(), ctime)
	mock.DefaultHero(hero)

	snapshot := dep.HeroSnapshot().NewSnapshot(hero)
	dep.HeroSnapshot().Online(snapshot)

	military := hero.Military()
	for i := 0; i < 5; i++ {
		c := entity.NewTestCaptain(military.NewCaptainId(), dep.Datas().RaceData().MinKeyData, dep.Datas(), military.GetOrCreateSoldierData,
			hero.LevelData, nil)
		military.AddCaptain(c)
	}

	var goodsFound bool
	var goodsId int32
	for _, g := range dep.Datas().GoodsConfig().SpecGoods {
		if g.SpecType == shared_proto.GoodsSpecType_GAT_HEBI {
			hero.Depot().AddGoods(g.Id, 1)
			goodsFound = true
			goodsId = u64.Int32(g.Id)
			break
		}
	}
	if !goodsFound {
		return
	}

	m := NewHebiModule(dep, ifacemock.MailModule, ifacemock.FightService, mock.MockTick())
	var roomId int32 = 0

	c2sChange1 := &hebi.C2SChangeCaptainProto{1}
	m.ProcessChangeCaptain(c2sChange1, ifacemock.HeroController)

	c2sCheckIn := &hebi.C2SCheckInRoomProto{roomId, goodsId}
	m.ProcessCheckInRoom(c2sCheckIn, ifacemock.HeroController)

	r := m.hebiManager.GetRoom(u64.FromInt32(roomId))
	Ω(r.hostId).Should(Equal(heroId))
	Ω(r.hostCaptain.Id).Should(Equal(int32(1)))

	c2sChange2 := &hebi.C2SChangeCaptainProto{2}
	m.ProcessChangeCaptain(c2sChange2, ifacemock.HeroController)

	Ω(r.hostCaptain.Id).Should(Equal(int32(2)))
}
