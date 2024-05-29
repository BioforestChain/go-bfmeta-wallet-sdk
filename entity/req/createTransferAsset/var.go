package createTransferAsset

// KVStorageInfo 结构定义
type KVStorageInfo struct {
	Key      string   `json:"key"`
	FileInfo FileInfo `json:"fileInfo"`
}

// FileInfo 结构定义
type FileInfo struct {
	Name string `json:"name"`
	Size int    `json:"size"`
}

// TransactionCommonParams 定义
type TransactionCommonParams struct {
	PublicKey               string            `json:"publicKey"`
	SecondPublicKey         string            `json:"secondPublicKey,omitempty"`
	RecipientId             string            `json:"recipientId,omitempty"`
	RangeType               int               `json:"rangeType,omitempty"`
	Range                   []string          `json:"range,omitempty"`
	Fee                     string            `json:"fee"`
	ApplyBlockHeight        int               `json:"applyBlockHeight"`
	Remark                  map[string]string `json:"remark,omitempty"`
	Dappid                  string            `json:"dappid,omitempty"`
	Lns                     string            `json:"lns,omitempty"`
	SourceIP                string            `json:"sourceIP,omitempty"`
	FromMagic               string            `json:"fromMagic,omitempty"`
	ToMagic                 string            `json:"toMagic,omitempty"`
	NumberOfEffectiveBlocks int               `json:"numberOfEffectiveBlocks,omitempty"`
	BinaryInfos             []KVStorageInfo   `json:"binaryInfos,omitempty"`
	Timestamp               int64             `json:"timestamp,omitempty"`
}

// TransactionCommonParamsWithRecipientId 定义
type TransactionCommonParamsWithRecipientId struct {
	TransactionCommonParams
	RecipientId string `json:"recipientId"`
}

// TransferAssetTransactionParams 定义
type TransferAssetTransactionParams struct {
	TransactionCommonParamsWithRecipientId
	SourceChainMagic string `json:"sourceChainMagic,omitempty"`
	SourceChainName  string `json:"sourceChainName,omitempty"`
	AssetType        string `json:"assetType,omitempty"`
	Amount           string `json:"amount"`
}
