package slice

import "strconv"

// ToInterfaceSli
/**
例如
[]int{} -> []interface{}
[]str{} -> []interface{}
*/
func ToInterfaceSli(slice interface{}) []interface{} {
	result := make([]interface{}, 0)
	switch slice.(type) {
	case []int:
		for _, t := range slice.([]int) {
			result = append(result, t)
		}
		return result
	case []int8:
		for _, t := range slice.([]int8) {
			result = append(result, t)
		}
		return result
	case []int16:
		for _, t := range slice.([]int16) {
			result = append(result, t)
		}
		return result
	case []int32:
		for _, t := range slice.([]int32) {
			result = append(result, t)
		}
		return result
	case []int64:
		for _, t := range slice.([]int64) {
			result = append(result, t)
		}
		return result
	case []uint:
		for _, t := range slice.([]uint) {
			result = append(result, t)
		}
		return result
	case []uint8:
		for _, t := range slice.([]uint8) {
			result = append(result, t)
		}
		return result
	case []uint16:
		for _, t := range slice.([]uint16) {
			result = append(result, t)
		}
		return result
	case []uint32:
		for _, t := range slice.([]uint32) {
			result = append(result, t)
		}
		return result
	case []uint64:
		for _, t := range slice.([]uint64) {
			result = append(result, t)
		}
		return result
	case []float32:
		for _, t := range slice.([]float32) {
			result = append(result, t)
		}
		return result
	case []float64:
		for _, t := range slice.([]float64) {
			result = append(result, t)
		}
		return result
	case []string:
		for _, t := range slice.([]string) {
			result = append(result, t)
		}
		return result
	}
	return nil
}

// StrSliToIntSli 字符串切片转 int 切片，如果元素不是int类型，则跳过，不报错
func StrSliToIntSli(slice []string) []int {
	intSli := make([]int, 0, len(slice))
	for _, s := range slice {
		i, err := strconv.Atoi(s)
		if err != nil {
			continue
		}
		intSli = append(intSli, i)
	}
	return intSli
}

// StrSliToFloat64Sli 字符串切片转 float64 切片，如果元素不是 float64 类型，则跳过，不报错
func StrSliToFloat64Sli(slice []string) []float64 {
	intSli := make([]float64, 0, len(slice))
	for _, s := range slice {
		i, err := strconv.ParseFloat(s, 64)
		if err != nil {
			continue
		}
		intSli = append(intSli, i)
	}
	return intSli
}

func StrSliToMap(slice []string) map[string]struct{} {
	res := make(map[string]struct{})
	for _, s := range slice {
		res[s] = struct{}{}
	}
	return res
}

// StrSliToMapCnt 字符串切片转map， value是key出现的次数
func StrSliToMapCnt(slice []string) map[string]int {
	res := make(map[string]int)
	for _, s := range slice {
		if _, ok := res[s]; ok {
			res[s]++
		} else {
			res[s] = 1
		}
	}
	return res
}
