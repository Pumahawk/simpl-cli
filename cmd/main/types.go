package main

import "github.com/pumahawk/simplcli/lib/application"

type CommandExec = func(application.Data, []string)

type SubCommand struct {
	Name string
	Exec CommandExec
}
