package task

import "time"

type Task struct {
	ID          int
	Title       string
	Description string
	Status      string
	CreatedAt   time.Time
}
