package transactionsResp

// GetTransactionsByBrowserResp 表示通过浏览器获取交易的返回数据结构体
type GetTransactionsByBrowserResp struct {
	Page     int                      `json:"page"`
	PageSize int                      `json:"pageSize"`
	Total    int                      `json:"total"`
	HasMore  bool                     `json:"hasMore"`
	DataList []TransactionInBlockJSON `json:"dataList"`
}

// TransactionAssetChangeJSON 表示交易资产变动数据结构体
type TransactionAssetChangeJSON struct {
	AccountType      int    `json:"accountType"`
	SourceChainMagic string `json:"sourceChainMagic"`
	AssetType        string `json:"assetType"`
	AssetPrealnum    string `json:"assetPrealnum"`
}

// AssetPrealnumJSON 如果有定义的话，可以在此补充对应的结构体定义
type AssetPrealnumJSON struct {
	RemainAssetPrealnum     string `json:"remainAssetPrealnum"`
	FrozenMainAssetPrealnum string `json:"frozenMainAssetPrealnum"`
}

// TransactionInBlockJSON 表示区块中的交易数据结构体
type TransactionInBlockJSON struct {
	TIndex                  int                          `json:"tIndex"`
	Height                  int                          `json:"height"`
	TransactionAssetChanges []TransactionAssetChangeJSON `json:"transactionAssetChanges"`
	AssetPrealnum           *AssetPrealnumJSON           `json:"assetPrealnum,omitempty"`
	Signature               string                       `json:"signature"`
	SignSignature           *string                      `json:"signSignature,omitempty"`
}

//baseApis2

// GetTransactionsResult 表示获取交易结果的结构体
type GetTransactionsResp struct {
	Trs              []TransactionInBlockJSON `json:"trs"`
	Count            int                      `json:"count"`
	CmdLimitPerQuery int                      `json:"cmdLimitPerQuery"`
}

type GetTransactionsResult struct {
	Success bool                `json:"success"`
	Result  GetTransactionsResp `json:"result"`
}

type GetTransactionsByBrowserResult struct {
	Success bool                         `json:"success"`
	Result  GetTransactionsByBrowserResp `json:"result"`
}
