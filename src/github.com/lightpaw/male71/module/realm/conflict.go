package realm

import (
	"github.com/lightpaw/logrus"
	"github.com/lightpaw/male7/config/blockdata"
	"github.com/lightpaw/male7/constants"
	"github.com/lightpaw/male7/entity/cb"
	"github.com/lightpaw/male7/entity/hexagon"
	"runtime/debug"
	"sync"
	"github.com/lightpaw/male7/util/u64"
	"math/rand"
	"github.com/lightpaw/male7/util/imath"
	"github.com/lightpaw/male7/config/regdata"
)

type conflict struct {
	regionData *regdata.RegionData

	mapData *blockdata.StitchedBlocks

	baseConflictCount map[cb.Cube]int // 每个地图块, 冲突的主城个数

	sync.Mutex
}

func newConfllict(regionData *regdata.RegionData) *conflict {
	return &conflict{
		regionData:        regionData,
		mapData:           regionData.Block,
		baseConflictCount: make(map[cb.Cube]int),
	}
}

// 是否可以在本地图迁城
// 需要判断占着的那个位置是不是你自己占的
// must be called under lock
func (c *conflict) doCanMoveBase(baseX, baseY, newX, newY int) (canMove bool) {

	newCube := cb.XYCube(newX, newY)
	if n, has := c.baseConflictCount[newCube]; has {
		// 看下是不是自己占着这个位置
		if n != 1 {
			return false
		}

		basePos := NeighborsInBaseConflictRange(baseX, baseY)
		if !cb.Contains(basePos, newCube) {
			return false
		}

	}

	return true
}

func (c *conflict) moveBaseIfCanMove(baseX, baseY, newX, newY int) (canMove bool) {

	if baseX == newX && baseY == newY {
		return true
	}

	c.Lock()
	defer c.Unlock()
	if !c.doCanMoveBase(baseX, baseY, newX, newY) {
		return false
	}

	c.doRemoveBase(baseX, baseY)
	c.doAddBase(newX, newY)
	return true
}

func (c *conflict) addBaseIfCanAdd(x, y int) (canAdd bool) {
	c.Lock()
	defer c.Unlock()
	return c.addBaseIfCanAddUnderLocker(x, y)
}

func (c *conflict) addBaseAnyway(x, y int) {
	c.Lock()
	defer c.Unlock()
	c.doAddBase(x, y)
}

func (c *conflict) addBaseIfCanAddUnderLocker(x, y int) (canAdd bool) {
	if _, has := c.baseConflictCount[cb.XYCube(x, y)]; has {
		return false
	}

	c.doAddBase(x, y)
	return true
}

func (c *conflict) moveBaseIfCanAdd(oldX, oldY, newX, newY int) (canAdd bool) {
	c.Lock()
	defer c.Unlock()

	if !c.doCanMoveBase(oldX, oldY, newX, newY) {
		return false
	}

	c.doAddBase(newX, newY)
	return true
}

// must be called under lock
func (c *conflict) doAddBase(x, y int) {
	pos := NeighborsInBaseConflictRange(x, y)

	for _, p := range pos {
		if v, has := c.baseConflictCount[p]; has {
			c.baseConflictCount[p] = v + 1
		} else {
			c.baseConflictCount[p] = 1
		}
	}
}

func (c *conflict) removeBase(x, y int) {
	c.Lock()
	defer c.Unlock()

	c.doRemoveBase(x, y)
}

// must be called under lock
func (c *conflict) doRemoveBase(x, y int) {
	pos := NeighborsInBaseConflictRange(x, y)

	for _, p := range pos {
		if v, has := c.baseConflictCount[p]; has {
			switch v {
			case 0:
				logrus.WithField("stack", string(debug.Stack())).Error("realm.conflict.removeBase时, count == 0")
				delete(c.baseConflictCount, p)
			case 1:
				delete(c.baseConflictCount, p)
			default:
				c.baseConflictCount[p] = v - 1
			}
		} else {
			x1, y1 := p.XY()
			logrus.WithField("x", x).WithField("y", y).
				WithField("x1", x1).WithField("y1", y1).
				WithField("stack", string(debug.Stack())).Error("realm.conflict.removeBase时, 竟然坐标在conflictCount中不存在")
		}
	}
}

func (r *Realm) randomNewBasePos(country uint64) (ok bool, x int, y int) {

	cbb := r.regionData.GetCountryBornBlock(country)

	radius := r.GetRadius()

	r.conflict.Lock()
	defer r.conflict.Unlock()

	// 中心点螺旋往外找
	cbb.Range(func(c cb.Cube) (toContinue bool) {
		ix, iy := c.XY()
		blockX, blockY := uint64(ix), uint64(iy)

		b := r.blockManager.getOrCreateBlock(cb.XYCube(int(blockX), int(blockY)))
		if b.GetHomeCount() >= r.levelData.Block.BlockData().NewHeroCrowdedCapcity {
			return true // skip
		}

		// 每个block尝试3次
		for i := 0; i < 3; i++ {
			rx, ry := r.mapData.RandomHomeXY(blockX, blockY)
			if !r.mapData.IsValidHomePosition(rx, ry) {
				return true
			}

			if r.isEdgeNotHomePos(rx, ry, radius) {
				continue
			}

			if r.conflict.addBaseIfCanAddUnderLocker(rx, ry) {
				ok = true
				x, y = rx, ry
				return false
			}
		}

		return true
	})

	return
}

