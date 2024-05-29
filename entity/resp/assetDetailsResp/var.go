package assetDetailsResp

// Asset 表示基础资产信息的结构体
type Asset struct {
	AssetType               string `json:"assetType"`               // 资产名称
	ApplyAddress            string `json:"applyAddress"`            // 发行地址
	GenesisAddress          string `json:"genesisAddress"`          // 创世地址
	SourceChainName         string `json:"sourceChainName"`         // 归属链名称
	IssuedAssetPrealnum     string `json:"issuedAssetPrealnum"`     // 发行总量
	RemainAssetPrealnum     string `json:"remainAssetPrealnum"`     // 当前总量
	FrozenMainAssetPrealnum string `json:"frozenMainAssetPrealnum"` // 冻结主资产量
	PublishTime             int64  `json:"publishTime"`             // 发行时间
	SourceChainMagic        string `json:"sourceChainMagic"`        // 跨链字段
}

// AssetInfo 表示包含地址数的资产信息结构体
type AssetInfo struct {
	Asset
	AddressQty int `json:"addressQty"` // 活跃地址数
}

// GetAssetDetailsResp 表示获取资产详细信息的响应结构体
type GetAssetDetailsResp struct {
	AssetInfo
	IconURL string `json:"iconUrl"` // 图标 URL
}

type GetAssetDetailsRespResult struct {
	Success bool                `json:"success"`
	Result  GetAssetDetailsResp `json:"result"`
}
