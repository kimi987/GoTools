package random_event

import (
	"github.com/lightpaw/male7/pb/shared_proto"
	"github.com/lightpaw/male7/config/resdata"
	"github.com/lightpaw/male7/util/weight"
	"github.com/lightpaw/male7/util/check"
	"github.com/lightpaw/logrus"
	"math/rand"
	"github.com/lightpaw/male7/entity/cb"
	"github.com/lightpaw/male7/config/blockdata"
	"github.com/lightpaw/male7/config/regdata"
)

// 随机事件

const(
	AllArea = 0
)

//gogen:config
type RandomEventDataDictionary struct {
	_ struct{} `singleton:"true"`
	_ struct{} `file:"随机事件/随机事件.txt"`

	// <季节，<地形，事件数据列表>>
	dictionary map[shared_proto.Season] (map[int] []*RandomEventData)
}

func (d *RandomEventDataDictionary) Init(filename string, configData interface{
	GetRandomEventDataArray() []*RandomEventData
	GetEventPositionArray() []*EventPosition
})  {
	areas := make(map[int] int)
	for _, data := range configData.GetEventPositionArray() {
		areas[data.TypeArea] = 0
	}
	
	d.dictionary = make(map[shared_proto.Season] (map[int] []*RandomEventData))
	d.dictionary[shared_proto.Season_SPRING] = make(map[int] []*RandomEventData)
	d.dictionary[shared_proto.Season_SUMMER] = make(map[int] []*RandomEventData)
	d.dictionary[shared_proto.Season_AUTUMN] = make(map[int] []*RandomEventData)
	d.dictionary[shared_proto.Season_WINTER] = make(map[int] []*RandomEventData)
	for area, _ := range areas {
		d.dictionary[shared_proto.Season_SPRING][area] = []*RandomEventData{}
		d.dictionary[shared_proto.Season_SUMMER][area] = []*RandomEventData{}
		d.dictionary[shared_proto.Season_AUTUMN][area] = []*RandomEventData{}
		d.dictionary[shared_proto.Season_WINTER][area] = []*RandomEventData{}
	}

	for _, data := range configData.GetRandomEventDataArray() {
		switch data.TypeSeason {
		case shared_proto.Season_SPRING, shared_proto.Season_SUMMER, shared_proto.Season_AUTUMN, shared_proto.Season_WINTER:
			if data.TypeArea == AllArea {
				for area, _ := range areas {
					d.dictionary[data.TypeSeason][area] = append(d.dictionary[data.TypeSeason][area], data)
				}
			} else {
				_, ok := areas[data.TypeArea]
				check.PanicNotTrue(ok, "%s 随机事件 %d 配置了 事件坐标.txt 中不存在的地形 %d", filename, data.Id, data.TypeArea)

				d.dictionary[data.TypeSeason][data.TypeArea] = append(d.dictionary[data.TypeSeason][data.TypeArea], data)
			}
		default:
			if data.TypeArea == AllArea {
				for area, _ := range areas {
					d.dictionary[shared_proto.Season_SPRING][area] = append(d.dictionary[shared_proto.Season_SPRING][area], data)
					d.dictionary[shared_proto.Season_SUMMER][area] = append(d.dictionary[shared_proto.Season_SUMMER][area], data)
					d.dictionary[shared_proto.Season_AUTUMN][area] = append(d.dictionary[shared_proto.Season_AUTUMN][area], data)
					d.dictionary[shared_proto.Season_WINTER][area] = append(d.dictionary[shared_proto.Season_WINTER][area], data)
				}
			} else{
				_, ok := areas[data.TypeArea]
				check.PanicNotTrue(ok, "%s 随机事件 %d 配置了 事件坐标.txt 中不存在的地形 %d", filename, data.Id, data.TypeArea)

				d.dictionary[shared_proto.Season_SPRING][data.TypeArea] = append(d.dictionary[shared_proto.Season_SPRING][data.TypeArea], data)
				d.dictionary[shared_proto.Season_SUMMER][data.TypeArea] = append(d.dictionary[shared_proto.Season_SUMMER][data.TypeArea], data)
				d.dictionary[shared_proto.Season_AUTUMN][data.TypeArea] = append(d.dictionary[shared_proto.Season_AUTUMN][data.TypeArea], data)
				d.dictionary[shared_proto.Season_WINTER][data.TypeArea] = append(d.dictionary[shared_proto.Season_WINTER][data.TypeArea], data)
			}
		}
	}

	for area, _ := range areas {
		if area == 3 { // 名城地形配置暂不处理
			continue
		}
		check.PanicNotTrue(len(d.dictionary[shared_proto.Season_SPRING][area]) > 0, "%s 随机事件 春季 地形 %d 没有数据!", filename, area)
		check.PanicNotTrue(len(d.dictionary[shared_proto.Season_SUMMER][area]) > 0, "%s 随机事件 夏季 地形 %d 没有数据!", filename, area)
		check.PanicNotTrue(len(d.dictionary[shared_proto.Season_AUTUMN][area]) > 0, "%s 随机事件 秋季 地形 %d 没有数据!", filename, area)
		check.PanicNotTrue(len(d.dictionary[shared_proto.Season_WINTER][area]) > 0, "%s 随机事件 冬季 地形 %d 没有数据!", filename, area)
	}
}

