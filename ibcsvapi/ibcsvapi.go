package ibcsvapi

import (
	"capitalbank/csv"
	"capitalbank/store"
	"time"
)

type IBankCSVAPI struct {
	Account     string
	Direction   int
	BankRegistr string
	Records     []csv.CSVRecord
}

func (i IBankCSVAPI) GetBalance(datefrom time.Time) ([]store.DataBalance, error) {
	return []store.DataBalance{}, nil
}

func (i IBankCSVAPI) GetState() (store.DataState, error) {
	return store.DataState{}, nil
}
