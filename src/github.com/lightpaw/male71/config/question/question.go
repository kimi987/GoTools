package question

import (
	"github.com/lightpaw/male7/config/resdata"
	"github.com/lightpaw/male7/util/check"
	"github.com/lightpaw/male7/util/u64"
)

//答题问题
//gogen:config
type QuestionData struct {
	_ struct{} `file:"答题/答题问题.txt"`
	_ struct{} `proto:"shared_proto.QuestionProto"`
	_ struct{} `protoconfig:"question"`

	Id          uint64   `validator:"int>0"`
	Question    string   `validator:"string>0"`
	RightAnswer string   `validator:"string>0"`
	WrongAnswer []string `validator:"string>0"`
}

func (d *QuestionData) Init(filename string) {
	check.PanicNotTrue(len(d.WrongAnswer) > 1, "%s 问题：%s 的错误答案必须大于1个.now:%s", filename, d.Id)
}

//答题名言
//gogen:config
type QuestionSayingData struct {
	_ struct{} `file:"答题/答题名言.txt"`
	_ struct{} `proto:"shared_proto.QuestionSayingProto"`
	_ struct{} `protoconfig:"question_saying"`

	Id      uint64
	Content string `validator:"string>0"`
	Author  string
}

//答题杂项
//gogen:config
type QuestionMiscData struct {
	_ struct{} `singleton:"true"`
	_ struct{} `file:"答题/答题杂项.txt"`
	_ struct{} `proto:"shared_proto.QuestionMiscProto"`
	_ struct{} `protoconfig:"question_misc"`

	MaxTimes      uint64 `validator:"int>0"`
	QuestionCount uint64 `validator:"int>0"`
}

func (d *QuestionMiscData) Init(filename string, configData interface {
	GetQuestionDataArray() []*QuestionData
}) {
	questionDataLen := len(configData.GetQuestionDataArray())
	check.PanicNotTrue(d.QuestionCount <= u64.FromInt(questionDataLen), "%s 每轮题数必须<=题库总数%s", filename, questionDataLen)
}

//答题奖励
//gogen:config
type QuestionPrizeData struct {
	_ struct{} `file:"答题/答题奖励.txt"`
	_ struct{} `proto:"shared_proto.QuestionPrizeProto"`
	_ struct{} `protoconfig:"question_prize"`

	Score uint64 `validator:"int" key:"true"`
	Prize *resdata.Prize
}

func (d *QuestionPrizeData) Init(filename string, configData interface {
	QuestionMiscData() *QuestionMiscData
}) {
	miscData := configData.QuestionMiscData()
	check.PanicNotTrue(d.Score >= 0 && d.Score <= miscData.QuestionCount, "%s 答对的题数必须>=0且<=每轮的题数%s", filename, miscData.QuestionCount)
}
