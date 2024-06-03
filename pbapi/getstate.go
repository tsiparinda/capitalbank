package pbapi

import (
	pb_models "capitalbank/pbapi/models"
	"capitalbank/store"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

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
