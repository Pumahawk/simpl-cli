package main

import "github.com/pumahawk/simpl-cli/lib/application"

type CommandExec = func(application.Data, []string)

type SubCommand struct {
	Name string
	Exec CommandExec
}
