package main

import (
	"bufio"
	"fmt"
	"log"
	"os/exec"
	"strings"
)

const (
	BUILDKIT_ENV_MISSING       = "BUILDKIT_HOST environment variable is not set."
	BUILDKIT_UNABLE_TO_CONNECT = "ERRO failed to get buildkit information."
)

func build() error {
	cmd := exec.Command("railpack", "build", ".")
	cmd.Dir = "/home/baby/workflow/projects/friday"
	stdout, err := cmd.StdoutPipe()
	stderr, err := cmd.StderrPipe()

	if err != nil {
		log.Fatalf("Error creating StdoutPipe: %v", err)
	}

	err = cmd.Start()
	if err != nil {
		log.Fatalf("Error starting command: %v", err)
	}

	scanner := bufio.NewScanner(stdout)
	errScanner := bufio.NewScanner(stderr)
	for scanner.Scan() {
		_ = scanner.Text()
	}

	for errScanner.Scan() {
		line := errScanner.Text()
		if strings.Contains(line, BUILDKIT_ENV_MISSING) || strings.Contains(line, BUILDKIT_UNABLE_TO_CONNECT) {
			fmt.Println("Issue with build kit")
			cmd.Process.Kill()
		}
	}

	if scanner.Err() != nil {
		log.Printf("Error reading from stdout: %v", scanner.Err())
		if cmd.Process != nil {
			cmd.Process.Kill()
		}
	}

	if errScanner.Err() != nil {
		log.Printf("Error reading from stderr: %v", errScanner.Err())
		if cmd.Process != nil {
			cmd.Process.Kill()
		}
	}

	err = cmd.Wait()

	if err != nil {
		log.Fatalf("Command finished with error: %v", err)
	}
	return nil
}

func main() {
	build()
}
