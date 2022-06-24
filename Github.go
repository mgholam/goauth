package goauth

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type GitHub struct {
	Config
}

func NewGitHub(c Config) GitHub {
	return GitHub{c}
}

func (c *GitHub) GetLoginURL() string {
	return fmt.Sprintf("https://github.com/login/oauth/authorize?client_id=%s&redirect_uri=%s", c.ClientID, c.CallbackURL)
}

func (c *GitHub) Authenticate(r *http.Request) (string, error) {

	err := r.ParseForm()
	if err != nil {
		return "", err
	}
	code := r.FormValue("code")
	reqURL := fmt.Sprintf("https://github.com/login/oauth/access_token?client_id=%s&client_secret=%s&code=%s", c.ClientID, c.ClientSecret, code)
	req, err := http.NewRequest(http.MethodPost, reqURL, nil)
	if err != nil {
		return "", err
	}
	req.Header.Set("accept", "application/json")

	// Send out the HTTP request
	httpClient := http.Client{}
	res, err := httpClient.Do(req)
	if err != nil {
		return "", err
	}

	// Parse the request body into the `oAuthAccessResponse` struct
	var t oAuthAccessResponse
	if err := json.NewDecoder(res.Body).Decode(&t); err != nil {
		return "", err
	}
	res.Body.Close()
	req, _ = http.NewRequest(http.MethodGet, "https://api.github.com/user", nil)
	req.Header.Set("accept", "application/json")
	req.Header.Set("Authorization", "token "+t.AccessToken)
	res, err = httpClient.Do(req)
	if err != nil {
		return "", err
	}
	b, _ := ioutil.ReadAll(res.Body)
	res.Body.Close()
	return string(b), nil
}
