package logic

import (
	"capitalbank/csv"
	"capitalbank/pbapi"
	"fmt"

	// "capitalbank/logger"

	"capitalbank/store"
)

func StartExchangePayments(csvrecords []csv.CSVRecord) error {
	fmt.Println("StartExchangePayments", csvrecords)
	payments := []store.Payment{}
	// send privat
	err := store.LoadPaymentsPrivat(&payments)
	if err != nil {
		return err
	}
	for _, p := range payments {
		if p.Token.Valid {
			//var privat api.BankAPI
			privat := pbapi.PrivatBankAPI{
				UserAgent:   "Додаток API",
				Token:       p.Token.String,
				ContentType: "application/json;charset=utf8",
				Account:     p.PayerAccount,
			}
			_, err := privat.SendPayment()
			if err == nil {
				//	store.SaveTransactions(tran)
			}
		}
	}

	// send IBankCSV
	// payments := make([]store.Payment, 0)
	// err := store.LoadPaymentsIBankCSV(&payments)
	// if err != nil {
	// 	return err
	// }
	// var iBankCSV api.BankAPI
	// iBankCSV = ibcsvapi.IBankCSVAPI{
	// 	Account:     a.Account,
	// 	BankRegistr: a.BankRegistr,
	// 	Direction:   a.Direction,
	// 	Records:     csvrecords,
	// }
	// tran, err := iBankCSV.GetTransactions()
	// if err == nil {
	// 	store.SaveTransactions(tran)
	// } else {
	// 	fmt.Println("StartExchangeTran: Error from GetTransactions iBank2UAcsv: ", err)
	// }

	return nil
}
