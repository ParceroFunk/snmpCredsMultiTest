// Package snmpmanager is an snmp implementation for quick testing based on
// multiple SNMP credentials, only for default port UDP 161
package snmpmanager

import (
	"fmt"
	// "log"
	"time"

	"github.com/gosnmp/gosnmp"
)

// constants for SNMP testing
const (
	SysName  string = "1.3.6.1.2.1.1.5.0"
	SysDescr string = "1.3.6.1.2.1.1.1.0"
)

// buildSNMPParams builds a *gosnmp.GoSNMP configured for either SNMPv2c or
// SNMPv3, based on version and the remaining credential fields.
//
// Expected cred layouts (after stripping the version field):
//
//	2c: [community]
//	3 : [username, authProtocol, authPassphrase, privProtocol, privPassphrase]
func buildSNMPParams(ip, version string, cred []string) (*gosnmp.GoSNMP, error) {
	// basic testing option
	timeout := 5 * time.Second
	retries := 2

	switch version {

	case "2c":
		if len(cred) != 1 {
			return nil, fmt.Errorf("2c requires a community string")
		}
		return &gosnmp.GoSNMP{
			Target:    ip,
			Port:      161,
			Version:   gosnmp.Version2c,
			Community: cred[0],
			Timeout:   timeout,
			Retries:   retries,
		}, nil

	case "3":
		if len(cred) != 5 {
			return nil, fmt.Errorf("v3 requires username, authProtocol, authPassphrase, privProtocol, privPassphrase")
		}
		username := cred[0]
		authProtoStr := cred[1]
		authPass := cred[2]
		privProtoStr := cred[3]
		privPass := cred[4]

		var authProtocol gosnmp.SnmpV3AuthProtocol
		switch authProtoStr {
		case "SHA":
			authProtocol = gosnmp.SHA
		case "SHA256":
			authProtocol = gosnmp.SHA256
		case "SHA512":
			authProtocol = gosnmp.SHA512
		case "MD5":
			authProtocol = gosnmp.MD5
		default:
			return nil, fmt.Errorf("unsupported auth protocol: %s", authProtoStr)
		}

		var privProtocol gosnmp.SnmpV3PrivProtocol
		switch privProtoStr {
		case "AES":
			privProtocol = gosnmp.AES
		case "AES256":
			privProtocol = gosnmp.AES256
		case "DES":
			privProtocol = gosnmp.DES
		default:
			return nil, fmt.Errorf("unsupported privacy protocol: %s", privProtoStr)
		}

		return &gosnmp.GoSNMP{
			Target:        ip,
			Port:          161,
			Version:       gosnmp.Version3,
			SecurityModel: gosnmp.UserSecurityModel,
			MsgFlags:      gosnmp.AuthPriv,
			Timeout:       timeout,
			Retries:       retries,
			SecurityParameters: &gosnmp.UsmSecurityParameters{
				UserName:                 username,
				AuthenticationProtocol:   authProtocol,
				AuthenticationPassphrase: authPass,
				PrivacyProtocol:          privProtocol,
				PrivacyPassphrase:        privPass,
			},
		}, nil

	default:
		return nil, fmt.Errorf("unsupported SNMP version: %s", version)
	}
}

// TestCredential builds the params, connects, runs Get, and closes the
// connection before returning. Because Close() is deferred *inside* this
// function (not in main's loop body), the socket is released at the end
// of every single iteration instead of accumulating until main() exits.
func TestCredential(ip string, cred []string, oids []string) (*gosnmp.SnmpPacket, error) {
	params, err := buildSNMPParams(ip, cred[0], cred[1:])
	if err != nil {
		return nil, fmt.Errorf("building params: %w", err)
	}

	err = params.Connect()
	if err != nil {
		return nil, fmt.Errorf("connect: %w", err)
	}
	defer params.Conn.Close()

	result, err := params.Get(oids)
	if err != nil {
		return nil, fmt.Errorf("get: %w", err)
	}

	return result, nil
}
