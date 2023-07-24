package common

import (
	"fmt"

	jsoniter "github.com/json-iterator/go"
)

// Struct2Json 结构体转为json
func Struct2Json(obj interface{}) string {
	str, err := jsoniter.Marshal(obj)
	if err != nil {
		panic(fmt.Sprintf("[Struct2Json]转换异常: %v", err))
	}
	return string(str)
}

// Json2Struct json转为结构体
func Json2Struct(str string, obj interface{}) {
	// 将json转为结构体
	err := jsoniter.Unmarshal([]byte(str), obj)
	if err != nil {
		panic(fmt.Sprintf("[Json2Struct]转换异常: %v", err))
	}
}

// JsonI2Struct json interface转为结构体
func JsonI2Struct(str interface{}, obj interface{}) {
	JsonStr := str.(string)
	Json2Struct(JsonStr, obj)
}
