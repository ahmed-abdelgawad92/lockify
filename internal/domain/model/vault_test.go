package model

import (
	"slices"
	"strings"
	"testing"
	"time"
)

func createTestVault(t *testing.T) *Vault {
	t.Helper()
	vault, err := NewVault("test", "test", "test")
	if err != nil {
		t.Fatalf("failed to create test vault: %v", err)
	}
	if vault == nil {
		t.Fatal("NewVault() should return a vault, got nil")
	}

	return vault
}

func TestNewVault(t *testing.T) {
	vault := createTestVault(t)

	if vault.Meta.Env != "test" {
		t.Errorf("expected env %q, got %q", "test", vault.Meta.Env)
	}
	if vault.Meta.FingerPrint != "test" {
		t.Errorf("expected fingerprint %q, got %q", "test", vault.Meta.FingerPrint)
	}
	if vault.Meta.Salt != "test" {
		t.Errorf("expected salt %q, got %q", "test", vault.Meta.Salt)
	}
	if len(vault.Entries) != 0 {
		t.Errorf("expected 0 entries, got %d", len(vault.Entries))
	}
}

func TestSetEntry(t *testing.T) {
	vault := createTestVault(t)

	vault.SetEntry("test_key", "test_value")
	if len(vault.Entries) != 1 {
		t.Errorf("expected 1 entry, got %d", len(vault.Entries))
	}
	if vault.Entries["test_key"].Value != "test_value" {
		t.Errorf("expected value %q, got %q", "test_value", vault.Entries["test_key"].Value)
	}
	if vault.Entries["test_key"].CreatedAt == "" {
		t.Errorf("expected created at, got empty")
	}
}

func TestGetEntry(t *testing.T) {
	vault := createTestVault(t)

	vault.SetEntry("test_key", "test_value")
	if len(vault.Entries) != 1 {
		t.Errorf("expected 1 entry, got %d", len(vault.Entries))
	}

	entry, err := vault.GetEntry("test_key")
	if err != nil {
		t.Fatalf("failed to get entry: %v", err)
	}
	if entry.Value != "test_value" {
		t.Errorf("expected value %q, got %q", "test_value", entry.Value)
	}
	if entry.CreatedAt == "" {
		t.Errorf("expected created at, got empty")
	}
	if entry.UpdatedAt == "" {
		t.Errorf("expected updated at, got empty")
	}
}

func TestSetEntryUpdateKey(t *testing.T) {
	vault := createTestVault(t)

	vault.SetEntry("test_key", "test_value")
	if len(vault.Entries) != 1 {
		t.Errorf("expected 1 entry, got %d", len(vault.Entries))
	}
	// Manually set timestamps to past time for testing
	entry := vault.Entries["test_key"]
	pastTime := time.Now().UTC().Add(-24 * time.Hour).Format(time.RFC3339)
	entry.CreatedAt = pastTime
	entry.UpdatedAt = pastTime
	vault.Entries["test_key"] = entry

	vault.SetEntry("test_key", "test_value2")
	entry, err := vault.GetEntry("test_key")
	if err != nil {
		t.Fatalf("failed to get entry: %v", err)
	}
	if entry.Value != "test_value2" {
		t.Errorf("expected value %q, got %q", "test_value2", entry.Value)
	}
	if entry.CreatedAt != pastTime {
		t.Errorf("expected created at to be %q, got %q", pastTime, entry.CreatedAt)
	}
	if entry.UpdatedAt == pastTime {
		t.Errorf("expected updated at not to be %q, got %q", pastTime, entry.UpdatedAt)
	}
}

func TestGetNonExistentEntry(t *testing.T) {
	vault := createTestVault(t)

	_, err := vault.GetEntry("non_existent_key")
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
	if err.Error() != "key \"non_existent_key\" not found" {
		t.Errorf("expected error %q, got %q", "key \"non_existent_key\" not found", err.Error())
	}
}

func TestGetEntryWithEmptyKey(t *testing.T) {
	vault := createTestVault(t)

	_, err := vault.GetEntry("")
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
	if err.Error() != "key cannot be empty" {
		t.Errorf("expected error %q, got %q", "key cannot be empty", err.Error())
	}
}

