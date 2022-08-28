package main

import (
	"fmt"
	"os"
	"os/exec"
)

var (
	command = "code -w" // "code" needs waiting for the files to be closed before returning.
	rootdir = "."
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
	files, err := listFiles()
	if err != nil {
		return err
	}
	if len(files) == 0 {
		fmt.Println("no files")
		return nil
	}

	tempFile, err := os.CreateTemp(tempdir, "textrn-")
	if err != nil {
		return err
	}
	defer os.Remove(tempFile.Name())

	for _, f := range files {
		fmt.Fprintln(tempFile, f)
	}

	command += " " + tempFile.Name()
	err = openEditor(command)
	if err != nil {
		return err
	}

	return nil
}

// except for directories
func listFiles() ([]string, error) {
	var fileNames []string
	files, _ := os.ReadDir(rootdir)
	for _, f := range files {
		if f.IsDir() {
			continue
		}
		fileNames = append(fileNames, f.Name())
	}
	return fileNames, nil
}

func openEditor(command string) error {
	cmd := exec.Command("sh", "-c", command)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
