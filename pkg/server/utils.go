package main

import (
	"encoding/json"
	"io/ioutil"
	"os"
)

func Min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func ReadFile(filePath string) ([]byte, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	return ioutil.ReadAll(file)
}

func WriteFile(filePath string, bytes []byte) error {
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()
	_, writeErr := file.Write(bytes)
	if writeErr != nil {
		return writeErr
	}
	return nil
}

func JsonLoadFromFile(filePath string, v interface{}) error {
	bytes, err := ReadFile(filePath)
	if err != nil {
		return err
	}

	json.Unmarshal(bytes, v)
	return nil
}

func JsonWriteToFile(filePath string, v interface{}) error {
	bytes, err := json.Marshal(v)
	if err != nil {
		return err
	}
	return WriteFile(filePath, bytes)
}
