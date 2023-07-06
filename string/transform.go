package string

import (
	"encoding/json"
	"strconv"
	"unsafe"
)

// StrToInt 字符串转数字，错误则返回默认值
func StrToInt(str string, defaultVal int) int {
	i, err := strconv.Atoi(str)
	if err != nil {
		return defaultVal
	}
	return i
}

// StrToFloat64 字符串转数字，错误则返回默认值
func StrToFloat64(str string, defaultVal float64) float64 {
	i, err := strconv.ParseFloat(str, 64)
	if err != nil {
		return defaultVal
	}
	return i
}

func ToJSON(a interface{}) string {
	b, err := json.Marshal(a)
	if err != nil {
		return ""
	}
	return BytesToStr(b)
}

// StrToBytes 字符串转byte数组, 使用unsafe.Pointer来转换不同类型的指针，没有底层数据的拷贝
func StrToBytes(s string) []byte {
	tmp1 := (*[2]uintptr)(unsafe.Pointer(&s))
	tmp2 := [3]uintptr{tmp1[0], tmp1[1], tmp1[1]}
	return *(*[]byte)(unsafe.Pointer(&tmp2))
}

// BytesToStr byte数组转字符串, 使用unsafe.Pointer来转换不同类型的指针，没有底层数据的拷贝
func BytesToStr(bytes []byte) string {
	return *(*string)(unsafe.Pointer(&bytes))
}
