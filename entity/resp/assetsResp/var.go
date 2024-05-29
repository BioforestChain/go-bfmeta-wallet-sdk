package assetsResp

// GetAssetsData 表示单个资产数据的结构体
type GetAssetsData struct {
	AssetType        string `json:"assetType"`
	SourceChainMagic string `json:"sourceChainMagic"`
	ApplyAddress     string `json:"applyAddress"`
	SourceChainName  string `json:"sourceChainName"`
	IconURL          string `json:"iconUrl"`
}

// GetAssetsResp 表示获取资产响应的结构体
type GetAssetsResp struct {
	Page     int             `json:"page"`
	PageSize int             `json:"pageSzie"`
	Total    int             `json:"total"`
	HasMore  bool            `json:"hasMore"`
	DataList []GetAssetsData `json:"dataList"`
}

type GetAssetsRespResult struct {
	Success bool          `json:"success"`
	Result  GetAssetsResp `json:"result"`
}
