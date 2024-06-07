package pbapi

import (
	"capitalbank/config"
	"capitalbank/logger"
	pb_models "capitalbank/pbapi/models"
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

func (a PrivatBankAPI) GetBalance(datefrom time.Time) ([]store.DataBalance, error) {
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

	dateTo := time.Now()
	//main cycle for receive all of packages
	for {
		responseData := pb_models.BalanceResponseData{}
		url, _ := a.CombineURL(
			pb_models.PbURL{
				URL:       "https://acp.privatbank.ua/api/statements/balance",
				Acc:       a.Account,
				StartDate: datefrom,
				EndDate:   dateTo,
				FollowId:  followId,
				Limit:     limit})
		req, _ := http.NewRequest("GET", url, nil)

		logger.Log.WithFields(logrus.Fields{
			"url": url,
		}).Debug("Request URL to take balance:", url)

		req.Header.Add("User-Agent", a.UserAgent)
		req.Header.Add("token", a.Token)
		req.Header.Add("Content-Type", a.ContentType)

		res, err := http.DefaultClient.Do(req)
		if err != nil {
			logger.Log.WithFields(logrus.Fields{
				"err": err,
			}).Warnf("The HTTP request failed with error %s\n", err)
			return []store.DataBalance{}, err
		} else {
			data, _ := io.ReadAll(res.Body)
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
				logger.Log.WithFields(result).Trace("GET: ", url)
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
	return databalance, nil
}
