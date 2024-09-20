package resp

import "github.com/BioforestChain/go-bfmeta-wallet-sdk/entity/jbase"

type CommonResult[T any, E any] struct {
	Success bool             `json:"success"`
	Result  T                `json:"result"`
	Error   E                `json:"error"`
	MinFee  jbase.JsonMinFee `json:"minFee"`
}
