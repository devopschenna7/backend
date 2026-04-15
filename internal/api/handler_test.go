package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"ecommerce-app/internal/store" // Adjust import path based on your module name
)

func setupTestHandler() (*Handler, *store.MemoryStore) {
	s := store.NewMemoryStore()
	return NewHandler(s), s
}

func TestHandleSignup(t *testing.T) {
	h, _ := setupTestHandler()

	payload := map[string]string{
		"firstName": "John", "lastName": "Smith", "email": "john@smith.com",
		"mobile": "9999999999", "username": "johnsmith", "password": "pass123",
	}
	body, _ := json.Marshal(payload)

	req, _ := http.NewRequest("POST", "/api/signup", bytes.NewBuffer(body))
	rr := httptest.NewRecorder()

	h.HandleSignup(rr, req)

	if status := rr.Code; status != http.StatusCreated {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusCreated)
	}
}

func TestHandleLogin(t *testing.T) {
	h, s := setupTestHandler()

	// Pre-populate store
	s.Save(store.User{Username: "testuser", Password: "correctpassword"})

	tests := []struct {
		name       string
		payload    map[string]string
		wantStatus int
	}{
		{
			name:       "Valid Login",
			payload:    map[string]string{"username": "testuser", "password": "correctpassword"},
			wantStatus: http.StatusOK,
		},
		{
			name:       "Invalid Password",
			payload:    map[string]string{"username": "testuser", "password": "wrongpassword"},
			wantStatus: http.StatusUnauthorized,
		},
		{
			name:       "Non-existent User",
			payload:    map[string]string{"username": "ghost", "password": "password"},
			wantStatus: http.StatusUnauthorized,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body, _ := json.Marshal(tt.payload)
			req, _ := http.NewRequest("POST", "/api/login", bytes.NewBuffer(body))
			rr := httptest.NewRecorder()

			h.HandleLogin(rr, req)

			if status := rr.Code; status != tt.wantStatus {
				t.Errorf("handler returned wrong status code: got %v want %v", status, tt.wantStatus)
			}
		})
	}
}
