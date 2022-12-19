package main

import (
	"encoding/json"
	"log"
	"os"
)

type PlaylistState struct {
	Id         string   `json:"id"`
	SnapshotId string   `json:"snapshot_id"`
	Name       string   `json:"name"`
	Tracks     []string `json:"tracks"`
}

type State struct {
	Playlists map[string]PlaylistState `json:"playlists"`
}

var defaultState = State{
	Playlists: map[string]PlaylistState{},
}

const stateFilePath = "data/state.json"

func LoadState() (State, error) {
	file, err := getOrCreateStateFile()
	defer file.Close()

	if err != nil {
		log.Fatal(err)
		return State{}, err
	}

	var state State
	err = ReadJson(file, &state)
	if err != nil {
		log.Fatal(err)
		return State{}, err
	}

	return state, nil
}

func SavePlaylistState(state PlaylistState) error {
	loadedState, err := LoadState()
	if err != nil {
		return err
	}
	loadedState.Playlists[state.Id] = state
	return SaveState(loadedState)
}

func SaveState(state State) error {
	file, fileErr := getOrCreateStateFile()

	if fileErr != nil {
		return fileErr
	}

	err := writeState(file, state)
	if err != nil {
		log.Fatal(err)
		return err
	}

	return nil
}


func writeState(file *os.File, state State) error {
	stateBytes, _ := json.MarshalIndent(state, "", "  ")
	_ = file.Truncate(0)
	_, _ = file.Write(stateBytes)

	return nil
}

func getOrCreateStateFile() (*os.File, error) {
	if FileExists(stateFilePath) {
		file, err := os.OpenFile(stateFilePath, os.O_RDWR, 0755)
		if err != nil {
			log.Fatal(err)
			return nil, err
		}
		return file, nil
	}

	file, err := os.Create(stateFilePath)

	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	err = writeState(file, defaultState)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	return file, nil
}