package main

import (
	"data_treatment/config_handler"
	"data_treatment/db_handler"
	"data_treatment/log_handler"
	"data_treatment/serial_handler"
	"data_treatment/shared"
	"os"
	"strings"

	"github.com/tarm/serial"
)

func main() {
	log_path_raw, err := config_handler.GetData(shared.Consts, "DEFAULT_LOG_PATH", config_handler.TYPE_STRING)
	if err != nil {
		log_handler.Fatal("DEFAULT_LOG_PATH must be string")
	}
	log_path := log_path_raw.(string)
	log_filename_raw, err := config_handler.GetData(shared.Consts, "LOG_FILENAME", config_handler.TYPE_STRING)
	if err != nil {
		log_handler.Fatal("LOG_FILENAME must be string")
	}
	log_filename := log_filename_raw.(string)
	log_file_ext_raw, err := config_handler.GetData(shared.Consts, "LOG_FILE_EXT", config_handler.TYPE_STRING)
	if err != nil {
		log_handler.Fatal("LOG_FILE_EXT must be string")
	}
	log_ext := log_file_ext_raw.(string)

	log_handler.CreateAsyncLogger(log_path, log_filename, log_ext)
	if len(os.Args) < 2 {
		log_handler.Fatal("You must specify the serial port")
	}
	conn := strings.TrimSpace(os.Args[1])
	c := &serial.Config{Name: conn, Baud: 9600}
	port, err := serial.OpenPort(c)
	if err != nil {
		log_handler.Fatal("Cannot open serial port: %s", err)
	}
	defer port.Close()
	log_handler.Success("Serial port = [%s]", conn)
	for {
		data, err := serial_handler.ReceiveDataFromPort(port)
		if err != nil {
			log_handler.Error("Error receiving data: %s", err)
			continue
		}
		db_handler.JsonToData(data)
	}
}
