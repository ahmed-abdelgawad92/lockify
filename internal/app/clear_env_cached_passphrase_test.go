package app

import (
	"context"
	"errors"
	"testing"
)

func TestClearEnvCachedPassphraseUseCase_Execute_Success(t *testing.T) {
	env := "test"
	clearCalled := false
	var clearedEnv string

	passphraseService := &mockPassphraseService{
		ClearFunc: func(ctx context.Context, env string) error {
			clearCalled = true
			clearedEnv = env
			return nil
		},
	}

	useCase := NewClearEnvCachedPassphraseUseCase(passphraseService)

	err := useCase.Execute(context.Background(), env)
	if err != nil {
		t.Fatalf("Execute() returned unexpected error: %v", err)
	}

	if !clearCalled {
		t.Error("Execute() should call Clear(), but it didn't")
	}

	if clearedEnv != env {
		t.Errorf("Execute() called Clear() with env %q, want %q", clearedEnv, env)
	}
}

func TestClearEnvCachedPassphraseUseCase_Execute_Error(t *testing.T) {
	passphraseService := &mockPassphraseService{
		ClearFunc: func(ctx context.Context, env string) error {
			return errors.New("clear error")
		},
	}

	useCase := NewClearEnvCachedPassphraseUseCase(passphraseService)

	err := useCase.Execute(context.Background(), "test")
	if err == nil {
		t.Fatal("Execute() with Clear error expected error, got nil")
	}

	if err.Error() != "clear error" {
		t.Errorf("Execute() error = %q, want %q", err.Error(), "clear error")
	}
}
