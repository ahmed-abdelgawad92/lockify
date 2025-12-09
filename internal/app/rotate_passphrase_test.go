package app

import (
	"context"
	"errors"
	"strings"
	"testing"

	"github.com/ahmed-abdelgawad92/lockify/internal/domain/model"
)

func TestRotatePassphraseUseCase_Execute_Success(t *testing.T) {
	env := "test"
	currentPassphrase := "old-passphrase"
	newPassphrase := "new-passphrase"
	currentSalt := "old-salt"
	newSalt := "new-salt"
	currentFingerprint := "old-fingerprint"
	newFingerprint := "new-fingerprint"

	vault, _ := model.NewVault(env, currentFingerprint, currentSalt)
	vault.SetEntry("key1", "encrypted-value-1")
	vault.SetEntry("key2", "encrypted-value-2")

	var savedVault *model.Vault
	decryptCallCount := 0
	encryptCallCount := 0

	vaultRepo := &mockVaultRepository{
		LoadFunc: func(ctx context.Context, env string) (*model.Vault, error) {
			return vault, nil
		},
		SaveFunc: func(ctx context.Context, vault *model.Vault) error {
			savedVault = vault
			return nil
		},
	}

	encryptionService := &mockEncryptionService{
		DecryptFunc: func(ciphertext, encodedSalt, passphrase string) ([]byte, error) {
			decryptCallCount++
			if encodedSalt != currentSalt {
				t.Errorf("Decrypt() called with salt %q, want %q", encodedSalt, currentSalt)
			}
			if passphrase != currentPassphrase {
				t.Errorf("Decrypt() called with passphrase %q, want %q", passphrase, currentPassphrase)
			}
			return []byte("decrypted-value"), nil
		},
		EncryptFunc: func(plaintext []byte, encodedSalt, passphrase string) (string, error) {
			encryptCallCount++
			if encodedSalt != newSalt {
				t.Errorf("Encrypt() called with salt %q, want %q", encodedSalt, newSalt)
			}
			if passphrase != newPassphrase {
				t.Errorf("Encrypt() called with passphrase %q, want %q", passphrase, newPassphrase)
			}
			return "new-encrypted-value", nil
		},
	}

	hashService := &mockHashService{
		VerifyFunc: func(hashedPassphrase, passphrase string) error {
			if hashedPassphrase != currentFingerprint {
				t.Errorf("Verify() called with fingerprint %q, want %q", hashedPassphrase, currentFingerprint)
			}
			if passphrase != currentPassphrase {
				t.Errorf("Verify() called with passphrase %q, want %q", passphrase, currentPassphrase)
			}
			return nil
		},
		GenerateSaltFunc: func(size int) (string, error) {
			if size != 16 {
				t.Errorf("GenerateSalt() called with size %d, want 16", size)
			}
			return newSalt, nil
		},
		HashFunc: func(passphrase string) (string, error) {
			if passphrase != newPassphrase {
				t.Errorf("Hash() called with passphrase %q, want %q", passphrase, newPassphrase)
			}
			return newFingerprint, nil
		},
	}

	useCase := NewRotatePassphraseUseCase(vaultRepo, encryptionService, hashService)

	err := useCase.Execute(context.Background(), env, currentPassphrase, newPassphrase)
	if err != nil {
		t.Fatalf("Execute() returned unexpected error: %v", err)
	}

	// Verify vault was saved with new salt and fingerprint
	if savedVault == nil {
		t.Fatal("Execute() should call Save() with the vault, but Save() was not called")
	}

	if savedVault.Meta.Salt != newSalt {
		t.Errorf("Execute() should update salt to %q, got %q", newSalt, savedVault.Meta.Salt)
	}

	if savedVault.Meta.FingerPrint != newFingerprint {
		t.Errorf("Execute() should update fingerprint to %q, got %q", newFingerprint, savedVault.Meta.FingerPrint)
	}

	// Verify all entries were re-encrypted
	if decryptCallCount != 2 {
		t.Errorf("Execute() should decrypt 2 entries, decrypted %d", decryptCallCount)
	}

	if encryptCallCount != 2 {
		t.Errorf("Execute() should encrypt 2 entries, encrypted %d", encryptCallCount)
	}

	// Verify entries have new encrypted values
	entry1, _ := savedVault.GetEntry("key1")
	if entry1.Value != "new-encrypted-value" {
		t.Errorf("Execute() should update entry1 value, got %q", entry1.Value)
	}

	entry2, _ := savedVault.GetEntry("key2")
	if entry2.Value != "new-encrypted-value" {
		t.Errorf("Execute() should update entry2 value, got %q", entry2.Value)
	}
}

