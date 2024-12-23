package login

import (
	"flag"
	"log"
	"net/http"

	app "github.com/pumahawk/simplcli/lib/application"
	"github.com/pumahawk/simplcli/lib/svc/auth"
)

func Exec(conf app.Data, args []string) {
	config := ReadConfigFlag(conf, args)
	userTokenC := StartLoginWebServer(config.AuthServer, config.Server.Port)
	userToken, ok := <-userTokenC
	if !ok {
		log.Fatal("Unable to read from token. Channel is closed.")
	}
	tokenInfo, err := auth.Tokenize(config.AuthServer, config.Server.Port, userToken)
	if err != nil {
		log.Fatalf("Unable to tokenize. %s", err.Error())
	}
	err = auth.SaveUserAuthData(conf, config.User, tokenInfo)
	if err != nil {
		log.Fatalf("Unable to encode token to stdout. %s", err.Error())
	}
}

func ReadConfigFlag(appData app.Data, args []string) (config ConfigFlags) {
	flags := flag.NewFlagSet("login", flag.ExitOnError)
	flags.StringVar(&config.Server.Port, "port", "8080", "Server port")
	flags.StringVar(&config.AuthServer.Host, "host", appData.KCHost, "Keycloak host")
	flags.StringVar(&config.AuthServer.ClientId, "client-id", "frontend-cli", "Client Id")
	flags.StringVar(&config.AuthServer.Realm, "realm", appData.KCRealm, "Keycloak realm")
	flags.StringVar(&config.User, "user", appData.User, "User session")
	flags.Parse(args)

	if config.AuthServer.Host == "" {
		log.Fatalln("Mandatory auth-host flag missing")
	}
	return
}

func StartLoginWebServer(authServer auth.AuthServer, localPort string) chan auth.UserToken {
	authInfo := auth.NewAuthInfo(authServer.Host, authServer.Realm)
	userTokenC := make(chan auth.UserToken)
	go func() {
		log.Println("Start login server")
		log.Println("Server: localhost:" + localPort)
		http.HandleFunc("GET /auth", func(w http.ResponseWriter, r *http.Request) {
			log.Println("Authentication...")
			w.WriteHeader(200)
			w.Write(AUTH_PAGE_HTML)
		})
		http.HandleFunc("GET /code", func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(200)
			w.Write(CODE_PAGE_HTML)
			userTokenC <- auth.UserToken{
				Code:         r.URL.Query().Get("code"),
				Iss:          r.URL.Query().Get("iss"),
				SessionState: r.URL.Query().Get("session_state"),
				State:        r.URL.Query().Get("state"),
			}
		})
		http.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
			log.Println("Login request.")
			http.Redirect(w, r, authInfo.ToURI(authServer, localPort), 302)
		})
		error := http.ListenAndServe("localhost:"+localPort, nil)
		log.Printf("Unable to start login server. %s", error.Error())
		close(userTokenC)
	}()
	return userTokenC
}
