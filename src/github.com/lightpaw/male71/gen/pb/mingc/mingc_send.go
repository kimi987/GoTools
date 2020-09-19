package mingc

import (
	"github.com/lightpaw/male7/gen/pb/util"
	"github.com/lightpaw/male7/pb/shared_proto"
	"github.com/lightpaw/pbutil"
)

var (
	pool           = pbutil.Pool
	newProtoMsg    = util.NewProtoMsg
	newCompressMsg = util.NewCompressMsg
	safeMarshal    = util.SafeMarshal
	_              = shared_proto.ErrIntOverflowConfig
)

type marshaler util.Marshaler

const (
	MODULE_ID = 42

	C2S_MINGC_LIST = 4

	C2S_VIEW_MINGC = 7

	C2S_MC_BUILD = 10

	C2S_MC_BUILD_LOG = 13

	C2S_MINGC_HOST_GUILD = 20
)

func NewS2cMingcListMsg(ver int32, mingcs *shared_proto.MingcsProto) pbutil.Buffer {
	msg := &S2CMingcListProto{
		Ver:    ver,
		Mingcs: mingcs,
	}
	return NewS2cMingcListProtoMsg(msg)
}

var s2c_mingc_list = [...]byte{42, 5} // 5
func NewS2cMingcListProtoMsg(object *S2CMingcListProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_mingc_list[:], "s2c_mingc_list")

}

// 服务器繁忙，请稍后再试
var ERR_MINGC_LIST_FAIL_SEVER_ERROR = pbutil.StaticBuffer{3, 42, 6, 1} // 6-1

func NewS2cViewMingcMsg(mingc *shared_proto.MingcProto) pbutil.Buffer {
	msg := &S2CViewMingcProto{
		Mingc: mingc,
	}
	return NewS2cViewMingcProtoMsg(msg)
}

var s2c_view_mingc = [...]byte{42, 8} // 8
func NewS2cViewMingcProtoMsg(object *S2CViewMingcProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_view_mingc[:], "s2c_view_mingc")

}

// 没有这座名城
var ERR_VIEW_MINGC_FAIL_INVALID_ID = pbutil.StaticBuffer{3, 42, 9, 1} // 9-1

func NewS2cMcBuildMsg(mc_id int32, new_level int32, new_support int32, new_daily_added_support int32, next_time int32) pbutil.Buffer {
	msg := &S2CMcBuildProto{
		McId:                 mc_id,
		NewLevel:             new_level,
		NewSupport:           new_support,
		NewDailyAddedSupport: new_daily_added_support,
		NextTime:             next_time,
	}
	return NewS2cMcBuildProtoMsg(msg)
}

var s2c_mc_build = [...]byte{42, 11} // 11
func NewS2cMcBuildProtoMsg(object *S2CMcBuildProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_mc_build[:], "s2c_mc_build")

}

// 君主等级不够
var ERR_MC_BUILD_FAIL_HERO_LEVEL_LIMIT = pbutil.StaticBuffer{3, 42, 12, 5} // 12-5

// 名城 id 错误
var ERR_MC_BUILD_FAIL_INVALID_MC_ID = pbutil.StaticBuffer{3, 42, 12, 3} // 12-3

// 没有次数了
var ERR_MC_BUILD_FAIL_NO_COUNT = pbutil.StaticBuffer{3, 42, 12, 1} // 12-1

// 在 cd 中
var ERR_MC_BUILD_FAIL_IN_CD = pbutil.StaticBuffer{3, 42, 12, 2} // 12-2

// 没有联盟
var ERR_MC_BUILD_FAIL_NO_GUILD = pbutil.StaticBuffer{3, 42, 12, 4} // 12-4

func NewS2cMcBuildLogMsg(logs *shared_proto.GuildMcBuildProto) pbutil.Buffer {
	msg := &S2CMcBuildLogProto{
		Logs: logs,
	}
	return NewS2cMcBuildLogProtoMsg(msg)
}

var s2c_mc_build_log = [...]byte{42, 14} // 14
func NewS2cMcBuildLogProtoMsg(object *S2CMcBuildLogProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_mc_build_log[:], "s2c_mc_build_log")

}

// 名城 id 错误
var ERR_MC_BUILD_LOG_FAIL_INVALID_MC_ID = pbutil.StaticBuffer{3, 42, 15, 1} // 15-1

var RESET_DAILY_MC_S2C = pbutil.StaticBuffer{2, 42, 16} // 16

func NewS2cMingcHostGuildMsg(guild *shared_proto.GuildSnapshotProto) pbutil.Buffer {
	msg := &S2CMingcHostGuildProto{
		Guild: guild,
	}
	return NewS2cMingcHostGuildProtoMsg(msg)
}

var s2c_mingc_host_guild = [...]byte{42, 21} // 21
func NewS2cMingcHostGuildProtoMsg(object *S2CMingcHostGuildProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_mingc_host_guild[:], "s2c_mingc_host_guild")

}

// 名城id错误
var ERR_MINGC_HOST_GUILD_FAIL_INVALID_MC_ID = pbutil.StaticBuffer{3, 42, 22, 1} // 22-1
