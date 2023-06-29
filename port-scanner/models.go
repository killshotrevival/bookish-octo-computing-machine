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

type PortScannerProbe struct {
	Probe        PortScannerProbeProbe        `json:"probe"`
	Rarity       PortScannerProbeRarity       `json:"rarity"`
	TotalWaitMs  PortScannerProbeTotalWaitMs  `json:"totalwaitms"`
	TcpWrappedMs PortScannerProbeTcpWrappedMs `json:"tcpwrappedms"`
	Ports        PortScannerProbePorts        `json:"ports"`
	SslPorts     PortScannerProbeSslports     `json:"sslports"`
	Matches      []PortScannerProbeMatches    `json:"matches"`
	SoftMatches  []PortScannerProbeMatches    `json:"softmatches"`
}

type PortScannerProbeTotalWaitMs struct {
	TotalWaitMs string `json:"totalwaitms"`
}

type PortScannerProbeTcpWrappedMs struct {
	TcpWrappedMs string `json:"tcpwrappedms"`
}

type PortScannerProbeProbe struct {
	Protocol    string `json:"protocol"`
	ProbeName   string `json:"probename"`
	ProbeString string `json:"probestring"`
}
type PortScannerProbeRarity struct {
	Rarity string `json:"rarity"`
}

type PortScannerProbePorts struct {
	Ports string `json:"ports"`
}

type PortScannerProbeSslports struct {
	SslPorts string `json:"sslports"`
}

type PortScannerProbeMatches struct {
	Service         string `json:"service"`
	Pattern         string `json:"pattern"`
	PatternCompiled string `json:"pattern_compiled"`
	VersionInfo     string `json:"versioninfo"`
}
