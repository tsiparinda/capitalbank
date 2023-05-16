package store

import (
	"capitalbank/db"
	"capitalbank/logger"
	"database/sql"

	_ "github.com/denisenkom/go-mssqldb"
	"github.com/sirupsen/logrus"
)

// ...

func SaveTransactions(data []DataTransaction) {

	for i, _ := range data {
		// Insert data into database
		_, err := db.DB.Exec("INSERT INTO bank_rawtran (Direction, BankRegistr, CntrCode, CntrName, DateTran, Comment, SumTran, ID, TranType, CntrAcc) VALUES (@p1, @p2, @p3, @p4, @p5, @p6, @p7, @p8, @p9, @p10)",
			sql.Named("p1", data[i].Direction),
			sql.Named("p2", data[i].BankRegistr),
			sql.Named("p3", data[i].CntrCode),
			sql.Named("p4", data[i].CntrName),
			sql.Named("p5", data[i].DateTran),
			sql.Named("p6", data[i].Comment),
			sql.Named("p7", data[i].SumTran),
			sql.Named("p8", data[i].ID),
			sql.Named("p9", data[i].TranType),
			sql.Named("p10", data[i].CntrAcc))
		if err != nil {
			logger.Log.WithFields(logrus.Fields{
				"animal": "walrus",
				"size":   10,
			}).Debugf("Error inserting data into database:", err.Error())
		}
	}

	for i, _ := range data {
		// Insert data into database
		_, err := db.DB.Exec("exec bank_AddTransaction @Direction, @BankRegistr, @CntrCode, @CntrName, @CntrAcc, @DateTran, @Comment, @SumTran, @ID, @TranType",
			sql.Named("Direction", data[i].Direction),
			sql.Named("BankRegistr", data[i].BankRegistr),
			sql.Named("CntrCode", data[i].CntrCode),
			sql.Named("CntrName", data[i].CntrName),
			sql.Named("CntrAcc", data[i].CntrAcc),
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
