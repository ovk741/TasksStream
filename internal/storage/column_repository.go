package storage

import "github.com/ovk741/TasksStream/internal/domain"

type ColumnRepository interface {
	Create(domain.Column)
	GetByBoardID(boardID string) []domain.Column
	GetByID(ColumnID string) (domain.Column, error)
	Update(column domain.Column) error
	Delete(columnID string) error
}
