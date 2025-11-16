package service

type Cache interface {
	// Set stores a value in cache
	Set(service, key, value string) error
	// Get retrieves a value from cache
	Get(service, key string) (string, error)
	// Delete removes a value from cache
	Delete(service, key string) error
	// DeleteAll removes all values for a service from cache
	DeleteAll(service string) error
}