func TestDeleteEntry(t *testing.T) {
	vault := createTestVault(t)

	vault.SetEntry("test_key", "test_value")
	if len(vault.Entries) != 1 {
		t.Errorf("expected 1 entry, got %d", len(vault.Entries))
	}

	vault.DeleteEntry("test_key")
	if len(vault.Entries) != 0 {
		t.Errorf("expected 0 entries, got %d", len(vault.Entries))
	}

	_, err := vault.GetEntry("test_key")
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
	if err.Error() != "key \"test_key\" not found" {
		t.Errorf("expected error %q, got %q", "key \"test_key\" not found", err.Error())
	}
}

func TestDeleteNonExistentEntry(t *testing.T) {
	vault := createTestVault(t)

	err := vault.DeleteEntry("non_existent_key")
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
	if err.Error() != "key \"non_existent_key\" not found" {
		t.Errorf("expected error %q, got %q", "key \"non_existent_key\" not found", err.Error())
	}
}

func TestListKeys(t *testing.T) {
	vault := createTestVault(t)

	vault.SetEntry("test_key", "test_value")
	if len(vault.Entries) != 1 {
		t.Errorf("expected 1 entry, got %d", len(vault.Entries))
	}

	keys := vault.ListKeys()
	if len(keys) != 1 {
		t.Errorf("expected 1 key, got %d", len(keys))
	}
	if keys[0] != "test_key" {
		t.Errorf("expected key %q, got %q", "test_key", keys[0])
	}
}

func TestListKeysEmpty(t *testing.T) {
	vault := createTestVault(t)

	keys := vault.ListKeys()
	if len(keys) != 0 {
		t.Errorf("expected 0 keys, got %d", len(keys))
	}
}

func TestListKeysMultiple(t *testing.T) {
	vault := createTestVault(t)

	vault.SetEntry("test_key", "test_value")
	vault.SetEntry("test_key2", "test_value2")

	keys := vault.ListKeys()
	if len(keys) != 2 {
		t.Errorf("expected 2 keys, got %d", len(keys))
	}
	if !slices.Contains(keys, "test_key") {
		t.Errorf("expected keys to contain %q", "test_key")
	}
	if !slices.Contains(keys, "test_key2") {
		t.Errorf("expected keys to contain %q", "test_key2")
	}
}

func TestErrorSetEntryWithEmptyKey(t *testing.T) {
	vault := createTestVault(t)

	err := vault.SetEntry("", "test_value")
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
	if err.Error() != "key cannot be empty" {
		t.Errorf("expected error %q, got %q", "key cannot be empty", err.Error())
	}
}

func TestSetEntryWithEmptyValue(t *testing.T) {
	vault := createTestVault(t)

	err := vault.SetEntry("test_key", "")
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
	if err.Error() != "encrypted value cannot be empty" {
		t.Errorf("expected error %q, got %q", "encrypted value cannot be empty", err.Error())
	}
}

func TestSetPassphrase(t *testing.T) {
	vault := createTestVault(t)

	vault.SetPassphrase("test_passphrase")
	if vault.Passphrase() != "test_passphrase" {
		t.Errorf("expected passphrase %q, got %q", "test_passphrase", vault.Passphrase())
	}
}

func TestSetPath(t *testing.T) {
	vault := createTestVault(t)

	vault.SetPath("test_path")
	if vault.Path() != "test_path" {
		t.Errorf("expected path %q, got %q", "test_path", vault.Path())
	}
}

func TestNewVault_ValidationErrors(t *testing.T) {
	tests := []struct {
		name        string
		env         string
		fingerprint string
		salt        string
		wantErr     string
	}{
		{
			name:        "empty env",
			env:         "",
			fingerprint: "test",
			salt:        "test",
			wantErr:     "environment cannot be empty",
		},
		{
			name:        "empty fingerprint",
			env:         "test",
			fingerprint: "",
			salt:        "test",
			wantErr:     "fingerprint cannot be empty",
		},
		{
			name:        "empty salt",
			env:         "test",
			fingerprint: "test",
			salt:        "",
			wantErr:     "salt cannot be empty",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := NewVault(tt.env, tt.fingerprint, tt.salt)
			if err == nil {
				t.Fatal("expected error, got nil")
			}
			if !strings.Contains(err.Error(), tt.wantErr) {
				t.Errorf("expected error to contain %q, got %q", tt.wantErr, err.Error())
			}
		})
	}
}
