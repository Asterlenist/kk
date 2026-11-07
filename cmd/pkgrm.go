package cmd

import (
	"fmt"

	"kk/systeminfo"
)

func PkgRemove(args []string) error {
	if err := validateNames(args, "usage: kk pkg-rem <package> [package...]"); err != nil {
		return err
	}
	info, err := systeminfo.Detect()
	if err != nil {
		return fmt.Errorf("detect system: %w", err)
	}
	if info.PackageManager == "" {
		return fmt.Errorf("no supported package manager found")
	}
	program, commandArgs, err := packageCommand(info.PackageManager, "remove", args)
	if err != nil {
		return err
	}
	return runCommand(program, commandArgs...)
}
