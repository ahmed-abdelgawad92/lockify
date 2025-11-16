package service

import "context"

// EncryptionService provides encryption and decryption operations for vault entries
type EncryptionService interface {
	// Encrypt encrypts plaintext and returns base64-encoded ciphertext
	Encrypt(ctx context.Context, plaintext []byte) (string, error)
	// Decrypt decrypts base64-encoded ciphertext and returns plaintext
	Decrypt(ctx context.Context, ciphertext string) ([]byte, error)
}