func (r *Realm) randomRebornBasePos() (ok bool, x int, y int) {

	//radius := r.GetRadius()
	//
	//r.conflict.Lock()
	//defer r.conflict.Unlock()
	//
	//// 找拥挤度最低的城池
	//if least := r.blockManager.getLeastHomeCountBlock(); least != nil {
	//	for i := 0; i < 3; i++ {
	//		// 随机3次
	//		rx, ry := r.mapData.RandomHomeXY(least.blockX, least.blockY)
	//
	//		if !r.mapData.IsValidHomePosition(rx, ry) {
	//			continue
	//		}
	//
	//		if r.isEdgeNotHomePos(rx, ry, radius) {
	//			continue
	//		}
	//
	//		if r.conflict.addBaseIfCanAddUnderLocker(rx, ry) {
	//			ok = true
	//			x, y = rx, ry
	//			return
	//		}
	//	}
	//}
	//
	// return

	// 改成随机规则
	return r.randomBasePos()
}

func (r *Realm) GetRadius() uint64 {
	return r.blockManager.radius.Load()
}

func (r *Realm) IsEdgeNotHomePos(x, y int) bool {
	return r.isEdgeNotHomePos(x, y, r.GetRadius())
}

func (r *Realm) isEdgeNotHomePos(x, y int, radius uint64) bool {

	ux, uy := u64.FromInt(x), u64.FromInt(y)
	if ux < r.config().EdgeNotHomeLen || uy < r.config().EdgeNotHomeLen {
		return true
	}

	rb := r.GetMapData().GetRadiusBlock(radius)
	if ux < rb.MinX+r.config().EdgeNotHomeLen || uy < rb.MinY+r.config().EdgeNotHomeLen || ux+r.config().EdgeNotHomeLen > rb.MaxX || uy+r.config().EdgeNotHomeLen > rb.MaxY {
		return true
	}

	return false
}

func (r *Realm) randomXiongNuBasePos(originX, originY, minRange, maxRange int) (ok bool, x int, y int) {

	radius := r.GetRadius()

	r.conflict.Lock()
	defer r.conflict.Unlock()

	minRange = imath.Min(minRange, maxRange)
	maxRange = imath.Max(minRange, maxRange)
	minIX := originX - minRange
	maxIX := originX + minRange
	minOX := originX - maxRange
	maxOX := originX + maxRange
	minIY := originY - minRange
	maxIY := originY + minRange
	minOY := originY - maxRange
	maxOY := originY + maxRange

	isValidPos := func(x, y int) bool {
		// 在内圈中，返回false
		if minIX < x && x < maxIX && minIY < y && y < maxIY {
			return false
		}

		// 在外圈之外，返回false
		if x < minOX || x > maxOX || y < minOY || y > maxOY {
			return false
		}

		if !r.mapData.IsValidHomePosition(x, y) {
			return false
		}

		if r.isEdgeNotHomePos(x, y, radius) {
			return false
		}

		return r.conflict.addBaseIfCanAddUnderLocker(x, y)
	}

	originBlockX, originBlockY := r.mapData.GetBlockByPos(originX, originY)

	// 在自己这个地图块中随机10次
	for i := 0; i < 10; i++ {
		rx, ry := r.mapData.RandomHomeXY(originBlockX, originBlockY)
		if isValidPos(rx, ry) {
			return true, rx, ry
		}
	}

	diffRange := imath.Max(maxRange-minRange, 5)
	isValidPosWithRange := func(x, y, ox, oy int) (bool, int, int) {
		if isValidPos(x+ox, y+oy) {
			return true, x + ox, y + oy
		}
		return false, 0, 0
	}

	// 做10次随机
	for i := 0; i < 10; i++ {
		d0 := minRange + rand.Intn(diffRange)
		d1 := rand.Intn(maxRange)
		if ok, rx, ry := isValidPosWithRange(originX, originY, d0, d1); ok {
			return true, rx, ry
		}
		if ok, rx, ry := isValidPosWithRange(originX, originY, -d0, d1); ok {
			return true, rx, ry
		}
		if ok, rx, ry := isValidPosWithRange(originX, originY, d0, -d1); ok {
			return true, rx, ry
		}
		if ok, rx, ry := isValidPosWithRange(originX, originY, -d0, -d1); ok {
			return true, rx, ry
		}

		if ok, rx, ry := isValidPosWithRange(originX, originY, d1, d0); ok {
			return true, rx, ry
		}
		if ok, rx, ry := isValidPosWithRange(originX, originY, -d1, d0); ok {
			return true, rx, ry
		}
		if ok, rx, ry := isValidPosWithRange(originX, originY, d1, -d0); ok {
			return true, rx, ry
		}
		if ok, rx, ry := isValidPosWithRange(originX, originY, -d1, -d0); ok {
			return true, rx, ry
		}
	}

	// 随机之后还找不到，遍历
	for x := originX - maxRange; x < originX+maxRange; x++ {
		for y := originY - maxRange; y < originY+maxRange; y++ {
			if isValidPos(x, y) {
				return true, x, y
			}
		}
	}

	return
}

