package entity

import (
	"github.com/lightpaw/logrus"
	"github.com/lightpaw/male7/config"
	"github.com/lightpaw/male7/config/goods"
	"github.com/lightpaw/male7/config/heroinit"
	"github.com/lightpaw/male7/pb/server_proto"
	"github.com/lightpaw/male7/pb/shared_proto"
	"github.com/lightpaw/male7/util/atomic"
	"github.com/lightpaw/male7/util/i64"
	"github.com/lightpaw/male7/util/imath"
	"github.com/lightpaw/male7/util/timeutil"
	"github.com/lightpaw/male7/util/u64"
	"math"
	"time"
	"github.com/lightpaw/male7/config/resdata"
	"github.com/lightpaw/male7/constants"
	"github.com/lightpaw/male7/util/collection"
)

func newDepot(initData *heroinit.HeroInitData) *Depot {
	d := &Depot{}
	d.goodsMap = make(map[uint64]uint64)
	d.baowuMap = make(map[uint64]uint64)

	d.baowuLogs = collection.NewRingList(initData.BaowuLogLimit)

	d.idGen = atomic.NewUint64(0)
	d.genIdGoodsMap = make(map[uint64]goods.GenIdGoods)
	d.tempGenIdGoodsMap = make(map[uint64]int64)

	d.maxDepotGenIdGoodsCapacity = []uint64{initData.MaxDepotEquipCapacity}

	d.tempDepotExpireDuration = initData.TempDepotExpireDuration

	d.genIdGoodsCountWithoutTempDepot = make([]uint64, len(d.maxDepotGenIdGoodsCapacity))
	d.genIdGoodsTotalCount = make([]uint64, len(d.maxDepotGenIdGoodsCapacity))

	d.nextCheckExpireGoodsTime = math.MaxInt64

	return d
}

// 背包
type Depot struct {
	// 普通物品
	goodsMap map[uint64]uint64 // key是id，value是个数

	// 宝物
	baowuMap map[uint64]uint64 // key是id，value是个数

	unlockBaowuData    *resdata.BaowuData
	unlockBaowuEndTime time.Time
	miaoBaowuTimes     uint64

	baowuLogs *collection.RingList

	// 一个格子一个的物品
	// 每次登陆重新生成id
	idGen             *atomic.Uint64
	genIdGoodsMap     map[uint64]goods.GenIdGoods
	tempGenIdGoodsMap map[uint64]int64 // key是生成的id, value是过期时间，临时背包

	maxDepotGenIdGoodsCapacity      []uint64 // 背包自增id的物品的最大容量
	genIdGoodsCountWithoutTempDepot []uint64 // 背包自增id的物品的数量(不包括临时背包)
	genIdGoodsTotalCount            []uint64 // 背包自增id的物品的数量(包括临时背包)

	tempDepotExpireDuration time.Duration // 临时背包物品的过期间隔

	nextCheckExpireGoodsTime      int64 // 下次检查物品过期的时间
	nextCheckMayExpiredGoodsCount int64 // 下次检查物品时可能过期的物品数量，只在removeGenIdGoods的时候判断那个物品的过期时间是不是跟检查时间一致，一致，就减少
}

func (d *Depot) GoodsMap() map[uint64]uint64 {
	return d.goodsMap
}

func (d *Depot) NewId() uint64 {
	return d.idGen.Inc()
}

