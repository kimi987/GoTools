package regdata

import (
	"github.com/lightpaw/male7/config/blockdata"
	"github.com/lightpaw/male7/constants"
	"github.com/lightpaw/male7/entity/hexagon"
	"github.com/lightpaw/male7/pb/shared_proto"
	"github.com/lightpaw/male7/util/check"
	"github.com/lightpaw/male7/util/u64"
	"github.com/lightpaw/logrus"
	"math"
	"github.com/lightpaw/male7/config/mingcdata"
	"github.com/lightpaw/male7/entity/cb"
	"github.com/lightpaw/male7/config/country"
)

func RegionDataID(regionType shared_proto.RegionType, level uint64) uint64 {
	return uint64(regionType)*10000 + level
}

// 地区分级
//gogen:config
type RegionData struct {
	_ struct{} `file:"地图/地区.txt"`
	_ struct{} `proto:"shared_proto.RegionDataProto"`
	_ struct{} `protoconfig:"region_data"`

	Id uint64 `head:"-,RegionDataID(%s.RegionType%c %s.Level)"`

	// 地区类型，1-蛮荒地区 2-主城地区 3-荣誉地区
	RegionType shared_proto.RegionType

	Level uint64

	BlockXLen    uint64 `protofield:"-"`
	BlockYLen    uint64 `protofield:"-"`
	BlockId      uint64 `head:"block" protofield:"-"`
	CenterBlockX uint64 `validator:"uint" default:"3" protofield:"-"`
	CenterBlockY uint64 `validator:"uint" default:"3" protofield:"-"`
	InitRadius   uint64 `default:"3" protofield:"-"`

	RandomMinRadius uint64 `validator:"uint" default:"3" protofield:"-"`
	RandomMaxRadius uint64 `validator:"uint" default:"3" protofield:"-"`

	randomBlocks []cb.Cube

	Block        *blockdata.StitchedBlocks `head:"-" protofield:"-"`
	SubBlockXLen uint64                    `head:"-"`
	SubBlockYLen uint64                    `head:"-"`

	BornMinRadius uint64 `validator:"uint" default:"3" protofield:"-"`
	BornMaxRadius uint64 `validator:"uint" default:"3" protofield:"-"`

	countryBornBlockMap map[uint64][]cb.Cube
	defaultBornBlock    []cb.Cube

	GuildMoveBaseMinRadius uint64 `default:"1"`
	GuildMoveBaseMaxRadius uint64 `default:"32"`

	// 刷怪列表
	monsters []*RegionMonsterData

	maxMonsterLevel uint64 // 打倒这个等级的怪可以晋级下一个等级的地图

	// 刷怪列表
	multiLevelMonsters []*RegionMultiLevelNpcData

	// 区域带
	Area []*RegionAreaData `head:"-" protofield:"-"`
}

const maxBlockLen = math.MaxUint8

func BlockSequence(x, y uint64) uint64 {
	return x<<8 | y
}

