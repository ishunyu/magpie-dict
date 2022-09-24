package testutils

import (
	"io/fs"
	"os"
	"path/filepath"
	"testing"

	"github.com/ishunyu/magpie-dict/pkg/internal/utils"
)

type TestPaths struct {
	RootPath   string
	DataPath   string
	TmpPath    string
	IndexPath  string
	BackupPath string
}

func Backup(t *testing.T, tPaths TestPaths) {
	err := os.Mkdir(tPaths.BackupPath, fs.ModeDir|0744)
	if err != nil {
		t.Fatal(err)
	}

	output, err := utils.ExecShell("cp", "-r", tPaths.TmpPath, tPaths.BackupPath)
	if err != nil {
		t.Fatal(err, output)
	}
}

func Restore(t *testing.T, tPaths TestPaths) {
	Clean(t, tPaths)

	output, err := utils.ExecShell("mv", tPaths.BackupPath+"/*", tPaths.RootPath)
	if err != nil {
		t.Fatal(err, string(output))
	}

	err = utils.RemoveDir(tPaths.BackupPath)
	if err != nil {
		t.Fatal(err)
	}
}

func Clean(t *testing.T, tPaths TestPaths) {
	CleanCache(t, tPaths)
	CleanManifest(t, tPaths)

	err := utils.RemoveDir(tPaths.TmpPath)
	if err != nil {
		t.Fatal(err)
	}
}

func CleanCache(t *testing.T, tPaths TestPaths) {
	blevePath := filepath.Join(tPaths.IndexPath, "bleve")

	err := utils.RemoveDir(blevePath)
	if err != nil {
		t.Fatal(err)
	}
}

func CleanManifest(t *testing.T, tPaths TestPaths) {
	manifestPath := filepath.Join(tPaths.IndexPath, "manifest.json")
	err := os.Remove(manifestPath)
	if err != nil && !os.IsNotExist(err) {
		t.Fatal(err)
	}
}