func (d *Depot) unmarshal(heroId int64, heroName string, proto *server_proto.HeroDepotServerProto, datas interface {
	GetBaowuData(id uint64) *resdata.BaowuData
	GetEquipmentData(id uint64) *goods.EquipmentData
	EquipmentRefinedData() *config.EquipmentRefinedDataConfig
}, ctime time.Time) {
	if proto == nil {
		return
	}

	u64.CopyMapTo(d.goodsMap, proto.GetGoods())
	u64.CopyMapTo(d.baowuMap, proto.GetBaowu())

	if proto.UnlockBaowu > 0 {
		d.unlockBaowuData = datas.GetBaowuData(proto.UnlockBaowu)
		if d.unlockBaowuData != nil {
			d.unlockBaowuEndTime = timeutil.Unix64(proto.UnlockBaowuEndTime)
		}
	}
	d.miaoBaowuTimes = proto.MiaoBaowuTimes

	for _, v := range proto.BaowuLog {
		d.baowuLogs.Add(v)
	}

	timeUnix := timeutil.Marshal64(ctime)

	for _, p := range proto.Equipments {
		data := datas.GetEquipmentData(p.GetDataId())
		if data == nil {
			// 丢东西了
			logrus.WithField("lost", "equipment").
				WithField("heroId", heroId).
				WithField("heroName", heroName).
				WithField("id", p.GetDataId()).Errorf("lost depot equipment")
			continue
		}

		expiredTime, haveExpireTime := proto.TempExpiredTime[p.GetId()]
		if haveExpireTime && expiredTime <= timeUnix {
			// 已经过期，不加载上来了
			continue
		}

		id := d.NewId()
		e := NewEquipment(id, data)

		e.unmarshal(p, datas)

		if haveExpireTime {
			d.addGenIdGoodsCapacityNotEnough(e, expiredTime)
		} else {
			d.addGenIdGoodsCapacityEnough(e)
		}
	}
}

func (d *Depot) encodeServer() *server_proto.HeroDepotServerProto {

	proto := &server_proto.HeroDepotServerProto{}

	proto.Goods = u64.CopyMap(d.goodsMap)
	proto.Baowu = u64.CopyMap(d.baowuMap)

	if d.unlockBaowuData != nil {
		proto.UnlockBaowu = d.unlockBaowuData.Id
		proto.UnlockBaowuEndTime = timeutil.Marshal64(d.unlockBaowuEndTime)
	}
	proto.MiaoBaowuTimes = d.miaoBaowuTimes

	d.baowuLogs.Range(func(v interface{}) (toContinue bool) {
		if log, ok := v.([]byte); ok {
			proto.BaowuLog = append(proto.BaowuLog, log)
		}

		return true
	})

	// 装备
	proto.Equipments = make([]*server_proto.EquipmentServerProto, 0, d.genIdGoodsTotalCount[goods.EQUIPMENT])
	for _, g := range d.genIdGoodsMap {
		switch v := g.(type) {
		case *Equipment:
			proto.Equipments = append(proto.Equipments, v.encodeServer())
		default:
			logrus.WithField("goods", g).Errorln("存在新的自增id的物品没有存到Proto中去")
		}
	}

	proto.TempExpiredTime = u64.CopyUi64Map(d.tempGenIdGoodsMap)

	return proto
}

func (d *Depot) EncodeClient() *shared_proto.HeroDepotProto {

	proto := &shared_proto.HeroDepotProto{}

	for k, v := range d.goodsMap {
		proto.Goods = append(proto.Goods, &shared_proto.Int32Pair{u64.Int32(k), u64.Int32(v)})
	}

	for k, v := range d.baowuMap {
		proto.Baowu = append(proto.Baowu, &shared_proto.Int32Pair{u64.Int32(k), u64.Int32(v)})
	}

	if d.unlockBaowuData != nil {
		proto.UnlockBaowu = u64.Int32(d.unlockBaowuData.Id)
		proto.UnlockBaowuEndTime = timeutil.Marshal32(d.unlockBaowuEndTime)
	}
	proto.MiaoBaowuTimes = u64.Int32(d.miaoBaowuTimes)

	d.baowuLogs.ReverseRange(func(v interface{}) (toContinue bool) {
		if log, ok := v.([]byte); ok {
			proto.BaowuLog = append(proto.BaowuLog, log)
			if len(proto.BaowuLog) >= constants.BaowuLogPreviewCount {
				return false
			}
		}

		return true
	})

	// 装备
	proto.Equipments = make([]*shared_proto.EquipmentProto, 0, d.genIdGoodsTotalCount[goods.EQUIPMENT])
	for _, g := range d.genIdGoodsMap {
		switch v := g.(type) {
		case *Equipment:
			proto.Equipments = append(proto.Equipments, v.EncodeClient())
		default:
			logrus.WithField("goods", g).Errorln("存在新的自增id的物品没有存到Proto中去")
		}
	}

	for k, v := range d.tempGenIdGoodsMap {
		proto.TempExpiredTime = append(proto.TempExpiredTime, &shared_proto.Int32Pair{u64.Int32(k), i64.Int32(v)})
	}

	return proto
}

