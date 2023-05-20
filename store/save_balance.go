package store

import (
	"capitalbank/db"
	"capitalbank/logger"
	"database/sql"

	_ "github.com/denisenkom/go-mssqldb"
	"github.com/sirupsen/logrus"
)

// ...

func SaveBalance(data []DataBalance) {

	// for i, _ := range data {
	// 	// Insert data into database
	// 	_, err := db.DB.Exec("INSERT INTO bank_balances (Direction, BankRegistr, CntrCode, CntrName, DateTran, Comment, SumTran, ID, TranType, CntrAcc) VALUES (@p1, @p2, @p3, @p4, @p5, @p6, @p7, @p8, @p9, @p10)",
	// 		sql.Named("p1", data[i].Direction),
	// 		sql.Named("p2", data[i].BankRegistr),
	// 		sql.Named("p3", data[i].CntrCode),
	// 		sql.Named("p4", data[i].CntrName),
	// 		sql.Named("p5", data[i].DateTran),
	// 		sql.Named("p6", data[i].Comment),
	// 		sql.Named("p7", data[i].SumTran),
	// 		sql.Named("p8", data[i].ID),
	// 		sql.Named("p9", data[i].TranType),
	// 		sql.Named("p10", data[i].CntrAcc))
	// 	if err != nil {
	// 		logger.Log.WithFields(logrus.Fields{
	// 			"animal": "walrus",
	// 			"size":   10,
	// 		}).Debugf("Error inserting data into database:", err.Error())
	// 	}
	// }

	for i, _ := range data {
		func insertIntoBankBalances(db *sql.DB, data DataBalance) error {
			_, err := db.Exec("EXEC dbo.InsertIntoBankBalances @BankRegistr=@p1, @Dpd=@p2, @Source=@p3, @Acc=@p4, @Currency=@p5, @BalanceIn=@p6, @BalanceInEq=@p7, @BalanceOut=@p8, @BalanceOutEq=@p9, @TurnoverDebt=@p10, @TurnoverDebtEq=@p11, @TurnoverCred=@p12, @TurnoverCredEq=@p13, @IsFinalBal=@p14", 
			data.BankRegistr, 
			data.Dpd, 
			data.Source, 
			data.Acc, 
			data.Currency, 
			data.BalanceIn, 
			data.BalanceInEq, 
			data.BalanceOut, 
			data.BalanceOutEq, 
			data.TurnoverDebt, 
			data.TurnoverDebtEq, 
			data.TurnoverCred, 
			data.TurnoverCredEq, 
			data.IsFinalBal)
			return err
		}
		
		// Insert data into database
		_, err := db.DB.Exec("exec bank_AddBalance @Direction, @BankRegistr, @CntrCode, @CntrName, @CntrAcc, @DateTran, @Comment, @SumTran, @ID, @TranType",
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
