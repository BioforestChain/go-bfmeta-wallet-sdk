package broadcastTra

type BroadcastTransactionParams struct {
	Signature     string `json:"signature"`
	SignSignature string `json:"signSignature,omitempty"`
	Buffer        string `json:"buffer"`
	IsOnChain     bool   `json:"isOnChain,omitempty"`
}
