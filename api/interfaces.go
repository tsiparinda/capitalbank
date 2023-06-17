package api

import (
	"capitalbank/store"
	"time"
)

// In this example, BankAPI is an interface with methods for getting the balance and transactions.
// For each bank, you would create a struct (like PrivatBankAPI) that implements BankAPI
// This way, you can write generic code that works with BankAPI interface, and it will be able to handle different banks
// APIs as long as they provide the necessary methods.

type BankAPI interface {
	GetState() (store.DataState, error)
	GetBalance(time.Time) ([]store.DataBalance, error)
	GetTransactions() ([]store.DataTransaction, error)
}
