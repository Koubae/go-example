package app

import (
	configLoader "filedb/internal/app/config"
	"fmt"
)

func App() {
	var config = configLoader.LoadConfigurations()
	fmt.Printf("configFile:%v\n", config)
}
