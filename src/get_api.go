package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"time"
)

// case 8, get new API key
func getAPIKey(promptForNewKey bool) string {
	if promptForNewKey {
		clearScreen()
		for {
			fmt.Println("Enter your API key from hashes.com/profile")
			var newAPIKey string
			fmt.Scanln(&newAPIKey)

			if verifyAPIKey(newAPIKey) {
				encryptedKey, err := encrypt([]byte(newAPIKey), encryptionKey)
				if err != nil {
					fmt.Printf("Error encrypting API key: %v\n", err)
					continue
				}
				err = ioutil.WriteFile(apiKeyFile, encryptedKey, 0644)
				if err != nil {
					fmt.Printf("Error writing encrypted API key to file: %v\n", err)
					continue
				}
				return newAPIKey
			} else {
				fmt.Println("API key not verified, try again.")
				time.Sleep(500 * time.Millisecond)
			}
		}
	}

	for {
		var apiKey string

		if _, err := os.Stat(apiKeyFile); err == nil {
			encryptedKey, err := ioutil.ReadFile(apiKeyFile)
			if err != nil {
				fmt.Printf("Error reading API key file: %v\n", err)
			} else {
				decryptedKey, err := decrypt(encryptedKey, encryptionKey)
				if err != nil {
					fmt.Printf("Error decrypting API key: %v\n", err)
				} else {
					apiKey = string(decryptedKey)
					if verifyAPIKey(apiKey) {
						return apiKey
					}
				}
			}
		}

		fmt.Println("Enter your API key from hashes.com/profile")
		fmt.Scanln(&apiKey)
		encryptedKey, err := encrypt([]byte(apiKey), encryptionKey)
		if err != nil {
			fmt.Printf("Error encrypting API key: %v\n", err)
			continue
		}
		err = ioutil.WriteFile(apiKeyFile, encryptedKey, 0644)
		if err != nil {
			fmt.Printf("Error writing encrypted API key to file: %v\n", err)
			continue
		}
		if verifyAPIKey(apiKey) {
			return apiKey
		}
	}
}