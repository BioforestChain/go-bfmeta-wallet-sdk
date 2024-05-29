package pkgTransferAssetResp

// SuccessPackageResult 表示成功的创建结果
type SuccessPackageResult struct {
	Success bool `json:"success"`
	Result  struct {
		Buffer string `json:"buffer"`
	} `json:"result"`
}

// FailurePackageResult 表示失败的创建结果
type FailurePackageResult struct {
	//Success bool `json:"success"`
	Error struct {
		Code    string `json:"code"`
		Message string `json:"message"`
	} `json:"error"`
	MinFee string `json:"minFee"`
}

// PackageResult 是 SuccessPackageResult 或 FailurePackageResult 的联合类型
type PackageResult struct {
	SuccessPackageResult
	FailurePackageResult
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
