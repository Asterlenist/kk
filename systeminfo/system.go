// Package systeminfo detects the operating system, Linux distribution, and
// an available native package manager.
package systeminfo

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
)

// Info describes the current operating system and its detected package manager.
type Info struct {
	OS             string `json:"os"`
	Distribution   string `json:"distribution,omitempty"`
	Version        string `json:"version,omitempty"`
	PackageManager string `json:"package_manager,omitempty"`
}

// Detect identifies the current operating system, distribution, and the first
// installed package manager appropriate for that system.
func Detect() (Info, error) {
	info := Info{OS: operatingSystemName(runtime.GOOS)}

	var release map[string]string
	if runtime.GOOS == "linux" {
		data, err := readOSRelease()
		if err != nil {
			return info, err
		}
		release = parseOSRelease(data)
		info.Distribution = release["PRETTY_NAME"]
		if info.Distribution == "" {
			info.Distribution = release["ID"]
		}
		info.Version = release["VERSION_ID"]
	}

	ids := strings.Fields(strings.ToLower(release["ID"] + " " + release["ID_LIKE"]))
	info.PackageManager = detectPackageManager(runtime.GOOS, ids, exec.LookPath)
	return info, nil
}

func operatingSystemName(goos string) string {
	switch goos {
	case "linux":
		return "Linux"
	case "darwin":
		return "macOS"
	case "windows":
		return "Windows"
	case "freebsd":
		return "FreeBSD"
	default:
		return goos
	}
}

func readOSRelease() ([]byte, error) {
	for _, path := range []string{"/etc/os-release", "/usr/lib/os-release"} {
		data, err := os.ReadFile(path)
		if err == nil {
			return data, nil
		}
		if !os.IsNotExist(err) {
			return nil, fmt.Errorf("read %s: %w", path, err)
		}
	}
	return nil, fmt.Errorf("could not find /etc/os-release or /usr/lib/os-release")
}

func parseOSRelease(data []byte) map[string]string {
	values := make(map[string]string)
	scanner := bufio.NewScanner(strings.NewReader(string(data)))
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		key, value, ok := strings.Cut(line, "=")
		if !ok {
			continue
		}
		key = strings.TrimSpace(key)
		if key == "" {
			continue
		}
		values[key] = unquoteOSReleaseValue(strings.TrimSpace(value))
	}
	return values
}

func unquoteOSReleaseValue(value string) string {
	if len(value) < 2 {
		return value
	}
	if value[0] == '"' && value[len(value)-1] == '"' {
		if unquoted, err := strconv.Unquote(value); err == nil {
			return unquoted
		}
	}
	if value[0] == '\'' && value[len(value)-1] == '\'' {
		return value[1 : len(value)-1]
	}
	return value
}

func detectPackageManager(goos string, ids []string, lookPath func(string) (string, error)) string {
	if goos != "linux" {
		return ""
	}

	preferences := make([]string, 0)
	for _, id := range ids {
		switch id {
		case "debian", "ubuntu", "linuxmint", "pop":
			preferences = append(preferences, "apt-get", "apt")
		case "fedora", "rhel", "centos", "rocky", "almalinux", "ol":
			preferences = append(preferences, "dnf", "yum")
		case "arch", "manjaro", "endeavouros", "artix", "garuda":
			preferences = append(preferences, "pacman")
		}
	}

	seen := make(map[string]bool, len(preferences))
	for _, name := range preferences {
		if seen[name] {
			continue
		}
		seen[name] = true
		if _, err := lookPath(name); err == nil {
			return name
		}
	}
	return ""
}
