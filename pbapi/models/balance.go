package pb_models

type Balance struct {
	Acc            string `json:"acc"`
	Currency       string `json:"currency"`
	BalanceIn      string `json:"balanceIn"`
	BalanceInEq    string `json:"balanceInEq"`
	BalanceOut     string `json:"balanceOut"`
	BalanceOutEq   string `json:"balanceOutEq"`
	TurnoverDebt   string `json:"turnoverDebt"`
	TurnoverDebtEq string `json:"turnoverDebtEq"`
	TurnoverCred   string `json:"turnoverCred"`
	TurnoverCredEq string `json:"turnoverCredEq"`
	BgfIBrnm       string `json:"bgfIBrnm"`
	Brnm           string `json:"brnm"`
	Dpd            string `json:"dpd"`
	NameACC        string `json:"nameACC"`
	State          string `json:"state"`
	Atp            string `json:"atp"`
	Flmn           string `json:"flmn"`
	DateOpenAccReg string `json:"date_open_acc_reg"`
	DateOpenAccSys string `json:"date_open_acc_sys"`
	DateCloseAcc   string `json:"date_close_acc"`
	IsFinalBal     bool   `json:"is_final_bal"`
}

type BalanceResponseData struct {
	Status        string    `json:"status"`
	Type          string    `json:"type"`
	ExistNextPage bool      `json:"exist_next_page"`
	NextPageID    string    `json:"next_page_id"`
	Balances      []Balance `json:"balances"`
}

func GetBalance() (int64, error) {
	return 0, nil
}

func GetTransactions() ([]Transaction, error) {
	return []Transaction{}, nil
}
