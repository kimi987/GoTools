package cb

import (
	"math/rand"
	"math"
)

type Cube uint64

func (c Cube) Scale(n int) Cube {
	x, y := c.XY()
	return XYCube(x*n, y*n)
}

func (c Cube) add(toAdd Cube) Cube {
	x1, y1 := c.XY()
	x2, y2 := toAdd.XY()
	return XYCube(x1+x2, y1+y2)
}

func (c Cube) AddXY(x, y int) (int, int) {
	x1, y1 := c.XY()
	return x1 + x, y1 + y
}

func XYCube(x, y int) Cube {
	return int2cube(x)<<16 | int2cube(y)
}

func XYCube2(x1, y1, x2, y2 int) Cube {
	return XYCube(x1+x2, y1+y2)
}

func XYCubeI32(x, y int32) Cube {
	return XYCube(int(x), int(y))
}

func (c Cube) XY() (int, int) {
	return CubeXY(c)
}

func (c Cube) XYI32() (int32, int32) {
	x, y := c.XY()
	return int32(x), int32(y)
}

const Max = math.MaxInt16
const Min = math.MinInt16

func CubeXY(c Cube) (int, int) {
	x := c >> 16
	y := c & 0xffff

	return cube2int(x), cube2int(y)
}

func int2cube(x int) Cube {
	if x < 0 {
		return Cube(-x)<<1 + 1
	} else {
		return Cube(x) << 1
	}
}

func cube2int(x Cube) int {
	if x&1 == 0 {
		return int(x >> 1)
	} else {
		return -int(x >> 1)
	}
}

func Contains(array []Cube, x Cube) bool {

	for _, v := range array {
		if v == x {
			return true
		}
	}

	return false
}

func RemoveIntersectCube(toRemoves, toAdds []Cube) ([]Cube, []Cube) {
	nr := len(toRemoves)
	na := len(toAdds)

	if nr == 0 || na == 0 {
		return toRemoves, toAdds
	}

	intersectCount := 0
loop:
	for i := 0; i < nr; i++ {
		v := toRemoves[i]

		n := intersectCount
		for j := n; j < na; j++ {
			if v == toAdds[j] {
				toRemoves[n], toRemoves[i] = toRemoves[i], toRemoves[n]
				toAdds[n], toAdds[j] = toAdds[j], toAdds[n]
				intersectCount++
				continue loop
			}
		}
	}

	n := intersectCount
	return toRemoves[n:], toAdds[n:]
}

func Mix(list []Cube) {
	// 循环设置每个位置，打乱顺序
	n := len(list)
	for i := range list {
		idx := n - i - 1
		swap := rand.Intn(n - i)
		list[idx], list[swap] = list[swap], list[idx]
	}
}

func Mix2(list []Cube) {
	// 循环设置每个位置，打乱顺序
	n := len(list)
	for i := range list {
		swap := i + rand.Intn(n-i)
		list[i], list[swap] = list[swap], list[i]
	}
}
