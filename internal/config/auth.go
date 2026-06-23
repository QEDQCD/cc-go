package config

import (
	"crypto/rand"
	"encoding/json"
	"errors"
	"log"
	"math/big"
	"os"
	"path/filepath"
)

const (
	DefaultAuthUsername = "admin"
	authPasswordLength  = 18
)

func AuthPath() (string, error) {
	dir, err := ConfigDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, "auth.json"), nil
}

func LoadAuth(legacy *AuthConfig) (AuthConfig, error) {
	path, err := AuthPath()
	if err != nil {
		return AuthConfig{}, err
	}

	data, err := os.ReadFile(path)
	if err != nil {
		if !errors.Is(err, os.ErrNotExist) {
			return AuthConfig{}, err
		}
		auth, createErr := createAuthFile(path, legacy)
		if createErr != nil {
			return AuthConfig{}, createErr
		}
		return auth, nil
	}

	var auth AuthConfig
	if err := json.Unmarshal(data, &auth); err != nil {
		return AuthConfig{}, err
	}
	return ensureAuthFields(auth), nil
}

func SaveAuth(auth AuthConfig) error {
	path, err := AuthPath()
	if err != nil {
		return err
	}
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}
	auth = ensureAuthFields(auth)
	data, err := json.MarshalIndent(auth, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0600)
}

func createAuthFile(path string, legacy *AuthConfig) (AuthConfig, error) {
	auth := AuthConfig{Username: DefaultAuthUsername}
	if legacy != nil && legacy.Password != "" {
		auth.Username = legacy.Username
		if auth.Username == "" {
			auth.Username = DefaultAuthUsername
		}
		auth.Password = legacy.Password
	} else {
		password, err := generatePassword()
		if err != nil {
			return AuthConfig{}, err
		}
		auth.Password = password
		log.Printf("auth: created %s with generated credentials (username: %s)", path, auth.Username)
	}

	if err := SaveAuth(auth); err != nil {
		return AuthConfig{}, err
	}
	return auth, nil
}

func ensureAuthFields(auth AuthConfig) AuthConfig {
	if auth.Username == "" {
		auth.Username = DefaultAuthUsername
	}
	return auth
}

func extractLegacyAuth(data []byte) *AuthConfig {
	var aux struct {
		Auth *AuthConfig `json:"auth"`
	}
	if err := json.Unmarshal(data, &aux); err != nil || aux.Auth == nil {
		return nil
	}
	if aux.Auth.Username == "" && aux.Auth.Password == "" {
		return nil
	}
	return aux.Auth
}

func generatePassword() (string, error) {
	const chars = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789!@#$%^&*"
	password := make([]byte, authPasswordLength)
	max := big.NewInt(int64(len(chars)))
	for i := range password {
		n, err := rand.Int(rand.Reader, max)
		if err != nil {
			return "", err
		}
		password[i] = chars[n.Int64()]
	}
	return string(password), nil
}
