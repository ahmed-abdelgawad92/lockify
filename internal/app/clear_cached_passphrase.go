package app

import (
	"context"

	"github.com/ahmed-abdelgawad92/lockify/internal/domain/service"
)

// ClearCachedPassphraseUc defines the interface for clearing all cached passphrases.
type ClearCachedPassphraseUc interface {
	Execute(context.Context) error
}

// ClearCachedPassphraseUseCase implements the use case for clearing all cached passphrases.
type ClearCachedPassphraseUseCase struct {
	passphraseService service.PassphraseService
}

// NewClearCachedPassphraseUseCase creates a new ClearCachedPassphraseUseCase instance.
func NewClearCachedPassphraseUseCase(
	passphraseService service.PassphraseService,
) ClearCachedPassphraseUc {
	return &ClearCachedPassphraseUseCase{passphraseService}
}

// Execute clears all cached passphrases from the system keyring.
func (useCase *ClearCachedPassphraseUseCase) Execute(ctx context.Context) error {
	return useCase.passphraseService.ClearAll(ctx)
}
