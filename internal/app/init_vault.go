package app

import (
	"context"

	"github.com/ahmed-abdelgawad92/lockify/internal/domain/model"
	"github.com/ahmed-abdelgawad92/lockify/internal/domain/service"
)

type InitializeVaultUseCase struct {
	vaultService service.VaultService
}

func NewInitializeVaultUseCase(vaultService service.VaultService) InitializeVaultUseCase {
	return InitializeVaultUseCase{vaultService}
}

func (useCase *InitializeVaultUseCase) Execute(ctx context.Context, env string) (*model.Vault, error) {
	return useCase.vaultService.Create(ctx, env)
}
