package portscanner

import (
	"encoding/json"
	"endgame/utils"
	"fmt"
	"os/exec"
	"regexp"
	"strconv"
	"strings"

	log "github.com/sirupsen/logrus"
)

// This function can be used for starting a port scanner on the context.target endpoint
func StartScan(scanData *utils.ScanData, host string) error {

	newLog := log.WithFields(log.Fields{
		"name": "port scanner",
	})
	newLog.Infof("Starting new port scanner on -> %s", host)

	////////////
	// As port scanning will run from subdomain takeover, there is no need to parse url  anymore.
	////////////

	// parsedURL, err := url.Parse(host)
	// if err != nil {
	// 	// panic(err)
	// 	return err
	// }

	cmd := exec.Command("furious", "-s", "connect", "-p", "1-65535", host)
	stdout, err := cmd.Output()
	if err != nil {
		newLog.Errorf("Error occurred while furious -> %s", err.Error())
		return err
	}

	furiousOutput := string(stdout)
	if strings.Contains(string(stdout), "no such host") {
		newLog.Errorf("Host %s was not found, can't run port scanner.", host)
		return nil
	}

	newLog.Infof("Response found -> %s", furiousOutput)

	// Parse port scan result to a map
	portScanResult, err := portScanResultToMap(furiousOutput, newLog)

	if err != nil {
		return err
	}

	// Get all high severity ports from the results
	highSeverityPorts := getHighSeverityPorts(portScanResult)

	fmt.Println(portScanResult)
	fmt.Println(highSeverityPorts)

	if len(portScanResult) > 0 {
		raiseAlerts(scanData, portScanResult, highSeverityPorts, newLog, host)
	}

	return nil
}

// This function iterates all ports and checks if any high severity port is exposed.
func getHighSeverityPorts(portScanResult map[int]string) (highSeverityPorts map[int]string) {

	highSeverityPorts = make(map[int]string)

	for port, service := range portScanResult {
		for _, highSeverityPort := range HighSeverityPorts {
			if port == highSeverityPort {
				highSeverityPorts[port] = service
			}
		}
	}

	return
}

func portScanResultToMap(furiousOutput string, newLog *log.Entry) (map[int]string, error) {

	portAndServiceRegex := regexp.MustCompile(`(\d+)/(?:tcp|udp)\s+OPEN\s+(\S+)\n`)
	portAndServiceMatches := portAndServiceRegex.FindAllStringSubmatch(furiousOutput, -1)

	portScanResult := make(map[int]string)

	for _, match := range portAndServiceMatches {
		port, err := strconv.Atoi(match[1])

		if err != nil {
			newLog.Errorf("Error occurred while converting port from string to int -> %s", err.Error())
			return make(map[int]string), err
		}

		service := match[2]
		if service == "http" || service == "https" {
			continue
		}
		portScanResult[port] = service
	}

	return portScanResult, nil
}

// This function will raise alerts using the alert details passed to it
func raiseAlert(scanData *utils.ScanData, name string, desc string, soln string, evid string, risk string, conf string, alertRef string, pluginId string, id int, auditPhase string, newLog *log.Entry) error {

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

	utils.SendRequestToWebhook(scanData, newLog, "alert", resp)

	return nil
}

// This function is to initialise alerts for portscanner service.
func raiseAlerts(scanData *utils.ScanData, portScanResult map[int]string, highSeverityPorts map[int]string, newLog *log.Entry, host string) error {

	newLog.Info("Raising low severity alert for detected ports.")
	portScanResultJSON, err := json.Marshal(portScanResult)

	if err != nil {
		newLog.Errorf("Error occurred while marshalling port scan result.")
	}

	var (
		id         int    = 1
		name       string = "[Recommendation] Review Open Ports"
		desc       string = "The security assessment has identified open ports on the target system. It is recommended to thoroughly review these open ports to ensure that only necessary services are accessible."
		soln       string = "Unnecessary or unused ports should be closed to reduce the attack surface and mitigate potential risks associated with unauthorized access or exploitation."
		evid       string = fmt.Sprintf("Port found on {%s} are %s", host, portScanResultJSON)
		risk       string = "Low"
		conf       string = "High"
		pluginId   string = "9027"
		alertRef   string = "portscanner_9027-1"
		auditPhase string = "tool"
	)

	raiseAlert(scanData, name, desc, soln, evid, risk, conf, alertRef, pluginId, id, auditPhase, newLog)

	if len(highSeverityPorts) < 1 {
		return nil
	}

	newLog.Info("Raising medium severity alert for critical ports.")
	highSeverityPortsJSON, err := json.Marshal(highSeverityPorts)

	if err != nil {
		newLog.Errorf("Error occurred while marshalling port scan result.")
	}

	id = 1
	name = "Target Has Ports Open for Critical Services"
	desc = "The security assessment has identified that the target has open ports for critical services. It is recommended that the client reviews all open ports and takes necessary steps to secure them."
	soln = "In particular, any unused ports should be shut down to reduce the attack surface and minimize the potential risk of unauthorized access or exploitation. Proper port management plays a crucial role in maintaining a secure network environment."
	evid = fmt.Sprintf("Port found on {%s} are %s", host, highSeverityPortsJSON)
	risk = "Medium"
	conf = "High"
	pluginId = "9027"
	alertRef = "portscanner_9027-2"
	auditPhase = "tool"

	raiseAlert(scanData, name, desc, soln, evid, risk, conf, alertRef, pluginId, id, auditPhase, newLog)

	return nil
}
