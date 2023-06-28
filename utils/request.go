package utils

import (
	"bytes"
	"encoding/json"
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

	newLog.Infof("Request status received -> %s for alert", resp.Status)

	if err != nil {
		newLog.Errorf("Error occurred while sending request on webhook -> %s", err.Error())
		return err
	}

	return nil
}
