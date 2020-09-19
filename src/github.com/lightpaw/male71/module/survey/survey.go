package survey

import (
	"github.com/lightpaw/logrus"
	"github.com/lightpaw/male7/entity"
	"github.com/lightpaw/male7/gen/iface"
	"github.com/lightpaw/male7/pb/shared_proto"
	"github.com/lightpaw/male7/service/conflict/heroservice/herolock"
	"github.com/lightpaw/rpc7"
	"github.com/lightpaw/male7/pb/rpcpb/login2game"
	"github.com/lightpaw/male7/gen/pb/survey"
)

func NewSurveyModule(datas iface.ConfigDatas, heroDataService iface.HeroDataService,
	timeService iface.TimeService, guildService iface.GuildService, mailModule iface.MailModule) *SurveyModule {
	m := &SurveyModule{
		datas:           datas,
		heroDataService: heroDataService,
		timeService:     timeService,
		guildService:    guildService,
		mailModule:      mailModule,
	}

	rpc7.Handle(login2game.NewCompleteQuestionnaireHandler(m.handleGiveSurveyPrize))

	return m
}

//gogen:iface
type SurveyModule struct {
	datas           iface.ConfigDatas
	heroDataService iface.HeroDataService
	timeService     iface.TimeService
	guildService    iface.GuildService
	mailModule      iface.MailModule
}

//gogen:iface
func (m *SurveyModule) ProcessComplete(proto *survey.C2SCompleteProto, hc iface.HeroController) {
	// nothing
}

func (m *SurveyModule) GmGiveSurveyPrize(heroId int64, surveyId string) {
	m.giveSurveyPrize(heroId, surveyId)
}

func (m *SurveyModule) handleGiveSurveyPrize(r *login2game.C2SCompleteQuestionnaireProto) (*login2game.S2CCompleteQuestionnaireProto, error) {
	ok := m.giveSurveyPrize(r.HeroId, r.Qnid)
	return &login2game.S2CCompleteQuestionnaireProto{
		Success: ok,
	}, nil
}

// 处理给问卷调查奖励
func (m *SurveyModule) giveSurveyPrize(heroId int64, surveyId string) (ok bool) {
	data := m.datas.GetSurveyData(surveyId)
	if data == nil {
		logrus.Debugf("问卷调查，答题数据没找到")
		return
	}

	//if heromodule.IsHeroLocked(heroId, m.heroDataService.Func, m.guildService.GetSnapshot, data.Condition) {
	//	logrus.Debugf("问卷调查，条件未解锁")
	//	return
	//}

	// 检查
	m.heroDataService.FuncWithSend(heroId, func(hero *entity.Hero, result herolock.LockResult) {
		heroSurvey := hero.Survey()
		if heroSurvey.IsCompleted(data) {
			logrus.Debugf("问卷调查，奖励已经领取了")
			return
		}

		result.Changed()
		result.Ok()

		heroSurvey.Complete(data)

		proto := m.datas.MailHelp().SurveyMail.NewTextMail(shared_proto.MailType_MailNormal)
		proto.Prize = data.PrizeProto

		ctime := m.timeService.CurrentTime()
		m.mailModule.SendProtoMail(heroId, proto, ctime)

		// 发送消息
		result.Add(data.CompleteMsg)

		ok = true
	})
	return
}
