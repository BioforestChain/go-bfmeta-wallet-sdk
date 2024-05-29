package createTransferAssetResp

// SuccessCreateResult 表示成功的创建结果
type SuccessCreateResult struct {
	Success bool `json:"success"`
	Result  struct {
		Buffer string `json:"buffer"`
	} `json:"result"`
}

// FailureCreateResult 表示失败的创建结果
type FailureCreateResult struct {
	Success bool `json:"success"`
	Error   struct {
		Code    string `json:"code"`
		Message string `json:"message"`
	} `json:"error"`
	MinFee string `json:"minFee"`
}

// CreateResult 是 SuccessCreateResult 或 FailureCreateResult 的联合类型
type CreateResult struct {
	SuccessCreateResult
	FailureCreateResult
}

// ApiFailureReturn 表示 API 调用失败的返回
type ApiFailureReturn struct {
	Success bool `json:"success"`
	Error   struct {
		Code        *string `json:"code,omitempty"`
		Message     string  `json:"message"`
		Description *string `json:"description,omitempty"`
	} `json:"error"`
}
