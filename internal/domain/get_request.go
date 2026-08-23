package domain

import (
	"fmt"
	"net/http"
	"time"
)

var httpClient = &http.Client{
	Timeout: 20 * time.Second,
}

func DoGetRequest(url, token string) (*http.Response, error) {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", "classroom-cli/1.0")

	res, err := httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}

	if res.StatusCode == http.StatusUnauthorized {
		_ = res.Body.Close()
		return nil, fmt.Errorf("authentication expired or invalid token (HTTP 401)")
	}

	if res.StatusCode == http.StatusForbidden {
		_ = res.Body.Close()
		return nil, fmt.Errorf("permission denied to access resource (HTTP 403)")
	}

	if res.StatusCode < 200 || res.StatusCode >= 300 {
		_ = res.Body.Close()
		return nil, fmt.Errorf("API request failed with HTTP %d %s", res.StatusCode, res.Status)
	}

	return res, nil
}
