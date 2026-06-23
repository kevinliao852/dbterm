package config

import (
	"github.com/cockroachdb/errors"
	"os"
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

	if config.path == "~/.dbterm/config.yaml" {
		config.path = os.ExpandEnv(config.path)
	}

	if _, err := os.Stat(config.path); os.IsNotExist(err) {
		// create the directory if it does not exist
		if err := os.MkdirAll("~/.dbterm", 0755); err != nil {
			return ErrConfigWriteFailed
		}

		// create the config file
		file, err := os.Create(config.path)
		if err != nil {
			return ErrConfigWriteFailed
		}
		defer file.Close()

	}

	// read yaml

	return nil
}
