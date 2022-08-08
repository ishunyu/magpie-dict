package main

import (
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

type TestPaths struct {
	rootPath   string
	dataPath   string
	tmpPath    string
	indexPath  string
	backupPath string
}

func TestBadPaths(t *testing.T) {
	tPaths := newTestPaths(t, "bad_paths")
	_, err := NewIndex(tPaths.dataPath, tPaths.indexPath)
	if err == nil {
		t.Fatalf("Expected error.")
	}
	t.Log(err)
}

func TestWithoutCacheWithoutManifest(t *testing.T) {
	testPath := "without_cache_without_manifest"
	tPaths := newTestPaths(t, testPath)
	clean(t, tPaths)
	defer clean(t, tPaths)

	doTest(t, tPaths)
}

func TestWithoutCacheWithManifest(t *testing.T) {
	testPath := "without_cache_with_manifest"
	tPaths := newTestPaths(t, testPath)
	cleanCache(t, tPaths)
	defer cleanCache(t, tPaths)

	doTest(t, tPaths)
}

func TestWithCacheWithoutManifest(t *testing.T) {
	tPaths := newTestPaths(t, "with_cache_without_manifest")
	backup(t, tPaths)
	defer restore(t, tPaths)
	cleanManifest(t, tPaths)
	defer cleanManifest(t, tPaths)

	doTest(t, tPaths)
}

func TestWithCacheWithManifest(t *testing.T) {
	testPath := "with_cache_with_manifest"
	tPaths := newTestPaths(t, testPath)
	backup(t, tPaths)
	defer restore(t, tPaths)

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

func newTestPaths(t *testing.T, testPathSub string) TestPaths {
	testPath := filepath.Join("testdata", "index", testPathSub)
	testPathAbs, err := filepath.Abs(testPath)
	if err != nil {
		t.Fatal(err)
	}

	return TestPaths{
		testPathAbs,
		filepath.Join(testPathAbs, "data"),
		filepath.Join(testPathAbs, "tmp"),
		filepath.Join(testPathAbs, "tmp", "index"),
		filepath.Join(testPathAbs, "bak"),
	}
}

func backup(t *testing.T, tPaths TestPaths) {
	err := os.Mkdir(tPaths.backupPath, fs.ModeDir|0744)
	if err != nil {
		t.Fatal(err)
	}

	output, err := execShell("cp", "-r", tPaths.tmpPath, tPaths.backupPath)
	if err != nil {
		t.Fatal(err, output)
	}
}

func restore(t *testing.T, tPaths TestPaths) {
	clean(t, tPaths)

	output, err := execShell("mv", tPaths.backupPath+"/*", tPaths.rootPath)
	if err != nil {
		t.Fatal(err, string(output))
	}

	err = removeDir(tPaths.backupPath)
	if err != nil {
		t.Fatal(err)
	}
}

func clean(t *testing.T, tPaths TestPaths) {
	cleanCache(t, tPaths)
	cleanManifest(t, tPaths)

	err := removeDir(tPaths.tmpPath)
	if err != nil {
		t.Fatal(err)
	}
}

func cleanCache(t *testing.T, tPaths TestPaths) {
	blevePath := filepath.Join(tPaths.indexPath, "bleve")

	err := removeDir(blevePath)
	if err != nil {
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

func removeDir(path string) error {
	err := os.RemoveAll(path)
	if err != nil {
		return err
	}

	err = os.Remove(path)
	if err != nil && !os.IsNotExist(err) {
		return err
	}

	return nil
}

func execShell(cmd ...string) ([]byte, error) {
	cmdStr := strings.Join(cmd, " ")
	return exec.Command("/bin/bash", "-c", cmdStr).CombinedOutput()
}
