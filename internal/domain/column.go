package domain

import "time"

type Column struct {
	ID        string    `json:"id"`
	Title     string    `json:"title"`
	Position  int       `json:"position"`
	BoardID   string    `json:"board_id"`
	CreatedAt time.Time `json:"created_at"`
}