func (d *RegionData) Init2(filename string, configDatas interface {
	GetRegionMonsterDataArray() []*RegionMonsterData
	GetRegionMultiLevelNpcDataArray() []*RegionMultiLevelNpcData
	GetMingcBaseDataArray() []*mingcdata.MingcBaseData
	GetBlockData(uint64) *blockdata.BlockData
	GetCountryDataArray() []*country.CountryData
	GetRegionAreaDataArray() []*RegionAreaData
}) {

	d.Area = configDatas.GetRegionAreaDataArray()

	// 跟Npc的ID有关系，不能随便改
	check.PanicNotTrue(d.BlockXLen <= maxBlockLen, "%s-%s 地区数据，配置的BlockXLen太大，必须 <= 255, 实际是 %d", d.Id, d.RegionType, d.BlockXLen)
	check.PanicNotTrue(d.BlockYLen <= maxBlockLen, "%s-%s 地区数据，配置的BlockYLen太大，必须 <= 255, 实际是 %d", d.Id, d.RegionType, d.BlockYLen)

	block := configDatas.GetBlockData(d.BlockId)
	check.PanicNotTrue(block != nil, "%s-%s 地区数据，配置的Block不存在", d.Id, d.RegionType)
	d.SubBlockXLen = block.XLen
	d.SubBlockYLen = block.YLen

	// 其他系统的阻挡点(不可建城)
	staticBlockMap := make(map[cb.Cube]struct{})
	for _, v := range configDatas.GetMingcBaseDataArray() {
		cubes := hexagon.SpiralRing(u64.Int(v.BaseX), u64.Int(v.BaseY), uint(v.Radius))
		for _, c := range cubes {
			if _, exist := staticBlockMap[c]; exist {
				x, y := c.XY()
				logrus.Panicf("名城%s-%s 野外坐标跟其他的名称存在重叠区域，重叠坐标(%d, %d)", v.Id, v.Name, x, y)
			}
			staticBlockMap[c] = struct{}{}
		}
	}

	stitchedBlock, err := blockdata.NewStitchedBlocks(d.BlockXLen, d.BlockYLen, d.CenterBlockX, d.CenterBlockY, block, staticBlockMap)
	if err != nil {
		logrus.WithError(err).Panicf("%s-%s 地区数据初始化拼接地图块数据失败", d.Id, d.RegionType)
	}
	d.Block = stitchedBlock

	// 初始就是最大的长度
	d.InitRadius = stitchedBlock.MaxRadius()

	var monsters []*RegionMonsterData
	var maxMonsterLevel uint64
	for _, m := range configDatas.GetRegionMonsterDataArray() {
		if m.RegionId == d.Id {
			monsters = append(monsters, m)
			maxMonsterLevel = u64.Max(maxMonsterLevel, m.Base.BaseLevel)

			check.PanicNotTrue(d.Block.IsWalkable(m.BaseX, m.BaseY),
				"%s-%s 地区数据，配置的怪物处于阻挡区域，怪物ID：%d，阻挡点：%d,%d",
				d.Id, d.RegionType, m.Id, m.BaseX, m.BaseY)
		}
	}
	d.monsters = monsters
	d.maxMonsterLevel = maxMonsterLevel

	for _, m1 := range d.monsters {
		for _, m2 := range d.monsters {
			if m1 == m2 {
				continue
			}

			check.PanicNotTrue(hexagon.Distance(m1.BaseX, m1.BaseY, m2.BaseX, m2.BaseY) > constants.BaseConflictRange,
				"%s-%s 地区数据，配置的怪物距离太近，怪物1：%d (%d,%d)，怪物2：%d (%d,%d)",
				d.Id, d.RegionType, m1.Id, m1.BaseX, m1.BaseY, m2.Id, m2.BaseX, m2.BaseY)
		}
	}

	var multiLevelMonsters []*RegionMultiLevelNpcData
	for _, v := range configDatas.GetRegionMultiLevelNpcDataArray() {
		if d.Id == v.RegionId {
			multiLevelMonsters = append(multiLevelMonsters, v)
		}
	}

	d.multiLevelMonsters = multiLevelMonsters

	for _, m1 := range d.multiLevelMonsters {
		for _, m2 := range d.multiLevelMonsters {
			if m1 == m2 {
				continue
			}

			check.PanicNotTrue(hexagon.Distance(int(m1.OffsetBaseX), int(m1.OffsetBaseY), int(m2.OffsetBaseX), int(m2.OffsetBaseY)) > constants.BaseConflictRange,
				"%s-%s 地区数据，配置的多等级怪物距离太近，怪物1：%d (%d,%d)，怪物2：%d (%d,%d)",
				d.Id, d.RegionType, m1.Id, m1.OffsetBaseX, m1.OffsetBaseY, m2.Id, m2.OffsetBaseX, m2.OffsetBaseY)

		}

		for _, m2 := range d.monsters {
			ox, oy := d.Block.GetOffsetByPos(m2.BaseX, m2.BaseY)

			check.PanicNotTrue(hexagon.Distance(int(m1.OffsetBaseX), int(m1.OffsetBaseY), ox, oy) > constants.BaseConflictRange,
				"%s-%s 地区数据，配置的多等级怪物与固定怪物距离太近，多等级怪物1：%d (%d,%d)，固定怪物2：%d (%d,%d)",
				d.Id, d.RegionType, m1.Id, m1.OffsetBaseX, m1.OffsetBaseY, m2.Id, ox, oy)
		}
	}

	minLen := u64.Min(d.Block.XLen, d.Block.YLen)
	check.PanicNotTrue(d.GuildMoveBaseMaxRadius < minLen*2, "联盟随机迁城令的最大半径必须小于地图宽度的一半 %v %v", d.GuildMoveBaseMaxRadius, minLen)
	check.PanicNotTrue(d.GuildMoveBaseMinRadius < d.GuildMoveBaseMaxRadius, "联盟随机迁城令的最小半径必须大于最大半径 %v %v", d.GuildMoveBaseMinRadius, d.GuildMoveBaseMaxRadius)

	check.PanicNotTrue(d.RandomMinRadius <= d.RandomMaxRadius,
		"%s-%s 地区数据，配置的随机迁城半径区域无效，必须d.RandomMinRadius(%d) <= d.RandomMaxRadius(%d)",
		d.Id, d.RegionType, d.RandomMinRadius, d.RandomMaxRadius)
	check.PanicNotTrue(d.RandomMaxRadius < stitchedBlock.MaxRadius(),
		"%s-%s 地区数据，配置的随机迁城半径区域无效，必须d.RandomMaxRadius(%d) < stitchedBlock.MaxRadius(%d)",
		d.Id, d.RegionType, d.RandomMaxRadius, stitchedBlock.MaxRadius())

	var randomBlocks []cb.Cube
	for i := d.RandomMinRadius; i < d.RandomMaxRadius; i++ {
		rb := stitchedBlock.GetRadiusBlock(i)
		if rb.Radius != i {
			break
		}
		randomBlocks = append(randomBlocks, rb.GetRingBlockXYs()...)
	}

	check.PanicNotTrue(len(randomBlocks) > 0,
		"%s-%s 地区数据，len(randomBlocks) == 0",
		d.Id, d.RegionType)

	d.randomBlocks = randomBlocks

	check.PanicNotTrue(d.BornMinRadius <= d.BornMaxRadius,
		"%s-%s 地区数据，配置的随机迁城半径区域无效，必须d.BornMinRadius(%d) <= d.BornMaxRadius(%d)",
		d.Id, d.RegionType, d.BornMinRadius, d.BornMaxRadius)
	check.PanicNotTrue(d.BornMaxRadius < stitchedBlock.MaxRadius(),
		"%s-%s 地区数据，配置的随机迁城半径区域无效，必须d.BornMaxRadius(%d) < stitchedBlock.MaxRadius(%d)",
		d.Id, d.RegionType, d.BornMaxRadius, stitchedBlock.MaxRadius())

	// 国家出生地初始化
	countryValidBlockMap := make(map[uint64]map[cb.Cube]struct{})
	for _, c := range configDatas.GetCountryDataArray() {
		// 在中心点的矩形区域

		func(centerX, centerY, radiusX, radiusY uint64) {
			minX := u64.Sub(centerX, radiusX)
			minY := u64.Sub(centerY, radiusY)
			maxX := u64.Min(centerX+radiusX+1, d.BlockXLen)
			maxY := u64.Min(centerY+radiusY+1, d.BlockYLen)

			validBlockMap := make(map[cb.Cube]struct{})
			for x := minX; x < maxX; x++ {
				for y := minY; y < maxY; y++ {
					validBlockMap[cb.XYCube(int(x), int(y))] = struct{}{}
				}
			}

			countryValidBlockMap[c.Id] = validBlockMap
		}(c.BornCenterX, c.BornCenterY, c.BornRadiusX, c.BornRadiusY)

	}

	ringValidBlockMap := make(map[cb.Cube]struct{})
	// 在内径和外径之间的区域
	for i := d.BornMinRadius; i <= d.BornMaxRadius; i++ {
		rb := stitchedBlock.GetRadiusBlock(i)
		if rb.Radius != i {
			break
		}

		for _, b := range rb.GetRingBlockXYs() {
			ringValidBlockMap[b] = struct{}{}
		}
	}

	countryBornBlockMap := make(map[uint64][]cb.Cube)
	for _, c := range configDatas.GetCountryDataArray() {
		check.PanicNotTrue(stitchedBlock.IsValidBlock(c.BornCenterX, c.BornCenterY),
			"%s-%s 配置的出生点中心位置(%d, %d)",
			c.Id, c.Name, c.BornCenterX, c.BornCenterY)

		selfValidBlockMap := countryValidBlockMap[c.Id]

		isValidIntBlock := func(x, y int) bool {
			xy := cb.XYCube(x, y)

			// 如果在自己国家的中心点区域，那么肯定可以
			if _, exist := selfValidBlockMap[xy]; exist {
				return true
			}

			// 不在环形区域，不行
			if _, exist := ringValidBlockMap[xy]; !exist {
				return false
			}

			// 如果在其他国家地区位置，不行
			for _, blockMap := range countryValidBlockMap {
				if _, exist := blockMap[xy]; exist {
					return false
				}
			}

			// 其他可以
			return true
		}

		maxRadius := u64.Max(
			u64.Max(c.BornCenterX, d.BlockXLen-c.BornCenterX-1),
			u64.Max(c.BornCenterY, d.BlockYLen-c.BornCenterY-1),
		)

		spiralBlockXys := blockdata.GetSpiralBlockXYs(int(c.BornCenterX), int(c.BornCenterY), int(maxRadius), isValidIntBlock)

		countryBornBlockMap[c.Id] = spiralBlockXys

		if len(d.defaultBornBlock) <= 0 {
			d.defaultBornBlock = spiralBlockXys
		}
	}

	d.countryBornBlockMap = countryBornBlockMap
}