func (d *RandomEventDataDictionary) CatchEventData(typeSeason shared_proto.Season, typeArea int) *RandomEventData {
	datas := d.dictionary[typeSeason][typeArea]
	return datas[rand.Intn(len(datas))]
}

//gogen:config
type RandomEventData struct {
	_ struct{} `file:"随机事件/随机事件.txt"`
	_ struct{} `proto:"shared_proto.RandomEventDataProto"`
	_ struct{} `protoconfig:"random_event"`

	Id                    uint64                    `validator:"int>0"` // 事件ID
	Title                 string                    `validator:"string>0"`
	Desc                  string                    `validator:"string>0"`
	Content               string
	Image                 string
	TypeSeason            shared_proto.Season       `validator:"int" default:"shared_proto.Season_InvalidSeason" protofield:"-"`
	//TypeArea              shared_proto.TypeArea     `validator:"int>0" protofield:"-"`
	TypeArea              int                       `validator:"int" default:"0" protofield:"-"`

	//shared_proto.EventOptionProto
	OptionDatas           []*EventOptionData         `head:"option" protofield:"-"`
	OptionsProto4Send []*shared_proto.EventOptionProto `head:"-" protofield:"-"`
}

func (d *RandomEventData) Init(filename string) {
	check.PanicNotTrue(len(d.OptionDatas) > 1, "%s 随机事件 %d 少于两个选项!", filename, d.Id)
	d.OptionsProto4Send = make([]*shared_proto.EventOptionProto, len(d.OptionDatas))
	for i, option := range d.OptionDatas {
		d.OptionsProto4Send[i] = option.encode4Init()
	}
}

//gogen:config
type EventOptionData struct {
	_ struct{} `file:"随机事件/选项.txt"`
	_ struct{} `proto:"shared_proto.EventOptionProto"`

	Id                    uint64                `validator:"int>0" protofield:"-"`
	Content               string                `validator:"string>0" protofield:"OptionText"`
	Success               string                `protofield:"SuccessText"`
	Failed                string                `protofield:"FailedText"`
	Cost                  *resdata.Cost         `default:"nullable" protofield:"OptionCost"` // 选项消耗，对应cost.txt

	FailedRate            uint64                `validator:"uint" default:"0" protofield:"-"` // 失敗概率

	SuccessPrize          *OptionPrize          `default:"nullable" protofield:"-"`
	FailedPrize           *OptionPrize          `default:"nullable" protofield:"-"`
}

func (d *EventOptionData) encode4Init() *shared_proto.EventOptionProto {
	var i interface{}
	i = d
	m, ok := i.(interface{
		Encode() *shared_proto.EventOptionProto
	})
	if !ok {
		logrus.Errorf("EventOptionData.Encode4Init() cast type fail")
	}
	return m.Encode()
}

//gogen:config
type OptionPrize struct {
	_ struct{} `file:"随机事件/选项奖励.txt"`

	Id         uint64                  `validator:"int>0"`
	GfAdd      uint64                  `validator:"uint"` // 每一级官府增加值，只有资源或君主经验生效果
	Prize      []*resdata.Prize
	Weight     []uint64

	randomer   *weight.WeightRandomer  `head:"-" protofield:"-"`
}

func (p *OptionPrize) Init(filename string) {
	prizeLen := len(p.Prize)
	check.PanicNotTrue(prizeLen > 0, "%s 选项奖励 %d 至少需要配一个奖励!", filename, p.Id)
	if prizeLen > 1 {
		check.PanicNotTrue(prizeLen == len(p.Weight), "%s 选项奖励 %d 权重数量与奖励数量不匹配!", filename, p.Id)

		weights := make([]uint64, prizeLen, prizeLen)
		copy(weights, p.Weight)

		if w, err := weight.NewWeightRandomer(weights); err != nil {
			logrus.WithError(err).Panicf("生成随机事件选项奖励权重机器失败")
		} else {
			p.randomer = w
		}
	}
}

func (p *OptionPrize) CatchPrize() *resdata.Prize {
	if p.randomer == nil {
		return p.Prize[0]
	}
	return p.Prize[p.randomer.RandomIndex()]
}

