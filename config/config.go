package config

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

type Config struct {
	MyAnimeList MyAnimeListConfig `yaml:"myanimelist"`
}

type MyAnimeListConfig struct {
	TokenType    string `yaml:"token_type"`
	AccessToken  string `yaml:"access_token"`
	RefreshToken string `yaml:"refresh_token"`
	ExpiresIn    int    `yaml:"expires_in"`
}

var config Config

func LoadConfig() error {
	configDir, err := os.UserConfigDir()
	if err != nil {
		return fmt.Errorf("failed to get user config dir: %w", err)
	}

	configPath := filepath.Join(configDir, AppName, "config.yaml")
	data, err := os.ReadFile(configPath)
	if err != nil {
		if os.IsNotExist(err) {
			config = Config{}
			return nil
		}
		return fmt.Errorf("failed to read config file: %w", err)
	}

	if err := yaml.Unmarshal(data, &config); err != nil {
		return fmt.Errorf("failed to unmarshal config file: %w", err)
	}

	return nil
}

func SaveConfig() error {
	configDir, err := os.UserConfigDir()
	if err != nil {
		return fmt.Errorf("failed to get user config dir: %w", err)
	}

	appConfigDir := filepath.Join(configDir, AppName)
	if err := os.MkdirAll(appConfigDir, 0755); err != nil {
		return fmt.Errorf("failed to create config dir: %w", err)
	}

	configPath := filepath.Join(appConfigDir, "config.yaml")
	data, err := yaml.Marshal(&config)
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	if err := os.WriteFile(configPath, data, 0644); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}

	return nil
}

func GetConfig() *Config {
	return &config
}
