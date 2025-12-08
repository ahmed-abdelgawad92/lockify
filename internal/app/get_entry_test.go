package app

import (
	"context"
	"encoding/base64"
	"fmt"
	"strings"
	"testing"

	"github.com/ahmed-abdelgawad92/lockify/internal/domain/model"
)

func TestGetEntryUseCase_Execute_Success(t *testing.T) {
	env := "test"
	key := "test-key"
	value := "test-value"
	salt := "test-salt"
	passphrase := "test-passphrase"

	vaultService := &mockVaultService{
		OpenFunc: func(ctx context.Context, env string) (*model.Vault, error) {
			savedVault, _ := model.NewVault(env, "test-fingerprint", salt)
			savedVault.SetPassphrase(passphrase)
			savedVault.SetEntry(key, base64.StdEncoding.EncodeToString([]byte(value)))
			return savedVault, nil
		},
	}

	encryptionService := &mockEncryptionService{
		DecryptFunc: func(ciphertext, encodedSalt, passphrase string) ([]byte, error) {
			decodedValue, _ := base64.StdEncoding.DecodeString(ciphertext)
			return []byte(decodedValue), nil
		},
	}

	useCase := NewGetEntryUseCase(vaultService, encryptionService)

	valueRetrieved, err := useCase.Execute(context.Background(), env, key)
	if err != nil {
		t.Fatalf("Execute() returned unexpected error: %v", err)
	}

	if value != valueRetrieved {
		t.Errorf("Execute() got %s, want %s", valueRetrieved, value)
	}
}

func TestGetEntryUseCase_Execute_EntryNotFound(t *testing.T) {
	env := "test"
	key := "test-key"
	salt := "test-salt"
	passphrase := "test-passphrase"

	vaultService := &mockVaultService{
		OpenFunc: func(ctx context.Context, env string) (*model.Vault, error) {
			savedVault, _ := model.NewVault(env, "test-fingerprint", salt)
			savedVault.SetPassphrase(passphrase)
			return savedVault, nil
		},
	}

	encryptionService := &mockEncryptionService{
		DecryptFunc: func(ciphertext, encodedSalt, passphrase string) ([]byte, error) {
			decodedValue, _ := base64.StdEncoding.DecodeString(ciphertext)
			return []byte(decodedValue), nil
		},
	}

	useCase := NewGetEntryUseCase(vaultService, encryptionService)

	_, err := useCase.Execute(context.Background(), env, key)
	if err == nil {
		t.Fatalf("Execute() should return non-existence error, got nil")
	}

	if !strings.Contains(err.Error(), fmt.Sprintf("key %q not found", key)) {
		t.Errorf("Execute() error = %q, want to contain '%s'", err.Error(), fmt.Sprintf("key %q not found", key))
	}
}
