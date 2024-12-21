package token

import (
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
	"net/url"

	app "github.com/pumahawk/simplcli/lib/application"
	"github.com/pumahawk/simplcli/lib/command/login"
)

func Exec(conf app.Data, args []string) {
	log.Panicln("Not yet implemented")
}

func Tokenize(authServer login.AuthServer, localPort string, token login.UserToken) (tokenInfo TokenInfo, err error) {
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
