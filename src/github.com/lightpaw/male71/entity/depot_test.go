package entity

import (
	"github.com/lightpaw/male7/config"
	"github.com/lightpaw/male7/config/data"
	"github.com/lightpaw/male7/config/goods"
	"github.com/lightpaw/male7/config/heroinit"
	"github.com/lightpaw/male7/config/resdata"
	"github.com/lightpaw/male7/pb/server_proto"
	"github.com/lightpaw/male7/pb/shared_proto"
	"github.com/lightpaw/male7/util/timeutil"
	"github.com/lightpaw/male7/util/u64"
	. "github.com/onsi/gomega"
	"math"
	"math/rand"
	"reflect"
	"sort"
	"testing"
	"time"
)

var goods1 *goods.GoodsData = &goods.GoodsData{Id: 1}
var goods2 *goods.GoodsData = &goods.GoodsData{Id: 2}

var equip1 *goods.EquipmentData = &goods.EquipmentData{Id: 1, Quality: &goods.EquipmentQualityData{Level: 1, LevelDatas: []*goods.EquipmentQualityLevelData{&goods.EquipmentQualityLevelData{Level: 1}}}, BaseStat: &data.SpriteStat{}, BaseStatProto: &shared_proto.SpriteStatProto{}}
var equip2 *goods.EquipmentData = &goods.EquipmentData{Id: 2, Quality: &goods.EquipmentQualityData{Level: 1, LevelDatas: []*goods.EquipmentQualityLevelData{&goods.EquipmentQualityLevelData{Level: 1}}}, BaseStat: &data.SpriteStat{}, BaseStatProto: &shared_proto.SpriteStatProto{}}

var captain1 *resdata.ResCaptainData = &resdata.ResCaptainData{Id: 1}
var captain2 *resdata.ResCaptainData = &resdata.ResCaptainData{Id: 2}

var initData = &heroinit.HeroInitData{
	MaxDepotEquipCapacity:   100,
	TempDepotExpireDuration: 24 * time.Hour,
}

func TestDepot_Goods(t *testing.T) {
	RegisterTestingT(t)

	depot := newDepot(initData)

	Expect(depot.GetGoodsCount(goods1.Id)).Should(BeEquivalentTo(0))
	Expect(depot.GetGoodsCount(goods2.Id)).Should(BeEquivalentTo(0))

	depot.AddGoods(goods1.Id, 3)
	Expect(depot.GetGoodsCount(goods1.Id)).Should(BeEquivalentTo(3))

	depot.AddGoods(goods2.Id, 13)
	Expect(depot.GetGoodsCount(goods2.Id)).Should(BeEquivalentTo(13))

	depot.RemoveGoods(goods1.Id, 2)
	Expect(depot.GetGoodsCount(goods1.Id)).Should(BeEquivalentTo(1))

	depot.RemoveGoods(goods2.Id, 11)
	Expect(depot.GetGoodsCount(goods2.Id)).Should(BeEquivalentTo(2))

	// 变成10个
	depot.AddGoods(goods1.Id, 9)
	// 变成20个
	depot.AddGoods(goods2.Id, 18)
	Expect(depot.HasEnoughGoods(goods1.Id, 7)).Should(BeTrue())
	Expect(depot.HasEnoughGoods(goods2.Id, 13)).Should(BeTrue())
	Expect(depot.HasEnoughGoods(goods1.Id, 17)).Should(BeFalse())
	Expect(depot.HasEnoughGoods(goods2.Id, 23)).Should(BeFalse())

	Expect(depot.HasEnoughGoodsArray([]uint64{goods1.Id, goods2.Id}, []uint64{3, 18})).Should(BeTrue())
	Expect(depot.HasEnoughGoodsArray([]uint64{goods1.Id, goods2.Id}, []uint64{13, 18})).Should(BeFalse())

	Depot_EncodeDecodeTest(depot, time.Now())
}

