package cmd

import (
	"log"
	"os"
	"os/exec"
)

func Update() error {
	cmd := exec.Command("curl", "-fsSL",
		"https://raw.githubusercontent.com/Asterlenist/kk/main/install.py")

	pipe := exec.Command("sudo", "python3")
	pipe.Stdin, _ = cmd.StdoutPipe()
	pipe.Stdout = os.Stdout
	pipe.Stderr = os.Stderr

	if err := pipe.Start(); err != nil {
		log.Fatal(err)
	}
	if err := cmd.Start(); err != nil {
		log.Fatal(err)
	}
	cmd.Wait()
	pipe.Wait()
	return nil
}
