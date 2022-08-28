package main

import (
	"fmt"
	"io/fs"
	"os"
	"os/exec"
)

var (
	command = "code -w" // "code" needs waiting for the files to be closed before returning.
	tempdir = "."
)

func main() {
	err := run()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run() error {
	fis, err := ls()
	if err != nil {
		return err
	}

	tempFile, err := os.CreateTemp(tempdir, "textrn-")
	if err != nil {
		return err
	}
	defer os.Remove(tempFile.Name())
	for _, fi := range fis {
		_, err = tempFile.Write([]byte(fmt.Sprintln(fi.Name())))
		if err != nil {
			return err
		}
	}

	command += " " + tempFile.Name()
	err = openEditor(command)
	if err != nil {
		return err
	}

	return nil
}

func ls() ([]fs.DirEntry, error) {
	dir, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	fis, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}
	return fis, nil
}

func openEditor(command string) error {
	cmd := exec.Command("sh", "-c", command)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
