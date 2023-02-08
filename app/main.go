package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"syscall"
)

// Usage: your_docker.sh run <image> <command> <arg1> <arg2> ...
func main() {
	command := os.Args[3]
	args := os.Args[4:len(os.Args)]

	// Create temporary directory to chroot in
	chrootDir, err := ioutil.TempDir("", "")
	if err != nil {
		panic(err)
	}
	defer os.RemoveAll(chrootDir)

	// Copy binary to run (no symlink possible)
	targetPath := filepath.Join(chrootDir, command)
	err = os.MkdirAll(filepath.Dir(targetPath), os.ModePerm)
	if err != nil {
		panic(err)
	}
	bytesRead, err := ioutil.ReadFile(command)
	if err != nil {
		panic(err)
	}
	err = ioutil.WriteFile(targetPath, bytesRead, os.ModePerm)
	if err != nil {
		panic(err)
	}

	// Run chroot
	err = syscall.Chroot(chrootDir)
	if err != nil {
		panic(err)
	}

	// Run the command
	cmd := exec.Command(command, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	err = cmd.Run()
	if err != nil {
		fmt.Printf("Err: %v", err)
		if ee, ok := err.(*exec.ExitError); ok {
			os.Exit(ee.ExitCode())
		} else {
			os.Exit(1)
		}
	}
}
