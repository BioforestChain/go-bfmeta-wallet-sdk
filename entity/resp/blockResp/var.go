package blockResp

type GetBlockResultResp struct {
	Success bool           `json:"success"`
	Result  GetBlockResult `json:"result"`
}

// GetBlockResult represents the structure of the basic API request result for blocks.
type GetBlockResult struct {
	Blocks           []BlockJSON `json:"blocks"`
	Count            int         `json:"count"`
	CmdLimitPerQuery int         `json:"cmdLimitPerQuery"`
}

// BlockJSON represents a block with transaction information.
type BlockJSON struct {
	BlockWithoutTransactionJSON
	TransactionInfo BlockTransactionInfoJSON `json:"transactionInfo"`
}

// BlockTransactionInfoJSON represents the transaction information in a block.
type BlockTransactionInfoJSON struct {
	StartTindex          int                    `json:"startTindex"`
	NumberOfTransactions int                    `json:"numberOfTransactions"`
	PayloadHash          string                 `json:"payloadHash"`
	PayloadLength        int                    `json:"payloadLength"`
	TotalAmount          string                 `json:"totalAmount"`
	TotalFee             string                 `json:"totalFee"`
	TransactionInBlocks  TransactionInBlockJSON `json:"transactionInBlocks"`
	StatisticInfo        StatisticInfoJSON      `json:"statisticInfo"`
}

// TransactionInBlockJSON represents a transaction inside a block.
type TransactionInBlockJSON []struct {
	TIndex                  int                          `json:"tIndex"`
	Height                  int                          `json:"height"`
	TransactionAssetChanges []TransactionAssetChangeJSON `json:"transactionAssetChanges"`
	AssetPrealnum           AssetPrealnumJSON            `json:"assetPrealnum,omitempty"`
	Signature               string                       `json:"signature"`
	SignSignature           string                       `json:"signSignature,omitempty"`
}

// SomeTransactionJSON represents some transaction JSON information.

// SomeTransactionJSON 表示某些交易信息的结构体
type SomeTransactionJSON[T any] struct {
	Transaction T `json:"transaction"`
}

// BlockWithoutTransactionJSON represents a block without transaction information.
type BlockWithoutTransactionJSON struct {
	Version                       int                           `json:"version"`
	Height                        int                           `json:"height"`
	BlockSize                     int                           `json:"blockSize"`
	Timestamp                     int                           `json:"timestamp"`
	Signature                     string                        `json:"signature"`
	SignSignature                 string                        `json:"signSignature,omitempty"`
	GeneratorPublicKey            string                        `json:"generatorPublicKey"`
	GeneratorSecondPublicKey      string                        `json:"generatorSecondPublicKey,omitempty"`
	GeneratorEquity               string                        `json:"generatorEquity"`
	PreviousBlockSignature        string                        `json:"previousBlockSignature"`
	Reward                        string                        `json:"reward"`
	Magic                         string                        `json:"magic"`
	BlockParticipation            string                        `json:"blockParticipation"`
	Remark                        map[string]string             `json:"remark"`
	Asset                         AssetJSON                     `json:"asset"`
	RoundOfflineGeneratersHashMap RoundOfflineGeneratersHashMap `json:"roundOfflineGeneratersHashMap"`
}

// RoundOfflineGeneratersHashMap represents the hash map of round offline generators.
type RoundOfflineGeneratersHashMap map[string]string

// Additional types used in the above structures (details not provided in the question) but necessary to compile:
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
} // This should be defined based on your specific transaction schema.
type TransactionAssetChangeJSON struct {
	AccountType      int    `json:"accountType"`
	SourceChainMagic string `json:"sourceChainMagic"`
	AssetType        string `json:"assetType"`
	AssetPrealnum    string `json:"assetPrealnum"`
}
type AssetPrealnumJSON struct {
	RemainAssetPrealnum     string `json:"remainAssetPrealnum"`
	FrozenMainAssetPrealnum string `json:"frozenMainAssetPrealnum"`
}
type StatisticInfoJSON struct {
	TotalFee                           string                                 `json:"totalFee"`
	TotalAsset                         string                                 `json:"totalAsset"`
	TotalChainAsset                    string                                 `json:"totalChainAsset"`
	TotalAccount                       int                                    `json:"totalAccount"`
	MagicAssetTypeTypeStatisticHashMap map[string]AssetTypeAssetStatisticJSON `json:"magicAssetTypeTypeStatisticHashMap"`
}
type AssetJSON struct{} // Define based on your schema.

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

// TransactionStorageJSON 表示交易存储信息的结构体
type TransactionStorageJSON struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}
