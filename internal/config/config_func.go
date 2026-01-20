package config

import (
	"encoding/json"
	"fmt"
	"os"
)

const configFileName = ".gatorconfig.json"

var homePath string
var configPath string

func InitPaths(mainWd string) {
	if homePath == "" {
		// homePath = mainWd
		newPath, err := os.UserHomeDir()
		if err != nil {
			fmt.Println(err)
		}

		homePath = newPath
	}

	configPath = homePath + "/" + configFileName
}

func Read() (Config, error) {
	var configData Config

	// fmt.Println(configPath)

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

// func getConfigPath() string {
// 	configPath := homePath + "/" + configFileName
// 	return configPath
// }

func write(cfg Config) error {
	// configPath := getConfigPath()

	jsonData, err := json.Marshal(cfg)
	if err != nil {
		return fmt.Errorf("Error encoding Config to json: %v", err)
	}

	err = os.WriteFile(configPath, jsonData, 0644)
	if err != nil {
		return fmt.Errorf("Error writing to file: %v", err)
	}

	return nil

}

func (c Config) SetUser(username string) error {
	c.Current_user_name = username
	err := write(c)
	return err
}
