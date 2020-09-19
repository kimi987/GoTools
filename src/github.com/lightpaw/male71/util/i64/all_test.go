package i64

import (
	"testing"
)

func TestAbs(t *testing.T) {
	input := []int64{
		0,
		-1,
		2,
		-4,
		24,
		48,
		-123,
	}

	want := []int64{
		0,
		1,
		2,
		4,
		24,
		48,
		123,
	}
	for i, n := range input {
		if abs := Abs(n); want[i] != abs {
			t.Errorf("Abs(%d) = %d, want %d", n, abs, want[i])
		}
	}

}

func BenchmarkAbs(b *testing.B) {
	for i := int64(0); i < int64(b.N); i++ {
		Abs(-i)
		Abs(i)
	}
}

func TestPow(t *testing.T) {
	input := []int64{
		0,
		1,
		2,
		3,
		4,
		5,
		6,
		7,
		8,
		9,
		10,
	}

	want5 := []int64{
		1,
		5,
		25,
		125,
		625,
		3125,
		15625,
		78125,
		390625,
		1953125,
		9765625,
	}

	want8 := []int64{
		1,
		8,
		64,
		512,
		4096,
		32768,
		262144,
		2097152,
		16777216,
		134217728,
		1073741824,
	}

	want10 := []int64{
		1,
		10,
		100,
		1000,
		10000,
		100000,
		1000000,
		10000000,
		100000000,
		1000000000,
		10000000000,
	}

	inputSpecial := []int64{
		-1,
		-2,
		-100,
	}

	wantSpecial := []int64{
		0,
		0,
		0,
	}
	for i, n := range input {
		if x := Pow(5, n); want5[i] != x {
			t.Errorf("Pow(%d) = %d, want %d", n, x, want5[i])
		}

		if x := Pow(8, n); want8[i] != x {
			t.Errorf("Pow(%d) = %d, want %d", n, x, want8[i])
		}

		if x := Pow(10, n); want10[i] != x {
			t.Errorf("Pow(%d) = %d, want %d", n, x, want10[i])
		}
	}

	for i, n := range inputSpecial {
		if x := Pow(5, n); wantSpecial[i] != x {
			t.Errorf("Pow(%d) = %d, want %d", n, x, wantSpecial[i])
		}
	}

}

func BenchmarkPow(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Pow(8, 9)
	}
}

func TestMax(t *testing.T) {
	inputX := []int64{
		1,
		-3,
		5,
		7,
		-9,
	}

	inputY := []int64{
		2,
		-4,
		6,
		-8,
		10,
	}

	want := []int64{
		2,
		-3,
		6,
		7,
		10,
	}

	for i, _ := range inputX {
		if n := Max(inputX[i], inputY[i]); n != want[i] {
			t.Errorf("Max(%d, %d) = %d, want %d", inputX[i], inputY[i], n, want[i])
		}
	}
}

func TestMin(t *testing.T) {
	inputX := []int64{
		1,
		-3,
		5,
		7,
		-9,
	}

	inputY := []int64{
		2,
		-4,
		6,
		-8,
		10,
	}

	want := []int64{
		1,
		-4,
		5,
		-8,
		-9,
	}

	for i, _ := range inputX {
		if n := Min(inputX[i], inputY[i]); n != want[i] {
			t.Errorf("Min(%d, %d) = %d, want %d", inputX[i], inputY[i], n, want[i])
		}
	}
}

func BenchmarkMax(b *testing.B) {
	var x, y int64 = 0, 100
	for i := 0; i < b.N; i++ {
		x += Max(x, y)
	}
}

func BenchmarkMin(b *testing.B) {
	var x, y int64 = 0, 100
	for i := 0; i < b.N; i++ {
		x += Min(x, y)
	}
}
