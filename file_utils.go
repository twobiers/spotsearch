package main

import (
	"encoding/json"
	"errors"
	"io"
	"os"
)

func ReadJson(file *os.File, v any) error {
	contentBytes, _ := io.ReadAll(file)
	return json.Unmarshal(contentBytes, &v)
}

func FileExists(path string) bool {
	_, statErr := os.Stat(path)
	if errors.Is(statErr, os.ErrNotExist) {
		return false
	}
	return true
}