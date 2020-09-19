package hexagon

import (
	. "github.com/lightpaw/male7/entity/cb"
	"github.com/lightpaw/male7/util/imath"
	. "github.com/onsi/gomega"
	"math/rand"
	"testing"
	"time"
)

func TestNeighbors(t *testing.T) {
	RegisterTestingT(t)

	Ω(Neighbors(0, 0)).Should(BeEquivalentTo(evenDirections))

	Ω(Neighbors(1, 1)).Should(BeEquivalentTo([]Cube{
		XYCube(2, 2), XYCube(2, 1), XYCube(1, 0),
		XYCube(0, 1), XYCube(0, 2), XYCube(1, 2)}))

	Ω(Neighbors(4, 3)).Should(BeEquivalentTo([]Cube{
		XYCube(5, 3), XYCube(5, 2), XYCube(4, 2),
		XYCube(3, 2), XYCube(3, 3), XYCube(4, 4)}))
}

func TestIsNeighbors(t *testing.T) {
	RegisterTestingT(t)

	for _, n := range evenDirections {
		nx, ny := n.XY()
		Ω(IsNeighbors(0, 0, nx, ny)).Should(BeTrue())
	}

	for _, n := range []Cube{
		XYCube(2, 2), XYCube(2, 1), XYCube(1, 0),
		XYCube(0, 1), XYCube(0, 2), XYCube(1, 2)} {
		nx, ny := n.XY()
		Ω(IsNeighbors(1, 1, nx, ny)).Should(BeTrue())
	}

	Ω(IsNeighbors(1, 1, 0, 0)).Should(BeFalse())
	Ω(IsNeighbors(1, 1, 2, 0)).Should(BeFalse())

	for _, n := range []Cube{
		XYCube(5, 3), XYCube(5, 2), XYCube(4, 2),
		XYCube(3, 2), XYCube(3, 3), XYCube(4, 4)} {
		nx, ny := n.XY()
		Ω(IsNeighbors(4, 3, nx, ny)).Should(BeTrue())
	}

	Ω(IsNeighbors(4, 3, 3, 4)).Should(BeFalse())
	Ω(IsNeighbors(4, 3, 5, 4)).Should(BeFalse())
}

func TestRing(t *testing.T) {
	RegisterTestingT(t)

	Ω(Ring(0, 0, 0)).Should(BeEquivalentTo([]Cube{XYCube(0, 0)}))
	Ω(Ring(1, 1, 0)).Should(BeEquivalentTo([]Cube{XYCube(1, 1)}))
	Ω(Ring(4, 3, 0)).Should(BeEquivalentTo([]Cube{XYCube(4, 3)}))

	Ω(Ring(0, 0, 1)).Should(ConsistOf(evenDirections))

	Ω(Ring(1, 1, 1)).Should(ConsistOf([]Cube{
		XYCube(2, 2), XYCube(2, 1), XYCube(1, 0),
		XYCube(0, 1), XYCube(0, 2), XYCube(1, 2)}))

	Ω(Ring(4, 3, 1)).Should(ConsistOf([]Cube{
		XYCube(5, 3), XYCube(5, 2), XYCube(4, 2),
		XYCube(3, 2), XYCube(3, 3), XYCube(4, 4)}))

	Ω(Ring(4, 3, 2)).Should(ConsistOf([]Cube{
		XYCube(2, 4), XYCube(3, 4),
		XYCube(4, 5), XYCube(5, 4),
		XYCube(6, 4), XYCube(6, 3),
		XYCube(6, 2), XYCube(5, 1),
		XYCube(4, 1), XYCube(3, 1),
		XYCube(2, 2), XYCube(2, 3)}))
}

func TestSpiralRing(t *testing.T) {
	RegisterTestingT(t)

	Ω(SpiralRing(0, 0, 0)).Should(BeEquivalentTo([]Cube{XYCube(0, 0)}))
	Ω(SpiralRing(1, 1, 0)).Should(BeEquivalentTo([]Cube{XYCube(1, 1)}))
	Ω(SpiralRing(4, 3, 0)).Should(BeEquivalentTo([]Cube{XYCube(4, 3)}))

	Ω(SpiralRing(0, 0, 1)).Should(ConsistOf(append([]Cube{XYCube(0, 0)}, evenDirections...)))

	Ω(SpiralRing(1, 1, 1)).Should(ConsistOf([]Cube{XYCube(1, 1),
		XYCube(2, 2), XYCube(2, 1), XYCube(1, 0),
		XYCube(0, 1), XYCube(0, 2), XYCube(1, 2)}))

	Ω(SpiralRing(4, 3, 1)).Should(ConsistOf([]Cube{XYCube(4, 3),
		XYCube(5, 3), XYCube(5, 2), XYCube(4, 2),
		XYCube(3, 2), XYCube(3, 3), XYCube(4, 4)}))

	Ω(SpiralRing(4, 3, 2)).Should(ConsistOf([]Cube{XYCube(4, 3), // 原点
		XYCube(5, 3), XYCube(5, 2), XYCube(4, 2),
		XYCube(3, 2), XYCube(3, 3), XYCube(4, 4), // 第一层
		XYCube(2, 4), XYCube(3, 4),
		XYCube(4, 5), XYCube(5, 4),
		XYCube(6, 4), XYCube(6, 3),
		XYCube(6, 2), XYCube(5, 1),
		XYCube(4, 1), XYCube(3, 1),
		XYCube(2, 2), XYCube(2, 3)}))
}

