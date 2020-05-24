package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
)

// Config used for storing app configuration
type Config struct {
	RootPath  string `json:"rootPath"`
	IndexPath string `json:"indexPath"`
}

// GetConfig returns config data based in json
func GetConfig() *Config {
	args := os.Args[1:]
	if len(args) == 0 {
		fmt.Println("Missing configuration file argument.")
		os.Exit(1)
	}

	jsonFile, err := os.Open(args[0])
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer jsonFile.Close()
	bytes, _ := ioutil.ReadAll(jsonFile)

	var config Config
	json.Unmarshal([]byte(bytes), &config)

	fmt.Printf("%+v\n", config)

	return &config
}

func (config *Config) GetDataPath() string {
	return filepath.Join(config.RootPath, "resource", "data", "EP21Outfile.csv")
}

func (config *Config) GetHtmlDir() string {
	return filepath.Join(config.RootPath, "static")
}
