package main

import (
	"capitalbank/config"
	"capitalbank/db"
	"capitalbank/logger"
	"capitalbank/logic"
	"fmt"
	"os"

	"github.com/sirupsen/logrus"
)

func main() {
	err := db.DB.Ping()
	if err != nil {
		fmt.Println("Database connection is not active")
		os.Exit(1)
	}

	// trace, debug, info, warn, error, fatal, panic
	loglevel := config.Config["logLevel"].(string)
	level, err := logrus.ParseLevel(loglevel)
	if err != nil {
		fmt.Printf("Error parsing level: %v\n", err)
		return
	}
	logger.Log.SetLevel(level)

	fields := make(map[string]interface{})
	fields["logLevel"] = logger.Log.GetLevel()
	// Add more fields dynamically...
	fields["location"] = "Auroville"
	logger.Log.WithFields(fields).Info("Program was started")

	logic.StartExchange()
}
