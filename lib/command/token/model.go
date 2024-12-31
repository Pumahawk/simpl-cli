package token

import "github.com/pumahawk/simpl-cli/lib/svc/auth"

type ConfigFlags struct {
	User       string
	Verbose    bool
	AuthServer auth.AuthServer
}
