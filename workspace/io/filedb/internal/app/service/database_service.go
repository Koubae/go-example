package service

import (
	"filedb/internal/app/model"
	"filedb/pkg/serialization"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"
)

type DatabaseService struct {
	storageDirPath string
}

func (s *DatabaseService) CreateDatabase(name string) model.Database {
	exists, err := s.DatabaseExists(name)
	if err != nil {
		log.Fatalf("Error checking if database %s exists, error %v\n", name, err)
	}

	if exists {
		log.Printf("Database %s already exists\n", name)
	} else {
		err := s._createDatabaseAt(name)
		if err != nil {
			log.Fatalf("Error creating database %s, error %v\n", name, err)
		}
	}

	manifest := s.loadManifestFile(name, err)
	database := model.Database{
		Manifest: manifest,
	}
	return database

}

func (s *DatabaseService) DatabaseExists(name string) (bool, error) {
	databasePath := s.databasePath(name)
	info, err := os.Stat(databasePath)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}

		log.Printf("Error while getting info of path %s, error %v\n", databasePath, err)
		return false, DatabaseReadPathError
	}

	if info.IsDir() {
		return true, nil
	}

	return false, InvalidDatabaseFormatAtPath
}

func (s *DatabaseService) _createDatabaseAt(name string) error {
	databasePath := s.databasePath(name)
	log.Printf("Creating database %s (path=%s)\n", name, databasePath)

	err := os.MkdirAll(databasePath, 0755)
	if err != nil {
		return DatabaseWritePathError
	}
	log.Printf("Database %s created\n", databasePath)

	err = s._createManifest(name, err)
	if err != nil {
		return err
	}
	return nil

}

func (s *DatabaseService) _createManifest(name string, err error) error {
	manifestPath := s.manifestPath(name)
	file, err := os.Create(manifestPath)
	if err != nil {
		return DatabaseManifestCreateError
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			log.Printf("Error closing manifest file %s, error %v\n", manifestPath, err)
		}
	}(file)

	now := time.Now().UTC()
	manifestContent := model.Manifest{
		Name:    name,
		Created: now,
		Updated: now,
	}

	jsonContent, err := serialization.JSONSerializePretty[model.Manifest](manifestContent)
	if err != nil {
		return DatabaseJSONManifestSerializationError
	}

	_, err = file.Write(jsonContent)
	if err != nil {
		return DatabaseManifestWriteError
	}

	log.Printf("Manifest file for database %s at path %s created\n", name, manifestPath)
	return nil
}

func (s *DatabaseService) loadManifestFile(name string, err error) model.Manifest {
	manifestPath := s.manifestPath(name)
	data, err := os.ReadFile(manifestPath)
	if err != nil {
		log.Fatalf("Error reading manifest file %s, error %v\n", manifestPath, err)
	}
	log.Printf("Manifest file %s content: %s\n", manifestPath, string(data))

	var manifest model.Manifest
	err = serialization.JSONDeserialize(data, &manifest)
	if err != nil {
		log.Fatalf("Error unmarshalling manifest file %s, error %v\n", manifestPath, err)
	}

	fmt.Printf("Loaded manifest: %+v\n", manifest)
	return manifest
}

func (s *DatabaseService) databasePath(name string) string {
	return filepath.Join(s.storageDirPath, name)
}

func (s *DatabaseService) manifestPath(name string) string {
	return filepath.Join(s.databasePath(name), "manifest.json")
}
