package military_data

import (
	"github.com/lightpaw/male7/config/resdata"
)

//gogen:config
type TrainingLevelData struct {
	_     struct{} `file:"军事/修炼馆等级.txt"`
	_     struct{} `protoconfig:"training_level"`
	_     struct{} `proto:"shared_proto.TrainingLevelProto"`
	Level uint64   `validator:"int>0"`
	Name  string
	Desc  string
	Coef  float64
	Cost  *resdata.Cost
}
