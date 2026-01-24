package domain

import "time"

type Board struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"Created_at"`
}
