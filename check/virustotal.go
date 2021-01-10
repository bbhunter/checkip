package check

import (
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"net/url"

	"github.com/jreisinger/checkip/util"
)

// Check fills in data for a given IP address from virustotal API. It returns
// false if the IP address is considered malicious.
func (vt *VirusTotal) Check(ipaddr net.IP) (bool, error) {
	apiKey, err := util.GetConfigValue("VIRUSTOTAL_API_KEY")
	if err != nil {
		return false, fmt.Errorf("can't call API: %w", err)
	}

	// curl --header "x-apikey:$VIRUSTOTAL_API_KEY" https://www.virustotal.com/api/v3/ip_addresses/1.1.1.1

	baseURL, err := url.Parse("https://www.virustotal.com/api/v3/ip_addresses/" + ipaddr.String())
	if err != nil {
		return false, err
	}

	req, err := http.NewRequest("GET", baseURL.String(), nil)
	if err != nil {
		return false, err
	}

	// Set request headers.
	req.Header.Set("x-apikey", apiKey)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return false, fmt.Errorf("%s", resp.Status)
	}

	if err := json.NewDecoder(resp.Body).Decode(vt); err != nil {
		return false, err
	}

	if vt.Data.Attributes.LastAnalysisStats.Malicious > 0 ||
		vt.Data.Attributes.LastAnalysisStats.Suspicious > 0 {
		return false, nil
	}

	return true, nil
}

// Name returns the name of the check.
func (vt *VirusTotal) Name() string {
	return fmt.Sprint("VirusTotal")
}

// String returns the output of the check.
func (vt *VirusTotal) String() string {
	return fmt.Sprintf("scannners results: %d malicious, %d suspicious, %d harmless", vt.Data.Attributes.LastAnalysisStats.Malicious, vt.Data.Attributes.LastAnalysisStats.Suspicious, vt.Data.Attributes.LastAnalysisStats.Harmless)
}