package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func validateNames(names []string, usage string) error {
	if len(names) == 0 {
		return fmt.Errorf("%s", usage)
	}
	for _, name := range names {
		if strings.TrimSpace(name) == "" || strings.HasPrefix(name, "-") {
			return fmt.Errorf("invalid name %q", name)
		}
	}
	return nil
}

func runCommand(program string, args ...string) error {
	command := exec.Command(program, args...)
	command.Stdin = os.Stdin
	command.Stdout = os.Stdout
	command.Stderr = os.Stderr
	if err := command.Run(); err != nil {
		return fmt.Errorf("%s: %w", program, err)
	}
	return nil
}

func packageCommand(manager, action string, packages []string) (string, []string, error) {
	program := manager
	var args []string

	switch manager {
	case "apt-get", "apt":
		args = []string{action, "-y"}
	case "dnf", "yum":
		operation := action
		if action == "purge" {
			operation = "remove"
		}
		args = []string{operation, "-y"}
	default:
		return "", nil, fmt.Errorf("unsupported package manager %q", manager)
	}

	return program, append(args, packages...), nil
}
