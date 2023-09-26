package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"

	"endgame/utils"

	log "github.com/sirupsen/logrus"
)

func main() {

	log.SetFormatter(&log.TextFormatter{TimestampFormat: "2006-01-02 15:04:05", FullTimestamp: true})

	newLog := log.WithFields(log.Fields{
		"name": "main.go",
	})

	isLocal := flag.Bool("local", false, "if running in local")
	flag.Parse()
	configuration := utils.ScanData{}

	if !*isLocal {
		utils.LoadScanData(&configuration, newLog)

	} else {
		scanId, ok := os.LookupEnv("SCAN_ID")
		if !ok {
			newLog.Panic("No scan id found exiting")
		}

		newLog.Infof("Loading scan data configuration for scan id -> %s", scanId)
		configFilePath := "config.json"

		file, _ := os.Open(configFilePath)
		decoder := json.NewDecoder(file)

		err := decoder.Decode(&configuration)
		if err != nil {
			newLog.Errorf("Please create `config.json` file in proper format")
			panic(fmt.Sprintf("Error occurred while reading config -> %s", err.Error()))
		}
		file.Close()
		configuration.Meta.ScanId = scanId
	}

	newLog.Info("Validating configuration loaded")
	err := utils.ValidateScanData(&configuration)
	if err != nil {
		panic(fmt.Sprintf("Invalid configuration found -> %s", err.Error()))
	}
	newLog.Info("Configuration loaded successfully, sending start scan request on webhook")

	defer func(configuration *utils.ScanData) {
		if err := recover(); err != nil {
			newLog.Errorf("Panic occurred in main thread -> %s", err)
			utils.SendRequestToSlack(configuration, newLog, err)
		}
		utils.SendCompleteScanRequest(configuration, newLog)
	}(&configuration)

	go utils.SendRequestToWebhook(&configuration, newLog, "scan.started", []byte(`{"reason":"Scan Started successfully"}`))
	go utils.SendStartScanRequest(&configuration, newLog)
	go utils.SendHealthWebhook(&configuration, newLog)

	err = StartScansInRoutine(&configuration)

	if err != nil {
		panic(fmt.Sprintf("Error occurred while running scans -> %s", err.Error()))
	}
	newLog.Info("All scans completed successfully exiting")
}
