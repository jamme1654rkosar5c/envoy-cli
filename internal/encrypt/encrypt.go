// Package encrypt provides utilities for encrypting and decrypting
// sensitive values within .env files using AES-GCM symmetric encryption.
package encrypt

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"strings"

	"github.com/yourusername/envoy-cli/internal/parser"
)

const encryptedPrefix = "enc:"

// deriveKey produces a 32-byte AES-256 key from a passphrase using SHA-256.
func deriveKey(passphrase string) []byte {
	hash := sha256.Sum256([]byte(passphrase))
	return hash[:]
}

// EncryptValue encrypts a plaintext string and returns a base64-encoded
// ciphertext prefixed with "enc:".
func EncryptValue(plaintext, passphrase string) (string, error) {
	key := deriveKey(passphrase)
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", fmt.Errorf("encrypt: create cipher: %w", err)
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", fmt.Errorf("encrypt: create GCM: %w", err)
	}
	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", fmt.Errorf("encrypt: generate nonce: %w", err)
	}
	ciphertext := gcm.Seal(nonce, nonce, []byte(plaintext), nil)
	return encryptedPrefix + base64.StdEncoding.EncodeToString(ciphertext), nil
}

// DecryptValue decrypts a value previously encrypted by EncryptValue.
// Returns an error if the value is not prefixed with "enc:".
func DecryptValue(encoded, passphrase string) (string, error) {
	if !strings.HasPrefix(encoded, encryptedPrefix) {
		return "", errors.New("decrypt: value is not encrypted")
	}
	data, err := base64.StdEncoding.DecodeString(strings.TrimPrefix(encoded, encryptedPrefix))
	if err != nil {
		return "", fmt.Errorf("decrypt: base64 decode: %w", err)
	}
	key := deriveKey(passphrase)
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", fmt.Errorf("decrypt: create cipher: %w", err)
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", fmt.Errorf("decrypt: create GCM: %w", err)
	}
	if len(data) < gcm.NonceSize() {
		return "", errors.New("decrypt: ciphertext too short")
	}
	nonce, ciphertext := data[:gcm.NonceSize()], data[gcm.NonceSize():]
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return "", fmt.Errorf("decrypt: open GCM: %w", err)
	}
	return string(plaintext), nil
}

// IsEncrypted reports whether a value has the encrypted prefix.
func IsEncrypted(value string) bool {
	return strings.HasPrefix(value, encryptedPrefix)
}

// EncryptFile returns a new EnvFile where all entry values are encrypted.
func EncryptFile(file parser.EnvFile, passphrase string) (parser.EnvFile, error) {
	out := parser.EnvFile{Path: file.Path}
	for _, entry := range file.Entries {
		enc, err := EncryptValue(entry.Value, passphrase)
		if err != nil {
			return parser.EnvFile{}, fmt.Errorf("encrypt file: key %q: %w", entry.Key, err)
		}
		out.Entries = append(out.Entries, parser.Entry{Key: entry.Key, Value: enc, Comment: entry.Comment})
	}
	return out, nil
}

// DecryptFile returns a new EnvFile where all encrypted values are decrypted.
func DecryptFile(file parser.EnvFile, passphrase string) (parser.EnvFile, error) {
	out := parser.EnvFile{Path: file.Path}
	for _, entry := range file.Entries {
		value := entry.Value
		if IsEncrypted(value) {
			dec, err := DecryptValue(value, passphrase)
			if err != nil {
				return parser.EnvFile{}, fmt.Errorf("decrypt file: key %q: %w", entry.Key, err)
			}
			value = dec
		}
		out.Entries = append(out.Entries, parser.Entry{Key: entry.Key, Value: value, Comment: entry.Comment})
	}
	return out, nil
}
