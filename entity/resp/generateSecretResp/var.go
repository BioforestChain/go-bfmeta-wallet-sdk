package generateSecretResp

type GenerateSecretResp struct {
	Secret string `json:"secret"`
}

type GenerateSecretRespResult struct {
	Success bool               `json:"success"`
	Result  GenerateSecretResp `json:"result"`
}
