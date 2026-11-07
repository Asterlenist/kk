package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/BurntSushi/toml"
)

const defaultConfig = `[logging]
enabled = true
path = "~/.local/share/kk/history.log"
max_size = 0
`

const minLogSize = 128

type Config struct {
	Logging LoggingConfig `toml:"logging"`
}

type LoggingConfig struct {
	Enabled bool   `toml:"enabled"`
	Path    string `toml:"path"`
	MaxSize int64  `toml:"max_size"`
}

func Load() (Config, error) {
	configPath, err := ConfigPath()
	if err != nil {
		return Config{}, err
	}
	if err := createDefault(configPath); err != nil {
		return Config{}, err
	}

	cfg := Config{
		Logging: LoggingConfig{
			Enabled: true,
			Path:    "~/.local/share/kk/history.log",
		},
	}
	metadata, err := toml.DecodeFile(configPath, &cfg)
	if err != nil {
		return Config{}, fmt.Errorf("read config %q: %w", configPath, err)
	}
	if keys := metadata.Undecoded(); len(keys) > 0 {
		return Config{}, fmt.Errorf("unknown config key %q", keys[0])
	}
	if cfg.Logging.Path == "" {
		return Config{}, fmt.Errorf("logging.path must not be empty")
	}
	if cfg.Logging.MaxSize < 0 || (cfg.Logging.MaxSize > 0 && cfg.Logging.MaxSize < minLogSize) {
		return Config{}, fmt.Errorf("logging.max_size must be 0 (unlimited) or at least %d bytes", minLogSize)
	}
	cfg.Logging.Path, err = expandHome(cfg.Logging.Path)
	if err != nil {
		return Config{}, err
	}
	return cfg, nil
}

func ConfigPath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("find home directory: %w", err)
	}
	return filepath.Join(home, ".config", "kk", "config.toml"), nil
}

func createDefault(path string) error {
	if _, err := os.Stat(path); err == nil {
		return nil
	} else if !os.IsNotExist(err) {
		return fmt.Errorf("check config %q: %w", path, err)
	}

	if err := os.MkdirAll(filepath.Dir(path), 0700); err != nil {
		return fmt.Errorf("create config directory: %w", err)
	}
	file, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0600)
	if err != nil {
		if os.IsExist(err) {
			return nil
		}
		return fmt.Errorf("create config %q: %w", path, err)
	}
	if _, err := file.WriteString(defaultConfig); err != nil {
		file.Close()
		return fmt.Errorf("write default config %q: %w", path, err)
	}
	if err := file.Close(); err != nil {
		return fmt.Errorf("close default config %q: %w", path, err)
	}
	return nil
}

func expandHome(path string) (string, error) {
	if path != "~" && !filepath.IsAbs(path) && len(path) >= 2 && path[:2] == "~/" {
		home, err := os.UserHomeDir()
		if err != nil {
			return "", fmt.Errorf("expand logging.path: %w", err)
		}
		return filepath.Join(home, filepath.FromSlash(path[2:])), nil
	}
	if path == "~" {
		return os.UserHomeDir()
	}
	if len(path) > 1 && path[:2] == "~/" {
		return "", fmt.Errorf("invalid logging.path %q", path)
	}
	return path, nil
}
