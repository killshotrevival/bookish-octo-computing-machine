package main

import (
	"endgame/utils"
	"sync"

	log "github.com/sirupsen/logrus"
)

// This function will call all the scan's declared under `scanlist.go` file via go-routine
func StartScansInRoutine(scanData *utils.ScanData) error {
	newLog := log.WithFields(log.Fields{
		"name": "scanner.go",
	})
	allAwsAuditFunctions, err := InitializeAuditFunctions()

	if err != nil {
		newLog.Errorf("Error occurred while generating aws audit functions list -> %s", err.Error())
		return err
	}

	newLog.Infof("Audits list generated successfully, will be running %d audits on the target", len(allAwsAuditFunctions))

	var wg = sync.WaitGroup{}
	maxGoroutines := 10
	guard := make(chan struct{}, maxGoroutines)

	for i := 0; i < len(allAwsAuditFunctions); i++ {
		guard <- struct{}{}
		wg.Add(1)
		go func(n func(utils.ScanData) error) {
			n(*scanData)
			<-guard
			wg.Done()
		}(allAwsAuditFunctions[i])
	}

	wg.Wait()

	return nil
}
