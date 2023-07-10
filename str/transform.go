package str

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"strconv"
	"unsafe"
)

// ToInt 字符串转数字，错误则返回默认值
func ToInt(str string, defaultVal int) int {
	i, err := strconv.Atoi(str)
	if err != nil {
		return defaultVal
	}
	return i
}

// ToFloat64 字符串转数字，错误则返回默认值
func ToFloat64(str string, defaultVal float64) float64 {
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
	return FromBytes(b)
}

// ToBytes 字符串转byte数组, 使用unsafe.Pointer来转换不同类型的指针，没有底层数据的拷贝
func ToBytes(s string) []byte {
	tmp1 := (*[2]uintptr)(unsafe.Pointer(&s))
	tmp2 := [3]uintptr{tmp1[0], tmp1[1], tmp1[1]}
	return *(*[]byte)(unsafe.Pointer(&tmp2))
}

// FromBytes byte数组转字符串, 使用unsafe.Pointer来转换不同类型的指针，没有底层数据的拷贝
func FromBytes(bytes []byte) string {
	return *(*string)(unsafe.Pointer(&bytes))
}

// MD5 获取字符串MD5值
func MD5(str string) string {
	h := md5.New()
	_, err := h.Write(ToBytes(str))
	if err != nil {
		return ""
	}
	return hex.EncodeToString(h.Sum(nil))
}