func (d *Depot) AddGoods(goodsId, toAdd uint64) uint64 {
	if toAdd > 0 {
		c := d.goodsMap[goodsId] + toAdd
		d.goodsMap[goodsId] = c
		return c
	}

	return d.goodsMap[goodsId]
}

func (d *Depot) RemoveGoods(goodsId, toRemove uint64) uint64 {
	c := u64.Sub(d.goodsMap[goodsId], toRemove)
	if c > 0 {
		d.goodsMap[goodsId] = c
	} else {
		delete(d.goodsMap, goodsId)
	}

	return c
}

func (d *Depot) HasEnoughGoodsArray(goodsId, count []uint64) bool {
	n := imath.Min(len(goodsId), len(count))
	for i := 0; i < n; i++ {
		if !d.HasEnoughGoods(goodsId[i], count[i]) {
			return false
		}
	}

	return true
}

func (d *Depot) HasEnoughGoods(goodsId, count uint64) bool {
	return d.goodsMap[goodsId] >= count
}

func (d *Depot) GetGoodsCount(goodsId uint64) uint64 {
	return d.goodsMap[goodsId]
}

func (d *Depot) AddBaowu(goodsId, toAdd uint64) uint64 {
	if toAdd > 0 {
		c := d.baowuMap[goodsId] + toAdd
		d.baowuMap[goodsId] = c
		return c
	}

	return d.baowuMap[goodsId]
}

func (d *Depot) RemoveBaowu(goodsId, toRemove uint64) uint64 {
	c := u64.Sub(d.baowuMap[goodsId], toRemove)
	if c > 0 {
		d.baowuMap[goodsId] = c
	} else {
		delete(d.baowuMap, goodsId)
	}

	return c
}

func (d *Depot) HasEnoughBaowuArray(goodsId, count []uint64) bool {
	n := imath.Min(len(goodsId), len(count))
	for i := 0; i < n; i++ {
		if !d.HasEnoughBaowu(goodsId[i], count[i]) {
			return false
		}
	}

	return true
}

func (d *Depot) HasEnoughBaowu(goodsId, count uint64) bool {
	return d.baowuMap[goodsId] >= count
}

func (d *Depot) GetBaowuCount(goodsId uint64) uint64 {
	return d.baowuMap[goodsId]
}

func (d *Depot) UnlockBaowu(data *resdata.BaowuData, ctime time.Time) {
	d.unlockBaowuData = data
	d.unlockBaowuEndTime = ctime.Add(data.UnlockDuration)
}

func (d *Depot) ClearUnlockBaowu() {
	d.unlockBaowuData = nil
	d.unlockBaowuEndTime = time.Time{}
}

func (d *Depot) IncMiaoBaowuTimes() {
	d.miaoBaowuTimes++
}

func (d *Depot) GetMiaoBaowuTimes() uint64 {
	return d.miaoBaowuTimes
}

func (d *Depot) GetUnlockBaowuData() *resdata.BaowuData {
	return d.unlockBaowuData
}

func (d *Depot) GetUnlockBaowuEndTime() time.Time {
	return d.unlockBaowuEndTime
}

func (d *Depot) RangeBaowu(f func(id, count uint64) (toContinue bool)) {
	for k, v := range d.baowuMap {
		if !f(k, v) {
			break
		}
	}
}

func (d *Depot) AddBaowuLog(log []byte) {
	d.baowuLogs.Add(log)
}

func (d *Depot) RangeBaowuLogWithStartIndex(startIndex int, f func(log []byte) (toContinue bool)) {
	d.baowuLogs.ReverseRangeWithStartIndex(imath.Max(startIndex, 0), func(v interface{}) (toContinue bool) {
		if log, ok := v.([]byte); ok {
			return f(log)
		}
		return true
	})
}

// 物品是否过期了
func (d *Depot) IsGenIdGoodsExpired(id uint64, ctime time.Time) (isExpired bool) {
	expireTime, ok := d.tempGenIdGoodsMap[id]
	if ok {
		isExpired = timeutil.Marshal64(ctime) >= expireTime
	}
	return
}

