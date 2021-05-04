package main

import (
	"encoding/json"
	"io/ioutil"
	"os"
)

type Config struct {
	Port            int    `json:"stream-port"`
	VideoDevice     string `json:"video-device"`
	Width           int    `json:"width"`
	Height          int    `json:"height"`
	SingleFrameMode bool   `json:"single-frame-mode"`
	MaxSockets      int    `json:"max-sockets"`
}

func parseConfigFile(path string) (Config, error) {

	// Open config file
	file, err := os.Open(path)
	if err != nil {
		return Config{}, err
	}
	defer file.Close()

	// Parse config
	bytes, err := ioutil.ReadAll(file)
	if err != nil {
		return Config{}, nil
	}
	var config Config
	json.Unmarshal(bytes, &config)

	return config, nil
}
