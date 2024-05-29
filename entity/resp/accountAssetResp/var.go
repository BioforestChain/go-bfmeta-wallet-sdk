package accountAssetResp

// AssetDetail 表示资产详细信息的结构体
type AssetDetail struct {
	SourceChainMagic string `json:"sourceChainMagic"`
	AssetType        string `json:"assetType"`
	SourceChainName  string `json:"sourceChainName"`
	AssetNumber      string `json:"assetNumber"`
}

// AssetsMap 表示嵌套的资产映射结构体
type AssetsMap map[string]map[string]AssetDetail

// GetAccountAssetResp 表示获取账户资产信息的返回数据结构体
type GetAccountAssetResp struct {
	Address        string    `json:"address"`
	Assets         AssetsMap `json:"assets"`
	ForgingRewards string    `json:"forgingRewards"`
	VotingRewards  string    `json:"votingRewards"`
}

type GetAccountAssetRespResult struct {
	Success bool                `json:"success"`
	Result  GetAccountAssetResp `json:"result"`
}

// all
// GetAllAccountAssetResp 表示获取所有账户资产响应的结构体
type GetAllAccountAssetResp []map[string]map[string]string

type GetAllAccountAssetRespResult struct {
	Success bool                   `json:"success"`
	Result  GetAllAccountAssetResp `json:"result"`
}
