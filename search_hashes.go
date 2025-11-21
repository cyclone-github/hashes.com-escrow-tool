package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net"
	"net/http"
	"os"
	"text/tabwriter"
)

// search hashes
func searchHashes(apiKey string, hashes []string) error {
	var requestBody bytes.Buffer
	writer := multipart.NewWriter(&requestBody)

	_ = writer.WriteField("key", apiKey)
	for _, hash := range hashes {
		fmt.Fprintf(os.Stderr, "Searching hashes.com for %s...\n", hash)
		_ = writer.WriteField("hashes[]", hash)
	}
	writer.Close()

	req, err := http.NewRequest("POST", "https://hashes.com/en/api/search", &requestBody)
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", writer.FormDataContentType())

	resp, err := httpClient.Do(req)
	if err != nil {
		var netErr net.Error
		if errors.As(err, &netErr) && netErr.Timeout() {
			fmt.Fprintln(os.Stderr, "Request timed out while searching hashes.")
			return nil // non-fatal
		}
		return err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	var results map[string]interface{}
	if err := json.Unmarshal(body, &results); err != nil {
		return err
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 1, ' ', 0)
	defer w.Flush()

	// found hashes
	if founds, ok := results["founds"]; ok {
		if foundsSlice, ok := founds.([]interface{}); ok && len(foundsSlice) > 0 {
			fmt.Fprintln(w, "Found:")
			fmt.Fprintln(w, "Hash\tPlain\tAlgorithm\t")
			for _, item := range foundsSlice {
				foundMap, ok := item.(map[string]interface{})
				if !ok {
					fmt.Fprintln(os.Stderr, "Unexpected item format.")
					continue
				}

				hash, _ := foundMap["hash"].(string)
				plain, _ := foundMap["plaintext"].(string)
				algo, _ := foundMap["algorithm"].(string)
				fmt.Fprintf(w, "%s\t%s\t%s\t\n", hash, plain, algo)
			}
		}
	}

	// unfound hashes
	if unfounds, ok := results["unfounds"]; ok {
		if unfoundsSlice, ok := unfounds.([]interface{}); ok && len(unfoundsSlice) > 0 {
			fmt.Fprintln(w, "\nNot Found:")
			for _, item := range unfoundsSlice {
				unfoundMap, ok := item.(map[string]interface{})
				if !ok {
					continue
				}
				hash, _ := unfoundMap["hash"].(string)
				fmt.Fprintf(w, "%s\t\n", hash)
			}
		}
	}

	// credits used
	if cost, ok := results["cost"].(float64); ok {
		fmt.Fprintf(w, "\nCredits used: %d\n", int(cost))
	}

	return nil
}
