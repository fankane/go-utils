package assert

import "reflect"

func IsNumeric(i interface{}) bool {
	if i == nil {
		return false
	}
	switch i.(type) {
	case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64, float32, float64:
		return true
	default:
		return false
	}
}

func IsInteger(i interface{}) bool {
	if i == nil {
		return false
	}
	switch i.(type) {
	case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64:
		return true
	default:
		return false
	}
	return false
}

func IsFloat(i interface{}) bool {
	if i == nil {
		return false
	}
	switch i.(type) {
	case float32, float64:
		return true
	default:
		return false
	}
}

func IsArrayOrSlice(i interface{}) bool {
	if i == nil {
		return false
	}
	return reflect.TypeOf(i).Kind() == reflect.Array || reflect.TypeOf(i).Kind() == reflect.Slice
}
