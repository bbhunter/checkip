package checkip

import (
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"strings"
)

// Shodan holds information about an IP address from shodan.io scan data.
type Shodan struct {
	Org   string `json:"org"`
	Data  data   `json:"data"`
	Os    string `json:"os"`
	Ports []int  `json:"ports"`
}

type data []struct {
	Product string `json:"product"`
	Version string `json:"version"`
	Port    int    `json:"port"`
}

// Check fills in Shodan data for a given IP address. Its get the data from
// https://api.shodan.io. It returns false if version of at least one listening
// service is known.
func (s *Shodan) Check(ipaddr net.IP) (bool, error) {
	apiKey, err := getConfigValue("SHODAN_API_KEY")
	if err != nil {
		return true, fmt.Errorf("can't call API: %w", err)
	}

	resp, err := http.Get(fmt.Sprintf("https://api.shodan.io/shodan/host/%s?key=%s", ipaddr, apiKey))
	if err != nil {
		return true, err
	}
	defer resp.Body.Close()

	if err := json.NewDecoder(resp.Body).Decode(&s); err != nil {
		return true, err
	}

	return s.isOK(), nil
}

func (s *Shodan) isOK() bool {
	return !s.gotServiceVersion()
}

func (s *Shodan) gotServiceVersion() bool {
	for _, d := range s.Data {
		if d.Version != "" && !okBanner(d.Product) {
			return true
		}
	}
	return false
}

func okBanner(product string) bool {
	okBannersPrefixes := []string{
		"OpenSSH",
	}
	for _, prefix := range okBannersPrefixes {
		if strings.HasPrefix(product, prefix) {
			return true
		}
	}
	return false
}

// String returns the result of the check.
func (s *Shodan) String() string {
	os := "OS unknown"
	if s.Os != "" {
		os = s.Os
	}

	var portInfo []string
	for _, d := range s.Data {
		var product string
		if d.Product != "" {
			product = d.Product
		}

		var version string
		if d.Version != "" {
			version = d.Version
		}

		if product == "" && version == "" {
			portInfo = append(portInfo, fmt.Sprintf("%d", d.Port))
		} else {
			portInfo = append(portInfo, fmt.Sprintf("%d (%s, %s)", d.Port, product, version))
		}
	}

	portStr := "port"
	if len(portInfo) != 1 {
		portStr += "s"
	}
	if len(portInfo) > 0 {
		portStr += ":"
	}

	return fmt.Sprintf("OS and ports\t%s, %d open %s %s", os, len(portInfo), portStr, strings.Join(portInfo, ", "))
}