//gogen:config
type EventPosition struct {
	_ struct{} `file:"随机事件/事件坐标.txt"`

	Id         uint64       `validator:"int>0"`
	PosX       int // 相对坐标，不是绝对坐标
	PosY       int
	TypeArea   int          `validator:"int>0"`
}

//gogen:config
type RandomEventPositionDictionary struct {
	_ struct{} `singleton:"true"`
	_ struct{} `file:"随机事件/事件坐标.txt"`

	dictionary  map[cb.Cube] int // <绝对坐标cb, TypeArea>
	array       []cb.Cube // 数组（随机取事件用）
	blockArray  [][]map[cb.Cube] int // 区块2维数组（用于检测老家区块事件数量）

	block       *blockdata.StitchedBlocks
}

func (d *RandomEventPositionDictionary) Init(filename string, configData interface {
	GetEventPositionArray() []*EventPosition
	GetRegionDataArray() []*regdata.RegionData
})  {
	d.dictionary = make(map[cb.Cube] int)

	regionData := configData.GetRegionDataArray()[0]
	d.block = regionData.Block
	d.blockArray = make([][]map[cb.Cube] int, regionData.BlockXLen)

	for blockX := uint64(0); blockX < regionData.BlockXLen; blockX++ {

		d.blockArray[blockX] = make([]map[cb.Cube] int, regionData.BlockYLen)

		for blockY := uint64(0); blockY < regionData.BlockYLen; blockY++ {

			d.blockArray[blockX][blockY] = make(map[cb.Cube] int)

			for _, eventPos := range configData.GetEventPositionArray() {

				x, y := regionData.Block.OffsetXY(blockX, blockY, eventPos.PosX, eventPos.PosY)
				cube := cb.XYCube(x, y)
				d.array = append(d.array, cube)
				d.dictionary[cube] = eventPos.TypeArea

				d.blockArray[blockX][blockY][cube] = eventPos.TypeArea
			}
		}
	}
}

func (d *RandomEventPositionDictionary) GetTypeArea(cube cb.Cube) int {
	if area, ok := d.dictionary[cube]; ok {
		return area
	}
	return 0 // 返回0 说明地形坐标不存在
}

func (d *RandomEventPositionDictionary) GetRandomPositions(n int) []cb.Cube {
	dLen := len(d.array)
	list := make([]cb.Cube, dLen)
	copy(list, d.array)
	if n < dLen {
		cb.Mix(list)
		return list[:n]
	}
	return list
}

// 区块补足检测
func (d *RandomEventPositionDictionary) CheckAndCatch4Block(list []cb.Cube, n int, baseX, baseY int) []cb.Cube {

	blockX, blockY := d.block.MustBlockByPos(baseX, baseY)

	passedCubes := make(map[cb.Cube] int)
	for _, cube := range list {
		x, y := cb.CubeXY(cube)
		bX, bY := d.block.GetBlockByPos(x, y)
		if bX == blockX && bY == blockY {
			passedCubes[cube] = 0
		}
	}

	m := d.blockArray[blockX][blockY]
	if n > len(m) {
		n = len(m)
	}
	if cubesLen := len(passedCubes); cubesLen < n {
		randList := []cb.Cube{}
		for cube, _ := range m {
			if len(passedCubes) <= 0 {
				randList = append(randList, cube)
			} else if _, ok := passedCubes[cube]; !ok {
				randList = append(randList, cube)
			} else {
				delete(passedCubes, cube)
			}
		}
		cb.Mix(randList)
		return randList[:n - cubesLen]
	}
	return nil
}

// 在除外的点位置随机
func (d *RandomEventPositionDictionary) GetRandomPositionsWithout(excludeList []cb.Cube, n int) []cb.Cube {
	excludeLen := len(excludeList)
	if excludeLen <= 0 {
		return d.GetRandomPositions(n)
	}
	excludeMap := make(map[cb.Cube] int, excludeLen)
	for _, cube := range excludeList {
		excludeMap[cube] = 0
	}
	list := []cb.Cube{}
	for _, cube := range d.array {
		if len(excludeMap) <= 0 {
			list = append(list, cube)
		} else if _, ok := excludeMap[cube]; !ok {
			list = append(list, cube)
		} else {
			delete(excludeMap, cube)
		}
	}
	if n < len(list) {
		cb.Mix(list)
		return list[:n]
	}
	return list
}

// 从传入的点中拆分出满足区块的点和不满足的
func (d *RandomEventPositionDictionary) SelectPositions4Block(list []cb.Cube, baseX, baseY int) (satisfyList []cb.Cube) {

	blockX, blockY := d.block.MustBlockByPos(baseX, baseY)

	for _, cube := range list {
		x, y := cb.CubeXY(cube)
		bX, bY := d.block.GetBlockByPos(x, y)
		if bX == blockX && bY == blockY {
			satisfyList = append(satisfyList, cube)
		}
	}

	return
}
