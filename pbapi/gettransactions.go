package pbapi

import (
	"capitalbank/config"
	"capitalbank/logger"
	models "capitalbank/pbapi/models"
	"capitalbank/store"
	"capitalbank/utils"
	"encoding/json"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
)

func (a PrivatBankAPI) GetTransactions() ([]store.DataTransaction, error) {
	// Implement the method for getting transactions from the PrivatBank API
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
		}).Info("Error loading reqdays from config:", err.Error())
	}

	var reqdays int
	// take requestdays from config file and convert to int
	if value, ok := config.Config["RequestDaysTrans"].(float64); ok {
		// The value is a float64, handle it accordingly
		reqdays = int(value)
	} else {
		logger.Log.WithFields(logrus.Fields{
			"reqdays": value,
		}).Info("Error loading reqdays from config:", err.Error())
		reqdays = 1
	}
	dateTo := time.Now()
	dateFrom := dateTo.AddDate(0, 0, -reqdays)
	//main cycle for receive all of packages
	for {
		responseData := models.TransactionResponseData{}
		url, _ := a.CombineURL(
			models.PbURL{
				URL:       "https://acp.privatbank.ua/api/statements/transactions",
				Acc:       a.Account,
				StartDate: dateFrom,
				EndDate:   dateTo,
				FollowId:  followId,
				Limit:     limit})
		req, _ := http.NewRequest("GET", url, nil)

		logger.Log.WithFields(logrus.Fields{
			"url": url,
		}).Debug("Request URL to take transactions:", url)

		req.Header.Add("User-Agent", a.UserAgent)
		req.Header.Add("token", a.Token)
		req.Header.Add("Content-Type", a.ContentType)

		res, err := http.DefaultClient.Do(req)
		if err != nil {
			logger.Log.WithFields(logrus.Fields{
				"err": err,
			}).Warnf("The HTTP request failed with error %s\n", err)
			return []store.DataTransaction{}, err
		} else {
			data, _ := io.ReadAll(res.Body) // was ioutil
			// Unmarshal the data into the struct
			json.Unmarshal(data, &responseData)
			if responseData.Status == "ERROR" {
				logger.Log.WithFields(logrus.Fields{
					"acc": a.Account,
					"err": err,
				}).Warnf("PrivatBankAPI.GetTransactions: Error status has got from Privatbank")
				return []store.DataTransaction{}, err
			}

			logger.Log.WithFields(logrus.Fields{
				"responseData": responseData,
			}).Trace("Requested transactions by account:", a.Account)

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
				logger.Log.WithFields(result).Trace("GET: ", url)
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
							NumDoc:      responseData.Transactions[i].NUM_DOC,
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
