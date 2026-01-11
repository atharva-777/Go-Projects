package store

import "testing"

func TestCreate(t *testing.T) {
	s, _ := New(":memory:")
	u, err := s.Create("https://example.com")
	if err != nil {
		t.Fatalf("create error: %v", err)
	}
	if u.Original != "https://example.com" {
		t.Fatal("mismatch")
	}
	if u.Code == "" {
		t.Fatal("code empty")
	}
}

func TestGet(t *testing.T) {
	s, _ := New(":memory:")
	created, _ := s.Create("https://example.com")
	got := s.Get(created.Code)
	if got == nil {
		t.Fatal("not found")
	}
	if got.Original != created.Original {
		t.Fatal("mismatch")
	}
}

func TestUpdate(t *testing.T) {
	s, _ := New(":memory:")
	u, _ := s.Create("https://example.com")
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
	s, _ := New(":memory:")
	u, _ := s.Create("https://example.com")
	ok := s.Delete(u.Code)
	if !ok {
		t.Fatal("delete failed")
	}
	if got := s.Get(u.Code); got != nil {
		t.Fatal("still exists")
	}
}

func TestIncrementVisits(t *testing.T) {
	s, _ := New(":memory:")
	u, _ := s.Create("https://example.com")
	err := s.IncrementVisits(u.Code)
	if err != nil {
		t.Fatalf("increment error: %v", err)
	}
	if got := s.Get(u.Code); got.Visits != 1 {
		t.Fatal("visits not incremented")
	}
}
