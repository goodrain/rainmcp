package utils

import (
	"encoding/json"
	"fmt"
)

// FormatJSON 将任意类型的数据格式化为格式化的 JSON 字符串
func FormatJSON(data interface{}) string {
	// 如果数据是字节切片，尝试解析为 JSON
	if bytes, ok := data.([]byte); ok {
		var jsonData interface{}
		if err := json.Unmarshal(bytes, &jsonData); err == nil {
			data = jsonData
		}
	}

	// 将数据转换为格式化的 JSON
	result, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		// 如果格式化失败，返回原始数据的字符串表示
		return fmt.Sprintf("%v", data)
	}
	return string(result)
}
