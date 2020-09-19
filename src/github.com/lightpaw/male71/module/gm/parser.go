package gm

import (
	"fmt"
	"strconv"
	"strings"
	"github.com/lightpaw/male7/pb/shared_proto"
)

func parseString(s string) string {
	return s
}

func parseInt32(s string) int32 {
	i, err := strconv.ParseInt(s, 10, 32)
	if err != nil {
		fmt.Printf("解析int32参数失败，%s, %v", s, i)
	}

	return int32(i)
}

func parseBool(s string) bool {
	b, err := strconv.ParseBool(s)
	if err != nil {
		fmt.Printf("bool，%s, %v", s, b)
	}

	return b
}

func parseBytes(s string) []byte {
	return []byte(s)
}

const sep = ","

func parseStringArray(s string) []string {
	return strings.Split(s, sep)
}

func parseInt32Array(s string) []int32 {
	strArray := parseStringArray(s)

	array := make([]int32, len(strArray))
	for i, v := range strArray {
		array[i] = parseInt32(v)
	}

	return array
}

func parseBoolArray(s string) []bool {
	strArray := parseStringArray(s)

	array := make([]bool, len(strArray))
	for i, v := range strArray {
		array[i] = parseBool(v)
	}

	return array
}

func parseBytesArray(s string) [][]byte {
	strArray := parseStringArray(s)

	array := make([][]byte, len(strArray))
	for i, v := range strArray {
		array[i] = parseBytes(v)
	}

	return array
}

func parseGuildMarkProto(str string) *shared_proto.GuildMarkProto {
	return &shared_proto.GuildMarkProto{}
}
