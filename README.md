# snmpCredsMultiTest

This is GoLang CLI tool for testing multiple SNMP credentials (independently from
their version) on list of multiple IP address, and writes out which
IP/credential pairs are reachable.

For each IP, credentials are tried in order until one succeeds; the tool then
moves on to the next IP. Results are written to a JSON file.

## Requirements

- Go 1.21+ (adjust to your actual `go.mod` version)
- [gosnmp](https://github.com/gosnmp/gosnmp)

## Installation

```bash
git clone https://github.com/ParceroFunk/snmpCredsMultiTest.git
cd snmpCredsMultiTest
go build -o snmpCredsMultiTest .
```

## Usage

```bash
./snmpCredsMultiTest [flags]
```

### Flags

| Flag             | Default            | Description                                   |
|------------------|---------------------|------------------------------------------------|
| `-d`, `-devices`  | `devices.txt`       | Input file with the list of device IPs (one per line) |
| `-c`, `-creds`    | `creds.txt`         | Input file with SNMP credentials (one per line) |
| `-o`, `-output`   | `reachables.json`   | Output JSON file with reachable device data     |
| `-h`, `-help`     | —                    | Show usage and exit                             |

Example:

```bash
./snmpCredsMultiTest -d hosts.txt -c snmp_creds.txt -o results.json
```

Run `./snmpCredsMultiTest -h` at any time to see this same list from the CLI.

## Input file formats

### Devices file (`-devices`)

One IP address per line:

```
192.168.1.100
192.168.1.200
```

### Credentials file (`-creds`)

One credential per line, space-separated. Fields depend on SNMP version:

**SNMPv2c** — `2c <community>`

```
2c public
```

**SNMPv3** — `3 <username> <authProtocol> <authPassphrase> <privProtocol> <privPassphrase>`

```
3 username SHA authPass123 AES privPass123
```

Supported auth protocols: `MD5`, `SHA`, `SHA256`
Supported privacy protocols: `DES`, `AES`

## Output

Reachable devices are written as JSON to the `-output` file, containing the
IP, the SNMP version/credential that worked, and the queried `sysName` /
`sysDescr` values for that device.

## Project structure

```
.
├── main.go              # entry point: wires config, filemanager, and snmpmanager together
├── config/               # CLI flag parsing (-devices/-creds/-output) and usage text
├── filemanager/          # reading input files line-by-line, writing JSON results
├── snmpmanager/          # builds gosnmp.GoSNMP params per credential, runs Connect/Get
└── snmpmodules/          # ReachableDevice type and credential-map helpers
```

## Notes

- Failed credential/IP combinations are logged and skipped — the tool keeps
  testing the remaining credentials and devices rather than exiting on the
  first failure.
- Each SNMP connection is opened and closed within a single test attempt, so
  sockets aren't held open across the full run.
