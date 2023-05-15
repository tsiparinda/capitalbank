package api

import "capitalbank/store"

// In this example, BankAPI is an interface with methods for getting the balance and transactions.
// Transaction is another interface with methods for getting the transaction ID and amount.

// For each bank, you would create a struct (like PrivatBankAPI) that implements BankAPI and another struct (like PrivatBankTransaction) that implements Transaction.
// This way, you can write generic code that works with BankAPI and Transaction interfaces, and it will be able to handle different banks
// APIs as long as they provide the necessary methods.
// Keep in mind that this is a simplified example and you would likely need to adjust it to fit your specific needs.

type BankAPI interface {
	GetState() (store.DataState, error)
	GetBalance() (int64, error)
	GetTransactions() ([]store.DataTransaction, error)
}

type Transaction interface {
	GetID() string
	GetAmount() int64
	// ... other common transaction methods ...
}

// type PrivatBankTransaction struct {
//     // ... PrivatBank specific fields ...
// }

// func (t PrivatBankTransaction) GetID() string {
//     // Return the ID from the PrivatBank specific field
//     return t.ID
// }

// func (t PrivatBankTransaction) GetAmount() float64 {
//     // Return the amount from the PrivatBank specific field
//     amount, _ := strconv.ParseFloat(t.SUM, 64)
//     return amount
// }

// Similar struct and methods can be defined for other bank's transaction

// type PrivatBankAPI struct {
//     // ... fields for connecting to the PrivatBank API ...
// }

// func (api PrivatBankAPI) GetBalance() (float64, error) {
//     // Implement the method for getting the balance from the PrivatBank API
//     // ...
// }

// func (api PrivatBankAPI) GetTransactions() ([]Transaction, error) {
//     // Implement the method for getting the transactions from the PrivatBank API
//     // This should return a slice of PrivatBankTransaction values, but as []Transaction
//     // ...
// }

// Similar struct and methods can be defined for other bank's API
