package value

import "fmt"

// FileFormat represents the format of a file (JSON or dotenv).
type FileFormat string

const (
	// JSON represents JSON file format.
	JSON FileFormat = "json"
	// DotEnv represents dotenv file format.
	DotEnv FileFormat = "dotenv"
)

// NewFileFormat creates a new FileFormat from a string value.
func NewFileFormat(value string) (FileFormat, error) {
	format := FileFormat(value)
	if !format.IsValid() {
		return "", fmt.Errorf("invalid file format %q: must be either %q or %q", value, JSON, DotEnv)
	}
	return format, nil
}

func (fileFormat FileFormat) String() string {
	return string(fileFormat)
}

// IsJSON checks if the file format is JSON.
func (fileFormat FileFormat) IsJSON() bool {
	return fileFormat == JSON
}

// IsDotEnv checks if the file format is dotenv.
func (fileFormat FileFormat) IsDotEnv() bool {
	return fileFormat == DotEnv
}

// IsValid checks if the file format is valid.
func (fileFormat FileFormat) IsValid() bool {
	return fileFormat.IsJSON() || fileFormat.IsDotEnv()
}
