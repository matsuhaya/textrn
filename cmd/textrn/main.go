package main

import (
	"bufio"
	"errors"
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
		fmt.Fprintf(os.Stdout, "No files")
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

	newFiles, err := scanTempFile(tempFile.Name())
	if err != nil {
		return err
	}

	err = renameFiles(files, newFiles)
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

func scanTempFile(tempFileName string) ([]string, error) {
	// os.Open used because file seek reset doesn't read tempFile by 'vim'.
	tempFile, err := os.Open(tempFileName)
	if err != nil {
		return nil, err
	}
	defer tempFile.Close()

	var newFiles []string
	usedName := make(map[string]bool)
	scanner := bufio.NewScanner(tempFile)
	for scanner.Scan() {
		newFileName := scanner.Text()
		if usedName[newFileName] {
			return nil, errors.New("duplicate file name specified")
		}
		usedName[newFileName] = true
		newFiles = append(newFiles, newFileName)
	}
	return newFiles, nil
}

func renameFiles(oldFiles, newFiles []string) error {
	if len(oldFiles) != len(newFiles) {
		return errors.New("the number of new and old files must match")
	}

	for i, f := range oldFiles {
		os.Rename(f, newFiles[i])
	}
	return nil
}
