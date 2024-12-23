package token

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	app "github.com/pumahawk/simplcli/lib/application"
	"github.com/pumahawk/simplcli/lib/svc/auth"
)

func Exec(conf app.Data, args []string) {
	flags := ReadConfigFlag(conf, args)
	tokenInfo, err := auth.LoadUserAuthData(conf, flags.User)
	if err != nil {
		log.Fatalf("Unable to load user auth data. %s", err.Error())
	}
	if tokenInfo.IsExpired() {
		tokenInfo, err = auth.ReloadToken(flags.AuthServer, tokenInfo)
		if err != nil {
			log.Fatalf("Unable to reload token. %s", err.Error())
		}
		tokenInfo.UpdateExpirationTime(time.Now())
		auth.SaveUserAuthData(conf, flags.User, tokenInfo)
	}

	if flags.Verbose {
		err = json.NewEncoder(os.Stdout).Encode(tokenInfo)
		if err != nil {
			log.Fatalf("Unable to print token")
		}
	} else {
		fmt.Println(tokenInfo.AccessToken)
	}
}

func ReadConfigFlag(appData app.Data, args []string) (config ConfigFlags) {
	flags := flag.NewFlagSet("token", flag.ExitOnError)
	flags.StringVar(&config.AuthServer.Host, "host", appData.KCHost, "Authentication server host")
	flags.StringVar(&config.AuthServer.Realm, "realm", appData.KCRealm, " Keycloak Realm")
	flags.StringVar(&config.AuthServer.ClientId, "client-id", "frontend-cli", "Client Id")
	flags.BoolVar(&config.Verbose, "v", false, "Verbose mode")
	flags.StringVar(&config.User, "user", appData.User, "User session")
	flags.Parse(args)

	if config.AuthServer.Host == "" {
		log.Fatalln("Mandatory --host flag missing")
	}
	return
}
