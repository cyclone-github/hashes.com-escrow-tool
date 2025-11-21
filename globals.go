package main

import (
	"net"
	"net/http"
	"time"
)

// global structs, constants and variables
type FoundHistory struct {
	ID          int    `json:"id"`
	BTC         string `json:"btc"`
	XMR         string `json:"xmr"`
	LTC         string `json:"ltc"`
	Date        string `json:"date"`
	TotalHashes int    `json:"totalHashes"`
	ValidHashes int    `json:"validHashes"`
	Status      string `json:"status"`
	Algorithm   string `json:"algorithm"`
	AlgorithmID int    `json:"algorithmId"`
}

type WalletBalances struct {
	BTC     string `json:"BTC"`
	XMR     string `json:"XMR"`
	LTC     string `json:"LTC"`
	Credits string `json:"credits"`
}

const (
	apiKeyFile          = "api_key.enc"
	base64StaticSeedKey = "NWl5cTk3RlEwZy9HODFBQTU3NF5lZU0lel0zSwo="
)

var (
	encryptionKey string

	httpClient = &http.Client{
		Timeout: 5 * time.Second,
	}

	netDialer = &net.Dialer{
		Timeout:   5 * time.Second,
		KeepAlive: 30 * time.Second,
	}
)
