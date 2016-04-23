package oauth

import (
	// "golang.org/x/oauth2"
	"github.com/nyelonong/finapimate/utils"
	"encoding/json"
	"net/http"
	"log"
	"encoding/base64"
	"net/url"
	"time"
	"bytes"
	"io/ioutil"
	"strings"
)

type AccessResponse struct {
	AccessToken string `json:"access_token"`
	TokenType string `json:"token_type"`
	Scope string `json:"scope"`
}

type OAuthModule struct {
	AccessToken string
}

var o *OAuthModule

func init() {
	o = &OAuthModule{}
}

func GetTime() string {
	now := time.Now().Format(time.RFC3339Nano)
	// hack time format.
	arr := strings.Split(now, ".")
	now = strings.Join([]string{arr[0], arr[1][len(arr[1])-9:]}, ".")
	return now
}

// Get access token from API if not yet existed.
// TODO: handle refresh token.
func GetAccessToken() (string, error) {

	if o.AccessToken != "" {
		return o.AccessToken, nil
	}

	token, err := getAccessTokenManually()
	if err != nil {
		log.Println(err)
		return "", err
	}

	o.AccessToken = token
	return token, nil
}

func getAccessTokenManually() (string, error) {
	auth := base64.StdEncoding.EncodeToString([]byte(utils.ClientId + ":" + utils.ClientSecret))

    data := url.Values{}
    data.Set("grant_type", "client_credentials")

    req, err := http.NewRequest("POST", utils.BcaHost + utils.BcaTokenPath, bytes.NewBufferString(data.Encode()))
    if err != nil {
    	log.Println(err)
    	return "", err
    }
    req.Header.Add("Authorization", "Basic: " + auth)
    req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

    client := &http.Client{}
    resp, err := client.Do(req)
    if err != nil {
    	log.Println(err)
    	return "", err
    }

    contents, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
		return "", err
	}

	var ar AccessResponse
	err = json.Unmarshal(contents, &ar)
	if err != nil {
		log.Println(err)
		return "", err
	}

	log.Println(ar.AccessToken)
	return ar.AccessToken, nil
}