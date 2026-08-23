package utils_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/classroom-cli/internal/utils"
)

func TestLoadConfigFromEnv(t *testing.T) {
	os.Setenv("CLASSROOM_CLIENT_ID", "test-client-id")
	os.Setenv("CLASSROOM_CLIENT_SECRET", "test-client-secret")
	defer func() {
		os.Unsetenv("CLASSROOM_CLIENT_ID")
		os.Unsetenv("CLASSROOM_CLIENT_SECRET")
	}()

	cfg, src, err := utils.LoadConfig()
	if err != nil {
		t.Fatalf("expected config to load from env, got err: %v", err)
	}

	if cfg.ClientId != "test-client-id" || cfg.ClientSecret != "test-client-secret" {
		t.Errorf("unexpected config values: %+v", cfg)
	}

	if src != "environment variables" {
		t.Errorf("expected source 'environment variables', got '%s'", src)
	}
}

func TestStoreAndReadToken(t *testing.T) {
	tempHome := t.TempDir()
	origHome := os.Getenv("HOME")
	os.Setenv("HOME", tempHome)
	defer os.Setenv("HOME", origHome)

	tokenJSON := []byte(`{"access_token":"sample-access","refresh_token":"sample-refresh","expires_in":3600}`)
	err := utils.StoreTokenData(tokenJSON)
	if err != nil {
		t.Fatalf("StoreTokenData failed: %v", err)
	}

	stored := utils.ReadTokenData()
	if stored.AccessToken != "sample-access" || stored.RefreshToken != "sample-refresh" {
		t.Errorf("unexpected stored token data: %+v", stored)
	}

	// Verify file mode
	info, err := os.Stat(filepath.Join(tempHome, ".classroom", "secret.u"))
	if err != nil {
		t.Fatalf("could not stat secret file: %v", err)
	}
	if perm := info.Mode().Perm(); perm != 0600 {
		t.Errorf("expected permissions 0600, got %v", perm)
	}
}
