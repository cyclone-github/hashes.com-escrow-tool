package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
)

// convert crypto to usd via hashes.com API
func toUSD(value float64, currency string) (map[string]interface{}, error) {
	if currency == "credits" {
		return map[string]interface{}{
			"currentprice": nil,
			"converted":    "N/A",
		}, nil
	}

	url := "https://hashes.com/en/api/conversion"
	resp, err := httpClient.Get(url)
	if err != nil {
		var netErr net.Error
		if errors.As(err, &netErr) && netErr.Timeout() {
			fmt.Fprintln(os.Stderr, "Request timed out while fetching USD conversion rates.")
			return nil, nil // non-fatal, caller handles missing data
		}
		return nil, fmt.Errorf("failed to send request: %v", err)
	}
	defer resp.Body.Close()

	var response struct {
		Success bool   `json:"success"`
		BTC     string `json:"BTC"`
		XMR     string `json:"XMR"`
		LTC     string `json:"LTC"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("failed to decode response: %v", err)
	}
	if !response.Success {
		return nil, fmt.Errorf("API returned success=false")
	}

	var currentPrice string
	switch strings.ToUpper(currency) {
	case "BTC":
		currentPrice = response.BTC
	case "XMR":
		currentPrice = response.XMR
	case "LTC":
		currentPrice = response.LTC
	default:
		return nil, fmt.Errorf("unsupported currency: %s", currency)
	}

	currentPriceFloat, err := strconv.ParseFloat(currentPrice, 64)
	if err != nil {
		return nil, fmt.Errorf("failed to parse current price: %v", err)
	}

	converted := fmt.Sprintf("$%.3f", value*currentPriceFloat)

	return map[string]interface{}{
		"currentprice": currentPrice,
		"converted":    converted,
	}, nil
}
