package store

import (
	"testing"
)

func TestMemoryStore(t *testing.T) {
	s := NewMemoryStore()
	user := User{
		FirstName: "Jane", LastName: "Doe", Email: "jane@test.com",
		Mobile: "1234567890", Username: "janedoe", Password: "password123",
	}

	// Test Save
	s.Save(user)

	// Validate user was saved
	s.mu.RLock()
	savedUser, exists := s.users["janedoe"]
	s.mu.RUnlock()
	if !exists {
		t.Fatalf("Expected user to be saved in store")
	}
	if savedUser.Email != "jane@test.com" {
		t.Errorf("Expected email jane@test.com, got %s", savedUser.Email)
	}

	// Test Exists Checks
	tests := []struct {
		name     string
		username string
		email    string
		mobile   string
		wantErr  string
	}{
		{"Duplicate Username", "janedoe", "new@test.com", "0000000000", "username taken"},
		{"Duplicate Email", "newuser", "jane@test.com", "0000000000", "email registered"},
		{"Duplicate Mobile", "newuser", "new@test.com", "1234567890", "mobile registered"},
		{"All Unique", "newuser", "new@test.com", "0000000000", ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := s.Exists(tt.username, tt.email, tt.mobile)
			if err != nil && err.Error() != tt.wantErr {
				t.Errorf("Exists() error = %v, wantErr %v", err, tt.wantErr)
			}
			if err == nil && tt.wantErr != "" {
				t.Errorf("Expected error %v, got nil", tt.wantErr)
			}
		})
	}
}
