package cb

import (
	. "github.com/onsi/gomega"
	"math/rand"
	"testing"
	"time"
)

func TestCubesMix(t *testing.T) {
	RegisterTestingT(t)
	rand.Seed(time.Now().UnixNano())

	list := []Cube{XYCube(0, 0), XYCube(1, 1), XYCube(2, 2), XYCube(3, 3)}
	mix := make([]Cube, len(list))
	copy(mix, list)
	mixCubes := Cubes(mix)

	mixCubes.Mix()
	Ω(mixCubes).Should(ConsistOf(list[0], list[1], list[2], list[3]))
}

func TestCubes(t *testing.T) {
	RegisterTestingT(t)

	list := []Cube{XYCube(0, 0), XYCube(1, 1), XYCube(2, 2), XYCube(3, 3)}
	cubes := Cubes(list)

	var actual []Cube
	cubes.Range(func(c Cube) (toContinue bool) {
		actual = append(actual, c)
		return true
	})

	Ω(actual).Should(Equal(list))

	actual = actual[0:0]
	cubes.RandomRange(func(c Cube) (toContinue bool) {
		actual = append(actual, c)
		return true
	})

	Ω(actual).Should(ConsistOf(list[0], list[1], list[2], list[3]))

	for i := 0; i < 10; i++ {
		Ω(list).Should(ContainElement(cubes.Random()))
	}
}
