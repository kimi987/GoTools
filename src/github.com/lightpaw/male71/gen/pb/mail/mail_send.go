package mail

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
	MODULE_ID = 8

	C2S_LIST_MAIL = 1

	C2S_DELETE_MAIL = 8

	C2S_KEEP_MAIL = 11

	C2S_COLLECT_MAIL_PRIZE = 14

	C2S_READ_MAIL = 20

	C2S_READ_MULTI = 24

	C2S_DELETE_MULTI = 26

	C2S_GET_MAIL = 28
)

func NewS2cListMailMsg(read int32, keep int32, report int32, report_tag int32, has_prize int32, collected int32, mail [][]byte) pbutil.Buffer {
	msg := &S2CListMailProto{
		Read:      read,
		Keep:      keep,
		Report:    report,
		ReportTag: report_tag,
		HasPrize:  has_prize,
		Collected: collected,
		Mail:      mail,
	}
	return NewS2cListMailProtoMsg(msg)
}

func NewS2cListMailMarshalMsg(read int32, keep int32, report int32, report_tag int32, has_prize int32, collected int32, mail [][]byte) pbutil.Buffer {
	msg := &S2CListMailProto{
		Read:      read,
		Keep:      keep,
		Report:    report,
		ReportTag: report_tag,
		HasPrize:  has_prize,
		Collected: collected,
		Mail:      mail,
	}
	return NewS2cListMailProtoMsg(msg)
}

var s2c_list_mail = [...]byte{8, 2} // 2
func NewS2cListMailProtoMsg(object *S2CListMailProto) pbutil.Buffer {

	return util.NewSnappyCompressMsg(object, s2c_list_mail[:], "s2c_list_mail")

}

// 无效的min_id
var ERR_LIST_MAIL_FAIL_INVALID_ID = pbutil.StaticBuffer{3, 8, 3, 2} // 3-2

// 无效的参数类型（0-2）
var ERR_LIST_MAIL_FAIL_INVALID_PARAM = pbutil.StaticBuffer{3, 8, 3, 4} // 3-4

// 服务器忙，请稍后再试
var ERR_LIST_MAIL_FAIL_SERVER_ERROR = pbutil.StaticBuffer{3, 8, 3, 1} // 3-1

func NewS2cReceiveMailMsg(mail []byte) pbutil.Buffer {
	msg := &S2CReceiveMailProto{
		Mail: mail,
	}
	return NewS2cReceiveMailProtoMsg(msg)
}

func NewS2cReceiveMailMarshalMsg(mail marshaler) pbutil.Buffer {
	msg := &S2CReceiveMailProto{
		Mail: safeMarshal(mail),
	}
	return NewS2cReceiveMailProtoMsg(msg)
}

var s2c_receive_mail = [...]byte{8, 4} // 4
func NewS2cReceiveMailProtoMsg(object *S2CReceiveMailProto) pbutil.Buffer {

	return util.NewSnappyCompressMsg(object, s2c_receive_mail[:], "s2c_receive_mail")

}

func NewS2cDeleteMailMsg(id []byte) pbutil.Buffer {
	msg := &S2CDeleteMailProto{
		Id: id,
	}
	return NewS2cDeleteMailProtoMsg(msg)
}

var s2c_delete_mail = [...]byte{8, 9} // 9
func NewS2cDeleteMailProtoMsg(object *S2CDeleteMailProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_delete_mail[:], "s2c_delete_mail")

}

// 发送的id无效
var ERR_DELETE_MAIL_FAIL_INVALID_ID = pbutil.StaticBuffer{3, 8, 10, 3} // 10-3

// 邮件有奖励可以领取，不能删除
var ERR_DELETE_MAIL_FAIL_NOT_EMPTY = pbutil.StaticBuffer{3, 8, 10, 1} // 10-1

// 服务器忙，请稍后再试
var ERR_DELETE_MAIL_FAIL_SERVER_ERROR = pbutil.StaticBuffer{3, 8, 10, 2} // 10-2

func NewS2cKeepMailMsg(id []byte, keep bool) pbutil.Buffer {
	msg := &S2CKeepMailProto{
		Id:   id,
		Keep: keep,
	}
	return NewS2cKeepMailProtoMsg(msg)
}

var s2c_keep_mail = [...]byte{8, 12} // 12
func NewS2cKeepMailProtoMsg(object *S2CKeepMailProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_keep_mail[:], "s2c_keep_mail")

}

// 发送的id无效
var ERR_KEEP_MAIL_FAIL_INVALID_ID = pbutil.StaticBuffer{3, 8, 13, 1} // 13-1

// 服务器忙，请稍后再试
var ERR_KEEP_MAIL_FAIL_SERVER_ERROR = pbutil.StaticBuffer{3, 8, 13, 2} // 13-2

