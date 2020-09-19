package blockdata

import (
	. "github.com/onsi/gomega"
	"testing"
)

func TestUint8_map_Int(t *testing.T) {
	RegisterTestingT(t)

	u8map := newUint8Map(10, 10, unwalkable)

	Ω(u8map.SetInt(-1, 0, unwalkable)).Should(BeFalse())
	Ω(u8map.SetInt(0, -1, unwalkable)).Should(BeFalse())
	Ω(u8map.SetInt(-1, -1, unwalkable)).Should(BeFalse())

	Ω(u8map.SetInt(10, 0, unwalkable)).Should(BeFalse())
	Ω(u8map.SetInt(0, 10, unwalkable)).Should(BeFalse())
	Ω(u8map.SetInt(10, 10, unwalkable)).Should(BeFalse())

	Ω(u8map.SetInt(0, 0, unwalkable)).Should(BeTrue())

	for x := 0; x < 10; x++ {
		for y := 0; y < 10; y++ {

			Ω(u8map.GetInt(x, y)).Should(Equal(u8map.Get(uint64(x), uint64(y))))

			if x == 0 && y == 0 {
				Ω(u8map.GetInt(x, y)).Should(Equal(uint8(unwalkable)))
			} else {
				Ω(u8map.GetInt(x, y)).Should(Equal(uint8(walkable)))
			}
		}
	}

	Ω(u8map.GetInt(-1, 0)).Should(Equal(uint8(unwalkable)))
	Ω(u8map.GetInt(0, -1)).Should(Equal(uint8(unwalkable)))
	Ω(u8map.GetInt(-1, -1)).Should(Equal(uint8(unwalkable)))

	Ω(u8map.GetInt(10, 0)).Should(Equal(uint8(unwalkable)))
	Ω(u8map.GetInt(0, 10)).Should(Equal(uint8(unwalkable)))
	Ω(u8map.GetInt(10, 10)).Should(Equal(uint8(unwalkable)))
}

func TestUint8_map_Uint(t *testing.T) {
	RegisterTestingT(t)

	u8map := newUint8Map(10, 10, unwalkable)

	Ω(u8map.Set(10, 0, unwalkable)).Should(BeFalse())
	Ω(u8map.Set(0, 10, unwalkable)).Should(BeFalse())
	Ω(u8map.Set(10, 10, unwalkable)).Should(BeFalse())

	Ω(u8map.Set(0, 0, unwalkable)).Should(BeTrue())

	for x := uint64(0); x < 10; x++ {
		for y := uint64(0); y < 10; y++ {

			Ω(u8map.Get(x, y)).Should(Equal(u8map.GetInt(int(x), int(y))))

			if x == 0 && y == 0 {
				Ω(u8map.Get(x, y)).Should(Equal(uint8(unwalkable)))
			} else {
				Ω(u8map.Get(x, y)).Should(Equal(uint8(walkable)))
			}
		}
	}

	Ω(u8map.Get(10, 0)).Should(Equal(uint8(unwalkable)))
	Ω(u8map.Get(0, 10)).Should(Equal(uint8(unwalkable)))
	Ω(u8map.Get(10, 10)).Should(Equal(uint8(unwalkable)))
}
