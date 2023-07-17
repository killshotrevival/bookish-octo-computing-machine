package utils

import "encoding/json"

type ScanData struct {
	Meta       Meta    `json:"meta"`
	Context    Context `json:"context"`
	ApiService string  `json:"api_service_name"`
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
	Target            string `json:"target"`
	ScanScopeCoverage string `json:"scopeCoverage"`
}

type SlackRequest struct {
	JsonBlock SlackRequestData `json:"json_block"`
}

type SlackRequestData struct {
	Text string `json:"text"`
}

type AlertContext struct {
	Tags  json.RawMessage `json:"tags"`
	Alert AlertBody       `json:"alert"`
}

type AlertBody struct {
	Id          int    `json:"id"`
	PluginId    string `json:"pluginId"`
	AlertRef    string `json:"alertRef"`
	AuditPhase  string `json:"auditPhase"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Solution    string `json:"solution"`
	Evidence    string `json:"evidence"`
	Risk        string `json:"risk"`
	Confidence  string `json:"confidence"`
}
