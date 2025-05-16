package shared

import (
	"data_treatment/config_handler"
	"log"
)

var Consts map[string]any

func init() {
	consts, err := config_handler.LoadConsts(config_handler.DEFAULT_CONFIG_PATH)
	if err != nil {
		log.Fatalf("Error while loading the configuration file: %s", err)
	}
	Consts = consts
}
