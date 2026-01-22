package model

import "time"

type Board struct {
	ID        int64
	Name      string
	CreatedAt time.Time
}
