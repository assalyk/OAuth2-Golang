package handlers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/davecgh/go-spew/spew"
	"golang.org/x/oauth2"
)

var (
	SpotifyOauthConf = &oauth2.Config{
		RedirectURL:  "http://localhost:8080/spotify/callback",
		ClientID:     "",
		ClientSecret: "",
		Scopes:       []string{"user-read-private", "user-read-email"},
		Endpoint: oauth2.Endpoint{
			AuthURL:  "https://accounts.spotify.com/authorize",
			TokenURL: "https://accounts.spotify.com/api/token",
		},
	}
)

func HandleSpotifyLogin(w http.ResponseWriter, r *http.Request) {
	url := SpotifyOauthConf.AuthCodeURL(oauthStateString)
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

func HandleSpotifyCallback(w http.ResponseWriter, r *http.Request) {
	state := r.FormValue("state")
	if state != oauthStateString {
		fmt.Printf("invalid oauth state, expected %q got %q", oauthStateString, state)
		return
	}

	code := r.FormValue("code")
	fmt.Printf("this is the code %s\n", code)
	token, err := SpotifyOauthConf.Exchange(oauth2.NoContext, code)
	if err != nil {
		fmt.Printf("SpotifyOauthConf.Exchange() failed with %q", err)
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}
	fmt.Printf("this is the token %s\n", token.AccessToken)
	var bearer = "Bearer " + token.AccessToken
	req, err := http.NewRequest("GET", "https://api.spotify.com/v1/me", nil)
	req.Header.Add("Authorization", bearer)
	client := &http.Client{}
	info, err := client.Do(req)
	if err != nil {
		fmt.Printf("Get: %q", err)
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}
	defer info.Body.Close()
	response, err := ioutil.ReadAll(info.Body)
	if err != nil {
		fmt.Printf("ReadAll: %q", err)
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}
	var m map[string]interface{}
	if err := json.Unmarshal(response, &m); err != nil {
		fmt.Printf("error unmarshalling response: %s", err)
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
	}
	spew.Dump(m)
	Fname := m["localizedFirstName"]
	Lname := m["localizedLastName"]
	fmt.Println(Fname, Lname)
	http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
}
