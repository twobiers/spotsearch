package main

import (
	"context"
	"log"

	"github.com/google/uuid"
	"github.com/zmb3/spotify/v2"
	spotifyauth "github.com/zmb3/spotify/v2/auth"

	"golang.org/x/oauth2"

	"net/http"
)

const (
	addr = ":8080"
	redirectUri = "http://localhost" + addr + "/callback"
)

var (
	auth  = spotifyauth.New(spotifyauth.WithRedirectURL(redirectUri), spotifyauth.WithScopes(spotifyauth.ScopeUserReadPrivate, spotifyauth.ScopeUserLibraryRead))
	state = uuid.New().String()
	ch    = make(chan *spotify.Client)
)

func CallbackHandler(w http.ResponseWriter, r *http.Request) {
	token, err := auth.Token(r.Context(), state, r)

	if err != nil {
		http.Error(w, "Couldn't get token", http.StatusNotFound)
		return
	}

	client := spotify.New(auth.Client(r.Context(), token))
	
	config, _ := GetConfiguration()

	config.Spotify.AccessToken = token.AccessToken
	config.Spotify.RefreshToken = token.RefreshToken
	config.Spotify.Expiry = token.Expiry
	config.Spotify.TokenType = token.TokenType
	
	_ = SaveConfiguration(config)

	ch <- client
}

func main() {
	client, err := tryLoadClient()
	token, tokenErr := client.Token()

	if tokenErr != nil || !token.Valid() {
		log.Println("A token couldn't be created from the stored configuration")
		
		// We're going to start the server only if we can't use the tokens from configuration
		http.HandleFunc("/callback", CallbackHandler)

		log.Printf("server is listening at %s...", addr)
		go http.ListenAndServe(addr, nil)

		url := auth.AuthURL(state)
		log.Println(url)
		client = <-ch
	} 
	
	user, err := client.CurrentUser(context.TODO())
	if err != nil {
		log.Fatal(err)
	}
	log.Println("You are logged in as:", user.ID)
}


func tryLoadClient() (*spotify.Client, error) {
	config, err := GetConfiguration()
	if err != nil {
		return nil, err
	}

	token := oauth2.Token{
		AccessToken:  config.Spotify.AccessToken,
		RefreshToken: config.Spotify.RefreshToken,
		TokenType:    config.Spotify.TokenType,
		Expiry:       config.Spotify.Expiry,
	}
	httpClient := auth.Client(context.Background(), &token)
	return spotify.New(httpClient, spotify.WithRetry(true)), nil
}