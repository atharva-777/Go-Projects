package store

import (
	"crypto/rand"
	"database/sql"
	"math/big"
	"time"

	_ "modernc.org/sqlite"
)

const alphabet = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

type URL struct {
	Code      string    `json:"code"`
	Original  string    `json:"original"`
	CreatedAt time.Time `json:"created_at"`
	Visits    int       `json:"visits"`
}

type Store struct {
	db *sql.DB
}

func New(dbPath string) (*Store, error) {
	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return nil, err
	}
	if err := db.Ping(); err != nil {
		return nil, err
	}
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS urls (
        code TEXT PRIMARY KEY,
        original TEXT NOT NULL,
        created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
        visits INTEGER DEFAULT 0
    )`)
	if err != nil {
		return nil, err
	}
	return &Store{db: db}, nil
}

func (s *Store) Create(original string) (*URL, error) {
	code := ""
	for {
		code = genCode(6)
		var exists bool
		err := s.db.QueryRow("SELECT 1 FROM urls WHERE code = ?", code).Scan(&exists)
		if err == sql.ErrNoRows {
			break
		} else if err != nil {
			return nil, err
		}
	}
	now := time.Now()
	_, err := s.db.Exec("INSERT INTO urls (code, original, created_at, visits) VALUES (?, ?, ?, ?)",
		code, original, now, 0)
	if err != nil {
		return nil, err
	}
	return &URL{Code: code, Original: original, CreatedAt: now, Visits: 0}, nil
}

func (s *Store) Get(code string) *URL {
	row := s.db.QueryRow("SELECT code, original, created_at, visits FROM urls WHERE code = ?", code)
	var u URL
	err := row.Scan(&u.Code, &u.Original, &u.CreatedAt, &u.Visits)
	if err == sql.ErrNoRows {
		return nil
	} else if err != nil {
		return nil
	}
	return &u
}

func (s *Store) Update(code, original string) bool {
	result, err := s.db.Exec("UPDATE urls SET original = ? WHERE code = ?", original, code)
	if err != nil {
		return false
	}
	rows, err := result.RowsAffected()
	return err == nil && rows > 0
}

func (s *Store) Delete(code string) bool {
	result, err := s.db.Exec("DELETE FROM urls WHERE code = ?", code)
	if err != nil {
		return false
	}
	rows, err := result.RowsAffected()
	return err == nil && rows > 0
}

func (s *Store) IncrementVisits(code string) error {
	_, err := s.db.Exec("UPDATE urls SET visits = visits + 1 WHERE code = ?", code)
	return err
}

func (s *Store) Close() error {
	return s.db.Close()
}

func genCode(n int) string {
	b := make([]byte, n)
	for i := 0; i < n; i++ {
		idx, _ := rand.Int(rand.Reader, big.NewInt(int64(len(alphabet))))
		b[i] = alphabet[idx.Int64()]
	}
	return string(b)
}
