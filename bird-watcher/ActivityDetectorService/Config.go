package main

import (
	"encoding/json"
	"io/ioutil"
	"os"
)

type Config struct {
	WebcamURL               string  `json:"webcam-url"`
	FrameBufferSize         int     `json:"frame-buffer-size"`
	ActivityThreshold       float32 `json:"activity-threshold"`
	ConsecutiveActiveFrames int     `json:"consecutive-active-frames-required"`
	OutputDir               string  `json:"output-dir"`
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
