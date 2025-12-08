package app

import (
	"context"
	"encoding/base64"
	"slices"
	"testing"

	"github.com/ahmed-abdelgawad92/lockify/internal/domain/model"
)

func TestListEntriesUseCase_Execute_Success(t *testing.T) {
	env := "test"
	key1 := "test-key-1"
	key2 := "test-key-2"
	value := "test-value"
	salt := "test-salt"
	passphrase := "test-passphrase"

	vaultService := &mockVaultService{
		OpenFunc: func(ctx context.Context, env string) (*model.Vault, error) {
			savedVault, _ := model.NewVault(env, "test-fingerprint", salt)
			savedVault.SetPassphrase(passphrase)
			savedVault.SetEntry(key1, base64.StdEncoding.EncodeToString([]byte(value)))
			savedVault.SetEntry(key2, base64.StdEncoding.EncodeToString([]byte(value)))
			return savedVault, nil
		},
	}

	useCase := NewListEntriesUseCase(vaultService)

	allKeys, err := useCase.Execute(context.Background(), env)
	if err != nil {
		t.Fatalf("Execute() returned unexpected error: %v", err)
	}

	if len(allKeys) != 2 {
		t.Fatalf("length of keys error, want: 2, got: %v", len(allKeys))
	}

	if !slices.Contains(allKeys, key1) {
		t.Errorf("keys should contain %v", key1)
	}

	if !slices.Contains(allKeys, key2) {
		t.Errorf("keys should contain %v", key2)
	}
}
