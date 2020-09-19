package realm

import (
	"github.com/lightpaw/male7/entity/cb"
	"github.com/lightpaw/male7/util/atomic"
	"sync"

	"github.com/lightpaw/logrus"
	"github.com/lightpaw/male7/config/blockdata"
	"github.com/lightpaw/male7/util/u64"
	"github.com/lightpaw/male7/config/regdata"
	"github.com/lightpaw/male7/constants"
	"github.com/lightpaw/male7/gen/iface"
)

func newBlockManager(mapData *blockdata.StitchedBlocks, initRadius uint64) *block_manager {
	return &block_manager{
		mapData:       mapData,
		radius:        atomic.NewUint64(initRadius),
		blockMap:      &sync.Map{},
		blockIndexMap: &sync.Map{},
	}
}

type block_manager struct {
	mapData *blockdata.StitchedBlocks

	radius *atomic.Uint64

	//blockMap map[cb.Cube]*block
	blockMap *sync.Map

	blockIndexMap *sync.Map
}

func (m *block_manager) addBase(x, y int, isHome bool) {
	b := m.getOrCreateBlockByBasePos(x, y)
	b.baseCount.Inc()
	if isHome {
		b.homeCount.Inc()
	}
}

func (m *block_manager) removeBase(x, y int, isHome bool) {
	b := m.getOrCreateBlockByBasePos(x, y)
	b.baseCount.Dec()
	if isHome {
		b.homeCount.Dec()
	}
}

func (m *block_manager) moveBase(oldX, oldY, newX, newY int, isHome bool) {
	oldBlockX, oldBlockY := m.mapData.MustBlockByPos(oldX, oldY)
	newBlockX, newBlockY := m.mapData.MustBlockByPos(newX, newY)

	if oldBlockX != newBlockX || oldBlockY != newBlockY {
		b := m.getOrCreateBlock(cb.XYCube(int(oldBlockX), int(oldBlockY)))
		b.baseCount.Dec()
		if isHome {
			b.homeCount.Dec()
		}

		b = m.getOrCreateBlock(cb.XYCube(int(newBlockX), int(newBlockY)))
		b.baseCount.Inc()
		if isHome {
			b.homeCount.Inc()
		}
	}
}

func (m *block_manager) getBlock(c cb.Cube) *block {
	if b, ok := m.blockMap.Load(c); ok {
		return b.(*block)
	}
	return nil
}

func (m *block_manager) getOrCreateBlockByBasePos(x, y int) *block {
	blockX, blockY := m.mapData.MustBlockByPos(x, y)
	return m.getOrCreateBlock(cb.XYCube(int(blockX), int(blockY)))
}

func (m *block_manager) getOrCreateBlock(c cb.Cube) *block {

	b, ok := m.blockMap.Load(c)
	if !ok {
		x, y := c.XY()
		b, _ = m.blockMap.LoadOrStore(c, newBlock(uint64(x), uint64(y)))
	}

	return b.(*block)
}

func (m *block_manager) getLeastHomeCountBlock() (least *block) {

	var homeCount uint64
	m.rangeBlock(func(b *block) (toContinue bool) {
		count := b.GetHomeCount()
		if least == nil || homeCount > count {
			least = b
			homeCount = count
		}
		return true
	})

	return
}

func (m *block_manager) rangeBlock(f func(b *block) (toContinue bool)) {
	m.blockMap.Range(func(key, value interface{}) bool {
		return f(value.(*block))
	})
}

func newBlock(blockX uint64, blockY uint64) *block {
	return &block{
		blockX:     blockX,
		blockY:     blockY,
		baseMap:    make(map[int64]cb.Cube),
		baozMap:    make(map[uint64]uint64),
		junTuanMap: make(map[uint64]uint64),
	}
}

