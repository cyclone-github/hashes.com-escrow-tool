package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

// case 4, hash identifier
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

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		return err
	}

	if success, ok := result["success"].(bool); ok && success {
		if algorithms, ok := result["algorithms"].([]interface{}); ok && len(algorithms) > 0 {
			fmt.Println("Possible Algorithms:")
			for _, algo := range algorithms {
				fmt.Printf("  %s\n", algo)
			}
		} else {
			fmt.Println("No algorithms found.")
		}
	} else {
		fmt.Println("No results found.")
	}

	return nil
}