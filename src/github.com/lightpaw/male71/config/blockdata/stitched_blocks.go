package blockdata

import (
	"github.com/lightpaw/male7/entity/cb"
	"github.com/pkg/errors"
	"github.com/lightpaw/male7/util/u64"
	"github.com/lightpaw/logrus"
	"math/rand"
	"github.com/lightpaw/pbutil"
	"github.com/lightpaw/male7/gen/pb/region"
)

func NewStitchedBlocks(blockXLen, blockYLen, centerX, centerY uint64, block *BlockData, staticBlockMap map[cb.Cube]struct{}) (*StitchedBlocks, error) {

	if centerX >= blockXLen || centerY >= blockYLen {
		return nil, errors.Errorf("地图块初始中心位置超出最大地图块范围, len(%d,%d) center(%d,%d)", blockXLen, blockYLen, centerX, centerY)
	}

	XLen := block.XLen * blockXLen
	YLen := block.YLen * blockYLen

	if XLen > cb.Max {
		return nil, errors.Errorf("地图块太大，导致XLen超出最大上限(32767)，block.XLen(%d) blockXLen(%d) XLen(%d)  ", block.XLen, blockXLen, XLen)
	}
	if YLen > cb.Max {
		return nil, errors.Errorf("地图块太大，导致YLen超出最大上限(32767)，block.YLen(%d) blockYLen(%d) YLen(%d)  ", block.YLen, blockYLen, YLen)
	}

	sb := &StitchedBlocks{
		blockXLen:      blockXLen,
		blockYLen:      blockYLen,
		blockCount:     blockXLen * blockYLen,
		block:          block,
		XLen:           XLen,
		YLen:           YLen,
		centerX:        centerX,
		centerY:        centerY,
		staticBlockMap: staticBlockMap,
	}

	sb.radiusBlocks = newRadiusBlocks(sb, centerX, centerY, sb.IsValidIntBlock)

	return sb, nil
}

type StitchedBlocks struct {
	blockXLen  uint64
	blockYLen  uint64
	blockCount uint64

	block *BlockData

	// 拼接后最终的大小
	XLen uint64
	YLen uint64

	centerX uint64
	centerY uint64

	staticBlockMap map[cb.Cube]struct{}

	radiusBlocks []*radius_blocks
}

// 出生区域，以(x,y)为中心，
type BornBlockInfo struct {
	sb *StitchedBlocks

	// 内半径（含）
	minRadius uint64

	// 外半径（含）
	maxRadius uint64

	// 中心位置
	centerX uint64

	centerY uint64

	// 矩形半径
	radiusX uint64

	radiusY uint64

	spiralBlockXYs []cb.Cube // 螺旋向外
}

func (bbi *BornBlockInfo) RangeBlock(f func(blockX, blockY uint64) (toContinue bool)) {
	for _, b := range bbi.spiralBlockXYs {
		x, y := b.XY()
		if !f(uint64(x), uint64(y)) {
			break
		}
	}
}

func GetSpiralBlockXYs(centerX, centerY, radius int, isValidIntBlock func(x, y int) bool) (blocks []cb.Cube) {
	return getSpiralBlockXYs(centerX, centerY, radius, isValidIntBlock)
}

func getSpiralBlockXYs(centerX, centerY, radius int, isValidIntBlock func(x, y int) bool) (blocks []cb.Cube) {
	return GetRingBlockXYs(centerX, centerY, 0, radius, isValidIntBlock)
}

func GetRingBlockXYs(centerX, centerY, minRadius, maxRadius int, isValidIntBlock func(x, y int) bool) (blocks []cb.Cube) {
	for i := minRadius; i <= maxRadius; i++ {
		blocks = append(blocks, getRingBlockXYs(centerX, centerY, i, isValidIntBlock)...)
	}
	return
}

