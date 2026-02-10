package dnslog

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
)

func HashAPIKey(key string) string {
	sum := sha256.Sum256([]byte(key))
	return hex.EncodeToString(sum[:])
}

func GenerateAPIKey() (string, string, error) {
	raw := make([]byte, 32)
	if _, err := rand.Read(raw); err != nil {
		return "", "", err
	}
	plain := hex.EncodeToString(raw)
	return plain, HashAPIKey(plain), nil
}
