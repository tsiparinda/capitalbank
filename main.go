package main

import (
	"capitalbank/db"
	"capitalbank/logger"
	"capitalbank/logic"
	"fmt"
	"os"
)

func main() {
	err := db.DB.Ping()
	if err != nil {
		fmt.Println("Database connection is not active")
		os.Exit(1)
	}

	fields := make(map[string]interface{})
	fields["logLevel"] = logger.Log.GetLevel()
	// Add more fields dynamically...
	fields["location"] = "Earth"
	logger.Log.WithFields(fields).Info("Program was started")

	logic.StartExchangeTran()
	logic.StartUpdateBalance()
}
