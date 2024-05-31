package asymmetricEncryptResp

type ResAsymmetricEncrypt struct {
	EncryptedMessage string `json:"encryptedMessage"`
	Nonce            string `json:"nonce"`
}
