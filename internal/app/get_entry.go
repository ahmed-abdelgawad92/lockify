package app

import (
	"context"

	"github.com/ahmed-abdelgawad92/lockify/internal/domain/service"
)

// GetEntryUc defines the interface for retrieving entries from the vault.
type GetEntryUc interface {
	Execute(ctx context.Context, env, key string) (string, error)
}

// GetEntryUseCase implements the use case for retrieving entries from the vault.
type GetEntryUseCase struct {
	vaultService      service.VaultServiceInterface
	encryptionService service.EncryptionService
}

// NewGetEntryUseCase creates a new GetEntryUseCase instance.
func NewGetEntryUseCase(vaultService service.VaultServiceInterface, encryptionService service.EncryptionService) GetEntryUc {
	return &GetEntryUseCase{vaultService, encryptionService}
}

// Execute retrieves and decrypts an entry from the vault.
func (useCase *GetEntryUseCase) Execute(ctx context.Context, env, key string) (string, error) {
	vault, err := useCase.vaultService.Open(ctx, env)
	if err != nil {
		return "", err
	}

	entry, err := vault.GetEntry(key)
	if err != nil {
		return "", err
	}

	value, err := useCase.encryptionService.Decrypt(entry.Value, vault.Meta.Salt, vault.Passphrase())
	if err != nil {
		return "", err
	}

	return string(value), nil
}