func (r *Realm) randomBlockBasePos(blockX, blockY uint64) (ok bool, x, y int) {
	radius := r.GetRadius()

	r.conflict.Lock()
	defer r.conflict.Unlock()

	// 随机遍历
	r.levelData.Block.RangeBlockHomeCubes(blockX, blockY, func(c cb.Cube) (toContinue bool) {
		rx, ry := c.XY()

		if !r.mapData.IsValidHomePosition(rx, ry) {
			return true
		}

		if r.isEdgeNotHomePos(rx, ry, radius) {
			return true
		}

		if r.conflict.addBaseIfCanAddUnderLocker(rx, ry) {
			ok = true
			x, y = rx, ry
			return false
		}
		return true
	})

	return
}

func (r *Realm) randomBasePos() (ok bool, x int, y int) {

	radius := r.GetRadius()

	r.conflict.Lock()
	defer r.conflict.Unlock()

	// 随机主城位置，环形区域随机
	r.regionData.RangeRandomBlocks(func(blockX, blockY uint64) (toContinue bool) {
		// 每个block尝试3次
		for i := 0; i < 3; i++ {
			rx, ry := r.mapData.RandomHomeXY(blockX, blockY)
			if !r.mapData.IsValidHomePosition(rx, ry) {
				return true
			}

			if r.isEdgeNotHomePos(rx, ry, radius) {
				continue
			}

			if r.conflict.addBaseIfCanAddUnderLocker(rx, ry) {
				ok = true
				x, y = rx, ry
				return false
			}
		}

		return true
	})

	return
}

func (r *Realm) RandomBasePos() (x int, y int) {

	radius := r.GetRadius()

	r.conflict.Lock()
	defer r.conflict.Unlock()

	ok := false
	r.mapData.RangeBlock(radius, blockdata.BlockRangeTypeRandom, func(blockX, blockY uint64) (toContinue bool) {

		// 每个block尝试3次
		for i := 0; i < 3; i++ {
			rx, ry := r.mapData.RandomHomeXY(blockX, blockY)
			if !r.mapData.IsValidHomePosition(rx, ry) {
				return true
			}

			if r.isEdgeNotHomePos(rx, ry, radius) {
				continue
			}

			if r.conflict.addBaseIfCanAddUnderLocker(rx, ry) {
				ok = true
				x, y = rx, ry
				return false
			}
		}

		return true
	})

	if !ok {
		// 每个block都来了一次，还不行？
		// 获取人数最少的block，遍历一次
		if least := r.blockManager.getLeastHomeCountBlock(); least != nil {
			r.mapData.RangeBlockHomeCubes(least.blockX, least.blockY, func(c cb.Cube) (toContinue bool) {
				rx, ry := c.XY()
				if !r.mapData.IsValidHomePosition(rx, ry) {
					return true
				}

				if r.isEdgeNotHomePos(rx, ry, radius) {
					return true
				}

				if r.conflict.addBaseIfCanAddUnderLocker(rx, ry) {
					ok = true
					x, y = rx, ry
					return false
				}

				return true
			})
		}

		if !ok {
			// 如果还找不到，随机，不管了（保证不会出现在边缘位置）
			ok = true

			x, y = r.mapData.RandomHomeXY(r.mapData.RandomBlock(radius))
			r.conflict.doAddBase(x, y)
		}
	}

	return
}

func (r *Realm) RandomAroundBase(centerX, centerY int) (x, y int, ok bool) {

	radius := r.GetRadius()
	f := func(x, y int) bool {
		if !r.mapData.IsValidHomePosition(x, y) {
			return false
		}

		if r.isEdgeNotHomePos(x, y, radius) {
			return false
		}

		return r.conflict.addBaseIfCanAddUnderLocker(x, y)
	}

	cube, ok := blockdata.RandomSingleRingCube(centerX, centerY, r.levelData.GuildMoveBaseMinRadius, r.levelData.GuildMoveBaseMaxRadius, f)
	if ok {
		x, y = cube.XY()
	}

	return
}

func (r *Realm) AroundBase(x1, y1, x2, y2 int) bool {
	return u64.FromInt(hexagon.Distance(x1, y1, x2, y2)) <= r.levelData.GuildMoveBaseMaxRadius
}

var (
	conflictCubeCacheLock = sync.Mutex{}
	conflictCubeCache     = make(map[cb.Cube][]cb.Cube)
)

func NeighborsInBaseConflictRange(x, y int) []cb.Cube {
	conflictCubeCacheLock.Lock()
	defer conflictCubeCacheLock.Unlock()

	cb := cb.XYCube(x, y)
	if result, has := conflictCubeCache[cb]; has {
		return result
	} else {
		result = hexagon.SpiralRing(x, y, constants.BaseConflictRange)
		conflictCubeCache[cb] = result
		return result
	}
}
