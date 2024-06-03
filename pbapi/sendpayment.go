package pbapi

import (
	"capitalbank/logger"
	"capitalbank/store"
	"capitalbank/utils"
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/sirupsen/logrus"
)

func (a PrivatBankAPI) SendPayment() (store.PaymentResponce, error) {
	// Implement the method for send payment across the PrivatBank API
	// This should return a responce from api
	// ...
	state, err := a.checkState()
	if !state {
		return store.PaymentResponce{}, err
	}

	// init vars
	datatrans := store.PaymentResponce{}
	//var followId string

	// limit := 10
	// // take requestdays from config file and convert to int
	// if value, ok := config.Config["RowsInPack"].(float64); ok {
	// 	// The value is a float64, handle it accordingly
	// 	limit = int(value)
	// } else {
	// 	logger.Log.WithFields(logrus.Fields{
	// 		"rowsinpack": value,
	// 	}).Infof("Error loading reqdays from config:", err.Error())
	// }

	// var reqdays int
	// // take requestdays from config file and convert to int
	// if value, ok := config.Config["RequestDaysTrans"].(float64); ok {
	// 	// The value is a float64, handle it accordingly
	// 	reqdays = int(value)
	// } else {
	// 	logger.Log.WithFields(logrus.Fields{
	// 		"reqdays": value,
	// 	}).Infof("Error loading reqdays from config:", err.Error())
	// 	reqdays = 1
	// }
	// dateTo := time.Now()
	// dateFrom := dateTo.AddDate(0, 0, -reqdays)
	//main cycle for receive all of packages
	for {
		responseData := store.PaymentResponce{}
		// url, _ := a.CombineURL(
		// 	models.PbURL{
		// 		URL: "https://acp.privatbank.ua/api/proxy/payment/create",
		// 		Acc: a.Account,
		// 		//StartDate: dateFrom,
		// 		//EndDate:   dateTo,
		// 		//FollowId:  followId,
		// 		//Limit:     limit
		// 	})
		url := "https://acp.privatbank.ua/api/proxy/payment/create"

		req, _ := http.NewRequest("GET", url, nil)

		logger.Log.WithFields(logrus.Fields{
			"url": url,
		}).Debugf("Request URL to send payment:", url)

		req.Header.Add("User-Agent", a.UserAgent)
		req.Header.Add("token", a.Token)
		req.Header.Add("Content-Type", a.ContentType)

		res, err := http.DefaultClient.Do(req)
		if err != nil {
			logger.Log.WithFields(logrus.Fields{
				"err": err,
			}).Warnf("The HTTP request failed with error %s\n", err)
			return store.PaymentResponce{}, err
		} else {
			data, _ := ioutil.ReadAll(res.Body)
			// Unmarshal the data into the struct
			json.Unmarshal(data, &responseData)
			if responseData.Status == "ERROR" {
				logger.Log.Warnf("PrivatBankAPI.GetBalance: Error status has got from Privatbank")
				return store.PaymentResponce{}, err
			}

			for i, _ := range responseData.Transactions {
				//save data to logs if debug level
				result, err := utils.StructToMap(responseData.Transactions[i])
				if err != nil {
					logger.Log.WithFields(logrus.Fields{
						"err": err,
					}).Warnf("Error when StructToMap balances %s\n", err)
					return store.PaymentResponce{}, err
				}
				result["bank"] = "privat"
				logger.Log.WithFields(result).Tracef("GET: ", url)
			}

			for i, _ := range responseData.Transactions {
				if responseData.Transactions[i].PR_PR == "r" && responseData.Transactions[i].FL_REAL == "r" {
					// summa in coins!!!
					// summa, _ := strconv.ParseInt(strings.ReplaceAll(responseData.Transactions[i].SUM_E, ".", ""), 10, 64)
					// datatrans = append(datatrans,
					// 	store.Payment{
					// 		Direction:   a.Direction,
					// 		BankRegistr: a.BankRegistr,
					// 		CntrCode:    responseData.Transactions[i].AUT_CNTR_CRF,
					// 		CntrName:    responseData.Transactions[i].AUT_CNTR_NAM,
					// 		CntrAcc:     responseData.Transactions[i].AUT_CNTR_ACC,
					// 		Comment:     responseData.Transactions[i].OSND,
					// 		DateTran:    responseData.Transactions[i].DAT_OD,
					// 		ID:          responseData.Transactions[i].ID,
					// 		TranType:    responseData.Transactions[i].TRANTYPE,
					// 		SumTran:     summa,
					// 	})
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
	//trancount := len(datatrans)
	//logger.Log.WithFields(logrus.Fields{}).Debugf("It was received %v transactions", trancount)
	return datatrans, nil
}
