package domain

import "time"

type Task struct {
	ID          string    `json:"id"`
	ColumnID    string    `json:"column_id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Position    int       `json:"position"`
	CreatedAt   time.Time `json:"created_at"`
}
