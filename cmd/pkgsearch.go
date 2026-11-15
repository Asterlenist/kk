package cmd

import (
	"fmt"
	"kk/systeminfo"
	"os"
	"os/exec"
)

func PkgSearch(args []string) {
	if len(args) == 0 {
		fmt.Println("Usage: kk pkg-search <query>")
		return
	}
	info, _ := systeminfo.Detect()
	var cmd *exec.Cmd
	switch info.PackageManager {
	case "pacman":
		cmd = exec.Command("pacman", "-Ss", args[0])
	case "apt":
		cmd = exec.Command("apt", "search", args[0])
	case "dnf":
		cmd = exec.Command("dnf", "search", args[0])
	}
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Run()
}