type block struct {
	blockX uint64
	blockY uint64

	baseCount atomic.Int64
	homeCount atomic.Int64

	// 这里也可以只存储id，从总的baseMap中取数据
	baseMap map[int64]cb.Cube // TODO

	// 宝藏怪物待刷新列表，刷新时候，根据个数，从最少的开始刷新
	baozMap map[uint64]uint64

	// 军团怪物待刷新列表，刷新时候，根据个数，从最少的开始刷新
	junTuanMap map[uint64]uint64
}

func (b *block) GetBaseCount() uint64 {
	return u64.FromInt64(b.baseCount.Load())
}

func (b *block) GetHomeCount() uint64 {
	return u64.FromInt64(b.homeCount.Load())
}

func (b *block) getNextToAddBaozNpc(baozDatas []*regdata.BaozNpcData) (*regdata.BaozNpcData, uint64) {

	// 找到个数最少的那个怪物
	var toAdd *regdata.BaozNpcData
	var minCount uint64
	for _, v := range baozDatas {
		count := b.baozMap[v.Id]
		if count < v.KeepCount {
			if toAdd == nil || count < minCount {
				toAdd = v
				minCount = count
			}
		}
	}

	return toAdd, minCount
}

func (b *block) getNextToAddJunTuanNpc(datas map[uint64]*regdata.JunTuanNpcPlaceData) (*regdata.JunTuanNpcPlaceData, uint64) {

	// 找到个数最少的那个怪物
	var toAdd *regdata.JunTuanNpcPlaceData
	var minCount uint64
	for _, v := range datas {
		count := b.junTuanMap[v.Group]
		if count < v.KeepCount {
			if toAdd == nil || count < minCount {
				toAdd = v
				minCount = count
			}
		}
	}

	return toAdd, minCount
}

func (b *block) getBaozCount(baozData *regdata.BaozNpcData) uint64 {
	return b.baozMap[baozData.Id]
}

func (b *block) increseBaozCount(baozData *regdata.BaozNpcData) {
	b.baozMap[baozData.Id]++
}

func (b *block) decreseBaozCount(baozData *regdata.BaozNpcData) {
	b.baozMap[baozData.Id] = u64.Sub(b.baozMap[baozData.Id], 1)
}

func (b *block) increseJunTuanCount(data *regdata.JunTuanNpcData) {
	b.junTuanMap[data.Group]++
}

func (b *block) decreseJunTuanCount(data *regdata.JunTuanNpcData) {
	b.junTuanMap[data.Group] = u64.Sub(b.junTuanMap[data.Group], 1)
}

//getOrCreateBlockIndex 获取或者创建索引
func (m *block_manager) getOrCreateBlockIndex(x, y int) *block_index {
	c := cb.XYCube(x/constants.RealmIndexBlockSize, y/constants.RealmIndexBlockSize)
	b, ok := m.blockIndexMap.Load(c)
	if !ok {
		b, _ = m.blockIndexMap.LoadOrStore(c, newBlockIndex())
	}

	return b.(*block_index)
}

//newBlockIndex 新建索引
func newBlockIndex() *block_index {
	return &block_index{
		herosMap:  &sync.Map{},
		troopsMap: &sync.Map{},
		basesMap:  &sync.Map{},
		ruinsMap:  &sync.Map{},
	}
}

//索引结构
type block_index struct {
	herosMap  *sync.Map
	troopsMap *sync.Map
	basesMap  *sync.Map
	ruinsMap  *sync.Map
}

//AddHeroIndex 添加英雄索引
func (m *block_manager) AddHeroIndex(x, y int, hc iface.HeroController) {
	//存储英雄
	b := m.getOrCreateBlockIndex(x, y)
	if b != nil {
		b.herosMap.Store(hc.Id(), hc)
		//添加反向索引
		temp := hc.SetBlockIndex(b)
		if btemp, ok := temp.(*block_index); ok && btemp != b {
			btemp.RemoveHeroIndex(hc)
		}
		logrus.Debugf("添加英雄索引 坐标[%d] [%d]", x, y)
	}
}

