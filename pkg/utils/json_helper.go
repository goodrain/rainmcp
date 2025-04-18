package utils

import (
	"encoding/json"
	"reflect"
	"strings"
)

// FieldInfo 表示字段信息，包含字段名称和描述
type FieldInfo struct {
	Name        string      `json:"name"`
	Description string      `json:"description,omitempty"`
	Type        string      `json:"type"`
	Value       interface{} `json:"value"`
}

// ObjectWithDescription 表示带有字段描述的对象
type ObjectWithDescription struct {
	Fields []FieldInfo `json:"fields"`
}

// MarshalJSONWithDescription 将对象序列化为包含字段描述的JSON
func MarshalJSONWithDescription(obj interface{}) ([]byte, error) {
	val := reflect.ValueOf(obj)
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}

	// 如果是切片或数组，处理每个元素
	if val.Kind() == reflect.Slice || val.Kind() == reflect.Array {
		result := make([]interface{}, val.Len())
		for i := 0; i < val.Len(); i++ {
			item := val.Index(i).Interface()
			if reflect.ValueOf(item).Kind() == reflect.Struct || 
			   (reflect.ValueOf(item).Kind() == reflect.Ptr && reflect.ValueOf(item).Elem().Kind() == reflect.Struct) {
				desc, err := structToObjectWithDescription(item)
				if err != nil {
					return nil, err
				}
				result[i] = desc
			} else {
				result[i] = item
			}
		}
		return json.MarshalIndent(result, "", "  ")
	}

	// 处理单个结构体
	if val.Kind() == reflect.Struct {
		return structToJSONWithDescription(obj)
	}

	// 处理map或其他类型
	return json.MarshalIndent(obj, "", "  ")
}

// 将结构体转换为带有字段描述的JSON
func structToJSONWithDescription(obj interface{}) ([]byte, error) {
	desc, err := structToObjectWithDescription(obj)
	if err != nil {
		return nil, err
	}
	return json.MarshalIndent(desc, "", "  ")
}

// 将结构体转换为带有字段描述的对象
func structToObjectWithDescription(obj interface{}) (map[string]interface{}, error) {
	val := reflect.ValueOf(obj)
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}
	typ := val.Type()

	result := make(map[string]interface{})

	for i := 0; i < val.NumField(); i++ {
		field := typ.Field(i)
		
		// 跳过未导出的字段
		if field.PkgPath != "" {
			continue
		}

		fieldVal := val.Field(i)
		jsonTag := field.Tag.Get("json")
		if jsonTag == "-" {
			continue
		}

		// 提取JSON字段名
		jsonName := field.Name
		parts := strings.Split(jsonTag, ",")
		if parts[0] != "" {
			jsonName = parts[0]
		}

		// 提取描述
		description := field.Tag.Get("description")

		// 处理嵌套结构体
		var fieldValue interface{}
		if fieldVal.Kind() == reflect.Struct || 
		   (fieldVal.Kind() == reflect.Ptr && !fieldVal.IsNil() && fieldVal.Elem().Kind() == reflect.Struct) {
			nestedObj, err := structToObjectWithDescription(fieldVal.Interface())
			if err != nil {
				return nil, err
			}
			fieldValue = nestedObj
		} else if fieldVal.Kind() == reflect.Slice || fieldVal.Kind() == reflect.Array {
			// 处理切片或数组
			sliceResult := make([]interface{}, fieldVal.Len())
			for j := 0; j < fieldVal.Len(); j++ {
				item := fieldVal.Index(j).Interface()
				if reflect.ValueOf(item).Kind() == reflect.Struct || 
				   (reflect.ValueOf(item).Kind() == reflect.Ptr && reflect.ValueOf(item).Elem().Kind() == reflect.Struct) {
					desc, err := structToObjectWithDescription(item)
					if err != nil {
						return nil, err
					}
					sliceResult[j] = desc
				} else {
					sliceResult[j] = item
				}
			}
			fieldValue = sliceResult
		} else {
			// 其他类型直接使用值
			fieldValue = fieldVal.Interface()
		}

		// 添加字段信息
		fieldInfo := map[string]interface{}{
			"value": fieldValue,
		}
		
		if description != "" {
			fieldInfo["description"] = description
		}
		
		result[jsonName] = fieldInfo
	}

	return result, nil
}
