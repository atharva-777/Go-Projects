package store

import (
    "crypto/rand"
    "math/big"
    "sync"
    "time"
)

const alphabet = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

type URL struct {
    Code      string    `json:"code"`
    Original  string    `json:"original"`
    CreatedAt time.Time `json:"created_at"`
    Visits    int       `json:"visits"`
}

type Store struct {
    mu   sync.RWMutex
    urls map[string]*URL
}

func New() *Store {
    return &Store{urls: make(map[string]*URL)}
}

func (s *Store) Create(original string) *URL {
    s.mu.Lock()
    defer s.mu.Unlock()
    code := ""
    for {
        code = genCode(6)
        if _, ok := s.urls[code]; !ok {
            break
        }
    }
    u := &URL{Code: code, Original: original, CreatedAt: time.Now(), Visits: 0}
    s.urls[code] = u
    return u
}

func (s *Store) Get(code string) *URL {
    s.mu.RLock()
    defer s.mu.RUnlock()
    if u, ok := s.urls[code]; ok {
        return &URL{Code: u.Code, Original: u.Original, CreatedAt: u.CreatedAt, Visits: u.Visits}
    }
    return nil
}

func (s *Store) Update(code, original string) bool {
    s.mu.Lock()
    defer s.mu.Unlock()
    if u, ok := s.urls[code]; ok {
        u.Original = original
        return true
    }
    return false
}

func (s *Store) Delete(code string) bool {
    s.mu.Lock()
    defer s.mu.Unlock()
    if _, ok := s.urls[code]; ok {
        delete(s.urls, code)
        return true
    }
    return false
}

func (s *Store) IncrementVisits(code string) {
    s.mu.Lock()
    defer s.mu.Unlock()
    if u, ok := s.urls[code]; ok {
        u.Visits++
    }
}

func genCode(n int) string {
    b := make([]byte, n)
    for i := 0; i < n; i++ {
        idx, _ := rand.Int(rand.Reader, big.NewInt(int64(len(alphabet))))
        b[i] = alphabet[idx.Int64()]
    }
    return string(b)
}
