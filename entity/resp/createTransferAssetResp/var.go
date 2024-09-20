package createTransferAssetResp

import (
	"github.com/BioforestChain/go-bfmeta-wallet-sdk/entity/jbase"
)

// SuccessCreateResult 表示成功的创建结果
type SuccessCreateResult struct {
	Buffer jbase.Base64StringBuffer `json:"buffer"`
}

// // CreateResult 是 SuccessCreateResult 或 FailureCreateResult 的联合类型
type CreateResult struct {
	Success bool                  `json:"success"`
	Result  SuccessCreateResult   `json:"result,omitempty"`
	Error   jbase.JsonCommonError `json:"error,omitempty"`
	MinFee  jbase.JsonMinFee      `json:"minFee,omitempty"`
}

// ApiFailureReturn 表示 API 调用失败的返回
type ApiFailureReturn struct {
	Success bool                  `json:"success"`
	Error   jbase.JsonCommonError `json:"error"`
}
