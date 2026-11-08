package cmd

import (
	"fmt"

	"kk/systeminfo"
)

func PkgUpgrade(args []string) error {
	if len(args) > 0 {
		return fmt.Errorf("usage: kk pkg-upg")
	}

	info, err := systeminfo.Detect()
	if err != nil {
		return fmt.Errorf("detect system: %w", err)
	}
	if info.PackageManager == "" {
		return fmt.Errorf("no supported package manager found")
	}

	program, commandArgs, err := packageCommand(info.PackageManager, "upgrade", nil)
	if err != nil {
		return err
	}
	return runCommand(program, commandArgs...)
}
