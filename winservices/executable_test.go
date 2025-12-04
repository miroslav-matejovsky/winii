package winservices

import (
	"os"
	"testing"
)

func TestFindServiceExecutable(t *testing.T) {
	// Get the path of the current test executable for testing existing file
	exePath, err := os.Executable()
	if err != nil {
		t.Fatalf("failed to get executable path: %v", err)
	}

	tests := []struct {
		name        string
		input       string
		expected    string
		expectError bool
	}{
		{
			name:        "Existing executable path",
			input:       exePath,
			expected:    exePath,
			expectError: false,
		},
		{
			name:        "Existing executable path with arguments",
			input:       exePath + " /service",
			expected:    exePath,
			expectError: false,
		},
		{
			name:        "Quoted existing executable path",
			input:       `"` + exePath + `" /service`,
			expected:    exePath,
			expectError: false,
		},
		{
			name:        "Non-existent path",
			input:       `C:\NonExistent\Path\app.exe`,
			expected:    "",
			expectError: true,
		},
		{
			name:        "Empty string",
			input:       "",
			expected:    "",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := findServiceExecutable(tt.input)
			if tt.expectError {
				if err == nil {
					t.Errorf("findServiceExecutable(%q) expected error, got nil", tt.input)
				}
			} else {
				if err != nil {
					t.Errorf("findServiceExecutable(%q) unexpected error: %v", tt.input, err)
				}
				if got != tt.expected {
					t.Errorf("findServiceExecutable(%q) = %q, want %q", tt.input, got, tt.expected)
				}
			}
		})
	}
}

func TestFindServiceExecutablePath(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "Windows path with exe and argument",
			input:    `C:\Company\Product\Subfolder\App.exe /service`,
			expected: `C:\Company\Product\Subfolder\App.exe`,
		},
		{
			name:     "Quoted Windows path with exe and argument",
			input:    `"C:\Company\Product\Some Server\App.exe" /service`,
			expected: `C:\Company\Product\Some Server\App.exe`,
		},
		{
			name:     "Windows path without exe and argument",
			input:    `C:\Company\Product\App Folder\SomeApp /service`,
			expected: `C:\Company\Product\App Folder\SomeApp`,
		},
		{
			name:     "Quoted path with Unix-style argument",
			input:    `"C:\Program Files\MyApp\app" -config=test`,
			expected: `C:\Program Files\MyApp\app`,
		},
		{
			name:     "Path with double-dash argument",
			input:    `C:\Some Path\program --verbose`,
			expected: `C:\Some Path\program`,
		},
		{
			name:     "Path without arguments",
			input:    `C:\Program Files\MyApp\app`,
			expected: `C:\Program Files\MyApp\app`,
		},
		{
			name:     "Empty string",
			input:    "",
			expected: ".", // filepath.Clean returns "." for empty string
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := extractExecutablePath(tt.input)
			if got != tt.expected {
				t.Errorf("extractExePath(%q) = %q, want %q", tt.input, got, tt.expected)
			}
		})
	}
}
