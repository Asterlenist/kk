package cmd

import (
	"fmt"
	"runtime"
	"strings"
)

func ServiceRemove(args []string) error {
	if runtime.GOOS != "linux" {
		return fmt.Errorf("service-rem is supported only on Linux")
	}
	if err := validateNames(args, "usage: kk service-rem <service> [service...]"); err != nil {
		return err
	}
	for _, service := range args {
		if strings.HasPrefix(service, "-") {
			return fmt.Errorf("invalid service name %q", service)
		}
	}
	commandArgs := append([]string{"disable", "--now"}, args...)
	return runCommand("systemctl", commandArgs...)
}
