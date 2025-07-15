package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
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
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
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

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 1, ' ', 0) // Removed tabwriter.Debug
	defer w.Flush()

	founds, found := results["founds"]
	if found && len(founds.([]interface{})) > 0 {
		fmt.Fprintln(w, "Found:")
		fmt.Fprintln(w, "Hash\tPlain\tAlgorithm\t")
		foundsSlice := founds.([]interface{})
		for _, foundItem := range foundsSlice {
			foundMap, ok := foundItem.(map[string]interface{})
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

	unfounds, unfound := results["unfounds"]
	if unfound && len(unfounds.([]interface{})) > 0 {
		fmt.Fprintln(w, "\nNot Found:")
		for _, unfoundItem := range unfounds.([]interface{}) {
			unfoundMap, _ := unfoundItem.(map[string]interface{})
			fmt.Fprintf(w, "%s\t\n", unfoundMap["hash"].(string))
		}
	}

	if cost, ok := results["cost"].(float64); ok {
		fmt.Fprintf(w, "\nCredits used: %d\n", int(cost))
	}

	return nil
}
