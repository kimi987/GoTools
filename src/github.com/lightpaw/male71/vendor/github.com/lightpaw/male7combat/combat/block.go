package combat

type BlockInfo interface {
	Walkable(x, y int) bool
	SetWalkable(x, y int)
	SetUnwalkable(x, y int)
}

func NewFastBlock(xlen, ylen int) BlockInfo {

	if xlen < 0 {
		xlen = 0
	}
	if ylen < 0 {
		ylen = 0
	}

	blockCount := xlen * ylen
	return &boolArrayBlockInfo{xLen: xlen, yLen: ylen, blocked: make([]bool, blockCount)}
}

type boolArrayBlockInfo struct {
	xLen    int
	yLen    int
	blocked []bool
}

func (b *boolArrayBlockInfo) index(x, y int) int {
	return x + b.xLen*y
}

func (b *boolArrayBlockInfo) Walkable(x, y int) bool {
	if x < 0 || x >= b.xLen || y < 0 || y >= b.yLen {
		return false
	}

	return !b.blocked[b.index(x, y)]
}

func (b *boolArrayBlockInfo) SetWalkable(x, y int) {
	if x < 0 || x >= b.xLen || y < 0 || y >= b.yLen {
		return
	}

	b.blocked[b.index(x, y)] = false
}

func (b *boolArrayBlockInfo) SetUnwalkable(x, y int) {
	if x < 0 || x >= b.xLen || y < 0 || y >= b.yLen {
		return
	}

	b.blocked[b.index(x, y)] = true
}
