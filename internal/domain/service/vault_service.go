package service

import (
	"context"
	"fmt"

	"github.com/apixify/lockify/internal/domain/model"
	"github.com/apixify/lockify/internal/domain/repository"
)

type VaultService struct {
	vaultRepo         repository.VaultRepository
	passphraseService PassphraseService
}

func NewVaultService(vaultRepo repository.VaultRepository, passphraseService PassphraseService) VaultService {
	return VaultService{vaultRepo, passphraseService}
}

func (vs *VaultService) Open(ctx context.Context, env string) (*model.Vault, error) {
	passphrase, err := vs.passphraseService.Get(ctx, env)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve passphrase: %w", err)
	}

	vault, err := vs.vaultRepo.Load(ctx, env)
	if err != nil {
		return nil, fmt.Errorf("failed to open vault for environment %s: %w", env, err)
	}

	if err = vs.passphraseService.Validate(ctx, vault, passphrase); err != nil {
		vs.passphraseService.Clear(ctx, env)
		return nil, fmt.Errorf("invalid credentials: %w", err)
	}

	vault.SetPassphrase(passphrase)

	return vault, nil
}

func (vs *VaultService) Save(ctx context.Context, vault *model.Vault) error {
	return vs.vaultRepo.Save(ctx, vault)
}
