package main

import (
	"io"
	"log"
	"os/exec"
)

func getInkscapeVersion(binPath string, helpArgs ...string) (string, error) {
	return runCommandWithArgumentsGetStdOut(binPath, helpArgs...)
}

func runCommandWithArgumentsGetStdOut(binPath string, args ...string) (string, error) {
	log.Println("runCommandWithArgumentsGetStdOut", binPath, args)
	out, err := exec.Command(binPath, args...).Output()
	return string(out), err
}

func runCommandReadStdin(r io.Reader, binPath string, args ...string) (string, error) {
	cmd := exec.Command(binPath, args...)
	stdin, err := cmd.StdinPipe()
	if err != nil {
		return "", err
	}
	var funcError error
	go func() {
		defer func() {
			if funcError = stdin.Close(); err != nil {
				return
			}
		}()
		_, funcError = io.Copy(stdin, r)
		return
	}()

	// If there is a problem with running inkscape...
	stdoutStderr, err := cmd.CombinedOutput()
	if err != nil {
		log.Printf("%s\n", stdoutStderr)
		return string(stdoutStderr), err
	}
	if funcError != nil {
		return "", funcError
	}
	return "", nil
}
