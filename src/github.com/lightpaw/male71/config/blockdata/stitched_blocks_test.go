package blockdata

import (
	"testing"
	. "github.com/onsi/gomega"
	"github.com/lightpaw/male7/entity/cb"
)

func TestRangeBlock(t *testing.T) {
	RegisterTestingT(t)

	block := &BlockData{
		XLen:           100,
		YLen:           100,
		BaseCountLimit: 50,
	}

	sb, err := NewStitchedBlocks(7, 7, 3, 3, block, make(map[cb.Cube]struct{}))
	Ω(err).Should(Succeed())

	Ω(getRingBlockXYs(3, 3, 0, sb.IsValidIntBlock)).Should(Equal([]cb.Cube{cb.XYCube(3, 3)}))

	Ω(getSpiralBlockXYs(3, 3, 1, sb.IsValidIntBlock)).Should(ConsistOf(
		cb.XYCube(2, 2), cb.XYCube(3, 2), cb.XYCube(4, 2),
		cb.XYCube(2, 3), cb.XYCube(3, 3), cb.XYCube(4, 3),
		cb.XYCube(2, 4), cb.XYCube(3, 4), cb.XYCube(4, 4),
	))

	Ω(getSpiralBlockXYs(3, 3, 2, sb.IsValidIntBlock)).Should(ConsistOf(
		cb.XYCube(1, 1), cb.XYCube(2, 1), cb.XYCube(3, 1), cb.XYCube(4, 1), cb.XYCube(5, 1),
		cb.XYCube(1, 2), cb.XYCube(2, 2), cb.XYCube(3, 2), cb.XYCube(4, 2), cb.XYCube(5, 2),
		cb.XYCube(1, 3), cb.XYCube(2, 3), cb.XYCube(3, 3), cb.XYCube(4, 3), cb.XYCube(5, 3),
		cb.XYCube(1, 4), cb.XYCube(2, 4), cb.XYCube(3, 4), cb.XYCube(4, 4), cb.XYCube(5, 4),
		cb.XYCube(1, 5), cb.XYCube(2, 5), cb.XYCube(3, 5), cb.XYCube(4, 5), cb.XYCube(5, 5),
	))

	Ω(sb.GetRound4BlockByPos(0, 0)).Should(Equal([]cb.Cube{
		cb.XYCube(0, 0), cb.XYCube(1, 0), cb.XYCube(0, 1), cb.XYCube(1, 1),
	}))

	Ω(sb.GetRound4BlockByPos(99, 99)).Should(Equal([]cb.Cube{
		cb.XYCube(0, 0), cb.XYCube(1, 0), cb.XYCube(0, 1), cb.XYCube(1, 1),
	}))

	Ω(sb.GetRound4BlockByPos(249, 250)).Should(Equal([]cb.Cube{
		cb.XYCube(2, 2), cb.XYCube(1, 2), cb.XYCube(2, 3), cb.XYCube(1, 3),
	}))

	Ω(sb.GetRound4BlockByPos(600, 600)).Should(Equal([]cb.Cube{
		cb.XYCube(6, 6), cb.XYCube(5, 6), cb.XYCube(6, 5), cb.XYCube(5, 5),
	}))

	Ω(sb.GetRound4BlockByPos(699, 699)).Should(Equal([]cb.Cube{
		cb.XYCube(6, 6), cb.XYCube(5, 6), cb.XYCube(6, 5), cb.XYCube(5, 5),
	}))

}

