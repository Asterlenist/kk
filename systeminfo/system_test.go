package systeminfo

import (
	"errors"
	"reflect"
	"testing"
)

func TestParseOSRelease(t *testing.T) {
	got := parseOSRelease([]byte(`
# comment
ID=ubuntu
ID_LIKE="debian"
PRETTY_NAME="Ubuntu 24.04 LTS"
VERSION_ID='24.04'
`))

	want := map[string]string{
		"ID":          "ubuntu",
		"ID_LIKE":     "debian",
		"PRETTY_NAME": "Ubuntu 24.04 LTS",
		"VERSION_ID":  "24.04",
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("parseOSRelease() = %#v, want %#v", got, want)
	}
}

func TestDetectPackageManagerPrefersDistributionManager(t *testing.T) {
	available := map[string]bool{"apt-get": true, "dnf": true}
	got := detectPackageManager("linux", []string{"ubuntu", "debian"}, func(name string) (string, error) {
		if available[name] {
			return "/usr/bin/" + name, nil
		}
		return "", errors.New("not found")
	})
	if got != "apt-get" {
		t.Fatalf("detectPackageManager() = %q, want apt-get", got)
	}
}

func TestDetectPackageManagerReturnsEmptyWhenUnavailable(t *testing.T) {
	got := detectPackageManager("linux", []string{"unknown"}, func(string) (string, error) {
		return "", errors.New("not found")
	})
	if got != "" {
		t.Fatalf("detectPackageManager() = %q, want empty", got)
	}
}

func TestDetectPackageManagerIgnoresUnsupportedManagers(t *testing.T) {
	available := map[string]bool{
		"pacman": true,
		"brew":   true,
		"winget": true,
		"pkg":    true,
	}
	for _, goos := range []string{"linux", "darwin", "windows", "freebsd"} {
		got := detectPackageManager(goos, []string{"arch"}, func(name string) (string, error) {
			if available[name] {
				return "/usr/bin/" + name, nil
			}
			return "", errors.New("not found")
		})
		if got != "" {
			t.Errorf("detectPackageManager(%q) = %q, want empty", goos, got)
		}
	}
}

func TestDetectPackageManagerDoesNotUseFallbackOnUnknownSystem(t *testing.T) {
	got := detectPackageManager("linux", []string{"unknown"}, func(name string) (string, error) {
		if name == "apt-get" {
			return "/usr/bin/apt-get", nil
		}
		return "", errors.New("not found")
	})
	if got != "" {
		t.Fatalf("detectPackageManager() = %q, want empty", got)
	}
}
