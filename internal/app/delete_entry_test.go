package app

import (
	"context"
	"encoding/base64"
	"fmt"
	"strings"
	"testing"

	"github.com/ahmed-abdelgawad92/lockify/internal/domain/model"
)

func TestDeleteEntryUseCase_Execute_Success(t *testing.T) {
env := "test"
	key := "test-key"
	value := "test-value"
	salt := "test-salt"
	passphrase := "test-passphrase"

	var savedVault *model.Vault
	vaultService := &mockVaultService{
		OpenFunc: func(ctx context.Context, env string) (*model.Vault, error) {
			savedVault, _ = model.NewVault(env, "test-fingerprint", salt)
			savedVault.SetPassphrase(passphrase)
			savedVault.SetEntry(key, base64.StdEncoding.EncodeToString([]byte(value)))
			return savedVault, nil
		},
	}

	useCase := NewDeleteEntryUseCase(vaultService)

	err := useCase.Execute(context.Background(), env, key)
	if err != nil {
		t.Fatalf("Execute() returned unexpected error: %v", err)
	}

	_, err = savedVault.GetEntry(key)
	if err == nil  {
		t.Fatalf("Execute() did not delete the key successfully")
	}
}

func TestDeleteEntryUseCase_Execute_EntryNotFound(t *testing.T) {
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

	useCase := NewDeleteEntryUseCase(vaultService)

	err := useCase.Execute(context.Background(), env, key)
	if err == nil {
		t.Fatalf("Execute() should return non-existence error, got nil")
	}

	if !strings.Contains(err.Error(), fmt.Sprintf("key %q not found", key)) {
		t.Errorf("Execute() error = %q, want to contain '%s'", err.Error(), fmt.Sprintf("key %q not found", key))
	}
}