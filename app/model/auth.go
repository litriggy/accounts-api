package model

// SignIn model info
//	@Description	로그인/회원가입 request Body struct
type SignIn struct {
	AccessToken string `json:"accessToken"`
}

// GoogleUserInfo model info
//	@Description	chrome extension auth info struct
type GoogleUserInfo struct {
	ID            string `json:"id"`
	Email         string `json:"email"`
	VerifiedEmail bool   `json:"verified_email"`
	Name          string `json:"name"`
	GivenName     string `json:"given_name"`
	FamilyName    string `json:"family_name"`
	Link          string `json:"link"`
	Picture       string `json:"picture"`
	Locale        string `json:"locale"`
	Hd            string `json:"hd"`
}

type AuthInfo struct {
	ID       string
	Email    string
	Picture  string
	Nickname string
}
