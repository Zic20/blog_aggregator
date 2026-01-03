package config

import (
	"encoding/json"
	"fmt"
	"os"
)

type Config struct {
	DBUrl           string `json:"db_url"`
	CurrentUserName string `json:"current_user_name"`
}

func (cfg *Config) SetUser(user string) error {
	cfg.CurrentUserName = user
	if err := write(*cfg); err != nil {
		return fmt.Errorf("Error setting user: %s", err)
	}
	return nil
}

const (
	configFileName = "/.gatorconfig.json"
)

func Read() (Config, error) {
	var configData Config
	configPath, err := getConfigFilePath()
	if err != nil {
		return configData, err
	}
	content, err := os.ReadFile(configPath)
	if err != nil {
		return configData, fmt.Errorf("Error reading config file: %s", err)
	}

	err = json.Unmarshal(content, &configData)
	if err != nil {
		return configData, fmt.Errorf("Error marshalling data: %s", err)
	}

	return configData, nil
}

func getConfigFilePath() (string, error) {
	configDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("Error loading home directory: %s", err)
	}

	return configDir + configFileName, nil
}
func write(cfg Config) error {
	data, err := json.Marshal(&cfg)
	if err != nil {
		return err
	}

	configPath, err := getConfigFilePath()
	if err != nil {
		return err
	}

	if err = os.WriteFile(configPath, data, 0644); err != nil {
		return err
	}
	return nil
}
