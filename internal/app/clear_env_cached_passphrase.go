package app

import (
	"context"

	"github.com/ahmed-abdelgawad92/lockify/internal/domain/service"
)

// ClearEnvCachedPassphraseUseCase implements the use case for clearing cached passphrase for a specific environment.
type ClearEnvCachedPassphraseUseCase struct {
	passphraseService service.PassphraseService
}

// NewClearEnvCachedPassphraseUseCase creates a new ClearEnvCachedPassphraseUseCase instance.
func NewClearEnvCachedPassphraseUseCase(
	passphraseService service.PassphraseService,
) ClearEnvCachedPassphraseUseCase {
	return ClearEnvCachedPassphraseUseCase{passphraseService}
}

// Execute clears the cached passphrase for the specified environment.
func (useCase *ClearEnvCachedPassphraseUseCase) Execute(ctx context.Context, env string) error {
	return useCase.passphraseService.Clear(ctx, env)
}
