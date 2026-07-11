package main

import (
	"fmt"
	"log"

	"github.com/ParceroFunk/snmpCredsMultiTest/snmpmanager"
	"github.com/ParceroFunk/snmpCredsMultiTest/snmpmodules"
)

func main() {
	fmt.Println("Starting SNMP testing with multiple creds on an IP list")

	// TODO: read these from files instead of hardcoding
	deviceIPs := []string{"192.168.1.100", "192.168.1.200"}
	snmpCreds := [][]string{
		{"2c", "public"},
		{"3", "username", "SHA", "authPass123", "AES", "privPass123"},
	}

	oids := []string{snmpmanager.SysName, snmpmanager.SysDescr}

	// start a slice of reachables
	var reachables []snmpmodules.ReachableDevice

	// iterate over each IP and each credential to test if one is working
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
			// Check result for reachable value constructor
			var response [2]string
			for index, r := range result.Variables {
				var ok bool
				response[index], ok = r.Value.(string)
				if !ok {
					log.Printf("[%s] response %v failed to parse for saving: %v", ip, r.Value, err)
				}
				fmt.Printf("  %s = %v\n", r.Name, r.Value)
			}
			// verify result for correct value alocation
			reachables = append(reachables, snmpmodules.New(ip, cred[0], response[0], response[1], credMap))
		}
	}

	fmt.Println(reachables)
}
