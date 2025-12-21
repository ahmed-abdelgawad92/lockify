package cache

import (
	"github.com/ahmed-abdelgawad92/lockify/internal/domain/service"
	"github.com/zalando/go-keyring"
)

// OSKeyring implements Cache using the OS keyring
type OSKeyring struct {
	service string
}

// NewOSKeyring creates a new OS keyring implementation
func NewOSKeyring(s string) service.Cache {
	return &OSKeyring{service: s}
}

// Set stores a value in the keyring
func (k *OSKeyring) Set(key, value string) error {
	return keyring.Set(k.service, key, value)
}

// Get retrieves a value from the keyring
func (k *OSKeyring) Get(key string) (string, error) {
	return keyring.Get(k.service, key)
}

// Delete removes a value from the keyring
func (k *OSKeyring) Delete(key string) error {
	return keyring.Delete(k.service, key)
}

// DeleteAll removes all values for a service from the keyring
func (k *OSKeyring) DeleteAll() error {
	return keyring.DeleteAll(k.service)
}
