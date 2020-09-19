package entity

import (
	"github.com/lightpaw/male7/pb/shared_proto"
	"github.com/lightpaw/male7/util/u64"
)

type HeroQuestion struct {
	UsedTimes     uint64
	EachQuestions []*EachQuestion
}

type EachQuestion struct {
	QuestionId    uint64
	QuestionState shared_proto.HeroQuestionState
	answer        uint64
}

func (hq *HeroQuestion) LastQuestion() *EachQuestion {
	size := len(hq.EachQuestions)
	if size <= 0 {
		return nil
	}
	return hq.EachQuestions[size-1]
}

func (hq *HeroQuestion) Start(qid uint64) bool {
	size := len(hq.EachQuestions)
	if size > 0 {
		return false
	}

	newQuestion := EachQuestion{QuestionId: qid, QuestionState: shared_proto.HeroQuestionState_QUESTION_WAIT}
	hq.EachQuestions = append(hq.EachQuestions, &newQuestion)
	hq.UsedTimes++
	return true
}

func (hq *HeroQuestion) Next(qid uint64) bool {
	lastQuestion := hq.LastQuestion()
	if lastQuestion == nil {
		return false
	}
	if lastQuestion.QuestionState == shared_proto.HeroQuestionState_QUESTION_WAIT {
		return false
	}

	newQuestion := EachQuestion{QuestionId: qid, QuestionState: shared_proto.HeroQuestionState_QUESTION_WAIT}
	hq.EachQuestions = append(hq.EachQuestions, &newQuestion)
	return true
}

func (hq *HeroQuestion) Answer(qid uint64, right bool, answer uint64) bool {
	lastQuestion := hq.LastQuestion()
	if lastQuestion == nil {
		return false
	}

	if lastQuestion.QuestionId != qid {
		return false
	}

	if lastQuestion.QuestionState != shared_proto.HeroQuestionState_QUESTION_WAIT {
		return false
	}

	lastQuestion.answer = answer

	if right {
		lastQuestion.QuestionState = shared_proto.HeroQuestionState_QUESTION_RIGHT
	} else {
		lastQuestion.QuestionState = shared_proto.HeroQuestionState_QUESTION_WRONG
	}

	return true
}

func (hq *HeroQuestion) Complated() {
	hq.EachQuestions = hq.EachQuestions[:0]
}

func (hq *HeroQuestion) ResetDaily() {
	hq.Complated()
	hq.UsedTimes = 0
}

func (hq *HeroQuestion) RightCount() uint64 {
	var rightCount uint64
	for _, q := range hq.EachQuestions {
		if q.QuestionState == shared_proto.HeroQuestionState_QUESTION_RIGHT {
			rightCount++
		}
	}
	return rightCount
}

func (hq *HeroQuestion) encodeClient(allRightAmount uint64) *shared_proto.HeroQuestionProto {
	out := hq.encodeServer()
	out.AllRightQuestionAmount = u64.Int32(allRightAmount)
	return out
}

func (hq *HeroQuestion) encodeServer() *shared_proto.HeroQuestionProto {
	out := &shared_proto.HeroQuestionProto{}
	out.UsedTimes = u64.Int32(hq.UsedTimes)
	for _, q := range hq.EachQuestions {
		qout := &shared_proto.HeroEachQuestionProto{}
		qout.CurrentQuestion = u64.Int32(q.QuestionId)
		qout.State = q.QuestionState
		qout.Answer = u64.Int32(q.answer)
		out.CurrentQuestion = append(out.CurrentQuestion, qout)
	}

	return out
}

func (hq *HeroQuestion) unmarshal(proto *shared_proto.HeroQuestionProto) {
	if proto == nil {
		return
	}

	hq.UsedTimes = u64.FromInt32(proto.UsedTimes)
	for _, eqp := range proto.CurrentQuestion {
		if eqp == nil {
			continue
		}
		eachQ := EachQuestion{QuestionId: u64.FromInt32(eqp.CurrentQuestion), QuestionState: eqp.State, answer:u64.FromInt32(eqp.Answer)}
		hq.EachQuestions = append(hq.EachQuestions, &eachQ)
	}
}

func NewHeroQuestion() *HeroQuestion {
	return &HeroQuestion{UsedTimes: 0, EachQuestions: []*EachQuestion{}}
}
