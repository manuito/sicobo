package application

import (
	"encoding/json"
	"fmt"
	"os"
)

/*
 * Pilote app configuration management.
 * Uses file "conf.json" for configuration definition
 */

// Configuration : Robot configuration holder
type Configuration struct {
	BingAPIkey       string `json:"bingAPIkey"`
	GoogleBookAPIKey string `json:"googleBookAPIKey"`
	LogLevel         string `json:"logLevel"`
	MongoDb          string `json:"mongoDb"`
}

func loadConfiguration() Configuration {
	file, _ := os.Open("conf.json")
	defer file.Close()
	decoder := json.NewDecoder(file)
	configuration := Configuration{}
	err := decoder.Decode(&configuration)
	if err != nil {
		fmt.Println("error at conf loading:", err)
	}
	return configuration
}
