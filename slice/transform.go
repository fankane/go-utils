package slice

// SliToInterfaces
/**
例如
[]int{} -> []interface{}
[]string{} -> []interface{}
*/
func SliToInterfaces(slice interface{}) []interface{} {
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
