package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"sync"
)

// secretVersion holds a single version's payload.
type secretVersion struct {
	Data []byte
}

// secret holds all versions of a secret.
type secret struct {
	Name     string
	Versions []secretVersion
}

// store is a thread-safe in-memory secret store keyed by full secret name
// (projects/{project}/secrets/{secret}).
type store struct {
	mu      sync.RWMutex
	secrets map[string]*secret
}

func newStore() *store {
	return &store{secrets: make(map[string]*secret)}
}

func main() {
	s := newStore()
	mux := http.NewServeMux()

	mux.HandleFunc("GET /healthz", handleHealthz)
	mux.HandleFunc("/v1/", s.handleV1)

	log.Println("mock-secrets listening on :8088")
	if err := http.ListenAndServe(":8088", mux); err != nil {
		log.Fatalf("server error: %v", err)
	}
}

func handleHealthz(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, "ok")
}

// handleV1 routes requests under /v1/projects/... to the correct handler.
func (s *store) handleV1(w http.ResponseWriter, r *http.Request) {
	// Strip the /v1/ prefix for easier parsing.
	path := strings.TrimPrefix(r.URL.Path, "/v1/")

	// Route: POST projects/{project}/secrets/{secret}:addVersion
	if r.Method == http.MethodPost && strings.HasSuffix(path, ":addVersion") {
		s.addVersion(w, r, strings.TrimSuffix(path, ":addVersion"))
		return
	}

	// Route: POST projects/{project}/secrets — create secret
	if r.Method == http.MethodPost {
		s.createSecret(w, r, path)
		return
	}

	// Route: GET projects/{project}/secrets/{secret}/versions/{version}:access
	if r.Method == http.MethodGet && strings.HasSuffix(path, ":access") {
		s.accessVersion(w, r, strings.TrimSuffix(path, ":access"))
		return
	}

	http.NotFound(w, r)
}

// createSecret handles POST /v1/projects/{project}/secrets?secretId=<id>
func (s *store) createSecret(w http.ResponseWriter, r *http.Request, path string) {
	// path = "projects/{project}/secrets"
	parts := strings.Split(path, "/")
	if len(parts) != 3 || parts[0] != "projects" || parts[2] != "secrets" {
		http.Error(w, "invalid path", http.StatusBadRequest)
		return
	}
	project := parts[1]
	secretID := r.URL.Query().Get("secretId")
	if secretID == "" {
		http.Error(w, "missing secretId query parameter", http.StatusBadRequest)
		return
	}

	name := fmt.Sprintf("projects/%s/secrets/%s", project, secretID)

	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.secrets[name]; exists {
		http.Error(w, "secret already exists", http.StatusConflict)
		return
	}

	s.secrets[name] = &secret{Name: name}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"name": name})
}

// addVersionRequest is the expected JSON body for adding a version.
type addVersionRequest struct {
	Payload struct {
		Data string `json:"data"` // base64-encoded
	} `json:"payload"`
}

// addVersion handles POST /v1/projects/{project}/secrets/{secret}:addVersion
func (s *store) addVersion(w http.ResponseWriter, r *http.Request, path string) {
	// path = "projects/{project}/secrets/{secret}"
	parts := strings.Split(path, "/")
	if len(parts) != 4 || parts[0] != "projects" || parts[2] != "secrets" {
		http.Error(w, "invalid path", http.StatusBadRequest)
		return
	}
	name := strings.Join(parts, "/")

	var req addVersionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	data, err := base64.StdEncoding.DecodeString(req.Payload.Data)
	if err != nil {
		http.Error(w, "invalid base64 payload", http.StatusBadRequest)
		return
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	sec, exists := s.secrets[name]
	if !exists {
		http.Error(w, "secret not found", http.StatusNotFound)
		return
	}

	sec.Versions = append(sec.Versions, secretVersion{Data: data})
	version := len(sec.Versions)
	versionName := fmt.Sprintf("%s/versions/%d", name, version)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"name": versionName})
}

// accessVersion handles GET /v1/projects/{project}/secrets/{secret}/versions/{version}:access
func (s *store) accessVersion(w http.ResponseWriter, _ *http.Request, path string) {
	// path = "projects/{project}/secrets/{secret}/versions/{version}"
	parts := strings.Split(path, "/")
	if len(parts) != 6 || parts[0] != "projects" || parts[2] != "secrets" || parts[4] != "versions" {
		http.Error(w, "invalid path", http.StatusBadRequest)
		return
	}

	secretName := strings.Join(parts[:4], "/")
	versionStr := parts[5]

	s.mu.RLock()
	defer s.mu.RUnlock()

	sec, exists := s.secrets[secretName]
	if !exists {
		http.Error(w, "secret not found", http.StatusNotFound)
		return
	}

	var idx int
	if versionStr == "latest" {
		if len(sec.Versions) == 0 {
			http.Error(w, "no versions", http.StatusNotFound)
			return
		}
		idx = len(sec.Versions) - 1
	} else {
		if _, err := fmt.Sscanf(versionStr, "%d", &idx); err != nil {
			http.Error(w, "invalid version", http.StatusBadRequest)
			return
		}
		idx-- // versions are 1-indexed
	}

	if idx < 0 || idx >= len(sec.Versions) {
		http.Error(w, "version not found", http.StatusNotFound)
		return
	}

	versionName := fmt.Sprintf("%s/versions/%d", secretName, idx+1)
	payload := base64.StdEncoding.EncodeToString(sec.Versions[idx].Data)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]any{
		"name": versionName,
		"payload": map[string]string{
			"data": payload,
		},
	})
}
