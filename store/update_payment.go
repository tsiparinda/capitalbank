package store

import (
	"capitalbank/db"
	"capitalbank/logger"
	"database/sql"
	"encoding/json"
	"os"

	_ "github.com/denisenkom/go-mssqldb"
	"github.com/sirupsen/logrus"
)

// ...

func UpdatePayment(p Payment, rsp PaymentResponse, rsperr error) {

	// Insert data into bank_payments_history and update payment's status
	// Encode payment struct into JSON
	rspJSON, err := json.Marshal(rsp)
	if err != nil {
		logger.Log.Error("UpdatePayment: Error Marshall response: ", err.Error())
		os.Exit(1)
	}

	logger.Log.WithFields(logrus.Fields{
		"id_key": p.ID_Key,
	}).Trace("Run sp.[bank_PaymentsResponseSave]")

	if _, err := db.DB.Exec("exec bank_PaymentsResponseSave @id_payment, @rsp, @err, @ref_num",
		sql.Named("id_payment", p.ID_Key),
		sql.Named("rsp", rspJSON),
		sql.Named("err", rsperr.Error()),
		sql.Named("ref_num", rsp.PaymentPackRef)); err != nil {
		logger.Log.WithFields(logrus.Fields{
			"id_key": p.ID_Key,
		}).Info("UpdatePayment: Error run sp.[bank_PaymentsResponseSave]: ", err.Error())
	}

}
