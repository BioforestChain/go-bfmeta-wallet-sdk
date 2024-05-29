package transactions

type GetAddressBalanceReq struct {
	Address   string `json:"address,omitempty"`
	Magic     string `json:"magic,omitempty"`
	AssetType string `json:"asset_type,omitempty"`
}

// GetTransactionsParams 表示获取交易参数的结构体
type GetTransactionsParams struct {
	Signature    string   `json:"signature,omitempty"`
	Height       int      `json:"height,omitempty"`
	MinHeight    int      `json:"minHeight,omitempty"`
	MaxHeight    int      `json:"maxHeight,omitempty"`
	SenderId     string   `json:"senderId,omitempty"`
	RecipientId  string   `json:"recipientId,omitempty"`
	Address      string   `json:"address,omitempty"`
	Type         []string `json:"type,omitempty"`
	StorageValue string   `json:"storageValue,omitempty"`
	Page         int      `json:"page,omitempty"`
	PageSize     int      `json:"pageSize,omitempty"`
	Sort         int      `json:"sort,omitempty"`
}
