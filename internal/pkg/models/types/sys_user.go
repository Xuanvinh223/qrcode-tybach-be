package types

type TokenInfo struct {
	AccessToken string `json:"accessToken"`
}

type UserInfo struct {
	Permissions []string `json:"permissions"`
	Username    string   `json:"username"`
	RealName    string   `json:"realname"`
	Avatar      string   `json:"avatar"`
	Email       string   `json:"email"`
}
