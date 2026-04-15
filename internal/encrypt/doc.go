// Package encrypt provides AES-GCM based encryption and decryption for
// sensitive values stored in .env files.
//
// # Overview
//
// Values are encrypted using a passphrase-derived AES-256 key and stored
// as base64-encoded ciphertext with the "enc:" prefix, making it easy to
// identify which entries are protected at a glance.
//
// # Usage
//
//	// Encrypt a single value
//	enc, err := encrypt.EncryptValue("my-secret", passphrase)
//
//	// Decrypt a previously encrypted value
//	plain, err := encrypt.DecryptValue(enc, passphrase)
//
//	// Encrypt all values in an EnvFile
//	encrypted, err := encrypt.EncryptFile(file, passphrase)
//
//	// Decrypt an EnvFile (skips non-encrypted values)
//	decrypted, err := encrypt.DecryptFile(encrypted, passphrase)
//
// # Security Notes
//
// Key derivation uses SHA-256; for production use consider PBKDF2 or Argon2.
// Each encryption call generates a fresh random nonce, so identical plaintexts
// produce different ciphertexts.
package encrypt
