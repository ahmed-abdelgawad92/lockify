package model

// Meta contains metadata about the vault including environment, salt, and fingerprint.
type Meta struct {
	Env         string `json:"env"`
	Salt        string `json:"salt"`
	FingerPrint string `json:"fingerprint"`
}
