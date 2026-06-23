package config

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/cockroachdb/errors"
)

type DBConfig struct {
}

var (
	ErrConfigNotFound    = errors.New("configuration file not found")
	ErrConfigWriteFailed = errors.New("failed to write configuration file")
)

// ConfigManager is a struct that manages configuration settings for the application.
// It can read the ~/.dbterm/config.yaml file and provide methods to access and modify configuration settings.
type ConfigManager struct {
	path   string   // Path to the configuration file
	config DBConfig // Configuration settings
}

// NewConfigManager creates a new ConfigManager instance with the specified path to the configuration file.
func NewConfigManager(path *string) *ConfigManager {

	var configPath string

	if path != nil {
		configPath = *path
	} else {
		configPath = "~/.dbterm/config.yaml"
	}

	return &ConfigManager{
		path:   configPath,
		config: DBConfig{},
	}
}

func (config *ConfigManager) GetPath() string {
	return config.path
}

func (config *ConfigManager) LoadConfig() error {
	if config.path == "" {
		return ErrConfigNotFound
	}

	resolvedPath, err := resolveHomePath(config.path)
	if err != nil {
		return ErrConfigWriteFailed
	}
	config.path = resolvedPath

	if _, err := os.Stat(config.path); err != nil {
		if !os.IsNotExist(err) {
			return ErrConfigWriteFailed
		}

		// create the directory if it does not exist
		if err := os.MkdirAll(filepath.Dir(config.path), 0755); err != nil {
			return ErrConfigWriteFailed
		}

		// create the config file
		file, err := os.OpenFile(config.path, os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0600)
		if err != nil {
			return ErrConfigWriteFailed
		}
		if err := file.Close(); err != nil {
			return ErrConfigWriteFailed
		}
	}

	// read yaml

	return nil
}

func resolveHomePath(path string) (string, error) {
	if path != "~" && !strings.HasPrefix(path, "~/") {
		return path, nil
	}

	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	if path == "~" {
		return home, nil
	}
	return filepath.Join(home, strings.TrimPrefix(path, "~/")), nil
}
