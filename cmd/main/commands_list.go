package main

import (
	"github.com/pumahawk/simplcli/lib/command/login"
)

var DefinedCommands = []SubCommand{
	{"login", login.Exec},
}
