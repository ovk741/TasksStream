package service

import (
	"errors"
	"time"

	"github.com/ovk741/TasksStream/internal/domain"
	"github.com/ovk741/TasksStream/internal/storage"
)

type ColumnService interface {
	Create(title string, boardID string) (domain.Column, error)
	GetByBoardID(boardID string) ([]domain.Column, error)
}

type columnService struct {
	columnRepo storage.ColumnRepository
	boardRepo  storage.BoardRepository
	generateID func() string
}

func NewColumnService(columnRepo storage.ColumnRepository, boardRepo storage.BoardRepository, generateID func() string) ColumnService {
	return &columnService{
		columnRepo: columnRepo,
		boardRepo:  boardRepo,
		generateID: generateID,
	}
}

func (s *columnService) Create(title string, boardID string) (domain.Column, error) {
	if title == "" {
		return domain.Column{}, errors.New("column title is empty")
	}

	if !s.boardRepo.Exists(boardID) {
		return domain.Column{}, errors.New("board not found")
	}

	column := domain.Column{
		ID:        s.generateID(),
		Title:     title,
		BoardID:   boardID,
		CreatedAt: time.Now(),
	}

	s.columnRepo.Create(column)

	return column, nil
}
func (s *columnService) GetByBoardID(boardID string) ([]domain.Column, error) {

	if !s.boardRepo.Exists(boardID) {
		return nil, errors.New("board not found")
	}
	return s.columnRepo.GetByBoardID(boardID), nil
}
