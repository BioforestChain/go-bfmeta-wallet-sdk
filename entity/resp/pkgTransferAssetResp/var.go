package pkgTransferAssetResp

import "github.com/BioforestChain/go-bfmeta-wallet-sdk/entity/jbase"

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
	Error  jbase.JsonCommonError `json:"error"`
	MinFee jbase.JsonMinFee      `json:"minFee"`
}

// PackageResult 是 SuccessPackageResult 或 FailurePackageResult 的联合类型
type PackageResult struct {
	SuccessPackageResult
	FailurePackageResult
}

// ApiFailureReturn 表示 API 调用失败的返回
type ApiFailureReturn struct {
	Success bool                  `json:"success"`
	Error   jbase.JsonCommonError `json:"error"`
}
