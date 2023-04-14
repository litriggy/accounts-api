package model

type RECVTransfer struct {
	Target        uint32                `json:"to"`            // example:"송금 대상 유저 Id 값"`
	ServiceID     uint32                `json:"serviceId"`     // example:"서비스 Id 값"`
	OnchainEvent  []OnTransactionDetail `json:"onChainEvents"` // example:"블록체인 거래"`
	OffchainEvent OffTransactionDetail  `json:"offChainEvent"` // example:"장부거래"`
	SecPw         string                `json:"secPw"`         // example:"2차 비밀번호"`
}

type TransferOnchain struct {
	Sender    []string `json:"sender"`
	Target    string   `json:"to"`
	Amount    []int64  `json:"amount"`
	ServiceID uint32   `json:"serviceId"`
}

type OffTransactionDetail struct {
	To     string `json:"to"`
	Amount int64  `json:"amount"`
}

type OnTransactionDetail struct {
	From   string `json:"from"`   // example:"보낼 db에 기록된 지갑 주소"`
	To     string `json:"to"`     // example:"대상 지갑 주소"`
	Amount int64  `json:"amount"` // example:"수량 wei 단위로 작성"`
	Txhash string `json:"txhash"` // example:"pk 보관 중이지 않은 지갑 일 경우 필수 아님" require:"false"`
}

//RESPONSES

type TransactionsList struct {
	TransactionID int32     `json:"transactionId"`
	FromID        int32     `json:"from"`
	ToID          int32     `json:"to"`
	Memo          string    `json:"memo"`
	TotalAmount   int64     `json:"amount"`
	TokenInfo     TokenInfo `json:"tokenInfo"`
}

type TokenInfo struct {
	Name       string `json:"name"`
	Symbol     string `json:"symbol"`
	Decimals   int32  `json:"decimals"`
	Image      string `json:"img"`
	IsNative   bool   `json:"isNative"`
	NetType    string `json:"netType"`
	WalletType string `json:"walletType"`
}

type TransactionDetail struct {
	From      string `json:"from"`
	To        string `json:"to"`
	Amount    int64  `json:"amount"`
	IsOnchain bool   `json:"isOnChain"`
	Txhash    string `json:"txhash"`
	Status    int32  `json:"status"`
}
