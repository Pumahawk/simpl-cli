package profile

import (
	"encoding/json"
	"flag"
	"log"
	"os"

	app "github.com/pumahawk/simplcli/lib/application"
	pr "github.com/pumahawk/simplcli/lib/svc/profile"
)

func Exec(conf app.Data, args []string) {
	flags := ReadConfigFlag(conf, args)
	var prFileName = pr.GetProfileFile(conf.DirData, flags.Name)
	if flags.Set {
		stpr, err := pr.LoadProfile(prFileName)
		if err != nil {
			stpr = pr.Info{}
		}
		if flags.User != "" {
			stpr.User = flags.User
		}
		if flags.Khost != "" {
			stpr.KeyCloakHost = flags.Khost
		}
		if flags.Krealm != "" {
			stpr.KeyCloakRealm = flags.Krealm
		}
		stpr.SaveProfile(prFileName)
	} else {
		stpr, err := pr.LoadProfile(prFileName)
		if err != nil {
			log.Fatalf("Unable to get profile info. %s", err.Error())
		}
		err = json.NewEncoder(os.Stdout).Encode(stpr)
		if err != nil {
			log.Fatalf("Unable to necode profile info")
		}
	}
}

func ReadConfigFlag(appData app.Data, args []string) (config ConfigFlags) {
	flags := flag.NewFlagSet("profile", flag.ExitOnError)
	flags.StringVar(&config.Name, "name", "", "Profile name")
	flags.BoolVar(&config.Set, "set", false, "Create or override profile")
	flags.StringVar(&config.User, "user", "", "Set user")
	flags.StringVar(&config.Khost, "keycloak-host", "", "Set Keycloak host")
	flags.StringVar(&config.Krealm, "keycloak-realm", "", "Set Keycloak realm")
	flags.Parse(args)
	if config.Name == "" {
		log.Fatalln("Invalid profile name")
	}
	return
}
