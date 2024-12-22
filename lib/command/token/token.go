package token

import (
	"encoding/json"
	"flag"
	"log"
	"os"
	"time"

	app "github.com/pumahawk/simplcli/lib/application"
	"github.com/pumahawk/simplcli/lib/svc"
)

func Exec(conf app.Data, args []string) {
	flags := ReadConfigFlag(args)
	tokenInfo, err := svc.LoadUserAuthData(conf, flags.User)
	if err != nil {
		log.Fatalf("Unable to load user auth data. %s", err.Error())
	}
	if time.Now().UnixMilli() > tokenInfo.TimeExiration {
		tokenInfo, err = svc.ReloadToken(flags.AuthServer, tokenInfo)
		if err != nil {
			log.Fatalf("Unable to reload token. %s", err.Error())
		}
		tokenInfo.UpdateExpirationTime(time.Now())
		svc.SaveUserAuthData(conf, flags.User, tokenInfo)
	}

	err = json.NewEncoder(os.Stdout).Encode(tokenInfo)
	if err != nil {
		log.Fatalf("Unable to print token")
	}
}

func ReadConfigFlag(args []string) (config ConfigFlags) {
	flags := flag.NewFlagSet("token", flag.ExitOnError)
	flags.StringVar(&config.AuthServer.Host, "host", "", "Authentication server host")
	flags.StringVar(&config.AuthServer.Realm, "realm", "authority", "Realm")
	flags.StringVar(&config.AuthServer.ClientId, "client-id", "frontend-cli", "Client Id")
	flags.StringVar(&config.Port, "port", "8080", "Redirect local server port")
	flags.StringVar(&config.User, "user", "default", "User session")
	flags.Parse(args)

	if config.AuthServer.Host == "" {
		log.Fatalln("Mandatory --host flag missing")
	}
	return
}
