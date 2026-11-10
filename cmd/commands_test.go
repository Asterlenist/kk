package cmd

import (
	"reflect"
	"testing"
)

func TestPackageCommand(t *testing.T) {
	tests := []struct {
		manager string
		action  string
		program string
		args    []string
	}{
		{"apt", "install", "apt", []string{"install", "-y", "curl"}},
		{"apt-get", "purge", "apt-get", []string{"purge", "-y", "curl"}},
		{"dnf", "install", "dnf", []string{"install", "-y", "curl"}},
		{"yum", "purge", "yum", []string{"remove", "-y", "curl"}},
		{"pacman", "install", "pacman", []string{"-S", "--needed", "--noconfirm", "curl"}},
		{"pacman", "remove", "pacman", []string{"-R", "--noconfirm", "curl"}},
		{"pacman", "purge", "pacman", []string{"-Rns", "--noconfirm", "curl"}},
		{"pacman", "update", "pacman", []string{"-Syu", "--noconfirm"}},
		{"pacman", "upgrade", "pacman", []string{"-Syu", "--noconfirm"}},
	}

	for _, test := range tests {
		packages := []string{"curl"}
		if test.action == "update" || test.action == "upgrade" {
			packages = nil
		}
		program, args, err := packageCommand(test.manager, test.action, packages)
		if err != nil {
			t.Fatalf("packageCommand(%q, %q): %v", test.manager, test.action, err)
		}
		if program != test.program || !reflect.DeepEqual(args, test.args) {
			t.Errorf("packageCommand(%q, %q) = (%q, %#v), want (%q, %#v)", test.manager, test.action, program, args, test.program, test.args)
		}
	}

	for _, manager := range []string{"zypper", "apk", "emerge", "xbps-install", "brew", "winget", "choco", "scoop", "pkg"} {
		if _, _, err := packageCommand(manager, "install", []string{"curl"}); err == nil {
			t.Errorf("packageCommand(%q) unexpectedly succeeded", manager)
		}
	}
	if _, _, err := packageCommand("pacman", "unknown", nil); err == nil {
		t.Error("packageCommand(\"pacman\", \"unknown\") unexpectedly succeeded")
	}
}
