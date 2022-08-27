package main

import (
	"fmt"
	"os"
)

func main() {
	err := run()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run() error {
	err := ls()
	if err != nil {
		return err
	}
	return nil
}

func ls() error {
	dir, err := os.Getwd()
	if err != nil {
		return err
	}

	fis, err := os.ReadDir(dir)
	if err != nil {
		return err
	}

	for _, fi := range fis {
		fmt.Println(fi.Name())
	}

	return nil
}
