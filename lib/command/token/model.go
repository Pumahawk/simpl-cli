package token

import "net/url"

type ConfigFlags struct {
	User       string
	Port       string
	AuthServer AuthServer
}

type AuthServer struct {
	Host     string
	ClientId string
	Realm    string
}

type TokenInfo struct {
	AccessToken      string `json:"access_token"`
	ExpiresIn        int    `json:"expires_in"`
	RefreshExpiresIn int    `json:"refresh_expires_in"`
	RefreshToken     string `json:"refresh_token"`
	TokenType        string `json:"token_type"`
	IdToken          string `json:"id_token"`
	NotBeforePolicy  int    `json:"not-before-policy"`
	SessionState     string `json:"session_state"`
	Scope            string `json:"scope"`
}

type TokenizeInfo struct {
	Code         string
	GrantType    string
	ClientId     string
	RedirectUri  string
	CodeVerifier string
}

func NewTokenizeInfo(code string, localPort string) TokenizeInfo {
	return TokenizeInfo{
		Code:         code,
		GrantType:    "authorization_code",
		ClientId:     "frontend-cli",
		RedirectUri:  "http://localhost:" + localPort + "/auth",
		CodeVerifier: "gd8PkFgqwnYZOJJrxuMDk0Rjk2q3hx6VYYpIas4KvsECpPBpMXttrxc8bsT9kPtM8w41IdkvvBJOfX4RqwJLSM1hgrgBv5t6",
	}
}

func (token TokenizeInfo) ToUrlValues(values *url.Values, localPort string) {
	values.Add("code", token.Code)
	values.Add("grant_type", "authorization_code")
	values.Add("client_id", "frontend-cli")
	values.Add("redirect_uri", "http://localhost:"+localPort+"/auth")
	values.Add("code_verifier", "gd8PkFgqwnYZOJJrxuMDk0Rjk2q3hx6VYYpIas4KvsECpPBpMXttrxc8bsT9kPtM8w41IdkvvBJOfX4RqwJLSM1hgrgBv5t6")
	return
}
