package entity

import (
	"github.com/lightpaw/logrus"
	"github.com/lightpaw/male7/pb/shared_proto"
	. "github.com/onsi/gomega"
	"testing"
)

func TestProgressType(t *testing.T) {
	RegisterTestingT(t)

	// duplicate
	for k := range increTaskTargetTypes {
		if _, exist := updateTaskProgressFuncs[k]; exist {
			logrus.Panicf("increTaskTargetTypes 和 updateTaskProgressFuncs都配置了，%v", k)
		}
	}

	// register invoke
	targetMap := map[shared_proto.TaskTargetType]string{
		shared_proto.TaskTargetType_InvalidTaskTargetType:   "InvalidTaskTargetType",
		shared_proto.TaskTargetType_TASK_TARGET_GUILD_LEVEL: "TASK_TARGET_GUILD_LEVEL",
	}

	hasError := false
	for k := range increTaskTargetTypes {
		if _, exist := targetMap[k]; exist {
			hasError = true
			logrus.Errorf("任务类型重复，increTaskTargetTypes 和 updateTaskProgressFuncs都配置了，%v", k)
		}

		targetMap[k] = shared_proto.TaskTargetType_name[int32(k)]
	}

	for k := range updateTaskProgressFuncs {
		if _, exist := targetMap[k]; exist {
			hasError = true
			logrus.Errorf("任务类型重复，increTaskTargetTypes 和 updateTaskProgressFuncs都配置了，%v", k)
		}

		targetMap[k] = shared_proto.TaskTargetType_name[int32(k)]
	}

	for k, v := range shared_proto.TaskTargetType_name {
		kt := shared_proto.TaskTargetType(k)
		if _, exist := targetMap[kt]; !exist {
			hasError = true
			logrus.Errorf("任务类型缺少处理函数，%v, %s", k, v)
		}
	}

	Ω(hasError).Should(BeFalse())
}
