package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/atharva-777/go-projects/url-shortner/store"
)

var s *store.Store

func init() {
	var err error
	s, err = store.New("urls.db")
	if err != nil {
		log.Fatalf("failed to initialize store: %v", err)
	}
}

func main() {
	mux := http.NewServeMux()
	// serve static UI at /ui/
	mux.Handle("/ui/", http.StripPrefix("/ui/", http.FileServer(http.Dir("static"))))

	mux.HandleFunc("/api/shorten", shortenHandler)
	mux.HandleFunc("/api/url/", apiURLHandler) // GET/PUT/DELETE
	mux.HandleFunc("/", rootHandler)           // redirect or health

	addr := ":8080"
	srv := &http.Server{Addr: addr, Handler: mux, ReadTimeout: 5 * time.Second}
	log.Printf("starting server on %s", addr)
	log.Fatal(srv.ListenAndServe())
}

// shortenHandler handles POST /api/shorten with JSON {"url":"..."}
func shortenHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	var req struct {
		URL string `json:"url"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.URL == "" {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}
	u, err := s.Create(req.URL)
	if err != nil {
		log.Printf("create error: %v", err)
		http.Error(w, "failed to create short URL", http.StatusInternalServerError)
		return
	}
	resp := map[string]interface{}{
		"code":      u.Code,
		"short_url": fmt.Sprintf("http://%s/%s", r.Host, u.Code),
		"original":  u.Original,
		"visits":    u.Visits,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

// apiURLHandler handles GET/PUT/DELETE for /api/url/{code}
func apiURLHandler(w http.ResponseWriter, r *http.Request) {
	// path: /api/url/{code} or /api/url/{code}/stats
	p := strings.TrimPrefix(r.URL.Path, "/api/url/")
	if p == "" || p == "/" {
		http.Error(w, "missing code", http.StatusBadRequest)
		return
	}
	parts := strings.Split(strings.Trim(p, "/"), "/")
	code := parts[0]
	if code == "" {
		http.Error(w, "missing code", http.StatusBadRequest)
		return
	}

	switch r.Method {
	case http.MethodGet:
		u := s.Get(code)
		if u == nil {
			http.NotFound(w, r)
			return
		}
		resp := map[string]interface{}{"code": u.Code, "original": u.Original, "visits": u.Visits, "created_at": u.CreatedAt}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	case http.MethodPut:
		var req struct {
			URL string `json:"url"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.URL == "" {
			http.Error(w, "invalid request body", http.StatusBadRequest)
			return
		}
		ok := s.Update(code, req.URL)
		if !ok {
			http.NotFound(w, r)
			return
		}
		w.WriteHeader(http.StatusNoContent)
	case http.MethodDelete:
		ok := s.Delete(code)
		if !ok {
			http.NotFound(w, r)
			return
		}
		w.WriteHeader(http.StatusNoContent)
	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}

// rootHandler serves redirect for /{code} or simple health for /
func rootHandler(w http.ResponseWriter, r *http.Request) {
	path := strings.Trim(r.URL.Path, "/")
	if path == "" {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("URL shortener is running"))
		return
	}
	// treat first segment as code
	code := strings.Split(path, "/")[0]
	u := s.Get(code)
	if u == nil {
		http.NotFound(w, r)
		return
	}
	if err := s.IncrementVisits(code); err != nil {
		log.Printf("increment visits error: %v", err)
	}
	http.Redirect(w, r, u.Original, http.StatusFound)
}
