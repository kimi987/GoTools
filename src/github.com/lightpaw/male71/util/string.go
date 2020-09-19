package util

import (
	"unicode/utf8"
	"unsafe"
	"bytes"
	"unicode"
	"strings"
)

// 快速把一个[]byte转换为string. 使string在背后直接使用这个[]byte作为它的数据而不再copy一份
// 调用之后不能再修改原始的[]byte中的数据, 不然会导致string被修改
func Byte2String(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}

// 快速把一个string转换为[]byte. 直接使用string背后的[]byte
// 调用之后不能修改[]byte中的数据, 不然会导致string被修改
func String2Byte(b string) []byte {
	return *(*[]byte)(unsafe.Pointer(&b))
}

func GetCharLen(s string) int {

	charLen := 0
	for _, r := range []rune(s) {
		n := utf8.RuneLen(r)
		switch n {
		case -1:
			// 当成最大的unicode字符处理
			charLen += 4
		case 1:
			charLen += 1
		default:
			// 多字节字符，一个当成2个
			charLen += 2
		}
	}

	return charLen
}

func TruncateCharLen(s string, maxCharLen int) string {
	runeArray := []rune(s)
	b := &bytes.Buffer{}

	charLen := 0
	for _, r := range runeArray {
		n := utf8.RuneLen(r)
		switch n {
		case -1:
			// 跳过这种字符
			continue
		case 1:
			charLen += 1
		default:
			// 多字节字符，一个当成2个
			charLen += 2
		}

		if charLen > maxCharLen {
			break
		}
		b.WriteRune(r)
	}

	return b.String()
}

// 替换非法的字符
func ReplaceInvalidChar(s string) string {
	runeArray := []rune(s)
	buffer := bytes.NewBuffer(make([]byte, 0, len(runeArray)))

	for _, r := range runeArray {
		if isValidRune(r, true) {
			buffer.WriteRune(r)
		}
	}

	return strings.TrimSpace(buffer.String())
}

func IsValidName(s string) bool {

	// 名字，不能空格，符号开头
	// 不能包含非法字符

	runeArray := []rune(s)
	lastIdx := len(runeArray) - 1
	for i, r := range runeArray {
		puntValid := !(i == 0 || i == lastIdx)

		// 开头必须是字母，或者数字，不能是符号开头
		if !isValidRune(r, puntValid) {
			return false
		}
	}

	return true
}

func HaveInvalidChar(s string) bool {
	runeArray := []rune(s)

	for _, r := range runeArray {
		if !isValidRune(r, true) {
			return true
		}
	}

	return false
}

func isValidRune(r rune, punctValid bool) bool {
	switch r {
	case '_', ' ':
		if !punctValid {
			return false
		}
	case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':

	default:
		if !unicode.IsLetter(r) {
			return false
		}
	}

	return true
}
