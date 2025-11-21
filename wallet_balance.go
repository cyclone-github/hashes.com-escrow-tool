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

// get wallet balance
func getWalletBalances(apiKey string) error {
	fmt.Fprintln(os.Stderr, "Wallet Balance:")
	fmt.Fprintln(os.Stderr)

	url := fmt.Sprintf("https://hashes.com/en/api/balance?key=%s", apiKey)

	resp, err := httpClient.Get(url)
	if err != nil {
		var netErr net.Error
		if errors.As(err, &netErr) && netErr.Timeout() {
			fmt.Fprintln(os.Stderr, "Request timed out while fetching wallet balance.")
			return nil // non-fatal
		}
		return fmt.Errorf("failed to send request: %v", err)
	}
	defer resp.Body.Close()

	var response struct {
		Success bool `json:"success"`
		WalletBalances
	}
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return fmt.Errorf("failed to decode response: %v", err)
	}

	if !response.Success {
		return fmt.Errorf("request was not successful")
	}

	walletBalances := response.WalletBalances
	writer := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', tabwriter.AlignRight|tabwriter.Debug)
	fmt.Fprintln(writer, "Crypto \t Coins \t USD")

	// Convert and print balances for BTC, XMR, LTC
	for _, crypto := range []struct {
		name  string
		value string
	}{
		{"BTC", walletBalances.BTC},
		{"XMR", walletBalances.XMR},
		{"LTC", walletBalances.LTC},
	} {
		amount, err := strconv.ParseFloat(crypto.value, 64)
		if err != nil {
			return fmt.Errorf("Error parsing %s amount: %v\n", crypto.name, err)
		}

		usdMap, err := toUSD(amount, crypto.name)
		if err != nil {
			return fmt.Errorf("Error converting %s: %v\n", crypto.name, err)
		}

		// if toUSD times out, show N/A
		converted := "N/A"
		if usdMap != nil {
			if v, ok := usdMap["converted"]; ok {
				if s, ok2 := v.(string); ok2 {
					converted = s
				} else {
					converted = fmt.Sprintf("%v", v)
				}
			}
		}

		fmt.Fprintf(writer, "%s \t %s \t %s\n", crypto.name, crypto.value, converted)
	}

	fmt.Fprintf(writer, "Credits \t %s\n", walletBalances.Credits)
	writer.Flush()

	return nil
}
