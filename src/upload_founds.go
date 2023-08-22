package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"time"
)

// case 1, upload founds
func uploadFounds(apiKey string) error {
	clearScreen()
	fmt.Println("Upload Founds:\n")
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

	fmt.Println("Enter algorithm ID:")
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
	_, err = io.Copy(filePart, file)
	if err != nil {
		return fmt.Errorf("An error occurred: failed to copy file: %v", err)
	}

	err = writer.Close()
	if err != nil {
		return fmt.Errorf("An error occurred: failed to close writer: %v", err)
	}

	url := "https://hashes.com/en/api/founds"
	req, err := http.NewRequest(http.MethodPost, url, &requestBody)
	if err != nil {
		return fmt.Errorf("An error occurred: failed to create request: %v", err)
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("An error occurred: failed to send request: %v", err)
	}
	defer resp.Body.Close()

	var response struct {
		Success bool `json:"success"`
	}
	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		return fmt.Errorf("An error occurred: failed to decode response: %v", err)
	}

	if !response.Success {
		return fmt.Errorf("Upload failed")
	}

	fmt.Println("Upload Successful")

	return nil
}