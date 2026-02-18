package storage

import "github.com/ovk741/TasksStream/internal/domain"

type BoardMemberRepository interface {
	Add(member domain.BoardMember) error
	GetRole(boardID, userID string) (domain.BoardRole, error)
	IsMember(boardID, userID string) (bool, error)
	Remove(boardID, userID string) error
	GetMembers(boardID string) ([]domain.BoardMember, error)
}
