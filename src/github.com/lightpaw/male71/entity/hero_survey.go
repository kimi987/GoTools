package entity

import (
	"github.com/lightpaw/male7/config"
	"github.com/lightpaw/male7/config/survey"
	"github.com/lightpaw/male7/pb/shared_proto"
	"github.com/lightpaw/logrus"
)

func NewHeroSurvey() *HeroSurvey {
	return &HeroSurvey{}
}

// 问卷调查
type HeroSurvey struct {
	completeSurvey []string // 完成了的问卷调查
}

func (h *HeroSurvey) IsCompleted(data *survey.SurveyData) bool {
	for _, id := range h.completeSurvey {
		if id == data.Id {
			return true
		}
	}
	return false
}

func (h *HeroSurvey) Complete(data *survey.SurveyData) {
	h.completeSurvey = append(h.completeSurvey, data.Id)
}

func (h *HeroSurvey) Encode() *shared_proto.HeroSurveyProto {
	proto := &shared_proto.HeroSurveyProto{}

	proto.CompeleteSurvey = h.completeSurvey

	return proto
}

func (h *HeroSurvey) unmarshal(proto *shared_proto.HeroSurveyProto, datas *config.ConfigDatas) {
	if proto == nil {
		return
	}

	h.completeSurvey = proto.CompeleteSurvey
	for _, id := range h.completeSurvey {
		if datas.GetSurveyData(id) == nil {
			logrus.Errorf("存在问卷调查数据找到不到: %d", id)
		}
	}
}
