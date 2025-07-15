package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"text/tabwriter"
)

// get found history
func getFoundHistory(apiKey string) error {
	url := fmt.Sprintf("https://hashes.com/en/api/uploads?key=%s", apiKey)

	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("failed to send request: %v", err)
	}
	defer resp.Body.Close()

	var response struct {
		Success bool           `json:"success"`
		List    []FoundHistory `json:"list"`
	}
	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		return fmt.Errorf("failed to decode response: %v", err)
	}

	if !response.Success {
		return fmt.Errorf("request was not successful")
	}

	fmt.Fprintln(os.Stderr, "Upload History (last 20):\n")
	writer := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', tabwriter.AlignRight|tabwriter.Debug)
	defer writer.Flush()

	fmt.Fprintln(writer, "ID\tDate/Time\tBTC (USD)\tXMR (USD)\tLTC (USD)\t-m\tTotal\tFound\tStatus")

	btcRate, _ := toUSD(1, "BTC")
	xmrRate, _ := toUSD(1, "XMR")
	ltcRate, _ := toUSD(1, "LTC")

	startIndex := len(response.List) - 20
	if startIndex < 0 {
		startIndex = 0
	}

	for i := startIndex; i < len(response.List); i++ {
		history := response.List[i]

		btc := parseFloat(history.BTC)
		xmr := parseFloat(history.XMR)
		ltc := parseFloat(history.LTC)

		btcUSD := fmt.Sprintf("$%.3f", btc*parseFloat(btcRate["currentprice"].(string)))
		xmrUSD := fmt.Sprintf("$%.3f", xmr*parseFloat(xmrRate["currentprice"].(string)))
		ltcUSD := fmt.Sprintf("$%.3f", ltc*parseFloat(ltcRate["currentprice"].(string)))

		fmt.Fprintf(writer, "%d\t%s\t%s\t%s\t%s\t%d\t%d\t%d\t%s\n",
			history.ID, history.Date, btcUSD, xmrUSD, ltcUSD, history.AlgorithmID, history.TotalHashes, history.ValidHashes, history.Status)
	}

	return nil
}
