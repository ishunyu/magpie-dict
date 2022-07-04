package main

import (
	"path/filepath"
	"testing"
)

func TestGetIndex(t *testing.T) {
	var currPath, err = filepath.Abs("testdata/index")
	if err != nil {
		t.Fatal()
	}
	var config = &Config{
		Hostname:        "localhost",
		Port:            9999,
		RootPath:        "",
		IndexPath:       filepath.Join(currPath, "tmp/index"),
		DataPath:        filepath.Join(currPath, "data"),
		TempPath:        filepath.Join(currPath, "tmp"),
		ComparePath:     "",
		CompareVenvPath: "",
	}

	var index = GetIndex(config)
	var results = index.Search("name", "testshow")

	if len(results) != 1 {
		t.Fatalf("len(results): %d", len(results))
	}

	var result = *results[0]
	if result.showID != "testshow" || result.filename != "01" || result.subID != 0 {
		t.Fatalf("result: %v", result)
	}
}
