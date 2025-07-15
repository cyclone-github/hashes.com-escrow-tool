package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

// hash identifier
func hashIdentifier(hash string, extended bool) error {
	url := "https://hashes.com/en/api/identifier?hash=" + hash
	if extended {
		url += "&extended=true"
	}

	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		return err
	}

	if success, ok := result["success"].(bool); ok && success {
		if algorithms, ok := result["algorithms"].([]interface{}); ok && len(algorithms) > 0 {
			fmt.Fprintln(os.Stderr, "Possible Algorithms:")
			for _, algo := range algorithms {
				fmt.Fprintf(os.Stderr, "  %s\n", algo)
			}
		} else {
			fmt.Fprintln(os.Stderr, "No algorithms found.")
		}
	} else {
		fmt.Fprintln(os.Stderr, "No results found.")
	}

	return nil
}
