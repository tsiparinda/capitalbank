package csv

import (
	"capitalbank/utils"
	"fmt"

	dec "github.com/shopspring/decimal"
)

type CSVRecord struct {
	AUT_MY_CRF             string      `json:"AUT_MY_CRF"`
	AUT_MY_MFO             string      `json:"AUT_MY_MFO"`
	AUT_MY_ACC             string      `json:"AUT_MY_ACC"`
	CURR                   string      `json:"CURR"`
	DATE_TIME_DAT_OD_TIM_P string      `json:"DATE_TIME_DAT_OD_TIM_P"`
	ID                     string      `json:"ID"`
	AUT_BANK_CRF           string      `json:"AUT_BANK_CRF"`
	AUT_BANK_MFO           string      `json:"AUT_BANK_MFO"`
	AUT_CNTR_ACC           string      `json:"AUT_CNTR_ACC"`
	AUT_CNTR_CRF           string      `json:"AUT_CNTR_CRF"`
	AUT_CNTR_NAM           string      `json:"AUT_CNTR_NAM"`
	NUM_DOC                string      `json:"NUM_DOC"`
	DAT_KL                 string      `json:"DAT_KL"`
	SUM_DT                 dec.Decimal `json:"SUM_DT"`
	SUM_CT                 dec.Decimal `json:"SUM_CT"`
	OSND                   string      `json:"OSND"`
	SUM_E                  dec.Decimal `json:"SUM_E"`
}

type CSVfiles struct {
	FileName string `json:"file_name"`
}

func (c CSVRecord) GetRecord() CSVRecord {
	return c
}

func (c CSVfiles) GetFile() CSVfiles {
	return c
}

func NewCSVRecord(line []string) *CSVRecord {
	//fmt.Printf("line: %v %v %v \n", line[13], line[14], line[16])
	//sumdt := dec.RequireFromString(line[13]).Round(2)
	sumdt, err := utils.Str2Dec(line[13], 2)
	if err != nil {
		fmt.Println("Failed to convert sumdt to decimal: ", err)
		return &CSVRecord{}
	}
	// sumct := dec.RequireFromString(line[14]).Round(2)
	sumct, err := utils.Str2Dec(line[14], 2)
	if err != nil {
		fmt.Println("Failed to convert sumct to decimal:", err)
		return &CSVRecord{}
	}
	//	sume := dec.RequireFromString(line[16]).Round(2)
	sume, err := utils.Str2Dec(line[16], 2)
	if err != nil {
		fmt.Println("Failed to convert sume to decimal:", err)
		return &CSVRecord{}
	}
	return &CSVRecord{
		AUT_MY_CRF:             line[0],
		AUT_MY_MFO:             line[1],
		AUT_MY_ACC:             line[2],
		CURR:                   line[3],
		DATE_TIME_DAT_OD_TIM_P: line[4],
		ID:                     line[5],
		AUT_BANK_CRF:           line[6],
		AUT_BANK_MFO:           line[7],
		AUT_CNTR_ACC:           line[8],
		AUT_CNTR_CRF:           line[9],
		AUT_CNTR_NAM:           line[10],
		NUM_DOC:                line[11],
		DAT_KL:                 line[12],
		SUM_DT:                 sumdt,
		SUM_CT:                 sumct,
		OSND:                   line[15],
		SUM_E:                  sume,
	}
}
