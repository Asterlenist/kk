package main

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

func TestRunLogsCommandAndCreatesConfig(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)
	previousArgs := os.Args
	os.Args = []string{"kk", "pkg-inst"}
	defer func() { os.Args = previousArgs }()

	if got := run(); got != 1 {
		t.Fatalf("run() = %d, want 1", got)
	}
	logPath := filepath.Join(home, ".local", "share", "kk", "history.log")
	data, err := os.ReadFile(logPath)
	if err != nil {
		t.Fatalf("read history: %v", err)
	}
	var got historyEntry
	if err := json.Unmarshal(data[:len(data)-1], &got); err != nil {
		t.Fatalf("decode history entry: %v", err)
	}
	if got.Command != "pkg-inst" || got.Status != "error" {
		t.Fatalf("history entry = %+v, want failed pkg-inst invocation", got)
	}
	if _, err := os.Stat(filepath.Join(home, ".config", "kk", "config.toml")); err != nil {
		t.Fatalf("default config was not created: %v", err)
	}
}

func TestRunDoesNotLogWhenDisabled(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)
	configPath := filepath.Join(home, ".config", "kk", "config.toml")
	if err := os.MkdirAll(filepath.Dir(configPath), 0700); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(configPath, []byte("[logging]\nenabled = false\n"), 0600); err != nil {
		t.Fatal(err)
	}
	previousArgs := os.Args
	os.Args = []string{"kk", "system"}
	defer func() { os.Args = previousArgs }()

	if got := run(); got != 0 {
		t.Fatalf("run() = %d, want 0", got)
	}
	logPath := filepath.Join(home, ".local", "share", "kk", "history.log")
	if _, err := os.Stat(logPath); !os.IsNotExist(err) {
		t.Fatalf("history file exists or stat failed: %v", err)
	}
}

type historyEntry struct {
	Command string `json:"command"`
	Status  string `json:"status"`
}