func NewS2cCollectMailPrizeMsg(id []byte) pbutil.Buffer {
	msg := &S2CCollectMailPrizeProto{
		Id: id,
	}
	return NewS2cCollectMailPrizeProtoMsg(msg)
}

var s2c_collect_mail_prize = [...]byte{8, 15} // 15
func NewS2cCollectMailPrizeProtoMsg(object *S2CCollectMailPrizeProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_collect_mail_prize[:], "s2c_collect_mail_prize")

}

// 发送的id无效
var ERR_COLLECT_MAIL_PRIZE_FAIL_INVALID_ID = pbutil.StaticBuffer{3, 8, 16, 1} // 16-1

// 这个邮件没有奖励
var ERR_COLLECT_MAIL_PRIZE_FAIL_NOT_PRIZE = pbutil.StaticBuffer{3, 8, 16, 2} // 16-2

// 服务器忙，请稍后再试
var ERR_COLLECT_MAIL_PRIZE_FAIL_SERVER_ERROR = pbutil.StaticBuffer{3, 8, 16, 3} // 16-3

func NewS2cReadMailMsg(id []byte) pbutil.Buffer {
	msg := &S2CReadMailProto{
		Id: id,
	}
	return NewS2cReadMailProtoMsg(msg)
}

var s2c_read_mail = [...]byte{8, 21} // 21
func NewS2cReadMailProtoMsg(object *S2CReadMailProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_read_mail[:], "s2c_read_mail")

}

// 发送的id无效
var ERR_READ_MAIL_FAIL_INVALID_ID = pbutil.StaticBuffer{3, 8, 22, 2} // 22-2

// 服务器忙，请稍后再试
var ERR_READ_MAIL_FAIL_SERVER_ERROR = pbutil.StaticBuffer{3, 8, 22, 1} // 22-1

func NewS2cNotifyMailCountMsg(has_prize_not_collected_count int32, has_report_not_readed_count int32, has_yw_report_not_readed_count int32, has_bz_report_not_readed_count int32, no_report_not_readed_count int32) pbutil.Buffer {
	msg := &S2CNotifyMailCountProto{
		HasPrizeNotCollectedCount: has_prize_not_collected_count,
		HasReportNotReadedCount:   has_report_not_readed_count,
		HasYwReportNotReadedCount: has_yw_report_not_readed_count,
		HasBzReportNotReadedCount: has_bz_report_not_readed_count,
		NoReportNotReadedCount:    no_report_not_readed_count,
	}
	return NewS2cNotifyMailCountProtoMsg(msg)
}

var s2c_notify_mail_count = [...]byte{8, 23} // 23
func NewS2cNotifyMailCountProtoMsg(object *S2CNotifyMailCountProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_notify_mail_count[:], "s2c_notify_mail_count")

}

func NewS2cReadMultiMsg(ids [][]byte, report bool, prize []byte) pbutil.Buffer {
	msg := &S2CReadMultiProto{
		Ids:    ids,
		Report: report,
		Prize:  prize,
	}
	return NewS2cReadMultiProtoMsg(msg)
}

func NewS2cReadMultiMarshalMsg(ids [][]byte, report bool, prize marshaler) pbutil.Buffer {
	msg := &S2CReadMultiProto{
		Ids:    ids,
		Report: report,
		Prize:  safeMarshal(prize),
	}
	return NewS2cReadMultiProtoMsg(msg)
}

var s2c_read_multi = [...]byte{8, 25} // 25
func NewS2cReadMultiProtoMsg(object *S2CReadMultiProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_read_multi[:], "s2c_read_multi")

}

func NewS2cDeleteMultiMsg(ids [][]byte, report bool) pbutil.Buffer {
	msg := &S2CDeleteMultiProto{
		Ids:    ids,
		Report: report,
	}
	return NewS2cDeleteMultiProtoMsg(msg)
}

var s2c_delete_multi = [...]byte{8, 27} // 27
func NewS2cDeleteMultiProtoMsg(object *S2CDeleteMultiProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_delete_multi[:], "s2c_delete_multi")

}

func NewS2cGetMailMsg(data []byte) pbutil.Buffer {
	msg := &S2CGetMailProto{
		Data: data,
	}
	return NewS2cGetMailProtoMsg(msg)
}

func NewS2cGetMailMarshalMsg(data marshaler) pbutil.Buffer {
	msg := &S2CGetMailProto{
		Data: safeMarshal(data),
	}
	return NewS2cGetMailProtoMsg(msg)
}

var s2c_get_mail = [...]byte{8, 29} // 29
func NewS2cGetMailProtoMsg(object *S2CGetMailProto) pbutil.Buffer {

	return util.NewSnappyCompressMsg(object, s2c_get_mail[:], "s2c_get_mail")

}

// 邮件不存在
var ERR_GET_MAIL_FAIL_MAIL_NOT_FOUND = pbutil.StaticBuffer{3, 8, 30, 1} // 30-1
