package main

import (
	"fmt"
	"kk/cmd"
	"os"
)

func main() {
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
		os.Exit(2)
	}

	command := os.Args[1]
	args := os.Args[2:]
	var err error

	switch command {
	case "pkg-inst":
		err = cmd.PkgInstall(args)
	case "pkg-rem":
		err = cmd.PkgRemove(args)
	case "pkg-pur":
		err = cmd.PkgPurge(args)
	case "flatpak-inst":
		err = cmd.FlatpakInstall(args)
	case "service-add":
		err = cmd.ServiceAdd(args)
	case "service-rem":
		err = cmd.ServiceRemove(args)
	case "system":
		err = cmd.SystemInfo()
	default:
		fmt.Fprintln(os.Stderr, "Unknown command:", command)
		os.Exit(2)
	}
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
