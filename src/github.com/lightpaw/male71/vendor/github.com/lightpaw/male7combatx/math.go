package combatx

import "math"

func Distance(x1, y1, x2, y2 int) int {
	dx := x1 - x2
	dy := y1 - y2
	switch {
	case dx == 0:
		return IAbs(dy)
	case dy == 0:
		return IAbs(dx)
	default:
		return int(math.Sqrt(float64(dx*dx + dy*dy)))
	}
}

func IsInRange(x1, y1, x2, y2, rg int) bool {
	dx := x1 - x2
	dy := y1 - y2
	return dx*dx+dy*dy <= rg*rg
}

func IMin(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func IMax(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func IAbs(a int) int {
	if a >= 0 {
		return a
	}
	return -a
}

func IMinMax(i, min, max int) int {
	return IMin(IMax(i, min), max)
}

func IDivide(a, b int) int {
	if a > 0 && b > 0 {
		return (a + b - 1) / b
	}
	return 0
}

func I32Min(a, b int32) int32 {
	if a < b {
		return a
	}
	return b
}

func I32Max(a, b int32) int32 {
	if a > b {
		return a
	}
	return b
}

func I32MinMax(i, min, max int32) int32 {
	return I32Min(I32Max(i, min), max)
}
