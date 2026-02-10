package dnslog

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"io"
	"strings"

	"github.com/genwilliam/dnslog_for_go/config"
)

const secretPrefix = "enc:"

var ErrSecretKeyRequired = errors.New("webhook_secret_key_required")

func EncryptWebhookSecret(plain string) (string, error) {
	if plain == "" {
		return "", nil
	}
	key, err := loadSecretKey()
	if err != nil {
		return "", err
	}
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}
	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}
	cipherText := gcm.Seal(nil, nonce, []byte(plain), nil)
	out := append(nonce, cipherText...)
	return secretPrefix + base64.StdEncoding.EncodeToString(out), nil
}

func DecryptWebhookSecret(stored string) (string, error) {
	if stored == "" {
		return "", nil
	}
	if !strings.HasPrefix(stored, secretPrefix) {
		return stored, nil
	}
	key, err := loadSecretKey()
	if err != nil {
		return "", err
	}
	raw, err := base64.StdEncoding.DecodeString(strings.TrimPrefix(stored, secretPrefix))
	if err != nil {
		return "", err
	}
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}
	if len(raw) < gcm.NonceSize() {
		return "", errors.New("ciphertext too short")
	}
	nonce := raw[:gcm.NonceSize()]
	cipherText := raw[gcm.NonceSize():]
	plain, err := gcm.Open(nil, nonce, cipherText, nil)
	if err != nil {
		return "", err
	}
	return string(plain), nil
}

func loadSecretKey() ([]byte, error) {
	cfg := config.Get()
	if cfg == nil || cfg.WebhookSecretKey == "" {
		return nil, ErrSecretKeyRequired
	}
	keyStr := cfg.WebhookSecretKey
	if b, err := base64.StdEncoding.DecodeString(keyStr); err == nil {
		if len(b) == 32 {
			return b, nil
		}
	}
	if b, err := hex.DecodeString(keyStr); err == nil {
		if len(b) == 32 {
			return b, nil
		}
	}
	return nil, errors.New("invalid webhook secret key")
}
