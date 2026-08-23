package auth

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/4rjxn/classroom/internal/models"
	"github.com/4rjxn/classroom/internal/utils"
	"github.com/pkg/browser"
)

const (
	tokenURL    = "https://oauth2.googleapis.com/token"
	authBaseURL = "https://accounts.google.com/o/oauth2/v2/auth"
	redirectURI = "http://localhost:4321"
)

var defaultScopes = []string{
	"https://www.googleapis.com/auth/classroom.courses.readonly",
	"https://www.googleapis.com/auth/classroom.courseworkmaterials.readonly",
	"https://www.googleapis.com/auth/classroom.coursework.me.readonly",
	"https://www.googleapis.com/auth/classroom.coursework.students.readonly",
	"https://www.googleapis.com/auth/classroom.announcements.readonly",
	"https://www.googleapis.com/auth/classroom.student-submissions.me.readonly",
}

type AuthResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int    `json:"expires_in"`
	TokenType    string `json:"token_type"`
	Error        string `json:"error,omitempty"`
	ErrorDesc    string `json:"error_description,omitempty"`
}

// Backward compatibility alias for any existing caller
type AuthResponce = AuthResponse

// OfflineGeneration checks if a stored refresh token exists and attempts to refresh it.
// If refresh fails or no token is found, it starts the interactive OAuth browser flow.
func OfflineGeneration(config models.Config) (string, error) {
	refreshToken := utils.ReadRefreshToken()
	if refreshToken != "" {
		token, err := refreshAccessToken(config, refreshToken)
		if err == nil && token != "" {
			return token, nil
		}
		// If refresh failed (e.g. revoked), proceed to GenerateToken
		fmt.Println("Stored token invalid or expired. Re-authenticating...")
	}
	return GenerateToken(config)
}

// OffileGeneration is kept for backward compatibility with old spelling
func OffileGeneration(config models.Config) string {
	token, _ := OfflineGeneration(config)
	return token
}

func refreshAccessToken(config models.Config, refreshToken string) (string, error) {
	client := &http.Client{Timeout: 15 * time.Second}
	res, err := client.PostForm(tokenURL, url.Values{
		"refresh_token": {refreshToken},
		"client_id":     {config.ClientId},
		"client_secret": {config.ClientSecret},
		"grant_type":    {"refresh_token"},
	})
	if err != nil {
		return "", fmt.Errorf("token refresh request failed: %w", err)
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read refresh response: %w", err)
	}

	if res.StatusCode != http.StatusOK {
		return "", fmt.Errorf("refresh token failed with status %d: %s", res.StatusCode, string(body))
	}

	var authResp AuthResponse
	if err := json.Unmarshal(body, &authResp); err != nil {
		return "", fmt.Errorf("failed to parse refresh response: %w", err)
	}

	if authResp.AccessToken == "" {
		return "", fmt.Errorf("no access token returned in refresh response")
	}

	// Update stored token (preserves refresh token if not returned)
	_ = utils.StoreTokenData(body)
	return authResp.AccessToken, nil
}

// GenerateToken performs the full OAuth 2.0 authorization code flow via local HTTP server.
func GenerateToken(config models.Config) (string, error) {
	authURL, err := url.Parse(authBaseURL)
	if err != nil {
		return "", fmt.Errorf("invalid auth base URL: %w", err)
	}

	q := authURL.Query()
	q.Add("client_id", config.ClientId)
	q.Add("redirect_uri", redirectURI)
	q.Add("response_type", "code")
	q.Add("access_type", "offline")
	q.Add("scope", strings.Join(defaultScopes, " "))
	q.Add("prompt", "consent")
	authURL.RawQuery = q.Encode()

	codeChan := make(chan string, 1)
	errChan := make(chan error, 1)

	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if authErr := r.URL.Query().Get("error"); authErr != "" {
			desc := r.URL.Query().Get("error_description")
			w.Header().Set("Content-Type", "text/html; charset=utf-8")
			w.WriteHeader(http.StatusBadRequest)
			_, _ = fmt.Fprintf(w, `<!DOCTYPE html><html><body style="font-family:system-ui;text-align:center;padding:40px;">
<h2 style="color:#ef4444;">Authorization Failed</h2>
<p>%s: %s</p>
<p>You may return to your terminal.</p>
</body></html>`, authErr, desc)
			errChan <- fmt.Errorf("OAuth error: %s (%s)", authErr, desc)
			return
		}

		code := r.URL.Query().Get("code")
		if code == "" {
			w.WriteHeader(http.StatusBadRequest)
			_, _ = fmt.Fprint(w, "No authorization code found.")
			return
		}

		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		_, _ = fmt.Fprint(w, `<!DOCTYPE html><html><body style="font-family:system-ui;text-align:center;padding:40px;">
<h2 style="color:#10b981;">✓ Authorization Successful</h2>
<p>Google Classroom CLI is now connected. You can close this tab and return to your terminal.</p>
</body></html>`)

		codeChan <- code
	})

	listener, err := net.Listen("tcp", ":4321")
	if err != nil {
		return "", fmt.Errorf("failed to start local auth callback server on :4321 (is another instance running?): %w", err)
	}

	server := &http.Server{Handler: mux}

	go func() {
		if serveErr := server.Serve(listener); serveErr != nil && serveErr != http.ErrServerClosed {
			errChan <- serveErr
		}
	}()

	fullAuthURL := authURL.String()
	fmt.Println("Opening browser for Google Classroom authorization...")
	fmt.Println("If the browser does not open automatically, visit this URL:")
	fmt.Printf("\n  %s\n\n", fullAuthURL)

	_ = browser.OpenURL(fullAuthURL)

	var code string
	select {
	case code = <-codeChan:
	case authErr := <-errChan:
		_ = server.Shutdown(context.Background())
		return "", authErr
	case <-time.After(3 * time.Minute):
		_ = server.Shutdown(context.Background())
		return "", fmt.Errorf("authorization timed out after 3 minutes")
	}

	// Shutdown callback server gracefully
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	_ = server.Shutdown(shutdownCtx)

	// Exchange authorization code for token
	client := &http.Client{Timeout: 15 * time.Second}
	res, err := client.PostForm(tokenURL, url.Values{
		"code":          {code},
		"client_id":     {config.ClientId},
		"client_secret": {config.ClientSecret},
		"redirect_uri":  {redirectURI},
		"grant_type":    {"authorization_code"},
	})
	if err != nil {
		return "", fmt.Errorf("token exchange request failed: %w", err)
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read token exchange response: %w", err)
	}

	if res.StatusCode != http.StatusOK {
		return "", fmt.Errorf("token exchange failed with status %d: %s", res.StatusCode, string(body))
	}

	var authResp AuthResponse
	if err := json.Unmarshal(body, &authResp); err != nil {
		return "", fmt.Errorf("failed to parse token response: %w", err)
	}

	if authResp.AccessToken == "" {
		return "", fmt.Errorf("empty access token in exchange response: %s", string(body))
	}

	_ = utils.StoreTokenData(body)
	return authResp.AccessToken, nil
}
