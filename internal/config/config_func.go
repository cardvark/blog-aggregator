package config

import (
	"encoding/json"
	"fmt"
	"os"
)

func Read(homePath string) (Config, error) {
	configPath := homePath + "/.gatorconfig.json"
	fmt.Println(configPath)

	var configData Config

	fileContent, err := os.ReadFile(configPath)

	if err != nil {
		return configData, fmt.Errorf("Error reading file content: %v", err)
	}

	err = json.Unmarshal(fileContent, &configData)
	if err != nil {
		return configData, fmt.Errorf("Error unmarshaling file content: %v", err)
	}

	return configData, nil

}
