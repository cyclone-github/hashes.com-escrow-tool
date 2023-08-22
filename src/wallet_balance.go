package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"text/tabwriter"
)

// case 5, get wallet balance
func getWalletBalances(apiKey string) error {
	url := fmt.Sprintf("https://hashes.com/en/api/balance?key=%s", apiKey)

	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("failed to send request: %v", err)
	}
	defer resp.Body.Close()

	var response struct {
		Success bool `json:"success"`
		WalletBalances
	}
	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		return fmt.Errorf("failed to decode response: %v", err)
	}

	if !response.Success {
		return fmt.Errorf("request was not successful")
	}

	walletBalances := response.WalletBalances
	writer := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', tabwriter.AlignRight|tabwriter.Debug)
	fmt.Fprintln(writer, "Crypto \t Coins \t USD")

	// Convert and print balances for BTC, XMR, LTC
	for _, crypto := range []struct{ name, value string }{
		{"BTC", walletBalances.BTC},
		{"XMR", walletBalances.XMR},
		{"LTC", walletBalances.LTC},
	} {
		amount, err := strconv.ParseFloat(crypto.value, 64)
		if err != nil {
			return fmt.Errorf("Error parsing %s amount: %v\n", crypto.name, err)
		}
		cryptoUSD, err := toUSD(amount, crypto.name)
		if err != nil {
			return fmt.Errorf("Error converting %s: %v\n", crypto.name, err)
		}
		fmt.Fprintf(writer, "%s \t %s \t %s\n", crypto.name, crypto.value, cryptoUSD["converted"])
	}

	fmt.Fprintf(writer, "Credits \t %s\n", walletBalances.Credits)
	writer.Flush()

	return nil
}