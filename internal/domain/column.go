package domain

import "time"

type Column struct {
	ID        string
	Title     string
	Position  int
	BoardID   string
	CreatedAt time.Time
}
