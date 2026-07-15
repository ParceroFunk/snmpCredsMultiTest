package main

import (
	"log"
	"regexp"
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

	// convert reachable devices to CSV with "ip,hostname", exclude headers
	csvData, err := getCSV(reachables, "ip_address,hostname,description", false)
	if err != nil {
		log.Printf("Failed to write CSV export: %v", err)
	}

	// filter by vendor
	var vendorData [][]string
	vendorRegEx := "[cC]isco"
	re := regexp.MustCompile(vendorRegEx)
	for _, row := range csvData {
		if re.MatchString(row[2]) {
			vendorData = append(vendorData, row)
		}
	}

	// save resulting CSV data into output file
	save(&fileMgr, vendorData)
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

// getCSV takes JSON-like data of reachable struct and returns a [][]string,
// gettin it ready for a CSV file writting. Helper function for case this project.
func getCSV(data []snmpmodules.ReachableDevice, fields string, header bool) ([][]string, error) {
	keys := strings.Split(fields, ",")
	csvData, err := utils.ToCSV(data, keys, header)
	return csvData, err
}

func save(fileMgr *filemanager.FileManager, data [][]string) {
	err := fileMgr.CSVWriteResult(data)
	if err != nil {
		log.Printf("Failed to write results to %v file: %v", fileMgr.OutputFilePath, err)
	}
}
