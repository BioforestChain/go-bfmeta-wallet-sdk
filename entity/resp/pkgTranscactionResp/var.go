package pkgTranscactionResp

import "github.com/BioforestChain/go-bfmeta-wallet-sdk/entity/jbase"

type SuccessPackageResult struct {
	Buffer string `json:"buffer"`
}

type PackageResult struct {
	Success bool                  `json:"success"`
	Result  SuccessPackageResult  `json:"result,omitempty"`
	Error   jbase.JsonCommonError `json:"error,omitempty"`
	MinFee  jbase.JsonMinFee      `json:"minFee,omitempty"`
}