func TestEdge(t *testing.T) {
	RegisterTestingT(t)
	block := &BlockData{
		XLen:           100,
		YLen:           100,
		BaseCountLimit: 50,
	}

	sb, err := NewStitchedBlocks(5, 4, 3, 3, block, make(map[cb.Cube]struct{}))
	Ω(err).Should(Succeed())

	rb := sb.GetRadiusBlock(0)
	Ω(rb.Radius).Should(Equal(uint64(0)))
	Ω(rb.MinBlockX).Should(Equal(uint64(3)))
	Ω(rb.MinBlockY).Should(Equal(uint64(3)))
	Ω(rb.MaxBlockX).Should(Equal(uint64(3)))
	Ω(rb.MaxBlockY).Should(Equal(uint64(3)))

	Ω(rb.MinX).Should(Equal(uint64(300)))
	Ω(rb.MinY).Should(Equal(uint64(300)))
	Ω(rb.MaxX).Should(Equal(uint64(399)))
	Ω(rb.MaxY).Should(Equal(uint64(399)))

	rb = sb.GetRadiusBlock(1)
	Ω(rb.Radius).Should(Equal(uint64(1)))
	Ω(rb.MinBlockX).Should(Equal(uint64(2)))
	Ω(rb.MinBlockY).Should(Equal(uint64(2)))
	Ω(rb.MaxBlockX).Should(Equal(uint64(4)))
	Ω(rb.MaxBlockY).Should(Equal(uint64(3)))

	Ω(rb.MinX).Should(Equal(uint64(200)))
	Ω(rb.MinY).Should(Equal(uint64(200)))
	Ω(rb.MaxX).Should(Equal(uint64(499)))
	Ω(rb.MaxY).Should(Equal(uint64(399)))

	rb = sb.GetRadiusBlock(2)
	Ω(rb.Radius).Should(Equal(uint64(2)))
	Ω(rb.MinBlockX).Should(Equal(uint64(1)))
	Ω(rb.MinBlockY).Should(Equal(uint64(1)))
	Ω(rb.MaxBlockX).Should(Equal(uint64(4)))
	Ω(rb.MaxBlockY).Should(Equal(uint64(3)))

	Ω(rb.MinX).Should(Equal(uint64(100)))
	Ω(rb.MinY).Should(Equal(uint64(100)))
	Ω(rb.MaxX).Should(Equal(uint64(499)))
	Ω(rb.MaxY).Should(Equal(uint64(399)))

	rb = sb.GetRadiusBlock(3)
	Ω(rb.Radius).Should(Equal(uint64(3)))
	Ω(rb.MinBlockX).Should(Equal(uint64(0)))
	Ω(rb.MinBlockY).Should(Equal(uint64(0)))
	Ω(rb.MaxBlockX).Should(Equal(uint64(4)))
	Ω(rb.MaxBlockY).Should(Equal(uint64(3)))

	Ω(rb.MinX).Should(Equal(uint64(0)))
	Ω(rb.MinY).Should(Equal(uint64(0)))
	Ω(rb.MaxX).Should(Equal(uint64(499)))
	Ω(rb.MaxY).Should(Equal(uint64(399)))

	rb = sb.GetRadiusBlock(10)
	Ω(rb.Radius).Should(Equal(uint64(3)))
	Ω(rb.MinBlockX).Should(Equal(uint64(0)))
	Ω(rb.MinBlockY).Should(Equal(uint64(0)))
	Ω(rb.MaxBlockX).Should(Equal(uint64(4)))
	Ω(rb.MaxBlockY).Should(Equal(uint64(3)))

	Ω(rb.MinX).Should(Equal(uint64(0)))
	Ω(rb.MinY).Should(Equal(uint64(0)))
	Ω(rb.MaxX).Should(Equal(uint64(499)))
	Ω(rb.MaxY).Should(Equal(uint64(399)))
}

func TestRangePos(t *testing.T) {
	RegisterTestingT(t)

	block := &BlockData{
		XLen:           100,
		YLen:           100,
		BaseCountLimit: 50,
	}

	for x := 0; x < int(block.XLen); x++ {
		for y := 0; y < int(block.YLen); y++ {
			block.possibleHomeCubes = append(block.possibleHomeCubes, cb.XYCube(x, y))
		}
	}

	sb, err := NewStitchedBlocks(7, 7, 3, 3, block, make(map[cb.Cube]struct{}))
	Ω(err).Should(Succeed())

	for bx := uint64(0); bx < 7; bx++ {
		for by := uint64(0); by < 7; by++ {
			sb.RangeBlockHomeCubes(bx, by, func(c cb.Cube) (toContinue bool) {

				x, y := c.XY()

				abx, aby := sb.MustBlockByPos(x, y)
				Ω(abx).Should(Equal(bx))
				Ω(aby).Should(Equal(by))
				return true
			})
		}
	}

}
