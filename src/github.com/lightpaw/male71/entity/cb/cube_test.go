package cb

import "testing"
import (
	. "github.com/onsi/gomega"
	"math/rand"
	"time"
)

func TestMix(t *testing.T) {
	RegisterTestingT(t)
	rand.Seed(time.Now().UnixNano())

	list := []Cube{XYCube(0, 0), XYCube(1, 1), XYCube(2, 2), XYCube(3, 3)}
	mix := make([]Cube, len(list))
	copy(mix, list)

	Ω(mix).Should(Equal(list))

	Mix(mix)
	Ω(mix).Should(ConsistOf(list[0], list[1], list[2], list[3]))

	Mix2(mix)
	Ω(mix).Should(ConsistOf(list[0], list[1], list[2], list[3]))
}

func TestCube(t *testing.T) {
	RegisterTestingT(t)

	len := 512

	cubeMap := map[Cube]int{}
	for x := -len; x < len; x++ {
		for y := -len; y < len; y++ {
			cube := XYCube(x, y)
			n, exist := cubeMap[cube]
			Ω(n).Should(Equal(0))
			Ω(exist).Should(BeFalse())

			cubeMap[cube]++

			newX, newY := cube.XY()
			Ω(newX).Should(Equal(x))
			Ω(newY).Should(Equal(y))
		}
	}

	for x := -len; x < len; x++ {
		for y := -len; y < len; y++ {

			cube := XYCube(x, y)
			n, exist := cubeMap[cube]
			Ω(n).Should(Equal(1))
			Ω(exist).Should(BeTrue())

		}
	}
}
