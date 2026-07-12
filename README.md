# snmpCredsMultiTest

This is GoLang CLI tool for testing multiple SNMP credentials (independently from
their version) on list of multiple IP address, and writes out which
IP/credential pairs are reachable.

For each IP, credentials are tried in order until one succeeds; the tool
then moves on to the next credential/device pair. IPs are tested in
parallel, bounded by a configurable worker limit. Results are written to a
JSON file.

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

| Flag              | Default            | Description                                             |
|-------------------|---------------------|-----------------------------------------------------------|
| `-d`, `-devices`   | `devices.txt`       | Input file with the list of device IPs (one per line)     |
| `-c`, `-creds`     | `creds.txt`         | Input file with SNMP credentials (one per line)           |
| `-o`, `-output`    | `reachables.json`   | Output JSON file with reachable device data                |
| `-w`, `-workers`   | `20`                 | How many IPs are tested in parallel                        |
| `-h`, `-help`      | —                    | Show usage and exit                                         |

Example:

```bash
./snmpCredsMultiTest -d hosts.txt -c snmp_creds.txt -o results.json -w 50
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

## Concurrency

IPs are tested in parallel, one goroutine per IP, bounded by `-workers`
(default 20):

- The main goroutine loops over the IP list. Before spawning each worker it
  acquires a slot from a buffered channel used as a semaphore; once
  `-workers` goroutines are already in flight, the loop blocks until one of
  them finishes and frees a slot. This is what actually caps concurrency.
- Each worker goroutine tries credentials against a single IP sequentially,
  stopping at the first one that succeeds, then releases its semaphore slot
  and marks itself done on a `sync.WaitGroup`.
- A separate goroutine waits on the `WaitGroup` and, once every worker has
  finished, closes the results channel — this is what lets the main
  goroutine's result-collection loop terminate instead of blocking forever.

Credentials within a single IP are always tested sequentially; only testing
*across* IPs is parallelized. Result order is not guaranteed to match the
order of the input devices file, since IPs may finish in any order.

## Project structure

```
.
├── main.go               # entry point: wires config, filemanager, and discovery together
├── config/                # CLI flag parsing (-devices/-creds/-output/-workers) and usage text
├── discovery/              # concurrent IP/credential test loop (goroutines, semaphore, WaitGroup)
├── filemanager/            # reading input files line-by-line, writing JSON results
├── snmpmanager/             # builds gosnmp.GoSNMP params per credential, runs Connect/Get
└── snmpmodules/             # ReachableDevice type and credential-map helpers
```

## Notes

- Failed credential/IP combinations are logged and skipped — the tool keeps
  testing the remaining credentials and devices rather than exiting on the
  first failure.
- Each SNMP connection is opened and closed within a single test attempt, so
  sockets aren't held open across the full run.
- Increasing `-workers` speeds up large scans but opens more sockets at
  once; tune it down on constrained networks or when scanning through slow
  links.
