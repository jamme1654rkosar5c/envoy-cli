package encrypt_test

import (
	"strings"
	"testing"

	"github.com/yourusername/envoy-cli/internal/encrypt"
	"github.com/yourusername/envoy-cli/internal/parser"
)

const testPassphrase = "super-secret-passphrase"

func makeEnvFile(entries []parser.Entry) parser.EnvFile {
	return parser.EnvFile{Path: ".env", Entries: entries}
}

func TestEncryptValue_ProducesEncPrefix(t *testing.T) {
	enc, err := encrypt.EncryptValue("hello", testPassphrase)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.HasPrefix(enc, "enc:") {
		t.Errorf("expected enc: prefix, got %q", enc)
	}
}

func TestEncryptDecrypt_RoundTrip(t *testing.T) {
	original := "my-secret-value"
	enc, err := encrypt.EncryptValue(original, testPassphrase)
	if err != nil {
		t.Fatalf("encrypt: %v", err)
	}
	dec, err := encrypt.DecryptValue(enc, testPassphrase)
	if err != nil {
		t.Fatalf("decrypt: %v", err)
	}
	if dec != original {
		t.Errorf("expected %q, got %q", original, dec)
	}
}

func TestDecryptValue_WrongPassphrase(t *testing.T) {
	enc, _ := encrypt.EncryptValue("secret", testPassphrase)
	_, err := encrypt.DecryptValue(enc, "wrong-passphrase")
	if err == nil {
		t.Error("expected error for wrong passphrase, got nil")
	}
}

func TestDecryptValue_NotEncrypted(t *testing.T) {
	_, err := encrypt.DecryptValue("plaintext", testPassphrase)
	if err == nil {
		t.Error("expected error for non-encrypted value")
	}
}

func TestIsEncrypted(t *testing.T) {
	if encrypt.IsEncrypted("plain") {
		t.Error("plain value should not be detected as encrypted")
	}
	enc, _ := encrypt.EncryptValue("val", testPassphrase)
	if !encrypt.IsEncrypted(enc) {
		t.Error("encrypted value should be detected as encrypted")
	}
}

func TestEncryptFile_AllValuesEncrypted(t *testing.T) {
	file := makeEnvFile([]parser.Entry{
		{Key: "DB_PASS", Value: "secret1"},
		{Key: "API_KEY", Value: "secret2"},
	})
	enc, err := encrypt.EncryptFile(file, testPassphrase)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	for _, entry := range enc.Entries {
		if !encrypt.IsEncrypted(entry.Value) {
			t.Errorf("key %q value not encrypted", entry.Key)
		}
	}
}

func TestDecryptFile_RestoresOriginalValues(t *testing.T) {
	original := makeEnvFile([]parser.Entry{
		{Key: "DB_PASS", Value: "secret1"},
		{Key: "API_KEY", Value: "secret2"},
	})
	encrypted, _ := encrypt.EncryptFile(original, testPassphrase)
	decrypted, err := encrypt.DecryptFile(encrypted, testPassphrase)
	if err != nil {
		t.Fatalf("decrypt file: %v", err)
	}
	for i, entry := range decrypted.Entries {
		if entry.Value != original.Entries[i].Value {
			t.Errorf("key %q: expected %q, got %q", entry.Key, original.Entries[i].Value, entry.Value)
		}
	}
}

func TestDecryptFile_SkipsPlainValues(t *testing.T) {
	file := makeEnvFile([]parser.Entry{
		{Key: "APP_ENV", Value: "production"},
	})
	out, err := encrypt.DecryptFile(file, testPassphrase)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out.Entries[0].Value != "production" {
		t.Errorf("expected plain value to pass through, got %q", out.Entries[0].Value)
	}
}
