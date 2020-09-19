package util

import (
	"math"
)

func Distance(x1, y1, x2, y2 int) int {
	dx := x1 - x2
	dy := y1 - y2
	return int(math.Sqrt(float64(dx*dx + dy*dy)))
}

func IsInRange(x1, y1, x2, y2, distance int) bool {
	dx := x1 - x2
	dy := y1 - y2
	return dx*dx+dy*dy <= distance*distance
}
