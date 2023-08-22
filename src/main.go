package main

import (
	"fmt"
	"strings"
	"time"
)

// version history
// v0.1.0; 2023-08-22.1100; initial github release

// main function
func main() {
	clearScreen()
	printCyclone()
	fmt.Println(" ######################################################################")
	fmt.Println("#              Cyclone's Hashes.com API Escrow Tool v0.1.0             #")
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
		fmt.Println("3.  Search Hashes")
		fmt.Println("4.  Hash Identifier")
		fmt.Println("5.  Wallet Balance")
		fmt.Println("6.  Show Profit")
		fmt.Println("7.  Withdrawal History")
		fmt.Println("8.  Enter New API")
		fmt.Println("9.  Remove API Key")
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
			// Search Hashes
			clearScreen()
			hashes := getHashesFromUser()
			if err := searchHashes(apiKey, hashes); err != nil {
				fmt.Printf("An error occurred: %v\n", err)
			}		
		case "4":
			// Hash Identifier
			clearScreen()
			hashes := getHashesFromUser()
			if len(hashes) == 0 {
				fmt.Println("No hash provided.")
				break
			}
			hash := hashes[0]
			if err := hashIdentifier(hash, true); err != nil {
				fmt.Printf("An error occurred: %v\n", err)
			}
		case "5":
			// Wallet Balance
			clearScreen()
			if err := getWalletBalances(apiKey); err != nil {
				fmt.Printf("An error occurred: %v\n", err)
			}
		case "6":
			// Show Profit
			clearScreen()
			if err := getProfit(apiKey); err != nil {
				fmt.Printf("An error occurred: %v\n", err)
			}
		case "7":
			// Withdrawal History
			clearScreen()
			if err := withdrawalHistory(apiKey); err != nil {
				fmt.Printf("An error occurred: %v\n", err)
			}
		case "8":
			// Enter New API
			getAPIKey(true)
		case "9":
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