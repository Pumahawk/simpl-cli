package main

import (
	"github.com/pumahawk/simplcli/lib/command/login"
	"log"
	"os"
)

var defaultTokenFilePath = os.TempDir() + "/tokenfile"
var tokenFile string

func main() {
	if len(os.Args) < 2 {
		log.Fatal("Missing subcommand")
	}
	subcommand := os.Args[1]
	switch subcommand {
	case "login":
		login.Exec(os.Args[2:])
	default:
		log.Fatalf("Subcommand mandatory. Parameter: %s", subcommand)
	}
}

func LoopUpEnvDefault(envName string, defValue string) string {
	val, exist := os.LookupEnv("SIMPLCLI_TOKENFILE")
	if exist {
		return val
	} else {
		return defValue
	}
}
