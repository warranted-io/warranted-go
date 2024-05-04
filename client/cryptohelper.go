package client

import (
	"crypto/hmac"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/hex"
	"errors"
	"hash"
)

// TimeSafeCompare does a time-safe comparison of two strings
func TimeSafeCompare(a, b string) bool {
	if len(a) != len(b) {
		return false
	}
	var result uint8
	for i := 0; i < len(a); i++ {
		result |= uint8(a[i] ^ b[i])
	}
	return result == 0
}

// CreateHMAC creates an HMAC of a JSON request
func CreateHMAC(url, jsonData, secretKey string, algorithm string) (string, error) {
	var h func() hash.Hash
	switch algorithm {
	case "sha256":
		h = sha256.New
	case "sha512":
		h = sha512.New
	default:
		return "", errors.New("unsupported hashing algorithm")
	}
	mac := hmac.New(h, []byte(secretKey))
	_, err := mac.Write([]byte(url + jsonData))
	if err != nil {
		return "", err
	}
	sum := mac.Sum(nil)
	return hex.EncodeToString(sum), nil
}
