package config

import (
	"os"
	"path/filepath"
)

func getConfigFilePath() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	path := filepath.Join(homeDir, configFileName)

	return path, nil
}
