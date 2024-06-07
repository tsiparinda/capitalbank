package logic

import (
	"capitalbank/api"
	"capitalbank/csv"
	"capitalbank/ibcsvapi"
	"capitalbank/logger"
	"capitalbank/pbapi"
	"fmt"

	// "capitalbank/logger"

	"capitalbank/store"

	"github.com/sirupsen/logrus"
)

func StartExchangeTran(csvrecords []csv.CSVRecord) {
	fmt.Println("StartExchangeTran", csvrecords)
	acc := []store.Account{}
	err := store.LoadAccounts(&acc)
	if err != nil {
		logger.Log.Error("StartExchangeTran: Error from LoadAccounts:", err.Error())
		return
	}
	//get and save transactions
	for _, a := range acc {
		switch a.Bank {
		case "privat":

			if a.Token.Valid {
				// var privat api.BankAPI
				privat := pbapi.PrivatBankAPI{
					UserAgent:   "Додаток API",
					Token:       a.Token.String,
					ContentType: "application/json;charset=utf8",
					Account:     a.Account,
					BankRegistr: a.BankRegistr,
					Direction:   a.Direction,
				}
				tran, err := privat.GetTransactions()
				if err == nil {
					store.SaveTransactions(tran)
				}
			}

		case "iBank2UAcsv":
			var iBankCSV api.BankAPI
			iBankCSV = ibcsvapi.IBankCSVAPI{
				Account:     a.Account,
				BankRegistr: a.BankRegistr,
				Direction:   a.Direction,
				Records:     csvrecords,
			}
			tran, err := iBankCSV.GetTransactions()
			if err == nil {
				store.SaveTransactions(tran)
			} else {
				logger.Log.WithFields(logrus.Fields{
					"Account": a.Account,
				}).Error("StartExchangeTran: Error from GetTransactions iBank2UAcsv:", err.Error())
				return
			}

		}
	}
}
