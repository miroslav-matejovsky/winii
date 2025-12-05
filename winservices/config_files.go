package winservices

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	fi "github.com/miroslav-matejovsky/winii/fileinfo"
)

var configFileExtensions = []string{
	".json",
	".xml",
	".ini",
	".config",
}

var skipDirs = []string{
	"c:\\windows",
}

type ServiceConfigFile struct {
	Path           string    // Path is the full path to the configuration file.
	Contents       string    // Contents holds the text content of the configuration file.
	CreationTime   time.Time // CreationTime is when the file was created.
	LastAccessTime time.Time // LastAccessTime is when the file was last accessed.
	LastWriteTime  time.Time // LastWriteTime is when the file was last modified.
}

func (s ServiceConfigFile) String() string {
	return s.Path
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
	// We are on Windows, so perform case-insensitive normalization
	normalizedDir := strings.ToLower(dir)
	// Check if the directory is a system directory to skip
	for _, skip := range skipDirs {
		if strings.HasPrefix(normalizedDir, skip) {
			return []ServiceConfigFile{}, nil
		}
	}

	var configFiles []ServiceConfigFile
	var firstError error

	err = filepath.Walk(dir, func(path string, info os.FileInfo, walkErr error) error {
		if walkErr != nil {
			if firstError == nil {
				firstError = fmt.Errorf("error walking the path %q: %w", path, walkErr)
			}
			return nil
		}
		if info.IsDir() {
			return nil
		}
		ext := strings.ToLower(filepath.Ext(info.Name()))
		for _, cfgExt := range configFileExtensions {
			if ext == cfgExt {
				data, err := os.ReadFile(path)
				if err != nil {
					if firstError == nil {
						firstError = fmt.Errorf("failed to read config file %q: %w", path, err)
					}
					return nil
				}
				timestamps, err := fi.GetFileTimestamps(path)
				if err != nil {
					if firstError == nil {
						firstError = fmt.Errorf("failed to get file times for %q: %w", path, err)
					}
					return nil
				}
				configFiles = append(configFiles, ServiceConfigFile{
					Path:           path,
					Contents:       string(data),
					CreationTime:   timestamps.CreationTime,
					LastAccessTime: timestamps.LastAccessTime,
					LastWriteTime:  timestamps.LastWriteTime,
				})
				break
			}
		}
		return nil
	})
	if err != nil {
		return configFiles, fmt.Errorf("failed to walk directory %q: %w", dir, err)
	}
	return configFiles, firstError
}
