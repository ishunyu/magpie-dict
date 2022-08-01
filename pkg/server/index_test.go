package main

import (
	"os"
	"path/filepath"
	"testing"
)

func TestGetIndex(t *testing.T) {
	testFolder := "get_index"
	defer clean(t, testFolder)

	index, err := NewIndex(getPaths(t, testFolder))
	if err != nil {
		t.Fatal(err)
	}

	var results = index.Search("name", "testshow")
	if len(results) != 1 {
		t.Fatalf("len(results): %d", len(results))
	}

	var result = *results[0]
	if result.String() != "testshow.01.0" {
		t.Fatalf("result: %v", result)
	}
}

func TestBadIndexPath(t *testing.T) {
	_, err := NewIndex(getPaths(t, "bad_data"))
	if err == nil {
		t.Fatalf("Expected error.")
	}
	t.Log(err)
}

func getPaths(t *testing.T, testFolder string) (string, string) {
	testPath := "testdata/" + testFolder
	absoluteTestPath, err := filepath.Abs(testPath)
	if err != nil {
		t.Fatal(err)
	}

	return filepath.Join(absoluteTestPath, "data"), filepath.Join(absoluteTestPath, "tmp/index")
}

func clean(t *testing.T, name string) {
	_, indexPath := getPaths(t, name)
	err := os.RemoveAll(indexPath)
	if err != nil {
		t.Fatal(err)
	}
}
