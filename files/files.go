package files

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
)

var (
	ErrUserHomeNotFound = errors.New("users home home not found")
)

func IsFullPath(path string) bool {
	absPath, err := filepath.Abs(path)
	if err != nil {
		return false
	}
	return absPath == filepath.Clean(path)
}

func GetUserConfigDir() (string, error) {
	if runtime.GOOS == "windows" {
		return os.Getenv("APPDATA"), nil
	}

	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("%w", ErrUserHomeNotFound)
	}
	return filepath.Join(home, ".config"), nil
}