func TestEventOffset(t *testing.T) {
	RegisterTestingT(t)

	// 偏移换换
	for i, n := range evenDirections {

		oddOffsetX, oddOffsetY := evenOffset2OddOffset(n.XY())

		odd := oddDirections[i]
		oddX, oddY := odd.XY()
		Ω(oddOffsetX).Should(Equal(oddX))
		Ω(oddOffsetY).Should(Equal(oddY))
	}

	// 偶数加偶数偏移
	Ω(ShiftEvenOffset(0, 0, 0, 0)).Should(Equal(XYCube(0, 0)))
	Ω(ShiftEvenOffset(0, 0, 1, 0)).Should(Equal(XYCube(1, 0)))
	Ω(ShiftEvenOffset(0, 0, 2, 0)).Should(Equal(XYCube(2, 0)))
	Ω(ShiftEvenOffset(0, 0, 0, 1)).Should(Equal(XYCube(0, 1)))
	Ω(ShiftEvenOffset(0, 0, 0, 2)).Should(Equal(XYCube(0, 2)))
	Ω(ShiftEvenOffset(0, 0, 1, 1)).Should(Equal(XYCube(1, 1)))
	Ω(ShiftEvenOffset(4, 3, 1, 3)).Should(Equal(XYCube(5, 6)))

	// 奇数加偶数偏移
	Ω(ShiftEvenOffset(1, 1, 0, 0)).Should(Equal(XYCube(1, 1)))
	Ω(ShiftEvenOffset(1, 1, 1, 0)).Should(Equal(XYCube(2, 2)))
	Ω(ShiftEvenOffset(1, 1, 2, 0)).Should(Equal(XYCube(3, 1)))
	Ω(ShiftEvenOffset(1, 1, 0, 1)).Should(Equal(XYCube(1, 2)))
	Ω(ShiftEvenOffset(1, 1, 0, 2)).Should(Equal(XYCube(1, 3)))
	Ω(ShiftEvenOffset(1, 1, 1, 1)).Should(Equal(XYCube(2, 3)))
	Ω(ShiftEvenOffset(5, 3, 1, 3)).Should(Equal(XYCube(6, 7)))

	Ω(ShiftEvenOffset(1, 1, -1, 0)).Should(Equal(XYCube(0, 2)))
	Ω(ShiftEvenOffset(1, 1, -2, 0)).Should(Equal(XYCube(-1, 1)))
	Ω(ShiftEvenOffset(1, 1, 0, -1)).Should(Equal(XYCube(1, 0)))
	Ω(ShiftEvenOffset(1, 1, 0, -2)).Should(Equal(XYCube(1, -1)))
	Ω(ShiftEvenOffset(1, 1, -1, -1)).Should(Equal(XYCube(0, 1)))
	Ω(ShiftEvenOffset(5, 3, -1, -3)).Should(Equal(XYCube(4, 1)))

	// 反向获取偏移
	// 偶数加偶数偏移
	Ω(EvenOffsetBetween(0, 0, 0, 0)).Should(Equal(XYCube(0, 0)))
	Ω(EvenOffsetBetween(0, 0, 1, 0)).Should(Equal(XYCube(1, 0)))
	Ω(EvenOffsetBetween(0, 0, 2, 0)).Should(Equal(XYCube(2, 0)))
	Ω(EvenOffsetBetween(0, 0, 0, 1)).Should(Equal(XYCube(0, 1)))
	Ω(EvenOffsetBetween(0, 0, 0, 2)).Should(Equal(XYCube(0, 2)))
	Ω(EvenOffsetBetween(0, 0, 1, 1)).Should(Equal(XYCube(1, 1)))
	Ω(EvenOffsetBetween(4, 3, 5, 6)).Should(Equal(XYCube(1, 3)))

	// 奇数加偶数偏移
	Ω(EvenOffsetBetween(1, 1, 1, 1)).Should(Equal(XYCube(0, 0)))
	Ω(EvenOffsetBetween(1, 1, 2, 2)).Should(Equal(XYCube(1, 0)))
	Ω(EvenOffsetBetween(1, 1, 3, 1)).Should(Equal(XYCube(2, 0)))
	Ω(EvenOffsetBetween(1, 1, 1, 2)).Should(Equal(XYCube(0, 1)))
	Ω(EvenOffsetBetween(1, 1, 1, 3)).Should(Equal(XYCube(0, 2)))
	Ω(EvenOffsetBetween(1, 1, 2, 3)).Should(Equal(XYCube(1, 1)))
	Ω(EvenOffsetBetween(5, 3, 6, 7)).Should(Equal(XYCube(1, 3)))

	Ω(EvenOffsetBetween(1, 1, 0, 2)).Should(Equal(XYCube(-1, 0)))
	Ω(EvenOffsetBetween(1, 1, -1, 1)).Should(Equal(XYCube(-2, 0)))
	Ω(EvenOffsetBetween(1, 1, 1, 0)).Should(Equal(XYCube(0, -1)))
	Ω(EvenOffsetBetween(1, 1, 1, -1)).Should(Equal(XYCube(0, -2)))
	Ω(EvenOffsetBetween(1, 1, 0, 1)).Should(Equal(XYCube(-1, -1)))
	Ω(EvenOffsetBetween(5, 3, 4, 1)).Should(Equal(XYCube(-1, -3)))

}