func TestRotatePassphraseUseCase_Execute_LoadError(t *testing.T) {
	vaultRepo := &mockVaultRepository{
		LoadFunc: func(ctx context.Context, env string) (*model.Vault, error) {
			return nil, errors.New("load error")
		},
	}

	useCase := NewRotatePassphraseUseCase(vaultRepo, &mockEncryptionService{}, &mockHashService{})

	err := useCase.Execute(context.Background(), "test", "old", "new")
	if err == nil {
		t.Fatal("Execute() with load error expected error, got nil")
	}

	if !strings.Contains(err.Error(), "failed to open vault for environment") {
		t.Errorf("Execute() error = %q, want to contain 'failed to open vault for environment'", err.Error())
	}
}

func TestRotatePassphraseUseCase_Execute_VerifyError(t *testing.T) {
	vault, _ := model.NewVault("test", "fingerprint", "salt")
	vaultRepo := &mockVaultRepository{
		LoadFunc: func(ctx context.Context, env string) (*model.Vault, error) {
			return vault, nil
		},
	}

	hashService := &mockHashService{
		VerifyFunc: func(hashedPassphrase, passphrase string) error {
			return errors.New("invalid passphrase")
		},
	}

	useCase := NewRotatePassphraseUseCase(vaultRepo, &mockEncryptionService{}, hashService)

	err := useCase.Execute(context.Background(), "test", "wrong", "new")
	if err == nil {
		t.Fatal("Execute() with invalid passphrase expected error, got nil")
	}

	if !strings.Contains(err.Error(), "invalid credentials") {
		t.Errorf("Execute() error = %q, want to contain 'invalid credentials'", err.Error())
	}
}

func TestRotatePassphraseUseCase_Execute_GenerateSaltError(t *testing.T) {
	vault, _ := model.NewVault("test", "fingerprint", "salt")
	vaultRepo := &mockVaultRepository{
		LoadFunc: func(ctx context.Context, env string) (*model.Vault, error) {
			return vault, nil
		},
	}

	hashService := &mockHashService{
		VerifyFunc: func(hashedPassphrase, passphrase string) error {
			return nil
		},
		GenerateSaltFunc: func(size int) (string, error) {
			return "", errors.New("salt error")
		},
	}

	useCase := NewRotatePassphraseUseCase(vaultRepo, &mockEncryptionService{}, hashService)

	err := useCase.Execute(context.Background(), "test", "old", "new")
	if err == nil {
		t.Fatal("Execute() with salt error expected error, got nil")
	}

	if !strings.Contains(err.Error(), "failed to generate salt") {
		t.Errorf("Execute() error = %q, want to contain 'failed to generate salt'", err.Error())
	}
}

func TestRotatePassphraseUseCase_Execute_HashError(t *testing.T) {
	vault, _ := model.NewVault("test", "fingerprint", "salt")
	vaultRepo := &mockVaultRepository{
		LoadFunc: func(ctx context.Context, env string) (*model.Vault, error) {
			return vault, nil
		},
	}

	hashService := &mockHashService{
		VerifyFunc: func(hashedPassphrase, passphrase string) error {
			return nil
		},
		GenerateSaltFunc: func(size int) (string, error) {
			return "new-salt", nil
		},
		HashFunc: func(passphrase string) (string, error) {
			return "", errors.New("hash error")
		},
	}

	useCase := NewRotatePassphraseUseCase(vaultRepo, &mockEncryptionService{}, hashService)

	err := useCase.Execute(context.Background(), "test", "old", "new")
	if err == nil {
		t.Fatal("Execute() with hash error expected error, got nil")
	}

	if !strings.Contains(err.Error(), "failed to hash the fingerprint") {
		t.Errorf("Execute() error = %q, want to contain 'failed to hash the fingerprint'", err.Error())
	}
}

