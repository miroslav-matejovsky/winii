package winservices

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// findServiceExecutable extracts the executable path from a given Windows service executable string,
// validates that it is not empty or the current directory (which are considered invalid),
// and verifies the file's existence on disk. If the path has no file extension, it appends ".exe"
// and checks for existence again.
//
// Parameters:
//   - winServiceExecutable: A string representing the Windows service executable path, potentially with arguments.
//
// Returns:
//   - A string representing the verified executable path.
//   - An error if the path is invalid, empty, the current directory, or the file does not exist.
func findServiceExecutable(winServiceExecutable string) (string, error) {
	path := extractExecutablePath(winServiceExecutable)
	if path == "" || path == "." {
		return "", fmt.Errorf("invalid service executable path: %q", winServiceExecutable)
	}
	if _, err := os.Stat(path); err == nil {
		return path, nil
	}
	ext := filepath.Ext(path)
	if ext == "" {
		path += ".exe"
		if _, err := os.Stat(path); err == nil {
			return path, nil
		}
	}
	return "", fmt.Errorf("service executable not found: %q", winServiceExecutable)
}

// extractExecutablePath extracts the executable path from a given Windows service executable string.
// It handles quoted paths (e.g., "C:\Path\app.exe") and unquoted paths with arguments,
// stopping at the first space followed by an argument marker ('/' or '-').
//
// Parameters:
//   - winServiceExecutable: A string representing the Windows service executable path, potentially with arguments.
//
// Returns:
//   - A cleaned string representing the extracted executable path. Returns "." for an empty input.
func extractExecutablePath(winServiceExecutable string) string {
	// Handle quoted path
	if strings.HasPrefix(winServiceExecutable, `"`) {
		if endQuote := strings.Index(winServiceExecutable[1:], `"`); endQuote != -1 {
			return filepath.Clean(winServiceExecutable[1 : endQuote+1])
		}
	}

	// If not quoted, find first argument marker (space followed by slash or hyphen)
	for i := 0; i < len(winServiceExecutable); i++ {
		if winServiceExecutable[i] == ' ' {
			// Check if next character is an argument marker
			if i+1 < len(winServiceExecutable) && (winServiceExecutable[i+1] == '/' || winServiceExecutable[i+1] == '-') {
				return filepath.Clean(winServiceExecutable[:i])
			}
		}
	}

	// If no argument markers found, return the whole string
	return filepath.Clean(winServiceExecutable)
}
