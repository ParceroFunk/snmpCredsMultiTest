package discovery

import (
	"log"
	"sync"

	"github.com/ParceroFunk/snmpCredsMultiTest/snmpmanager"
	"github.com/ParceroFunk/snmpCredsMultiTest/snmpmodules"
)

// Run tests every credential against every IP and returns the devices
// that responded successfully. IPs are tested concurrently (bounded by
// maxGoRoutines); credentials within a single IP are tested
// sequentally, stoping at the first one that succeeds.
func Run(ips []string, creds [][]string, maxGoRoutines int) []snmpmodules.ReachableDevice {
	oids := []string{snmpmanager.SysName, snmpmanager.SysDescr}

	resultsCh := make(chan snmpmodules.ReachableDevice, len(ips))
	sem := make(chan struct{}, maxGoRoutines)

	var wg sync.WaitGroup
	for _, ip := range ips {
		wg.Add(1)
		sem <- struct{}{} // acquire a slot; blocks once maxConcurrency is in flight

		go func(ip string) {
			defer wg.Done()
			defer func() { <-sem }() // release the slot

			if device, ok := testIP(ip, creds, oids); ok {
				resultsCh <- device
			}
		}(ip)
	}

	// Close the results channel once every goroutine has finished, so the
	// range below terminates instead of blocking forever.
	go func() {
		wg.Wait()
		close(resultsCh)
	}()

	var reachables []snmpmodules.ReachableDevice
	for device := range resultsCh {
		reachables = append(reachables, device)
	}

	return reachables
}

// testIP tries each credential against a single IP, sequentially, stopping
// at the first one that succeeds. Returns ok=false if none worked.
func testIP(ip string, creds [][]string, oids []string) (snmpmodules.ReachableDevice, bool) {
	for _, cred := range creds {
		result, err := snmpmanager.TestCredential(ip, cred, oids)
		if err != nil {
			log.Printf("[%s] cred %v failed: %v", ip, cred, err)
			continue
		}

		log.Printf("[%s] cred %v succeeded", ip, cred)

		credMap, err := snmpmodules.CreateCredMapFromCredLength(cred[1:])
		if err != nil {
			log.Printf("[%s] cred %v failed to parse for saving: %v", ip, cred, err)
		}

		var response [2]string
		for index, variable := range result.Variables {
			switch value := variable.Value.(type) {
			case []byte:
				response[index] = string(value)
			case string:
				response[index] = value
			default:
				log.Printf("[%s] response for %v failed to parse for saving. unknown type", ip, variable.Name)
			}
		}

		return snmpmodules.New(ip, response[0], response[1], cred[0], credMap), true
	}

	return snmpmodules.ReachableDevice{}, false
}
