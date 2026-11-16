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
		switch action {
		case "install":
			args = []string{"install", "-y"}
		case "remove":
			args = []string{"remove", "-y"}
		case "purge":
			args = []string{"purge", "-y"}
		case "update": // -up
			args = []string{"update"}
		case "upgrade": // -UP
			args = []string{"upgrade", "-y"}
		default:
			return "", nil, fmt.Errorf("unsupported package action %q for %s", action, manager)
		}
	case "dnf", "yum":
		switch action {
		case "install":
			args = []string{"install", "-y"}
		case "remove", "purge": // dnf не имеет purge, используем remove
			args = []string{"remove", "-y"}
		case "update": // -up
			// Явно refresh, чтобы обновить кэш, если нужно
			args = []string{"refresh"}
		case "upgrade": // -UP
			args = []string{"update", "-y"}
		default:
			return "", nil, fmt.Errorf("unsupported package action %q for %s", action, manager)
		}
	case "pacman":
		switch action {
		case "install":
			args = []string{"-S", "--needed", "--noconfirm"}
		case "remove":
			args = []string{"-R", "--noconfirm"}
		case "purge":
			args = []string{"-Rns", "--noconfirm"}
		case "update": // -up
			args = []string{"-Sy"}
		case "upgrade": // -UP
			args = []string{"-Syu", "--noconfirm"}
		default:
			return "", nil, fmt.Errorf("unsupported package action %q for pacman", action)
		}
	default:
		return "", nil, fmt.Errorf("unsupported package manager %q", manager)
	}

	// Для apt update/refresh не нужны имена пакетов
	if action == "update" && (manager == "apt-get" || manager == "apt" || manager == "dnf" || manager == "yum") {
		return program, args, nil
	}

	return program, append(args, packages...), nil
}
