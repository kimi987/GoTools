package u64

func Plus(x, y uint64) uint64 {
	return x + y
}

func Sub(x, y uint64) uint64 {
	if x < y {
		return 0
	}

	return x - y
}

func AddInt(o uint64, toAdd int) uint64 {
	if toAdd < 0 {
		return SubInt(o, -toAdd)
	}

	return o + uint64(toAdd)
}

func SubInt(o uint64, toReduce int) uint64 {
	if toReduce < 0 {
		return AddInt(o, -toReduce)
	}

	return Sub(o, uint64(toReduce))
}

func Plus2Float64(o, toAdd uint64) float64 {
	return float64(o + toAdd)
}

func Sub2Float64(o, toReduce uint64) float64 {
	if o < toReduce {
		return -float64(toReduce - o)
	} else {
		return float64(o - toReduce)
	}
}

// 此方法就是简单的数值计算
func Multi(d uint64, multi float64) uint64 {
	if d == 0 || multi <= 0 {
		return 0
	}

	return uint64(float64(d) * multi)
}

// 这个方法跟Multi比较，在d不是非常大，multi精度大于3位小数，处理了这个范围内计算出来的值，精度损失的问题
// 比如（20 * 0.05 != 1）的问题，用于一些策划配置的系数计算
// 数字太大会失真，浮点数小数超过3位会失真
func MultiCoef(d uint64, multi float64) uint64 {
	if d == 0 || multi <= 0 {
		return 0
	}

	return d * f1000(multi) / 1000
}

func f1000(f float64) uint64 {
	return uint64(f*1000 + 0.0001)
}

func MultiF64(d uint64, multi float64) uint64 {
	if d == 0 || multi <= 0 {
		return 0
	}

	fd := float64(d)
	return uint64((multi + (1 / (fd * 10))) * fd)
}

func Division2Float64(x, y uint64) float64 {
	if y == 0 {
		return 0
	}

	return float64(x) / float64(y)
}
