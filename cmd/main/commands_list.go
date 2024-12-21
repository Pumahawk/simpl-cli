package main

import (
	"github.com/pumahawk/simplcli/lib/command/login"
	"github.com/pumahawk/simplcli/lib/command/token"
)

var DefinedCommands = []SubCommand{
	{"login", login.Exec},
	{"token", token.Exec},
}
