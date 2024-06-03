package store

import (
	"capitalbank/db"
	"capitalbank/logger"

	_ "github.com/denisenkom/go-mssqldb"
)

func LoadPaymentsPrivat(payments *[]Payment) error {

	// Select data from database
	rows, err := db.DB.Query("SELECT document_number, payment_date, payment_destination, payer_account, recipient_account, recipient_nceo, Amount, Token  FROM bank_payments2send_privat")
	if err != nil {
		logger.Log.Errorf("Error loading payments from database:", err.Error())
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var p Payment
		// Scan each column into the corresponding field of an Account. Adjust this line as needed based on your table structure.
		err = rows.Scan(&p.DocumentNumber, &p.PaymentDate, &p.PaymentDestination, &p.PayerAccount, &p.RecipientAccount, &p.RecipientNceo, &p.PaymentAmount, &p.Token.String)
		if err != nil {
			logger.Log.Errorf("Error scanning accounts rows:", err.Error())
			return err
		}
		*payments = append(*payments, p)
	}

	// Check for errors from iterating over rows.
	if err := rows.Err(); err != nil {
		logger.Log.Errorf("Error iterating accounts rows:", err.Error())
		return err
	}

	return nil
}
