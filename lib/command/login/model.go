package login

type UserToken struct {
	Code string
}

type AuthInfo struct {
	Host string
	Path string
	ClientId string
	RedirectUri string
	State string
	ResponseMode string
	ResponseType string
	Scope string
	Nonce string
	CodeChallenge string
	CodeChallengeMethod string
}

func NewAuthInfo(host string) AuthInfo {
	return AuthInfo{
		Host: host,
		Path: "/realms/authority/protocol/openid-connect/auth",
		ClientId: "frontend-cli",
		RedirectUri:  "http%3A%2F%2Flocalhost%3A8080%2Fauth",
		State: "29f2c56a-d4df-49cf-87dc-a870669a61ab",
		ResponseMode: "fragment",
		ResponseType: "code",
		Scope: "openid",
		Nonce: "bb0c8a3e-9667-4ec8-8dde-10c0a52f40be",
		CodeChallenge: "h65k6qficd4zde4UqCMzo-kieGPBUF830XOEnbmzIhM",
		CodeChallengeMethod: "S256",
	}
}

func (info *AuthInfo) ToURI() string {
	return  info.Host +
	info.Path +
	"?client_id=" + info.ClientId +
	"&redirect_uri=" + info.RedirectUri +
	"&state=" + info.State +
	"&response_mode=" + info.ResponseMode +
	"&response_type=" + info.ResponseType +
	"&scope=" + info.Scope +
	"&nonce=" + info.Nonce +
	"&code_challenge=" + info.CodeChallenge +
	"&code_challenge_method=" + info.CodeChallengeMethod
}
