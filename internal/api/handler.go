package api

import (
	"encoding/json"
	"net/http"

	"backend/internal/store" // Imports your store package
)

type Handler struct {
	db *store.MemoryStore
}

func NewHandler(db *store.MemoryStore) *Handler {
	return &Handler{db: db}
}

func (h *Handler) HandleCheckAvailability(w http.ResponseWriter, r *http.Request) {
	username := r.URL.Query().Get("username")
	if err := h.db.Exists(username, "", ""); err != nil && err.Error() == "username taken" {
		http.Error(w, "Username taken", http.StatusConflict)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (h *Handler) HandleSignup(w http.ResponseWriter, r *http.Request) {
	var u store.User
	if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
		http.Error(w, "Invalid payload", http.StatusBadRequest)
		return
	}

	if err := h.db.Exists(u.Username, u.Email, u.Mobile); err != nil {
		http.Error(w, err.Error(), http.StatusConflict)
		return
	}

	h.db.Save(u)
	w.WriteHeader(http.StatusCreated)
}

func (h *Handler) HandleLogin(w http.ResponseWriter, r *http.Request) {
	var creds struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&creds); err != nil {
		http.Error(w, "Invalid payload", http.StatusBadRequest)
		return
	}

	// Capture both the user object and the boolean
	user, ok := h.db.Authenticate(creds.Username, creds.Password)
	if !ok {
		http.Error(w, "Invalid username or password", http.StatusUnauthorized)
		return
	}

	// Send back a JSON object containing the user's name
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"message":   "Login successful",
		"firstName": user.FirstName,
		"lastName":  user.LastName,
	})
}
