package app

import (
	"context"
	"errors"
	"strings"
	"testing"

	"github.com/ahmed-abdelgawad92/lockify/internal/domain/model"
)

func TestInitializeVaultUseCase_Execute_Success(t *testing.T) {
	env := "test"
	fingerprint := "test-fingerprint"
	salt := "test-salt"

	expectedVault, _ := model.NewVault(env, fingerprint, salt)

	vaultService := &mockVaultService{
		CreateFunc: func(ctx context.Context, env string) (*model.Vault, error) {
			if env != "test" {
				t.Errorf("Create() called with env %q, want %q", env, "test")
			}
			return expectedVault, nil
		},
	}

	useCase := NewInitializeVaultUseCase(vaultService)

	vault, err := useCase.Execute(context.Background(), env)

	if err != nil {
		t.Fatalf("Execute() returned unexpected error: %v", err)
	}

	if vault == nil {
		t.Fatal("Execute() should return a vault, but got nil")
	}

	if vault.Meta.Env != env {
		t.Errorf("Execute() returned vault with env %q, want %q", vault.Meta.Env, env)
	}

	if vault.Meta.FingerPrint != fingerprint {
		t.Errorf("Execute() returned vault with fingerprint %q, want %q", vault.Meta.FingerPrint, fingerprint)
	}

	if vault.Meta.Salt != salt {
		t.Errorf("Execute() returned vault with salt %q, want %q", vault.Meta.Salt, salt)
	}
}

func TestInitializeVaultUseCase_Execute_CreateError(t *testing.T) {
	env := "test"
	expectedError := errors.New("vault already exists")

	vaultService := &mockVaultService{
		CreateFunc: func(ctx context.Context, env string) (*model.Vault, error) {
			return nil, expectedError
		},
	}

	useCase := NewInitializeVaultUseCase(vaultService)

	vault, err := useCase.Execute(context.Background(), env)

	if err == nil {
		t.Fatalf("Execute() should return error, got nil")
	}

	if vault != nil {
		t.Fatalf("Execute() should return nil vault on error, got %v", vault)
	}

	if !strings.Contains(err.Error(), expectedError.Error()) {
		t.Errorf("Execute() error = %q, want to contain %q", err.Error(), expectedError.Error())
	}
}
