package hexagon

import (
	. "github.com/lightpaw/male7/entity/cb"
	"github.com/lightpaw/male7/util/imath"
	"math"
	"sort"
)

// http://www.redblobgames.com/grids/hexagons/
// TODO 只是处理用到的，以后再将4种情况整合 odd-rbt even-rbt odd-q event-q
// 当前使用的是 odd-q

var (
	// 逆时针旋转， 旋转开始于 index=4的那个方向
	ringStartIndex = 4
	evenDirections = []Cube{
		XYCube(1, 0), XYCube(1, -1), XYCube(0, -1),
		XYCube(-1, -1), XYCube(-1, 0), XYCube(0, 1),
	}

	oddDirections = []Cube{
		XYCube(1, 1), XYCube(1, 0), XYCube(0, -1),
		XYCube(-1, 0), XYCube(-1, 1), XYCube(0, 1),
	}
)

func Neighbors(x, y int) []Cube {

	directions := oddDirections
	if x&1 == 0 {
		directions = evenDirections
	}

	nb := make([]Cube, 0, len(directions))
	for _, grid := range directions {
		ox, oy := grid.XY()
		nbx := x + ox
		nby := y + oy
		nb = append(nb, XYCube(nbx, nby))
	}

	return nb
}

func IsNeighbors(x, y, nx, ny int) bool {
	directions := oddDirections
	if x&1 == 0 {
		directions = evenDirections
	}

	for _, grid := range directions {
		ox, oy := grid.XY()
		nbx := x + ox
		nby := y + oy
		if nbx == nx && nby == ny {
			return true
		}
	}

	return false
}

// 获取一环坐标
func Ring(x, y int, radius uint) []Cube {
	r := int(radius)
	ring := make([]Cube, 0, 6*r)
	ring = putRing(ring, x, y, r)
	return ring
}

// 螺旋区域坐标，包含XY
func SpiralRing(x, y int, radius uint) []Cube {
	r := int(radius)

	n := 1
	for i := 1; i <= r; i++ {
		n += i * 6
	}

	ring := make([]Cube, 0, n)
	for i := 0; i <= r; i++ {
		ring = putRing(ring, x, y, i)
	}

	return ring
}

func putRing(ring []Cube, x, y, radius int) []Cube {
	if radius <= 0 {
		ring = append(ring, XYCube(x, y))
		return ring
	}

	for i := 0; i < radius; i++ {
		x, y = directionNeighbor(x, y, ringStartIndex)
	}

	for i := 0; i < 6; i++ {
		for ra := 0; ra < radius; ra++ {
			ring = append(ring, XYCube(x, y))

			x, y = directionNeighbor(x, y, i)
		}
	}

	return ring
}

// 以x,y为中心点获取某个方向的相邻点
func directionNeighbor(x, y, di int) (int, int) {
	if x&1 == 0 {
		return evenDirections[di].AddXY(x, y)
	} else {
		return oddDirections[di].AddXY(x, y)
	}
}

// 偶数偏移
func ShiftEvenOffset(originX, originY, evenOffsetX, eventOffsetY int) Cube {

	if originX&1 == 1 {
		oddOffsetX, oddOffsetY := evenOffset2OddOffset(evenOffsetX, eventOffsetY)
		return XYCube2(originX, originY, oddOffsetX, oddOffsetY)
	}

	return XYCube2(originX, originY, evenOffsetX, eventOffsetY)
}

func evenOffset2OddOffset(evenOffsetX, eventOffsetY int) (oddOffsetX, oddOffsetY int) {
	if evenOffsetX&1 == 1 {
		// 奇数偏移列，Y+1
		return evenOffsetX, eventOffsetY + 1
	}

	return evenOffsetX, eventOffsetY
}

func EvenOffsetBetween(originX, originY, targetX, targetY int) Cube {
	if originX&1 == 1 {
		oddOffsetX, oddOffsetY := targetX-originX, targetY-originY
		return XYCube(oddOffset2EvenOffset(oddOffsetX, oddOffsetY))
	}

	return XYCube(targetX-originX, targetY-originY)
}

func oddOffset2EvenOffset(oddOffsetX, oddOffsetY int) (evenOffsetX, eventOffsetY int) {
	if oddOffsetX&1 == 1 {
		// 奇数偏移列，Y-1
		return oddOffsetX, oddOffsetY - 1
	}

	return oddOffsetX, oddOffsetY
}

func Distance(ox1, oy1, ox2, oy2 int) int {
	x1, y1, z1 := oddq_to_cube(ox1, oy1)
	x2, y2, z2 := oddq_to_cube(ox2, oy2)

	return cube_distance(x1, y1, z1, x2, y2, z2)
}

func oddq_to_cube(ox, oy int) (x, y, z int) {
	x = ox
	z = oy - (ox-(ox&1))/2
	y = -x - z

	return
}

func cube_distance(x1, y1, z1, x2, y2, z2 int) int {
	return imath.Max(imath.Max(imath.Abs(x1-x2), imath.Abs(y1-y2)), imath.Abs(z1-z2))
}

func SpiralSort(cubes []Cube, centX, centY int) {

	scores := make([]int, len(cubes))
	for i, c := range cubes {
		x, y := c.XY()
		scores[i] = Distance(x, y, centX, centY)
	}

	// 排序
	sort.Sort(&spiral_slice{
		cubes:  cubes,
		scores: scores,
	})
}

type spiral_slice struct {
	cubes  []Cube
	scores []int
}

func (a *spiral_slice) Len() int { return len(a.cubes) }
func (a *spiral_slice) Swap(i, j int) {
	a.cubes[i], a.cubes[j] = a.cubes[j], a.cubes[i]
	a.scores[i], a.scores[j] = a.scores[j], a.scores[i]
}
func (a *spiral_slice) Less(i, j int) bool { return a.scores[i] < a.scores[j] }

func OffsetDistance(x1, y1, x2, y2 int, size float64) float64 {
	if x1 == x2 && y1 == y2 {
		return 0
	}

	px1, py1 := oddq_offset_to_pixel(x1, y1, size)
	px2, py2 := oddq_offset_to_pixel(x2, y2, size)
	dx, dy := px1-px2, py1-py2

	return math.Sqrt(dx*dx + dy*dy)
}

var sqrt3 = math.Sqrt(3)

func oddq_offset_to_pixel(cx, cy int, size float64) (px, py float64) {
	px = size * 3 / 2 * float64(cx)
	py = size * sqrt3 * (float64(cy) + 0.5*float64(cx&1))
	return
}
