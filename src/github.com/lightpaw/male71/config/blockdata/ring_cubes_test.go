package blockdata

import (
	"github.com/lightpaw/male7/entity/cb"
	. "github.com/onsi/gomega"
	"testing"
	"fmt"
)

func TestRingCubes(t *testing.T) {
	RegisterTestingT(t)

	cubes := NewSingleRingCubes(0, 0, 0, 1, true, func(x, y int) bool {
		return x >= 0 && y >= 0
	})

	var arr []cb.Cube
	cubes.Range(func(c cb.Cube) (toContinue bool) {
		x, y := c.XY()
		Ω(x >= 0 && y >= 0).Should(BeTrue())

		arr = append(arr, c)
		return true
	})

	Ω(arr).Should(ConsistOf(cb.XYCube(0, 0)))

	cubes = NewSingleRingCubes(4, 5, 5, 10, true, func(x, y int) bool {
		return x >= 0 && y >= 0
	})

	arr = arr[0:0]
	cubes.Range(func(c cb.Cube) (toContinue bool) {
		x, y := c.XY()
		Ω(x >= 0 && y >= 0).Should(BeTrue())

		arr = append(arr, c)
		return true
	})

	cubes = NewSingleRingCubes(4, 5, 5, 5, true, func(x, y int) bool {
		return x >= 0 && y >= 0
	})

	arr = arr[0:0]
	cubes.Range(func(c cb.Cube) (toContinue bool) {
		x, y := c.XY()
		Ω(x >= 0 && y >= 0).Should(BeTrue())

		arr = append(arr, c)
		return true
	})

	Ω(arr).Should(BeEmpty())

	cubes = NewSingleRingCubes(4, 5, 5, 1, true, func(x, y int) bool {
		return x >= 0 && y >= 0
	})

	arr = arr[0:0]
	cubes.Range(func(c cb.Cube) (toContinue bool) {
		x, y := c.XY()
		Ω(x >= 0 && y >= 0).Should(BeTrue())

		arr = append(arr, c)
		return true
	})

	Ω(arr).Should(BeEmpty())
}

func TestGetSpiralBlockXYs(t *testing.T) {

	xys := GetSpiralBlockXYs(2, 2, 2, func(x, y int) bool {
		return true
	})
	for _, xy := range xys {
		x, y := xy.XY()
		fmt.Println(x, y)
	}

}
