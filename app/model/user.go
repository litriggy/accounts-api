package model

// Soft integration request body
//	@Description	서명으로 지갑 등록 request Body struct
type AddSoftWallet struct {
	WalletAddr string `json:"walletAddr" example:"ex) 0x00..."`
	Signature  string `json:"signature" example:"ex) 0x00..."`
	Salt       string `json:"salt" example:"ex) 165...."`
	WalletType string `json:"walletType" example:"ex) kaikas, metamask"`
}

// Hard integration request body
//	@Description	개인 키로 지갑 등록 request Body struct
type AddHardWallet struct {
	PrivateKey string `json:"privateKey" example:"ex) 027737b5...."`
	WalletType string `json:"walletType" example:"ex) eth, sol, bit, apt"`
	SecPw      string `json:"secPw" example:"ex) 000000"`
}

type AddSecondPass struct {
	SecPw string `json:"secPw"`
}

//resp

type UserInfo struct {
	ID       int32  `json:"userId"`
	Nickname string `json:"nickname"`
	Email    string `json:"email"`
	Picture  string `json:"picture"`
}
