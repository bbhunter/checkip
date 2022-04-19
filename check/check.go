// Package check contains functions that can check an IP address.
package check

import (
	"github.com/jreisinger/checkip"
)

// Passive checks don't interact directly with the target IP address.
var Passive = []checkip.Check{
	AbuseIPDB,
	BlockList,
	CinsScore,
	DBip,
	DnsMX,
	DnsName,
	Firehol,
	IPSum,
	IPtoASN,
	MaxMind,
	OTX,
	PhishStats,
	Shodan,
	ThreadCrowd,
	UrlScan,
	VirusTotal,
}

// Active checks interact with the target IP address.
var Active = []checkip.Check{
	Ping,
	TcpPorts,
}

// na returns "n/a" if s is empty.
func na(s string) string {
	if s == "" {
		return "n/a"
	}
	return s
}

// nonEmpty returns strings that are not empty.
func nonEmpty(strings ...string) []string {
	var ss []string
	for _, s := range strings {
		if s != "" {
			ss = append(ss, s)
		}
	}
	return ss
}
