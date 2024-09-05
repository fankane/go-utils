package gql

import (
	"fmt"
	"reflect"
	"strings"
	"time"
)

const jsonTagIgnore = "-"

// GetTableBasicFieldsByStruct 获取基础类型的表字段
// 如果有json标签，则名称使用json标签里面的；否则使用字段名；JSON tag - 表示忽略此字段
func GetTableBasicFieldsByStruct(obj interface{}) (map[string]ValidType, error) {
	rv := reflect.ValueOf(obj)
	if rv.Kind() == reflect.Ptr {
		rv = rv.Elem() // 如果传入的是指针，则解引用
	}
	if rv.Kind() != reflect.Struct {
		return nil, fmt.Errorf("expected a struct, got %v", rv.Kind())
	}
	result := make(map[string]ValidType)
	for i := 0; i < rv.NumField(); i++ {
		field := rv.Type().Field(i)
		jsonTag := field.Tag.Get("json")
		if jsonTag == jsonTagIgnore {
			continue
		}
		name := field.Name
		if strings.TrimSpace(jsonTag) != "" {
			name = jsonTag
		}
		baseVT := baseFieldParse(field)
		if baseVT != nil {
			result[name] = baseVT
		} else if field.Type.Kind() == reflect.Slice || field.Type.Kind() == reflect.Array {
			sliceVT := baseFieldParse(field.Type.Elem().Field(0))
			if sliceVT != nil {
				listVT := baseToList(sliceVT)
				if listVT != nil {
					result[name] = listVT
				}
			}
		}
	}
	return result, nil
}

func baseFieldParse(field reflect.StructField) ValidType {
	if field.Type == reflect.TypeOf(time.Time{}) {
		return DateTime
	}
	switch field.Type.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return Int
	case reflect.String:
		return String
	case reflect.Float32, reflect.Float64:
		return Float
	case reflect.Bool:
		return Boolean
	default:
		return nil
	}
}

func baseToList(v ValidType) ValidType {
	switch v {
	case Int:
		return ListInt
	case Float:
		return ListFloat
	case Boolean:
		return ListBoolean
	case String:
		return ListString
	case DateTime:
		return ListDateTime
	default:
		return nil
	}
}
