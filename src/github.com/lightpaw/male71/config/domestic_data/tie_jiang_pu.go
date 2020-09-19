package domestic_data

import (
	"github.com/lightpaw/male7/config/goods"
	"github.com/lightpaw/male7/util/check"
	"github.com/lightpaw/male7/util/u64"
	"time"
	"github.com/lightpaw/male7/util/imath"
	"github.com/lightpaw/pbutil"
	"github.com/lightpaw/male7/gen/pb/domestic"
)

//gogen:config
type TieJiangPuLevelData struct {
	_     struct{} `file:"内政/铁匠铺等级.txt"`
	_     struct{} `protoconfig:"tie_jiang_pu_level"`
	_     struct{} `proto:"shared_proto.TieJiangPuLevelProto"`
	Level uint64   `validator:"int>0"`

	// 铁匠铺
	MaxForgingTimes          uint64                 `validator:"int>0"`                                                                                                                                   // 最大能够的锻造次数
	RecoveryForgingDuration  time.Duration                                                                                                                                                                // 锻造次数恢复间隔 XXX时间恢复一次
	CanForgingEquipPos       []uint64               `validator:"int>0,notAllNil" head:"can_forging_equip_pos" protofield:"CanForgingEquipPos`                                                             // 能够锻造的装备的位置
	CanForgingEquip          []*goods.EquipmentData `validator:"int>0,duplicate,notAllNil" head:"can_forging_equip_id" protofield:"CanForgingEquip,config.U64a2I32a(goods.GetEquipmentDataKeyArray(%s))"` // 能够锻造的装备
	LockedCanForgingEquipPos []uint64               `head:"-" protofield:"LockedCanForgingEquipPos`                                                                                                       // 未解锁的能够锻造的装备的位置
	LockedCanForgingEquip    []*goods.EquipmentData `head:"-" protofield:"LockedCanForgingEquip,config.U64a2I32a(goods.GetEquipmentDataKeyArray(%s))"`                                                    // 未解锁的能够锻造的装备
	LockedEquipNeedLevel     []uint64               `head:"-" protofield:"LockedEquipNeedLevel`                                                                                                           // 未解锁的装备对应需要的铁匠铺的等级
	CanOneKeyForging         bool                                                                                                                                                                         // 能否一键锻造

	newForgingEquipPos    []uint64
	newForgingEquipPosMsg pbutil.Buffer
}

func (e *TieJiangPuLevelData) GetNewForgingEquipPos() []uint64 {
	return e.newForgingEquipPos
}

func (e *TieJiangPuLevelData) GetNewForgingEquipPosMsg() pbutil.Buffer {
	return e.newForgingEquipPosMsg
}

func (e *TieJiangPuLevelData) GetForgingEquip(idx uint64) *goods.EquipmentData {
	if idx < 0 || idx >= uint64(len(e.CanForgingEquip)) {
		return nil
	}

	return e.CanForgingEquip[idx]
}

func (d *TieJiangPuLevelData) Init(filename string, configDatas interface {
	GetTieJiangPuLevelDataArray() []*TieJiangPuLevelData
}) {
	check.PanicNotTrue(d.RecoveryForgingDuration > 0, "铁匠铺等级数据 %d，恢复一次的时间[%v]必须 > 0 并且 <= 1", d.Level, d.RecoveryForgingDuration)

	check.PanicNotTrue(len(d.CanForgingEquip) > 0, "铁匠铺等级数据 %d，配置的能够打造的装备数量[%v]必须 > 0!", d.Level, len(d.CanForgingEquip))
	check.PanicNotTrue(len(d.CanForgingEquip) == len(d.CanForgingEquipPos), "铁匠铺等级数据 %d，配置的装备位置数量必须跟能够打造的装备数量一致[%v]必须 > 0!", d.Level, len(d.CanForgingEquip), len(d.CanForgingEquipPos))

	array := configDatas.GetTieJiangPuLevelDataArray()
	for level := d.Level + 1; level <= uint64(len(array)); level++ {
		nextLevel := array[level-1]
		if nextLevel == nil {
			break
		}

		// 在下一级中，有我本级解锁跟未解锁的位置

		for idx, pos := range nextLevel.CanForgingEquipPos {
			if u64.Contain(d.CanForgingEquipPos, pos) {
				// 有了
				continue
			}

			if u64.Contain(d.LockedCanForgingEquipPos, pos) {
				// 有了
				continue
			}

			d.LockedCanForgingEquipPos = append(d.LockedCanForgingEquipPos, pos)
			d.LockedCanForgingEquip = append(d.LockedCanForgingEquip, nextLevel.CanForgingEquip[idx])
			d.LockedEquipNeedLevel = append(d.LockedEquipNeedLevel, nextLevel.Level)
		}
	}

	if d.Level-1 > 0 {
		prevLevel := d.Level - 1
		prevLevelData := array[prevLevel-1]
		if prevLevelData.CanOneKeyForging {
			check.PanicNotTrue(d.CanOneKeyForging, "铁匠铺等级数据 %d，前置等级[%d]能够一键锻造，但是下一级[%d]不可以一键锻造", d.Level, prevLevelData.Level, d.Level)
		}

		n := imath.Min(len(prevLevelData.CanForgingEquipPos), len(prevLevelData.CanForgingEquip))

	out:
		for i, e := range d.CanForgingEquip {
			pos := d.CanForgingEquipPos[i]
			for i := 0; i < n; i++ {

				prevPos := prevLevelData.CanForgingEquipPos[i]
				if prevPos == pos {
					if e == prevLevelData.CanForgingEquip[i] {
						continue out
					}
				}
			}

			// 是个新位置
			d.newForgingEquipPos = append(d.newForgingEquipPos, pos)
		}

		if len(d.newForgingEquipPos) > 0 {
			d.newForgingEquipPosMsg = domestic.NewS2cUpdateNewForgingPosMsg(u64.Int32Array(d.newForgingEquipPos)).Static()
		}
	}
}
