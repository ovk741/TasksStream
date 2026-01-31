package storage

import "github.com/ovk741/TasksStream/internal/domain"

type ColumnRepository interface {
	Create(domain.Column) (domain.Column, error)
	GetByBoardID(boardID string) ([]domain.Column, error)
	GetByID(ColumnID string) (domain.Column, error)
	Update(column domain.Column) (domain.Column, error)
	Delete(columnID string) error
	Move(columnID string, position int) (domain.Column, error)
}
