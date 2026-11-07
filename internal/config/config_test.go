package config

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestLoadCreatesDefaultConfig(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)

	cfg, err := Load()
	if err != nil {
		t.Fatal(err)
	}
	wantPath := filepath.Join(home, ".local", "share", "kk", "history.log")
	if !cfg.Logging.Enabled || cfg.Logging.Path != wantPath || cfg.Logging.MaxSize != 0 {
		t.Fatalf("unexpected defaults: %+v", cfg.Logging)
	}
	if _, err := os.Stat(filepath.Join(home, ".config", "kk", "config.toml")); err != nil {
		t.Fatalf("default config was not created: %v", err)
	}
}

func TestLoadReadsLoggingSettings(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)
	configPath := filepath.Join(home, ".config", "kk", "config.toml")
	if err := os.MkdirAll(filepath.Dir(configPath), 0700); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(configPath, []byte("[logging]\nenabled = false\npath = \"~/custom.log\"\nmax_size = 1024\n"), 0600); err != nil {
		t.Fatal(err)
	}

	cfg, err := Load()
	if err != nil {
		t.Fatal(err)
	}
	if cfg.Logging.Enabled || cfg.Logging.Path != filepath.Join(home, "custom.log") || cfg.Logging.MaxSize != 1024 {
		t.Fatalf("unexpected settings: %+v", cfg.Logging)
	}
}

func TestLoadRejectsInvalidLimit(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)
	configPath := filepath.Join(home, ".config", "kk", "config.toml")
	if err := os.MkdirAll(filepath.Dir(configPath), 0700); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(configPath, []byte("[logging]\nmax_size = 1\n"), 0600); err != nil {
		t.Fatal(err)
	}

	_, err := Load()
	if err == nil || !strings.Contains(err.Error(), "logging.max_size") {
		t.Fatalf("expected invalid limit error, got %v", err)
	}
}
