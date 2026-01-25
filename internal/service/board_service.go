package service

import (
	"errors"
	"time"

	"github.com/ovk741/TasksStream/internal/domain"
	"github.com/ovk741/TasksStream/internal/storage"
)

type BoardService interface {
	Create(name string) (domain.Board, error)
	GetAll() ([]domain.Board, error)
}

type boardService struct {
	boardRepo  storage.BoardRepository
	generateID func() string
}

func NewBoardService(boardRepo storage.BoardRepository, generateID func() string) BoardService {
	return &boardService{
		boardRepo:  boardRepo,
		generateID: generateID,
	}

}

func (s *boardService) Create(name string) (domain.Board, error) {
	if name == "" {
		return domain.Board{}, errors.New("board name is empty")
	}

	board := domain.Board{
		ID:        s.generateID(),
		Name:      name,
		CreatedAt: time.Now(),
	}

	s.boardRepo.Create(board)

	return board, nil
}
func (s *boardService) GetAll() ([]domain.Board, error) {
	if s == nil {
		return nil, errors.New("boards not found")
	}
	return s.boardRepo.GetAll(), nil
}
