package main

import (
	"encoding/base64"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
	"time"
)

/*
Cyclone's Hashes.com API Escrow Tool
github.com/cyclone-github/hashes.com-escrow-tool

This tool requires an API key from hashes.com
'Search Hashes' requires credits
See hashes.com for more info

version history
v0.1.0; 2023-08-18.1630
	initial release
v0.1.1; 2023-08-19.1945
	added withdrawal history
v0.1.1; 2023-08-20.1630
	modified withdrawal history output with tabwriter
v0.1.2; 2023-08-24.1540
	added download left lists
v0.1.3; 2024-01-14.1600
	cleaned/updated code
	changed last nth history from 10 to 20
	updated API key encryption logic
v1.0.0; 2024-04-14.1500
	fixed download deduplication
	release v1.0.0
v1.0.1; 2025-07-12.2145
	maintenance release
v1.1.0; 2025-07-14.0900
	replace Kraken crypto to USD conversion with https://hashes.com/en/api/conversion
v1.1.1; 2025-07-15.0950
	updated print statements to use stdout / stderr where applicable
	beta: added escrow websocket monitor (option #9)
	beta: added -websocket flag to start tool in escrow websocket monitor mode
v1.1.2; 2025-11-21
	fixed redundant new line logic
	added http timeouts
*/

// main function
func main() {
	cycloneFlag := flag.Bool("cyclone", false, "Display coded message")
	versionFlag := flag.Bool("version", false, "Display version")
	wsMode := flag.Bool("websocket", false, "start in websocket monitor mode (menu option 9)")
	flag.Parse()

	if *cycloneFlag {
		codedBy := "Q29kZWQgYnkgY3ljbG9uZSA7KQo="
		decoded, _ := base64.StdEncoding.DecodeString(codedBy)
		fmt.Fprintln(os.Stderr, string(decoded))
		return
	}
	if *versionFlag {
		version := "Cyclone's Hashes.com API Escrow Tool v1.1.2; 2025-11-21"
		fmt.Fprintln(os.Stderr, version)
		return
	}
	// check for API key
	getEncryptionKey()

	if *wsMode {
		apiKey := getAPIKey(false)
		if err := monitorWebsocket(apiKey); err != nil {
			fmt.Fprintf(os.Stderr, "An error occurred: %v\n", err)
		}
		return // exit
	}

	clearScreen()
	printCyclone()
	fmt.Fprintln(os.Stderr, " ######################################################################")
	fmt.Fprintln(os.Stderr, "#              Cyclone's Hashes.com API Escrow Tool v1.1.2             #")
	fmt.Fprintln(os.Stderr, "#           github.com/cyclone-github/hashes.com-escrow-tool           #")
	fmt.Fprintln(os.Stderr, "#            This tool requires an API key from hashes.com             #")
	fmt.Fprintln(os.Stderr, "#                   'Search Hashes' requires credits                   #")
	fmt.Fprintln(os.Stderr, "#                     See hashes.com for more info                     #")
	fmt.Fprintln(os.Stderr, " ######################################################################")
	fmt.Fprintln(os.Stderr)
	apiKey := getAPIKey(false)
	for {
		// CLI Menu
		time.Sleep(100 * time.Millisecond)
		fmt.Fprintln(os.Stderr)
		fmt.Fprintln(os.Stderr, "Select an option:")
		fmt.Fprintln(os.Stderr, "1.  Upload Founds")
		fmt.Fprintln(os.Stderr, "2.  Upload History")
		fmt.Fprintln(os.Stderr, "3.  Download Left Lists")
		fmt.Fprintln(os.Stderr, "4.  Search Hashes")
		fmt.Fprintln(os.Stderr, "5.  Hash Identifier")
		fmt.Fprintln(os.Stderr, "6.  Wallet Balance")
		fmt.Fprintln(os.Stderr, "7.  Show Profit")
		fmt.Fprintln(os.Stderr, "8.  Withdrawal History")
		fmt.Fprintln(os.Stderr, "9.  Monitor Escrow WebSocket")
		fmt.Fprintln(os.Stderr, "n.  Enter New API")
		fmt.Fprintln(os.Stderr, "r.  Remove API Key")
		fmt.Fprintln(os.Stderr, "c.  Clear Screen")
		fmt.Fprintln(os.Stderr, "q.  Quit")
		var choice string
		fmt.Scanln(&choice)

		switch strings.ToLower(choice) {
		case "1":
			// Upload Founds
			clearScreen()
			if err := uploadFounds(apiKey); err != nil {
				fmt.Fprintf(os.Stderr, "An error occurred: %v\n", err)
			}
		case "2":
			// Show Upload History
			clearScreen()
			if err := getFoundHistory(apiKey); err != nil {
				fmt.Fprintf(os.Stderr, "An error occurred: %v\n", err)
			}
		case "3":
			// Download Left Lists
			clearScreen()
			if err := downloadLeftList(apiKey); err != nil {
				fmt.Fprintf(os.Stderr, "An error occurred: %v\n", err)
			}
		case "4":
			// Search Hashes
			clearScreen()
			hashPlaintext := pasteHashes("Search Hashes.com")
			hashes := strings.Split(hashPlaintext, "\n")
			if err := searchHashes(apiKey, hashes); err != nil {
				fmt.Fprintf(os.Stderr, "An error occurred: %v\n", err)
			}
		case "5":
			// Hash Identifier
			clearScreen()
			hashPlaintext := pasteHashes("Hashes.com Identifier")
			hashes := strings.Split(hashPlaintext, "\n")
			if len(hashes) == 0 {
				fmt.Fprintln(os.Stderr, "No hash provided.")
				break
			}
			for _, hash := range hashes {
				if err := hashIdentifier(hash, true); err != nil {
					fmt.Fprintf(os.Stderr, "An error occurred: %v\n", err)
				}
			}
		case "6":
			// Wallet Balance
			clearScreen()
			if err := getWalletBalances(apiKey); err != nil {
				fmt.Fprintf(os.Stderr, "An error occurred: %v\n", err)
			}
		case "7":
			// Show Profit
			clearScreen()
			if err := getProfit(apiKey); err != nil {
				fmt.Fprintf(os.Stderr, "An error occurred: %v\n", err)
			}
		case "8":
			// Withdrawal History
			clearScreen()
			if err := withdrawalHistory(apiKey); err != nil {
				fmt.Fprintf(os.Stderr, "An error occurred: %v\n", err)
			}
		case "9":
			// Monitor Escrow WebSocket // beta
			clearScreen()
			log.Println("Starting Hashes.com Escrow WebSocket Monitor\nPress CTR+C to quit")
			fmt.Fprintln(os.Stderr)
			if err := monitorWebsocket(apiKey); err != nil {
				fmt.Fprintf(os.Stderr, "An error occurred: %v\n", err)
			}
		case "n":
			// Enter New API
			getAPIKey(true)
		case "r":
			// Remove API Key
			clearScreen()
			removeAPIKey()
		case "c":
			// clear screen
			clearScreen()
		case "q":
			// exit program
			return
		default:
			fmt.Fprintln(os.Stderr, "Invalid choice, please try again.")
			time.Sleep(500 * time.Millisecond)
		}
	}
}
