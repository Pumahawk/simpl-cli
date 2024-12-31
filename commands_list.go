package main

import (
	"github.com/pumahawk/simpl-cli/lib/command/login"
	"github.com/pumahawk/simpl-cli/lib/command/profile"
	"github.com/pumahawk/simpl-cli/lib/command/token"
)

var DefinedCommands = []SubCommand{
	{"login", login.Exec},
	{"token", token.Exec},
	{"profile", profile.Exec},
}
