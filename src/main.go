package main

import (
	"fmt"
	"strings"
	"time"
)

// version history
// v0.1.0; 2023-08-22.1100; initial github release
// v0.1.2; 2023-08-24.1540; added download left lists

// main function
func main() {
	clearScreen()
	printCyclone()
	fmt.Println(" ######################################################################")
	fmt.Println("#              Cyclone's Hashes.com API Escrow Tool v0.1.2             #")
	fmt.Println("#            This tool requires an API key from hashes.com             #")
	fmt.Println("#                   'Search Hashes' requires credits                   #")
	fmt.Println("#                     See hashes.com for more info                     #")
	fmt.Println(" ######################################################################\n")
	// check for API key
	getEncryptionKey()
	apiKey := getAPIKey(false)

	for {
		// CLI Menu
		time.Sleep(100 * time.Millisecond)
		fmt.Println("\nSelect an option:")
		fmt.Println("1.  Upload Founds")
		fmt.Println("2.  Upload History")
		fmt.Println("3.  Download Left Lists")
		fmt.Println("4.  Search Hashes")
		fmt.Println("5.  Hash Identifier")
		fmt.Println("6.  Wallet Balance")
		fmt.Println("7.  Show Profit")
		fmt.Println("8.  Withdrawal History")
		fmt.Println("n.  Enter New API")
		fmt.Println("r.  Remove API Key")
		fmt.Println("c.  Clear Screen")
		fmt.Println("q.  Quit")
		var choice string
		fmt.Scanln(&choice)

		switch strings.ToLower(choice) {
		case "1":
			// Upload Founds
			clearScreen()
			if err := uploadFounds(apiKey); err != nil {
				fmt.Printf("An error occurred: %v\n", err)
			}
		case "2":
			// Show Upload History
			clearScreen()
			if err := getFoundHistory(apiKey); err != nil {
				fmt.Printf("An error occurred: %v\n", err)
			}
		case "3":
			// Download Left Lists
			clearScreen()
			if err := downloadLeftList(apiKey); err != nil {
				fmt.Printf("An error occurred: %v\n", err)
			}
		case "4":
			// Search Hashes
			clearScreen()
			hashPlaintext := pasteHashes()
			hashes := strings.Split(hashPlaintext, "\n")
			if err := searchHashes(apiKey, hashes); err != nil {
				fmt.Printf("An error occurred: %v\n", err)
			}
		case "5":
			// Hash Identifier
			clearScreen()
			hashPlaintext := pasteHashes()
			hashes := strings.Split(hashPlaintext, "\n")
			if len(hashes) == 0 {
				fmt.Println("No hash provided.")
				break
			}
			for _, hash := range hashes {
				if err := hashIdentifier(hash, true); err != nil {
					fmt.Printf("An error occurred: %v\n", err)
				}
			}
		case "6":
			// Wallet Balance
			clearScreen()
			if err := getWalletBalances(apiKey); err != nil {
				fmt.Printf("An error occurred: %v\n", err)
			}
		case "7":
			// Show Profit
			clearScreen()
			if err := getProfit(apiKey); err != nil {
				fmt.Printf("An error occurred: %v\n", err)
			}
		case "8":
			// Withdrawal History
			clearScreen()
			if err := withdrawalHistory(apiKey); err != nil {
				fmt.Printf("An error occurred: %v\n", err)
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
			fmt.Println("Invalid choice, please try again.")
			time.Sleep(500 * time.Millisecond)
		}
	}
}
