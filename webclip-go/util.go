package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
)

// CopyFile copies a single file
func CopyFile(src, dst string) error {
	sourceFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer sourceFile.Close()

	destFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer destFile.Close()

	_, err = io.Copy(destFile, sourceFile)
	return err
}

// CopyDir copies a directory recursively
func CopyDir(src string, dst string) error {
	srcInfo, err := os.Stat(src)
	if err != nil {
		return err
	}

	if err := os.MkdirAll(dst, srcInfo.Mode()); err != nil {
		return err
	}

	entries, err := ioutil.ReadDir(src)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		sourcePath := filepath.Join(src, entry.Name())
		destPath := filepath.Join(dst, entry.Name())

		if entry.IsDir() {
			if err := CopyDir(sourcePath, destPath); err != nil {
				return err
			}
		} else {
			if err := CopyFile(sourcePath, destPath); err != nil {
				return err
			}
		}
	}
	return nil
}

// Clear clears the terminal
func Clear() {
	var cmd *exec.Cmd
	if runtime.GOOS == "windows" {
		cmd = exec.Command("cls")
	} else {
		cmd = exec.Command("clear")
	}
	cmd.Stdout = os.Stdout
	cmd.Run()
}

// Printh (Helper for printing error headers, though Ruby used printe)
func Printh(message string) {
	fmt.Println("")
	fmt.Println("=== ERROR ===")
	fmt.Println(message)
	fmt.Println("==============")
	fmt.Println("")
}
