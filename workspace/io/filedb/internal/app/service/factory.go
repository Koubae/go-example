package service

import "filedb/internal/app/config"

func NewDatabaseService() *DatabaseService {
	storageDirPath := config.GetStorageDirectory()

	return &DatabaseService{
		storageDirPath: storageDirPath,
	}
}
