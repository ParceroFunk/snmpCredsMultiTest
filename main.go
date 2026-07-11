package main

import (
	"fmt"
	"log"
	"strings"

	"github.com/ParceroFunk/snmpCredsMultiTest/config"
	"github.com/ParceroFunk/snmpCredsMultiTest/filemanager"
	"github.com/ParceroFunk/snmpCredsMultiTest/snmpmanager"
	"github.com/ParceroFunk/snmpCredsMultiTest/snmpmodules"
)

func main() {
	// Parse inputs from CLI flags
	cfg := config.Parse()

	log.Printf("Starting SNMP testing with multiple creds on an IP list")

	// Read devices from file with 1 IP per line
	fileMgr := filemanager.New(
		cfg.DeviceFile, // InputFilePath
		cfg.OutputFile, // OutputFilePath
	)
	deviceIPs, err := fileMgr.ReadLines() // []string with IPs
	if err != nil {
		panic(fmt.Errorf("failed to read devices IPs file: %v", err.Error()))
	}
	log.Printf("Found %d devices on %q file", len(deviceIPs), fileMgr.InputFilePath)

	// Read credentials from a file with 1 credential per line
	fileMgr.InputFilePath = cfg.CredsFile
	snmpCredsLines, err := fileMgr.ReadLines() // []string with SNMP creds
	if err != nil {
		panic(fmt.Errorf("failed to read SNMP credential file: %v", err.Error()))
	}
	log.Printf("Found %d SNMP credentials on %q file", len(snmpCredsLines), fileMgr.InputFilePath)
	// Devide credentials by space for SNMP struct creation
	var snmpCreds [][]string
	for _, line := range snmpCredsLines {
		snmpCreds = append(snmpCreds, strings.Split(line, " "))
	}

	oids := []string{snmpmanager.SysName, snmpmanager.SysDescr}

	// start a slice of reachables
	var reachables []snmpmodules.ReachableDevice

	// iterate over each IP and each credential to test if one is working
IPLoop:
	for _, ip := range deviceIPs {
		for _, cred := range snmpCreds {
			result, err := snmpmanager.TestCredential(ip, cred, oids)
			if err != nil {
				log.Printf("[%s] cred %v failed: %v", ip, cred, err)
				continue // keep testing other creds/IPs instead of Fatalf-ing out
			}

			fmt.Printf("[%s] cred %v succeeded:\n", ip, cred)

			// Add successfull devices to a slice of Reachables struct
			credMap, err := snmpmodules.CreateCredMapFromCredLength(cred[1:])
			if err != nil {
				log.Printf("[%s] cred %v failed to parse for saving: %v", ip, cred, err)
			}

			// Check OID result for reachable struct value constructor
			var response [2]string
			for index, variable := range result.Variables {
				switch value := variable.Value.(type) {
				case []byte:
					response[index] = string(value)
				case string:
					response[index] = value
				default:
					log.Printf("[%s] response for %v failed to parse for saving. unkown type", ip, variable.Name)
				}
			}

			// verify result for correct value alocation
			reachables = append(reachables, snmpmodules.New(ip, cred[0], response[0], response[1], credMap))

			// Continue with next IP when credential worked and was saved successfully
			continue IPLoop
		}
	}

	fmt.Println(reachables)
	err = fileMgr.WriteResult(reachables)
	if err != nil {
		log.Printf("Failed to write results to %v file: %v", fileMgr.OutputFilePath, err.Error())
	}
}
