package cmd

import (
	"fmt"
	"os"
	"path/filepath"
)

func Update() error {
	if os.Getuid() != 0 {
		return fmt.Errorf("this command must be run as root: use sudo or switch to the root user")
	}

	repoDir, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("resolve working directory: %w", err)
	}
	if _, err := os.Stat(filepath.Join(repoDir, ".git")); err == nil {
		if err := runCommand("git", "pull", "--ff-only"); err != nil {
			return fmt.Errorf("update repository: %w", err)
		}
	} else {
		tempDir, err := os.MkdirTemp("", "kk-update-")
		if err != nil {
			return fmt.Errorf("create temporary directory: %w", err)
		}
		defer os.RemoveAll(tempDir)

		repoTarget := filepath.Join(tempDir, "kk")
		if err := runCommand("git", "clone", "https://github.com/Asterlenist/kk.git", repoTarget); err != nil {
			return fmt.Errorf("clone project repository: %w", err)
		}
		repoDir = repoTarget
	}

	binaryPath := filepath.Join(repoDir, "kk")
	if err := runCommand("go", "build", "-o", binaryPath, repoDir); err != nil {
		return fmt.Errorf("build updated binary: %w", err)
	}
	if err := runCommand("install", "-m", "0755", binaryPath, "/usr/local/bin/kk"); err != nil {
		return fmt.Errorf("install binary: %w", err)
	}
	return nil
}
