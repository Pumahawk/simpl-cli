package auth

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"time"

	app "github.com/pumahawk/simplcli/lib/application"
)

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
	tokenInfo.UpdateExpirationTime(time.Now())
	return
}

func getUserAuthFileName(appData app.Data, user string) string {
	return appData.DirData + "/" + user + ".json"
}

func LoadUserAuthData(appData app.Data, user string) (tokenInfo TokenInfo, err error) {
	fileName := getUserAuthFileName(appData, user)
	file, err := os.Open(fileName)
	if err != nil {
		err = errors.New(fmt.Sprintf("Unable to open user file %s", fileName))
		return
	}
	err = json.NewDecoder(file).Decode(&tokenInfo)
	return
}

func SaveUserAuthData(appData app.Data, user string, tokenInfo TokenInfo) error {
	log.Printf("Save token. User %s", user)
	fileName := getUserAuthFileName(appData, user)
	file, err := os.Create(fileName)
	if err != nil {
		return errors.New(fmt.Sprintf("Unable to create user file %s", fileName))
	}
	err = json.NewEncoder(file).Encode(tokenInfo)
	return err
}

func ReloadToken(authServer AuthServer, tokenInfo TokenInfo) (result TokenInfo, err error) {
	log.Println("Reload Token")
	urlValues := url.Values{}
	urlValues.Add("refresh_token", tokenInfo.RefreshToken)
	urlValues.Add("grant_type", "refresh_token")
	urlValues.Add("client_id", authServer.ClientId)
	resp, err := http.PostForm(authServer.Host+"/realms/"+authServer.Realm+"/protocol/openid-connect/token", urlValues)
	if err != nil {
		return
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		log.Printf("Invalid response code of tocanization. %s", resp.Status)
		err = errors.New("Invalid response code")
		body, errReadBody := io.ReadAll(resp.Body)
		if errReadBody != nil {
			err = errors.Join(err, errReadBody)
			return
		}
		log.Printf("Body: %s", body)
		return
	}
	json.NewDecoder(resp.Body).Decode(&result)
	return
}
