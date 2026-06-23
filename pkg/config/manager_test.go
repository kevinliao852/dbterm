package config

import (
	"errors"
	"os"
	"path/filepath"
	"testing"
)

func TestLoadConfigExpandsDefaultHomePath(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)

	manager := NewConfigManager(nil)
	if err := manager.LoadConfig(); err != nil {
		t.Fatalf("LoadConfig() error = %v", err)
	}

	want := filepath.Join(home, ".dbterm", "config.yaml")
	if manager.GetPath() != want {
		t.Fatalf("GetPath() = %q, want %q", manager.GetPath(), want)
	}

	info, err := os.Stat(want)
	if err != nil {
		t.Fatalf("config file was not created: %v", err)
	}
	if info.Mode().Perm() != 0600 {
		t.Fatalf("config permissions = %o, want 600", info.Mode().Perm())
	}
}

func TestLoadConfigCreatesParentForExplicitPath(t *testing.T) {
	path := filepath.Join(t.TempDir(), "nested", "config.yaml")
	manager := NewConfigManager(&path)

	if err := manager.LoadConfig(); err != nil {
		t.Fatalf("LoadConfig() error = %v", err)
	}
	if _, err := os.Stat(path); err != nil {
		t.Fatalf("config file was not created: %v", err)
	}
}

func TestLoadConfigRejectsEmptyPath(t *testing.T) {
	path := ""
	manager := NewConfigManager(&path)

	if err := manager.LoadConfig(); !errors.Is(err, ErrConfigNotFound) {
		t.Fatalf("LoadConfig() error = %v, want %v", err, ErrConfigNotFound)
	}
}
