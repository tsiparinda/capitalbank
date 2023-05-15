package pbapi

import (
	"capitalbank/config"
	"capitalbank/logger"
	pb_models "capitalbank/pbapi/models"
	"capitalbank/store"
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

func (a PrivatBankAPI) GetBalance() (int64, error) {
	// Implement the method for getting the balance from the PrivatBank API
	// ...
	return 0, nil
}

func (a PrivatBankAPI) GetTransactions() ([]store.DataTransaction, error) {
	// Implement the method for getting the transactions from the PrivatBank API
	// This should return a slice of PrivatBankTransaction values, but as []Transaction
	// ...

	state, err := a.GetState()
	if err != nil {
		fmt.Println("PrivatBankAPI.GetTransactions Error running GetState")
		return []store.DataTransaction{}, err
	}

	if state.Phase != "WRK" {
		fmt.Println("PrivatBankAPI.GetTransactions Privanbank is not in the WRK State!")
		return []store.DataTransaction{}, nil
	}

	datatrans := []store.DataTransaction{}
	// Create a new instance of ResponseData

	var followId string

	limit := 10
	dateTo := time.Now()
	var intreqdays int
	// take requestdays from config file and convert to int
	if reqdays, ok := config.Config["requestDays"].(float64); ok {
		// The value is a float64, handle it accordingly
		intreqdays = int(reqdays)
	} else {
		logger.Log.WithFields(logrus.Fields{
			"reqdays": reqdays,
		}).Fatalf("Error loading reqdays from config:", err.Error())
	}

	dateFrom := dateTo.AddDate(0, 0, -intreqdays)

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

		//fmt.Printf("url:", url)

		req.Header.Add("User-Agent", a.UserAgent)
		req.Header.Add("token", a.Token)
		req.Header.Add("Content-Type", a.ContentType)

		res, err := http.DefaultClient.Do(req)
		if err != nil {
			fmt.Printf("The HTTP request failed with error %s\n", err)
			return []store.DataTransaction{}, err
		} else {
			data, _ := ioutil.ReadAll(res.Body)

			// Unmarshal the data into the struct
			json.Unmarshal(data, &responseData)

			if responseData.Status == "ERROR" {
				err = fmt.Errorf("PrivatBankAPI.GetTransactions: Error has got from Privatbank")
				fmt.Println(err.Error())
				//return []api.Transaction{}, err
			}
			// Unmarshal the data into the struct
			// b, err := json.MarshalIndent(responseData, "", "  ")
			// if err != nil {
			// 	log.Println(err)
			// 	return []store.DataTransaction{}, err
			// }
			// fmt.Println(string(b))
			//fmt.Println(responseData)

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

		datastate := store.DataState{Status: responseData.Status, Type: responseData.Type, Phase: responseData.Settings.Phase}
		//db.SaveState(datastate)

		return datastate, nil
	}
}
