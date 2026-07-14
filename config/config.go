// Package config handles CLI flag parsing for snmpCredsMultiTest.
package config

import (
	"flag"
	"fmt"
	"os"
	"strconv"
	"text/tabwriter"
)

// Config holds the resolved file paths for a single run of the tool.
type Config struct {
	DeviceFile     string
	CredsFile      string
	OutputFile     string
	MaxConcurrency int
}

const (
	defaultDeviceFile = "devices.txt"
	usageDeviceFile   = "input file with devices IP list"
	defaultCredsFile  = "creds.txt"
	usageCredsFile    = "input file with SNMP credentials list"
	defaultOutputFile = "reachables.csv"
	usageOutputFile   = "output JSON file with reachable devices data"
	// Concurrency values
	defaultMaxConcurrency = 20
	usageMaxConcurrency   = "how many IPs are tested in parallel"
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

	flag.IntVar(&cfg.MaxConcurrency, "workers", defaultMaxConcurrency, usageMaxConcurrency)
	flag.IntVar(&cfg.MaxConcurrency, "w", defaultMaxConcurrency, usageMaxConcurrency+" (shorthand)")

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
		short, long, flagType, def, usage string
	}{
		{"d", "devices", fmt.Sprintf("%T", defaultDeviceFile), defaultDeviceFile, usageDeviceFile},
		{"c", "creds", fmt.Sprintf("%T", defaultCredsFile), defaultCredsFile, usageCredsFile},
		{"o", "output", fmt.Sprintf("%T", defaultOutputFile), defaultOutputFile, usageOutputFile},
		{"w", "workers", fmt.Sprintf("%T", defaultMaxConcurrency), strconv.Itoa(defaultMaxConcurrency), usageMaxConcurrency},
	}

	w := tabwriter.NewWriter(os.Stderr, 0, 4, 2, ' ', 0)
	for _, r := range rows {
		fmt.Fprintf(w, "  -%s, -%s\t%s\t%s (default %q)\n", r.short, r.long, r.flagType, r.usage, r.def)
	}
	w.Flush()
}
