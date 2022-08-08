package utils

import (
	"os"
	"os/exec"
	"strings"
)

func RemoveDir(path string) error {
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

func ExecShell(cmd ...string) ([]byte, error) {
	cmdStr := strings.Join(cmd, " ")
	return exec.Command("/bin/bash", "-c", cmdStr).CombinedOutput()
}
