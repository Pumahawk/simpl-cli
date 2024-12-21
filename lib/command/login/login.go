package login

import (
	"encoding/json"
	"errors"
	"flag"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
)

func Exec(args []string) {
	flags := flag.NewFlagSet("login", flag.ExitOnError)
	authInfo := NewAuthInfo("https://t1.authority.dev.aruba-simpl.cloud/auth")
	flags.Parse(args)
	userTokenC := StartLoginWebServer(authInfo)
	userToken, ok := <-userTokenC
	if !ok {
		log.Fatal("Unable to read from token. Channel is closed.")
	}
	token, err := Tokenize(userToken)
	if err != nil {
		log.Fatalf("Unable to tokenize. %s", err.Error())
	}
	err = json.NewEncoder(os.Stdout).Encode(token)
	if err != nil {
		log.Fatalf("Unable to encode token to stdout. %s", err)
	}
}

func StartLoginWebServer(authInfo AuthInfo) chan UserToken {
	userTokenC := make(chan UserToken)
	go func() {
		log.Println("Start login server")
		log.Println("Server: localhost:8080")
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
			http.Redirect(w, r, authInfo.ToURI(), 301)
		})
		error := http.ListenAndServe("localhost:8080", nil)
		log.Printf("Unable to start login server. %s", error.Error())
		close(userTokenC)
	}()
	return userTokenC
}

func Tokenize(token UserToken) (tokenInfo TokenInfo, err error) {
	values := url.Values{}
	NewTokenizeInfo(token.Code).ToUrlValues(&values)
	log.Println("Tokenize...")
	r, err := http.PostForm("https://t1.authority.dev.aruba-simpl.cloud/auth/realms/authority/protocol/openid-connect/token", values)
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