func getRingBlockXYs(centerX, centerY, radius int, isValidIntBlock func(x, y int) bool) (blocks []cb.Cube) {

	addValidBlock := func(blocks []cb.Cube, x, y int) []cb.Cube {
		if isValidIntBlock(x, y) {
			blocks = append(blocks, cb.XYCube(x, y))
		}
		return blocks
	}

	if radius <= 0 {
		// 特殊情况
		blocks = addValidBlock(blocks, centerX, centerY)
		return
	}

	// 左右上下各一点
	blocks = addValidBlock(blocks, centerX, centerY+radius)
	blocks = addValidBlock(blocks, centerX, centerY-radius)
	blocks = addValidBlock(blocks, centerX+radius, centerY)
	blocks = addValidBlock(blocks, centerX-radius, centerY)

	// i := range 1 .. r-1，每次8个点
	// (r, i), (-r, i), (r, -i), (-r, -i), (i, r), (i, -r), (-i, r), (-i, -r)
	for i := 1; i < radius; i++ {
		blocks = addValidBlock(blocks, centerX+radius, centerY+i)
		blocks = addValidBlock(blocks, centerX-radius, centerY+i)
		blocks = addValidBlock(blocks, centerX+radius, centerY-i)
		blocks = addValidBlock(blocks, centerX-radius, centerY-i)

		blocks = addValidBlock(blocks, centerX+i, centerY+radius)
		blocks = addValidBlock(blocks, centerX+i, centerY-radius)
		blocks = addValidBlock(blocks, centerX-i, centerY+radius)
		blocks = addValidBlock(blocks, centerX-i, centerY-radius)
	}

	// 4个角，(r, r), (-r, -r), (-r, r), (r, -r)
	blocks = addValidBlock(blocks, centerX+radius, centerY+radius)
	blocks = addValidBlock(blocks, centerX+radius, centerY-radius)
	blocks = addValidBlock(blocks, centerX-radius, centerY+radius)
	blocks = addValidBlock(blocks, centerX-radius, centerY-radius)

	return
}

func newRadiusBlocks(sb *StitchedBlocks, centerX, centerY uint64, isValidIntBlock func(x, y int) bool) []*radius_blocks {

	block := sb.block

	maxRadius := u64.Max(
		u64.Max(centerX, sb.blockXLen-centerX-1),
		u64.Max(centerY, sb.blockYLen-centerY-1),
	)

	radius := make([]*radius_blocks, maxRadius+1)
	for i := uint64(0); i <= maxRadius; i++ {

		ring := getRingBlockXYs(int(centerX), int(centerY), int(i), isValidIntBlock)
		spiral := getSpiralBlockXYs(int(centerX), int(centerY), int(i), isValidIntBlock)

		blockMap := make(map[cb.Cube]struct{})
		for _, c := range spiral {
			blockMap[c] = struct{}{}
		}

		minBlockX := u64.Sub(centerX, i)
		minBlockY := u64.Sub(centerY, i)
		maxBlockX := u64.Min(centerX+i, sb.blockXLen-1)
		maxBlockY := u64.Min(centerY+i, sb.blockYLen-1)

		radius[i] = &radius_blocks{
			Radius: i,

			MinBlockX: minBlockX,
			MinBlockY: minBlockY,
			MaxBlockX: maxBlockX,
			MaxBlockY: maxBlockY,

			MinX: block.XLen * minBlockX,
			MinY: block.YLen * minBlockY,
			MaxX: block.XLen*(maxBlockX+1) - 1,
			MaxY: block.YLen*(maxBlockY+1) - 1,

			ringBlockXYs:        ring,
			spiralBlockXYs:      spiral,
			blockMap:            blockMap,
			AutoExpandBaseCount: block.AutoExpandBaseCount * uint64(len(spiral)),
			updateMapRadiusMsg:  region.NewS2cUpdateMapRadiusMsg(u64.Int32(centerX), u64.Int32(centerY), u64.Int32(i)).Static(),
		}
	}

	return radius
}

