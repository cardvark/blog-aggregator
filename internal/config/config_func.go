package config

import (
	"encoding/json"
	"fmt"
	"os"
)

const configFileName = ".gatorconfig.json"

var homePath string

func InitHomePath(mainWd string) {
	if homePath == "" {
		homePath = mainWd
	}
}

func Read() (Config, error) {
	var configData Config

	configPath := getConfigPath()

	fmt.Println(configPath)

	fileContent, err := os.ReadFile(configPath)

	if err != nil {
		return configData, fmt.Errorf("Error reading file content: %v\n", err)
	}

	err = json.Unmarshal(fileContent, &configData)
	if err != nil {
		return configData, fmt.Errorf("Error unmarshaling file content: %v\n", err)
	}

	return configData, nil

}

func getConfigPath() string {
	configPath := homePath + "/" + configFileName
	return configPath
}

// func SetUser(homePath string)
