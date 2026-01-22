package model

import "time"

type Card struct {
	ID          int64
	Title       string
	Description string
	Position    string
	ColumnID    int64
	CreateAt    time.Time
}
