package service

import "filedb/internal/app/config"

func NewDatabaseService() *DatabaseService {
	storageDirPath := config.GetStorageDirectory()

	return &DatabaseService{
		storageDirPath: storageDirPath,
	}
}

func NewTableService(databasePath string, table string) *TableService {
	return &TableService{
		databasePath: databasePath,
		table:        table,
	}
}
