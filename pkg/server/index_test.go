package main

import (
	"path/filepath"
	"testing"

	"github.com/ishunyu/magpie-dict/pkg/internal/testutils"
)

func TestBadPaths(t *testing.T) {
	tPaths := newTestPaths(t, "bad_paths")
	_, err := NewIndex(tPaths.DataPath, tPaths.IndexPath)
	if err == nil {
		t.Fatalf("Expected error.")
	}
	t.Log(err)
}

func TestWithoutCacheWithoutManifest(t *testing.T) {
	testPath := "without_cache_without_manifest"
	tPaths := newTestPaths(t, testPath)
	testutils.Clean(t, tPaths)
	defer testutils.Clean(t, tPaths)

	doTest(t, tPaths)
}

func TestWithoutCacheWithManifest(t *testing.T) {
	testPath := "without_cache_with_manifest"
	tPaths := newTestPaths(t, testPath)
	testutils.CleanCache(t, tPaths)
	defer testutils.CleanCache(t, tPaths)

	doTest(t, tPaths)
}

func TestWithCacheWithoutManifest(t *testing.T) {
	tPaths := newTestPaths(t, "with_cache_without_manifest")
	testutils.Backup(t, tPaths)
	defer testutils.Restore(t, tPaths)
	testutils.CleanManifest(t, tPaths)
	defer testutils.CleanManifest(t, tPaths)

	doTest(t, tPaths)
}

func TestWithCacheWithManifest(t *testing.T) {
	testPath := "with_cache_with_manifest"
	tPaths := newTestPaths(t, testPath)
	testutils.Backup(t, tPaths)
	defer testutils.Restore(t, tPaths)

	doTest(t, tPaths)
}

func doTest(t *testing.T, tPaths testutils.TestPaths) {
	index, err := NewIndex(tPaths.DataPath, tPaths.IndexPath)
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

func newTestPaths(t *testing.T, testPathSub string) testutils.TestPaths {
	testPath := filepath.Join("testdata", "index", testPathSub)
	testPathAbs, err := filepath.Abs(testPath)
	if err != nil {
		t.Fatal(err)
	}

	return testutils.TestPaths{
		RootPath:   testPathAbs,
		DataPath:   filepath.Join(testPathAbs, "data"),
		TmpPath:    filepath.Join(testPathAbs, "tmp"),
		IndexPath:  filepath.Join(testPathAbs, "tmp", "index"),
		BackupPath: filepath.Join(testPathAbs, "bak"),
	}
}
