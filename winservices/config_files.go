package winservices

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

var configFileExtensions = []string{
	".json",
	".xml",
	".ini",
	".config",
}

type ServiceConfigFile struct {
	Path     string
	Contents string

	CreationTime   time.Time
	LastAccessTime time.Time
	LastWriteTime  time.Time
}

// collectServiceConfigFiles walks the provided directory, finds files with known
// configuration extensions, reads their contents and returns a slice of ServiceConfigFile.
// The dir parameter must point to an existing directory.
// Returns an error if dir is empty, does not exist, is not a directory, or if reading files fails.
func collectServiceConfigFiles(dir string) ([]ServiceConfigFile, error) {
	if dir == "" {
		return nil, fmt.Errorf("dir is empty")
	}

	// Validate that dir exists and is a directory
	info, err := os.Stat(dir)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("dir does not exist: %q", dir)
		}
		return nil, fmt.Errorf("failed to stat dir %q: %w", dir, err)
	}
	if !info.IsDir() {
		return nil, fmt.Errorf("path is not a directory: %q", dir)
	}

	var configFiles []ServiceConfigFile

	err = filepath.Walk(dir, func(path string, info os.FileInfo, walkErr error) error {
		if walkErr != nil {
			return fmt.Errorf("error walking the path %q: %w", path, walkErr)
		}
		if info.IsDir() {
			return nil
		}
		ext := strings.ToLower(filepath.Ext(info.Name()))
		for _, cfgExt := range configFileExtensions {
			if ext == cfgExt {
				data, err := os.ReadFile(path)
				if err != nil {
					return fmt.Errorf("failed to read config file %q: %w", path, err)
				}
				times, err := getFileTime(path)
				if err != nil {
					return fmt.Errorf("failed to get file times for %q: %w", path, err)
				}
				configFiles = append(configFiles, ServiceConfigFile{
					Path:           path,
					Contents:       string(data),
					CreationTime:   times.CreationTime,
					LastAccessTime: times.LastAccessTime,
					LastWriteTime:  times.LastWriteTime,
				})
				break
			}
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return configFiles, nil
}
