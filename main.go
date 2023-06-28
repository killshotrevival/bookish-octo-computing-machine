package main

import (
	"encoding/json"
	"os"

	"endgame/utils"

	log "github.com/sirupsen/logrus"
)

func main() {

	newLog := log.WithFields(log.Fields{
		"name": "main.go",
	})

	scanId, ok := os.LookupEnv("SCAN_ID")

	if !ok {
		newLog.Panic("No scan id found exiting")
		log.Exit(1)
	}

	newLog.Infof("Loading scan data configuration for scan id -> %s", scanId)
	configFilePath := "config.json"

	file, _ := os.Open(configFilePath)
	decoder := json.NewDecoder(file)

	configuration := utils.ScanData{}
	err := decoder.Decode(&configuration)
	if err != nil {
		newLog.Panicf("Error occurred while reading config -> %s", err.Error())
		newLog.Panicf("Please create `config.json` file in proper format")
		log.Exit(1)
	}
	file.Close()

	configuration.Meta.ScanId = scanId

	newLog.Info("Validating configuration loaded")
	err = utils.ValidateScanData(&configuration)
	if err != nil {
		newLog.Panicf("Invalid configuration found -> %s", err.Error())
		log.Exit(1)
	}
	newLog.Info("Configuration loaded successfully, sending start scan request on webhook")

	utils.SendRequestToWebhook(&configuration, newLog, "scan.started", []byte(`{"reason":"Scan Started successfully"}`))

	err = StartScansInRoutine(&configuration)

	if err != nil {
		newLog.Panicf("Error occurred while running scans -> %s", err.Error())
		log.Exit(1)
	}

	newLog.Info("All scans completed successfully exiting")
}
