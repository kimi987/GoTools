package survey

import (
	"fmt"
	"github.com/lightpaw/male7/config/data"
	"github.com/lightpaw/male7/config/resdata"
	"github.com/lightpaw/male7/gen/pb/survey"
	"github.com/lightpaw/male7/pb/shared_proto"
	"github.com/lightpaw/pbutil"
	"strings"
)

// 问卷调查
//gogen:config
type SurveyData struct {
	_ struct{} `file:"问卷调查/问卷调查.txt"`
	_ struct{} `proto:"shared_proto.SurveyDataProto"`
	_ struct{} `protoconfig:"SurveyDatas"`

	Id          string                   // 问卷id，ID越大，越是没答过的越是优先答题
	Name        string                   // 名字
	Icon        string                   // 图标
	Url         string                   // 问卷链接
	Condition   *data.UnlockCondition    `type:"sub"`                     // 解锁条件
	Prize       *resdata.Prize           `protofield:"-"`                 // 奖励
	PrizeProto  *shared_proto.PrizeProto `head:"-" protofield:"Prize,%s"` // 奖励proto
	CompleteMsg pbutil.Buffer            `head:"-" protofield:"-"`        // 完成消息
}

func (data *SurveyData) Init() {
	data.CompleteMsg = survey.NewS2cCompleteMsg(data.Id).Static()
	data.PrizeProto = data.Prize.Encode4Init()
	if strings.Contains(data.Url, "?") {
		data.Url += fmt.Sprintf("&serverid={{serverid}}&rid=%s", data.Id)
	} else {
		data.Url += fmt.Sprintf("?serverid={{serverid}}&rid=%s", data.Id)
	}
}
