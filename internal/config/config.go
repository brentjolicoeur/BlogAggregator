package config

import (
	"encoding/json"
	"os"
)

type Config struct {
	DBURL           string `json:"db_url"`
	CurrentUserName string `json:"current_user_name"`
}

func Read() (Config, error) {
	path, err := getConfigFilePath()
	if err != nil {
		return Config{}, err
	}

	content, err := os.ReadFile(path)
	if err != nil {
		return Config{}, err
	}

	var config Config
	if err := json.Unmarshal(content, &config); err != nil {
		return Config{}, err
	}
	return config, nil
}

func write(cfg Config) error {
	path, err := getConfigFilePath()
	if err != nil {
		return err
	}
	data, err := json.Marshal(cfg)
	if err != nil {
		return err
	}
	err = os.WriteFile(path, data, 0644)

	return err
}

func (c *Config) SetUser(name string) error {
	c.CurrentUserName = name

	err := write(*c)

	return err
}
