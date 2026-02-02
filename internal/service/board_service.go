package service

import (
	"time"

	"github.com/ovk741/TasksStream/internal/domain"
	"github.com/ovk741/TasksStream/internal/storage"
)

type BoardService interface {
	Create(name string) (domain.Board, error)
	GetAll() ([]domain.Board, error)
	Update(boardID string, name string) (domain.Board, error)
	Delete(boardID string) error
}

type boardService struct {
	boardRepo  storage.BoardRepository
	columnRepo storage.ColumnRepository
	taskRepo   storage.TaskRepository
	generateID func() string
}

func NewBoardService(
	boardRepo storage.BoardRepository,
	columnRepo storage.ColumnRepository,
	taskRepo storage.TaskRepository,
	generateID func() string,
) BoardService {
	return &boardService{
		boardRepo:  boardRepo,
		generateID: generateID,
		taskRepo:   taskRepo,
	}

}

func (s *boardService) Create(name string) (domain.Board, error) {
	if name == "" {
		return domain.Board{}, domain.ErrInvalidInput
	}

	board := domain.Board{
		ID:        s.generateID(),
		Name:      name,
		CreatedAt: time.Now(),
	}

	if _, err := s.boardRepo.Create(board); err != nil {
		return domain.Board{}, err
	}

	return board, nil
}
func (s *boardService) GetAll() ([]domain.Board, error) {

	board, err := s.boardRepo.GetAll()
	if err != nil {
		return nil, err
	}

	return board, nil
}

func (s *boardService) Update(boardID string, name string) (domain.Board, error) {
	if boardID == "" || name == "" {
		return domain.Board{}, domain.ErrInvalidInput
	}
	board, err := s.boardRepo.GetByID(boardID)
	if err != nil {
		return domain.Board{}, err
	}

	board.Name = name

	if _, err := s.boardRepo.Update(board); err != nil {
		return domain.Board{}, err
	}

	return board, nil

}

func (s *boardService) Delete(boardID string) error {
	if boardID == "" {
		return domain.ErrInvalidInput
	}

	_, err := s.boardRepo.GetByID(boardID)
	if err != nil {
		return err
	}

	return s.boardRepo.Delete(boardID)
}
