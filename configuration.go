package main

import (
	"encoding/json"
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

const configFilePath = "data/configuration.json"

func GetConfiguration() (Config, error) {
	file, err := getOrCreateConfigurationFile()
	if err != nil {
		log.Fatal(err)
		return Config{}, err
	}
	defer file.Close()

	var config Config
	err = ReadJson(file, &config)
	
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
	if FileExists(configFilePath) {
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