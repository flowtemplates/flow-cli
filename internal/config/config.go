package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

// Config struct defining expected fields
type Config struct {
	TemplatesFolder string `json:"templatesFolder" yaml:"templatesFolder"`
}

// fileExists checks if a file exists
func fileExists(filename string) bool {
	_, err := os.Stat(filename)
	return err == nil
}

// findConfig searches for the config file in the current directory and its parents
func findConfig(baseName string) (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("failed to get pwd: %w", err)
	}

	extensions := []string{".json", ".yaml", ".yml"}

	for {
		for _, ext := range extensions {
			configPath := filepath.Join(dir, baseName+ext)
			if fileExists(configPath) {
				return configPath, nil
			}
		}

		parentDir := filepath.Dir(dir)
		if parentDir == dir {
			break
		}
		dir = parentDir
	}

	return "", fmt.Errorf("config file %s[.json/.yaml/.yml] not found", baseName)
}

// ReadConfigFile reads and parses a JSON or YAML file into the provided struct
func ReadConfigFile(filename string, v *Config) error {
	data, err := os.ReadFile(filename)
	if err != nil {
		return fmt.Errorf("failed to read file: %w", err)
	}

	// Detect file type and parse accordingly
	switch filepath.Ext(filename) {
	case ".json":
		err = json.Unmarshal(data, v)
	case ".yaml", ".yml":
		err = yaml.Unmarshal(data, v)
	default:
		return fmt.Errorf("unsupported file format: %s", filename)
	}

	if err != nil {
		return fmt.Errorf("failed to parse %s: %w", filename, err)
	}

	return nil
}

// GetConfig searches for a config file, reads, and parses it into a Config struct
func GetConfig(baseName string) (*Config, error) {
	configPath, err := findConfig(baseName)
	if err != nil {
		return nil, err
	}

	var config Config
	if err := ReadConfigFile(configPath, &config); err != nil {
		return nil, err
	}

	return &config, nil
}
