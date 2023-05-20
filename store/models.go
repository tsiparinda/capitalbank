package store

import (
	"database/sql"
	"time"
)

// this data is loaded from db as list of bank's account and their parameters
type Account struct {
	Direction    int
	Account      string
	Bank         string
	Token        sql.NullString
	BankRegistr  string
	BalanceState []DataCheckBalance
}

type DataCheckBalance struct {
	Dpd         time.Time // date
	BankRegistr string    // registr of bank account in the capital2010
	FBAD        bool      // state  of bank account's saldo in the capital2010
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
	CntrAcc     string //  Account
	DateTran    string // date in our account
	Comment     string // comment
	SumTran     int64  // SUMMA in national currency
	ID          string // unique number
	TranType    string // D, C
}

type DataBalance struct {
	Direction      int    // `json:"Capital_Direction"`
	BankRegistr    string // registr of bank account in the capital2010
	Acc            string // account number
	Currency       string // UAH
	BalanceIn      int64  // balance on begin of day
	BalanceInEq    int64  // in national currence
	BalanceOut     int64
	BalanceOutEq   int64
	TurnoverDebt   int64
	TurnoverDebtEq int64
	TurnoverCred   int64
	TurnoverCredEq int64
	// BgfIBrnm       string
	// Brnm           string
	Dpd     string // date
	NameACC string // name of account
	// State          string // 1
	// Atp            string // D
	// Flmn           string // DN
	// DateOpenAccReg string
	// DateOpenAccSys string
	// DateCloseAcc   string
	IsFinalBal bool   // !!!
	Source     string // C- Capital; B- Bank
}
