package ibcsvapi

import (
	"capitalbank/logger"
	"capitalbank/store"

	dec "github.com/shopspring/decimal"
	"github.com/sirupsen/logrus"
)

func (a IBankCSVAPI) GetTransactions() ([]store.DataTransaction, error) {
	datatrans := []store.DataTransaction{}
	records := a.Records
	//fmt.Println("from ibcsvapi: ", records)
	for i, _ := range records {
		//fmt.Println("from ibcsvapi: ", records[i].AUT_MY_ACC, a.Account)
		//check if account is equal
		if records[i].AUT_MY_ACC == a.Account {
			summ := dec.Decimal{}
			trantype := ""
			if records[i].SUM_CT.GreaterThan(records[i].SUM_DT) {
				summ = records[i].SUM_CT.Mul(dec.NewFromInt(100))
				trantype = "C"
			} else {
				summ = records[i].SUM_DT.Mul(dec.NewFromInt(100))
				trantype = "D"
			}
			// create unique identifier of transaction
			id := records[i].AUT_MY_ACC + "#" + records[i].NUM_DOC + "#" + records[i].AUT_CNTR_CRF
			datatrans = append(datatrans,
				store.DataTransaction{
					Direction:   a.Direction,
					BankRegistr: a.BankRegistr,
					CntrCode:    records[i].AUT_CNTR_CRF,
					CntrName:    records[i].AUT_CNTR_NAM,
					CntrAcc:     records[i].AUT_CNTR_ACC,
					Comment:     records[i].OSND,
					DateTran:    records[i].DATE_TIME_DAT_OD_TIM_P,
					ID:          id,
					TranType:    trantype,
					SumTran:     summ.BigInt().Int64(),
				})
		}
	}

	trancount := len(datatrans)
	logger.Log.WithFields(logrus.Fields{}).Debugf("It was received %v transactions", trancount)
	return datatrans, nil
}
