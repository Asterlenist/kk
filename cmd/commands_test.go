package cmd

import (
	"reflect"
	"testing"
)

func TestPackageCommandSupportsOnlyAPTAndDNF(t *testing.T) {
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
	}

	for _, test := range tests {
		program, args, err := packageCommand(test.manager, test.action, []string{"curl"})
		if err != nil {
			t.Fatalf("packageCommand(%q, %q): %v", test.manager, test.action, err)
		}
		if program != test.program || !reflect.DeepEqual(args, test.args) {
			t.Errorf("packageCommand(%q, %q) = (%q, %#v), want (%q, %#v)", test.manager, test.action, program, args, test.program, test.args)
		}
	}

	for _, manager := range []string{"pacman", "zypper", "apk", "emerge", "xbps-install", "brew", "winget", "choco", "scoop", "pkg"} {
		if _, _, err := packageCommand(manager, "install", []string{"curl"}); err == nil {
			t.Errorf("packageCommand(%q) unexpectedly succeeded", manager)
		}
	}
}
