package main

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"

	log "github.com/sirupsen/logrus"
)

// Config used for storing app configuration
type Config struct {
	Hostname        string `json:"hostname"`
	Port            int    `json:"port"`
	RootPath        string `json:"rootPath"`
	IndexPath       string `json:"indexPath"`
	DataPath        string `json:"dataPath"`
	TempPath        string `json:"tmpPath"`
	ComparePath     string `json:"comparePath"`
	CompareVenvPath string `json:"compareVenvPath"`
}

// GetConfig returns config data based in json
func GetConfig() *Config {
	args := os.Args[1:]
	if len(args) == 0 {
		log.Error("Missing configuration file argument.")
		os.Exit(1)
	}

	jsonFile, err := os.Open(args[0])
	if err != nil {
		log.Error(err)
		os.Exit(1)
	}
	defer jsonFile.Close()
	bytes, _ := ioutil.ReadAll(jsonFile)

	var config Config
	json.Unmarshal([]byte(bytes), &config)

	configStr, _ := json.MarshalIndent(config, "", "  ")

	log.Infof(string(configStr))

	return &config
}

func (config *Config) GetDataPath() string {
	return filepath.Join(config.RootPath, "resource", "data")
}

func (config *Config) GetHtmlDir() string {
	return filepath.Join(config.RootPath, "static")
}

func (config *Config) GetPort() int {
	if config.Port == 0 {
		return 8090
	}
	return config.Port
}
