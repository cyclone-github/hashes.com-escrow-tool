package main

import (
	"fmt"
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

	toolName        = "Cyclone's Hashes.com API Escrow Tool"
	toolVersion     = "v1.1.4"
	toolVersionDate = "2026-05-22"

	hashesBaseURL    = "https://hashes.com"
	hashesAPIBaseURL = hashesBaseURL + "/en/api"

	maxHistoryEntries     = 20
	maxSearchHashes       = 250
	downloadWorkerCount   = 10
	leftListChannelBuffer = 100

	httpResponseTimeout = 10 * time.Second
	downloadBodyTimeout = 120 * time.Second
	dialKeepAlive       = 30 * time.Second
)

var (
	encryptionKey string

	httpClient = &http.Client{
		Timeout: httpResponseTimeout,
	}

	downloadHTTPClient = &http.Client{
		Transport: &http.Transport{
			DialContext:           netDialer.DialContext,
			TLSHandshakeTimeout:   httpResponseTimeout,
			ResponseHeaderTimeout: httpResponseTimeout,
		},
	}

	netDialer = &net.Dialer{
		Timeout:   httpResponseTimeout,
		KeepAlive: dialKeepAlive,
	}
)

func formatToolVersion() string {
	return fmt.Sprintf("%s %s; %s", toolName, toolVersion, toolVersionDate)
}
