package bai_zhan_service

import (
	"github.com/lightpaw/male7/config/bai_zhan_data"
	"github.com/lightpaw/male7/gen/iface"
	"github.com/lightpaw/male7/module/bai_zhan/bai_zhan_objs"
	"github.com/lightpaw/male7/service/heromodule"
	"github.com/lightpaw/male7/util/event"
	"time"
)

func NewBaiZhanService() *BaiZhanService {
	s := &BaiZhanService{}

	s.baiZhanObjs = bai_zhan_objs.NewBaiZhanObjs()
	s.eventQueue = event.NewEventQueue(1024, time.Second*3, "BaizhanEvent")

	heromodule.RegisterHeroOnlineListener(s)

	return s
}

//gogen:iface
type BaiZhanService struct {
	baiZhanObjs *bai_zhan_objs.BaiZhanObjs
	eventQueue  *event.EventQueue
}

func (s *BaiZhanService) OnHeroOnline(hc iface.HeroController) {
	levelData := s.getHistoryMaxJunXianLevelData(hc.Id())
	if levelData != nil {
		hc.Send(levelData.MaxJunXianLevelChangedMsg)
	}
}

func (s *BaiZhanService) GetBaiZhanObj(id int64) bai_zhan_objs.RHeroBaiZhanObj {
	obj := s.baiZhanObjs.GetBaiZhanObj(id)
	if obj == nil {
		return nil
	}

	return obj
}

func (s *BaiZhanService) GetJunXianLevel(id int64) uint64 {
	obj := s.GetBaiZhanObj(id)
	if obj == nil {
		return 0
	}

	return obj.LevelData().Level
}

func (s *BaiZhanService) GetHistoryMaxJunXianLevel(id int64) uint64 {
	levelData := s.getHistoryMaxJunXianLevelData(id)
	if levelData == nil {
		return 0
	}

	return levelData.Level
}

func (s *BaiZhanService) getHistoryMaxJunXianLevelData(id int64) *bai_zhan_data.JunXianLevelData {
	obj := s.GetBaiZhanObj(id)
	if obj == nil {
		return nil
	}

	return obj.HistoryMaxJunXianLevelData()
}

func (s *BaiZhanService) GetPoint(id int64) uint64 {
	obj := s.GetBaiZhanObj(id)
	if obj == nil {
		return 0
	}

	return obj.Point()
}

func (s *BaiZhanService) TimeOutFunc(f bai_zhan_objs.BaiZhanObjsFunc) bool {
	return s.eventQueue.TimeoutFunc(true, func() {
		f(s.baiZhanObjs)
	})
}

func (s *BaiZhanService) Func(f bai_zhan_objs.BaiZhanObjsFunc) {
	s.eventQueue.Func(true, func() {
		f(s.baiZhanObjs)
	})
}

func (s *BaiZhanService) Stop(saveFuc bai_zhan_objs.BaiZhanObjsFunc) {
	s.eventQueue.Stop()
	saveFuc(s.baiZhanObjs)
}
