package subdomaintakeover

import (
	"encoding/json"
	"endgame/utils"
	"io/ioutil"
	"net/url"
	"os/exec"
	"strings"

	"github.com/haccer/subjack/subjack"
	log "github.com/sirupsen/logrus"
)

// This function can be used for starting subdomain takeover scan on the host of context.target endpoint
func StartScan(scanData utils.ScanData) error {

	newLog := log.WithFields(log.Fields{
		"name": "subdomain takeover",
	})
	newLog.Info("Starting subdomain takeover scan.")

	parsedURL, err := url.Parse(scanData.Context.Target)
	if err != nil {
		panic(err)
	}

	domain := parsedURL.Host

	if scanData.Context.ScanScopeCoverage == "full_domain" {

		domainArray := strings.Split(parsedURL.Host, ".")
		domain = domainArray[len(domainArray)-2] + "." + domainArray[len(domainArray)-1]
	}

	newLog.Infof("Finding subdomains for -> %s", domain)

	cmd := exec.Command("subfinder", "-d", domain, "-silent")
	stdout, err := cmd.Output()

	if err != nil {
		newLog.Panicf("Error occurred while running subfinder -> %s", err.Error())
		return err
	}

	subdomains := strings.Split(string(stdout), "\n")

	// subdomains = append(subdomains, strings.Split(string(stdout), "\n")...)
	vulnerableSubdomains := make(map[string]string)

	var fingerprints []subjack.Fingerprints
	config, _ := ioutil.ReadFile("/home/astra/Desktop/fingerprints.json")
	json.Unmarshal(config, &fingerprints)

	newLog.Infof("%d subdomains found", len(subdomains))
	newLog.Infof("subdomains found -> %s", subdomains)
	for i := 0; i < len(subdomains); i++ {
		subdomain := subdomains[i]
		newLog.Infof("%d. Testing subdomain takeover on %s", i, subdomain)
		service := subjack.Identify(subdomain, false, false, 10, fingerprints)

		if service != "" {
			newLog.Infof("\n[ALERT] Subdomain takeover possible on %s\n", subdomain)
			service = strings.ToLower(service)
			vulnerableSubdomains[subdomain] = service
		}
	}

	if len(vulnerableSubdomains) > 0 {
		raiseAlerts(scanData, vulnerableSubdomains, newLog)
	}

	return nil
}

// This function will raise alerts using the alert details passed to it
func raiseAlert(scanData utils.ScanData, name string, desc string, soln string, evid string, risk string, conf string, alertRef string, pluginId string, id int, auditPhase string, newLog *log.Entry) error {

	newAlertBody := utils.AlertBody{
		Name:        name,
		Description: desc,
		Solution:    soln,
		Evidence:    evid,
		Risk:        risk,
		Confidence:  conf,
		AlertRef:    alertRef,
		PluginId:    pluginId,
		Id:          id,
		AuditPhase:  auditPhase,
	}

	newAlertContext := utils.AlertContext{
		Alert: newAlertBody,
		Tags:  []byte(`{"fetchFromAlert": true}`),
	}

	resp, err := json.Marshal(newAlertContext)

	if err != nil {
		newLog.Errorf("Error occurred while marshalling alert context")
	}

	utils.SendRequestToWebhook(&scanData, newLog, "alert", resp)

	return nil
}

// This function is to initialise alerts for portscanner service.
func raiseAlerts(scanData utils.ScanData, vulnerableSubdomains map[string]string, newLog *log.Entry) error {

	newLog.Info("Raising high severity alerts for subdomain takeover.")
	vulnerableSubdomainsJSON, err := json.Marshal(vulnerableSubdomains)

	if err != nil {
		newLog.Errorf("Error occurred while marshalling subdomain takeover result.")
	}

	var (
		id         int    = 1
		name       string = "[CRITICAL] Subdomain Takeover Possible"
		desc       string = "A subdomain takeover vulnerability has been identified, indicating the possibility of an attacker gaining control over a subdomain that is no longer in use or improperly configured. Subdomain takeover occurs when an external entity is able to take control of a subdomain, potentially leading to malicious activities such as phishing, data theft, or unauthorized access."
		soln       string = "Properly configure DNS settings such as CNAME records and remove any unused or obsolete subdomains."
		evid       string = string(vulnerableSubdomainsJSON)
		risk       string = "High"
		conf       string = "High"
		pluginId   string = "1"
		alertRef   string = "subover_1"
		auditPhase string = "tool"
	)

	raiseAlert(scanData, name, desc, soln, evid, risk, conf, alertRef, pluginId, id, auditPhase, newLog)

	return nil
}
