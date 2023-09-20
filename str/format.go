package str

import (
	"fmt"
	"unicode/utf8"
)

const StepSize = 1024.0

var sufIdx = []string{
	"B",
	"KB",
	"MB",
	"GB",
	"TB",
	"PB",
	"EB",
	"ZB",
	"YB",
}

// FormatFileSize 按所给字节数转换相应存储单位 保留两位小数点
func FormatFileSize(bytes float64) string {
	for i := 0; i < len(sufIdx); i++ {
		if bytes < StepSize {
			return fmt.Sprintf("%.2f %s", bytes, sufIdx[i])
		}
		if i == len(sufIdx)-1 {
			break
		}
		bytes = bytes / StepSize
	}
	return fmt.Sprintf("%.2f %s", bytes, sufIdx[len(sufIdx)-1])
}

// GetStrIndex 获取字符串指定下标的字符, idx 从0开始计算
func GetStrIndex(s string, idx int) string {
	if len(s) <= idx {
		return ""
	}
	// 先将 s 转为 []rune ,防止例如中文下的 s[idx:idx+1] 乱码
	return string([]rune(s)[idx : idx+1])
}

// SliceOfChar 将字符串每个字符抽出来，组成字符串数组 eg: "hi中国" -> ["h", "i", "中", "国"]
func SliceOfChar(s string) []string {
	res := make([]string, 0)
	for _, c := range s {
		res = append(res, fmt.Sprintf("%c", c))
	}
	return res
}

// LenOfUTF8 字符串在utf8编码下的长度，一个中文算1
func LenOfUTF8(s string) int {
	return utf8.RuneCountInString(s)
}

// SubOfUTF8 带中文字符串切割
func SubOfUTF8(s string, start, end int) string {
	if start < 0 || start >= end || start >= LenOfUTF8(s) || end > LenOfUTF8(s) {
		return ""
	}
	return string([]rune(s)[start:end])
}
