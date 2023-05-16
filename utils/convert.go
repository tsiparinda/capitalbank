package utils

import (
	"encoding/json"

	"golang.org/x/text/encoding/charmap"
)

func StructToMap(item interface{}) (map[string]interface{}, error) {
	data, err := json.Marshal(item)
	if err != nil {
		return nil, err
	}
	var result map[string]interface{}
	err = json.Unmarshal(data, &result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func utf8ToWin1251(input string) (string, error) {
	decoder := charmap.Windows1251.NewDecoder()
	output, err := decoder.String(input)
	if err != nil {
		return "", err
	}
	return output, nil
}
