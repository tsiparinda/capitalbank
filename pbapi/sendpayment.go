package pbapi

import (
	"bytes"
	"capitalbank/store"
	"capitalbank/utils"
	"encoding/json"
	"errors"
	"io"
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
	// Convert the JSON data to CP-1251 encoding
	paymentcp1251, err := utils.ConvertToCP1251(paymentJSON)
	if err != nil {
		//fmt.Println("PrivatBankAPI.SendPayment Error converting to CP-1251:", err) //!!!!!!!
		return store.PaymentResponse{}, err
	}
	//fmt.Println("test JSON", paymentcp1251)
	// Create HTTP request with payment JSON in the body
	url := "https://acp.privatbank.ua/api/proxy/payment/create"
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(paymentcp1251))
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

	var responseData store.PaymentResponse
	data, _ := io.ReadAll(res.Body) // was ioutil
	// Unmarshal the data into the struct
	json.Unmarshal(data, &responseData)

	//fmt.Println("PrivatBankAPI.SendPayment responseData:  ", responseData) //!!!!
	if responseData.ResponseStatus == "ERROR" {
		return responseData, errors.New("error status received from PrivatBank API")
	}

	return responseData, nil
}
