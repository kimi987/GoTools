package config

import (
	"bytes"
	"github.com/lightpaw/logrus"
	"github.com/pkg/errors"
	"strconv"
	"strings"
)

type ObjectParser struct {
	dataMap map[string][]string

	line int
}

func (p *ObjectParser) Line() int {
	return p.line
}

func (p *ObjectParser) String(key string) string {
	sa := p.OriginStringArray(key)

	if len(sa) > 0 {
		return replaceNewLine(sa[0])
	}

	return ""
}

func replaceNewLine(s string) string {
	s = strings.Replace(s, "\\n", "\n", -1)
	s = strings.Replace(s, "\\r", "\r", -1)
	return s
}

func (p *ObjectParser) Int(key string) int {
	s := p.String(key)
	if len(s) == 0 {
		return 0
	}

	i, err := strconv.Atoi(s)
	if err != nil {
		return 0
	}

	return i
}

func (p *ObjectParser) Int64(key string) int64 {
	s := p.String(key)
	if len(s) == 0 {
		return 0
	}

	i, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return 0
	}

	return i
}

func (p *ObjectParser) Uint64(key string) uint64 {
	s := p.String(key)
	if len(s) == 0 {
		return 0
	}

	i, err := strconv.ParseUint(s, 10, 64)
	if err != nil {
		return 0
	}

	return i
}

func (p *ObjectParser) Float64(key string) float64 {
	s := p.String(key)
	if len(s) == 0 {
		return 0
	}

	i, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return 0
	}

	return i
}

func (p *ObjectParser) Bool(key string) bool {
	s := p.String(key)
	if len(s) == 0 {
		return false
	}

	i, err := strconv.ParseBool(s)
	if err != nil {
		return false
	}

	return i
}

func (p *ObjectParser) KeyExist(key string) bool {
	_, ok := p.dataMap[strings.ToLower(key)]
	return ok
}

func (p *ObjectParser) OriginStringArray(key string) []string {
	return p.dataMap[strings.ToLower(key)]
}

func (p *ObjectParser) IntArray(key, sep string, nullable bool) []int {
	sa := p.StringArray(key, sep, nullable)

	out := make([]int, 0, len(sa))
	for _, s := range sa {

		v, err := strconv.Atoi(s)
		if err != nil {
			logrus.Errorf("配置解析错误(之前不是检查过类型吗...)，IntArray %s %s, %s", key, sep, sa)

			out = append(out, -1)
			continue
		}

		out = append(out, v)
	}

	return out
}

func (p *ObjectParser) Int64Array(key, sep string, nullable bool) []int64 {
	sa := p.StringArray(key, sep, nullable)

	out := make([]int64, 0, len(sa))
	for _, s := range sa {

		v, err := strconv.ParseInt(s, 10, 64)
		if err != nil {
			logrus.Errorf("配置解析错误(之前不是检查过类型吗...)，Int64Array %s %s, %s", key, sep, sa)

			out = append(out, -1)
			continue
		}

		out = append(out, v)
	}

	return out
}

func (p *ObjectParser) Uint64Array(key, sep string, nullable bool) []uint64 {
	sa := p.StringArray(key, sep, nullable)

	out := make([]uint64, 0, len(sa))
	for _, s := range sa {

		v, err := strconv.ParseUint(s, 10, 64)
		if err != nil {
			logrus.Errorf("配置解析错误(之前不是检查过类型吗...)，Uint64Array %s %s, %s", key, sep, sa)

			out = append(out, 0)
			continue
		}

		out = append(out, v)
	}

	return out
}

func (p *ObjectParser) Float64Array(key, sep string, nullable bool) []float64 {
	sa := p.StringArray(key, sep, nullable)

	out := make([]float64, 0, len(sa))
	for _, s := range sa {

		v, err := strconv.ParseFloat(s, 64)
		if err != nil {
			logrus.Errorf("配置解析错误(之前不是检查过类型吗...)，Float64Array %s %s, %s", key, sep, sa)

			out = append(out, -1)
			continue
		}

		out = append(out, v)
	}

	return out
}

