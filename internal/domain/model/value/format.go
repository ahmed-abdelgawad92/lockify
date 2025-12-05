package value

import "fmt"

type FileFormat string

const (
	Json   FileFormat = "json"
	DotEnv FileFormat = "dotenv"
)

func NewFileFormat(value string) (FileFormat, error) {
	format := FileFormat(value)
	if !format.IsValid() {
		return "", fmt.Errorf("invalid file format %q: must be either %q or %q", value, Json, DotEnv)
	}
	return format, nil
}

func (fileFormat FileFormat) String() string {
	return string(fileFormat)
}

func (fileFormat FileFormat) IsJson() bool {
	return fileFormat == Json
}

func (fileFormat FileFormat) IsDotEnv() bool {
	return fileFormat == DotEnv
}

func (fileFormat FileFormat) IsValid() bool {
	return fileFormat.IsJson() || fileFormat.IsDotEnv()
}
