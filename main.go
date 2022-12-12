package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"sync"
	"time"
)

func mustEnv(key string) string {
	v, exists := os.LookupEnv(key)
	if !exists {
		log.Fatalf("Missing ENV %q", key)
	}

	return v
}

const (
	ListenAddr = ":8080"
)

var (
	Token   = mustEnv("TOKEN")
	mtx     sync.RWMutex
	timeout *time.Timer
	state   SpaceApi
)

func main() {
	if err := readState(); err != nil {
		log.Printf("failed reading initial state: %v", err)
	}

	state.Reset()

	http.HandleFunc("/spaceapi", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPut:
			handlePost(w, r)
		case http.MethodGet:
			handleGet(w, r)
		default:
			http.Error(w, "Method not implemented.", http.StatusNotImplemented)
		}
	})

	log.Println("Listening on", ListenAddr)
	log.Fatal(http.ListenAndServe(ListenAddr, nil))
}

func handleGet(w http.ResponseWriter, r *http.Request) {
	mtx.RLock()
	defer mtx.RUnlock()

	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(state); err != nil {
		// Ignore errors because they only happen when connection is lost
		return
	}
}

func handlePost(w http.ResponseWriter, r *http.Request) {
	if !isAuthorized(w, r) {
		return
	}

	var s SpaceApi
	if err := json.NewDecoder(r.Body).Decode(&s); err != nil {
		http.Error(w, "invalid data", http.StatusBadRequest)
		return
	}

	mtx.Lock()
	defer mtx.Unlock()

	state = s
	if timeout != nil {
		timeout.Stop()
	}

	timeout = time.NewTimer(5 * time.Minute)
	go func() {
		for range timeout.C {
			state.Reset()
		}
	}()

	if err := writeState(); err != nil {
		http.Error(w, "failed persisting state", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func readState() error {
	f, err := os.OpenFile("data.json", os.O_RDONLY, 0755)
	if err != nil {
		return fmt.Errorf("opening file: %v", err)
	}

	if err := json.NewDecoder(f).Decode(&state); err != nil {
		return fmt.Errorf("encoding state: %v", err)
	}

	return nil
}

func writeState() error {
	f, err := os.OpenFile("data.json", os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		return fmt.Errorf("opening file: %v", err)
	}

	if err := json.NewEncoder(f).Encode(state); err != nil {
		return fmt.Errorf("encoding state: %v", err)
	}

	return nil
}

func isAuthorized(w http.ResponseWriter, r *http.Request) bool {
	if _, exists := r.Header["Authorization"]; !exists {
		http.Error(w, "Needs authorization", http.StatusUnauthorized)
		return false
	}

	var t string
	_, err := fmt.Sscanf(r.Header.Get("Authorization"), "Token %s", &t)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return false
	}

	if t != Token {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return false
	}

	return true
}
