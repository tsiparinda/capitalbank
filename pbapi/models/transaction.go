package pb_models

import "strconv"

type Transaction struct {
	AUT_MY_CRF               string `json:"AUT_MY_CRF"`
	AUT_MY_MFO               string `json:"AUT_MY_MFO"`
	AUT_MY_ACC               string `json:"AUT_MY_ACC"`
	AUT_MY_NAM               string `json:"AUT_MY_NAM"`
	AUT_MY_MFO_NAME          string `json:"AUT_MY_MFO_NAME"`
	AUT_MY_MFO_CITY          string `json:"AUT_MY_MFO_CITY"`
	AUT_CNTR_CRF             string `json:"AUT_CNTR_CRF"`
	AUT_CNTR_MFO             string `json:"AUT_CNTR_MFO"`
	AUT_CNTR_ACC             string `json:"AUT_CNTR_ACC"`
	AUT_CNTR_NAM             string `json:"AUT_CNTR_NAM"`
	AUT_CNTR_MFO_NAME        string `json:"AUT_CNTR_MFO_NAME"`
	AUT_CNTR_MFO_CITY        string `json:"AUT_CNTR_MFO_CITY"`
	CCY                      string `json:"CCY"`
	FL_REAL                  string `json:"FL_REAL"`
	PR_PR                    string `json:"PR_PR"`
	DOC_TYP                  string `json:"DOC_TYP"`
	NUM_DOC                  string `json:"NUM_DOC"`
	DAT_KL                   string `json:"DAT_KL"`
	DAT_OD                   string `json:"DAT_OD"`
	OSND                     string `json:"OSND"`
	SUM                      string `json:"SUM"`
	SUM_E                    string `json:"SUM_E"`
	REF                      string `json:"REF"`
	REFN                     string `json:"REFN"`
	TIM_P                    string `json:"TIM_P"`
	DATE_TIME_DAT_OD_TIM_P   string `json:"DATE_TIME_DAT_OD_TIM_P"`
	ID                       string `json:"ID"`
	TRANTYPE                 string `json:"TRANTYPE"`
	DLR                      string `json:"DLR"`
	TECHNICAL_TRANSACTION_ID string `json:"TECHNICAL_TRANSACTION_ID"`
}

type TransactionResponseData struct {
	Status        string        `json:"status"`
	Type          string        `json:"type"`
	ExistNextPage bool          `json:"exist_next_page"`
	NextPageID    string        `json:"next_page_id"`
	Transactions  []Transaction `json:"transactions"`
}

func (t Transaction) GetID() string {
	// Return the ID from the PrivatBank specific field
	return t.ID
}

func (t Transaction) GetAmount() int64 {
	// Return the amount from the PrivatBank specific field
	amount, _ := strconv.ParseInt(t.SUM, 10, 64)
	return amount
}
