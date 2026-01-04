package store

import (
	"encoding/json"
	"io"
	"os"
	"sync"
	"time"

	"github.com/atharva-777/go-projects/task-tracker/task"
)

const (
	Todo       = "todo"
	InProgress = "inprogress"
	Done       = "done"
)

type Store struct {
	path  string
	mu    sync.Mutex
	tasks []task.Task
	next  int
}

func (s *Store) load() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	f, err := os.OpenFile(s.path, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	defer f.Close()

	var ts []task.Task
	dec := json.NewDecoder(f)
	if err := dec.Decode(&ts); err != nil && err != io.EOF {
		// ignore decode errors for empty file
	}

	s.tasks = ts
	s.next = 1
	for _, t := range s.tasks {
		if t.ID >= s.next {
			s.next = t.ID + 1
		}
	}
	return nil
}

func New(path string) *Store {
	s := &Store{path: path}
	_ = s.load()
	return s
}

// saveLocked writes tasks to disk and MUST be called with s.mu held.
func (s *Store) saveLocked() error {
	f, err := os.Create(s.path)
	if err != nil {
		return err
	}
	defer f.Close()

	enc := json.NewEncoder(f)
	enc.SetIndent("", "  ")
	return enc.Encode(s.tasks)
}

// save is the public save that grabs the lock and calls saveLocked.
func (s *Store) save() error {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.saveLocked()
}

func (s *Store) Add(title, desc string) task.Task {
	s.mu.Lock()
	defer s.mu.Unlock()
	t := task.Task{
		ID: s.next, Title: title, Description: desc, Status: Todo, CreatedAt: time.Now(),
	}

	s.next++
	s.tasks = append(s.tasks, t)
	_ = s.saveLocked()
	return t
}

func (s *Store) Update(id int, title, desc string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()

	for i := range s.tasks {
		if s.tasks[i].ID == id {
			s.tasks[i].Title = title
			s.tasks[i].Description = desc
			_ = s.saveLocked()
			return true
		}
	}
	return false
}

func (s *Store) Delete(id int) bool {
	s.mu.Lock()
	defer s.mu.Unlock()

	for i := range s.tasks {
		if s.tasks[i].ID == id {
			s.tasks = append(s.tasks[:i], s.tasks[i+1:]...)
			_ = s.saveLocked()
			return true
		}
	}
	return false
}

func (s *Store) SetStatus(id int, status string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()

	for i := range s.tasks {
		if s.tasks[i].ID == id {
			s.tasks[i].Status = status
			_ = s.saveLocked()
			return true
		}
	}
	return false
}

func (s *Store) List(filter string) []task.Task {
	s.mu.Lock()
	defer s.mu.Unlock()
	if filter == "all" {
		return append([]task.Task(nil), s.tasks...)
	}

	var out []task.Task

	for _, t := range s.tasks {
		if (filter == "todo" && t.Status == Todo) ||
			(filter == "inprogress" && t.Status == InProgress) ||
			(filter == "done" && t.Status == Done) {
			out = append(out, t)
		}
	}
	return out
}
