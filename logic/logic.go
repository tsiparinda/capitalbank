package logic

import (
	"capitalbank/api"
	"capitalbank/pbapi"
	"capitalbank/store"
)

func StartExchange() error {
	acc := []store.Account{}
	err := store.LoadAccounts(&acc)
	if err != nil {
		return err
	}
	for _, a := range acc {
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

				tran, err := privat.GetTransactions()
				if err == nil {
					store.SaveTransactions(tran)
				}
			}
		}
	}
	return nil
}
