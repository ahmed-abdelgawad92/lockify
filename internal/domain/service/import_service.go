package service

import "io"

// ImportService defines the interface for importing entries from various file formats.
type ImportService interface {
	FromJSON(r io.Reader) (map[string]string, error)
	FromDotEnv(r io.Reader) (map[string]string, error)
}
