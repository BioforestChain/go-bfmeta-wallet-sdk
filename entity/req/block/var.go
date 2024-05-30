package block

// 基础 API 请求参数的结构体

// 获取区块参数的结构体
type GetBlockParams struct {
	Signature string `json:"signature,omitempty"`
	Height    int    `json:"height,omitempty"`
	Page      int    `json:"page,omitempty"`
	PageSize  int    `json:"pageSize,omitempty"`
}
