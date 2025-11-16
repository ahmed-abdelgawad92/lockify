package crypto

import (
	"context"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io"
	"runtime"

	"github.com/apixify/lockify/internal/config"
	"github.com/apixify/lockify/internal/domain/service"
	"golang.org/x/crypto/argon2"
)

// EncryptionService implements domain.EncryptionService
type AESEncryptionService struct {
	aead cipher.AEAD
	cfg  config.EncryptionConfig
	key  []byte
	salt []byte
}

// NewAESEncryptionService creates a new encryption service instance
func NewAESEncryptionService(encodedSalt string, passphrase string, cfg config.EncryptionConfig) (service.EncryptionService, error) {
	if encodedSalt == "" {
		return nil, fmt.Errorf("salt cannot be empty")
	}
	if passphrase == "" {
		return nil, fmt.Errorf("passphrase cannot be empty")
	}

	salt, err := base64.StdEncoding.DecodeString(encodedSalt)
	if err != nil {
		return nil, fmt.Errorf("invalid salt encoding: %w", err)
	}
	if len(salt) == 0 {
		return nil, fmt.Errorf("salt cannot be empty")
	}

	key := deriveKey([]byte(passphrase), salt, cfg)

	block, err := aes.NewCipher(key)
	if err != nil {
		clearBytes(key, salt)
		return nil, fmt.Errorf("failed to create cipher: %w", err)
	}

	aead, err := cipher.NewGCM(block)
	if err != nil {
		clearBytes(key, salt)
		return nil, fmt.Errorf("failed to create GCM: %w", err)
	}

	service := &AESEncryptionService{aead, cfg, key, salt}
	runtime.SetFinalizer(service, (*AESEncryptionService).cleanup)

	return service, nil
}

// Encrypt encrypts plaintext and returns base64-encoded ciphertext
func (e *AESEncryptionService) Encrypt(ctx context.Context, plaintext []byte) (string, error) {
	if plaintext == nil {
		return "", fmt.Errorf("plaintext cannot be nil")
	}

	nonce := make([]byte, e.cfg.NonceSize)
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", fmt.Errorf("failed to generate nonce: %w", err)
	}

	ciphertext := e.aead.Seal(nil, nonce, plaintext, nil)
	result := append(nonce, ciphertext...)
	encoded := base64.StdEncoding.EncodeToString(result)

	clearBytes(nonce, ciphertext)

	return encoded, nil
}

// Decrypt decrypts base64-encoded ciphertext and returns plaintext
func (e *AESEncryptionService) Decrypt(ctx context.Context, ciphertext string) ([]byte, error) {
	if ciphertext == "" {
		return nil, fmt.Errorf("ciphertext cannot be empty")
	}

	raw, err := base64.StdEncoding.DecodeString(ciphertext)
	if err != nil {
		return nil, fmt.Errorf("invalid ciphertext encoding: %w", err)
	}

	if err := e.validateCiphertextLength(raw); err != nil {
		return nil, err
	}

	// Extract nonce and ciphertext
	nonce := raw[:e.cfg.NonceSize]
	ciphertextBytes := raw[e.cfg.NonceSize:]
	plaintext, err := e.aead.Open(nil, nonce, ciphertextBytes, nil)
	if err != nil {
		clearBytes(nonce, ciphertextBytes)
		return nil, fmt.Errorf("decryption failed: %w", err)
	}

	clearBytes(nonce, ciphertextBytes)

	return plaintext, nil
}

// validateCiphertextLength checks if the ciphertext meets the minimum length requirement
// The minimum length is nonce size + AEAD overhead (authentication tag)
func (e *AESEncryptionService) validateCiphertextLength(ciphertext []byte) error {
	minLen := e.cfg.NonceSize + e.aead.Overhead()
	if len(ciphertext) < minLen {
		return fmt.Errorf("ciphertext too short: expected at least %d bytes, got %d", minLen, len(ciphertext))
	}
	return nil
}

// deriveKey derives a key from a passphrase using Argon2id
func deriveKey(passphrase []byte, salt []byte, cfg config.EncryptionConfig) []byte {
	return argon2.IDKey(passphrase, salt, cfg.ArgonTime, cfg.ArgonMemory, cfg.ArgonThreads, cfg.KeyLength)
}

// cleanup clears sensitive data from memory
func (e *AESEncryptionService) cleanup() {
	if e.key != nil {
		clearBytes(e.key)
	}
	if e.salt != nil {
		clearBytes(e.salt)
	}
}

// clearBytes clears sensitive data from memory by overwrites a byte slice with zeros
func clearBytes(b ...[]byte) {
	if b == nil {
		return
	}
	for _, eachB := range b {
		for i := range eachB {
			eachB[i] = 0
		}
	}
	runtime.KeepAlive(b)
}
