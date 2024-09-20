package broadcastResp

import "github.com/BioforestChain/go-bfmeta-wallet-sdk/entity/jbase"

type BroadcastRespResult[T any] struct {
	Success bool                             `json:"success"`
	Result  TransactionApiSuccessReturn[any] `json:"result"`
	MinFee  jbase.JsonMinFee                 `json:"minFee"`
	Error   jbase.JsonCommonError            `json:"error"`
}

// TransactionApiSuccessReturn 继承 ApiSuccessReturn 结构体，并增加 minFee 字段
type TransactionApiSuccessReturn[T any] struct {
	MinFee jbase.JsonMinFee `json:"minFee"`
}
