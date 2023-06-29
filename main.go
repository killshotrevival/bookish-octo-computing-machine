package main

import (
	"encoding/json"
	"flag"
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
			newLog.Panicf("Error occurred while reading config -> %s", err.Error())
			newLog.Panicf("Please create `config.json` file in proper format")
			log.Exit(1)
		}
		file.Close()
		configuration.Meta.ScanId = scanId
	}

	newLog.Info("Validating configuration loaded")
	err := utils.ValidateScanData(&configuration)
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
