package main

import (
	"fmt"
	"os"
	"strings"

	"kk/cmd"
	"kk/internal/config"
	"kk/internal/history"
)

func main() {
	os.Exit(run())
}

func run() int {
	cfg, err := config.Load()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return 1
	}

	if len(os.Args) < 2 {
		fmt.Print(`Usage: kk <command> [args]
To print help for a specific command, run "kk <command> --help"
`)
		logCommand(cfg, "kk", fmt.Errorf("no command provided"))
		return 2
	}

	command := os.Args[1]
	args := os.Args[2:]
	var commandErr error
	exitCode := 0

	switch command {
	case "--help", "-h", "help":
		fmt.Print(`Usage: kk <command> [args]

		kk -i <package> [package...]  - installs packages
		kk -s <package> [package...]  - searches the packages
		kk -up                        - updates the package index
		kk -UP                        - upgrades the installed packages
		kk -r <package> [package...]  - removes packages
		kk -R <package> [package...]  - purges packages and configuration
		kk -fi <app-id> [app-id...]   - installs apps from Flathub
		kk -Sa <service> [service...] - enables and starts systemd services
		kk -Sr <service> [service...] - disables and stops systemd services
		kk system                     - displays the operating system and package manager
		kk update                     - updates kk to the latest version
		kk doctor                     - checks the system
`)
	case "-i":
		commandErr = cmd.PkgInstall(args)
	case "-s":
		cmd.PkgSearch(args)
	case "-up":
		commandErr = cmd.PkgUpdate(args)
	case "-UP":
		commandErr = cmd.PkgUpgrade(args)
	case "-r":
		commandErr = cmd.PkgRemove(args)
	case "-R":
		commandErr = cmd.PkgPurge(args)
	case "-fi":
		commandErr = cmd.FlatpakInstall(args)
	case "-Sa":
		commandErr = cmd.ServiceAdd(args)
	case "-Sr":
		commandErr = cmd.ServiceRemove(args)
	case "system":
		commandErr = cmd.SystemInfo()
	case "update":
		commandErr = cmd.Update()
	case "doctor":
		cmd.Doctor()
	default:
		fmt.Fprintln(os.Stderr, "Unknown command:", command)
		commandErr = fmt.Errorf("unknown command %q", command)
		exitCode = 2
	}

	logCommand(cfg, strings.Join(os.Args[1:], " "), commandErr)
	if commandErr != nil {
		if exitCode == 0 {
			fmt.Fprintln(os.Stderr, commandErr)
			exitCode = 1
		}
	}
	return exitCode
}

func logCommand(cfg config.Config, command string, commandErr error) {
	if !cfg.Logging.Enabled {
		return
	}
	if err := history.Append(cfg.Logging.Path, cfg.Logging.MaxSize, command, commandErr); err != nil {
		fmt.Fprintln(os.Stderr, "write command history:", err)
	}
}
