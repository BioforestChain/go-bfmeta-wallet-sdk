package asymmetricDecryptResp

type ResAsymmetricDecrypt struct {
	EncryptedMessage string `json:"encryptedMessage"`
	Nonce            string `json:"nonce"`
}
