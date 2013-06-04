package main

import (
	"os/exec"
	"path/filepath"
)

// source is http url
func restoreDir(source string, dest string) error {
	var cmd *exec.Cmd

	restoreDirBin, _ := filepath.Abs("bin/restore-dir")
	cmd = exec.Command(restoreDirBin, source)
	cmd.Dir = dest

	return cmd.Run()
}
