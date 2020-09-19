package msg

import (
	"io"
	"github.com/pkg/errors"
)

func ReadMsg(data []byte) (int, int, []byte, error) {
	moduleId, newData, err := ReadVarint(data)
	if err != nil {
		return 0, 0, nil, errors.Wrap(err, "readMsg moduleId, ")
	}

	msgId, newData, err := ReadVarint(newData)
	if err != nil {
		return 0, 0, nil, errors.Wrap(err, "readMsg msgId, ")
	}

	return int(moduleId), int(msgId), newData, nil
}

// Read a varint as length
func ReadVarint(b []byte) (int, []byte, error) {
	ln := len(b)
	if ln == 0 {
		return 0, nil, io.EOF
	}
	n1 := b[0]
	if n1 <= 127 {
		return int(n1), b[1:], nil
	}

	if ln == 1 {
		return 0, nil, io.EOF
	}
	n2 := b[1]
	if n2 <= 127 {
		return decodeVarint2(n1, n2), b[2:], nil
	}

	if ln == 2 {
		return 0, nil, io.EOF
	}
	n3 := b[2]

	return decodeVarint3(n1, n2, n3), b[3:], nil
}

func decodeVarint2(n1, n2 byte) int {
	return (int(n2) << 7) | (int(n1) & 0x7f)
}

// first bit of the third byte is also considered as actiondata
func decodeVarint3(n1, n2, n3 byte) int {
	return (int(n3) << 14) | ((int(n2) & 0x7f) << 7) | (int(n1) & 0x7f)
}
