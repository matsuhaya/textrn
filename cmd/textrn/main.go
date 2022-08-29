package main

import (
	"bufio"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
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

	usedOldName := make(map[string]int)
	for i, fn := range fileNames {
		fmt.Fprintln(tempFile, fn)
		usedOldName[fn] = i
	}

	command += " " + tempFile.Name()
	err = openEditor(command)
	if err != nil {
		return err
	}

	newFileNames, usedNewName, err := scanTempFile(tempFile.Name())
	if err != nil {
		return err
	}

	fileNames, newFileNames, err = replaceUsedFileNameToUniq(&fileNames, &newFileNames, usedOldName, usedNewName)
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

func scanTempFile(tempFileName string) ([]string, map[string]bool, error) {
	// os.Open used because file seek reset doesn't read tempFile by 'vim'.
	tempFile, err := os.Open(tempFileName)
	if err != nil {
		return nil, nil, err
	}
	defer tempFile.Close()

	var newFiles []string
	usedNewName := make(map[string]bool)
	scanner := bufio.NewScanner(tempFile)
	for scanner.Scan() {
		newFileName := scanner.Text()
		if usedNewName[newFileName] {
			return nil, nil, errors.New("duplicate file name specified")
		}
		usedNewName[newFileName] = true
		newFiles = append(newFiles, newFileName)
	}
	return newFiles, usedNewName, nil
}

// if newFileName was used old filename, replace uniq temp file name.
// after that, add new rename list "temp -> newFileName".
// ex.)
// before:
// a						->	b
// b						->	c
//
// after:
// a						->	temp-123456
// b						->	c
// temp-123456	->	b
func replaceUsedFileNameToUniq(oldFileNames, newFileNames *[]string, usedOldName map[string]int, usedNewName map[string]bool) ([]string, []string, error) {
	for newIndex, newFileName := range *newFileNames {
		if oldIndex, ok := usedOldName[newFileName]; ok && newIndex < oldIndex {
			var tempFineName string
			isUniq := false
			for i := 0; i < 1000; i++ {
				tempFineName = genTempFileName("temp-")
				if _, ok := usedOldName[tempFineName]; !ok && !usedNewName[tempFineName] {
					isUniq = true
				}
				if isUniq {
					break
				}
			}
			if !isUniq {
				return nil, nil, errors.New("failed to generate uniq file name")
			}

			(*newFileNames)[newIndex] = tempFineName
			*oldFileNames = append(*oldFileNames, tempFineName)
			*newFileNames = append(*newFileNames, newFileName)
		}
	}
	return *oldFileNames, *newFileNames, nil
}

func genTempFileName(prefix string) string {
	randBytes := make([]byte, 16)
	rand.Read(randBytes)
	return prefix + hex.EncodeToString(randBytes)
}

func renameFiles(oldFiles, newFiles []string) error {
	if len(oldFiles) != len(newFiles) {
		return errors.New("the number of new and old files must match")
	}

	for i, f := range oldFiles {
		err := os.Rename(filepath.Join(rootdir, f), filepath.Join(rootdir, newFiles[i]))
		if err != nil {
			return err
		}
	}
	return nil
}
