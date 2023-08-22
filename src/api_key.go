package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
)

// generate unique encryption key
func getEncryptionKey() {
	decodedSeed, err := base64.StdEncoding.DecodeString(base64StaticSeedKey)
	if err != nil {
		fmt.Printf("Error decoding Base64 static seed key: %v\n")
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

	resp, err := http.Get(url)
	if err != nil {
		fmt.Printf("Failed to send request: %v\n", err)
		return false
	}
	defer resp.Body.Close()

	var response struct {
		Success bool `json:"success"`
	}
	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		fmt.Printf("Failed to decode response: %v\n", err)
		return false
	}

	if response.Success {
		fmt.Println("API key verified")
		return true
	}

	//fmt.Println("API key not verified, try again")
	return false
}
