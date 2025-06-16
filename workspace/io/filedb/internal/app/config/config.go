package config

import (
	"gopkg.in/yaml.v3"
	"log"
	"os"
	"path/filepath"
)

const configDirectoryName string = "config"
const configFilename string = "config.yml"

type ConfigurationSchema struct {
	Version string `yaml:"version"`
	Name    string `yaml:"name"`
}

// config ConfigSchema
// Also as reference we could declare a 'struct' or config as so if we feel lazy
// var Config map[string]interface{} = make(map[string]interface{})
var config ConfigurationSchema

func LoadConfigurations() ConfigurationSchema {
	var appDirectory string = GetAppDirectory()
	configPath := filepath.Join(appDirectory, configDirectoryName, configFilename)
	configFile, err := os.ReadFile(configPath)
	if err != nil {
		log.Fatalf("Error reading configuration file at path %s: %s\n", configPath, err)
	}

	err = yaml.Unmarshal(configFile, &config)
	if err != nil {
		log.Fatalf("Error unmashalling yaml configuration file at path %s: %s\n", configPath, err)
	}

	log.Printf("Configuration file %s loaded in %s\n", configFilename, configPath)
	return config
}

func GetConfig() ConfigurationSchema {
	return config
}

func GetAppDirectory() string {
	appDirectory, err := os.Getwd()
	if err != nil {
		log.Fatal("Error getting current working directory: %s\n", err)
	}
	return appDirectory
}
