package main

import (
	"os"
	"path/filepath"
	"testing"
)

type TestPaths struct {
	dataPath  string
	tmpPath   string
	indexPath string
}

func TestBadPaths(t *testing.T) {
	tPaths := testPaths(t, "bad_paths")
	_, err := NewIndex(tPaths.dataPath, tPaths.indexPath)
	if err == nil {
		t.Fatalf("Expected error.")
	}
	t.Log(err)
}

func TestWithoutCacheWithoutManifest(t *testing.T) {
	testPath := "without_cache_without_manifest"
	tPaths := testPaths(t, testPath)
	clean(t, tPaths)
	defer clean(t, tPaths)

	doTest(t, tPaths)
}

func TestWithoutCacheWithManifest(t *testing.T) {
	testPath := "without_cache_with_manifest"
	tPaths := testPaths(t, testPath)
	cleanCache(t, tPaths)
	defer cleanCache(t, tPaths)

	doTest(t, tPaths)
}

func TestWithCacheWithoutManifest(t *testing.T) {
	tPaths := testPaths(t, "with_cache_without_manifest")
	cleanManifest(t, tPaths)
	defer cleanManifest(t, tPaths)

	doTest(t, tPaths)
}

func TestWithCacheWithManifest(t *testing.T) {
	testPath := "with_cache_with_manifest"
	tPaths := testPaths(t, testPath)

	doTest(t, tPaths)
}

func doTest(t *testing.T, tPaths TestPaths) {
	index, err := NewIndex(tPaths.dataPath, tPaths.indexPath)
	if err != nil {
		t.Fatal(err)
	}

	var results = index.Search("name", "show")
	if len(results) != 1 {
		t.Fatalf("len(results): %d", len(results))
	}

	var result = *results[0]
	if result.String() != "show.01.0" {
		t.Fatalf("result: %v", result)
	}
}

func testPaths(t *testing.T, testPathSub string) TestPaths {
	testPath := filepath.Join("testdata", "index", testPathSub)
	testPathAbs, err := filepath.Abs(testPath)
	if err != nil {
		t.Fatal(err)
	}

	return TestPaths{
		filepath.Join(testPathAbs, "data"),
		filepath.Join(testPath, "tmp"),
		filepath.Join(testPathAbs, "tmp", "index"),
	}
}

func clean(t *testing.T, tPaths TestPaths) {
	cleanCache(t, tPaths)
	cleanManifest(t, tPaths)

	err := os.RemoveAll(tPaths.tmpPath)
	if err != nil {
		t.Fatal(err)
	}

	err = os.Remove(tPaths.tmpPath)
	if err != nil && !os.IsNotExist(err) {
		t.Fatal(err)
	}
}

func cleanCache(t *testing.T, tPaths TestPaths) {
	blevePath := filepath.Join(tPaths.indexPath, "bleve")

	err := os.RemoveAll(blevePath)
	if err != nil {
		t.Fatal(err)
	}

	err = os.Remove(blevePath)
	if err != nil && !os.IsNotExist(err) {
		t.Fatal(err)
	}
}

func cleanManifest(t *testing.T, tPaths TestPaths) {
	manifestPath := filepath.Join(tPaths.indexPath, "manifest.json")
	err := os.Remove(manifestPath)
	if err != nil && !os.IsNotExist(err) {
		t.Fatal(err)
	}
}
