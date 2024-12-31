package login

import (
	"github.com/pumahawk/simpl-cli/lib/svc/auth"
)

type ConfigFlags struct {
	User       string
	Server     LocalServer
	AuthServer auth.AuthServer
}

type LocalServer struct {
	Port string
}
