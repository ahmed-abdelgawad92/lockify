package app

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/apixify/lockify/internal/domain"
	"github.com/apixify/lockify/internal/domain/model/value"
	"github.com/apixify/lockify/internal/domain/repository"
	"github.com/apixify/lockify/internal/domain/service"
)

type ExportEnvUseCase struct {
	vaultRepo         repository.VaultRepository
	passphraseService service.PassphraseService
	encryptionService service.EncryptionService
	logger            domain.Logger
}

func NewExportEnvUseCase(
	vaultRepo repository.VaultRepository,
	passphraseService service.PassphraseService,
	encryptionService service.EncryptionService,
	logger domain.Logger,
) ExportEnvUseCase {
	return ExportEnvUseCase{vaultRepo, passphraseService, encryptionService, logger}
}

func (useCase *ExportEnvUseCase) Execute(ctx context.Context, env string, exportFormat value.FileFormat) error {
	if !exportFormat.IsValid() {
		return fmt.Errorf("format must be either %s or %s. %s is given", value.Json, value.DotEnv, exportFormat)
	}
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

	if exportFormat.IsDotEnv() {
		for k, v := range vault.Entries {
			decryptedVal, _ := useCase.encryptionService.Decrypt(v.Value, vault.Meta.Salt, passphrase)
			useCase.logger.Output("%s=%s\n", k, decryptedVal)
		}
	} else {
		mappedEntries := make(map[string]string)
		for k, v := range vault.Entries {
			decryptedVal, _ := useCase.encryptionService.Decrypt(v.Value, vault.Meta.Salt, passphrase)
			mappedEntries[k] = string(decryptedVal)
		}

		data, _ := json.MarshalIndent(mappedEntries, "", "  ")
		useCase.logger.Output(string(data))
	}

	return nil
}
