package app

import (
	"context"
	"errors"
	"testing"
)

func TestClearCachedPassphraseUseCase_Execute_Success(t *testing.T) {
	clearAllCalled := false
	passphraseService := &mockPassphraseService{
		ClearAllFunc: func(ctx context.Context) error {
			clearAllCalled = true
			return nil
		},
	}

	useCase := NewClearCachedPassphraseUseCase(passphraseService)

	err := useCase.Execute(context.Background())
	if err != nil {
		t.Fatalf("Execute() returned unexpected error: %v", err)
	}

	if !clearAllCalled {
		t.Error("Execute() should call ClearAll(), but it didn't")
	}
}

func TestClearCachedPassphraseUseCase_Execute_Error(t *testing.T) {
	passphraseService := &mockPassphraseService{
		ClearAllFunc: func(ctx context.Context) error {
			return errors.New("clear all error")
		},
	}

	useCase := NewClearCachedPassphraseUseCase(passphraseService)

	err := useCase.Execute(context.Background())
	if err == nil {
		t.Fatal("Execute() with ClearAll error expected error, got nil")
	}

	if err.Error() != "clear all error" {
		t.Errorf("Execute() error = %q, want %q", err.Error(), "clear all error")
	}
}
