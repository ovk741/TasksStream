package service

import (
	"time"

	"github.com/ovk741/TasksStream/internal/domain"
	"github.com/ovk741/TasksStream/internal/storage"
)

type ColumnService interface {
	Create(title string, boardID string) (domain.Column, error)
	GetByBoardID(boardID string) ([]domain.Column, error)
	Update(columnID string, title string) (domain.Column, error)
	Move(columnID string, position int) (domain.Column, error)
	Delete(columnID string) error
}

type columnService struct {
	columnRepo storage.ColumnRepository
	boardRepo  storage.BoardRepository
	taskRepo   storage.TaskRepository
	generateID func() string
}

func NewColumnService(
	columnRepo storage.ColumnRepository,
	boardRepo storage.BoardRepository,
	taskRepo storage.TaskRepository,
	generateID func() string,
) ColumnService {
	return &columnService{
		columnRepo: columnRepo,
		boardRepo:  boardRepo,
		taskRepo:   taskRepo,
		generateID: generateID,
	}
}

func (s *columnService) Create(title string, boardID string) (domain.Column, error) {
	if title == "" {
		return domain.Column{}, ErrInvalidInput
	}

	_, err := s.boardRepo.GetByID(boardID)
	if err != nil {
		return domain.Column{}, ErrNotFound
	}
	columns, err := s.columnRepo.GetByBoardID(boardID)
	if err != nil {
		return domain.Column{}, err
	}

	column := domain.Column{
		ID:        s.generateID(),
		Title:     title,
		BoardID:   boardID,
		Position:  len(columns),
		CreatedAt: time.Now(),
	}

	s.columnRepo.Create(column)

	return column, nil
}
func (s *columnService) GetByBoardID(boardID string) ([]domain.Column, error) {

	_, err := s.boardRepo.GetByID(boardID)
	if err != nil {
		return nil, ErrNotFound
	}
	column, err := s.columnRepo.GetByBoardID(boardID)
	if err != nil {
		return nil, err
	}
	return column, nil
}

func (s *columnService) Update(columnID string, title string) (domain.Column, error) {
	if columnID == "" || title == "" {
		return domain.Column{}, ErrInvalidInput
	}
	column, err := s.columnRepo.GetByID(columnID)
	if err != nil {
		return domain.Column{}, err
	}

	column.Title = title

	if _, err := s.columnRepo.Update(column); err != nil {
		return domain.Column{}, err
	}

	return column, nil
}

func (s *columnService) Delete(columnID string) error {

	if columnID == "" {
		return ErrInvalidInput
	}

	_, err := s.columnRepo.GetByID(columnID)
	if err != nil {
		return err
	}

	return s.columnRepo.Delete(columnID)

}

func (s *columnService) Move(columnID string, position int) (domain.Column, error) {
	if columnID == "" || position < 0 {
		return domain.Column{}, ErrInvalidInput
	}

	return s.columnRepo.Move(columnID, position)
}
