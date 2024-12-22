package login

// "https://t1.authority.dev.aruba-simpl.cloud/auth"

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
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
	tokenInfo, err := Tokenize(config.AuthServer, config.Server.Port, userToken)
	if err != nil {
		log.Fatalf("Unable to tokenize. %s", err.Error())
	}
	err = saveUserAuthData(conf, config.User, tokenInfo)
	if err != nil {
		log.Fatalf("Unable to encode token to stdout. %s", err.Error())
	}
}

func ReadConfigFlag(args []string) (config ConfigFlags) {
	flags := flag.NewFlagSet("login", flag.ExitOnError)
	flags.StringVar(&config.Server.Port, "port", "8080", "Server port")
	flags.StringVar(&config.AuthServer.Host, "auth-host", "", "Authentication server host")
	flags.StringVar(&config.AuthServer.ClientId, "auth-client-id", "frontend-cli", "Client Id")
	flags.StringVar(&config.AuthServer.Realm, "realm", "authority", "Keycloak realm")
	flags.StringVar(&config.User, "user", "default", "User session")
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

func saveUserAuthData(appData app.Data, user string, tokenInfo TokenInfo) error {
	fileName := appData.DirData + "/" + user + ".json"
	file, err := os.Create(fileName)
	if err != nil {
		return errors.New(fmt.Sprintf("Unable to open user file %s", fileName))
	}
	err = json.NewEncoder(file).Encode(tokenInfo)
	return err
}

func Tokenize(authServer AuthServer, localPort string, token UserToken) (tokenInfo TokenInfo, err error) {
	values := url.Values{}
	NewTokenizeInfo(token.Code, localPort).ToUrlValues(&values, localPort)
	log.Println("Tokenize...")
	r, err := http.PostForm(authServer.Host+"/realms/"+authServer.Realm+"/protocol/openid-connect/token", values)
	if err != nil {
		return
	}
	defer r.Body.Close()
	log.Println("Tokenize status: ", r.Status)
	log.Println("Tokenize size: ", r.Header.Get("content-length"))
	body, err := io.ReadAll(r.Body)
	if err != nil {
		return tokenInfo, err
	}
	if r.StatusCode >= 200 && r.StatusCode < 300 {
		err = json.Unmarshal(body, &tokenInfo)
	} else {
		log.Printf("Body: %s", body)
		err = errors.New("Bad tokenization. " + r.Status)
	}
	return
}
