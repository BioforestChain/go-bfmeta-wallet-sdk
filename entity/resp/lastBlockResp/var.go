package lastBlockResp

// LastBlockInfo 包含最新区块信息的结构体
type LastBlockInfo[T any] struct {
	Version                  int                      `json:"version"`
	Height                   int                      `json:"height"`
	Timestamp                int64                    `json:"timestamp"`
	BlockSize                int                      `json:"blockSize"`
	GeneratorPublicKey       string                   `json:"generatorPublicKey"`
	GeneratorSecondPublicKey string                   `json:"generatorSecondPublicKey,omitempty"`
	GeneratorEquity          string                   `json:"generatorEquity"`
	PreviousBlockSignature   string                   `json:"previousBlockSignature"`
	Reward                   string                   `json:"reward"`
	Magic                    string                   `json:"magic"`
	BlockParticipation       string                   `json:"blockParticipation"`
	Signature                string                   `json:"signature"`
	SignSignature            string                   `json:"signSignature,omitempty"`
	Remark                   map[string]string        `json:"remark"`
	TransactionInfo          BlockTransactionInfoJSON `json:"transactionInfo"`
	Asset                    T                        `json:"asset"`
}

type GetLastBlockInfoRespResult struct {
	Success bool               `json:"success"`
	Result  LastBlockInfo[any] `json:"result"`
}

// GetLastBlockResult 表示获取最后一个区块结果的结构体
type GetLastBlockResult LastBlockInfo[any]

// BlockTransactionInfoJSON 包含区块交易信息的结构体
type BlockTransactionInfoJSON struct {
	StartTindex          int    `json:"startTindex"`
	NumberOfTransactions int    `json:"numberOfTransactions"`
	PayloadHash          string `json:"payloadHash"`
	PayloadLength        int    `json:"payloadLength"`
	TotalAmount          string `json:"totalAmount"`
	TotalFee             string `json:"totalFee"`
	//TransactionInBlocks  []TransactionInBlockJSON `json:"transactionInBlocks"`
	//todo Use of generic type without instantiation
	TransactionInBlocks TransactionInBlockJSON `json:"transactionInBlocks"`

	StatisticInfo StatisticInfoJSON `json:"statisticInfo"`
}

// TransactionInBlockJSON 表示区块中交易信息的结构体
type TransactionInBlockJSON []struct {
	TIndex                  int                          `json:"tIndex"`
	Height                  int                          `json:"height"`
	TransactionAssetChanges []TransactionAssetChangeJSON `json:"transactionAssetChanges"`
	AssetPrealnum           *AssetPrealnumJSON           `json:"assetPrealnum,omitempty"`
	Signature               string                       `json:"signature"`
	SignSignature           *string                      `json:"signSignature,omitempty"`
	Transaction             interface{}                  `json:"transaction"`
}

// SomeTransactionJSON 表示某些交易信息的结构体
type SomeTransactionJSON[T any] struct {
	Transaction T `json:"transaction"`
}

// TransactionJSON 基础交易信息的结构体
type TransactionJSON[AssetJSON any] struct {
	Version               int                     `json:"version"`
	Type                  string                  `json:"type"`
	SenderId              string                  `json:"senderId"`
	SenderPublicKey       string                  `json:"senderPublicKey"`
	SenderSecondPublicKey *string                 `json:"senderSecondPublicKey,omitempty"`
	RecipientId           *string                 `json:"recipientId,omitempty"`
	RangeType             int                     `json:"rangeType"`
	Range                 []string                `json:"range"`
	Fee                   string                  `json:"fee"`
	Timestamp             int64                   `json:"timestamp"`
	Dappid                *string                 `json:"dappid,omitempty"`
	Lns                   *string                 `json:"lns,omitempty"`
	SourceIP              *string                 `json:"sourceIP,omitempty"`
	FromMagic             string                  `json:"fromMagic"`
	ToMagic               string                  `json:"toMagic"`
	ApplyBlockHeight      int                     `json:"applyBlockHeight"`
	EffectiveBlockHeight  int                     `json:"effectiveBlockHeight"`
	Signature             string                  `json:"signature"`
	SignSignature         *string                 `json:"signSignature,omitempty"`
	Remark                map[string]string       `json:"remark"`
	Asset                 AssetJSON               `json:"asset"`
	Storage               *TransactionStorageJSON `json:"storage,omitempty"`
	StorageKey            *string                 `json:"storageKey,omitempty"`
	StorageValue          *string                 `json:"storageValue,omitempty"`
	Nonce                 int                     `json:"nonce"`
}

// TransactionStorageJSON 表示交易存储信息的结构体
type TransactionStorageJSON struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

// StatisticInfoJSON 包含统计信息的结构体
type StatisticInfoJSON struct {
	TotalFee                           string                                 `json:"totalFee"`
	TotalAsset                         string                                 `json:"totalAsset"`
	TotalChainAsset                    string                                 `json:"totalChainAsset"`
	TotalAccount                       int                                    `json:"totalAccount"`
	MagicAssetTypeTypeStatisticHashMap map[string]AssetTypeAssetStatisticJSON `json:"magicAssetTypeTypeStatisticHashMap"`
}

// AssetTypeAssetStatisticJSON 表示资产类型统计信息的结构体
type AssetTypeAssetStatisticJSON struct {
	AssetTypeTypeStatisticHashMap map[string]AssetStatisticJSON `json:"assetTypeTypeStatisticHashMap"`
}

// AssetStatisticJSON 包含资产统计信息的结构体
type AssetStatisticJSON struct {
	TypeStatisticHashMap map[string]CountAndAmountStatisticJSON `json:"typeStatisticHashMap"`
	Total                CountAndAmountStatisticJSON            `json:"total"`
}

// CountAndAmountStatisticJSON 表示交易计数和数量统计信息的结构体
type CountAndAmountStatisticJSON struct {
	ChangeAmount     string `json:"changeAmount"`
	ChangeCount      int    `json:"changeCount"`
	MoveAmount       string `json:"moveAmount"`
	TransactionCount int    `json:"transactionCount"`
}

// TransactionAssetChangeJSON 表示交易资产变更的结构体
type TransactionAssetChangeJSON struct {
	AccountType      int    `json:"accountType"`
	SourceChainMagic string `json:"sourceChainMagic"`
	AssetType        string `json:"assetType"`
	AssetPrealnum    string `json:"assetPrealnum"`
}

// AssetPrealnumJSON 表示资产的数量结构体
type AssetPrealnumJSON struct {
	RemainAssetPrealnum     string `json:"remainAssetPrealnum"`
	FrozenMainAssetPrealnum string `json:"frozenMainAssetPrealnum"`
}
