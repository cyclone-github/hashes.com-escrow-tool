package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"
	"text/tabwriter"
)

// case 3, search hashes
func searchHashes(apiKey string, hashes []string) error {
	var requestBody bytes.Buffer
	writer := multipart.NewWriter(&requestBody)
	_ = writer.WriteField("key", apiKey)
	for _, hash := range hashes {
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

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	var results map[string]interface{}
	if err := json.Unmarshal(body, &results); err != nil {
		return err
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 1, ' ', tabwriter.Debug)
	defer w.Flush()

	founds, found := results["founds"]
	if !found {
		fmt.Println("No results found.")
		return nil
	}

	foundsSlice, ok := founds.([]interface{})
	if !ok {
		fmt.Println("Unexpected results format.")
		return nil
	}

	if len(foundsSlice) > 0 {
		for _, foundItem := range foundsSlice {
			foundMap, ok := foundItem.(map[string]interface{})
			if !ok {
				fmt.Println("Unexpected item format.")
				continue
			}

			for key, value := range foundMap {
				if key != "salt" || fmt.Sprint(value) != "" {
					fmt.Fprintf(w, "%s:\t%s\n", key, fmt.Sprint(value))
				}
			}
			fmt.Println()
		}
	} else {
		unfounds, unfound := results["unfounds"]
		if unfound {
			unfoundHashes, ok := unfounds.([]interface{})
			if ok && len(unfoundHashes) > 0 {
				fmt.Println("Hash not found.")
			}
		}
	}

	return nil
}