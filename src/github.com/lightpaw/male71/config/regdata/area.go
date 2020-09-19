package regdata

import (
	"github.com/lightpaw/male7/entity/cb"
	"github.com/lightpaw/male7/util/check"
	"github.com/lightpaw/male7/util/imath"
)

//gogen:config
type AreaData struct {
	_ struct{} `file:"地图/区块链.txt"` // 起一个霸气的名字
	_ struct{} `protogen:"true"`
	_ struct{} `protoconfig:"-"`

	Id uint64 `protofield:"-"`

	// 根据中心位置设置环，最小半径为0，表示矩形或者圆
	CenterX   int `validator:"int>=0"`
	CenterY   int `validator:"int>=0"`
	MinRadius int `validator:"int>=0"`
	MaxRadius int `validator:"int>=0"`

	// 包含
	IncludeX []int `protofield:"-"`
	IncludeY []int `protofield:"-"`

	// 不包含
	ExcludeX []int `protofield:"-"`
	ExcludeY []int `protofield:"-"`

	isInit       bool
	validCubeMap map[cb.Cube]struct{}
}

func (d *AreaData) Init(filename string) {

	check.PanicNotTrue(len(d.IncludeX) == len(d.IncludeY), "%s 配置的区域[%d]的IncludeX和IncludeY的长度不一样", filename, d.Id)
	check.PanicNotTrue(len(d.ExcludeX) == len(d.ExcludeY), "%s 配置的区域[%d]的ExcludeX和ExcludeY的长度不一样", filename, d.Id)

	d.GetValidCubeMap()
}

func (d *AreaData) IsValidPos(x, y int) bool {
	return d.IsValidCube(cb.XYCube(x, y))
}

func (d *AreaData) IsValidCube(c cb.Cube) bool {
	m := d.GetValidCubeMap()
	_, exist := m[c]
	return exist
}

func (d *AreaData) GetValidCubeMap() map[cb.Cube]struct{} {
	if !d.isInit {
		d.doInit()
	}

	return d.validCubeMap
}

func (d *AreaData) doInit() {
	d.isInit = true

	validCubeMap := make(map[cb.Cube]struct{})
	if d.CenterX|d.CenterY|d.MinRadius|d.MaxRadius != 0 {
		for x := d.CenterX - d.MaxRadius; x <= d.CenterX+d.MaxRadius; x ++ {
			isValidX := imath.Abs(x-d.CenterX) >= d.MinRadius
			for y := d.CenterY - d.MaxRadius; y <= d.CenterY+d.MaxRadius; y ++ {
				isValidY := imath.Abs(y-d.CenterY) >= d.MinRadius
				if isValidX || isValidY {
					validCubeMap[cb.XYCube(x, y)] = struct{}{}
				}
			}
		}
	}

	if n := imath.Min(len(d.IncludeX), len(d.IncludeY)); n > 0 {
		for i, x := range d.IncludeX {
			y := d.IncludeY[i]
			validCubeMap[cb.XYCube(x, y)] = struct{}{}
		}
	}

	if n := imath.Min(len(d.ExcludeX), len(d.ExcludeY)); n > 0 {
		for i, x := range d.ExcludeX {
			y := d.ExcludeY[i]
			delete(validCubeMap, cb.XYCube(x, y))
		}
	}

	d.validCubeMap = validCubeMap
}
