// Package snmpmodules is a declation for SNMP working devices management
package snmpmodules

import "fmt"

type ReachableDevice struct {
	IP          string            `json:"ip_address"`
	Hostname    string            `json:"hostname"`
	Description string            `json:"description"`
	SNMPVersion string            `json:"SNMP_version"`
	SNMPCred    map[string]string `json:"SNMP_cred"`
}

func CreateCredMapFromCredLength(credSlice []string) (map[string]string, error) {
	credLength := len(credSlice)
	switch credLength {
	case 1:
		return map[string]string{
			"community": credSlice[0],
		}, nil
	case 5:
		return map[string]string{
			"username":     credSlice[0],
			"authProtocol": credSlice[1],
			"authPassword": credSlice[2],
			"privProtocol": credSlice[3],
			"privPassword": credSlice[4],
		}, nil
	default:
		return nil, fmt.Errorf("credential length not recognized for %v", credSlice)
	}
}

func New(ip, hostname, description, version string, credMap map[string]string) ReachableDevice {
	return ReachableDevice{
		IP:          ip,
		Hostname:    hostname,
		Description: description,
		SNMPVersion: version,
		SNMPCred:    credMap,
	}
}
