package portscanner

import (
	"encoding/json"
	"endgame/utils"
	"net"
	"net/url"
	"os"
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

	// cmd := exec.Command("rustscan", "-a", parsedURL.Host, "-r", "1-65535", "-u", "5000", "-g")
	// stdout, err := cmd.Output()

	// if err != nil {
	// 	newLog.Panicf("Error occurred while rust scan -> %s", err.Error())
	// 	return err
	// }

	stdout := "103.48.51.180 -> [21,80,53,110,143,443,587,1232,3306,8172,8443,8880,49670,64000]"

	newLog.Infof("Response found -> %s", stdout)

	re := regexp.MustCompile(`\[([\d,]+)\]`)
	portFoundString := re.FindStringSubmatch(string(stdout))

	portFoundArray := strings.Split(portFoundString[1], ",")

	if len(portFoundArray) < 1 {
		newLog.Info("No port found for scanning")
	}

	// newLog.Info("RustScan might have overloaded the server, let it rest.")
	// time.Sleep(10 * time.Second)

	err = serviceDetection(parsedURL.Host, &portFoundArray, newLog)
	if err != nil {
		newLog.Errorf("Error occurred while service detection -> %s", err.Error())
		return err
	}

	return nil
}

// This function will be used for doing service detection on a specific port
func serviceDetection(host string, portFoundArray *[]string, newLog *log.Entry) error {
	newLog.Info("Starting service detection")

	file, _ := os.Open("resource/portscanner_service_probes.json")
	decoder := json.NewDecoder(file)

	portScannerProbeList := []PortScannerProbe{}
	err := decoder.Decode(&portScannerProbeList)
	if err != nil {
		newLog.Panicf("Error occurred while reading config -> %s", err.Error())
		newLog.Panicf("Please create `config.json` file in proper format")
		log.Exit(1)
	}
	file.Close()

	newLog.Infof("Scanner Probe loaded -> %d", len(portScannerProbeList))

	for _, port := range *portFoundArray {
		newLog.Infof("Checking for port -> %s", port)
		err := checkForPort(host, port, &portScannerProbeList, newLog)
		if err != nil {
			newLog.Errorf("Error occurred while further examining port -> %s", err.Error())
		}
	}

	return nil
}

func checkForPort(host string, port string, portScannerProbeList *[]PortScannerProbe, newLog *log.Entry) error {

	for _, excludedPort := range ExcludedList {
		if port == excludedPort {
			newLog.Infof("%s is in excluded list, continuing without service detection.", port)
			return nil
		}
	}

	newLog.Info("Running NULL probe test...")
	addr, err := net.LookupIP(host)
	if err != nil {
		newLog.Errorf("Error occurred while fetching IP from hostname -> %s", err.Error())
		return err
	}
	bannerMatcher(host, addr[0].String(), port, []PortScannerProbe{(*portScannerProbeList)[0]}, "NULL", newLog)
	return nil

}

func bannerMatcher(host string, addr string, port string, probes []PortScannerProbe, phase string, newLog *log.Entry) error {
	var rarityInt, timeout int
	var err error
	var banner string
	for _, probe := range probes {
		newLog.Infof("Current probe : %s", probe.Probe.ProbeName)

		if probe.Rarity.Rarity != "" {
			rarityInt, err = strconv.Atoi(probe.Rarity.Rarity)
			if err != nil {
				newLog.Errorf("Error occurred converting rarity to int -> %s", err.Error())
				continue
			}
			if rarityInt > 5 && phase == "Excluded" {
				return nil
			}
		}

		if phase == "NULL" {
			timeout = 10
		} else if probe.TotalWaitMs.TotalWaitMs != "" {
			timeout, err = strconv.Atoi(probe.TotalWaitMs.TotalWaitMs)
			if err != nil {
				newLog.Errorf("Error occurred while converting total wait ms to int -> %s", err.Error())
				continue
			}

			timeout *= 1000
		} else {
			timeout = 6
		}
		banner, err = socketConnector(addr, port, probe.Probe.ProbeString, timeout, 0, newLog)
		if err != nil {
			continue
		}

		newLog.Infof("Banner received ->  %s", banner)
	}

	return nil
}

func socketConnector(addr string, port string, probeString string, timeout int, retry int, newLog *log.Entry) (string, error) {
	newLog.Infof("Connecting to %s with timeout %d", addr+":"+port, timeout)
	connection, err := net.DialTimeout("tcp", addr+":"+port, time.Duration(timeout)*time.Second)
	if err != nil {
		newLog.Errorf("Error occurred while connecting to server -> %s", err.Error())
		if err == os.ErrDeadlineExceeded {
			if retry > 0 {
				return socketConnector(addr, port, probeString, timeout+5, retry, newLog)
			} else {
				return "SOCKET_TIMEOUT_EXCEPTION", err
			}
		} else {
			return "SOCKET_EXCEPTION", err
		}
	}
	newLog.Info("Connection to server successful")

	connection.SetDeadline(time.Now().Add(time.Duration(timeout * int(time.Second))))
	defer connection.Close()

	newLog.Info("Sending Data to server")
	_, err = connection.Write([]byte(probeString))
	if err != nil {
		newLog.Errorf("Error occurred while sending data to server -> %s", err.Error())
		if err == os.ErrDeadlineExceeded {
			if retry > 0 {
				return socketConnector(addr, port, probeString, timeout+5, retry-1, newLog)
			} else {
				return "SOCKET_TIMEOUT_EXCEPTION", err
			}
		} else {
			return "SOCKET_EXCEPTION", err
		}
	}

	newLog.Info("Fetching data from server")
	buffer := make([]byte, 1024)
	mLen, err := connection.Read(buffer)
	if err != nil {
		newLog.Errorf("Error occurred while reading data from server -> %s", err.Error())
		if err == os.ErrDeadlineExceeded {
			if retry > 0 {
				return socketConnector(addr, port, probeString, timeout+5, retry-1, newLog)
			} else {
				return "SOCKET_TIMEOUT_EXCEPTION", err
			}
		} else {
			return "SOCKET_EXCEPTION", err
		}
	}

	return strings.Trim(string(buffer[:mLen]), "\r\n"), nil
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
