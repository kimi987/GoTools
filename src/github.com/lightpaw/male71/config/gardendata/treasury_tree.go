package gardendata

import (
	"github.com/lightpaw/male7/pb/shared_proto"
	"github.com/lightpaw/male7/config/resdata"
	"github.com/lightpaw/male7/config/i18n"
)

//gogen:config
type TreasuryTreeData struct {
	_ struct{} `file:"花园/摇钱树.txt"`
	_ struct{} `proto:"shared_proto.TreasuryTreeDataProto"`
	_ struct{} `protoconfig:"TreasuryTreeData"`

	Id uint64 `head:"-,uint64(%s.Season)" protofield:"-"`

	// 季节
	Season shared_proto.Season

	Prize *resdata.Prize

	Desc *i18n.I18nRef
}
