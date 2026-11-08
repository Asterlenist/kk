package cmd

import (
	"fmt"
	"os"

	"kk/systeminfo"
)

func PkgUpdate(args []string) error {
	if len(args) > 0 {
		return fmt.Errorf("usage: kk pkg-upd")
	}
	if os.Getuid() != 0 {
		return fmt.Errorf("this command must be run as root: use sudo or switch to the root user")
	}

	info, err := systeminfo.Detect()
	if err != nil {
		return fmt.Errorf("detect system: %w", err)
	}
	if info.PackageManager == "" {
		return fmt.Errorf("no supported package manager found")
	}
	program, commandArgs, err := packageCommand(info.PackageManager, "update", nil)
	if err != nil {
		return err
	}
	return runCommand(program, commandArgs...)
}
