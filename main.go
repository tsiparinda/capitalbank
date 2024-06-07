package main

import (
	"capitalbank/csv"
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

	var records []csv.CSVRecord
	var allfiles []csv.CSVfiles
	var delfiles bool = true
	records, allfiles, err = csv.LoadCSVfiles(records, allfiles)
	if err != nil {
		fmt.Println("Error loading csv files: ", err)
		delfiles = false
	}

	logic.StartExchangeTran(records)
	//logic.StartExchangePayments()
	logic.StartUpdateBalance()
	if delfiles {
		csv.DeleteCSVfiles(allfiles)
	}

}