// 是否有足够的容量
func (d *Depot) HasEnoughGenIdGoodsCapacity(goodsType goods.GoodsType, addCount uint64) (haveEnough bool) {
	return addCount <= d.getGenIdGoodsCapacity(goodsType)
}

func (d *Depot) getGenIdGoodsCapacity(genIdGoodsType goods.GoodsType) (capacity uint64) {
	return u64.Sub(d.maxDepotGenIdGoodsCapacity[genIdGoodsType], d.genIdGoodsCountWithoutTempDepot[genIdGoodsType])
}

func (d *Depot) getTmpGoodsCount(genIdGoodsType goods.GoodsType) (haveCapacity uint64) {
	return u64.Sub(d.genIdGoodsTotalCount[genIdGoodsType], d.genIdGoodsCountWithoutTempDepot[genIdGoodsType])
}

func (d *Depot) AddGenIdGoods(g goods.GenIdGoods, ctime time.Time) (expireTime int64) {
	if d.HasEnoughGenIdGoodsCapacity(g.GoodsData().GoodsType(), 1) {
		d.addGenIdGoodsCapacityEnough(g)
	} else {
		// 不够，加到临时背包去
		expireTime = timeutil.Marshal64(ctime.Add(d.tempDepotExpireDuration))
		d.addGenIdGoodsCapacityNotEnough(g, expireTime)
	}

	return
}

// 背包容量足够
func (d *Depot) addGenIdGoodsCapacityEnough(g goods.GenIdGoods) {
	d.genIdGoodsMap[g.Id()] = g

	// 增加不带过期时间的物品的数量
	d.genIdGoodsCountWithoutTempDepot[g.GoodsData().GoodsType()]++
	// 增加所有物品的数量
	d.genIdGoodsTotalCount[g.GoodsData().GoodsType()]++
}

// 背包容量不足够
func (d *Depot) addGenIdGoodsCapacityNotEnough(g goods.GenIdGoods, expireTime int64) {
	d.genIdGoodsMap[g.Id()] = g
	// 不够，加到临时背包去
	d.tempGenIdGoodsMap[g.Id()] = expireTime

	// 增加所有物品的数量
	d.genIdGoodsTotalCount[g.GoodsData().GoodsType()]++

	if d.nextCheckExpireGoodsTime > expireTime {
		d.nextCheckExpireGoodsTime = expireTime
		d.nextCheckMayExpiredGoodsCount = 1
	} else if d.nextCheckExpireGoodsTime == expireTime {
		d.nextCheckMayExpiredGoodsCount++
	}
}

func (d *Depot) getGenIdGoods(id uint64) goods.GenIdGoods {
	return d.genIdGoodsMap[id]
}

func (d *Depot) WalkGenIdGoods(walkFunc func(goods goods.GenIdGoods)) {
	for _, goods := range d.genIdGoodsMap {
		walkFunc(goods)
	}
}

// 获得没有过期的装备
func (d *Depot) GetNotExpiredGenIdGoods(id uint64, ctime time.Time) (g goods.GenIdGoods, haveExpireTime bool) {
	g = d.getGenIdGoods(id)
	if g == nil {
		return nil, false
	}

	expireTime, ok := d.tempGenIdGoodsMap[g.Id()]
	if !ok {
		// 没过期时间
		return g, false
	}

	if timeutil.Marshal64(ctime) >= expireTime {
		// 过期了
		return nil, false
	}

	return g, true
}

func (d *Depot) RemoveGenIdGoods(toRemove uint64) {
	g := d.genIdGoodsMap[toRemove]
	if g == nil {
		return
	}

	delete(d.genIdGoodsMap, toRemove)

	// 减少总的数量
	d.genIdGoodsTotalCount[g.GoodsData().GoodsType()]--

	if expireTime, ok := d.tempGenIdGoodsMap[toRemove]; ok {
		delete(d.tempGenIdGoodsMap, toRemove)

		if expireTime == d.nextCheckExpireGoodsTime {
			d.nextCheckMayExpiredGoodsCount--
		}
	} else {
		// 减少不带过期时间的数量
		d.genIdGoodsCountWithoutTempDepot[g.GoodsData().GoodsType()]--
	}
}

