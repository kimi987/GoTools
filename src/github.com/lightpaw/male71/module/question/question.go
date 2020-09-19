package question

import (
	"github.com/lightpaw/male7/gen/iface"
	"github.com/lightpaw/male7/entity"
	"github.com/lightpaw/male7/service/conflict/heroservice/herolock"
	"github.com/lightpaw/male7/util/u64"
	"github.com/lightpaw/logrus"
	"github.com/lightpaw/male7/pb/shared_proto"
	"github.com/lightpaw/male7/service/heromodule"
	"github.com/lightpaw/male7/gen/pb/question"
	"github.com/lightpaw/male7/pb/server_proto"
	"github.com/lightpaw/male7/service/operate_type"
)

func NewQuestionModule(dep iface.ServiceDep) *QuestionModule {
	m := &QuestionModule{}
	m.dep = dep
	m.datas = dep.Datas()
	m.timeService = dep.Time()
	m.broadcast = dep.Broadcast()

	return m
}

//gogen:iface
type QuestionModule struct {
	dep         iface.ServiceDep
	datas       iface.ConfigDatas
	timeService iface.TimeService
	broadcast   iface.BroadcastService
}

//gogen:iface
func (m *QuestionModule) ProcessStart(proto *question.C2SStartProto, hc iface.HeroController) {
	hc.FuncWithSend(func(hero *entity.Hero, result herolock.LockResult) {
		qid := u64.FromInt32(proto.Id)

		if m.datas.QuestionData().Get(qid) == nil {
			result.Add(question.ERR_START_FAIL_INVALID_ID)
			return
		}

		if hero.Question().UsedTimes >= m.datas.QuestionMiscData().MaxTimes {
			logrus.Debugf("答题，次数用完了")
			result.Add(question.ERR_START_FAIL_NO_TIMES)
			return
		}

		if hero.Question().LastQuestion() != nil {
			logrus.Debugf("答题，正在答题中")
			result.Add(question.ERR_START_FAIL_IN_QUESTION)
			return
		}

		if !hero.Question().Start(qid) {
			logrus.Debugf("答题-开始，前面验证通过了还失败")
			result.Add(question.ERR_START_FAIL_IN_QUESTION)
			return
		}
		result.Add(question.START_S2C)

		result.Changed()
		result.Ok()
	})
}

//gogen:iface
func (m *QuestionModule) ProcessAnswer(proto *question.C2SAnswerProto, hc iface.HeroController) {
	hc.FuncWithSend(func(hero *entity.Hero, result herolock.LockResult) {
		lastQuestion := hero.Question().LastQuestion()
		if lastQuestion == nil {
			logrus.Debugf("答题还没开始")
			result.Add(question.ERR_ANSWER_FAIL_NOT_START)
			return
		}

		qid := u64.FromInt32(proto.Id)
		if lastQuestion.QuestionId != qid {
			result.Add(question.ERR_ANSWER_FAIL_INVALID_QID)
			return
		}

		aid := u64.FromInt32(proto.Answer)
		if aid <= 0 {
			result.Add(question.ERR_ANSWER_FAIL_INVALID_AID)
			return
		}

		if lastQuestion.QuestionState != shared_proto.HeroQuestionState_QUESTION_WAIT {
			result.Add(question.ERR_ANSWER_FAIL_ALREADY_ANSWERED)
			return
		}

		if !hero.Question().Answer(qid, proto.Right, aid) {
			logrus.Debugf("答题-回答,前面验证通过了还失败")
			result.Add(question.ERR_ANSWER_FAIL_ALREADY_ANSWERED)
			return
		}
		result.Add(question.NewS2cAnswerMsg(u64.Int32(qid)))

		result.Changed()
		result.Ok()

		if proto.Right {
			hero.HistoryAmount().IncreaseOne(server_proto.HistoryAmountType_QuestionRightAmount)
			heromodule.UpdateTaskProgress(hero, result, shared_proto.TaskTargetType_TASK_TARGET_ALL_RIGHT_QUESTION_AMOUNT)
		}

		m.dep.Tlog().TlogAnswerFlow(hero, qid, proto.Right)
	})
}

//gogen:iface
func (m *QuestionModule) ProcessNext(proto *question.C2SNextProto, hc iface.HeroController) {
	hc.FuncWithSend(func(hero *entity.Hero, result herolock.LockResult) {
		if len(hero.Question().EachQuestions) >= u64.Int(m.datas.QuestionMiscData().QuestionCount) {
			result.Add(question.ERR_NEXT_FAIL_ENOUGH)
			return
		}

		lastQuestion := hero.Question().LastQuestion()
		if lastQuestion == nil || lastQuestion.QuestionState == shared_proto.HeroQuestionState_QUESTION_WAIT {
			result.Add(question.ERR_NEXT_FAIL_LAST_NOT_ANSWER)
			return
		}

		qid := u64.FromInt32(proto.Id)
		if m.datas.QuestionData().Get(qid) == nil {
			result.Add(question.ERR_NEXT_FAIL_INVALID_QID)
			return
		}

		if !hero.Question().Next(qid) {
			logrus.Debugf("答题-下一题,前面验证通过了还失败")
			result.Add(question.ERR_NEXT_FAIL_INVALID_QID)
			return
		}
		result.Add(question.NEXT_S2C)

		result.Changed()
		result.Ok()
	})
}

//gogen:iface
func (m *QuestionModule) ProcessGetPrize(proto *question.C2SGetPrizeProto, hc iface.HeroController) {
	hc.FuncWithSend(func(hero *entity.Hero, result herolock.LockResult) {
		hq := hero.Question()
		if u64.FromInt(len(hq.EachQuestions)) < m.datas.QuestionMiscData().QuestionCount {
			result.Add(question.ERR_GET_PRIZE_FAIL_NOT_FINISH)
			return
		}

		if hq.LastQuestion() == nil || hq.LastQuestion().QuestionState == shared_proto.HeroQuestionState_QUESTION_WAIT {
			result.Add(question.ERR_GET_PRIZE_FAIL_NOT_FINISH)
			return
		}

		rightCount := hq.RightCount()
		//if u64.FromInt32(proto.Score) != rightCount {
		//	result.Add(question.ERR_GET_PRIZE_FAIL_SCORE_ERR)
		//	return
		//}

		prize := m.datas.QuestionPrizeData().Get(rightCount)
		if prize != nil && prize.Prize != nil {
			ctime := m.timeService.CurrentTime()
			hctx := heromodule.NewContext(m.dep, operate_type.QuestionPrize)
			heromodule.AddPrize(hctx, hero, result, prize.Prize, ctime)
		}

		hq.Complated()
		result.Add(question.GET_PRIZE_S2C)

		heromodule.IncreTaskProgressOne(hero, result, shared_proto.TaskTargetType_TASK_TARGET_ACTIVE_START_QUESTION_COUNT)
		hero.HistoryAmount().IncreaseOne(server_proto.HistoryAmountType_AccumStartQuestion)
		heromodule.UpdateTaskProgress(hero, result, shared_proto.TaskTargetType_TASK_TARGET_ACCUM_START_QUESTION)

		result.Changed()
		result.Ok()
	})
}
