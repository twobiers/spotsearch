package client

import (
	"context"
	"log"

	"github.com/google/uuid"
	data "github.com/twobiers/spotsearch/internal/pkg/data"
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
	auth  = spotifyauth.New(spotifyauth.WithRedirectURL(redirectUri), spotifyauth.WithScopes(
		spotifyauth.ScopeUserReadPrivate, 
		spotifyauth.ScopeUserLibraryRead,
		spotifyauth.ScopePlaylistReadPrivate,
		spotifyauth.ScopePlaylistReadCollaborative,
		))
	state = uuid.New().String()
	ch    = make(chan *spotify.Client)
	client *spotify.Client
)

func CallbackHandler(w http.ResponseWriter, r *http.Request) {
	token, err := auth.Token(r.Context(), state, r)

	if err != nil {
		http.Error(w, "Couldn't get token", http.StatusNotFound)
		return
	}

	client := spotify.New(auth.Client(r.Context(), token))
	
	config, _ := data.GetConfiguration()

	config.Spotify.AccessToken = token.AccessToken
	config.Spotify.RefreshToken = token.RefreshToken
	config.Spotify.Expiry = token.Expiry
	config.Spotify.TokenType = token.TokenType
	
	_ = data.SaveConfiguration(config)

	ch <- client
}

func Authenticate() *spotify.Client {		
	// We're going to start the server only if we can't use the tokens from configuration
	http.HandleFunc("/callback", CallbackHandler)

	log.Printf("Server is listening at %s...", addr)
	go http.ListenAndServe(addr, nil)

	url := auth.AuthURL(state)
	log.Println("Please visit the following URL to authenticate: " + url)
	client = <-ch
	return client
}

func TestAuth() bool {
	client, err := GetSpotifyClient()
	token, tokenErr := client.Token()

	if err != nil || tokenErr != nil || !token.Valid() {
		return false
	}

	return true
}

func GetSpotifyClient() (*spotify.Client, error) {
	if client != nil {
		_, err := client.Token()
		if err != nil {
			return client, nil
		}
	}

	config, err := data.GetConfiguration()
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