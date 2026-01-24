package storage

import "github.com/ovk741/TasksStream/internal/domain"

type BoardRepository interface {
	Create(domain.Board)
	GetAll() []domain.Board
}
