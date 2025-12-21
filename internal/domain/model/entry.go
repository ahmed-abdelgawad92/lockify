package model

// Entry represents a single encrypted entry in the vault.
type Entry struct {
	Value     string `json:"value"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}
