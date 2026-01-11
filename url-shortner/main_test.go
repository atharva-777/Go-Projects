package main

import (
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/atharva-777/go-projects/url-shortner/store"
)

func init() {
	var err error
	s, err = store.New(":memory:")
	if err != nil {
		panic(err)
	}
}

func TestShortenHandler(t *testing.T) {
	req := httptest.NewRequest("POST", "/api/shorten",
		strings.NewReader(`{"url":"https://example.com"}`))
	w := httptest.NewRecorder()
	shortenHandler(w, req)
	if w.Code != 200 {
		t.Fatalf("expected 200, got %d", w.Code)
	}
}

func TestRootHandlerHealthCheck(t *testing.T) {
	req := httptest.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()
	rootHandler(w, req)
	if w.Code != 200 {
		t.Fatalf("expected 200, got %d", w.Code)
	}
}
