package pb_models


type SettingsData struct {
	Status   string `json:"status"`
	Type     string `json:"type"`
	Settings struct {
		Phase               string   `json:"phase"`
		DatesWithoutOperDay []string `json:"dates_without_oper_day"`
		Today               string   `json:"today"`
		LastDay             string   `json:"lastday"`
		WorkBalance         string   `json:"work_balance"`
		ServerDateTime      string   `json:"server_date_time"`
		DateFinalStatement  string   `json:"date_final_statement"`
	} `json:"settings"`
}
