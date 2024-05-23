package assert

import (
	"reflect"
)

var integerKinds = []reflect.Kind{
	reflect.Int,
	reflect.Int8,
	reflect.Int16,
	reflect.Int32,
	reflect.Int64,
	reflect.Uint,
	reflect.Uint8,
	reflect.Uint16,
	reflect.Uint32,
	reflect.Uint64,
}

var floatKinds = []reflect.Kind{
	reflect.Float32,
	reflect.Float64,
}

var arraySliceKinds = []reflect.Kind{
	reflect.Array,
	reflect.Slice,
}

const (
	Integer = "integer"
	Float   = "float"
	AS      = "array_slice"
)

var typeKind = map[string][]reflect.Kind{
	Integer: integerKinds,
	Float:   floatKinds,
	AS:      arraySliceKinds,
}

func IsNumeric(i interface{}) bool {
	return IsInteger(i) || IsFloat(i)
}

func IsInteger(i interface{}) bool {
	return assert(i, Integer)
}

func IsFloat(i interface{}) bool {
	return assert(i, Float)
}

func IsArrayOrSlice(i interface{}) bool {
	return assert(i, AS)
}

// IsTargetKind 是kind类型，或者kind类型的指针
func IsTargetKind(i interface{}, kind reflect.Kind) bool {
	if i == nil {
		return false
	}
	return isTargetOrTargetPtr(i, []reflect.Kind{kind})
}

func assert(i interface{}, t string) bool {
	if i == nil {
		return false
	}
	kindList, ok := typeKind[t]
	if !ok {
		return false
	}
	return isTargetOrTargetPtr(i, kindList)
}

func isTargetOrTargetPtr(i interface{}, kindList []reflect.Kind) bool {
	if reflect.TypeOf(i).Kind() == reflect.Ptr {
		return inKinds(reflect.ValueOf(i).Elem().Kind(), kindList)
	}
	return inKinds(reflect.TypeOf(i).Kind(), kindList)
}

func inKinds(target reflect.Kind, slice []reflect.Kind) bool {
	for _, t := range slice {
		if t == target {
			return true
		}
	}
	return false
}
