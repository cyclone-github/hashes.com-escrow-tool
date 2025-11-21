package main

import (
	"context"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net"
	"os"
)

// generate unique encryption key
func getEncryptionKey() {
	decodedSeed, err := base64.StdEncoding.DecodeString(base64StaticSeedKey)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error decoding Base64 static seed key: %v\n", err)
		return
	}

	staticSeedKey := string(decodedSeed)

	interfaces, err := net.Interfaces()
	if err != nil || len(interfaces) == 0 {
		hash := sha256.Sum256([]byte(staticSeedKey))
		encryptionKey = hex.EncodeToString(hash[:])
		return
	}

	var macAddr string
	for _, i := range interfaces {
		if i.HardwareAddr != nil && len(i.HardwareAddr.String()) > 0 {
			macAddr = i.HardwareAddr.String()
			break
		}
	}

	seed := macAddr + staticSeedKey
	hash := sha256.Sum256([]byte(seed))
	encryptionKey = hex.EncodeToString(hash[:])
}

// encrypt key file
func encrypt(data []byte, passphrase string) ([]byte, error) {
	block, err := aes.NewCipher(getSHA256Hash(passphrase))
	if err != nil {
		return nil, err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}

	ciphertext := gcm.Seal(nonce, nonce, data, nil)
	return ciphertext, nil
}

// decrypt key file
func decrypt(data []byte, passphrase string) ([]byte, error) {
	key := getSHA256Hash(passphrase)
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonceSize := gcm.NonceSize()
	if len(data) < nonceSize {
		return nil, fmt.Errorf("ciphertext too short")
	}

	nonce, ciphertext := data[:nonceSize], data[nonceSize:]
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, err
	}
	return plaintext, nil
}

// sha256Hash function for encrypt /decrypt file func
func getSHA256Hash(text string) []byte {
	hash := sha256.Sum256([]byte(text))
	return hash[:]
}

// verify hashes.com API key
func verifyAPIKey(apiKey string) bool {
	url := fmt.Sprintf("https://hashes.com/en/api/balance?key=%s", apiKey)

	resp, err := httpClient.Get(url)
	if err != nil {
		var netErr net.Error
		switch {
		case errors.Is(err, context.DeadlineExceeded):
			fmt.Fprintln(os.Stderr, "Request timed out while verifying API key.")
		case errors.As(err, &netErr) && netErr.Timeout():
			fmt.Fprintln(os.Stderr, "Request timed out while verifying API key.")
		default:
			fmt.Fprintf(os.Stderr, "Failed to send request: %v\n", err)
		}
		return false
	}
	defer resp.Body.Close()

	var response struct {
		Success bool `json:"success"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to decode response: %v\n", err)
		return false
	}

	if response.Success {
		fmt.Fprintln(os.Stderr, "API key verified")
		return true
	}

	return false
}
