package must

import (
	"github.com/lightpaw/logrus"
	"github.com/lightpaw/pbutil"
)

var empty = pbutil.StaticBuffer([]byte{})

func Msg(data pbutil.Buffer, err error) pbutil.Buffer {
	if err != nil {
		logrus.WithError(err).Errorf("must.Msg fail")
		return pbutil.Empty
	}
	return data
}
