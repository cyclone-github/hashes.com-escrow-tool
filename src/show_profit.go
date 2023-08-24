package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"text/tabwriter"
)

// show profit
func getProfit(apiKey string) error {
	fmt.Println("Total Profit:\n")
	url := fmt.Sprintf("https://hashes.com/en/api/profit?key=%s", apiKey)

	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("failed to send request: %v", err)
	}
	defer resp.Body.Close()

	var response struct {
		Success  bool              `json:"success"`
		Currency map[string]string `json:"currency"`
	}
	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		return fmt.Errorf("failed to decode response: %v", err)
	}

	if !response.Success {
		return fmt.Errorf("request was not successful")
	}

	usd := make(map[string]string)
	for currency, value := range response.Currency {
		if valueFloat, err := strconv.ParseFloat(value, 64); err == nil {
			usdValue, _ := toUSD(valueFloat, currency)
			usdValueConverted, ok := usdValue["converted"]
			if ok {
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
