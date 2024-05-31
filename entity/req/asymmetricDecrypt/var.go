package asymmetricDecrypt

type Req struct {
	EncryptedMessage []byte `json:"encryptedMessage"`
	EncryptPK        []byte `json:"encryptPk"`
	DecryptSK        []byte `json:"decryptSk"`
	Nonce            []byte `json:"nonce,omitempty"`
}
