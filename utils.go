package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

// parse string to float64
func parseFloat(value string) float64 {
	f, _ := strconv.ParseFloat(value, 64)
	return f
}

// paste hashes
func pasteHashes(action string) string {
	fmt.Fprintln(os.Stderr, action)
	fmt.Fprintln(os.Stderr, "Paste one hash per line (press Enter twice to finish):")

	reader := bufio.NewReader(os.Stdin)
	var hashPlaintexts []string

	for {
		line, err := reader.ReadString('\n')
		if err != nil || strings.TrimSpace(line) == "" {
			break
		}
		hashPlaintexts = append(hashPlaintexts, strings.TrimSpace(line))
	}

	return strings.Join(hashPlaintexts, "\n")
}

// select file
func selectFile() (string, string) {
	files, err := filepath.Glob("*.txt")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading files: %v\n", err)
		return "", ""
	}

	for {
		fmt.Fprintln(os.Stderr, "Select a file to upload:")
		for i, file := range files {
			fmt.Fprintf(os.Stderr, "%d. %s\n", i+1, file)
		}
		fmt.Fprintln(os.Stderr, "p. Paste hash:plaintext")
		fmt.Fprintln(os.Stderr, "c. Custom file path")
		fmt.Fprintln(os.Stderr, "m. Return to Menu")

		var input string
		fmt.Scanln(&input)
		input = strings.ToLower(input)

		if input == "m" {
			clearScreen()
			return "", "" // return empty strings to return to the menu
		}

		if input == "p" {
			return "PASTE", pasteHashes("")
		}

		if input == "c" {
			fmt.Fprintln(os.Stderr, "Enter the full path of the file:")
			var customFilePath string
			fmt.Scanln(&customFilePath)
			return customFilePath, ""
		}

		choice, err := strconv.Atoi(input)
		if err == nil && choice > 0 && choice <= len(files) {
			return files[choice-1], ""
		}

		fmt.Fprintln(os.Stderr, "Invalid selection. Try again.")
	}
}
