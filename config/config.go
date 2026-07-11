// Package config handles CLI flag parsing for snmpCredsMultiTest.
package config

import (
	"flag"
	"fmt"
	"os"
	"text/tabwriter"
)

// Config holds the resolved file paths for a single run of the tool.
type Config struct {
	DeviceFile string
	CredsFile  string
	OutputFile string
}

const (
	defaultDeviceFile = "devices.txt"
	usageDeviceFile   = "input file with devices IP list"
	defaultCredsFile  = "creds.txt"
	usageCredsFile    = "input file with SNMP credentials list"
	defaultOutputFile = "reachables.json"
	usageOutputFile   = "output JSON file with reachable devices data"
)

// Parse registers the CLI flags, parses os.Args, and returns the resolved
// Config. It must be called once, near the top of main().
func Parse() Config {
	var cfg Config

	flag.StringVar(&cfg.DeviceFile, "devices", defaultDeviceFile, usageDeviceFile)
	flag.StringVar(&cfg.DeviceFile, "d", defaultDeviceFile, usageDeviceFile+" (shorthand)")

	flag.StringVar(&cfg.CredsFile, "creds", defaultCredsFile, usageCredsFile)
	flag.StringVar(&cfg.CredsFile, "c", defaultCredsFile, usageCredsFile+" (shorthand)")

	flag.StringVar(&cfg.OutputFile, "output", defaultOutputFile, usageOutputFile)
	flag.StringVar(&cfg.OutputFile, "o", defaultOutputFile, usageOutputFile+" (shorthand)")

	flag.Usage = printUsage
	flag.Parse()

	return cfg
}

// printUsage prints a custom help message grouping each flag's long and
// short form on a single line, e.g. "-c, -creds string   <usage>".
func printUsage() {
	fmt.Fprintf(os.Stderr, "snmpCredsMultiTest: test SNMP credentials across a list of devices\n\n")
	fmt.Fprintf(os.Stderr, "Usage:\n")

	rows := []struct {
		short, long, def, usage string
	}{
		{"d", "devices", defaultDeviceFile, usageDeviceFile},
		{"c", "creds", defaultCredsFile, usageCredsFile},
		{"o", "output", defaultOutputFile, usageOutputFile},
	}

	w := tabwriter.NewWriter(os.Stderr, 0, 4, 2, ' ', 0)
	for _, r := range rows {
		fmt.Fprintf(w, "  -%s, -%s string\t%s (default %q)\n", r.short, r.long, r.usage, r.def)
	}
	w.Flush()
}
