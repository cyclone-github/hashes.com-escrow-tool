package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"text/tabwriter"
)

// withdrawal history
func withdrawalHistory(apiKey string) error {
	fmt.Fprintln(os.Stderr, "Withdrawal History (last 20):\n")
	url := "https://hashes.com/en/api/withdrawals?key=" + apiKey
	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("An error occurred: failed to send request: %v", err)
	}
	defer resp.Body.Close()

	var response struct {
		Success bool `json:"success"`
		List    []struct {
			ID          string `json:"id"`
			Amount      string `json:"amount"`
			AfterFee    string `json:"afterFee"`
			Transaction string `json:"transaction"`
			Currency    string `json:"currency"`
			Destination string `json:"destination"`
			Date        string `json:"date"`
		} `json:"list"`
	}

	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		return fmt.Errorf("An error occurred: failed to decode response: %v", err)
	}

	end := 20
	if end > len(response.List) {
		end = len(response.List)
	}
	last20 := response.List[:end]

	for i, j := 0, len(last20)-1; i < j; i, j = i+1, j-1 {
		last20[i], last20[j] = last20[j], last20[i]
	}

	uniqueCurrencies := map[string]float64{}
	for _, withdrawal := range last20 {
		currency := withdrawal.Currency
		if _, exists := uniqueCurrencies[currency]; !exists {
			rate, err := toUSD(1, currency)
			if err != nil {
				return fmt.Errorf("An error occurred: %v", err)
			}
			currentPrice, _ := strconv.ParseFloat(rate["currentprice"].(string), 64)
			uniqueCurrencies[currency] = currentPrice
		}
	}

	writer := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', tabwriter.Debug)
	fmt.Fprintln(writer, "Date\t Amount USD \t After Fee \t ID \t Crypto \t Transaction ID \t Destination Wallet")

	for _, withdrawal := range last20 {
		amount, _ := strconv.ParseFloat(withdrawal.Amount, 64)
		afterFee, _ := strconv.ParseFloat(withdrawal.AfterFee, 64)
		conversionRate := uniqueCurrencies[withdrawal.Currency]
		destinationWallet := withdrawal.Destination

		fmt.Fprintf(writer, "%s \t $%.3f \t $%.3f \t %s \t %s \t %s \t %s\n",
			withdrawal.Date, amount*conversionRate, afterFee*conversionRate, withdrawal.ID, withdrawal.Currency, withdrawal.Transaction, destinationWallet)
	}

	writer.Flush()
	return nil
}
