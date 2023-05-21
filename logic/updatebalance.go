package logic

import (
	"capitalbank/api"
	"capitalbank/config"
	"capitalbank/logger"
	"capitalbank/utils"
	"fmt"
	"time"

	// "capitalbank/logger"
	"capitalbank/pbapi"
	"capitalbank/store"

	"github.com/sirupsen/logrus"
)

func StartUpdateBalance() error {
	// get accounts
	acc := []store.Account{}
	err := store.LoadAccounts(&acc)
	if err != nil {
		return err
	}

	var reqdays int
	// take requestdays from config file and convert to int
	if value, ok := config.Config["RequestDaysBalance"].(float64); ok {
		// The value is a float64, handle it accordingly
		reqdays = int(value)
	} else {
		logger.Log.WithFields(logrus.Fields{
			"reqdays": value,
		}).Infof("Error loading reqdays from config:", err.Error())
		reqdays = 1
	}
	//	reqdays = 30

	dateTo, _ := utils.GetShortDate(time.Now())
	dateFrom := dateTo.AddDate(0, 0, -reqdays)
	//fmt.Printf("dateFrom : ", dateTo, dateFrom, reqdays)
	//return nil
	// fill the account struct by information about balance's state by date
	err = store.LoadCheckBalance(dateFrom, &acc)
	if err != nil {
		return err
	}
	//fmt.Printf("acc: \n", acc)
	// cycle by accounts
	for _, a := range acc {
		// calc first maximum early date with bad saldo and set it as datefrom
		earliestBADdate, err := EarliestFBADTrue(&a.BalanceState)
		//fmt.Printf("earliestBADdate, a: %v, %v \n", earliestBADdate, a)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			return err
		}

		switch a.Bank {
		case "privat":
			if a.Token.Valid == true {
				var privat api.BankAPI
				privat = pbapi.PrivatBankAPI{
					UserAgent:   "Додаток API",
					Token:       a.Token.String,
					ContentType: "application/json;charset=utf8",
					Account:     a.Account,
					BankRegistr: a.BankRegistr,
					Direction:   a.Direction,
				}
				// get balance by dates from bank's server
				bal, _ := privat.GetBalance(earliestBADdate)
				// fmt.Printf("bal %v \n", bal)

				//save balance
				store.SaveBalance(bal)

				// trace info to logs
				// result, err := utils.StructToMap(bal)
				// if err != nil {
				// 	fmt.Printf(err.Error())
				// }
				// logger.Log.WithFields(result).Tracef("GetBalance: ")

				// if err == nil {
				// 	store.SaveTransactions(tran)
				// }
			}
		}
	}
	return nil
}
