package config

import (
	"github.com/pkg/errors"
	"math/big"
	"regexp"
	"strconv"
	"strings"
)

func (p *ObjectParser) IsEmpty(validate map[string]*Validator) bool {
	for k := range validate {
		// 过滤所有空值的行
		sa := p.OriginStringArray(k)
		for _, v := range sa {
			if len(v) > 0 {
				return false
			}
		}
	}

	return true
}

func (p *ObjectParser) Validate(filename string, validate map[string]*Validator) error {

	for k, v := range validate {
		sa := p.OriginStringArray(k)

		err := v.CheckRow(sa, filename, k, p.line)
		if err != nil {
			return err
		}
	}

	return nil
}

func (v *Validator) Check(s string, filename, k string) error {
	return v.CheckRow([]string{s}, filename, k, 0)
}

func (v *Validator) CheckRow(sa []string, filename, k string, line int) error {
	if len(sa) == 0 {
		// 如果有默认值，跳过
		if len(v.defaults) <= 0 {
			return errors.Errorf("校验配置文件[%s](行数: %v)，字段[%s]不存在", filename, line, k)
		}

		if line > 3 {
			// 行数从3开始，第一行检查默认值
			return nil
		}

		sa = v.defaults
	}

	if !v.array {
		if len(sa) > 1 {
			return errors.Errorf("校验配置文件[%s](行数: %v)，字段[%s]只允许存在一个（实际个数：%v）, %s", filename, line, k, len(sa), sa)
		}

		if v.separator != "" {
			sa = strings.Split(sa[0], v.separator)
		}
	}

	if !v.duplicate && Duplicate(sa) {
		return errors.Errorf("校验配置文件[%s](行数: %v)，字段[%s]配置的值不允许重复, %s", filename, line, k, sa)
	}

	if v.count > 0 && len(sa) != v.count {
		return errors.Errorf("校验配置文件[%s](行数: %v)，字段[%s]配置的值个数不匹配，必须是[%v]个，实际是[%v]个, %s", filename, line, k, v.count, len(sa), sa)
	}

	hasAnyNil := false
	hasAnyNotNil := false
	sum := new(big.Float)
	for _, s := range sa {
		if !v.notAllNil && !v.allNilOrNot {
			if v.notNil && len(s) == 0 {
				return errors.Errorf("校验配置文件[%s](行数: %v)，字段[%s]不允许配置空值, %s", filename, line, k, sa)
			}
		}

		if len(s) > 0 {
			hasAnyNotNil = true
			if v.regexp != nil {
				if !v.regexp.MatchString(s) {
					return errors.Errorf("校验配置文件[%s](行数: %v)，字段[%s]必须是%s [pattern: %s], %s", filename, line, k, v.tips, v.regexpPattern, sa)
				}
			}

			if len(v.whiteList) > 0 {
				inWhiteList := false
				for _, w := range v.whiteList {
					s1 := s
					s2 := w
					if !v.matchCase {
						s1 = strings.ToUpper(s1)
						s2 = strings.ToUpper(s2)
					}

					if s1 == s2 {
						inWhiteList = true
						break
					}
				}

				if !inWhiteList {
					return errors.Errorf("校验配置文件[%s](行数: %v)，字段[%s]必须配置下面列表中的任意一个(大小写敏感: %v), 当前配置: %s 允许列表: %s", filename, line, k, v.matchCase, sa, v.whiteList)
				}
			}

			if v.sum != nil {
				f, _, err := new(big.Float).Parse(s, 10)
				if err != nil {
					return errors.Wrapf(err, "校验配置文件[%s](行数: %v)，字段[%s]配置了求和，因此必须是数字, 当前配置: %s", filename, line, k, sa)
				}
				sum.Add(sum, f)
			}
		} else {
			hasAnyNil = true
		}
	}

	if v.sum != nil && v.sum.Cmp(sum) != 0 {
		return errors.Errorf("校验配置文件[%s](行数: %v)，字段[%s]配置了求和，但是求和结果不一致, %s 相加得到:%v sum: %v", filename, line, k, sa, sum, v.sum)
	}

	if v.notAllNil && !hasAnyNotNil {
		return errors.Errorf("校验配置文件[%s](行数: %v)，字段[%s]不允许全部配置空值, %s", filename, line, k, sa)
	}

	if v.allNilOrNot && hasAnyNil && hasAnyNotNil {
		return errors.Errorf("校验配置文件[%s](行数: %v)，字段[%s]不允许一部分配置，另一部分不配置（要么全部配置，或者全部不配置）, %s", filename, line, k, sa)
	}

	return nil
}

func Duplicate(in []string) bool {
	for s, x := range in {
		if len(x) == 0 {
			continue
		}

		for i := s + 1; i < len(in); i++ {
			if x == in[i] {
				return true
			}
		}
	}
	return false
}

