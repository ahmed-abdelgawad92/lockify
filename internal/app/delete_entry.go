package app

import (
	"context"
	"fmt"

	"github.com/apixify/lockify/internal/domain/repository"
	"github.com/apixify/lockify/internal/domain/service"
)

type DeleteEntryUseCase struct {
	vaultRepo         repository.VaultRepository
	passphraseService service.PassphraseService
}

func NewDeleteEntryUseCase(vaultRepo repository.VaultRepository, passphraseService service.PassphraseService) DeleteEntryUseCase {
	return DeleteEntryUseCase{vaultRepo, passphraseService}
}

func (useCase *DeleteEntryUseCase) Execute(ctx context.Context, env, key string) error {
	passphrase, err := useCase.passphraseService.Get(ctx, env)
	if err != nil {
		return fmt.Errorf("failed to retrieve passphrase: %w", err)
	}

	vault, err := useCase.vaultRepo.Load(ctx, env)
	if err != nil {
		return fmt.Errorf("failed to open vault for environment %s: %w", env, err)
	}

	if err = useCase.passphraseService.Validate(ctx, vault, passphrase); err != nil {
		useCase.passphraseService.Clear(ctx, env)
		return fmt.Errorf("invalid credentials: %w", err)
	}

	if err = vault.DeleteEntry(key); err != nil {
		return fmt.Errorf("failed to delete key %s: %w", key, err)
	}

	return useCase.vaultRepo.Save(ctx, vault)
}
