package main

import (
	"flag"
	app "github.com/pumahawk/simplcli/lib/application"
	"log"
	"os"
)

func main() {
	data, args := readArgs()
	if len(args) < 1 {
		log.Fatal("Missing subcommand")
	}
	subcommand := args[0]
	for _, command := range DefinedCommands {
		if command.Name == subcommand {
			command.Exec(data, args[1:])
			os.Exit(0)
		}
	}
	log.Fatalf("Subcommand mandatory. Parameter: %s", subcommand)
}

func readArgs() (app.Data, []string) {
	appData := app.Data{}
	flag.StringVar(&appData.DirData, "dir-data", os.TempDir(), "Configuration directory")
	flag.Parse()
	return appData, flag.Args()
}
