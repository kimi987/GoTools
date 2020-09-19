package gm

import (
	"fmt"
	"github.com/lightpaw/male7/entity"
	"github.com/lightpaw/male7/gen/iface"
	"github.com/lightpaw/male7/gen/pb/gm"
	"github.com/lightpaw/male7/service/conflict/heroservice/herolock"
)

func (m *GmModule) newZhanJiangGmGroup() *gm_group {
	return &gm_group{
		tab: "过关斩将",
		handler: []*gm_handler{
			newHeroIntHandler("通关到X关", "0", m.passZhanJiangToGuanQia),
			newHeroIntHandler("通关到X章", "0", m.passZhanJiangToChapter),
			newHeroStringHandler("全部通关", "", m.passAllZhanJiang),
		},
	}
}

func (m *GmModule) passZhanJiangToGuanQia(amount int64, hero *entity.Hero, result herolock.LockResult, hc iface.HeroController) {
	c := m.datas.GetZhanJiangGuanQiaData(uint64(amount))
	if c == nil {
		hc.Send(gm.NewS2cGmMsg(fmt.Sprintf("关卡数据没找到: %d", amount)))
		return
	}

	heroZhanJiang := hero.ZhanJiang()
	for i := int64(0); i < amount; i++ {
		if c == nil {
			break
		}

		heroZhanJiang.Pass(c)
		result.Add(c.PassMsg)

		c = c.Prev
	}
}

func (m *GmModule) passZhanJiangToChapter(amount int64, hero *entity.Hero, result herolock.LockResult, hc iface.HeroController) {
	c := m.datas.GetZhanJiangChapterData(uint64(amount))
	if c == nil {
		hc.Send(gm.NewS2cGmMsg(fmt.Sprintf("章节数据没找到: %d", amount)))
		return
	}

	heroZhanJiang := hero.ZhanJiang()
	for i := int64(0); i < amount; i++ {
		if c == nil {
			break
		}

		for _, d := range c.ZhanJiangDatas {
			heroZhanJiang.Pass(d)
			result.Add(d.PassMsg)
		}

		c = c.PreChapter
	}
}

func (m *GmModule) passAllZhanJiang(amount string, hero *entity.Hero, result herolock.LockResult, hc iface.HeroController) {
	heroZhanJiang := hero.ZhanJiang()
	for _, c := range m.datas.GetZhanJiangChapterDataArray() {
		for _, d := range c.ZhanJiangDatas {
			heroZhanJiang.Pass(d)
			result.Add(d.PassMsg)
		}
	}
}