type Validator struct {
	array bool // true表示字段可以多个

	duplicate bool // true表示允许重复，默认不允许

	notNil bool // true表示不允许空值，默认允许

	notAllNil bool // true表示至少配置一个

	matchCase bool // true表示大小写必须匹配，默认忽略大小写

	count       int
	allNilOrNot bool // true表示允许全部为nil或者全部不为nil

	sum *big.Float

	regexp        *regexp.Regexp // 正则表达式校验
	regexpPattern string
	tips          string

	whiteList []string
	defaults  []string

	separator string // 分隔符

	mapEntitySeparator string // 二级分隔符
}

func newRegexpValidator(pattern, tips string) *Validator {
	return newRegexpValidator0(pattern, tips, false)
}

func newNotNilRegexpValidator(pattern, tips string) *Validator {
	return newRegexpValidator0(pattern, tips, true)
}

func newRegexpValidator0(pattern, tips string, notNil bool) *Validator {
	return &Validator{
		regexp:        regexp.MustCompile(pattern),
		regexpPattern: pattern,
		tips:          tips,
		notNil:        notNil,
	}
}

func ParseValidator(str, sep string, array bool, wlist []string, defaults []string) *Validator {

	parts := strings.Split(strings.ToLower(str), ",")
	if len(parts) == 0 {

		if len(wlist) == 0 && len(defaults) == 0 {
			return stringValidator
		}

		out := *stringValidator
		out.whiteList = wlist
		out.defaults = defaults
		return &out
	}

	out := validatorMap[key(parts[0], sep, array)]
	if out != nil {
		if len(parts) == 1 && len(wlist) == 0 && len(defaults) == 0 {
			return out
		}

		// copy
		t := *out
		out = &t
	} else {
		out = newRegexpValidator(parts[0], parts[0])

		out.array = array
		if len(sep) > 0 {
			out.separator = sep
			out.array = false
		}
	}

	parseParts(out, parts[1:]...)

	out.whiteList = wlist
	out.defaults = defaults

	return out
}

func parseParts(t *Validator, parts ...string) *Validator {
	for _, opt := range parts {
		switch opt {
		case "d", "dup", "duplicate":
			t.duplicate = true
		case "notnil", "not nil", "notnull", "not null":
			t.notNil = true
		case "notallnil", "not all nil", "notallnull", "not all null":
			t.notAllNil = true
		case "allnilornot":
			t.allNilOrNot = true
		case "case":
			t.matchCase = true
		default:
			switch {
			case strings.HasPrefix(opt, "tips"):
				t.tips = opt[5:]
			case strings.HasPrefix(opt, "count"):
				t.count, _ = strconv.Atoi(opt[6:])
			case strings.HasPrefix(opt, "sum"):
				t.sum, _, _ = new(big.Float).Parse(opt[4:], 10)
			}
		}
	}

	return t
}

func key(name, sep string, array bool) string {

	out := name

	if len(sep) > 0 {
		if sep == "array" {
			panic("sep == \"array\"")
		}

		out += "_" + sep
	} else if array {
		out += "_array"
	}

	return out
}

var (
	validatorMap = map[string]*Validator{
		"int":        intValidator,
		"uint":       uintValidator,
		"int>=0":     uintValidator,
		">=0":        uintValidator,
		"int>0":      uintLt0Validator,
		">0":         uintLt0Validator,
		"string":     stringValidator,
		"string>0":   stringNotEmptyValidator,
		"float64":    float64Validator,
		"float64>=0": float64Le0Validator,
		"float64>0":  float64Lt0Validator,
		"bool":       boolValidator,
	}

	stringValidator         = &Validator{}
	stringNotEmptyValidator = &Validator{notNil: true}
	intValidator            = newRegexpValidator("^(0|[-]?[1-9][0-9]*)$", "整数（不支持小数）")
	uintValidator           = newRegexpValidator("^(0|[1-9][0-9]*)$", "正整数（不支持小数）>=0")
	uintLt0Validator        = newNotNilRegexpValidator("^([1-9][0-9]*)$", "正整数（不支持小数）>0")
	float64Validator        = newRegexpValidator(`^[-]?\d+(\.\d+)?$`, "数字（支持小数）")
	float64Le0Validator     = newRegexpValidator(`^\d+(\.\d+)?$`, "正数（支持小数）>=0")
	float64Lt0Validator     = newNotNilRegexpValidator(`^(0\.\d*[1-9]|[1-9][0-9]*(\.\d*[1-9])?)$`, "正数（支持小数）>0")
	boolValidator           = newRegexpValidator(`^0|1|true|false$`, "布尔值，0|1|true|false")
)

func init() {
	keys := make([]string, 0, len(validatorMap))
	for k, _ := range validatorMap {
		keys = append(keys, k)
	}

	for _, name := range keys {
		vd := validatorMap[name]
		arrVd := *vd
		arrVd.array = true
		validatorMap[key(name, "", true)] = &arrVd

		sep := ";"
		sepVd := *vd
		sepVd.separator = sep
		validatorMap[key(name, sep, false)] = &sepVd

	}
}
