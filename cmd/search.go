package cmd

import (
	"context"
	"log"
	"strings"

	"github.com/spf13/cobra"
	client "github.com/twobiers/spotsearch/internal/pkg/client"
	data "github.com/twobiers/spotsearch/internal/pkg/data"
	"github.com/zmb3/spotify/v2"
)

var searchCmd = &cobra.Command{
	Use:   "search",
	Run: func(cmd *cobra.Command, args []string) {
	  search(strings.Join(args, " "))
	},
  }

func init() {
	rootCmd.AddCommand(searchCmd)
}

func search(query string) {
	state, _ := data.LoadState()
	client := client.GetSpotifyClient()
	result, err := client.Search(context.Background(), query, spotify.SearchTypeTrack)
	if err != nil {
		log.Fatal(err)
	}

	if result.Tracks.Total == 0 {
		log.Println("No track found")
	}

	track := result.Tracks.Tracks[0]

	containingPlaylists := make([]data.PlaylistState, 0)
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
