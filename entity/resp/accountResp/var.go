package accountResp

// EquityInfo 表示股权信息的结构体
type EquityInfo struct {
	Round       int    `json:"round"`
	Equity      string `json:"equity"`
	FixedEquity string `json:"fixedEquity"`
}

// GetAccountInfoResp 表示获取账户信息的返回数据结构体
type GetAccountInfoResp struct {
	Address         string     `json:"address"`
	PublicKey       string     `json:"publicKey"`
	SecondPublicKey string     `json:"secondPublicKey"`
	IsDelegate      bool       `json:"isDelegate"`
	IsAcceptVote    bool       `json:"isAcceptVote"`
	AccountStatus   int        `json:"accountStatus"`
	EquityInfo      EquityInfo `json:"equityInfo"`
}

type GetAccountInfoRespResult struct {
	Success bool               `json:"success"`
	Result  GetAccountInfoResp `json:"result"`
	//Result interface{} `json:"result"`
}
