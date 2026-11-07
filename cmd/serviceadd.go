package cmd

import (
	"fmt"
	"runtime"
	"strings"
)

func ServiceAdd(args []string) error {
	if runtime.GOOS != "linux" {
		return fmt.Errorf("service-add is supported only on Linux")
	}
	if err := validateNames(args, "usage: kk service-add <service> [service...]"); err != nil {
		return err
	}
	for _, service := range args {
		if strings.HasPrefix(service, "-") {
			return fmt.Errorf("invalid service name %q", service)
		}
	}
	commandArgs := append([]string{"enable", "--now"}, args...)
	return runCommand("systemctl", commandArgs...)
}
