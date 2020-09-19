package compress

import (
	"bytes"
	"compress/gzip"
	"io/ioutil"
)

func GzipCompress(buf []byte) []byte {
	return gzipCompressionWithLevel(buf, gzip.DefaultCompression)
}

func GzipUncompress(buf []byte) ([]byte, error) {
	r, err := gzip.NewReader(bytes.NewReader(buf))
	if err != nil {
		return nil, err
	}
	defer r.Close()
	return ioutil.ReadAll(r)
}

func gzipCompressionWithLevel(buf []byte, level int) []byte {
	var b bytes.Buffer
	w, _ := gzip.NewWriterLevel(&b, level)
	w.Write(buf)
	w.Close()
	return b.Bytes()
}
