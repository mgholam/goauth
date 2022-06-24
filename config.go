package goauth

type Config struct {
	ClientID     string
	ClientSecret string
	CallbackURL  string
}

type OAuthAccessResponse struct {
	AccessToken string `json:"access_token"`
}
