package accountAsset

type GetAccountAssetParams struct {
	Address string `json:"address"`
}

// GetAllAccountAssetReq 表示获取所有账户资产请求的结构体
type GetAllAccountAssetReq struct {
	Filter map[string]string `json:"filter"`
}
