package cmd

import (
	"fmt"
	"kk/systeminfo"
	"os/exec"
)

func Doctor() {
	fmt.Print("Checking the system...\n")

	// systemd check
	if _, err := exec.LookPath("systemctl"); err != nil {
		fmt.Println("[WARN] systemD is not found!")
	} else {
		fmt.Println("[OK] systemD is found")
	}

	// pkg manager check
	info, err := systeminfo.Detect()
	if err != nil || info.PackageManager == "" {
		fmt.Println("[FAIL] Package Manager is not found!")
	} else {
		fmt.Printf("[OK] Package Manager: %s\n", info.PackageManager)
	}

	// sudo check
	if _, err := exec.LookPath("sudo"); err != nil {
		fmt.Println("[FAIL] Sudo is not found!")
	} else {
		fmt.Println("[OK] Sudo is found")
	}

	// internet check
	if err := exec.Command("ping", "-c", "1", "archlinux.org").Run(); err != nil {
		fmt.Print("[WARN] No internet (or ping is blocked)")
	} else {
		fmt.Println("[OK] Internet is working")
	}
}
