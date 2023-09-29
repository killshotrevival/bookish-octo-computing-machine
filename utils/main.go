package utils

import (
	"os"
	"strconv"
	"time"

	log "github.com/sirupsen/logrus"
)

// This function will be used for loading all the required scan data from environment variables
func LoadScanData(scanData *ScanData, newLog *log.Entry) {
	var err error
	newLog.Info("Starting scan data loading")
	value, ok := os.LookupEnv("SCAN_AUDIT_ID")
	if ok {
		scanData.Meta.AuditId = value
	} else {
		newLog.Panic("Audit Id env not present")
	}

	value, ok = os.LookupEnv("SCAN_JOB_ID")
	if ok {
		scanData.Meta.JobId = value
	} else {
		newLog.Panic("Job Id env not present")
	}

	value, ok = os.LookupEnv("SCAN_WEBHOOK_TOKEN")
	if ok {
		scanData.Meta.WebhookToken = value
	} else {
		newLog.Panic("Webhook token env not present")
	}

	value, ok = os.LookupEnv("SCAN_WEBHOOK_URL")
	if ok {
		scanData.Meta.WebhookUrl = value
	} else {
		newLog.Panic("Webhook url env not present")
	}

	value, ok = os.LookupEnv("DAST_API_SVC_NAME")
	if !ok {
		newLog.Panic("Support server env not present")
	} else {
		scanData.ApiService = value
	}

	value, ok = os.LookupEnv("SCAN_ID")
	if !ok {
		newLog.Panic("No scan id found exiting")
	} else {
		scanData.Meta.ScanId = value
	}

	value, ok = os.LookupEnv("SCAN_TARGET")
	if !ok {
		newLog.Panic("No scan target found exiting")
	} else {
		scanData.Context.Target = value
	}

	value, ok = os.LookupEnv("MAX_SCAN_DURATION")
	if !ok {
		scanData.Meta.MaxScanDuration = 18000 // 5 hours in seconds
	} else {
		scanData.Meta.MaxScanDuration, err = strconv.ParseInt(value, 10, 64)

		if err != nil {
			newLog.Panic("Could not convert max scan duration to integer")
		}
	}

	value, ok = os.LookupEnv("SCAN_SCOPE_COVERAGE")
	if !ok {
		scanData.Context.ScanScopeCoverage = "full_domain"
	} else {
		scanData.Context.ScanScopeCoverage = value
	}

	newLog.Info("All data loaded successfully")
}

// This function can be used for validating the scan data loaded
func ValidateScanData(scanData *ScanData) error {
	return nil
}

// This function can be used for sending health check alert on webhook every 10 seconds.
func SendHealthWebhook(configuration *ScanData, newLog *log.Entry) {
	newLog.Info("Starting health check webhook go routine")
	uptimeTicker := time.NewTicker(10 * time.Second)

	for {
		newLog.Info("Sending health check webhook")
		<-uptimeTicker.C
		SendRequestToWebhook(configuration, newLog, "scan.health", []byte(`{"reason":"Alive and healthy"}`))
	}
}
