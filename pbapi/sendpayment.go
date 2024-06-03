package pbapi

import (
	"bytes"
	"capitalbank/store"
	"encoding/json"
	"errors"
	"net/http"
)

func (a PrivatBankAPI) SendPayment(payment store.Payment) (store.PaymentResponse, error) {
	// Implement the method for sending payment across the PrivatBank API
	// This should return a response from the API

	// Encode payment struct into JSON
	paymentJSON, err := json.Marshal(payment)
	if err != nil {
		return store.PaymentResponse{}, err
	}

	// Create HTTP request with payment JSON in the body
	url := "https://acp.privatbank.ua/api/proxy/payment/create"
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(paymentJSON))
	if err != nil {
		return store.PaymentResponse{}, err
	}

	// Add headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", a.UserAgent)
	req.Header.Set("token", a.Token)

	// Send HTTP request
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return store.PaymentResponse{}, err
	}
	defer res.Body.Close()

	// Decode response body into struct
	var responseData store.PaymentResponse
	if err := json.NewDecoder(res.Body).Decode(&responseData); err != nil {
		return store.PaymentResponse{}, err
	}

	if responseData.ResponseStatus == "ERROR" {
		return store.PaymentResponse{}, errors.New("error status received from PrivatBank API")
	}

	return responseData, nil
}

// func (a PrivatBankAPI) SendPaymentold(payment store.Payment) (store.PaymentResponse, error) {
// 	// Implement the method for send payment across the PrivatBank API
// 	// This should return a responce from api
// 	// ...
// 	state, err := a.checkState()
// 	if !state {
// 		return store.PaymentResponse{}, err
// 	}

// 	// Encode payment struct into JSON
// 	paymentJSON, err := json.Marshal(payment)
// 	if err != nil {
// 		return store.PaymentResponse{}, err
// 	}

// 	// init vars
// 	datatrans := store.PaymentResponse{}
// 	responseData := store.PaymentResponse{}

// 	url := "https://acp.privatbank.ua/api/proxy/payment/create"

// 	req, _ := http.NewRequest("GET", url, nil)

// 	logger.Log.WithFields(logrus.Fields{
// 		"url": url,
// 	}).Debugf("Request URL to send payment:", url)

// 	req.Header.Add("User-Agent", a.UserAgent)
// 	req.Header.Add("token", a.Token)
// 	req.Header.Add("Content-Type", a.ContentType)

// 	res, err := http.DefaultClient.Do(req)
// 	if err != nil {
// 		logger.Log.WithFields(logrus.Fields{
// 			"err": err,
// 		}).Warnf("The HTTP request failed with error %s\n", err)
// 		return store.PaymentResponse{}, err
// 	} else {
// 		data, _ := io.ReadAll(res.Body)
// 		// Unmarshal the data into the struct
// 		json.Unmarshal(data, &responseData)
// 		if responseData.ResponseStatus == "ERROR" {
// 			logger.Log.Warnf("PrivatBankAPI.SendPayment: Error status has got from Privatbank")
// 			return store.PaymentResponse{}, err
// 		}
// 	}
// 	return datatrans, nil
// }
