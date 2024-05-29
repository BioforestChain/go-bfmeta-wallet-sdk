package address

type Params struct {
	Address   string `json:"address"`
	Magic     string `json:"magic"`
	AssetType string `json:"assetType,omitempty"`
}
