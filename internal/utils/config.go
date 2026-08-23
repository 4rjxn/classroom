package utils

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/BurntSushi/toml"
	"github.com/classroom-cli/internal/models"
)

// LoadConfig attempts to load the configuration from several standard locations or environment variables.
func LoadConfig() (models.Config, string, error) {
	var cfg models.Config

	// 1. Check environment variables first
	envID := os.Getenv("CLASSROOM_CLIENT_ID")
	envSecret := os.Getenv("CLASSROOM_CLIENT_SECRET")
	if envID != "" && envSecret != "" {
		cfg.ClientId = envID
		cfg.ClientSecret = envSecret
		return cfg, "environment variables", nil
	}

	// 2. Candidate file paths
	home, _ := os.UserHomeDir()
	candidates := []string{
		"./config.toml",
	}
	if home != "" {
		candidates = append(candidates,
			filepath.Join(home, ".config", "classroom", "config.toml"),
			filepath.Join(home, ".classroom", "config.toml"),
		)
	}

	for _, path := range candidates {
		if _, err := os.Stat(path); err == nil {
			var parsed models.Config
			if _, decodeErr := toml.DecodeFile(path, &parsed); decodeErr != nil {
				return cfg, "", fmt.Errorf("error reading config file '%s': %w\nPlease ensure 'client_id' and 'client_secret' are correctly formatted", path, decodeErr)
			}
			if parsed.ClientId != "" && parsed.ClientSecret != "" {
				return parsed, path, nil
			}
			return cfg, "", fmt.Errorf("config file '%s' is missing 'client_id' or 'client_secret'", path)
		}
	}

	helpMsg := fmt.Sprintf(`Google Classroom CLI Configuration Not Found.

Please create a 'config.toml' file in one of the following locations:
  • ./config.toml
  • ~/.config/classroom/config.toml
  • ~/.classroom/config.toml

Or export the following environment variables:
  export CLASSROOM_CLIENT_ID="your-client-id.apps.googleusercontent.com"
  export CLASSROOM_CLIENT_SECRET="your-client-secret"

Example 'config.toml' format:
  client_id = "your-google-oauth-client-id"
  client_secret = "your-google-oauth-client-secret"

To obtain Google OAuth Credentials:
  1. Go to Google Cloud Console: https://console.cloud.google.com/
  2. Create a Project and enable the Google Classroom API.
  3. Configure OAuth Consent Screen (add scope: classroom.courses.readonly, etc.).
  4. Create OAuth 2.0 Client ID of type "Desktop app" or "Web application" with redirect URI http://localhost:4321.
`)

	return cfg, "", errors.New(helpMsg)
}
