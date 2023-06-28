package portscanner

import (
	"encoding/json"
	"endgame/utils"

	log "github.com/sirupsen/logrus"
)

func StartScan(scanData utils.ScanData) error {

	newLog := log.WithFields(log.Fields{
		"name": "port scanner",
	})
	newLog.Info("Starting new port scanner")

	newLog.Info("Raising alert for high severity port")
	newAlertBody := AlertBody{
		Name:        "Target Has Ports Open for Critical Services",
		Description: "The security assessment has identified that the target has open ports for critical services. It is recommended that the client reviews all open ports and takes necessary steps to secure them.",
		Solution:    "In particular, any unused ports should be shut down to reduce the attack surface and minimize the potential risk of unauthorized access or exploitation. Proper port management plays a crucial role in maintaining a secure network environment.",
		Evidence:    "These ports were found to be open: 22, 21",
		Risk:        "Medium",
		Confidence:  "High",
		AlertRef:    "portscanner_9027-1",
		PluginId:    "9027",
		Id:          1,
		AuditPhase:  "tool",
	}

	newAlertContext := AlertContext{
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
