package blockdata

import (
	"github.com/lightpaw/config"
	"github.com/lightpaw/logrus"
	"github.com/lightpaw/male7/entity/cb"
	"github.com/lightpaw/male7/entity/hexagon"
	"github.com/lightpaw/male7/pb/shared_proto"
	"github.com/lightpaw/male7/util/check"
	"github.com/lightpaw/male7/util/u64"
)

const walkable = 0
const unwalkable = 1
const blockPath = "地图阻挡/"

//gogen:config
type BlockData struct {
	_ struct{} `file:"地图阻挡/阻挡块.txt"`

	Id   uint64
	Name string

	XLen uint64
	YLen uint64

	AutoExpandBaseCount uint64 // 地区扩张所需城池数

	// 拥挤容量，超过这个值表示拥挤，新建城池时候，达到这个值不再进去
	NewHeroCrowdedCapcity uint64

	BaseCountLimit uint64 // 地图块城池个数限制（超过这个值，限制刷怪）

	EdgeNotHomeLen uint64 `validator:"uint" default:"0"`

	ProtoFileName string
	protoBytes    []byte

	block            *uint8_map
	surroundingBlock *uint8_map // block为中心，周围一圈都不能是不能走的

	// 所有可能的建主城的坐标
	possibleHomeCubes cb.Cubes

	// 以中心点为来几个环
	CenterX uint64   `head:"-,%s.XLen/2"`
	CenterY uint64   `head:"-,%s.YLen/2"`
	Radius  []uint64 `default:"10,30,50"`

	centerRingCubes   []cb.Cubes
	centerSpiralCubes []cb.Cube
}

func (b *BlockData) Init(gos *config.GameObjects, filename string) {

	b.block = newUint8Map(b.XLen, b.YLen, unwalkable)
	b.surroundingBlock = newUint8Map(b.XLen, b.YLen, unwalkable)

	protoPath := blockPath + b.ProtoFileName
	b.protoBytes = gos.Bytes(protoPath)
	check.PanicNotTrue(len(b.protoBytes) > 0, "%s 中 %d-%s 配置的地图阻挡块数据文件没找到，路径: %s", filename, b.Id, b.Name, protoPath)

	proto := &shared_proto.BlockInfoProto{}
	err := proto.Unmarshal(b.protoBytes)
	if err != nil {
		logrus.WithError(err).Panic("解析block.proto配置文件失败")
	}

	if len(proto.X) != len(proto.Y) {
		logrus.Panic("解析block.proto配置文件失败，len(proto.X) != len(proto.Y)")
	}

	for i, x32 := range proto.X {
		x := int(x32)
		y := int(proto.Y[i])

		ux, uy := uint64(x), uint64(y)
		if x < 0 || y < 0 || ux >= b.XLen || uy >= b.YLen {
			logrus.WithField("x", x).
				WithField("y", y).
				WithField("xLen", b.XLen).
				WithField("yLen", b.YLen).
				WithField("path", protoPath).
				Panic("解析block.proto配置文件失败，proto中的长度，超出地图设置的最大长度")
		}

		b.block.Set(ux, uy, unwalkable)
		b.surroundingBlock.Set(ux, uy, unwalkable)
		for _, c := range hexagon.Neighbors(x, y) {
			nx, ny := c.XY()
			if nx >= 0 && ny >= 0 {
				b.surroundingBlock.Set(u64.FromInt(nx), u64.FromInt(ny), unwalkable)
			}
		}
	}

	// 边缘一圈，不能是障碍格子
	//for x := uint64(0); x < b.XLen; x++ {
	//	check.PanicNotTrue(b.block.Get(x, 0) == walkable, "解析block.proto配置文件，地图边缘一圈都不能设置阻挡区域, 阻挡格：%d,%d", x, 0)
	//	check.PanicNotTrue(b.block.Get(x, b.YLen-1) == walkable, "解析block.proto配置文件，地图边缘一圈都不能设置阻挡区域, 阻挡格：%d,%d", x, b.YLen-1)
	//}
	//for y := uint64(0); y < b.YLen; y++ {
	//	check.PanicNotTrue(b.block.Get(0, y) == walkable, "解析block.proto配置文件，地图边缘一圈都不能设置阻挡区域, 阻挡格：%d,%d", 0, y)
	//	check.PanicNotTrue(b.block.Get(b.XLen-1, y) == walkable, "解析block.proto配置文件，地图边缘一圈都不能设置阻挡区域, 阻挡格：%d,%d", b.XLen-1, y)
	//}

	check.PanicNotTrue(b.EdgeNotHomeLen < b.XLen, "%s 中 %d-%s 配置的X长度太小，XLen: %d，必须>%d", filename, b.Id, b.Name, b.XLen, b.EdgeNotHomeLen)
	check.PanicNotTrue(b.EdgeNotHomeLen < b.YLen, "%s 中 %d-%s 配置的Y长度太小，YLen: %d，必须>%d", filename, b.Id, b.Name, b.YLen, b.EdgeNotHomeLen)

	// 缓存所有的建
	// 找到所有可建主城点, 然后随机打乱顺序
	var possibleHomeCubes cb.Cubes
	for x := b.EdgeNotHomeLen; x < b.XLen-b.EdgeNotHomeLen; x++ {
		for y := b.EdgeNotHomeLen; y < b.YLen-b.EdgeNotHomeLen; y++ {
			ix, iy := u64.Int(x), u64.Int(y)
			if b.IsValidBasePosition(ix, iy) {
				possibleHomeCubes = append(possibleHomeCubes, cb.XYCube(ix, iy))
			}
		}
	}

	if len(possibleHomeCubes) == 0 {
		logrus.WithField("filename", filename).WithField("path", protoPath).Panic("地图中竟然没有一个可建主城点??")
	}

	possibleHomeCubes.Mix()
	b.possibleHomeCubes = possibleHomeCubes

	// 围绕中心点，可行走位置
	cubes := NewMultiRingCubes(int(b.CenterX), int(b.CenterY), true, b.IsValidHomePosition, b.Radius...)
	check.PanicNotTrue(len(cubes) > 0, "%s 中 %d-%s 没找到中心点环形建城坐标", filename, b.Id, b.Name)
	for i, c := range cubes {
		check.PanicNotTrue(len(c) > 0, "%s 中 %d-%s 中心点环形建城坐标，第%d环没数据", filename, b.Id, b.Name, i+1)
	}

	b.centerRingCubes = cubes

	// 围绕中心点螺旋排序
	centerSpiralCubes := make([]cb.Cube, len(possibleHomeCubes))
	copy(centerSpiralCubes, possibleHomeCubes)
	hexagon.SpiralSort(centerSpiralCubes, int(b.CenterX), int(b.CenterY))

	b.centerSpiralCubes = centerSpiralCubes
}

