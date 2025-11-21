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

// withdrawal history
func withdrawalHistory(apiKey string) error {
	fmt.Fprintln(os.Stderr, "Withdrawal History (last 20):")
	fmt.Fprintln(os.Stderr)

	url := "https://hashes.com/en/api/withdrawals?key=" + apiKey
	resp, err := httpClient.Get(url)
	if err != nil {
		var netErr net.Error
		if errors.As(err, &netErr) && netErr.Timeout() {
			fmt.Fprintln(os.Stderr, "Request timed out while fetching withdrawal history.")
			return nil // non-fatal
		}
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

	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return fmt.Errorf("An error occurred: failed to decode response: %v", err)
	}

	if !response.Success {
		return fmt.Errorf("An error occurred: request was not successful")
	}

	end := 20
	if end > len(response.List) {
		end = len(response.List)
	}
	last20 := response.List[:end]

	// reverse to show most recent first
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

			// handle nil or missing rate
			if rate == nil {
				fmt.Fprintf(os.Stderr, "Warning: missing conversion rate for currency %s, using 0.000\n", currency)
				uniqueCurrencies[currency] = 0
				continue
			}

			priceVal, ok := rate["currentprice"]
			if !ok || priceVal == nil {
				fmt.Fprintf(os.Stderr, "Warning: missing current price for currency %s, using 0.000\n", currency)
				uniqueCurrencies[currency] = 0
				continue
			}

			priceStr, ok := priceVal.(string)
			if !ok {
				fmt.Fprintf(os.Stderr, "Warning: unexpected price type for currency %s, using 0.000\n", currency)
				uniqueCurrencies[currency] = 0
				continue
			}

			currentPrice, err := strconv.ParseFloat(priceStr, 64)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Warning: failed to parse price for currency %s, using 0.000\n", currency)
				uniqueCurrencies[currency] = 0
				continue
			}

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

		fmt.Fprintf(
			writer,
			"%s \t $%.3f \t $%.3f \t %s \t %s \t %s \t %s\n",
			withdrawal.Date,
			amount*conversionRate,
			afterFee*conversionRate,
			withdrawal.ID,
			withdrawal.Currency,
			withdrawal.Transaction,
			destinationWallet,
		)
	}

	writer.Flush()
	return nil
}
