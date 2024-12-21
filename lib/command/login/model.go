package login

import (
	app "github.com/pumahawk/simplcli/lib/application"
)

type LoginDataCommand struct {
	AppData app.Data
}

type ConfigFlags struct {
	Server     LocalServer
	AuthServer AuthServer
}

type AuthServer struct {
	Host     string
	ClientId string
	Realm    string
}

type LocalServer struct {
	Port string
}

type UserToken struct {
	Code         string
	Iss          string
	SessionState string
	State        string
}

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

func NewAuthInfo(host string) AuthInfo {
	return AuthInfo{
		Path:                "/realms/authority/protocol/openid-connect/auth",
		State:               "29f2c56a-d4df-49cf-87dc-a870669a61ab",
		ResponseMode:        "fragment",
		ResponseType:        "code",
		Scope:               "openid",
		Nonce:               "bb0c8a3e-9667-4ec8-8dde-10c0a52f40be",
		CodeChallenge:       "dyqvDKwOIdLE50mxt6o7_jDj8IkNAdcKi554hvCGFEQ",
		CodeChallengeMethod: "S256",
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
