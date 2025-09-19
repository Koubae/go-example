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

func (s *DatabaseService) CreateDatabaseIfNotExists(name string) (*model.Database, error) {
	return s.CreateDatabase(name, false)

}

func (s *DatabaseService) CreateDatabase(name string, raiseIfExists bool) (*model.Database, error) {
	exists, err := s.DatabaseExists(name)
	if err != nil {
		return nil, err
	}

	if exists {
		log.Printf("Database %s already exists\n", name)
		if raiseIfExists {
			return nil, DatabaseDuplicate
		}
	} else {
		err := s._createDatabaseAt(name)
		if err != nil {
			return nil, err
		}
	}

	manifest, err := s.loadManifestFile(name, err)
	if err != nil {
		return nil, err
	}

	database := model.NewDatabase(*manifest, s.databasePath(name))
	return database, nil
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

func (s *DatabaseService) loadManifestFile(name string, err error) (*model.Manifest, error) {
	manifestPath := s.manifestPath(name)
	data, err := os.ReadFile(manifestPath)
	if err != nil {
		log.Printf("Error reading manifest file %s, error %v\n", manifestPath, err)
		return nil, DatabaseManifestLoadError
	}
	log.Printf("Manifest file %s content: \n%s\n", manifestPath, string(data))

	var manifest model.Manifest
	err = serialization.JSONDeserialize(data, &manifest)
	if err != nil {
		log.Printf("Error unmarshalling manifest file %s, error %v\n", manifestPath, err)
		return nil, DatabaseManifestLoadError
	}

	fmt.Printf("Loaded manifest: %+v\n", manifest)
	return &manifest, nil
}

func (s *DatabaseService) databasePath(name string) string {
	return filepath.Join(s.storageDirPath, name)
}

func (s *DatabaseService) manifestPath(name string) string {
	return filepath.Join(s.databasePath(name), "manifest.json")
}
