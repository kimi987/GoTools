package bai_zhan

import (
	"github.com/lightpaw/male7/config/bai_zhan_data"
	"github.com/lightpaw/male7/pb/shared_proto"
	"math/rand"
)

func newBaiZhanLevelDatas(levelDatas []*bai_zhan_data.JunXianLevelData) []*bai_zhan_level_data {
	result := make([]*bai_zhan_level_data, len(levelDatas))
	for idx, data := range levelDatas {
		result[idx] = &bai_zhan_level_data{
			levelData:  data,
			mirrorObjs: make([]mirror_obj, 0, 64),
		}
	}
	return result
}

// 等级数据
type bai_zhan_level_data struct {
	levelData  *bai_zhan_data.JunXianLevelData // 军衔等级数据
	mirrorObjs []mirror_obj                    // 该军衔等级的所有镜像
}

func (d *bai_zhan_level_data) addMirror(obj mirror_obj) {
	d.mirrorObjs = append(d.mirrorObjs, obj)
}

func (d *bai_zhan_level_data) clearMirrors() {
	d.mirrorObjs = d.mirrorObjs[:0]
}

// 随机一个目标
func (d *bai_zhan_level_data) randomTarget(heroId int64, challengeTimes uint64) (targetId int64, targetProto *shared_proto.CombatPlayerProto, isDefenderNpc bool) {
	levelData := d.levelData

	strongMatchNpcGuardCaptain := levelData.StrongMatchNpcGuardCaptain
	if challengeTimes < uint64(len(strongMatchNpcGuardCaptain)) {
		// 在强制匹配范围内
		guard := strongMatchNpcGuardCaptain[challengeTimes]
		return guard.GetNpcId(), guard.GetPlayer(), true
	}

	totalTargetCaptainCount := len(levelData.NpcGuardCaptain) + len(d.mirrorObjs)
	randomIndex := rand.Intn(totalTargetCaptainCount)
	if randomIndex < len(levelData.NpcGuardCaptain) {
		guard := levelData.NpcGuardCaptain[randomIndex]
		return guard.GetNpcId(), guard.GetPlayer(), true
	}

	mirrorObj := d.mirrorObjs[randomIndex-len(levelData.NpcGuardCaptain)]
	if mirrorObj.Id() != heroId {
		return mirrorObj.Id(), mirrorObj.CombatMirror(), false
	}

	// 减少一个再来一次
	totalTargetCaptainCount--

	randomIndex = rand.Intn(totalTargetCaptainCount)
	if randomIndex < len(levelData.NpcGuardCaptain) {
		guard := levelData.NpcGuardCaptain[randomIndex]
		return guard.GetNpcId(), guard.GetPlayer(), true
	}

	mirrorObj = d.mirrorObjs[randomIndex-len(levelData.NpcGuardCaptain)]
	if mirrorObj.Id() == heroId {
		// 下一个一定有的
		mirrorObj = d.mirrorObjs[randomIndex-len(levelData.NpcGuardCaptain)+1]
	}

	return mirrorObj.Id(), mirrorObj.CombatMirror(), false
}

type mirror_obj interface {
	Id() int64
	CombatMirror() *shared_proto.CombatPlayerProto
}
