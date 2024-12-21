package main

import "github.com/pumahawk/simplcli/lib/command/login"

var DefinedCommands = []SubCommand{
	LoginCommand,
}

var LoginCommand = SubCommand{
	Name: "login",
	Exec: login.Exec,
}
