package app

import (
	"filedb/internal/app/config"
	"filedb/internal/app/service"
	"fmt"
	"log"
)

func App() {
	conf := config.LoadConfigurations()
	databaseService := service.NewDatabaseService()

	database, err := databaseService.CreateDatabaseIfNotExists(conf.DatabaseConfig.Name)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Database: %+v\n", database)

	tableName := "users" // todo make more tables
	tableService := service.NewTableService(database.Path(), tableName)

	fmt.Println(tableService)

}
