package store

import (
	"capitalbank/db"
	"capitalbank/logger"
	"database/sql"

	_ "github.com/denisenkom/go-mssqldb"
	"github.com/sirupsen/logrus"

	"golang.org/x/text/encoding/charmap"
)

// ...

func utf8ToWin1251(input string) (string, error) {
	decoder := charmap.Windows1251.NewDecoder()
	output, err := decoder.String(input)
	if err != nil {
		return "", err
	}
	return output, nil
}

func SaveTransactions(data []DataTransaction) {

	// Connect to database
	// db, err := sql.Open("sqlserver", "server=bold;database=capital2010;Integrated Security=SSPI")
	// if err != nil {
	// 	fmt.Println("Error opening database:", err.Error())
	// }
	// defer db.Close()

	for i, _ := range data {
		//cntrName, _ := utf8ToWin1251(data[i].CntrName)
		// Insert data into database
		_, err := db.DB.Exec("INSERT INTO bank_rawtran (Direction, BankRegistr, CntrCode, CntrName, DateTran, Comment, SumTran, ID, TranType) VALUES (@p1, @p2, @p3, @p4, @p5, @p6, @p7, @p8, @p9)",
			sql.Named("p1", data[i].Direction),
			sql.Named("p2", data[i].BankRegistr),
			sql.Named("p3", data[i].CntrCode),
			sql.Named("p4", data[i].CntrName),
			sql.Named("p5", data[i].DateTran),
			sql.Named("p6", data[i].Comment),
			sql.Named("p7", data[i].SumTran),
			sql.Named("p8", data[i].ID),
			sql.Named("p9", data[i].TranType))
		if err != nil {
			logger.Log.WithFields(logrus.Fields{
				"animal": "walrus",
				"size":   10,
			}).Debugf("Error inserting data into database:", err.Error())
		}
	}

	for i, _ := range data {
		//comment, _ := utf8ToWin1251(data[i].Comment)
		// if err != nil {
		// 	// handle error
		// }

		// save data to logs if debug level
		// fields := make(map[string]interface{})
		// fields["Direction"] = data[i].Direction
		// fields["BankRegistr"] = data[i].BankRegistr
		// fields["CntrCode"] = data[i].CntrCode
		// fields["CntrName"] = data[i].CntrName
		// fields["DateTran"] = data[i].DateTran
		// fields["Comment"] = data[i].Comment
		// fields["SumTran"] = data[i].SumTran
		// fields["ID"] = data[i].ID
		// fields["TranType"] = data[i].TranType
		// logger.Log.WithFields(fields).Debug("Try to save record")

		// Insert data into database
		_, err := db.DB.Exec("exec bank_AddTransaction @Direction, @BankRegistr, @CntrCode, @CntrName, @DateTran, @Comment, @SumTran, @ID, @TranType",
			sql.Named("Direction", data[i].Direction),
			sql.Named("BankRegistr", data[i].BankRegistr),
			sql.Named("CntrCode", data[i].CntrCode),
			sql.Named("CntrName", data[i].CntrName),
			sql.Named("DateTran", data[i].DateTran),
			sql.Named("Comment", data[i].Comment),
			sql.Named("SumTran", data[i].SumTran),
			sql.Named("ID", data[i].ID),
			sql.Named("TranType", data[i].TranType))
		if err != nil {
			logger.Log.WithFields(logrus.Fields{
				"ID": data[i].ID,
			}).Errorf("Error inserting data into database:", err.Error())
		}
	}

}
