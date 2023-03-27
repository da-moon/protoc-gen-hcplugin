package utils

import (
	"os"
	"path/filepath"
	"strings"
)

// FromRelativePath -
func FromRelativePath(arg string) (string, error) {
	abs, err := filepath.Abs(arg)
	if err != nil {
		return "", err
	}
	return strings.TrimPrefix(abs, os.Getenv("GOPATH")+"/src/"), nil

}
