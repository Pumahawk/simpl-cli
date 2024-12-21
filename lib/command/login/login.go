package login

// "https://t1.authority.dev.aruba-simpl.cloud/auth"

import (
	"encoding/json"
	"flag"
	"log"
	"net/http"
	"os"

	app "github.com/pumahawk/simplcli/lib/application"
)

func Exec(conf app.Data, args []string) {
	config := ReadConfigFlag(args)
	userTokenC := StartLoginWebServer(config.AuthServer, config.Server.Port)
	userToken, ok := <-userTokenC
	if !ok {
		log.Fatal("Unable to read from token. Channel is closed.")
	}
	err := json.NewEncoder(os.Stdout).Encode(userToken)
	if err != nil {
		log.Fatalf("Unable to encode token to stdout. %s", err)
	}
}

func ReadConfigFlag(args []string) (config ConfigFlags) {
	flags := flag.NewFlagSet("login", flag.ExitOnError)
	flags.StringVar(&config.Server.Port, "port", "8080", "Server port")
	flags.StringVar(&config.AuthServer.Host, "auth-host", "", "Authentication server host")
	flags.StringVar(&config.AuthServer.ClientId, "auth-client-id", "frontend-cli", "Client Id")
	flags.StringVar(&config.AuthServer.Realm, "realm", "authority", "Keycloak realm")
	flags.StringVar(&config.AuthServer.Realm, "user", "default", "User session")
	flags.Parse(args)

	if config.AuthServer.Host == "" {
		log.Fatalln("Mandatory auth-host flag missing")
	}
	return
}

func StartLoginWebServer(authServer AuthServer, localPort string) chan UserToken {
	authInfo := NewAuthInfo(authServer.Host)
	userTokenC := make(chan UserToken)
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
			userTokenC <- UserToken{
				Code:         r.URL.Query().Get("code"),
				Iss:          r.URL.Query().Get("iss"),
				SessionState: r.URL.Query().Get("session_state"),
				State:        r.URL.Query().Get("state"),
			}
		})
		http.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
			log.Println("Login request.")
			http.Redirect(w, r, authInfo.ToURI(authServer, localPort), 301)
		})
		error := http.ListenAndServe("localhost:"+localPort, nil)
		log.Printf("Unable to start login server. %s", error.Error())
		close(userTokenC)
	}()
	return userTokenC
}
