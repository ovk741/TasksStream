package storage

import "github.com/ovk741/TasksStream/internal/domain"

type BoardRepository interface {
	Create(domain.Board) (domain.Board, error)
	GetAll() ([]domain.Board, error)
	GetByID(boardID string) (domain.Board, error)
	Update(board domain.Board) (domain.Board, error)
	Delete(boardID string) error
}
