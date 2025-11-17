package app

import (
	"context"
	"fmt"

	"github.com/apixify/lockify/internal/domain/repository"
	"github.com/apixify/lockify/internal/domain/service"
)

type AddEntryUseCase struct {
	vaultRepo         repository.VaultRepository
	encryptionService service.EncryptionService
	passphraseService service.PassphraseService
}

type AddEntryDTO struct {
	Env   string
	Key   string
	Value string
}

func NewAddEntryUseCase(
	vaultRepo repository.VaultRepository,
	encryptionService service.EncryptionService,
	passphraseService service.PassphraseService,
) AddEntryUseCase {
	return AddEntryUseCase{vaultRepo, encryptionService, passphraseService}
}

func (useCase *AddEntryUseCase) Execute(ctx context.Context, dto AddEntryDTO) error {
	passphrase, err := useCase.passphraseService.Get(ctx, dto.Env)
	if err != nil {
		return fmt.Errorf("failed to retrieve passphrase: %w", err)
	}

	vault, err := useCase.vaultRepo.Load(ctx, dto.Env)
	if err != nil {
		return fmt.Errorf("failed to open vault for environment %s: %w", dto.Env, err)
	}

	if err = useCase.passphraseService.Validate(ctx, vault, passphrase); err != nil {
		useCase.passphraseService.Clear(ctx, dto.Env)
		return fmt.Errorf("invalid credentials: %w", err)
	}

	encryptedValue, err := useCase.encryptionService.Encrypt([]byte(dto.Value), vault.Meta.Salt, passphrase)
	if err != nil {
		return fmt.Errorf("failed to encrypt value: %w", err)
	}

	err = vault.SetEntry(dto.Key, encryptedValue)
	if err != nil {
		return fmt.Errorf("failed to set entry: %w", err)
	}

	return useCase.vaultRepo.Save(ctx, vault)
}
