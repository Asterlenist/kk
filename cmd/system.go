package cmd

import (
	"fmt"

	"kk/systeminfo"
)

func SystemInfo() error {
	info, err := systeminfo.Detect()
	if err != nil {
		return err
	}

	fmt.Println("Operating system:", info.OS)
	if info.Distribution != "" {
		fmt.Println("Distribution:", info.Distribution)
	}
	if info.Version != "" {
		fmt.Println("Version:", info.Version)
	}
	if info.PackageManager == "" {
		fmt.Println("Package manager: not found")
	} else {
		fmt.Println("Package manager:", info.PackageManager)
	}
	return nil
}
