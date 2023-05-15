package main

import (
	"fmt"
	"os"
	"capitalbank/config"
	"capitalbank/db"
	"capitalbank/logger"
	"capitalbank/logic"
)

func main() {
	err := db.DB.Ping()
	if err != nil {
		fmt.Println("Database connection is not active")
		os.Exit(1)
	}
	
	// trace, debug, info, warn, error, fatal, panic
	loglevel := config.Config["logLevel"].(string)
	logger.SetLogLevel(loglevel)

	fields := make(map[string]interface{})
	fields["logLevel"] = logger.Log.GetLevel()
	// Add more fields dynamically...
	fields["location"] = "Auroville"
	logger.Log.WithFields(fields).Info("Program was started")

	logic.GetParams()
}
