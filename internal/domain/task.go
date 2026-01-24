package domain

import "time"

type Task struct {
	ID          string
	ColumnID    string
	Title       string
	Description string
	Position    int
	CreatedAt   time.Time
}