//RemoveHeroIndex 移除英雄索引
func (m *block_manager) RemoveHeroIndex(hc iface.HeroController) {
	temp := hc.GetBlockIndex()
	if btemp, ok := temp.(*block_index); ok {
		btemp.RemoveHeroIndex(hc)
		logrus.Debugf("移除英雄索引")
	}
}

func (b *block_index) RemoveHeroIndex(hc iface.HeroController) {
	//删除英雄
	b.herosMap.Delete(hc.Id())
}

//rangeHero 遍历英雄
func (m *block_manager) rangeHero(x, y int, f func(hc iface.HeroController) (toContinue bool)) {
	b := m.getOrCreateBlockIndex(x, y)
	if b != nil {
		b.herosMap.Range(func(key, value interface{}) bool {
			return f(value.(iface.HeroController))
		})
	}
}

//AddTroopIndex 添加军队索引
func (m *block_manager) AddTroopIndex(x, y int, t *troop) {
	//存储军队
	b := m.getOrCreateBlockIndex(x, y)
	if b != nil {
		b.troopsMap.Store(t.Id(), t)
		//添加反向索引
		t.realmTracks = append(t.realmTracks, b)
	}
}

//移除军队索引
func (m *block_manager) RemoveTroopIndex(x, y int, t *troop) {
	//移除军队索引
	b := m.getOrCreateBlockIndex(x, y)
	if b != nil {
		b.troopsMap.Delete(t.Id())
	}
}

//移除军队索引
func (b *block_index) RemoveTroopIndex(t *troop) {
	b.troopsMap.Delete(t.Id())
}

//rangeTroop 遍历军队
func (m *block_manager) rangeTroop(x, y int, f func(t *troop) (toContinue bool)) {
	b := m.getOrCreateBlockIndex(x, y)
	if b != nil {
		b.troopsMap.Range(func(key, value interface{}) bool {
			return f(value.(*troop))
		})
	}
}

//AddBaseIndex 添加基地索引
func (m *block_manager) AddBaseIndex(x, y int, base *baseWithData) {
	b := m.getOrCreateBlockIndex(x, y)
	if b != nil {
		b.basesMap.Store(base.Id(), base)
	}
}

//移除基地索引
func (m *block_manager) RemoveBaseIndex(x, y int, base *baseWithData) {
	b := m.getOrCreateBlockIndex(x, y)
	if b != nil {
		b.basesMap.Delete(base.Id())
	}
}

//rangeBase 遍历基地
func (m *block_manager) rangeBase(x, y int, f func(base *baseWithData) (toContinue bool)) {
	b := m.getOrCreateBlockIndex(x, y)
	if b != nil {
		b.basesMap.Range(func(key, value interface{}) bool {
			return f(value.(*baseWithData))
		})
	}
}

//AddRuinIndex 添加墓碑索引
func (m *block_manager) AddRuinIndex(x, y int, ruin *ruinsBase) {
	b := m.getOrCreateBlockIndex(x, y)
	if b != nil {
		b.ruinsMap.Store(ruin, 1)
	}
}

//移除墓碑索引
func (m *block_manager) RemoveRuinIndex(x, y int, ruin *ruinsBase) {
	b := m.getOrCreateBlockIndex(x, y)
	if b != nil {
		b.ruinsMap.Delete(ruin)
	}
}

//移除墓碑索引
func (b *block_index) RemoveRuinIndex(ruin *ruinsBase) {
	b.ruinsMap.Delete(ruin)
}

//rangeRuin 遍历墓地
func (m *block_manager) rangeRuin(x, y int, f func(ruin *ruinsBase) (toContinue bool)) {
	b := m.getOrCreateBlockIndex(x, y)
	if b != nil {
		b.ruinsMap.Range(func(key, value interface{}) bool {
			return f(key.(*ruinsBase))
		})
	}
}
