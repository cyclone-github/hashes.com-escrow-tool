package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

// convert crypto to usd via Kraken API
func toUSD(value float64, currency string) (map[string]interface{}, error) {
	if currency == "credits" {
		return map[string]interface{}{"currentprice": nil, "converted": "N/A"}, nil
	}

	url := fmt.Sprintf("https://api.kraken.com/0/public/Ticker?pair=%susd", currency)
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %v", err)
	}
	defer resp.Body.Close()

	var response struct {
		Result map[string]struct {
			A []string `json:"a"`
		} `json:"result"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("failed to decode response: %v", err)
	}

	var currentPrice string
	switch strings.ToUpper(currency) {
	case "BTC":
		currentPrice = response.Result["XXBTZUSD"].A[0]
	case "XMR":
		currentPrice = response.Result["XXMRZUSD"].A[0]
	case "LTC":
		currentPrice = response.Result["XLTCZUSD"].A[0]
	}

	currentPriceFloat, err := strconv.ParseFloat(currentPrice, 64)
	if err != nil {
		return nil, fmt.Errorf("failed to parse current price: %v", err)
	}

	converted := fmt.Sprintf("$%.3f", value*currentPriceFloat)

	return map[string]interface{}{"currentprice": currentPrice, "converted": converted}, nil
}