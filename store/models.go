package store

import "database/sql"

// this data is loaded from db as list of bank's account and their parameters
type Account struct {
	Direction   int
	Account     string
	Bank        string
	Token       sql.NullString
	BankRegistr string
}

// some banks can return state of exchange in own format and this universal state we can save to db
type DataState struct {
	Status             string
	Type               string
	Phase              string
	Today              string
	LastDay            string
	WorkBalance        string
	ServerDateTime     string
	DateFinalStatement string
}

// info about one transaction which ready to save to db
type DataTransaction struct {
	Direction   int    // `json:"Capital_Direction"`
	BankRegistr string // registr of bank account in the capital2010
	CntrCode    string // OKPO, INN...
	CntrName    string // Nmae
	DateTran    string // date in our account
	Comment     string // comment
	SumTran     int64  // SUMMA in national currency
	ID          string // unique number
	TranType    string // D, C
}