func TestDepot_GenIdGoods(t *testing.T) {
	RegisterTestingT(t)

	depot := newDepot(initData)

	ctime := time.Now()

	var toRemoveEquipAndGemIds []uint64

	// 填满装备背包
	for i := uint64(0); i < depot.maxDepotGenIdGoodsCapacity[goods.EQUIPMENT]; i++ {
		Expect(depot.HasEnoughGenIdGoodsCapacity(goods.EQUIPMENT, 1)).Should(BeTrue())

		var expireTime int64

		newId := depot.NewId()

		if i%2 == 0 {
			expireTime = depot.AddGenIdGoods(newEqiupment(newId, equip1), ctime)
		} else {
			expireTime = depot.AddGenIdGoods(newEqiupment(newId, equip2), ctime)
		}

		Expect(expireTime).Should(BeEquivalentTo(0))
		Expect(depot.getTmpGoodsCount(goods.EQUIPMENT)).Should(BeEquivalentTo(0))
		Expect(depot.getGenIdGoodsCapacity(goods.EQUIPMENT)).Should(BeEquivalentTo(depot.maxDepotGenIdGoodsCapacity[goods.EQUIPMENT] - i - 1))

		if i < 3 {
			toRemoveEquipAndGemIds = append(toRemoveEquipAndGemIds, newId)
		}
	}
	Expect(depot.HasEnoughGenIdGoodsCapacity(goods.EQUIPMENT, 1)).Should(BeFalse())
	Expect(depot.nextCheckMayExpiredGoodsCount).Should(BeEquivalentTo(0))
	Expect(depot.nextCheckExpireGoodsTime).Should(BeEquivalentTo(math.MaxInt64))

	// 加过期装备
	for i := uint64(0); i < depot.maxDepotGenIdGoodsCapacity[goods.EQUIPMENT]; i++ {
		Expect(depot.HasEnoughGenIdGoodsCapacity(goods.EQUIPMENT, 1)).Should(BeFalse())

		addTime := ctime.Add(time.Second * time.Duration(rand.Int63n(50)))

		expireTime := timeutil.Marshal64(addTime.Add(depot.tempDepotExpireDuration))
		mayExpireGoodsCount := int64(0)

		if expireTime > depot.nextCheckExpireGoodsTime {
			// 新加的比那个要大
			mayExpireGoodsCount = depot.nextCheckMayExpiredGoodsCount
		} else if expireTime == depot.nextCheckExpireGoodsTime {
			// 新加的跟那个相等
			mayExpireGoodsCount = depot.nextCheckMayExpiredGoodsCount + 1
		} else {
			// 新加的跟那个小
			mayExpireGoodsCount = 1
		}

		newId := depot.NewId()
		realExpireTime := depot.AddGenIdGoods(newEqiupment(newId, equip1), addTime)

		Expect(expireTime).Should(BeEquivalentTo(realExpireTime))
		Expect(depot.nextCheckMayExpiredGoodsCount).Should(BeEquivalentTo(mayExpireGoodsCount))

		Expect(depot.getTmpGoodsCount(goods.EQUIPMENT)).Should(BeEquivalentTo(i + 1))
		Expect(depot.getGenIdGoodsCapacity(goods.EQUIPMENT)).Should(BeEquivalentTo(0))

		if i < 3 {
			toRemoveEquipAndGemIds = append(toRemoveEquipAndGemIds, newId)
		}
	}

	Depot_EncodeDecodeTest(depot, ctime)

	// 看30秒内应该有多少装备过期
	checkExpiredTime := ctime.Add(depot.tempDepotExpireDuration).Add(time.Second * 30)
	checkExpiredTimeUnix := timeutil.Marshal64(checkExpiredTime)

	var expireIds []uint64

	tmpEquipCount := depot.getTmpGoodsCount(goods.EQUIPMENT)
	expireEquipCount := uint64(0)

	for id, expireTime := range depot.tempGenIdGoodsMap {
		if expireTime > checkExpiredTimeUnix {
			// 不过期
			continue
		}

		g := depot.getGenIdGoods(id)
		Expect(g).Should(Not(BeNil()))

		expireIds = append(expireIds, g.Id())

		if g.GoodsData().GoodsType() == goods.EQUIPMENT {
			expireEquipCount++
		}
	}

	removeIds := depot.RemoveExpiredGoods(checkExpiredTime)

	u64.Sort(expireIds)
	u64.Sort(removeIds)

	Expect(depot.nextCheckExpireGoodsTime > checkExpiredTimeUnix).Should(BeTrue())

	Expect(len(removeIds)).Should(BeEquivalentTo(len(expireIds)))
	Expect(reflect.DeepEqual(removeIds, expireIds)).Should(BeTrue())

	Expect(tmpEquipCount - expireEquipCount).Should(BeEquivalentTo(depot.getTmpGoodsCount(goods.EQUIPMENT)))

	// 再来一次
	removeIds = depot.RemoveExpiredGoods(checkExpiredTime)
	// 都被移走了，长度为空
	Expect(len(removeIds)).Should(BeZero())

	equipCount := 0
	for id, expireTime := range depot.tempGenIdGoodsMap {
		if expireTime != depot.nextCheckExpireGoodsTime {
			continue
		}

		g := depot.getGenIdGoods(id)
		Expect(g).Should(Not(BeNil()))

		if g.GoodsData().GoodsType() == goods.EQUIPMENT {
			equipCount++
		}
	}

	oldTempGenIdGoodsMap := u64.CopyUi64Map(depot.tempGenIdGoodsMap)

	for _, id := range toRemoveEquipAndGemIds {
		depot.RemoveGenIdGoods(id)
		Expect(oldTempGenIdGoodsMap[id]).Should(BeZero())
	}

	moveEquipToDepotIds := depot.MoveTmpGoodsToDepotIfDepotHaveSlot(goods.EQUIPMENT, ctime)

	for _, id := range moveEquipToDepotIds {
		// 验证删除的都是临时背包里面的东西
		_, ok := oldTempGenIdGoodsMap[id]
		Expect(ok).Should(BeTrue())
		Expect(depot.tempGenIdGoodsMap[id]).Should(BeZero())
	}

	clearNoExpireTimeGoods(depot)

	moveEquipToDepotIds = depot.MoveTmpGoodsToDepotIfDepotHaveSlot(goods.EQUIPMENT, ctime)

	for _, id := range moveEquipToDepotIds {
		// 验证删除的都是临时背包里面的东西
		_, ok := oldTempGenIdGoodsMap[id]
		Expect(ok).Should(BeTrue())
		Expect(depot.tempGenIdGoodsMap[id]).Should(BeZero())
	}
}

