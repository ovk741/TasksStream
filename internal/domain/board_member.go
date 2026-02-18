package domain

import "time"

type BoardRole string

const (
	BoardRoleOwner  BoardRole = "owner"
	BoardRoleEditor BoardRole = "editor"
	BoardRoleViewer BoardRole = "viewer"
)

type BoardMember struct {
	ID        string
	BoardID   string
	UserID    string
	Role      BoardRole
	CreatedAt time.Time
}
