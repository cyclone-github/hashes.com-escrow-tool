package main

import (
	"fmt"
	"os"
	"strings"
)

// case 9, remove API key
func removeAPIKey() {
	fmt.Println("Are you sure you want to remove your API key? (y/n):")
	var confirmation string
	fmt.Scanln(&confirmation)

	if strings.ToLower(confirmation) == "y" {
		fileInfo, err := os.Stat(apiKeyFile)
		if err != nil {
			fmt.Printf("Error getting API key file info: %v\n", err)
			return
		}

		file, err := os.OpenFile(apiKeyFile, os.O_WRONLY, 0644)
		if err != nil {
			fmt.Printf("Error opening API key file: %v\n", err)
			return
		}

		zeroBytes := make([]byte, fileInfo.Size())
		_, err = file.Write(zeroBytes)
		if err != nil {
			fmt.Printf("Error overwriting API key file with zeroes: %v\n", err)
			return
		}

		file.Close()

		err = os.Remove(apiKeyFile)
		if err != nil {
			fmt.Printf("Error removing API key file: %v\n", err)
		} else {
			fmt.Println("API key removed successfully. Exiting...")
			os.Exit(0)
		}
	} else {
		fmt.Println("API key removal canceled.")
	}
}