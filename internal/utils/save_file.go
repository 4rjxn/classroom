package utils

import (
	"encoding/json"
	"os"
	"path/filepath"
)

func StoreRefreshToken(data []byte) {
	home, err := os.UserHomeDir()
	if err != nil {
		panic("os err no home")
	}
	dirPath := filepath.Join(home, ".classroom")
	if err = os.Mkdir(dirPath, 0755); err != nil {
		panic("cannot create folder err")
	}
	if err = os.WriteFile(filepath.Join(dirPath, "secret.u"), data, 0644); err != nil {
		panic("cannot WriteFile")
	}

}

type refreshToken struct {
	Token string `json:"refresh_token"`
}

func ReadRefreshToken() string {
	home, err := os.UserHomeDir()
	if err != nil {
		panic("os err no home")
	}
	fileName := filepath.Join(home, ".classroom", "secret.u")
	data, err := os.ReadFile(fileName)
	if err != nil {
		return ""
	}
	val := refreshToken{}
	json.Unmarshal(data, &val)
	return val.Token

}
