package di

import (
	"github.com/ahmed-abdelgawad92/lockify/internal/app"
	"github.com/ahmed-abdelgawad92/lockify/internal/config"
	"github.com/ahmed-abdelgawad92/lockify/internal/domain"
	"github.com/ahmed-abdelgawad92/lockify/internal/domain/repository"
	"github.com/ahmed-abdelgawad92/lockify/internal/domain/service"
	"github.com/ahmed-abdelgawad92/lockify/internal/domain/storage"
	"github.com/ahmed-abdelgawad92/lockify/internal/infrastructure/cache"
	"github.com/ahmed-abdelgawad92/lockify/internal/infrastructure/crypto"
	"github.com/ahmed-abdelgawad92/lockify/internal/infrastructure/fs"
	"github.com/ahmed-abdelgawad92/lockify/internal/infrastructure/logger"
)

var (
	vaultConfig      = config.DefaultVaultConfig()
	encryptionConfig = config.DefaultEncryptionConfig()
	log              = logger.New()
)

func getHashService() service.HashService {
	return crypto.NewBcryptHashService()
}

func getCacheService() service.Cache {
	return cache.NewOSKeyring("lockify")
}

func getPassphraseService() service.PassphraseService {
	return crypto.NewPassphraseService(getCacheService(), getHashService(), vaultConfig.PassphraseEnv)
}

func getEncryptionService() service.EncryptionService {
	return crypto.NewAESEncryptionService(encryptionConfig)
}

func getFileSystemStorage() storage.FileSystem {
	return fs.NewOSFileSystem()
}

func getVaultRepository() repository.VaultRepository {
	return fs.NewFileVaultRepository(getFileSystemStorage(), vaultConfig)
}

func getVaultService() service.VaultService {
	return service.NewVaultService(getVaultRepository(), getPassphraseService(), getHashService())
}

func getImportService() service.ImportService {
	return fs.NewFsImportService()
}

func GetLogger() domain.Logger {
	return log
}

func BuildAddEntry() app.AddEntryUseCase {
	return app.NewAddEntryUseCase(getVaultService(), getEncryptionService())
}

func BuildClearCachedPassphrase() app.ClearCachedPassphraseUseCase {
	return app.NewClearCachedPassphraseUseCase(getPassphraseService())
}

func BuildClearEnvCachedPassphrase() app.ClearEnvCachedPassphraseUseCase {
	return app.NewClearEnvCachedPassphraseUseCase(getPassphraseService())
}

func BuildDeleteEntry() app.DeleteEntryUseCase {
	return app.NewDeleteEntryUseCase(getVaultService())
}

func BuildExportEnv() app.ExportEnvUseCase {
	return app.NewExportEnvUseCase(getVaultService(), getEncryptionService(), GetLogger())
}

func BuildGetEntry() app.GetEntryUseCase {
	return app.NewGetEntryUseCase(getVaultService(), getEncryptionService())
}

func BuildInitializeVault() app.InitializeVaultUseCase {
	return app.NewInitializeVaultUseCase(getVaultService())
}

func BuildListEntries() app.ListEntriesUseCase {
	return app.NewListEntriesUseCase(getVaultService())
}

func BuildRotatePassphrase() app.RotatePassphraseUseCase {
	return app.NewRotatePassphraseUseCase(getVaultRepository(), getEncryptionService(), getHashService())
}

func BuildImportEnv() app.ImportEnvUseCase {
	return app.NewImportEnvUseCase(getVaultService(), getImportService(), getEncryptionService(), GetLogger())
}
