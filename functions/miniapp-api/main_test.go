package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	_ "github.com/lib/pq"
)

func TestGetUserTagsHandler(t *testing.T) {
	// Skip if DATABASE_URL is not set
	if os.Getenv("DATABASE_URL") == "" {
		t.Skip("DATABASE_URL not set, skipping database test")
	}

	// Initialize test database
	testDB, err := initDB()
	if err != nil {
		t.Fatalf("Failed to connect to test database: %v", err)
	}
	defer testDB.Close()

	// Create router
	router := setupRoutes(testDB)

	// Test without authorization header
	t.Run("Missing Authorization", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/api/user/tags", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if w.Code != http.StatusUnauthorized {
			t.Errorf("Expected status %d, got %d", http.StatusUnauthorized, w.Code)
		}

		var response APIResponse
		if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
			t.Fatalf("Failed to unmarshal response: %v", err)
		}

		if response.Success {
			t.Error("Expected success to be false")
		}
	})

	// Test with invalid authorization
	t.Run("Invalid Authorization", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/api/user/tags", nil)
		req.Header.Set("Authorization", "Bearer invalid-data")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if w.Code != http.StatusUnauthorized {
			t.Errorf("Expected status %d, got %d", http.StatusUnauthorized, w.Code)
		}
	})

	// Test CORS headers
	t.Run("CORS Headers", func(t *testing.T) {
		req := httptest.NewRequest("OPTIONS", "/api/user/tags", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
		}

		expectedHeaders := map[string]string{
			"Access-Control-Allow-Origin":  "*",
			"Access-Control-Allow-Methods": "GET, POST, PUT, DELETE, OPTIONS",
			"Access-Control-Allow-Headers": "Origin, Content-Type, Authorization",
		}

		for header, expectedValue := range expectedHeaders {
			if actual := w.Header().Get(header); actual != expectedValue {
				t.Errorf("Expected header %s: %s, got: %s", header, expectedValue, actual)
			}
		}
	})
}

func TestGetUserTagsWithCounts(t *testing.T) {
	// Skip if DATABASE_URL is not set
	if os.Getenv("DATABASE_URL") == "" {
		t.Skip("DATABASE_URL not set, skipping database test")
	}

	testDB, err := initDB()
	if err != nil {
		t.Fatalf("Failed to connect to test database: %v", err)
	}
	defer testDB.Close()

	// Test with non-existent user (should return empty slice, not error)
	tags, err := getUserTagsWithCounts(testDB, 999999)
	if err != nil {
		t.Errorf("Expected no error for non-existent user, got: %v", err)
	}

	if len(tags) != 0 {
		t.Errorf("Expected empty tags for non-existent user, got %d tags", len(tags))
	}
}
