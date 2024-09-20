package broadcastTra

import "github.com/BioforestChain/go-bfmeta-wallet-sdk/entity/jbase"

type BroadcastTransactionParams struct {
	Signature     jbase.HexStringBuffer    `json:"signature"`
	SignSignature jbase.HexStringBuffer    `json:"signSignature,omitempty"`
	Buffer        jbase.Base64StringBuffer `json:"buffer"`
	IsOnChain     bool                     `json:"isOnChain,omitempty"`
}
