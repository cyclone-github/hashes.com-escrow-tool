package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"os"
	"text/tabwriter"
)

// get found history
func getFoundHistory(apiKey string) error {
	url := fmt.Sprintf("https://hashes.com/en/api/uploads?key=%s", apiKey)

	resp, err := httpClient.Get(url)
	if err != nil {
		var netErr net.Error
		if errors.As(err, &netErr) && netErr.Timeout() {
			fmt.Fprintln(os.Stderr, "Request timed out while fetching upload history.")
			return nil
		}
		return fmt.Errorf("failed to send request: %v", err)
	}
	defer resp.Body.Close()

	var response struct {
		Success bool           `json:"success"`
		List    []FoundHistory `json:"list"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return fmt.Errorf("failed to decode response: %v", err)
	}

	if !response.Success {
		return fmt.Errorf("request was not successful")
	}

	fmt.Fprintln(os.Stderr, "Upload History (last 20):")
	fmt.Fprintln(os.Stderr)

	writer := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', tabwriter.AlignRight|tabwriter.Debug)
	defer writer.Flush()

	fmt.Fprintln(writer, "ID\tDate/Time\tBTC (USD)\tXMR (USD)\tLTC (USD)\t-m\tTotal\tFound\tStatus")

	// USD conversions
	btcRate, _ := toUSD(1, "BTC")
	xmrRate, _ := toUSD(1, "XMR")
	ltcRate, _ := toUSD(1, "LTC")

	startIndex := len(response.List) - 20
	if startIndex < 0 {
		startIndex = 0
	}

	for i := startIndex; i < len(response.List); i++ {
		h := response.List[i]

		btc := parseFloat(h.BTC)
		xmr := parseFloat(h.XMR)
		ltc := parseFloat(h.LTC)

		btcUSD := fmt.Sprintf("$%.3f", btc*parseFloat(btcRate["currentprice"].(string)))
		xmrUSD := fmt.Sprintf("$%.3f", xmr*parseFloat(xmrRate["currentprice"].(string)))
		ltcUSD := fmt.Sprintf("$%.3f", ltc*parseFloat(ltcRate["currentprice"].(string)))

		fmt.Fprintf(
			writer,
			"%d\t%s\t%s\t%s\t%s\t%d\t%d\t%d\t%s\n",
			h.ID,
			h.Date,
			btcUSD,
			xmrUSD,
			ltcUSD,
			h.AlgorithmID,
			h.TotalHashes,
			h.ValidHashes,
			h.Status,
		)
	}

	return nil
}
