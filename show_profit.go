package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"os"
	"strconv"
	"text/tabwriter"
)

// show profit
func getProfit(apiKey string) error {
	fmt.Fprintln(os.Stderr, "Total Profit:")
	fmt.Fprintln(os.Stderr)

	url := fmt.Sprintf("https://hashes.com/en/api/profit?key=%s", apiKey)

	resp, err := httpClient.Get(url)
	if err != nil {
		var netErr net.Error
		if errors.As(err, &netErr) && netErr.Timeout() {
			fmt.Fprintln(os.Stderr, "Request timed out while fetching profit.")
			return nil // non-fatal
		}
		return fmt.Errorf("failed to send request: %v", err)
	}
	defer resp.Body.Close()

	var response struct {
		Success  bool              `json:"success"`
		Currency map[string]string `json:"currency"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return fmt.Errorf("failed to decode response: %v", err)
	}

	if !response.Success {
		return fmt.Errorf("request was not successful")
	}

	usd := make(map[string]string)
	for currency, value := range response.Currency {
		if valueFloat, err := strconv.ParseFloat(value, 64); err == nil {
			usdValue, _ := toUSD(valueFloat, currency)
			if usdValueConverted, ok := usdValue["converted"]; ok {
				usd[currency] = fmt.Sprintf("%v", usdValueConverted)
			}
		}
	}

	writer := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', tabwriter.AlignRight|tabwriter.Debug)
	fmt.Fprintln(writer, "Crypto \t Coins \t USD")

	for currency, amount := range response.Currency {
		fmt.Fprintf(writer, "%s \t %s \t %s\n", currency, amount, usd[currency])
	}

	writer.Flush()
	return nil
}
