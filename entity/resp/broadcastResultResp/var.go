package broadcastResultResp

import (
	"github.com/BioforestChain/go-bfmeta-wallet-sdk/entity/jbase"
)

type SuccessBroadcastResult struct {
	Buffer string `json:"buffer"`
}

type BroadcastResult struct {
	Success bool                   `json:"success"`
	Result  SuccessBroadcastResult `json:"result,omitempty"`
	Error   jbase.JsonCommonError  `json:"error,omitempty"`
	MinFee  jbase.JsonMinFee       `json:"minFee,omitempty"`
}
