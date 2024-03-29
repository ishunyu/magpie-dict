package main

import (
	"encoding/json"
	"fmt"
	"os"
)

// Config used for storing app configuration
type Config struct {
	Hostname        string `json:"hostname"`
	Port            int    `json:"port"`
	LoggingPath     string `json:"loggingPath"`
	HtmlPath        string `json:"htmlPath"`
	DataPath        string `json:"dataPath"`
	IndexPath       string `json:"indexPath"`
	ComparePath     string `json:"comparePath"`
	CompareVenvPath string `json:"compareVenvPath"`
	TmpPath         string `json:"tmpPath"`
}

// GetConfig returns config data based in json
func GetConfig() *Config {
	args := os.Args[1:]
	if len(args) == 0 {
		fmt.Println("Missing configuration file argument.")
		os.Exit(1)
	}

	configFilePath := args[0]
	fmt.Println("Loading config from ", configFilePath)
	var config Config
	err := JsonLoadFromFile(configFilePath, &config)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	configStr, _ := json.MarshalIndent(config, "", "  ")
	fmt.Printf("Config loaded: %s\n", string(configStr))

	return &config
}
