package login

import (
	"flag"
	"log"
	"net/http"
)

func Exec(args []string) {
	flags := flag.NewFlagSet("login", flag.ExitOnError)
	authInfo := NewAuthInfo("https://t1.authority.dev.aruba-simpl.cloud/auth")
	flags.Parse(args)
	userTokenC := StartLoginWebServer(authInfo)
	token, ok := <-userTokenC
	if !ok {
		log.Fatal("Unable to read from token. Channel is closed.")
	}
	log.Printf("Token: %s", token)
}

func StartLoginWebServer(authInfo AuthInfo) chan UserToken {
	userTokenC := make(chan UserToken)
	go func() {
		log.Println("Start login server")
		http.HandleFunc("GET /auth", func(w http.ResponseWriter, r *http.Request) {
			log.Println("Authentication...")
			w.WriteHeader(200)
			w.Write(AUTH_PAGE_HTML)
		})
		http.HandleFunc("GET /code", func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(200)
			w.Write(CODE_PAGE_HTML)
			userTokenC <- UserToken{
				Code: r.URL.Query().Get("code"),
				Iss: r.URL.Query().Get("iss"),
				SessionState: r.URL.Query().Get("session_state"),
				State: r.URL.Query().Get("state"),
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

