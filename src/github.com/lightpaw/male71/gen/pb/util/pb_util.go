package util

import (
	"bytes"
	"compress/gzip"
	"errors"
	"github.com/golang/snappy"
	"github.com/lightpaw/logrus"
	"github.com/lightpaw/pbutil"
	"io/ioutil"
)

var (
	pool = pbutil.Pool
)

type Marshaler interface {
	Marshal() (dAtA []byte, err error)
}

func SafeMarshal(m Marshaler) []byte {
	data, err := m.Marshal()
	if err != nil {
		logrus.WithError(err).Errorf("safe.Marshal fail")
	}
	return data
}

type proto interface {
	Size() int
	MarshalTo([]byte) (int, error)
}

const (
	skipCompressLen = 64
)

func NewProtoMsg(object proto, head []byte, msgName string) pbutil.Buffer {
	if IsDebug {
		return tryCompressMsg(object, head, msgName)
	} else {
		m, _ := newProtoMsg(object, head, msgName)
		return m
	}
}

func tryCompressMsg(object proto, head []byte, msgName string) pbutil.Buffer {

	// 消息长度 > 127字节
	result, n, cn := newCompressMsgReturnLen(object, head, msgName, s2c_gzip_compress_msg[:], gzipBestCompression)

	if cn+DiffByte < n {
		_, defn, defcn := newCompressMsgReturnLen(object, head, msgName, s2c_gzip_compress_msg[:], gzipDefaultCompression)
		_, snapn, snapcn := newCompressMsgReturnLen(object, head, msgName, s2c_snappy_compress_msg[:], snappyCompression)

		logrus.WithField("n", n).WithField("cn", cn).WithField("rate", float64(n-cn)/float64(n)).
			WithField("defn", defn).WithField("defcn", defcn).WithField("defrate", float64(defn-defcn)/float64(defn)).
			WithField("snapn", snapn).WithField("snapcn", snapcn).WithField("snaprate", float64(snapn-snapcn)/float64(snapn)).
			WithField("head", head).Warn("消息优化，发现压缩效率高的消息")
	}

	return result
}

func newProtoMsg(object proto, head []byte, msgName string) (pbutil.Buffer, []byte) {
	headLen := len(head)
	_size := headLen + object.Size()
	result, buf := allocBuff(_size)

	copy(buf, head)
	if _, err := object.MarshalTo(buf[headLen:]); err != nil {
		result.Free()
		logrus.WithError(err).Errorf("%s.Marshal fail", msgName)
		return pbutil.Empty, nil
	} else {
		return result, buf
	}
}

const (
	compress_module_id          = 0
	snappy_compress_sequence_id = 0
	gzip_compress_sequence_id   = 1
)

func NewCompressMsg(object proto, head []byte, msgName string) pbutil.Buffer {
	return NewSnappyCompressMsg(object, head, msgName)
}

var s2c_snappy_compress_msg = [...]byte{compress_module_id, snappy_compress_sequence_id}

func snappyCompression(buf []byte) []byte {
	return snappy.Encode(nil, buf)
}

func NewSnappyCompressMsg(object proto, head []byte, msgName string) pbutil.Buffer {
	return newCompressMsg(object, head, msgName, s2c_snappy_compress_msg[:], snappyCompression)
}

var s2c_gzip_compress_msg = [...]byte{compress_module_id, gzip_compress_sequence_id}

func gzipBestCompression(buf []byte) []byte {
	return gzipCompressionWithLevel(buf, gzip.BestCompression)
}

func gzipDefaultCompression(buf []byte) []byte {
	return gzipCompressionWithLevel(buf, gzip.DefaultCompression)
}

func gzipCompressionWithLevel(buf []byte, level int) []byte {
	var b bytes.Buffer
	w, _ := gzip.NewWriterLevel(&b, level)
	w.Write(buf)
	w.Close()
	return b.Bytes()
}

func NewGzipCompressMsg(object proto, head []byte, msgName string) pbutil.Buffer {
	return newCompressMsg(object, head, msgName, s2c_gzip_compress_msg[:], gzipBestCompression)
}

func newCompressMsg(object proto, head []byte, msgName string, compressHead []byte, compressFunc func(buf []byte) []byte) pbutil.Buffer {
	result, n, cn := newCompressMsgReturnLen(object, head, msgName, compressHead, compressFunc)
	if IsDebug && n < cn {
		logrus.WithField("buf", n).WithField("compress_buf", cn).WithField("head", head).Warn("消息优化，发现压缩效率低的消息")
	}

	return result
}

func newCompressMsgReturnLen(object proto, head []byte, msgName string, compressHead []byte, compressFunc func(buf []byte) []byte) (pbutil.Buffer, int, int) {

	result, buf := newProtoMsg(object, head, msgName)
	n := len(buf)
	if n < skipCompressLen || n <= 0 {
		return result, 0, 0
	}

	compressBuffer := compressFunc(buf)
	cn := len(compressBuffer)
	if cn < n {
		// 压缩后比原来的小，使用压缩后数据
		defer result.Free()
		return NewBytesMsg(compressBuffer, compressHead), n, cn
	}

	return result, n, cn
}

func NewBytesMsg(data, head []byte) pbutil.Buffer {
	headLen := len(head)
	_size := headLen + len(data)
	result, buf := allocBuff(_size)

	copy(buf, head)
	copy(buf[headLen:], data)
	return result
}

func allocBuff(_size int) (result *pbutil.RecycleBuffer, buf []byte) {
	switch {
	case _size <= 127:
		result = pool.Alloc(_size + 1)
		buf = result.Buffer()
		buf[0] = uint8(_size)
		buf = buf[1:]
	case _size <= 16383:
		result = pool.Alloc(_size + 2)
		buf = result.Buffer()
		buf[0], buf[1] = (0x80 | uint8(_size&0x7f)), uint8(_size>>7)
		buf = buf[2:]
	default:
		result = pool.Alloc(_size + 3)
		buf = result.Buffer()
		buf[0], buf[1], buf[2] = (0x80 | uint8(_size&0x7f)), (0x80 | uint8((_size>>7)&0x7f)), uint8(_size>>14)
		buf = buf[3:]
	}

	return
}

func IsCompressMsg(moduleID, sequenceID int) bool {
	if compress_module_id == moduleID {
		switch sequenceID {
		case snappy_compress_sequence_id, gzip_compress_sequence_id:
			return true
		}
	}
	return false
}

var errUnkownCompressSequenceId = errors.New("Unkown compress msg sequence")

func UncompressMsg(sequenceID int, data []byte) ([]byte, error) {
	switch sequenceID {
	case snappy_compress_sequence_id:
		return snappy.Decode(nil, data)
	case gzip_compress_sequence_id:
		r, err := gzip.NewReader(bytes.NewReader(data))
		if err != nil {
			return nil, err
		}
		defer r.Close()
		uncomress, err := ioutil.ReadAll(r)
		return uncomress, err
	default:
		return nil, errUnkownCompressSequenceId
	}
}
