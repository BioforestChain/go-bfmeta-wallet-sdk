package pkgTranscactionResp

type SuccessPackageResult struct {
	Buffer string `json:"buffer"`
}

type FailurePackageResult struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

type PackageResult struct {
	Success bool                 `json:"success"`
	Result  SuccessPackageResult `json:"result,omitempty"`
	Error   FailurePackageResult `json:"error,omitempty"`
	MinFee  int                  `json:"minFee,omitempty"`
}
