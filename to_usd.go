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

type conversionRates struct {
	BTC string
	XMR string
	LTC string
}

func fetchConversionRates() (*conversionRates, error) {
	resp, err := httpClient.Get(hashesAPIBaseURL + "/conversion")
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

	return &conversionRates{
		BTC: response.BTC,
		XMR: response.XMR,
		LTC: response.LTC,
	}, nil
}

func conversionPrice(rates *conversionRates, currency string) (float64, error) {
	if rates == nil {
		return 0, nil
	}

	var currentPrice string
	switch strings.ToUpper(currency) {
	case "BTC":
		currentPrice = rates.BTC
	case "XMR":
		currentPrice = rates.XMR
	case "LTC":
		currentPrice = rates.LTC
	default:
		return 0, fmt.Errorf("unsupported currency: %s", currency)
	}

	return strconv.ParseFloat(currentPrice, 64)
}

func currentPriceFromRate(rate map[string]interface{}) float64 {
	if rate == nil {
		return 0
	}

	priceVal, ok := rate["currentprice"]
	if !ok || priceVal == nil {
		return 0
	}

	priceStr, ok := priceVal.(string)
	if !ok {
		return 0
	}

	return parseFloat(priceStr)
}

// convert crypto to usd via hashes.com API
func toUSD(value float64, currency string) (map[string]interface{}, error) {
	if currency == "credits" {
		return map[string]interface{}{
			"currentprice": nil,
			"converted":    "N/A",
		}, nil
	}

	rates, err := fetchConversionRates()
	if err != nil {
		return nil, err
	}
	if rates == nil {
		return nil, nil
	}

	currentPriceFloat, err := conversionPrice(rates, currency)
	if err != nil {
		return nil, err
	}

	currentPrice := ""
	switch strings.ToUpper(currency) {
	case "BTC":
		currentPrice = rates.BTC
	case "XMR":
		currentPrice = rates.XMR
	case "LTC":
		currentPrice = rates.LTC
	}

	converted := fmt.Sprintf("$%.3f", value*currentPriceFloat)

	return map[string]interface{}{
		"currentprice": currentPrice,
		"converted":    converted,
	}, nil
}
