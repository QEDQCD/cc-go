package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadAuthCreatesFile(t *testing.T) {
	origHome := os.Getenv("HOME")
	origUserProfile := os.Getenv("USERPROFILE")
	tmpDir := t.TempDir()
	os.Setenv("HOME", tmpDir)
	os.Setenv("USERPROFILE", tmpDir)
	defer func() {
		os.Setenv("HOME", origHome)
		os.Setenv("USERPROFILE", origUserProfile)
	}()

	auth, err := LoadAuth(nil)
	if err != nil {
		t.Fatalf("LoadAuth failed: %v", err)
	}
	if auth.Username != DefaultAuthUsername {
		t.Fatalf("expected username %s, got %s", DefaultAuthUsername, auth.Username)
	}
	if len(auth.Password) < authPasswordLength {
		t.Fatalf("expected generated password length >= %d, got %d", authPasswordLength, len(auth.Password))
	}

	path, err := AuthPath()
	if err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(path); err != nil {
		t.Fatalf("auth file not created: %v", err)
	}
}

func TestLoadAuthMigratesLegacyConfig(t *testing.T) {
	origHome := os.Getenv("HOME")
	origUserProfile := os.Getenv("USERPROFILE")
	tmpDir := t.TempDir()
	os.Setenv("HOME", tmpDir)
	os.Setenv("USERPROFILE", tmpDir)
	defer func() {
		os.Setenv("HOME", origHome)
		os.Setenv("USERPROFILE", origUserProfile)
	}()

	legacy := &AuthConfig{Username: "admin", Password: "legacy-pass-123"}
	auth, err := LoadAuth(legacy)
	if err != nil {
		t.Fatalf("LoadAuth failed: %v", err)
	}
	if auth.Password != legacy.Password {
		t.Fatalf("expected migrated password %q, got %q", legacy.Password, auth.Password)
	}

	path := filepath.Join(tmpDir, ".cc-go", "auth.json")
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if string(data) == "" {
		t.Fatal("expected auth file content")
	}
}
