package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"
)

// download left lists
func downloadLeftList(apiKey string) error {
	fmt.Print("Download Left Lists:\n\n")

	var response struct {
		Success bool
		List    []struct {
			AlgorithmID int
			LeftList    string
		}
	}

	go func() {
		url := "https://hashes.com/en/api/jobs?key=" + apiKey
		resp, err := http.Get(url)
		if err != nil {
			fmt.Printf("error sending request: %v\n", err)
			return
		}
		defer resp.Body.Close()
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			fmt.Printf("error reading response: %v\n", err)
			return
		}
		if err := json.Unmarshal(body, &response); err != nil {
			fmt.Printf("error unmarshalling response: %v\n", err)
			return
		}
	}()

	fmt.Print("Enter hash mode (ex: 0 for MD5) or CTRL+C to cancel: ")
	var hashMode string
	fmt.Scanln(&hashMode)
	fmt.Println("This may appear to hang on large left lists...")

	for !response.Success {
		time.Sleep(100 * time.Millisecond)
	}

	var totalLeftLists int
	for _, item := range response.List {
		if fmt.Sprintf("%d", item.AlgorithmID) == hashMode {
			totalLeftLists++
		}
	}

	var completedDownloads int
	var progressMu sync.Mutex

	printProgressBar := func() {
		progressMu.Lock()
		defer progressMu.Unlock()
		percentage := float64(completedDownloads) / float64(totalLeftLists) * 100
		fmt.Printf("\rProgress: [")
		for i := 0; i < int(percentage/5); i++ {
			fmt.Print("=")
		}
		for i := int(percentage / 5); i < 20; i++ {
			fmt.Print(" ")
		}
		fmt.Printf("] %.2f%% completed", percentage)
	}

	var downloadedHashes sync.Map

	leftListChan := make(chan string, 100)

	downloadLeftList := func() {
		for leftListURL := range leftListChan {
			resp, err := http.Get("https://hashes.com" + leftListURL)
			if err != nil {
				fmt.Printf("error downloading left list: %v\n", err)
				continue
			}
			defer resp.Body.Close()
			content, err := io.ReadAll(resp.Body)
			if err != nil {
				fmt.Printf("error reading left list content: %v\n", err)
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
			downloadLeftList()
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
				outFile.WriteString(hash + "\n")
				totalUniqueHashes++
			}
		}
		return true
	})

	fmt.Printf("\n\nUnique hashes:\t%d\n", totalUniqueHashes)
	fmt.Printf("Saved to file:\t%s\n", outFileName)
	time.Sleep(1 * time.Second)
	return nil
}
