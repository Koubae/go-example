package service

import (
	"fmt"
	"path/filepath"
)

type TableService struct {
	databasePath string
	table        string
}

func (s TableService) TableFileName() string {
	return fmt.Sprintf("table__%s.json", s.table)
}

func (s TableService) tablePath() string {
	return filepath.Join(s.databasePath, "tables", s.TableFileName())
}
