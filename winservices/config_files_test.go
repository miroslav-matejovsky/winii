package winservices

import (
	"os"
	"path/filepath"
	"testing"
)

func TestCollectServiceConfigFiles_Success(t *testing.T) {
	dir := t.TempDir()

	// create files
	files := map[string]string{
		"app.json":                         `{"name":"app"}`,
		"settings.XML":                     `<settings>ok</settings>`, // uppercase ext
		"readme.txt":                       `not a config`,
		filepath.Join("sub", "config.ini"): `key=value`,
	}
	for name, content := range files {
		p := filepath.Join(dir, name)
		if err := os.MkdirAll(filepath.Dir(p), 0o755); err != nil {
			t.Fatalf("mkdirall failed: %v", err)
		}
		if err := os.WriteFile(p, []byte(content), 0o644); err != nil {
			t.Fatalf("write file failed: %v", err)
		}
	}

	out, err := collectServiceConfigFiles(dir)
	if err != nil {
		t.Fatalf("collectServiceConfigFiles returned error: %v", err)
	}

	// expected config files: app.json, settings.XML, sub/config.ini
	expected := map[string]string{
		filepath.Join(dir, "app.json"):          `{"name":"app"}`,
		filepath.Join(dir, "settings.XML"):      `<settings>ok</settings>`,
		filepath.Join(dir, "sub", "config.ini"): `key=value`,
	}

	if len(out) != len(expected) {
		t.Fatalf("unexpected number of config files: got %d want %d", len(out), len(expected))
	}

	found := make(map[string]string, len(out))
	for _, cf := range out {
		found[cf.Path] = cf.Contents
	}

	for p, want := range expected {
		got, ok := found[p]
		if !ok {
			t.Errorf("expected file %q not found in result", p)
			continue
		}
		if got != want {
			t.Errorf("contents mismatch for %q: got %q want %q", p, got, want)
		}
	}
}

func TestCollectServiceConfigFiles_Errors(t *testing.T) {
	// empty dir param
	if _, err := collectServiceConfigFiles(""); err == nil {
		t.Error("expected error for empty dir, got nil")
	}

	// non-existent dir
	non := filepath.Join(t.TempDir(), "no-such-dir")
	if _, err := collectServiceConfigFiles(non); err == nil {
		t.Error("expected error for non-existent dir, got nil")
	}

	// path is a file
	tmp := t.TempDir()
	p := filepath.Join(tmp, "file.txt")
	if err := os.WriteFile(p, []byte("x"), 0o644); err != nil {
		t.Fatalf("write file failed: %v", err)
	}
	if _, err := collectServiceConfigFiles(p); err == nil {
		t.Error("expected error for path that is not a directory, got nil")
	}
}
