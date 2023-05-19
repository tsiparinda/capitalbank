package pbapi

import (
	"capitalbank/config"
	"capitalbank/logger"
	pb_models "capitalbank/pbapi/models"
	"capitalbank/store"
	"capitalbank/utils"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
)

type PrivatBankAPI struct {
	UserAgent   string
	Token       string
	ContentType string
	Account     string
	Direction   int
	BankRegistr string
}

func (a PrivatBankAPI) checkState() (bool, error) {
	state, err := a.GetState()
	if err != nil {
		fmt.Println("PrivatBankAPI.GetBalance Error running GetState")
		return false, err
	}

	if state.Phase != "WRK" || state.WorkBalance != "N" {
		return false, fmt.Errorf("PrivatBankAPI.GetBalance Privanbank is not in the WRK State!")
	}
	return true, nil
}

func (a PrivatBankAPI) GetBalance() ([]store.DataBalance, error) {
	// Implement the method for getting the balance from the PrivatBank API
	// ...
	state, err := a.checkState()
	if !state {
		return []store.DataBalance{}, err
	}
	// init vars
	databalance := []store.DataBalance{}
	var followId string

	limit := 10
	// take requestdays from config file and convert to int
	if value, ok := config.Config["RowsInPack"].(float64); ok {
		// The value is a float64, handle it accordingly
		limit = int(value)
	} else {
		logger.Log.WithFields(logrus.Fields{
			"rowsinpack": value,
		}).Infof("Error loading reqdays from config:", err.Error())
	}

	var reqdays int
	// take requestdays from config file and convert to int
	if value, ok := config.Config["RequestDays"].(float64); ok {
		// The value is a float64, handle it accordingly
		reqdays = int(value)
	} else {
		logger.Log.WithFields(logrus.Fields{
			"reqdays": value,
		}).Infof("Error loading reqdays from config:", err.Error())
		reqdays = 1
	}
	reqdays = 3
	dateTo := time.Now()
	dateFrom := dateTo.AddDate(0, 0, -reqdays)
	//	dateTo = dateTo.AddDate(0, 0, 1)
	//main cycle for receive all of packages
	for {
		responseData := pb_models.BalanceResponseData{}
		url, _ := a.CombineURL(
			PbURL{
				URL:       "https://acp.privatbank.ua/api/statements/balance",
				Acc:       a.Account,
				StartDate: dateFrom,
				EndDate:   dateTo,
				FollowId:  followId,
				Limit:     limit})
		req, _ := http.NewRequest("GET", url, nil)

		//	fmt.Printf("url:", url)

		logger.Log.WithFields(logrus.Fields{
			"url": url,
		}).Debugf("Request URL to take transactions:", url)

		req.Header.Add("User-Agent", a.UserAgent)
		req.Header.Add("token", a.Token)
		req.Header.Add("Content-Type", a.ContentType)

		res, err := http.DefaultClient.Do(req)
		if err != nil {
			//fmt.Printf("The HTTP request failed with error %s\n", err)
			logger.Log.WithFields(logrus.Fields{
				"err": err,
			}).Warnf("The HTTP request failed with error %s\n", err)
			return []store.DataBalance{}, err
		} else {
			data, _ := ioutil.ReadAll(res.Body)

			// Unmarshal the data into the struct
			json.Unmarshal(data, &responseData)

			if responseData.Status == "ERROR" {
				logger.Log.Warnf("PrivatBankAPI.GetBalance: Error status has got from Privatbank")
				return []store.DataBalance{}, err
			}

			for i, _ := range responseData.Balances {
				//save data to logs if debug level
				result, err := utils.StructToMap(responseData.Balances[i])
				if err != nil {
					logger.Log.WithFields(logrus.Fields{
						"err": err,
					}).Warnf("Error when StructToMap balances %s\n", err)
					return []store.DataBalance{}, err
				}
				result["bank"] = "privat"
				logger.Log.WithFields(result).Tracef("GET: ", url)
			}

			for i, _ := range responseData.Balances {
				// summ in coins!!!
				balanceIn, _ := strconv.ParseInt(strings.ReplaceAll(responseData.Balances[i].BalanceIn, ".", ""), 10, 64)
				balanceInEq, _ := strconv.ParseInt(strings.ReplaceAll(responseData.Balances[i].BalanceInEq, ".", ""), 10, 64)
				balanceOut, _ := strconv.ParseInt(strings.ReplaceAll(responseData.Balances[i].BalanceOut, ".", ""), 10, 64)
				balanceOutEq, _ := strconv.ParseInt(strings.ReplaceAll(responseData.Balances[i].BalanceOutEq, ".", ""), 10, 64)
				turnoverDebt, _ := strconv.ParseInt(strings.ReplaceAll(responseData.Balances[i].TurnoverDebt, ".", ""), 10, 64)
				turnoverDebtEq, _ := strconv.ParseInt(strings.ReplaceAll(responseData.Balances[i].TurnoverDebtEq, ".", ""), 10, 64)
				turnoverCred, _ := strconv.ParseInt(strings.ReplaceAll(responseData.Balances[i].TurnoverCred, ".", ""), 10, 64)
				turnoverCredEq, _ := strconv.ParseInt(strings.ReplaceAll(responseData.Balances[i].TurnoverCredEq, ".", ""), 10, 64)
				databalance = append(databalance,
					store.DataBalance{
						Direction:      a.Direction,
						BankRegistr:    a.BankRegistr,
						Acc:            a.Account,
						Currency:       responseData.Balances[i].Currency,
						BalanceIn:      balanceIn,
						BalanceInEq:    balanceInEq,
						BalanceOut:     balanceOut,
						BalanceOutEq:   balanceOutEq,
						TurnoverDebt:   turnoverDebt,
						TurnoverDebtEq: turnoverDebtEq,
						TurnoverCred:   turnoverCred,
						TurnoverCredEq: turnoverCredEq,
						Dpd:            responseData.Balances[i].Dpd,
						NameACC:        responseData.Balances[i].NameACC,
						IsFinalBal:     responseData.Balances[i].IsFinalBal,
						Source:         "B",
					})

			}
		}
		// If there is no next page, break the loop
		if !responseData.ExistNextPage {
			break
		}
		// Update followId for the next request
		followId = responseData.NextPageID

	}
	// trancount := len(databalance)
	// logger.Log.WithFields(logrus.Fields{}).Debugf("It was received %v transactions", trancount)

	return databalance, nil

}

