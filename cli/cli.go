// Package cli contains functions for running checks from command-line.
package cli

import (
	"encoding/json"
	"fmt"
	"log"
	"math"
	"net"
	"os"
	"sort"
	"sync"

	"github.com/jreisinger/checkip/check"
)

// Run runs checks concurrently against the ippaddr.
func Run(checks []check.Check, ipaddr net.IP) (Results, []error) {
	var results Results
	var errors []error

	var wg sync.WaitGroup
	for _, chk := range checks {
		wg.Add(1)
		go func(c check.Check) {
			defer wg.Done()
			r, err := c(ipaddr)
			if err != nil {
				errors = append(errors, err)
				return
			}
			results = append(results, r)
		}(chk)
	}
	wg.Wait()

	return results, errors
}

// Results are generic or security information provided by a Check.
type Results []check.Result

// PrintJSON prints all results in JSON.
func (rs Results) PrintJSON(ipaddr net.IP) {
	// if len(rs) == 0 {
	// 	return
	// }

	_, _, prob := rs.maliciousStats()

	out := struct {
		IpAddr        net.IP  `json:"ipaddr"`
		MaliciousProb string  `json:"malicious_prob"`
		Check         Results `json:"checks"`
	}{
		ipaddr,
		fmt.Sprintf("%.2f", prob),
		rs,
	}

	enc := json.NewEncoder(os.Stdout)
	if err := enc.Encode(out); err != nil {
		log.Fatal(err)
	}
}

// SortByName sorts Results by name.
func (rs Results) SortByName() {
	sort.Slice(rs, func(i, j int) bool {
		return rs[i].Name < rs[j].Name
	})
}

// PrintSummary prints summary results from Info and InfoSec checks.
func (rs Results) PrintSummary() {
	for _, r := range rs {
		// To avoid "invalid memory address or nil pointer dereference"
		// runtime error and printing empty summary info.
		if r.Info == nil || r.Info.Summary() == "" {
			continue
		}

		if r.Type == check.TypeInfo || r.Type == check.TypeInfoSec {
			fmt.Printf("%-15s %s\n", r.Name, r.Info.Summary())
		}
	}
}

// PrintMalicious prints how many of the InfoSec and Sec checks consider the IP
// address to be malicious.
func (rs Results) PrintMalicious() {
	total, malicious, prob := rs.maliciousStats()
	msg := fmt.Sprintf("%-15s %.0f%% (%d/%d) ",
		"malicious", math.Round(prob*100), malicious, total)
	switch {
	case prob >= 0.50:
		msg += `🚫`
	case prob >= 0.15:
		msg += `🤏`
	default:
		msg += `✅`
	}
	fmt.Println(msg)
}

func (rs Results) maliciousStats() (total, malicious int, prob float64) {
	for _, r := range rs {
		// if r.Info == nil {
		// 	continue
		// }
		if r.Type == check.TypeSec || r.Type == check.TypeInfoSec {
			total++
			if r.Malicious {
				malicious++
			}
		}
	}
	prob = float64(malicious) / float64(total)
	return total, malicious, prob
}
