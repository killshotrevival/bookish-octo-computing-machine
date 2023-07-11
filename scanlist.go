package main

import (
	"endgame/utils"

	portscanner "endgame/port-scanner"
	subdomaintakeover "endgame/subdomain-takeover"
)

// This function can be used for generating a list of all audit functions to be executed
func InitializeAuditFunctions() ([]func(utils.ScanData) error, error) {
	var allAwsAuditFunctions []func(utils.ScanData) error

	allAwsAuditFunctions = append(allAwsAuditFunctions,
		portscanner.StartScan,
		subdomaintakeover.StartScan,
	)

	return allAwsAuditFunctions, nil
}
