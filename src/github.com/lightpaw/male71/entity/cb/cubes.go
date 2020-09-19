package cb

import "math/rand"

type Cubes []Cube

func (cs Cubes) Mix() {
	Mix(cs)
}

func (cs Cubes) Random() Cube {
	if n := len(cs); n > 0 {
		return cs[rand.Intn(n)]
	}

	return 0
}

func (cs Cubes) Range(f func(c Cube) (toContinue bool)) (toContinue bool) {
	for _, c := range cs {
		if !f(c) {
			return false
		}
	}

	return true
}

func (cs Cubes) RandomRange(f func(c Cube) (toContinue bool)) (toContinue bool) {

	split := rand.Intn(len(cs))
	for i := split; i < len(cs); i++ {
		if !f(cs[i]) {
			return false
		}
	}

	for i := 0; i < split; i++ {
		if !f(cs[i]) {
			return false
		}
	}

	return true
}
