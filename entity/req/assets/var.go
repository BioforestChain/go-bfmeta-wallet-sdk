package assets

// PaginationOptions 表示分页选项的结构体
type PaginationOptions struct {
	Page      int    `json:"page"`
	PageSize  int    `json:"pageSize"`
	AssetType string `json:"assetType,omitempty"`
}
