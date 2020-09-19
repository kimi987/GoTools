package compress

import "github.com/golang/snappy"

func SnappyCompress(buf []byte) []byte {
	return snappy.Encode(nil, buf)
}

func SnappyUncompress(buf []byte) ([]byte, error) {
	return snappy.Decode(nil, buf)
}
