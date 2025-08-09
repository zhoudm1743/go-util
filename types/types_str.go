package types

import (
	"fmt"
	"strconv"
	"strings"
	"unicode"
	"unicode/utf8"

	"golang.org/x/exp/constraints"
)

type XStr string

func Str(s string) XStr {
	return XStr(s)
}

func Strings(ss ...string) XStr {
	return XStr(strings.Join(ss, ""))
}

func StrX[T constraints.Ordered](v T) XStr {
	return XStr(fmt.Sprint(v))
}

func (s XStr) Len() int {
	return len(s)
}

func (s XStr) Size() int {
	return utf8.RuneCountInString(s.String())
}

func (s XStr) IsEmpty() bool {
	return s.Len() == 0
}

func (s XStr) Int() int {
	i, _ := strconv.Atoi(string(s))
	return i
}

func (s XStr) Int64() int64 {
	return int64(s.Int())
}

func (s XStr) Uint() uint {
	return uint(s.Int())
}

func (s XStr) Uint64() uint64 {
	return uint64(s.Int())
}

func (s XStr) Float() float64 {
	f, _ := strconv.ParseFloat(string(s), 64)
	return f
}

func (s XStr) Bool() bool {
	b, _ := strconv.ParseBool(string(s))
	return b
}

func (s XStr) String() string {
	return string(s)
}

func (s XStr) Bytes() []byte {
	return []byte(s)
}

// FirstUpper 首字母大写
func (s XStr) FirstUpper() XStr {
	if s.Len() == 0 {
		return s
	}
	runes := []rune(s)
	if unicode.IsLetter(runes[0]) {
		runes[0] = unicode.ToUpper(runes[0])
	}
	return XStr(runes)
}

// FirstLower 首字母小写
func (s XStr) FirstLower() XStr {
	if s.Len() == 0 {
		return s
	}
	runes := []rune(s)
	if unicode.IsLetter(runes[0]) {
		runes[0] = unicode.ToLower(runes[0])
	}
	return XStr(runes)
}

func (s XStr) LastUpper() XStr {
	if s.Len() == 0 {
		return s
	}
	runes := []rune(s)
	if unicode.IsLetter(runes[len(runes)-1]) {
		runes[len(runes)-1] = unicode.ToUpper(runes[len(runes)-1])
	}
	return XStr(runes)
}

func (s XStr) LastLower() XStr {
	if s.Len() == 0 {
		return s
	}
	runes := []rune(s)
	if unicode.IsLetter(runes[len(runes)-1]) {
		runes[len(runes)-1] = unicode.ToLower(runes[len(runes)-1])
	}
	return XStr(runes)
}

func (s XStr) Upper() XStr {
	if s.Len() == 0 {
		return s
	}
	runes := []rune(s)
	for i := len(runes) - 1; i > 0; i-- {
		if unicode.IsLetter(runes[i]) {
			runes[i] = unicode.ToUpper(runes[i])
		}
	}
	return XStr(runes)
}

func (s XStr) Lower() XStr {
	if s.Len() == 0 {
		return s
	}
	runes := []rune(s)
	for i := len(runes) - 1; i > 0; i-- {
		if unicode.IsLetter(runes[i]) {
			runes[i] = unicode.ToLower(runes[i])
		}
	}
	return XStr(runes)
}

// Camel2Snake 驼峰转蛇形
func (s XStr) Camel2Snake() string {
	runes := []rune(s.RemoveSpace())
	var ns strings.Builder
	for _, r := range runes {
		if unicode.IsUpper(r) {
			ns.WriteString("_")
		}
		ns.WriteRune(unicode.ToLower(r))
	}
	return XStr(ns.String()).Trim("_", TrimLeft)
}

func (s XStr) Snake2BigCamel() string {
	runes := []rune(s.RemoveSpace())

	if len(runes) == 0 {
		return ""
	}

	if len(runes) == 1 {
		return string(unicode.ToUpper(runes[0]))
	}

	var ns strings.Builder
	if runes[0] != '_' {
		ns.WriteRune(runes[0])
	}
	for i := 1; i < len(runes); i++ {
		if runes[i] == '_' {
			continue
		}

		if runes[i-1] == '_' {
			ns.WriteRune(unicode.ToUpper(runes[i]))
		} else {
			ns.WriteRune(unicode.ToLower(runes[i]))
		}
	}
	return XStr(ns.String()).FirstUpper().String()
}

func (s XStr) Snake2LittleCamel() string {
	runes := []rune(s.RemoveSpace())

	if len(runes) == 0 {
		return ""
	}

	if len(runes) == 1 {
		return string(unicode.ToLower(runes[0]))
	}

	var ns strings.Builder
	if runes[0] != '_' {
		ns.WriteRune(runes[0])
	}
	for i := 1; i < len(runes); i++ {
		if runes[i] == '_' {
			continue
		}

		if runes[i-1] == '_' {
			ns.WriteRune(unicode.ToUpper(runes[i]))
		} else {
			ns.WriteRune(unicode.ToLower(runes[i]))
		}
	}
	return XStr(ns.String()).FirstLower().String()
}

// Append 追加多个字符串到s尾部
func (s XStr) Append(ss ...string) XStr {
	if len(ss) == 0 {
		return s
	}

	var str strings.Builder
	str.WriteString(s.String())
	for i := 0; i < len(ss); i++ {
		str.WriteString(ss[i])
	}

	return XStr(str.String())
}

// Added 将 str 插入到 s 头部
func (s XStr) Added(str string) XStr {
	return XStr(strings.Join([]string{str, s.String()}, ""))
}

// Split 将s以sep字符分割为数组/切片
func (s XStr) Split(sep string) []string {
	return strings.Split(string(s), sep)
}

type TrimType int

const (
	TrimLeft TrimType = iota
	TrimRight
)

func (s XStr) Trim(cutest string, direction ...TrimType) string {
	if len(direction) == 1 {
		switch direction[0] {
		case TrimLeft:
			return strings.TrimLeft(string(s), cutest)
		case TrimRight:
			return strings.TrimRight(string(s), cutest)
		}
	}
	return strings.Trim(string(s), cutest)
}

func (s XStr) TrimSpace() string {
	return strings.TrimSpace(string(s))
}

func (s XStr) ReplaceAll(old, new string) string {
	return strings.ReplaceAll(string(s), old, new)
}

func (s XStr) RemoveSpace() string {
	return s.ReplaceAll(" ", "")
}
