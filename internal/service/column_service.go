package service

import (
	"errors"
	"time"

	"github.com/ovk741/TasksStream/internal/domain"
	"github.com/ovk741/TasksStream/internal/storage"
)

type ColumnService interface {
	Create(userID, title string, boardID string) (domain.Column, error)
	GetByBoardID(userID, boardID string) ([]domain.Column, error)
	Update(userID, columnID string, title string) (domain.Column, error)
	Move(userID, columnID string, position int) (domain.Column, error)
	Delete(userID, columnID string) error
}

type columnService struct {
	columnRepo      storage.ColumnRepository
	boardRepo       storage.BoardRepository
	boardMemberRepo storage.BoardMemberRepository
	taskRepo        storage.TaskRepository
	generateID      func() string
}

func NewColumnService(
	columnRepo storage.ColumnRepository,
	boardRepo storage.BoardRepository,
	boardMemberRepo storage.BoardMemberRepository,
	taskRepo storage.TaskRepository,
	generateID func() string,
) ColumnService {
	return &columnService{
		columnRepo:      columnRepo,
		boardRepo:       boardRepo,
		boardMemberRepo: boardMemberRepo,
		taskRepo:        taskRepo,
		generateID:      generateID,
	}
}

func (s *columnService) Create(userID, title string, boardID string) (domain.Column, error) {
	if title == "" || boardID == "" {
		return domain.Column{}, domain.ErrInvalidInput
	}

	if err := s.requireBoardAccess(boardID, userID); err != nil {
		return domain.Column{}, err
	}

	_, err := s.boardRepo.GetByID(boardID)
	if err != nil {
		return domain.Column{}, domain.ErrNotFound
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

	if _, err := s.columnRepo.Create(column); err != nil {
		return domain.Column{}, err
	}

	return column, nil
}
func (s *columnService) GetByBoardID(userID, boardID string) ([]domain.Column, error) {

	_, err := s.boardRepo.GetByID(boardID)
	if err != nil {
		return nil, domain.ErrNotFound
	}

	if err := s.requireBoardAccess(boardID, userID); err != nil {
		return nil, err
	}

	column, err := s.columnRepo.GetByBoardID(boardID)
	if err != nil {
		return nil, err
	}
	return column, nil
}

func (s *columnService) Update(userID, columnID string, title string) (domain.Column, error) {
	if columnID == "" || title == "" {
		return domain.Column{}, domain.ErrInvalidInput
	}
	column, err := s.columnRepo.GetByID(columnID)
	if err != nil {
		return domain.Column{}, err
	}

	if err := s.requireBoardAccess(column.BoardID, userID); err != nil {
		return domain.Column{}, err
	}

	column.Title = title

	if _, err := s.columnRepo.Update(column); err != nil {
		return domain.Column{}, err
	}

	return column, nil
}

func (s *columnService) Delete(userID, columnID string) error {

	if columnID == "" {
		return domain.ErrInvalidInput
	}

	column, err := s.columnRepo.GetByID(columnID)
	if err != nil {
		return err
	}

	if err := s.requireBoardAccess(column.BoardID, userID); err != nil {
		return err
	}

	return s.columnRepo.Delete(columnID)

}

func (s *columnService) Move(userID, columnID string, position int) (domain.Column, error) {
	if columnID == "" || position < 0 {
		return domain.Column{}, domain.ErrInvalidInput
	}

	column, err := s.columnRepo.GetByID(columnID)
	if err != nil {
		return domain.Column{}, err
	}

	if err := s.requireBoardAccess(column.BoardID, userID); err != nil {
		return domain.Column{}, err
	}

	return s.columnRepo.Move(columnID, position)
}

func (s *columnService) requireBoardAccess(
	boardID, userID string,
) error {

	_, err := s.boardMemberRepo.GetRole(boardID, userID)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			return domain.ErrForbidden
		}
		return err
	}

	return nil
}
