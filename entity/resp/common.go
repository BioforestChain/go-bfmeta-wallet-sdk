package resp

type CommonResult[T any, E any] struct {
	Success bool   `json:"success"`
	Result  T      `json:"result"`
	Error   E      `json:"error"`
	MinFee  string `json:"minFee"`
}