func TestDistance(t *testing.T) {
	RegisterTestingT(t)

	x1, y1, z1 := oddq_to_cube(0, 0)
	Ω(x1).Should(Equal(0))
	Ω(y1).Should(Equal(0))
	Ω(z1).Should(Equal(0))

	xx := []int{0, 10, 234, -255, -100}
	yy := []int{0, 10, -321, 32000, -200}

	for i, x := range xx {
		y := yy[i]

		x1, y1, z1 = oddq_to_cube(x, y)

		for i := 0; i < 100; i++ {
			for _, c := range Ring(x, y, uint(i)) {
				cx, cy := c.XY()
				Ω(Distance(x, y, cx, cy)).Should(Equal(i))

				x2, y2, z2 := oddq_to_cube(cx, cy)
				Ω(cube_distance(x1, y1, z1, x2, y2, z2)).Should(Equal(i))
				Ω(cube_distance2(x1, y1, z1, x2, y2, z2)).Should(Equal(i))
			}
		}
	}

	Ω(cube_distance(0, 0, 0, 0, 0, 0)).Should(Equal(0))
	Ω(cube_distance(1, 0, 0, 0, 0, 0)).Should(Equal(1))
	Ω(cube_distance(-1, 0, 0, 0, 0, 0)).Should(Equal(1))
	Ω(cube_distance(1, 0, 0, 1, 0, 0)).Should(Equal(0))
	Ω(cube_distance(-1, 0, 0, 1, 0, 0)).Should(Equal(2))
	Ω(cube_distance(1, 2, -3, 0, 0, 0)).Should(Equal(3))
}

func TestSpiralSort(t *testing.T) {
	RegisterTestingT(t)
	rand.Seed(time.Now().UnixNano())

	cubes := []Cube{XYCube(1, 1), XYCube(10, 10), XYCube(100, 100), XYCube(1000, 1000)}

	actual := make([]Cube, len(cubes))
	copy(actual, cubes)
	Mix(actual)

	SpiralSort(actual, 0, 0)
	Ω(actual).Should(Equal([]Cube{XYCube(1, 1), XYCube(10, 10), XYCube(100, 100), XYCube(1000, 1000)}))

	SpiralSort(actual, 1001, 1001)
	Ω(actual).Should(Equal([]Cube{XYCube(1000, 1000), XYCube(100, 100), XYCube(10, 10), XYCube(1, 1)}))

	SpiralSort(actual, 30, 30)
	Ω(actual).Should(Equal([]Cube{XYCube(10, 10), XYCube(1, 1), XYCube(100, 100), XYCube(1000, 1000)}))

}

func BenchmarkDistance1(b *testing.B) {
	for i := 0; i < b.N; i++ {
		cube_distance(-1, 200, 30, -400, -500, 60)
	}
}

func cube_distance2(x1, y1, z1, x2, y2, z2 int) int {
	return (imath.Abs(x1-x2) + imath.Abs(y1-y2) + imath.Abs(z1-z2)) / 2
}

func BenchmarkDistance2(b *testing.B) {
	for i := 0; i < b.N; i++ {
		cube_distance2(-1, 200, 30, -400, -500, 60)
	}
}

func BenchmarkOddqToCube(b *testing.B) {
	for i := 0; i < b.N; i++ {
		oddq_to_cube(100, -100)
	}
}

func BenchmarkDistance(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Distance(100, -100, -200, 300)
	}
}
