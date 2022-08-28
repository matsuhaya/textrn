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
	fileNames, err := listFiles()
	if err != nil {
		return err
	}
	if len(fileNames) == 0 {
		fmt.Fprintf(os.Stdout, "No files")
		return nil
	}

	tempFile, err := os.CreateTemp(tempdir, "textrn-")
	if err != nil {
		return err
	}
	defer os.Remove(tempFile.Name())

	usedOldName := make(map[string]bool)
	for _, fn := range fileNames {
		fmt.Fprintln(tempFile, fn)
		usedOldName[fn] = true
	}

	command += " " + tempFile.Name()
	err = openEditor(command)
	if err != nil {
		return err
	}

	newFileNames, err := scanTempFile(tempFile.Name(), usedOldName)
	if err != nil {
		return err
	}

	err = renameFiles(fileNames, newFileNames)
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

func scanTempFile(tempFileName string, usedOldName map[string]bool) ([]string, error) {
	// os.Open used because file seek reset doesn't read tempFile by 'vim'.
	tempFile, err := os.Open(tempFileName)
	if err != nil {
		return nil, err
	}
	defer tempFile.Close()

	var newFiles []string
	usedNewName := make(map[string]bool)
	scanner := bufio.NewScanner(tempFile)
	for scanner.Scan() {
		newFileName := scanner.Text()
		if usedNewName[newFileName] {
			return nil, errors.New("duplicate file name specified")
		} else if usedOldName[newFileName] {
			return nil, errors.New("can not rename to the current used file or directly name")
		}
		usedNewName[newFileName] = true
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
