package createAccountResp

type CreateAccountResp struct {
	Address   string `json:"address"`
	PublicKey string `json:"publicKey"`
	SecretKey string `json:"secretKey"`
}

type CreateAccountRespResult struct {
	Success bool              `json:"success"`
	Result  CreateAccountResp `json:"result"`
}
