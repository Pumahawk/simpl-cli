package main

import (
	"flag"
	"log"
	"os"

	app "github.com/pumahawk/simplcli/lib/application"
	"github.com/pumahawk/simplcli/lib/svc/profile"
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
	flag.StringVar(&appData.DirData, "user", os.TempDir(), "Configuration directory")
	p := flag.String("profile", "", "Profile name")
	flag.StringVar(&appData.KCHost, "keycloak-host", "", "Keycloak host")
	flag.StringVar(&appData.KCRealm, "keycloak-realm", "authority", "Keycloak realm")
	flag.Parse()

	if *p != "" {
		profile, err := profile.LoadProfile(profile.GetProfileFile(appData.DirData, *p))
		if err != nil {
			log.Fatalf("Unable to load profile %s. %s", *p, err.Error())
		}
		mapProfileToAppData(&appData, profile)
	}
	return appData, flag.Args()
}

func mapProfileToAppData(appData *app.Data, profile profile.Info) {
	if appData.User == "" {
		appData.User = profile.User
	}
	if appData.KCHost == "" {
		appData.KCHost = profile.KeyCloakHost
	}
	if appData.KCRealm == "" {
		appData.KCRealm = profile.KeyCloakRealm
	}
}
