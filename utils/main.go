package utils

import (
	"os"

	log "github.com/sirupsen/logrus"
)

// This function will be used for loading all the required scan data from environment variables
func LoadScanData(scanData *ScanData, newLog *log.Entry) {
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

	newLog.Info("All data loaded successfully")
}

// This function can be used for validating the scan data loaded
func ValidateScanData(scanData *ScanData) error {
	return nil
}
