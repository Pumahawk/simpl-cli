package auth

import (
	"net/url"
	"time"
)

type AuthInfo struct {
	Path                string
	State               string
	ResponseMode        string
	ResponseType        string
	Scope               string
	Nonce               string
	CodeChallenge       string
	CodeChallengeMethod string
}

type UserToken struct {
	Code         string
	Iss          string
	SessionState string
	State        string
}

type AuthServer struct {
	Host     string
	ClientId string
	Realm    string
}

type TokenizeInfo struct {
	Code         string
	GrantType    string
	ClientId     string
	RedirectUri  string
	CodeVerifier string
}

type TokenInfo struct {
	TimeExiration    int64
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

func NewAuthInfo(host string, realm string) AuthInfo {
	return AuthInfo{
		Path:                "/realms/" + realm + "/protocol/openid-connect/auth",
		State:               "29f2c56a-d4df-49cf-87dc-a870669a61ab",
		ResponseMode:        "fragment",
		ResponseType:        "code",
		Scope:               "openid",
		Nonce:               "bb0c8a3e-9667-4ec8-8dde-10c0a52f40be",
		CodeChallenge:       "dyqvDKwOIdLE50mxt6o7_jDj8IkNAdcKi554hvCGFEQ",
		CodeChallengeMethod: "S256",
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

func NewTokenizeInfo(code string, localPort string) TokenizeInfo {
	return TokenizeInfo{
		Code:         code,
		GrantType:    "authorization_code",
		ClientId:     "frontend-cli",
		RedirectUri:  "http://localhost:" + localPort + "/auth",
		CodeVerifier: "gd8PkFgqwnYZOJJrxuMDk0Rjk2q3hx6VYYpIas4KvsECpPBpMXttrxc8bsT9kPtM8w41IdkvvBJOfX4RqwJLSM1hgrgBv5t6",
	}
}

func (info *AuthInfo) ToURI(authServer AuthServer, localPort string) string {
	return authServer.Host +
		info.Path +
		"?client_id=" + authServer.ClientId +
		"&redirect_uri=http%3A%2F%2Flocalhost%3A" + string(localPort) + "%2Fauth" +
		"&state=" + info.State +
		"&response_mode=" + info.ResponseMode +
		"&response_type=" + info.ResponseType +
		"&scope=" + info.Scope +
		"&nonce=" + info.Nonce +
		"&code_challenge=" + info.CodeChallenge +
		"&code_challenge_method=" + info.CodeChallengeMethod
}

func (tokenInfo *TokenInfo) UpdateExpirationTime(t time.Time) {
	tokenInfo.TimeExiration = t.Add(time.Second * time.Duration(tokenInfo.ExpiresIn)).Unix()
}

func (tokenInfo *TokenInfo) IsExpired() bool {
	return time.Now().Unix() > tokenInfo.TimeExiration
}
