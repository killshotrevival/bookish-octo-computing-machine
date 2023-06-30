package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	log "github.com/sirupsen/logrus"
)

// This function can be used for sending a request on the webhook
func SendRequestToWebhook(scanData *ScanData, newLog *log.Entry, event string, context []byte) error {
	newLog.Infof("Sending request for event -> %s", event)
	tempRequest := WebhookRequest{}

	tempRequest.Meta = RequestMeta{
		AuditId:      scanData.Meta.AuditId,
		JobId:        scanData.Meta.JobId,
		WebhookToken: scanData.Meta.WebhookToken,
		ScanId:       scanData.Meta.ScanId,
		Event:        event,
		Hostname:     "K8s",
	}
	tempRequest.Context = context

	postBody, _ := json.Marshal(tempRequest)
	responseBody := bytes.NewBuffer(postBody)

	resp, err := http.Post(scanData.Meta.WebhookUrl, "application/json", responseBody)

	newLog.Infof("Response status received -> %s for alert", resp.Status)

	if err != nil {
		newLog.Errorf("Error occurred while sending request on webhook -> %s", err.Error())
		return err
	}

	return nil
}

// This function can be used for triggering an alert message to slack
func SendRequestToSlack(scanData *ScanData, newLog *log.Entry, errorString interface{}) error {
	newLog.Info("Sending alert on slack")

	erroMessageString := fmt.Sprintf(`:pray-intensifies: *Error occurred inside endgame pod*
Error Message: %s
Scan Id: %s`, errorString, scanData.Meta.ScanId)
	tempSlackRequestData := SlackRequestData{Text: erroMessageString}

	tempSlackRequest := SlackRequest{JsonBlock: tempSlackRequestData}

	postBody, err := json.Marshal(tempSlackRequest)
	if err != nil {
		newLog.Errorf("Error occurred while marshalling the slack request body -> %s", err.Error())
		return err
	}
	responseBody := bytes.NewBuffer(postBody)

	resp, err := http.Post("http://"+scanData.ApiService+"/api/dast/slack_request", "application/json", responseBody)

	newLog.Infof("Response status received -> %s for slack alert", resp.Status)

	if err != nil {
		newLog.Errorf("Error occurred while sending alert on slack > %s", err.Error())
		return err
	}

	return nil
}

// This function can be used for triggering scan complete request
func SendCompleteScanRequest(scanData *ScanData, newLog *log.Entry) {
	newLog.Info("Sending complete scan request")
	client := &http.Client{}

	req, err := http.NewRequest("DELETE", "http://"+scanData.ApiService+"/api/endgame/"+scanData.Meta.ScanId, nil)
	if err != nil {
		panic(fmt.Sprintf("Error occurred while creating complete scan request -> %s", err.Error()))
	}
	newLog.Info("Request created successfully")
	resp, err := client.Do(req)
	newLog.Infof("Complete scan request status received -> %s", resp.Status)
	if err != nil {
		panic(fmt.Sprintf("Error occurred while sending complete scan request -> %s", err.Error()))
	}
	defer resp.Body.Close()
}
