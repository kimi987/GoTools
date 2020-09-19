package gm

import (
	"fmt"
	"github.com/lightpaw/male7/config/goods"
	"github.com/lightpaw/male7/entity"
	"github.com/lightpaw/male7/gen/iface"
)

func (m *GmModule) newPrintEquipsGmGroup() *gm_group {
	printEquipmentGroup := &gm_group{
		tab: "打印装备",
	}
	for _, data := range m.datas.EquipmentData().Array {
		func(data *goods.EquipmentData) {
			handler := newIntHandler(data.Name, "", func(amount int64, hc iface.HeroController) {

				e := entity.NewEquipment(1, data)
				fmt.Println(e.EncodeClient().String())
				nextLevel := e.LevelData().NextLevel()
				if nextLevel != nil {
					fmt.Println("currentLevelStat", nextLevel.CurrentLevelStat)
					fmt.Println("baseStat", data.BaseStat)
					fmt.Println(data.CalculateUpgradeLevelStat(e.LevelData(), nextLevel))
				}

			})

			printEquipmentGroup.handler = append(printEquipmentGroup.handler, handler)
		}(data)
	}

	return printEquipmentGroup
}
