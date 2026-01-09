package store

import "testing"

func TestCreate(t *testing.T) {
	s := New()
	u := s.Create(("https://exmaple.com"))
	if u.Original != "https://exmaple.com" {
		t.Fatal("mismatch")
	}
	if u.Code == "" {
		t.Fatal("code empty")
	}
}

func TestGet(t *testing.T) {
	s := New()
	created := s.Create("https://exmaple.com")
	got := s.Get(created.Code)
	if got == nil {
		t.Fatal("not found")
	}
	if got.Original != created.Original {
		t.Fatal("mismatch")
	}
}

func TestUpdate(t *testing.T) {
	s := New()
	u := s.Create("https://example.com")
	ok := s.Update(u.Code, "https://new.com")
	if !ok {
		t.Fatal("update failed")
	}
	got := s.Get(u.Code)
	if got.Original != "https://new.com" {
		t.Fatal("not updated")
	}
}

func TestDelete(t *testing.T) {
	s := New()
	u := s.Create("https://example.com")
	ok := s.Delete(u.Code)
	if !ok {
		t.Fatal("delete failed")
	}
	if got := s.Get(u.Code); got != nil {
		t.Fatal("still exists")
	}
}

func TestIncrementVisits(t *testing.T) {
	s := New()
	u := s.Create("https://example.com")
	s.IncrementVisits(u.Code)
	if got := s.Get(u.Code); got.Visits != 1 {
		t.Fatal("visits not increamented")
	}
}