func Depot_EncodeDecodeTest(depot *Depot, ctime time.Time) {
	proto := depot.encodeServer()

	// 必须要排下序，不然id生成就不对了
	sort.Sort(equipArray(proto.Equipments))

	newDepot := newDepot(initData)
	newDepot.unmarshal(0, "", proto, &configDatas{}, ctime)

	Expect(reflect.DeepEqual(depot, newDepot)).Should(BeTrue())
}

type equipArray []*server_proto.EquipmentServerProto

func (e equipArray) Len() int {
	return len(e)
}

func (e equipArray) Less(i, j int) bool {
	return e[i].Id < e[j].Id
}

func (e equipArray) Swap(i, j int) {
	e[i], e[j] = e[j], e[i]
}

func clearNoExpireTimeGoods(depot *Depot) {
	for id := range depot.genIdGoodsMap {
		_, ok := depot.tempGenIdGoodsMap[id]
		if ok {
			// 有过期时间，不删除
			continue
		}

		depot.RemoveGenIdGoods(id)
	}

}

func newEqiupment(id uint64, data *goods.EquipmentData) *Equipment {
	return NewEquipment(id, data)
}

type configDatas struct {
}

func (c *configDatas) GetGoodsData(id uint64) *goods.GoodsData {
	if id == goods1.Id {
		return goods1
	} else if id == goods2.Id {
		return goods2
	} else {
		return nil
	}
}

func (c *configDatas) GetEquipmentData(id uint64) *goods.EquipmentData {
	if id == equip1.Id {
		return equip1
	} else if id == equip2.Id {
		return equip2
	} else {
		return nil
	}
}

func (c *configDatas) GetCaptainSoulData(id uint64) *resdata.ResCaptainData {
	if id == captain1.Id {
		return captain1
	} else if id == captain2.Id {
		return captain2
	} else {
		return nil
	}
}

func (c *configDatas) EquipmentRefinedData() *config.EquipmentRefinedDataConfig {
	return nil
}

func (c *configDatas) GetBaowuData(id uint64) *resdata.BaowuData {
	return nil
}
