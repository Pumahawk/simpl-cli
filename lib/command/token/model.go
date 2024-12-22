package token

import "github.com/pumahawk/simplcli/lib/svc"

type ConfigFlags struct {
	User       string
	Port       string
	AuthServer svc.AuthServer
}
