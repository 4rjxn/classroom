package domain_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/4rjxn/classroom/internal/domain"
	"github.com/4rjxn/classroom/internal/models"
)

func TestDoGetRequest(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader != "Bearer valid-token" {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		resp := models.CourseResponse{
			Courses: []models.CourseModel{
				{Id: "101", Name: "Math 101", Sub: "Calculus"},
			},
		}
		_ = json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	// Test authorized request
	res, err := domain.DoGetRequest(server.URL, "valid-token")
	if err != nil {
		t.Fatalf("expected successful request, got: %v", err)
	}
	_ = res.Body.Close()

	// Test unauthorized request
	_, err = domain.DoGetRequest(server.URL, "invalid-token")
	if err == nil {
		t.Fatal("expected 401 error, got nil")
	}
}
