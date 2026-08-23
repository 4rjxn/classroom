package utils

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

type TokenData struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	TokenType    string `json:"token_type,omitempty"`
	ExpiresIn    int    `json:"expires_in,omitempty"`
}

func getCredentialsDir() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("could not determine user home directory: %w", err)
	}
	return filepath.Join(home, ".classroom"), nil
}

func StoreTokenData(data []byte) error {
	dirPath, err := getCredentialsDir()
	if err != nil {
		return err
	}
	if err := os.MkdirAll(dirPath, 0700); err != nil {
		return fmt.Errorf("could not create credentials directory %s: %w", dirPath, err)
	}

	secretFile := filepath.Join(dirPath, "secret.u")

	// If data already contains a refresh_token, we write directly.
	// If new data has no refresh_token (e.g. during a refresh request), preserve existing refresh_token.
	var incoming TokenData
	if err := json.Unmarshal(data, &incoming); err == nil {
		if incoming.RefreshToken == "" {
			existing := ReadTokenData()
			if existing.RefreshToken != "" {
				incoming.RefreshToken = existing.RefreshToken
				if merged, err := json.Marshal(incoming); err == nil {
					data = merged
				}
			}
		}
	}

	if err := os.WriteFile(secretFile, data, 0600); err != nil {
		return fmt.Errorf("could not write secret file: %w", err)
	}
	return nil
}

func StoreRefreshToken(data []byte) {
	_ = StoreTokenData(data)
}

func ReadTokenData() TokenData {
	dirPath, err := getCredentialsDir()
	if err != nil {
		return TokenData{}
	}
	secretFile := filepath.Join(dirPath, "secret.u")
	data, err := os.ReadFile(secretFile)
	if err != nil {
		return TokenData{}
	}
	var val TokenData
	_ = json.Unmarshal(data, &val)
	return val
}

func ReadRefreshToken() string {
	return ReadTokenData().RefreshToken
}

func ClearToken() error {
	dirPath, err := getCredentialsDir()
	if err != nil {
		return err
	}
	secretFile := filepath.Join(dirPath, "secret.u")
	return os.Remove(secretFile)
}
