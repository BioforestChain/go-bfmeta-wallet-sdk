package jbase

import (
	"encoding/json"
	"strconv"
)

// 定义一个自定义类型 JsonMinFee
type JsonMinFee string

// 实现自定义的 UnmarshalJSON 方法
func (s *JsonMinFee) UnmarshalJSON(data []byte) error {
	// 如果字段是一个字符串，直接解析
	if data[0] == '"' {
		return json.Unmarshal(data, (*string)(s))
	}

	// 如果字段是一个数字，将其转换为字符串
	var number float64
	if err := json.Unmarshal(data, &number); err != nil {
		return err
	}

	*s = JsonMinFee(strconv.FormatFloat(number, 'f', -1, 64))
	return nil
}
