package login

import (
	"github.com/pumahawk/simplcli/lib/svc"
)

type ConfigFlags struct {
	User       string
	Server     LocalServer
	AuthServer svc.AuthServer
}

type LocalServer struct {
	Port string
}
