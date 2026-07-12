package main

import (
	"log"
	"strings"

	"github.com/ParceroFunk/snmpCredsMultiTest/config"
	"github.com/ParceroFunk/snmpCredsMultiTest/discovery"
	"github.com/ParceroFunk/snmpCredsMultiTest/filemanager"
	"github.com/ParceroFunk/snmpCredsMultiTest/snmpmodules"
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

	save(&fileMgr, reachables)
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

func save(fileMgr *filemanager.FileManager, data []snmpmodules.ReachableDevice) {
	err := fileMgr.WriteResult(data)
	if err != nil {
		log.Printf("Failed to write results to %v file: %v", fileMgr.OutputFilePath, err)
	}
}
