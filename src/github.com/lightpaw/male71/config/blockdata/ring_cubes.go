package blockdata

import (
	"github.com/lightpaw/male7/entity/cb"
	"github.com/lightpaw/male7/entity/hexagon"
)

func NewMultiRingCubes(x, y int, mix bool, validate func(x, y int) bool, radius ...uint64) []cb.Cubes {

	var arr []cb.Cubes

	minRaduis := uint64(0)
	for _, r := range radius {
		if r > 0 && minRaduis <= r {
			arr = append(arr, NewSingleRingCubes(x, y, minRaduis, r+1, mix, validate))
			minRaduis = r + 1
		}
	}

	return arr
}

func NewSingleRingCubes(x, y int, minRadius, maxRadius uint64, mix bool, validate func(x, y int) bool) cb.Cubes {

	var cubes cb.Cubes
	for i := minRadius; i < maxRadius; i++ {
		for _, c := range hexagon.Ring(x, y, uint(i)) {
			if validate(c.XY()) {
				cubes = append(cubes, c)
			}
		}
	}

	if mix {
		cubes.Mix()
	}

	return cubes
}

func RandomSingleRingCube(x, y int, minRadius, maxRadius uint64, validate func(x, y int) bool) (c cb.Cube, ok bool) {
	for startRadius, endRadius := minRadius, minRadius+1; endRadius <= maxRadius; {
		c, ok = newSingleRingCube(x, y, startRadius, endRadius, validate)
		startRadius++
		endRadius++
		if ok {
			return
		}
	}
	return
}

func newSingleRingCube(x, y int, minRadius, maxRadius uint64, validate func(x, y int) bool) (c cb.Cube, ok bool) {
	for i := minRadius; i < maxRadius; i++ {
		cbs := hexagon.Ring(x, y, uint(i))
		cb.Mix(cbs)
		for _, rc := range cbs {
			if validate(rc.XY()) {
				x, y = rc.XY()
				return rc, true
			}
		}
	}

	return
}
