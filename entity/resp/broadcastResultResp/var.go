package broadcastResultResp

import (
	"encoding/json"
	"strconv"
)

type SuccessBroadcastResult struct {
	Buffer string `json:"buffer"`
}

type FailureBroadcastResult struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

type BroadcastResult struct {
	Success bool                   `json:"success"`
	Result  SuccessBroadcastResult `json:"result,omitempty"`
	Error   FailureBroadcastResult `json:"error,omitempty"`
	MinFee  StringOrNumber         `json:"minFee,omitempty"`
}

// 定义一个自定义类型 StringOrNumber
type StringOrNumber string

// 实现自定义的 UnmarshalJSON 方法
func (s *StringOrNumber) UnmarshalJSON(data []byte) error {
	// 如果字段是一个字符串，直接解析
	if data[0] == '"' {
		return json.Unmarshal(data, (*string)(s))
	}

	// 如果字段是一个数字，将其转换为字符串
	var number float64
	if err := json.Unmarshal(data, &number); err != nil {
		return err
	}

	*s = StringOrNumber(strconv.FormatFloat(number, 'f', -1, 64))
	return nil
}
