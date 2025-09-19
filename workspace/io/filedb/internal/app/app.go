package app

import (
	"filedb/internal/app/config"
	"filedb/internal/app/service"
	"log"
)

func App() {
	conf := config.LoadConfigurations()
	databaseService := service.NewDatabaseService()

	database := databaseService.CreateDatabase(conf.DatabaseConfig.Name)
	log.Printf("Database: %+v\n", database)

}
