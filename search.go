package main

import (
	"context"
	"log"

	"github.com/zmb3/spotify/v2"
)

func Search(query string) {
	state, _ := LoadState()
	client := GetSpotifyClient()
	result, err := client.Search(context.Background(), query, spotify.SearchTypeTrack)
	if err != nil {
		log.Fatal(err)
	}

	if result.Tracks.Total == 0 {
		log.Println("No track found")
	}

	track := result.Tracks.Tracks[0]

	containingPlaylists := make([]PlaylistState, 0)
	for _, value := range state.Playlists {
		if contains(value.Tracks, track.ID.String()) {
			containingPlaylists = append(containingPlaylists, value)
		}
	}

	if len(containingPlaylists) == 0 {
		log.Println("Not found in any playlist")
		return
	}

	log.Println("Track was found in the following playlists:")
	for _, playlist := range containingPlaylists {
		log.Print(playlist.Name)
	}
}

func contains(s []string, str string) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}

	return false
}
