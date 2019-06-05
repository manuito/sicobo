package application

import (
	"log"
)

var debug = State.Config.LogLevel == "DEBUG"
var info = State.Config.LogLevel == "INFO"

// Debug : simple debug entrypoint
func Debug(v ...interface{}) {
	if debug {
		log.Println("[DEBUG] ", v)
	}
}

// Info : simple debug entrypoint
func Info(v ...interface{}) {
	if debug || info {
		log.Println("[INFO] ", v)
	}
}

// Error : some failure (always display)
func Error(v ...interface{}) {
	log.Fatal("[ERROR] ", v)
}
