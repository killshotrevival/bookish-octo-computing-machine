package utils

import "encoding/json"

type ScanData struct {
	Meta    Meta    `json:"meta"`
	Context Context `json:"context"`
}

type WebhookRequest struct {
	Meta    RequestMeta     `json:"meta"`
	Context json.RawMessage `json:"context"`
}

type RequestMeta struct {
	Event        string `json:"event"`
	AuditId      string `json:"auditId"`
	JobId        string `json:"jobId"`
	WebhookToken string `json:"webhookToken"`
	ScanId       string `json:"scanId"`
	Hostname     string `json:"hostname"`
}

type Meta struct {
	AuditId      string `json:"auditId"`
	JobId        string `json:"jobId"`
	WebhookUrl   string `json:"webhookUrl"`
	WebhookToken string `json:"webhookToken"`
	ScanId       string `json:"scanId"`
}

type Context struct {
	Target string `json:"target"`
}

func ValidateScanData(scanData *ScanData) error {
	return nil
}
