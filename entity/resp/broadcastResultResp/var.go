package broadcastResultResp

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
	MinFee  string                 `json:"minFee,omitempty"`
}