func (a PrivatBankAPI) GetTransactions() ([]store.DataTransaction, error) {
	// Implement the method for getting the transactions from the PrivatBank API
	// This should return a slice of PrivatBankTransaction values, but as []Transaction
	// ...
	state, err := a.checkState()
	if !state {
		return []store.DataTransaction{}, err
	}

	// init vars
	datatrans := []store.DataTransaction{}
	var followId string

	limit := 10
	// take requestdays from config file and convert to int
	if value, ok := config.Config["RowsInPack"].(float64); ok {
		// The value is a float64, handle it accordingly
		limit = int(value)
	} else {
		logger.Log.WithFields(logrus.Fields{
			"rowsinpack": value,
		}).Infof("Error loading reqdays from config:", err.Error())
	}

	var reqdays int
	// take requestdays from config file and convert to int
	if value, ok := config.Config["RequestDays"].(float64); ok {
		// The value is a float64, handle it accordingly
		reqdays = int(value)
	} else {
		logger.Log.WithFields(logrus.Fields{
			"reqdays": value,
		}).Infof("Error loading reqdays from config:", err.Error())
		reqdays = 1
	}
	dateTo := time.Now()
	dateFrom := dateTo.AddDate(0, 0, -reqdays)
	//	dateTo = dateTo.AddDate(0, 0, 1)
	//main cycle for receive all of packages
	for {
		responseData := pb_models.TransactionResponseData{}
		url, _ := a.CombineURL(
			PbURL{
				URL:       "https://acp.privatbank.ua/api/statements/transactions",
				Acc:       a.Account,
				StartDate: dateFrom,
				EndDate:   dateTo,
				FollowId:  followId,
				Limit:     limit})
		req, _ := http.NewRequest("GET", url, nil)

		//	fmt.Printf("url:", url)

		logger.Log.WithFields(logrus.Fields{
			"url": url,
		}).Debugf("Request URL to take transactions:", url)

		req.Header.Add("User-Agent", a.UserAgent)
		req.Header.Add("token", a.Token)
		req.Header.Add("Content-Type", a.ContentType)

		res, err := http.DefaultClient.Do(req)
		// contentType := res.Header.Get("Content-Type")
		// fmt.Printf("encoding: ", contentType)
		if err != nil {
			logger.Log.WithFields(logrus.Fields{
				"err": err,
			}).Warnf("The HTTP request failed with error %s\n", err)
			return []store.DataTransaction{}, err
		} else {
			data, _ := ioutil.ReadAll(res.Body)

			// Unmarshal the data into the struct
			json.Unmarshal(data, &responseData)

			if responseData.Status == "ERROR" {
				logger.Log.Warnf("PrivatBankAPI.GetBalance: Error status has got from Privatbank")
				return []store.DataTransaction{}, err
			}

			for i, _ := range responseData.Transactions {
				//save data to logs if debug level
				result, err := utils.StructToMap(responseData.Transactions[i])
				if err != nil {
					logger.Log.WithFields(logrus.Fields{
						"err": err,
					}).Warnf("Error when StructToMap balances %s\n", err)
					return []store.DataTransaction{}, err
				}
				result["bank"] = "privat"
				logger.Log.WithFields(result).Tracef("GET: ", url)
			}

			for i, _ := range responseData.Transactions {
				if responseData.Transactions[i].PR_PR == "r" && responseData.Transactions[i].FL_REAL == "r" {
					// summa in coins!!!
					summa, _ := strconv.ParseInt(strings.ReplaceAll(responseData.Transactions[i].SUM_E, ".", ""), 10, 64)
					datatrans = append(datatrans,
						store.DataTransaction{
							Direction:   a.Direction,
							BankRegistr: a.BankRegistr,
							CntrCode:    responseData.Transactions[i].AUT_CNTR_CRF,
							CntrName:    responseData.Transactions[i].AUT_CNTR_NAM,
							CntrAcc:     responseData.Transactions[i].AUT_CNTR_ACC,
							Comment:     responseData.Transactions[i].OSND,
							DateTran:    responseData.Transactions[i].DAT_OD,
							ID:          responseData.Transactions[i].ID,
							TranType:    responseData.Transactions[i].TRANTYPE,
							SumTran:     summa,
						})
				}
			}
		}
		// If there is no next page, break the loop
		if !responseData.ExistNextPage {
			break
		}
		// Update followId for the next request
		followId = responseData.NextPageID

	}
	trancount := len(datatrans)
	logger.Log.WithFields(logrus.Fields{}).Debugf("It was received %v transactions", trancount)

	return datatrans, nil
}

func (a PrivatBankAPI) GetState() (store.DataState, error) {

	url := "https://acp.privatbank.ua/api/statements/settings"
	req, _ := http.NewRequest("GET", url, nil)

	req.Header.Add("User-Agent", a.UserAgent)
	req.Header.Add("token", a.Token)
	req.Header.Add("Content-Type", a.ContentType)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Printf("The HTTP request failed with error %s\n", err)
		return store.DataState{}, err
	} else {
		data, _ := ioutil.ReadAll(res.Body)

		// Create a new instance of ResponseData
		responseData := pb_models.SettingsData{}

		// Unmarshal the data into the struct
		json.Unmarshal(data, &responseData)

		// fmt.Println(string(data))

		datastate := store.DataState{
			Status:             responseData.Status,
			Type:               responseData.Type,
			Phase:              responseData.Settings.Phase,
			Today:              responseData.Settings.Today,
			LastDay:            responseData.Settings.LastDay,
			WorkBalance:        responseData.Settings.WorkBalance,
			ServerDateTime:     responseData.Settings.ServerDateTime,
			DateFinalStatement: responseData.Settings.DateFinalStatement,
		}
		//db.SaveState(datastate)

		return datastate, nil
	}
}
