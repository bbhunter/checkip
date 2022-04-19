package check

import (
	"encoding/csv"
	"encoding/json"
	"net"
	"os"
	"strconv"

	"github.com/jreisinger/checkip"
)

type phishstats struct {
	score float64 // 0-2 likely, 2-4 suspicious, 4-6 phishing, 6-10 omg phishing!
	url   string
}

func (ps phishstats) Summary() string {
	return na(ps.url)
}

func (ps phishstats) Json() ([]byte, error) {
	return json.Marshal(ps)
}

func PhishStats(ipaddr net.IP) (checkip.Result, error) {
	result := checkip.Result{
		Name: "phishstats.info",
		Type: checkip.TypeInfoSec,
	}

	file, err := getDbFilesPath("phish_score.csv")
	if err != nil {
		return result, err
	}
	url := "https://phishstats.info/phish_score.csv"
	if err := updateFile(file, url, ""); err != nil {
		return result, err
	}

	ps, err := getPhishStats(file, ipaddr)
	if err != nil {
		return result, err
	}
	result.Info = ps
	if ps.score > 2 {
		result.Malicious = true
	}

	return result, nil
}

func getPhishStats(csvFile string, ipaddr net.IP) (phishstats, error) {
	var ps phishstats

	f, err := os.Open(csvFile)
	if err != nil {
		return ps, err
	}

	csvReader := csv.NewReader(f)
	csvReader.Comment = '#'
	records, err := csvReader.ReadAll()
	if err != nil {
		return ps, err
	}

	for _, fields := range records {
		if ipaddr.String() == fields[3] {
			score, err := strconv.ParseFloat(fields[1], 64)
			if err != nil {
				return ps, err
			}
			ps.score = score
			ps.url = fields[2]
			break
		}
	}

	return ps, nil
}
