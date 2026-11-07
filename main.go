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
	    	kk pkg-inst <package> [package...] - installs packages using the native package manager
	    	kk pkg-rem <package> [package...]  - removes packages
	    	kk pkg-pur <package> [package...]  - purges packages and configuration where supported
	    	kk flatpak-inst <app-id>            - installs an app from Flathub
	    	kk service-add <service> [...]      - enables and starts systemd services
	    	kk service-rem <service> [...]      - disables and stops systemd services
    		kk system                    - displays the operating system and package manager
`)
		logCommand(cfg, "kk", fmt.Errorf("no command provided"))
		return 2
	}

	command := os.Args[1]
	args := os.Args[2:]
	var commandErr error
	exitCode := 0

	switch command {
	case "pkg-inst":
		commandErr = cmd.PkgInstall(args)
	case "pkg-rem":
		commandErr = cmd.PkgRemove(args)
	case "pkg-pur":
		commandErr = cmd.PkgPurge(args)
	case "flatpak-inst":
		commandErr = cmd.FlatpakInstall(args)
	case "service-add":
		commandErr = cmd.ServiceAdd(args)
	case "service-rem":
		commandErr = cmd.ServiceRemove(args)
	case "system":
		commandErr = cmd.SystemInfo()
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
