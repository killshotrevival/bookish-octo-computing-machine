package portscanner

import (
	"encoding/json"
	"endgame/utils"
	"net/url"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"
)

// This function can be used for starting a port scanner on the context.target endpoint
func StartScan(scanData utils.ScanData) error {

	newLog := log.WithFields(log.Fields{
		"name": "port scanner",
	})
	newLog.Info("Starting new port scanner")

	parsedURL, err := url.Parse(scanData.Context.Target)
	if err != nil {
		panic(err)
	}
	newLog.Infof("Domain name found -> %s", parsedURL.Host)

	cmd := exec.Command("rustscan", "-a", parsedURL.Host, "-r", "1-65535", "-u", "5000", "-g")
	stdout, err := cmd.Output()

	if err != nil {
		newLog.Panicf("Error occurred while rust scan -> %s", err.Error())
		return err
	}

	newLog.Infof("Response found -> %s", stdout)

	re := regexp.MustCompile(`\[([\d,]+)\]`)
	portFoundString := re.FindStringSubmatch(string(stdout))

	portFoundArray := strings.Split(portFoundString[1], ",")

	if len(portFoundArray) < 1 {
		newLog.Info("No port found for scanning")
	}

	newLog.Info("RustScan might have overloaded the server, let it rest.")
	time.Sleep(10 * time.Second)

	err = serviceDetection(&portFoundArray, newLog)
	if err != nil {
		newLog.Errorf("Error occurred while service detection -> %s", err.Error())
		return err
	}

	return nil
}

// This function will be used for doing service detection on a specific port
func serviceDetection(portFoundArray *[]string, newLog *log.Entry) error {
	newLog.Info("Starting service detection")

	file, _ := os.Open("resource/portscanner_service_probes.json")
	decoder := json.NewDecoder(file)

	configuration := []PortScannerProbe{}
	err := decoder.Decode(&configuration)
	if err != nil {
		newLog.Panicf("Error occurred while reading config -> %s", err.Error())
		newLog.Panicf("Please create `config.json` file in proper format")
		log.Exit(1)
	}
	file.Close()

	newLog.Infof("Scanner Probe loaded -> %d", len(configuration))

	for _, port := range *portFoundArray {
		newLog.Infof("Checking for port -> %s", port)
		err := checkForPort(port, newLog)
		if err != nil {
			newLog.Errorf("Error occurred while further examining port -> %s", err.Error())
		}
	}

	return nil
}

func checkForPort(port string, newLog *log.Entry) error {
	for _, excludedPort := range ExcludedList {
		if port == excludedPort {
			newLog.Infof("%s is in excluded list, continuing without service detection.", port)
			return nil
		}
	}

	newLog.Info("Running NULL probe test...")
	return nil

}

func bannerMatcher(host string, addr string, port int, probes []PortScannerProbe, phase string, newLog *log.Entry) error {
	var rarityInt int
	var err error
	for _, probe := range probes {
		newLog.Infof("Current probe : %s", probe.Probe.ProbeName)
		rarityInt, err = strconv.Atoi(probe.Rarity.Rarity)
		if err != nil {
			newLog.Errorf("Error occurred converting rarity to int -> %s", err.Error())
			continue
		}

		if probe.Rarity.Rarity != "" && rarityInt > 5 && phase == "Excluded" {
			return nil
		}

	}

	return nil
}

// This function can be used for raising alerts in port scanning
func raiseAlert(scanData utils.ScanData, newLog *log.Entry) error {
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