func (d *RegionData) GetMonsters() []*RegionMonsterData {
	return d.monsters
}

func (d *RegionData) RangeRandomBlocks(f func(blockX, blockY uint64) (toContinue bool)) {

	random := make([]cb.Cube, len(d.randomBlocks))
	copy(random, d.randomBlocks)
	cb.Mix(random)

	for _, xy := range random {
		x, y := xy.XY()
		if !f(uint64(x), uint64(y)) {
			return
		}
	}
}

func (d *RegionData) GetCountryBornBlock(country uint64) cb.Cubes {
	if info := d.countryBornBlockMap[country]; len(info) > 0 {
		return info
	}
	return d.defaultBornBlock
}

func (d *RegionData) GetMultiLevelMonsters() []*RegionMultiLevelNpcData {
	return d.multiLevelMonsters
}

func (d *RegionData) GetMaxMonsterLevel() uint64 {
	return d.maxMonsterLevel
}

func (d *RegionData) GetAreaByPos(x, y int) *RegionAreaData {
	return d.GetAreaByBlock(d.Block.GetIntBlockByPos(x, y))
}

func (d *RegionData) GetAreaByBlock(bx, by int) *RegionAreaData {
	for _, v := range d.Area {
		if v.Area.IsValidPos(bx, by) {
			return v
		}
	}
	return nil
}

type level_rank []*RegionData

func (a level_rank) Len() int           { return len(a) }
func (a level_rank) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a level_rank) Less(i, j int) bool { return a[i].Level < a[j].Level }