func (p *ObjectParser) BoolArray(key, sep string, nullable bool) []bool {
	sa := p.StringArray(key, sep, nullable)

	out := make([]bool, 0, len(sa))
	for _, s := range sa {

		v, err := strconv.ParseBool(s)
		if err != nil {
			logrus.Errorf("配置解析错误(之前不是检查过类型吗...)，BoolArray %s %s, %s", key, sep, sa)

			out = append(out, false)
			continue
		}

		out = append(out, v)
	}

	return out
}

func (p *ObjectParser) StringArray(key, sep string, nullable bool) []string {
	in := p.OriginStringArray(key)
	if len(in) == 0 {
		return nil
	}

	out := in
	if sep != "" {
		if len(in) > 1 {
			logrus.Errorf("配置解析错误(之前不是检查过类型吗...)，StringArray len(in) > 1, %s, %s", in, sep)
		}

		out = strings.Split(in[0], sep)
	}

	if nullable {
		return out
	}

	newOut := make([]string, 0, len(out))
	for _, v := range out {
		if len(v) > 0 {
			newOut = append(newOut, replaceNewLine(v))
		}
	}

	return newOut
}

func strArr2IntArr(in []string) ([]int, error) {

	out := make([]int, len(in))
	for i, s := range in {
		v, err := strconv.Atoi(s)
		if err != nil {
			return nil, errors.Wrapf(err, "strArr2IntArr atoi error, %s", s)
		}

		out[i] = v
	}

	return out, nil
}

func NewObjectParser(heads, fields []string, line int) *ObjectParser {
	p := &ObjectParser{
		dataMap: make(map[string][]string),
		line:    line,
	}

	for j := 0; j < len(fields); j++ {
		fields[j] = strings.TrimSpace(fields[j])
		fields[j] = strings.Trim(fields[j], "\"")

		key := strings.ToLower(heads[j])
		arr := p.dataMap[key]
		if len(arr) == 0 {
			p.dataMap[key] = []string{fields[j]}
		} else {
			p.dataMap[key] = append(arr, fields[j])
		}
	}

	return p
}

func ParseList(filename, fileContent string) ([]*ObjectParser, error) {

	if len(fileContent) == 0 {
		return make([]*ObjectParser, 0), nil
	}

	fileContent = deleteHeadRN(fileContent)

	as := strings.Split(fileContent, "\r\n")
	if len(as) <= 1 {
		as = strings.Split(fileContent, "\n")
		if len(as) <= 1 {
			as = strings.Split(fileContent, "\r")
			if len(as) <= 1 {
				return nil, errors.Errorf("%s 格式不正确，请复制一个正常文件，然后使用excel来编辑保存", filename)
			}
		}
	}

	heads := strings.Split(as[1], "\t")

	for i := 0; i < len(heads); i++ {
		heads[i] = strings.TrimSpace(heads[i])
		heads[i] = strings.Trim(heads[i], "\"")
	}

	parsers := make([]*ObjectParser, 0)

	for i := 2; i < len(as); i++ {
		if len(strings.TrimSpace(as[i])) == 0 {
			// empty line
			continue
		}

		line := i + 1
		fields := strings.Split(as[i], "\t")
		if len(heads) < len(fields) {
			return nil, errors.Errorf("%s 存在head之外的行，line: %s", filename, line)
		}

		parsers = append(parsers, NewObjectParser(heads, fields, line))
	}

	return parsers, nil
}

func deleteHeadRN(origin string) string {
	// 将双引号中间的换行符删掉
	array := strings.Split(origin, "\"")
	if len(array) <= 1 {
		return origin
	}

	buf := bytes.Buffer{}

	for i := 0; i < len(array); i++ {
		s := array[i]
		if len(s) == 0 {
			continue
		}

		if i%2 == 1 {
			s = strings.Replace(s, "\r", "", -1)
			s = strings.Replace(s, "\n", "", -1)
		}

		if i > 0 {
			buf.WriteString("\"")
		}

		buf.WriteString(s)
	}

	return buf.String()

}
