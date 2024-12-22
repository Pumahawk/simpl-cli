package token

import "github.com/pumahawk/simplcli/lib/svc/auth"

type ConfigFlags struct {
	User       string
	Verbose    bool
	AuthServer auth.AuthServer
}
