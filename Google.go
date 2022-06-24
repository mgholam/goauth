package goauth

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

const googletokenurl = `https://oauth2.googleapis.com/token`

type Google struct {
	Config
}

func NewGoogle(c Config) Google {
	return Google{c}
}

func (c *Google) GetLoginURL() string {
	return "https://accounts.google.com/o/oauth2/auth?client_id=" + c.ClientID +
		"&redirect_uri=" + c.CallbackURL +
		"&scope=profile email&response_type=code&state=pseudo-random"
}

func (c *Google) Authenticate(r *http.Request) (string, error) {
	err := r.ParseForm()
	if err != nil {
		return "", err
	}
	code := r.FormValue("code")

	varrs := fmt.Sprintf("client_id=%s&client_secret=%s&code=%s&grant_type=authorization_code&redirect_uri=%s", c.ClientID, c.ClientSecret, code, c.CallbackURL)
	req, err := http.NewRequest(http.MethodPost, googletokenurl, strings.NewReader(varrs))
	if err != nil {
		return "", err
	}
	// We set this header since we want the response
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("accept", "application/json")

	// // Send out the HTTP request
	httpClient := http.Client{}
	res, err := httpClient.Do(req)
	if err != nil {
		return "", err
	}

	// Parse the request body into the `OAuthAccessResponse` struct
	var t OAuthAccessResponse
	if err := json.NewDecoder(res.Body).Decode(&t); err != nil {
		return "", err
	}
	res.Body.Close()
	res, _ = http.Get("https://www.googleapis.com/oauth2/v3/userinfo?access_token=" + t.AccessToken)
	b, _ := ioutil.ReadAll(res.Body)
	res.Body.Close()
	return string(b), nil
}
