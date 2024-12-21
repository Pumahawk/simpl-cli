package token

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
	"github.com/pumahawk/simplcli/lib/command/login"
)

func Exec(conf app.Data, args []string) {
	flags := ReadConfigFlag(args)
	userToken, err := loadUserToken(conf, flags.User)
	if err != nil {
		log.Fatalf("Unable to read user token. %s", err.Error())
	}
	tokenInfo, err := Tokenize(flags.AuthServer, flags.Port, userToken)
	if err != nil {
		log.Fatalf("Unable to tokenize. %s", err.Error())
	}
	err = json.NewEncoder(os.Stdout).Encode(tokenInfo)
	if err != nil {
		log.Fatalf("Unable to encode token. %s", err.Error())
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

func Tokenize(authServer AuthServer, localPort string, token login.UserToken) (tokenInfo TokenInfo, err error) {
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

func loadUserToken(appData app.Data, user string) (userToken login.UserToken, err error) {
	fileName := appData.DirData + "/" + user + ".json"
	file, err := os.Open(fileName)
	if err != nil {
		err = errors.New(fmt.Sprintf("Unable to open user file %s. %s", fileName, err.Error()))
		return
	}
	err = json.NewDecoder(file).Decode(&userToken)
	return
}
