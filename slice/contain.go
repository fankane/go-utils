package slice

func InInterfaceSlice(target interface{}, slice []interface{}) bool {
	for _, t := range slice {
		if t == target {
			return true
		}
	}
	return false
}

// InInterfaces 数组里存在某个值
func InInterfaces(target interface{}, slice interface{}) bool {
	switch target.(type) {
	case int:
		if sli, ok := slice.([]int); ok {
			return InInts(target.(int), sli)
		}
	case int8:
		if sli, ok := slice.([]int8); ok {
			return InInt8s(target.(int8), sli)
		}
	case int16:
		if sli, ok := slice.([]int16); ok {
			return InInt16s(target.(int16), sli)
		}
	case int32:
		if sli, ok := slice.([]int32); ok {
			return InInt32s(target.(int32), sli)
		}
	case int64:
		if sli, ok := slice.([]int64); ok {
			return InInt64s(target.(int64), sli)
		}
	case uint:
		if sli, ok := slice.([]uint); ok {
			return InUints(target.(uint), sli)
		}
	case uint8:
		if sli, ok := slice.([]uint8); ok {
			return InUint8s(target.(uint8), sli)
		}
	case uint16:
		if sli, ok := slice.([]uint16); ok {
			return InUint16s(target.(uint16), sli)
		}
	case uint32:
		if sli, ok := slice.([]uint32); ok {
			return InUint32s(target.(uint32), sli)
		}
	case uint64:
		if sli, ok := slice.([]uint64); ok {
			return InUint64s(target.(uint64), sli)
		}
	case string:
		if sli, ok := slice.([]string); ok {
			return InStrings(target.(string), sli)
		}
	case float32:
		if sli, ok := slice.([]float32); ok {
			return InFloat32(target.(float32), sli)
		}
	case float64:
		if sli, ok := slice.([]float64); ok {
			return InFloat64(target.(float64), sli)
		}
	}
	return false
}

// InStrings 数组里存在某个值
func InStrings(target string, slice []string) bool {
	for _, s := range slice {
		if target == s {
			return true
		}
	}
	return false
}

// InInts 数组里存在某个值
func InInts(target int, slice []int) bool {
	for _, s := range slice {
		if target == s {
			return true
		}
	}
	return false
}

// InInt8s 数组里存在某个值
func InInt8s(str int8, slice []int8) bool {
	for _, s := range slice {
		if str == s {
			return true
		}
	}
	return false
}

// InInt16s 数组里存在某个值
func InInt16s(str int16, slice []int16) bool {
	for _, s := range slice {
		if str == s {
			return true
		}
	}
	return false
}

// InInt32s 数组里存在某个值
func InInt32s(str int32, slice []int32) bool {
	for _, s := range slice {
		if str == s {
			return true
		}
	}
	return false
}

// InInt64s 数组里存在某个值
func InInt64s(target int64, slice []int64) bool {
	for _, s := range slice {
		if target == s {
			return true
		}
	}
	return false
}

// InUints 数组里存在某个值
func InUints(target uint, slice []uint) bool {
	for _, s := range slice {
		if target == s {
			return true
		}
	}
	return false
}

// InUint8s 数组里存在某个值
func InUint8s(str uint8, slice []uint8) bool {
	for _, s := range slice {
		if str == s {
			return true
		}
	}
	return false
}

// InUint16s 数组里存在某个值
func InUint16s(str uint16, slice []uint16) bool {
	for _, s := range slice {
		if str == s {
			return true
		}
	}
	return false
}

// InUint32s 数组里存在某个值
func InUint32s(str uint32, slice []uint32) bool {
	for _, s := range slice {
		if str == s {
			return true
		}
	}
	return false
}

// InUint64s 数组里存在某个值
func InUint64s(target uint64, slice []uint64) bool {
	for _, s := range slice {
		if target == s {
			return true
		}
	}
	return false
}

// InFloat32 数组里存在某个值
func InFloat32(str float32, slice []float32) bool {
	for _, s := range slice {
		if str == s {
			return true
		}
	}
	return false
}

// InFloat64 数组里存在某个值
func InFloat64(target float64, slice []float64) bool {
	for _, s := range slice {
		if target == s {
			return true
		}
	}
	return false
}
