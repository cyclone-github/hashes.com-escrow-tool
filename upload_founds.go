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
	"time"
)

// upload founds
func uploadFounds(apiKey string) error {
	fmt.Fprintln(os.Stderr, "Upload Founds:")
	fmt.Fprintln(os.Stderr)

	filePath, hashPlaintext := selectFile()

	var file io.Reader
	var filename string

	if filePath == "PASTE" {
		file = bytes.NewBufferString(hashPlaintext)
		filename = time.Now().Format("150405") + ".txt"
	} else if filePath == "" {
		return nil
	} else {
		f, err := os.Open(filePath)
		if err != nil {
			return fmt.Errorf("An error occurred: failed to open file: %v", err)
		}
		defer f.Close()
		file = f
		filename = filePath
	}

	fmt.Fprintln(os.Stderr, "Enter algorithm ID:")
	var algo string
	fmt.Scanln(&algo)

	var requestBody bytes.Buffer
	writer := multipart.NewWriter(&requestBody)

	_ = writer.WriteField("key", apiKey)
	_ = writer.WriteField("algo", algo)

	filePart, err := writer.CreateFormFile("userfile", filename)
	if err != nil {
		return fmt.Errorf("An error occurred: failed to create form file: %v", err)
	}
	if _, err = io.Copy(filePart, file); err != nil {
		return fmt.Errorf("An error occurred: failed to copy file: %v", err)
	}

	if err = writer.Close(); err != nil {
		return fmt.Errorf("An error occurred: failed to close writer: %v", err)
	}

	url := "https://hashes.com/en/api/founds"
	req, err := http.NewRequest(http.MethodPost, url, &requestBody)
	if err != nil {
		return fmt.Errorf("An error occurred: failed to create request: %v", err)
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())

	resp, err := httpClient.Do(req)
	if err != nil {
		var netErr net.Error
		if errors.As(err, &netErr) && netErr.Timeout() {
			fmt.Fprintln(os.Stderr, "Request timed out while uploading founds.")
			return nil // non-fatal
		}
		return fmt.Errorf("An error occurred: failed to send request: %v", err)
	}
	defer resp.Body.Close()

	var response struct {
		Success bool `json:"success"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return fmt.Errorf("An error occurred: failed to decode response: %v", err)
	}

	if !response.Success {
		return fmt.Errorf("Upload failed")
	}

	fmt.Fprintln(os.Stderr, "Upload Successful")

	return nil
}
