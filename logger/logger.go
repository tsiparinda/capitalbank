package logger

import (
	"capitalbank/config"
	"capitalbank/db"
	"database/sql"
	"encoding/json"
	"fmt"

	_ "github.com/denisenkom/go-mssqldb"
	"github.com/sirupsen/logrus"
	"golang.org/x/text/encoding/charmap"
)

// HOW TO USE
// simple using fields
// 	logger.Log.WithFields(logrus.Fields{
// "ID": data[i].ID,
// }).Errorf("Error inserting data into database:", err.Error())
//
// dynamic fields
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
// or
// for i, _ := range responseData.Transactions {
// 	//save data to logs if debug level
// 	result, err := utils.StructToMap(responseData.Transactions[i])
// 	if err != nil {
// 		fmt.Printf(err.Error())
// 	}
// 	result["bank"] = "privat"
// 	logger.Log.WithFields(result).Tracef("GET: ", url)
// }

type MSSQLHook struct{}

func NewMSSQLHook() *MSSQLHook {
	return &MSSQLHook{}
}

func (h *MSSQLHook) Fire(entry *logrus.Entry) error {

	params, err := json.Marshal(entry.Data) // Convert the Fields map to a JSON string
	if err != nil {
		fmt.Println("Logger.Fire: Error Marshall entry.Data")
		return err
	}

	message, _ := EncodeWindows1251([]uint8(entry.Message))
	if err != nil {
		fmt.Println("Logger.Fire: Error EncodeWindows1251 Message")
		return err
	}
	params, err = EncodeWindows1251(params)
	if err != nil {
		fmt.Println("Logger.Fire: Error EncodeWindows1251 Params")
		return err
	}

	query := `INSERT INTO bank_logs (loglevel, message, params, time)VALUES (@p1, @p2, @p3, @p4)`
	_, err = db.DB.Exec(query,
		sql.Named("p1", entry.Level.String()),
		sql.Named("p2", message),
		sql.Named("p3", params),
		sql.Named("p4", entry.Time))
	return err
}

func (h *MSSQLHook) Levels() []logrus.Level {
	return logrus.AllLevels
}

var Log = logrus.New()

func init() {
	hook := NewMSSQLHook()
	// hook, err := NewMSSQLHook("server=bold;database=capital2010;Integrated Security=SSPI")
	// if err != nil {
	// 	Log.Fatal(err)
	// }

	Log.AddHook(hook)
	// Log.SetFormatter(&logrus.TextFormatter{
	// 	DisableColors: true,
	// 	ForceColors:   false,
	// })

	Log.SetFormatter(&logrus.JSONFormatter{})

	// trace, debug, info, warn, error, fatal, panic
	loglevel := config.Config["logLevel"].(string)
	level, err := logrus.ParseLevel(loglevel)
	if err != nil {
		fmt.Printf("Error parsing level: %v\n", err)
		loglevel = "Warn"
		level, _ = logrus.ParseLevel(loglevel)
	}
	Log.SetLevel(level)

}

func SetLogLevel(level string) {
	switch level {
	case "Debug":
		Log.SetLevel(logrus.DebugLevel)
	case "Info":
		Log.SetLevel(logrus.InfoLevel)
	case "Warn":
		Log.SetLevel(logrus.WarnLevel)
	case "Error":
		Log.SetLevel(logrus.ErrorLevel)
	case "Fatal":
		Log.SetLevel(logrus.FatalLevel)
	case "Panic":
		Log.SetLevel(logrus.PanicLevel)
	default:
		Log.SetLevel(logrus.InfoLevel)
	}
}

func utf8ToWin1251(input string) (string, error) {
	decoder := charmap.Windows1251.NewEncoder()
	output, err := decoder.String(input)
	if err != nil {
		return "", err
	}
	return output, nil
}

func EncodeWindows1251(ba []uint8) ([]uint8, error) {
	enc := charmap.Windows1251.NewEncoder()
	out, err := enc.String(string(ba))
	if err != nil {
		return []uint8(""), err
	}
	return []uint8(out), nil
}
