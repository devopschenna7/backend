package main

import (
	"log"
	"net/http"
	"time"

	"backend/internal/api"
	"backend/internal/store"
)

func main() {
	s := store.NewMemoryStore()
	h := api.NewHandler(s)

	mux := http.NewServeMux()

	mux.HandleFunc("/api/signup/check-availability", h.HandleCheckAvailability)
	mux.HandleFunc("/api/signup", h.HandleSignup)
	mux.HandleFunc("/api/login", h.HandleLogin)

	// Middleware for Logging and CORS
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// --- 1. LOG THE INCOMING REQUEST ---
		start := time.Now()
		log.Printf("--> %s %s", r.Method, r.URL.Path)

		// --- 2. HANDLE CORS ---
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		// If it's a preflight request, exit early with a 200 OK
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		// --- 3. SERVE THE ACTUAL REQUEST ---
		mux.ServeHTTP(w, r)

		// Optional: Log how long the request took
		log.Printf("<-- %s %s completed in %v", r.Method, r.URL.Path, time.Since(start))
	})

	log.Println("Backend listening on :8080")
	log.Fatal(http.ListenAndServe(":8080", handler))
}
