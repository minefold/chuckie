package main

import (
	"archive/zip"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func restoreDir(source string, dest string) error {
	var cmd *exec.Cmd

	// source is:
	// 	  "508b5c0ab75f04080000007b/data.tar.lzo"
	// or "s3://bucket/508b5c0ab75f04080000007b/data.tar.lzo"

	if !strings.Contains(source, "s3://") {
		source = "s3://minefold-production/worlds/" + source
	}

	restoreDirBin, _ := filepath.Abs("bin/restore-dir")
	cmd = exec.Command(restoreDirBin, source)
	cmd.Dir = dest

	return cmd.Run()
}

func zipPath(zf io.Writer, root string) error {
	zipW := zip.NewWriter(zf)

	err := filepath.Walk(root,
		func(path string, info os.FileInfo, err error) error {
			if !info.IsDir() {
				fileName := path[len(root)+1:]
				zipF, err := zipW.Create(fileName)
				if err != nil {
					return err
				}

				f, err := os.Open(path)
				if err != nil {
					return err
				}
				_, err = io.Copy(zipF, f)
				return err
			}
			return nil
		})

	if err != nil {
		zipW.Close()
		return err
	}

	return zipW.Close()
}