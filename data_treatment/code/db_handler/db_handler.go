package db_handler

import (
	"data_treatment/config_handler"
	"data_treatment/log_handler"
	"data_treatment/mail_handler"
	"data_treatment/shared"
	"database/sql"
	"os"
	"path/filepath"
	"reflect"

	"github.com/goccy/go-json"
	_ "github.com/mattn/go-sqlite3"
)

var credentials map[string]string
var receivers []string

func init() {
	val, err := config_handler.GetData(shared.Consts, "DEFAULT_MAIL_PATH", config_handler.TYPE_STRING)
	if err != nil {
		log_handler.Error("Error loading the config: %s", err)
		os.Exit(1)
	}
	mail_path := val.(string)
	raw_cred, err := config_handler.LoadConsts(mail_path)
	credentials = makeStringMap(raw_cred)
	if err != nil {
		log_handler.Error("Error loading credentials: %s", err)
		os.Exit(1)
	}

	val, err = config_handler.GetData(shared.Consts, "DEFAULT_SEND_TO_PATH", config_handler.TYPE_STRING)
	if err != nil {
		log_handler.Error("Error loading the send to config: %s", err)
		os.Exit(1)
	}
	send_to_path := val.(string)
	receivers, err = config_handler.LoadReceivers(send_to_path)
	if err != nil {
		log_handler.Error("Error loading receivers")
		os.Exit(1)
	}
}

func connectToDB(db_file string) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", db_file)
	if err != nil {
		return nil, err
	}
	create_table := `
    CREATE TABLE IF NOT EXISTS winery_data (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        temperature REAL,
        humidity REAL,
        date TIMESTAMP DEFAULT CURRENT_TIMESTAMP
    );`
	_, err = db.Exec(create_table)
	if err != nil {
		db.Close()
		return nil, err
	}
	return db, nil
}

func insertData(db *sql.DB, temperature float64, humidity float64) {
	insert_query := `INSERT INTO winery_data (temperature, humidity) VALUES (?, ?)`
	_, err := db.Exec(insert_query, temperature, humidity)
	if err != nil {
		log_handler.Error("Error inserting data: %s", err)
		return
	}
	log_handler.Success("Data inserted successfully [temperature = %.2f, humidity = %.2f]", temperature, humidity)
}

func JsonToData(data string) {
	var result map[string]any
	if err := json.Unmarshal([]byte(data), &result); err != nil {
		log_handler.Error("Error parsing JSON: %s", err)
		return
	}
	if result["presence"] == true {
		mail_handler.SendMail(credentials["MAIL_ACCOUNT"], credentials["MAIL_PASSWORD"], receivers, credentials["SMTP_SERVER"], "Alerta de presencia en la bodega.", "Se ha detectado presencia en la bodega.")
	}
	if result["fire"] == true {
		mail_handler.SendMail(credentials["MAIL_ACCOUNT"], credentials["MAIL_PASSWORD"], receivers, credentials["SMTP_SERVER"], "Alerta de fuego en la bodega.", "Hay fuego en la bodega.\nPuede ser 1 falso positivo.")
	}
	temp, ok_temp := validateType(result, "temperature", float64(0.0))
	hum, ok_hum := validateType(result, "humidity", float64(0.0))
	if ok_temp && ok_hum {
		val, err := config_handler.GetData(shared.Consts, "DEFAULT_DB_PATH", config_handler.TYPE_STRING)
		if err != nil {
			log_handler.Error("Error loading the config: %s", err)
			os.Exit(1)
		}
		err = os.MkdirAll(val.(string), 0755)
		if err != nil {
			log_handler.Error("Error creating database folder: %s", err)
			return
		}
		filename, err := config_handler.GetData(shared.Consts, "DB_FILENAME", config_handler.TYPE_STRING)
		if err != nil {
			log_handler.Error("Error loading the config: %s", err)
			os.Exit(1)
		}
		full_path := filepath.Join(val.(string), filename.(string))
		db, err := connectToDB(full_path)
		if err != nil {
			log_handler.Error("Error connecting to database: %s", err)
			return
		}
		defer db.Close()
		insertData(db, temp.(float64), hum.(float64))
	}
}

func validateType(data map[string]any, field string, expected_type any) (any, bool) {
	value, exists := data[field]
	if !exists || value == nil {
		return nil, false
	}
	if reflect.TypeOf(value) != reflect.TypeOf(expected_type) {
		return nil, false
	}
	return value, true
}

func makeStringMap(data map[string]any) map[string]string {
	credentials := make(map[string]string)
	for k, v := range data {
		str_val, ok := v.(string)
		if !ok {
			log_handler.Error("Expected string in credentials: type = \"%T\" for key = \"%s\"", v, k)
			os.Exit(1)
		}
		credentials[k] = str_val
	}
	return credentials
}