func TestRotatePassphraseUseCase_Execute_DecryptError(t *testing.T) {
	vault, _ := model.NewVault("test", "fingerprint", "salt")
	vault.SetEntry("key1", "encrypted-value")
	vaultRepo := &mockVaultRepository{
		LoadFunc: func(ctx context.Context, env string) (*model.Vault, error) {
			v, _ := model.NewVault(env, "fingerprint", "salt")
			v.SetEntry("key1", "encrypted-value")
			return v, nil
		},
	}

	encryptionService := &mockEncryptionService{
		DecryptFunc: func(ciphertext, encodedSalt, passphrase string) ([]byte, error) {
			return nil, errors.New("decrypt error")
		},
	}

	hashService := &mockHashService{
		VerifyFunc: func(hashedPassphrase, passphrase string) error {
			return nil
		},
		GenerateSaltFunc: func(size int) (string, error) {
			return "new-salt", nil
		},
		HashFunc: func(passphrase string) (string, error) {
			return "new-fingerprint", nil
		},
	}

	useCase := NewRotatePassphraseUseCase(vaultRepo, encryptionService, hashService)

	err := useCase.Execute(context.Background(), "test", "old", "new")
	if err == nil {
		t.Fatal("Execute() with decrypt error expected error, got nil")
	}

	if !strings.Contains(err.Error(), "failed to decrypt key") {
		t.Errorf("Execute() error = %q, want to contain 'failed to decrypt key'", err.Error())
	}
}

func TestRotatePassphraseUseCase_Execute_EncryptError(t *testing.T) {
	vault, _ := model.NewVault("test", "fingerprint", "salt")
	vault.SetEntry("key1", "encrypted-value")
	vaultRepo := &mockVaultRepository{
		LoadFunc: func(ctx context.Context, env string) (*model.Vault, error) {
			v, _ := model.NewVault(env, "fingerprint", "salt")
			v.SetEntry("key1", "encrypted-value")
			return v, nil
		},
	}

	encryptionService := &mockEncryptionService{
		DecryptFunc: func(ciphertext, encodedSalt, passphrase string) ([]byte, error) {
			return []byte("decrypted"), nil
		},
		EncryptFunc: func(plaintext []byte, encodedSalt, passphrase string) (string, error) {
			return "", errors.New("encrypt error")
		},
	}

	hashService := &mockHashService{
		VerifyFunc: func(hashedPassphrase, passphrase string) error {
			return nil
		},
		GenerateSaltFunc: func(size int) (string, error) {
			return "new-salt", nil
		},
		HashFunc: func(passphrase string) (string, error) {
			return "new-fingerprint", nil
		},
	}

	useCase := NewRotatePassphraseUseCase(vaultRepo, encryptionService, hashService)

	err := useCase.Execute(context.Background(), "test", "old", "new")
	if err == nil {
		t.Fatal("Execute() with encrypt error expected error, got nil")
	}

	if !strings.Contains(err.Error(), "failed to encrypt key") {
		t.Errorf("Execute() error = %q, want to contain 'failed to encrypt key'", err.Error())
	}
}

func TestRotatePassphraseUseCase_Execute_SaveError(t *testing.T) {
	vault, _ := model.NewVault("test", "fingerprint", "salt")
	vaultRepo := &mockVaultRepository{
		LoadFunc: func(ctx context.Context, env string) (*model.Vault, error) {
			return vault, nil
		},
		SaveFunc: func(ctx context.Context, vault *model.Vault) error {
			return errors.New("save error")
		},
	}

	hashService := &mockHashService{
		VerifyFunc: func(hashedPassphrase, passphrase string) error {
			return nil
		},
		GenerateSaltFunc: func(size int) (string, error) {
			return "new-salt", nil
		},
		HashFunc: func(passphrase string) (string, error) {
			return "new-fingerprint", nil
		},
	}

	useCase := NewRotatePassphraseUseCase(vaultRepo, &mockEncryptionService{}, hashService)

	err := useCase.Execute(context.Background(), "test", "old", "new")
	if err == nil {
		t.Fatal("Execute() with save error expected error, got nil")
	}

	if err.Error() != "save error" {
		t.Errorf("Execute() error = %q, want %q", err.Error(), "save error")
	}
}
