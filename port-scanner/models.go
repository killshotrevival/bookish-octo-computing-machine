package portscanner

import "encoding/json"

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
