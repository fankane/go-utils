package str

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
)

const (
	JPObject = "object"
	JPArray  = "array"
)

type JsonProperty struct {
	Type       string                   `json:"type"`                 // object, array, int ,float 等
	Properties map[string]*JsonProperty `json:"properties,omitempty"` // type = object  时有值
	Items      *JsonProperty            `json:"items,omitempty"`      // type = array 时有值
	ItemLen    int                      `json:"item_len,omitempty"`   // type = array 时有值
}

func ParseJSONProperty(jsonStr string) (*JsonProperty, error) {
	if !json.Valid([]byte(jsonStr)) {
		return nil, fmt.Errorf("not json string")
	}
	var jsonData interface{}
	err := json.Unmarshal([]byte(jsonStr), &jsonData)
	if err != nil {
		return nil, fmt.Errorf("unmarshal json err:%s", err)
	}

	valOf := reflect.ValueOf(jsonData)
	jsonP := &JsonProperty{}
	if err = parseObject(valOf, jsonP); err != nil {
		return nil, err
	}
	return jsonP, nil
}

func parseObject(valOf reflect.Value, result *JsonProperty) error {
	switch valOf.Kind() {
	case reflect.Map:
		result.Type = JPObject
		result.Properties = make(map[string]*JsonProperty)
		return parseMapObject(valOf, result)
	case reflect.Slice:
		result.Type = JPArray
		result.Items = &JsonProperty{}
		result.ItemLen = valOf.Len()
		return parseSliceObject(valOf, result.Items)
	case reflect.String, reflect.Bool:
		result.Type = valOf.Kind().String()
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		result.Type = "int"
	case reflect.Float32, reflect.Float64:
		if !strings.Contains(fmt.Sprintf("%v", valOf.Interface()), ".") {
			result.Type = "int"
			return nil
		}
		result.Type = "float"
	default:
		return fmt.Errorf("parseObject unsupport kind:%s, key:%s", valOf.Kind(), valOf.Interface())
	}
	return nil
}

func parseMapObject(valOf reflect.Value, result *JsonProperty) error {
	keys := valOf.MapKeys()
	for _, key := range keys {
		val := reflect.ValueOf(valOf.MapIndex(key).Interface())
		switch val.Kind() {
		case reflect.Map:
			tempPro := &JsonProperty{
				Type:       JPObject,
				Properties: map[string]*JsonProperty{},
			}
			result.Properties[key.Interface().(string)] = tempPro
			parseMapObject(val, tempPro)
		case reflect.Slice:
			tempPro := &JsonProperty{
				Type:    JPArray,
				Items:   &JsonProperty{},
				ItemLen: val.Len(),
			}
			result.Properties[key.Interface().(string)] = tempPro
			parseSliceObject(val, tempPro.Items)
		case reflect.String, reflect.Bool:
			result.Properties[key.Interface().(string)] = &JsonProperty{Type: val.Kind().String()}
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			result.Properties[key.Interface().(string)] = &JsonProperty{Type: "int"}
		case reflect.Float32, reflect.Float64:
			if !strings.Contains(fmt.Sprintf("%v", val.Interface()), ".") {
				result.Properties[key.Interface().(string)] = &JsonProperty{Type: "int"}
				continue
			}
			result.Properties[key.Interface().(string)] = &JsonProperty{Type: "float"}
		default:
			return fmt.Errorf("parseMapObject unsupport kind:%s, key:%s", val.Kind(), key.Interface())
		}
	}
	return nil
}

func parseSliceObject(valOf reflect.Value, result *JsonProperty) error {
	if valOf.Len() == 0 {
		return fmt.Errorf("empty slice")
	}
	//只读取数组第一个元素的结构
	firstEle := valOf.Index(0)
	return parseObject(reflect.ValueOf(firstEle.Interface()), result)
}
