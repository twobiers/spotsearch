package main

import (
	"encoding/json"
	"errors"
	"io"
	"log"
	"os"
	"time"
)

type SpotifyConfig struct {
	AccessToken  string    `json:"access-token"`
	RefreshToken string    `json:"refresh-token"`
	TokenType    string    `json:"token-type"`
	Expiry       time.Time `json:"expiry"`
}

type Config struct {
	Spotify      SpotifyConfig    `json:"spotify"`
}

var defaultConfig = Config{
	Spotify:      SpotifyConfig{},
}

const configFilePath = "configuration.json"

func GetConfiguration() (Config, error) {
	file, err := getOrCreateConfigurationFile()
	if err != nil {
		log.Fatal(err)
		return Config{}, err
	}
	defer file.Close()

	configContentBytes, _ := io.ReadAll(file)
	var config Config
	err = json.Unmarshal(configContentBytes, &config)
	
	if err != nil {
		log.Fatal(err)
		return Config{}, err
	}

	return config, nil
}

func SaveConfiguration(config Config) error {
	file, fileErr := getOrCreateConfigurationFile()

	if fileErr != nil {
		return fileErr
	}

	err := writeConfiguration(file, config)
	if err != nil {
		log.Fatal(err)
		return err
	}

	return nil
}

func writeConfiguration(file *os.File, config Config) error {
	configBytes, _ := json.MarshalIndent(config, "", "  ")
	_ = file.Truncate(0)
	_, _ = file.Write(configBytes)

	return nil
}

func getOrCreateConfigurationFile() (*os.File, error) {
	if configExists() {
		file, err := os.OpenFile(configFilePath, os.O_RDWR, 0755)
		if err != nil {
			log.Fatal(err)
			return nil, err
		}
		return file, nil
	}

	file, err := os.Create(configFilePath)

	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	err = writeConfiguration(file, defaultConfig)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	return file, nil
}

func configExists() bool {
	_, statErr := os.Stat(configFilePath)
	if errors.Is(statErr, os.ErrNotExist) {
		return false
	}
	return true
}