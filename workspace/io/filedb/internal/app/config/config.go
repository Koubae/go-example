package config

import (
	"log"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

const configDirectoryName string = "config"
const configFilename string = "config.yml"

type DatabaseConfig struct {
	Name string `yaml:"name"`
}

type ConfigurationSchema struct {
	Version        string `yaml:"version"`
	Name           string `yaml:"name"`
	DatabaseConfig `yaml:"database"`
}

// config ConfigSchema
// Also as reference we could declare a 'struct' or config as so if we feel lazy
// var Config map[string]interface{} = make(map[string]interface{})
var config ConfigurationSchema

func LoadConfigurations() ConfigurationSchema {
	appDirectory := GetAppDirectory()
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
		log.Fatalf("Error getting current working directory: %s\n", err)
	}
	return appDirectory
}
