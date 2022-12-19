package main

import (
	"context"
	"log"

	"github.com/google/uuid"
	"github.com/zmb3/spotify/v2"
	spotifyauth "github.com/zmb3/spotify/v2/auth"

	"net/http"
)

const port = 8080
const addr = ":8080"
const redirectUri = "http://localhost:8080/callback"

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
	log.Println("Login completed")
	ch <- client
}

func main() {
	http.HandleFunc("/callback", CallbackHandler)

	log.Printf("server is listening at %s...", addr)
    go http.ListenAndServe(addr, nil)

	url := auth.AuthURL(state)
	log.Println(url)

	client := <-ch
	
	user, err := client.CurrentUser(context.TODO())
	if err != nil {
		log.Fatal(err)
	}
	log.Println("You are logged in as:", user.ID)
}