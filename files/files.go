package files

import (
	"errors"
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

var (
	ErrUserHomeNotFound = errors.New("Users home directory not found")
)

func IsFullPath(path string) bool {
	absPath, err := filepath.Abs(path)
	if err != nil {
		return false
	}
	return absPath == filepath.Clean(path)
}

func PathExists(path string) bool {
	if _, err := os.Stat(path); err != nil {
		if os.IsNotExist(err) {
			return false
		} else {
			return true
		}
	}
	return false
}

func GetUserConfigDir() (string, bool) {
	var configpath string
	if runtime.GOOS == "windows" {
		configpath = filepath.Join(os.Getenv("APPDATA"), "wallmancer")
		if PathExists(configpath) {
			return configpath, true
		} else {
			return configpath, false
		}
	}

	home, err := os.UserHomeDir()
	if err != nil {
		configpath = filepath.Join("/tmp")
	}
	configpath = filepath.Join(home, ".config", "wallmancer")

	if PathExists(configpath) {
		return configpath, true
	} else {
		return configpath, false
	}
}

func DefaultConfigFilepath() (string, bool) {
	cfg_dir, _ := GetUserConfigDir()
	cfg_path := filepath.Join(cfg_dir, "config.yml")
	exists := PathExists(cfg_path)
	return cfg_path, exists
}

func GetCacheDir() (string, error) {
	var cacheDir string

	switch runtime.GOOS {
	case "windows":
		cacheDir = filepath.Join(os.Getenv("LOCALAPPDATA"), "wallmancer", "cache")
	case "darwin":
		home, err := os.UserHomeDir()
		if err != nil {
			return "", fmt.Errorf("%w: %v", ErrUserHomeNotFound, err)
		}
		cacheDir = filepath.Join(home, "Library", "Caches", "wallmancer")
	default: // Linux/Unix
		home, err := os.UserHomeDir()
		if err != nil {
			return "", fmt.Errorf("%w: %v", ErrUserHomeNotFound, err)
		}
		cacheDir = filepath.Join(home, ".cache", "wallmancer")
	}

	if err := os.MkdirAll(cacheDir, 0700); err != nil {
		return "", fmt.Errorf("failed to create cache directory: %w", err)
	}

	return cacheDir, nil
}

func WriteStringToCache(filename string, str string) error {
	// Get or create cache directory
	cacheDir, err := GetCacheDir()
	if err != nil {
		return err
	}

	// Ensure the full path exists
	fullPath := filepath.Join(cacheDir, filename)
	dirPath := filepath.Dir(fullPath)

	// Create parent directories if they don't exist
	if err := os.MkdirAll(dirPath, 0755); err != nil {
		return fmt.Errorf("failed to create cache directories: %w", err)
	}

	// Write file with user-only permissions (0600)
	if err := os.WriteFile(fullPath, []byte(str), 0600); err != nil {
		return fmt.Errorf("failed to write cache file: %w", err)
	}

	return nil
}

func WriteSliceToCache(filename string, lines []string) error {
	cacheDir, err := GetCacheDir()
	if err != nil {
		return err
	}

	fullPath := filepath.Join(cacheDir, filename)
	if err := os.MkdirAll(filepath.Dir(fullPath), 0755); err != nil {
		return fmt.Errorf("failed to create directories: %w", err)
	}

	// Create a buffer to build content
	var content strings.Builder
	for _, line := range lines {
		content.WriteString(line)
		content.WriteString("\n")
	}

	// Write all content at once with proper permissions
	if err := os.WriteFile(fullPath, []byte(content.String()), 0600); err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}

	return nil
}

func IsFileFresh(filepath string, expirySeconds int) bool {
	fileInfo, err := os.Stat(filepath)
	if err != nil {
		if os.IsNotExist(err) {
			return false
		}
		return false
	}

	modifiedTime := fileInfo.ModTime()
	expiryDuration := time.Duration(expirySeconds) * time.Second
	expiryTime := modifiedTime.Add(expiryDuration)

	return time.Now().Before(expiryTime)
}

func ReadFromCache(path string) (string, error) {
	cache, err := GetCacheDir()
	if err != nil {
		return "", fmt.Errorf("%w", err)
	}
	content, err := os.ReadFile(filepath.Join(cache, path))
	if err != nil {
		return "", fmt.Errorf("failed to read file: %w", err)
	}
	return string(content), nil
}

func ReadLine(filename string) (string, error) {
	content, err := os.ReadFile(filename)
	if err != nil {
		return "", fmt.Errorf("failed to read file: %w", err)
	}
	return string(content), nil
}

func GetRandomLine(filename string) (string, error) {
	// Read entire file
	content, err := os.ReadFile(filename)
	if err != nil {
		return "", fmt.Errorf("failed to read file: %w", err)
	}

	// Split into lines
	lines := strings.Split(string(content), "\n")

	// Filter out empty lines
	var nonEmptyLines []string
	for _, line := range lines {
		if strings.TrimSpace(line) != "" {
			nonEmptyLines = append(nonEmptyLines, line)
		}
	}

	if len(nonEmptyLines) == 0 {
		return "", fmt.Errorf("file is empty or contains only blank lines")
	}

	return nonEmptyLines[rand.Intn(len(nonEmptyLines))], nil
}
