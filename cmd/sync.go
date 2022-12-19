package cmd

import (
	"context"
	"log"

	"github.com/spf13/cobra"
	client "github.com/twobiers/spotsearch/internal/pkg/client"
	data "github.com/twobiers/spotsearch/internal/pkg/data"
	"github.com/zmb3/spotify/v2"
)

var syncCmd = &cobra.Command{
	Use:   "sync",
	Run: func(cmd *cobra.Command, args []string) {
	  synchronize()
	},
  }

func init() {
	rootCmd.AddCommand(syncCmd)
}

func synchronize() {
	state, err := data.LoadState()

	if err != nil {
		log.Fatal(err)
	}

	outdated := getOutdatedPlaylists(state)

	for _, p := range outdated {
		id := p.ID.String()
		items, err := fetchPlaylistItems(p.ID)
		if err != nil {
			log.Println("Could not fetch playlist " + id)
			continue
		}

		log.Printf("Fetched %d items for playlist '%s'", len(items), id)

		itemIds := make([]string, len(items))
		for i, e := range items {
			if e.Track.Track != nil {
				itemIds[i] = e.Track.Track.ID.String()
			}
		}

		if state.Playlists == nil {
			state.Playlists = make(map[string]data.PlaylistState)
		}

		state.Playlists[id] = data.PlaylistState{
			Id: id,
			SnapshotId: p.SnapshotID,
			Name: p.Name,
			Tracks: itemIds,
		}
	}

	data.SaveState(state)
}

func getOutdatedPlaylists(state data.State) ([]spotify.SimplePlaylist) {
	playlists, err := getPlaylists()
	if err != nil {
		log.Fatal(err)
	}
	
	outdated := make([]spotify.SimplePlaylist, 0)
	for _, playlist := range playlists {
		id := playlist.ID
		p, prs := state.Playlists[id.String()]
		if !prs || p.SnapshotId != playlist.SnapshotID {
			log.Printf("Playlist '%s' is outdated (Local Snapshot: %s, New Snapshot: %s)\n", id.String(), p.SnapshotId, playlist.SnapshotID)
			outdated = append(outdated, playlist)
		}
	}
	log.Printf("Identified %d outdated playlists", len(outdated))
	return outdated
}

func fetchPlaylistItems(id spotify.ID) ([]spotify.PlaylistItem, error) {
	client := client.GetSpotifyClient()

	total := 20
	offset := 0
	items := make([]spotify.PlaylistItem, 0)

	client.GetPlaylistItems(context.Background(), id)

	for total != len(items) {
		itemPage, err := client.GetPlaylistItems(context.Background(), id, spotify.Limit(50), spotify.Offset(offset))
		if err != nil {
			return nil, err
		}

		items = append(items, itemPage.Items...)
		total = itemPage.Total
		offset += 50
	}

	return items, nil
}

func getPlaylists() ([]spotify.SimplePlaylist, error) {
	client := client.GetSpotifyClient()

	total := 20
	offset := 0
	items := make([]spotify.SimplePlaylist, 0)

	for total != len(items) {
		playlistPage, err := client.CurrentUsersPlaylists(context.Background(), spotify.Limit(50), spotify.Offset(offset))
		if err != nil {
			return nil, err
		}

		items = append(items, playlistPage.Playlists...)
		total = playlistPage.Total
		offset += 50
	}

	return items, nil
}