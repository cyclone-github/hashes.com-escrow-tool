package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net"
	"os"
	"strings"
	"sync"
	"time"
)

// download left lists
func downloadLeftList(apiKey string) error {
	fmt.Fprintln(os.Stderr, "Download Left Lists:")
	fmt.Fprintln(os.Stderr)

	var response struct {
		Success bool `json:"success"`
		List    []struct {
			AlgorithmID int    `json:"algorithmId"`
			LeftList    string `json:"leftList"`
		} `json:"list"`
	}

	// fetch jobs list (with timeout via global httpClient)
	url := "https://hashes.com/en/api/jobs?key=" + apiKey
	resp, err := httpClient.Get(url)
	if err != nil {
		var netErr net.Error
		if errors.As(err, &netErr) && netErr.Timeout() {
			fmt.Fprintln(os.Stderr, "Request timed out while fetching jobs list.")
		} else {
			fmt.Fprintf(os.Stderr, "Error sending request: %v\n", err)
		}
		return nil
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading response: %v\n", err)
		return nil
	}
	if err := json.Unmarshal(body, &response); err != nil {
		fmt.Fprintf(os.Stderr, "Error unmarshalling response: %v\n", err)
		return nil
	}
	if !response.Success {
		fmt.Fprintln(os.Stderr, "API did not return success for jobs list.")
		return nil
	}

	fmt.Fprint(os.Stderr, "Enter hash mode (ex: 0 for MD5) or CTRL+C to cancel: ")
	var hashMode string
	fmt.Scanln(&hashMode)
	fmt.Fprintln(os.Stderr, "This may appear to hang on large left lists...")

	// count how many left lists match this mode
	var totalLeftLists int
	for _, item := range response.List {
		if fmt.Sprintf("%d", item.AlgorithmID) == hashMode {
			totalLeftLists++
		}
	}

	if totalLeftLists == 0 {
		fmt.Fprintf(os.Stderr, "No left lists found for hash mode %s.\n", hashMode)
		return nil
	}

	var completedDownloads int
	var progressMu sync.Mutex

	printProgressBar := func() {
		progressMu.Lock()
		defer progressMu.Unlock()

		if totalLeftLists == 0 {
			return
		}

		percentage := float64(completedDownloads) / float64(totalLeftLists) * 100
		fmt.Fprint(os.Stderr, "\rProgress: [")
		for i := 0; i < int(percentage/5); i++ {
			fmt.Fprint(os.Stderr, "=")
		}
		for i := int(percentage / 5); i < 20; i++ {
			fmt.Fprint(os.Stderr, " ")
		}
		fmt.Fprintf(os.Stderr, "] %.2f%% completed", percentage)
	}

	var downloadedHashes sync.Map
	leftListChan := make(chan string, 100)

	downloadLeftListWorker := func() {
		for leftListURL := range leftListChan {
			resp, err := httpClient.Get("https://hashes.com" + leftListURL)
			if err != nil {
				var netErr net.Error
				if errors.As(err, &netErr) && netErr.Timeout() {
					fmt.Fprintf(os.Stderr, "\nRequest timed out for left list %s\n", leftListURL)
				} else {
					fmt.Fprintf(os.Stderr, "\nError downloading left list %s: %v\n", leftListURL, err)
				}
				continue
			}

			content, err := io.ReadAll(resp.Body)
			resp.Body.Close()
			if err != nil {
				fmt.Fprintf(os.Stderr, "\nError reading left list content for %s: %v\n", leftListURL, err)
				continue
			}

			downloadedHashes.Store(leftListURL, string(content))

			progressMu.Lock()
			completedDownloads++
			progressMu.Unlock()
			printProgressBar()
		}
	}

	printProgressBar()

	var wg sync.WaitGroup
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			downloadLeftListWorker()
		}()
	}

	for _, item := range response.List {
		if fmt.Sprintf("%d", item.AlgorithmID) == hashMode {
			leftListChan <- item.LeftList
		}
	}

	close(leftListChan)
	wg.Wait()

	outFileName := fmt.Sprintf("hashes_%s_left.txt", hashMode)
	outFile, err := os.Create(outFileName)
	if err != nil {
		return fmt.Errorf("error creating output file: %v", err)
	}
	defer outFile.Close()

	uniqueHashes := make(map[string]bool)
	var totalUniqueHashes int
	downloadedHashes.Range(func(_, value interface{}) bool {
		hashes := value.(string)
		for _, hash := range strings.Split(hashes, "\n") {
			if hash != "" && !uniqueHashes[hash] {
				uniqueHashes[hash] = true
				_, _ = outFile.WriteString(hash + "\n")
				totalUniqueHashes++
			}
		}
		return true
	})

	fmt.Fprintf(os.Stderr, "\n\nUnique hashes:\t%d\n", totalUniqueHashes)
	fmt.Fprintf(os.Stderr, "Saved to file:\t%s\n", outFileName)
	time.Sleep(1 * time.Second)
	return nil
}
