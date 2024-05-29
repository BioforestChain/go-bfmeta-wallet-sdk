package broadcastResp

type BroadcastRespResult[T any] struct {
	Success bool                             `json:"success"`
	Result  TransactionApiSuccessReturn[any] `json:"result"`
	MinFee  string                           `json:"minFee"`
	Error   Error                            `json:"error"`
}

// TransactionApiSuccessReturn 继承 ApiSuccessReturn 结构体，并增加 minFee 字段
type TransactionApiSuccessReturn[T any] struct {
	MinFee string `json:"minFee"`
}
type Error struct {
	Code        string `json:"code,omitempty"`        //  可选字段
	Message     string `json:"message"`               // 必填字段
	Description string `json:"description,omitempty"` //  可选字段
}
