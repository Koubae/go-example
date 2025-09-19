package app

import (
	"encoding/json"
	"filedb/internal/app/config"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"
)

type Manifest struct {
	Name    string    `json:"name"`
	Created time.Time `json:"created"` // time.RFC3339Nano
	Updated time.Time `json:"updated"`
}

func App() {
	conf := config.LoadConfigurations()

	storageDirectory := config.GetStorageDirectory()

	databasePath := filepath.Join(storageDirectory, conf.DatabaseConfig.Name)
	info, err := os.Stat(databasePath)
	createDatabase := false
	if err != nil {
		if os.IsNotExist(err) {
			createDatabase = true
			log.Printf("Path %s does not exist\n", databasePath)
		} else {
			log.Fatalf("Error while getting info of path %s, error %v\n", databasePath, err)

		}
	} else {
		if info.IsDir() {
			log.Printf("Path %s is a directory\n", databasePath)
		} else {
			log.Printf("Path %s is a file\n", databasePath)
		}
	}

	if createDatabase {
		log.Printf("Creating database %s\n", databasePath)
		err := os.MkdirAll(databasePath, 0755)
		if err != nil {
			log.Fatalf("Error creating database %s, error %v\n", databasePath, err)
		} else {
			log.Printf("Database %s created\n", databasePath)
		}
	}

	manifestPath := filepath.Join(databasePath, "manifest.json")
	_, err = os.Stat(manifestPath)
	createManifest := false
	if err != nil {
		if !os.IsNotExist(err) {
			log.Fatalf("Error while getting info of path %s, error %v\n", manifestPath, err)
		} else {
			createManifest = true
			log.Printf("Manifest file %s does not exist\n", manifestPath)
		}
	} else {
		log.Printf("Manifest file %s exists\n", manifestPath)
	}

	if createManifest {
		file, err := os.Create(manifestPath)
		if err != nil {
			log.Fatalf("Error creating manifest file %s, error %v\n", manifestPath, err)
		}
		defer func(file *os.File) {
			err := file.Close()
			if err != nil {
				log.Fatalf("Error closing manifest file %s, error %v\n", manifestPath, err)
			}
		}(file)

		// created := time.Now().UTC().Format(time.RFC3339Nano)
		created := time.Now().UTC()
		manifestContent := Manifest{
			Name:    conf.DatabaseConfig.Name,
			Created: created,
			Updated: created,
		}

		jsonContent, err := json.MarshalIndent(manifestContent, "", "	")
		if err != nil {
			log.Fatalf("Error marshalling manifest content, error %v\n", err)
		}
		_, err = file.WriteString(string(jsonContent))
		if err != nil {
			log.Fatalf("Error writing manifest content to file %s, error %v\n", manifestPath, err)
		}
		log.Printf("Manifest file %s created\n", manifestPath)
	}

	data, err := os.ReadFile(manifestPath)
	if err != nil {
		log.Fatalf("Error reading manifest file %s, error %v\n", manifestPath, err)
	}
	log.Printf("Manifest file %s content: %s\n", manifestPath, string(data))

	var manifest Manifest
	err = json.Unmarshal(data, &manifest)
	if err != nil {
		log.Fatalf("Error unmarshalling manifest file %s, error %v\n", manifestPath, err)
	}

	fmt.Printf("Loaded manifest: %+v\n", manifest)
}