// 删掉过期物品
func (d *Depot) RemoveExpiredGoods(ctime time.Time) (removeIds []uint64) {
	timeUnix := timeutil.Marshal64(ctime)

	if timeUnix < d.nextCheckExpireGoodsTime {
		return
	}

	removeIds = make([]uint64, 0, d.nextCheckMayExpiredGoodsCount)
	var nextCheckTime int64 = math.MaxInt64
	var nextMayExpireGoodsCount int64 = 0

	for id, expireTime := range d.tempGenIdGoodsMap {
		// 检查过期时间
		if timeUnix >= expireTime {
			// 过期了
			d.RemoveGenIdGoods(id)

			removeIds = append(removeIds, id)
		} else {
			if nextCheckTime > expireTime {
				nextCheckTime = expireTime
				nextMayExpireGoodsCount = 1
			} else if nextCheckTime == expireTime {
				nextMayExpireGoodsCount++
			}
		}
	}

	// 设置下次时间
	d.nextCheckExpireGoodsTime = nextCheckTime
	d.nextCheckMayExpiredGoodsCount = nextMayExpireGoodsCount

	return
}

// 在背包有空位时移动临时背包中的物品到背包中
func (d *Depot) MoveTmpGoodsToDepotIfDepotHaveSlot(goodsType goods.GoodsType, ctime time.Time) (moveToDepotIds []uint64) {
	if len(d.tempGenIdGoodsMap) <= 0 {
		return
	}

	// 临时背包还有物品
	capacity := d.getGenIdGoodsCapacity(goodsType)

	if capacity <= 0 {
		// 空间不够
		return
	}

	tmpGoodsCount := d.getTmpGoodsCount(goodsType)
	if tmpGoodsCount <= 0 {
		// 物品不够
		return
	}

	if capacity > tmpGoodsCount {
		capacity = tmpGoodsCount
	}

	timeUnix := timeutil.Marshal64(ctime)

	moveToDepotIds = make([]uint64, 0, capacity)
	moveToDepotExpireTime := make([]int64, 0, capacity)

	for id, expireTime := range d.tempGenIdGoodsMap {
		if timeUnix >= expireTime {
			// 过期了
			continue
		}

		g := d.getGenIdGoods(id)
		if g == nil {
			logrus.Errorln("我去，竟然有物品在临时背包，但是实际背包中又没有!%d", id)
			continue
		}

		if g.GoodsData().GoodsType() != goodsType {
			continue
		}

		if uint64(len(moveToDepotIds)) < capacity {
			// 还有空间
			moveToDepotIds = append(moveToDepotIds, id)
			moveToDepotExpireTime = append(moveToDepotExpireTime, expireTime)
			continue
		}

		replaceIdx := -1
		maxExpireTimeDiff := int64(0)

		for i := 0; i < len(moveToDepotIds); i++ {
			if moveToDepotExpireTime[i] <= expireTime {
				continue
			}

			if replaceIdx < 0 {
				replaceIdx = i
				maxExpireTimeDiff = moveToDepotExpireTime[i] - expireTime
			} else {
				diff := moveToDepotExpireTime[i] - expireTime
				if diff > maxExpireTimeDiff {
					// 那个比这个更晚过期
					replaceIdx = i
					maxExpireTimeDiff = diff
				}
			}
		}

		if replaceIdx >= 0 {
			moveToDepotExpireTime[replaceIdx] = expireTime
			moveToDepotIds[replaceIdx] = id
		}
	}

	for _, id := range moveToDepotIds {
		delete(d.tempGenIdGoodsMap, id)
		d.genIdGoodsCountWithoutTempDepot[goodsType]++
	}

	return moveToDepotIds
}

func (d *Depot) resetDaily(datas interface {
	GetGoodsData(u uint64) *goods.GoodsData
}) {
	d.miaoBaowuTimes = 0

	for id := range d.goodsMap {
		g := datas.GetGoodsData(id)
		if g == nil {
			continue
		}
		if g.SpecType == shared_proto.GoodsSpecType_GAT_HEBI {
			delete(d.goodsMap, id)
		}
	}
}
