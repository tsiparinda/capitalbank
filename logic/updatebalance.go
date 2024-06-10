package logic

import (
	"capitalbank/config"
	"capitalbank/logger"
	"capitalbank/utils"
	"time"

	// "capitalbank/logger"
	"capitalbank/pbapi"
	"capitalbank/store"

	"github.com/sirupsen/logrus"
)

func StartUpdateBalance() {
	// get accounts
	acc := []store.Account{}
	err := store.LoadAccounts(&acc)
	if err != nil {
		logger.Log.Info("StartUpdateBalance: Error from LoadAccounts:", err.Error())
		return
	}

	var reqdays int
	// take requestdays from config file and convert to int
	if value, ok := config.Config["RequestDaysBalance"].(float64); ok {
		// The value is a float64, handle it accordingly
		reqdays = int(value)
	} else {
		logger.Log.Info("StartUpdateBalance: Error loading reqdays from config:", err.Error())
		reqdays = 1
	}

	dateTo, _ := utils.GetShortDate(time.Now())
	dateFrom := dateTo.AddDate(0, 0, -reqdays)
	// fill the account struct by information about balance's state by date
	err = store.LoadCheckBalance(dateFrom, &acc)
	if err != nil {
		logger.Log.WithFields(logrus.Fields{
			"dateFrom": dateFrom,
			"acc":      acc,
		}).Info("StartUpdateBalance: Error from LoadCheckBalance:", err.Error())
		return
	}
	// cycle by accounts
	for _, a := range acc {
		// calc first maximum early date with bad saldo and set it as datefrom
		earliestBADdate, err := EarliestFBADTrue(&a.BalanceState)
		if err != nil {
			logger.Log.WithFields(logrus.Fields{
				"acc": a,
			}).Info("StartUpdateBalance: Error from EarliestFBADTrue:", err.Error())
			return
		}

		switch a.Bank {
		case "privat":
			if a.Token.Valid {
				privat := pbapi.PrivatBankAPI{
					UserAgent:   "Додаток API",
					Token:       a.Token.String,
					ContentType: "application/json;charset=utf8",
					Account:     a.Account,
					BankRegistr: a.BankRegistr,
					Direction:   a.Direction,
				}
				// get balance by dates from bank's server
				bal, _ := privat.GetBalance(earliestBADdate)

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
}
