package util

import (
	"crypto/md5"
	"encoding/hex"
)

func Md5String(b []byte) string {
	sum := md5.Sum(b)
	return hex.EncodeToString(sum[:])
}
