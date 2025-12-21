package app

import (
	"context"

	"github.com/ahmed-abdelgawad92/lockify/internal/domain/service"
)

// ListEntriesUc defines the interface for listing entries in the vault.
type ListEntriesUc interface {
	Execute(ctx context.Context, env string) ([]string, error)
}

// ListEntriesUseCase implements the use case for listing entries in the vault.
type ListEntriesUseCase struct {
	vaultService service.VaultServiceInterface
}

// NewListEntriesUseCase creates a new ListEntriesUseCase instance.
func NewListEntriesUseCase(vaultService service.VaultServiceInterface) ListEntriesUc {
	return &ListEntriesUseCase{vaultService}
}

// Execute lists all entry keys in the vault for the specified environment.
func (useCase *ListEntriesUseCase) Execute(ctx context.Context, env string) ([]string, error) {
	vault, err := useCase.vaultService.Open(ctx, env)
	if err != nil {
		return nil, err
	}

	keys := make([]string, 0, len(vault.Entries))
	for k := range vault.Entries {
		keys = append(keys, k)
	}

	return keys, nil
}