func (b *BlockData) GetProtoBytes() []byte {
	return b.protoBytes
}

//func (b *BlockData) IsWalkable(x, y int) bool {
//	return b.block.GetInt(x, y) == walkable
//}

func (b *BlockData) isSurroundingWalkable(x, y int) bool {
	return b.surroundingBlock.GetInt(x, y) == walkable
}

func (b *BlockData) IsValidBasePosition(x, y int) bool {
	return b.isSurroundingWalkable(x, y)
}

func (b *BlockData) IsHomeArea(x, y int) bool {
	if x < 0 || y < 0 {
		return false
	}

	ux, uy := uint64(x), uint64(y)
	return b.isHomeArea(ux, uy)
}

func (b *BlockData) isHomeArea(ux, uy uint64) bool {
	return ux >= b.EdgeNotHomeLen && ux+b.EdgeNotHomeLen < b.XLen &&
		uy >= b.EdgeNotHomeLen && uy+b.EdgeNotHomeLen < b.YLen
}

// 这个位置是否可以放老家. 必须是可走且不在地图边缘10格范围内
func (b *BlockData) IsValidHomePosition(x, y int) bool {
	if x < 0 || y < 0 {
		return false
	}

	ux, uy := uint64(x), uint64(y)
	return b.isHomeArea(ux, uy) && b.IsValidBasePosition(x, y)
}

func (b *BlockData) GetPossibleHomeCubes() cb.Cubes {
	return b.possibleHomeCubes
}

func (b *BlockData) GetCenterRingCubes() []cb.Cubes {
	return b.centerRingCubes
}

func (b *BlockData) GetCenterSpiralCubes() []cb.Cube {
	return b.centerSpiralCubes
}
