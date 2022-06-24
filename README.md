# goauth
A no dependancy oauth client. 


## Sample code
- on success will return the json object for the login user
```go
package main

import (
	"fmt"
	"net/http"

	"github.com/mgholam/goauth"
)

func main() {
	goog := goauth.NewGoogle(goauth.Config{
		ClientID:     "",
		ClientSecret: "",
		CallbackURL:  "http://localhost:3000/auth/google/callback",
	})

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`<a href="/login">Login with Google</a>`))
	})

	// login redirect to google oauth
	http.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, goog.GetLoginURL(), http.StatusTemporaryRedirect)
	})

	// google callback handler
	http.HandleFunc("/auth/google/callback", func(w http.ResponseWriter, r *http.Request) {
		json, err := goog.Authenticate(r)
		if err != nil {
			fmt.Println(err)
			return
		}
		w.Header().Set("content-type", "application/json")
		w.Write([]byte(json))
	})

	http.ListenAndServe(":3000", nil)
} 
```


