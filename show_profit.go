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

	url := fmt.Sprintf("%s/profit?key=%s", hashesAPIBaseURL, apiKey)

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

	rates, err := fetchConversionRates()
	if err != nil {
		return err
	}

	usd := make(map[string]string)
	for currency, value := range response.Currency {
		if valueFloat, err := strconv.ParseFloat(value, 64); err == nil {
			if rates == nil {
				continue
			}

			currentPriceFloat, err := conversionPrice(rates, currency)
			if err != nil {
				continue
			}

			usd[currency] = fmt.Sprintf("$%.3f", valueFloat*currentPriceFloat)
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