type radius_blocks struct {
	// 半径
	Radius uint64

	MinBlockX, MinBlockY uint64
	MaxBlockX, MaxBlockY uint64

	MinX, MinY uint64
	MaxX, MaxY uint64

	ringBlockXYs   []cb.Cube // 这个半径这环的blockXY
	spiralBlockXYs []cb.Cube // 螺旋向外

	blockMap map[cb.Cube]struct{}

	// 最大可承载主城个数，超过这个值，扩展地图
	AutoExpandBaseCount uint64

	updateMapRadiusMsg pbutil.Buffer
}

func (rb *radius_blocks) GetRingBlockXYs() []cb.Cube {
	return rb.ringBlockXYs
}

func (rb *radius_blocks) ContainsBlock(blockX, blockY uint64) bool {
	_, exist := rb.blockMap[cb.XYCube(int(blockX), int(blockY))]
	return exist
}

func (rb *radius_blocks) GetUpdateMapRadiusMsg() pbutil.Buffer {
	return rb.updateMapRadiusMsg
}

func (sb *StitchedBlocks) MaxRadius() uint64 {
	return uint64(len(sb.radiusBlocks))
}

func (sb *StitchedBlocks) BlockData() *BlockData {
	return sb.block
}

type BlockRangeType uint8
type CubeRangeType uint8

const (
	BlockRangeTypeRandom       BlockRangeType = iota
	BlockRangeTypeCenterSpiral

	CubeRangeTypeRandom       CubeRangeType = iota
	CubeRangeTypeCenterSpiral
)

func (sb *StitchedBlocks) IsWalkable(x, y int) bool {
	if x < 0 || y < 0 {
		return false
	}

	ux, uy := uint64(x), uint64(y)
	ux = ux % sb.block.XLen
	uy = uy % sb.block.YLen

	return sb.block.block.Get(ux, uy) == walkable
}

func (sb *StitchedBlocks) isSurroundingWalkable(x, y int) bool {
	ux, uy := uint64(x), uint64(y)
	ux = ux % sb.block.XLen
	uy = uy % sb.block.YLen

	return sb.block.surroundingBlock.Get(ux, uy) == walkable
}

func (b *StitchedBlocks) IsValidBasePosition(x, y int) bool {
	// 地图边缘一圈，不让放
	if x <= 0 || y <= 0 {
		return false
	}

	ux, uy := uint64(x), uint64(y)
	if ux >= b.XLen-1 || uy >= b.YLen-1 {
		return false
	}

	if _, exist := b.staticBlockMap[cb.XYCube(x, y)]; exist {
		return false
	}

	return b.isSurroundingWalkable(x, y)
}

func (b *StitchedBlocks) IsValidHomePosition(x, y int) bool {
	return b.IsValidBasePosition(x, y)
}

func (b *StitchedBlocks) GetRadiusBlock(radius uint64) *radius_blocks {
	n := len(b.radiusBlocks)
	if radius < uint64(n) {
		return b.radiusBlocks[radius]
	}
	return b.radiusBlocks[n-1]
}

func (b *StitchedBlocks) IsValidBlock(x, y uint64) bool {
	return x < b.blockXLen && y < b.blockYLen
}

func (b *StitchedBlocks) IsValidIntBlock(x, y int) bool {
	return x >= 0 && y >= 0 && uint64(x) < b.blockXLen && uint64(y) < b.blockYLen
}

func (b *StitchedBlocks) RandomBlock(radius uint64) (uint64, uint64) {
	xys := b.GetRadiusBlock(radius).spiralBlockXYs
	x, y := xys[rand.Intn(len(xys))].XY()
	return uint64(x), uint64(y)
}

