package util

import (
	"math/rand"
	"github.com/tinylib/msgp/msgp"
	"encoding/base64"
)

func NewCpOrderId(sid uint32, heroId int64, productId uint64, moneyFen uint64, ctime int64) string {
	return newCpOrderIdV1(sid, heroId, productId, moneyFen, ctime)
}

const (
	VersionBitCount = 24

	V1Bit    = 1 << VersionBitCount
	RandMask = (1 << VersionBitCount) - 1
)

func newCpOrderIdV1(sid uint32, heroId int64, productId uint64, moneyFen uint64, ctime int64) string {
	// version(1) + random(3)
	randNum := rand.Uint32()&RandMask | V1Bit

	// 5 + 3 + 5 + 5 + 5 + 5
	b := make([]byte, 0, 28)
	b = msgp.AppendUint32(b, randNum)
	b = msgp.AppendUint32(b, sid)
	b = msgp.AppendInt64(b, heroId)
	b = msgp.AppendUint64(b, productId)
	b = msgp.AppendUint64(b, moneyFen)
	b = msgp.AppendInt64(b, ctime)

	return base64.RawURLEncoding.EncodeToString(b)
}

func NewCpOrderSign(cpOrderId, key string) string {
	return Md5String([]byte(cpOrderId + "_LP_" + key))
}

func ParseCpOrderId(orderId string) (version, randNum, sid uint32, heroId int64, productId uint64, moneyFen uint64, ctime int64, err error) {
	b, err := base64.RawURLEncoding.DecodeString(orderId)
	if err != nil {
		return
	}

	randNum, b, err = msgp.ReadUint32Bytes(b)
	if err != nil {
		return
	}

	sid, b, err = msgp.ReadUint32Bytes(b)
	if err != nil {
		return
	}

	heroId, b, err = msgp.ReadInt64Bytes(b)
	if err != nil {
		return
	}

	productId, b, err = msgp.ReadUint64Bytes(b)
	if err != nil {
		return
	}

	moneyFen, b, err = msgp.ReadUint64Bytes(b)
	if err != nil {
		return
	}

	ctime, b, err = msgp.ReadInt64Bytes(b)
	if err != nil {
		return
	}

	return randNum >> VersionBitCount, randNum, sid, heroId, productId, moneyFen, ctime, nil
}
