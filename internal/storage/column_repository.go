package storage

import "github.com/ovk741/TasksStream/internal/domain"

type ColumnRepository interface {
	Create(domain.Column)
	GetByBoardID(boardID string) []domain.Column
	Exists(ColumnID string) bool
}