func (b *StitchedBlocks) RangeBlock(radius uint64, blockRangeType BlockRangeType, f func(blockX, blockY uint64) (toContinue bool)) {
	var xys []cb.Cube
	switch blockRangeType {
	case BlockRangeTypeRandom:
		spiral := b.GetRadiusBlock(radius).spiralBlockXYs

		random := make([]cb.Cube, len(spiral))
		copy(random, spiral)
		cb.Mix(random)

		xys = random

	case BlockRangeTypeCenterSpiral:
		xys = b.GetRadiusBlock(radius).spiralBlockXYs
	default:
		logrus.Error("StitchedBlocks.RangeCubeBlock unkown BlockRangeType")
	}

	for _, xy := range xys {
		x, y := xy.XY()
		if !f(uint64(x), uint64(y)) {
			return
		}
	}
}

func (b *StitchedBlocks) GetOffsetByPos(x, y int) (offsetX, offsetY int) {
	return int(uint64(x) % b.block.XLen), int(uint64(y) % b.block.YLen)
}

func (b *StitchedBlocks) GetBlockByPos(x, y int) (blockX, blockY uint64) {
	return uint64(x) / b.block.XLen, uint64(y) / b.block.YLen
}

func (b *StitchedBlocks) GetIntBlockByPos(x, y int) (blockX, blockY int) {
	return x / int(b.block.XLen), y / int(b.block.YLen)
}

func (b *StitchedBlocks) MustBlockByPos(x, y int) (blockX, blockY uint64) {
	blockX, blockY = b.GetBlockByPos(x, y)
	return u64.Min(blockX, b.blockXLen-1), u64.Min(blockY, b.blockYLen-1)
}

func (b *StitchedBlocks) OffsetXY(blockX, blockY uint64, x, y int) (int, int) {
	offsetX := int(blockX * b.block.XLen)
	offsetY := int(blockY * b.block.YLen)

	return x + int(offsetX), y + int(offsetY)
}

func (b *StitchedBlocks) RandomHomeXY(blockX, blockY uint64) (int, int) {
	x, y := b.block.GetPossibleHomeCubes().Random().XY()
	return b.OffsetXY(blockX, blockY, x, y)
}

func (b *StitchedBlocks) RangeBlockHomeCubes(blockX, blockY uint64, f func(c cb.Cube) (toContinue bool)) (toContinue bool) {

	offsetX := int(blockX * b.block.XLen)
	offsetY := int(blockY * b.block.YLen)

	return b.block.GetPossibleHomeCubes().RandomRange(func(c cb.Cube) (toContinue bool) {
		return f(cb.XYCube(c.AddXY(offsetX, offsetY)))
	})

}

func (b *StitchedBlocks) GetRound4BlockByPos(x, y int) (blockPos []cb.Cube) {
	ux, uy := uint64(x), uint64(y)

	curX, curY := ux/b.block.XLen, uy/b.block.YLen
	blockPos = append(blockPos, cb.XYCube(int(curX), int(curY)))

	// x坐标上相邻的块
	var isLessThanHalfX bool
	if curX <= 0 {
		isLessThanHalfX = false
	} else if curX >= b.blockXLen-1 {
		isLessThanHalfX = true
	} else {
		isLessThanHalfX = (ux%b.block.XLen)*2 < b.block.XLen
	}

	// y坐标上相邻的块
	var isLessThanHalfY bool
	if curY <= 0 {
		isLessThanHalfY = false
	} else if curY >= b.blockYLen-1 {
		isLessThanHalfY = true
	} else {
		isLessThanHalfY = (uy%b.block.YLen)*2 < b.block.YLen
	}

	// xy坐标对应的块

	var diffX, diffY uint64
	if isLessThanHalfX {
		diffX = u64.Sub(curX, 1)
	} else {
		diffX = curX + 1
	}

	if isLessThanHalfY {
		diffY = u64.Sub(curY, 1)
	} else {
		diffY = curY + 1
	}

	blockPos = append(blockPos, cb.XYCube(int(diffX), int(curY)))
	blockPos = append(blockPos, cb.XYCube(int(curX), int(diffY)))
	blockPos = append(blockPos, cb.XYCube(int(diffX), int(diffY)))

	return
}
