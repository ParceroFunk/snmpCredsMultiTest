package main

import (
	"log"
	"strings"

	"github.com/ParceroFunk/snmpCredsMultiTest/config"
	"github.com/ParceroFunk/snmpCredsMultiTest/discovery"
	"github.com/ParceroFunk/snmpCredsMultiTest/filemanager"
	"github.com/ParceroFunk/snmpCredsMultiTest/snmpmodules"
	"github.com/ParceroFunk/snmpCredsMultiTest/utils"
)

func main() {
	// Parse inputs from CLI flags
	cfg := config.Parse()

	log.Printf("Starting SNMP testing with multiple creds on an IP list")

	fileMgr := filemanager.New(cfg.DeviceFile, cfg.OutputFile)

	// Create slice of IPs and slice of creds from config data
	deviceIPs, snmpCreds := getTestInputs(&fileMgr, cfg)

	// loop over IP and creds for getting a []snmpmodules.ReachableDevice
	reachables := discovery.Run(deviceIPs, snmpCreds, cfg.MaxConcurrency)

	// save(&fileMgr, reachables)

	err := exportCSV(&fileMgr, reachables, "ip,sysName", false)
	if err != nil {
		log.Printf("Failed to write CSV export: %v", err)
	}
}

func getTestInputs(fileMgr *filemanager.FileManager, cfg config.Config) ([]string, [][]string) {
	// Read devices from file with 1 IP per line
	deviceIPs, err := fileMgr.ReadLines() // []string with IPs
	if err != nil {
		log.Fatalf("failed to read devices IPs file: %v", err)
	}
	log.Printf("Found %d devices on %q file", len(deviceIPs), fileMgr.InputFilePath)

	// Read credentials from a file with 1 credential per line
	fileMgr.InputFilePath = cfg.CredsFile
	snmpCredsLines, err := fileMgr.ReadLines() // []string with SNMP creds
	if err != nil {
		log.Fatalf("failed to read SNMP credential file: %v", err)
	}
	log.Printf("Found %d SNMP credentials on %q file", len(snmpCredsLines), fileMgr.InputFilePath)

	// Divide credentials by space for SNMP struct creation
	var snmpCreds [][]string
	for _, line := range snmpCredsLines {
		snmpCreds = append(snmpCreds, strings.Split(line, " "))
	}
	return deviceIPs, snmpCreds
}

// exportCSV opens path via fileMgr and writes data to it as CSV, including
// only the comma-separated fields, in order.
func exportCSV(fileMgr *filemanager.FileManager, data []snmpmodules.ReachableDevice, fields string, header bool) error {
	w, err := fileMgr.OpenWriter()
	if err != nil {
		return err
	}
	defer w.Close()

	keys := strings.Split(fields, ",")
	return utils.ToCSV(w, data, keys, header)
}
