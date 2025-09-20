package service

import (
	"encoding/json"
	"errors"
	"filedb/internal/app/model"
	"filedb/pkg/serialization"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"
)

type TableService struct {
	databasePath string
	table        string
}

func (s *TableService) CreateTableIfNotExists() (*model.Table[any], error) {
	return s.CreateTable(false)
}

func (s *TableService) CreateTable(raiseIfExists bool) (*model.Table[any], error) {
	exists, err := s.TableExists()
	if err != nil {
		return nil, err
	}

	if exists {
		log.Printf("Table %s already exists\n", s.table)
		if raiseIfExists {
			return nil, TableDuplicate
		}
	} else {
		err := s._createTable()
		if err != nil {
			return nil, err
		}
	}

	table, err := s.loadTableManifest()
	if err != nil {
		return nil, err
	}
	return table, nil

}

func (s *TableService) TableExists() (bool, error) {
	tablePath := s.tablePath()

	_, err := os.Stat(tablePath)

	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return false, nil
		}
		log.Printf("Error while getting info of path %s, error %v\n", tablePath, err)
		return false, TableReadPathError
	}

	return true, nil
}

func (s *TableService) TableFileName() string {
	return fmt.Sprintf("table__%s.json", s.table)
}

func (s *TableService) tablePath() string {
	return filepath.Join(s.databasePath, "tables", s.TableFileName())
}

func (s *TableService) _createTable() error {
	tablePath := s.tablePath()
	log.Printf("Creating table %s (path=%s)\n", s.table, tablePath)

	file, closer, err := serialization.CreateFileWithSafeCloseDeferrable(tablePath)
	if err != nil {
		log.Printf("ERROR: %v\n", err)
		return TableFileCreateError
	}
	defer closer()

	now := time.Now().UTC()
	manifest := model.Manifest{
		Name:    s.table,
		Created: now,
		Updated: now,
	}

	table := model.Table[any]{
		Manifest: manifest,
		Records:  make(map[string]model.Record[any]),
	}
	jsonContent, err := serialization.JSONSerializePretty[model.Table[any]](table)
	if err != nil {
		return TableJSONSerializationError
	}

	if _, err = file.Write(jsonContent); err != nil {
		return TableJSONWriteError
	}

	log.Printf("JSON file for table %s at path %s created\n", s.table, tablePath)
	return nil
}

func (s *TableService) loadTableManifest() (*model.Table[any], error) {
	tablePath := s.tablePath()
	log.Printf("Creating table %s (path=%s)\n", s.table, tablePath)

	file, closer, err := serialization.OpenFileWithSafeCloseDeferrable(tablePath)
	if err != nil {
		log.Printf("ERROR: %v\n", err)
		return nil, err
	}
	defer closer()

	// ---------------------------
	// Read root value and validate is a JSON object {} (It could be an Array [])
	// ---------------------------
	decoder := json.NewDecoder(file)
	t, err := decoder.Token()
	if err != nil {
		return nil, errors.New(fmt.Sprintf("invalid JSON error while reading token, error: %v", err))
	}
	if d, ok := t.(json.Delim); !ok || d != '{' {
		return nil, errors.New(fmt.Sprintf("invalid JSON, root must be object, expected open brace { got '%s'", t))
	}

	// ---------------------------
	// Read up to manifest key and the stop
	// ---------------------------
	var manifest model.Manifest
	for decoder.More() {
		t, err := decoder.Token()
		if err != nil {
			return nil, errors.New(fmt.Sprintf("invalid JSON error while reading token, error: %v", err))
		}

		key := t.(string)
		if key == "manifest" {
			if err := decoder.Decode(&manifest); err != nil {
				return nil, errors.New(fmt.Sprintf("invalid JSON error while decoding manifest, error: %v", err))
			}
			break
		}
	}

	table := &model.Table[any]{
		Manifest: manifest,
		Records:  nil,
	}
	return table, nil
}
