package service

import (
	"errors"
	"time"

	"github.com/ovk741/TasksStream/internal/domain"
	"github.com/ovk741/TasksStream/internal/storage"
)

type BoardService interface {
	Create(userID, name string) (domain.Board, error)
	GetAll(userID string) ([]domain.Board, error)
	Update(userID, boardID, name string) (domain.Board, error)
	Delete(userID, boardID string) error
}

type boardService struct {
	boardRepo       storage.BoardRepository
	columnRepo      storage.ColumnRepository
	taskRepo        storage.TaskRepository
	boardMemberRepo storage.BoardMemberRepository
	generateID      func() string
}

func NewBoardService(
	boardRepo storage.BoardRepository,
	columnRepo storage.ColumnRepository,
	taskRepo storage.TaskRepository,
	boardMemberRepo storage.BoardMemberRepository,
	generateID func() string,
) BoardService {
	return &boardService{
		boardRepo:       boardRepo,
		columnRepo:      columnRepo,
		taskRepo:        taskRepo,
		boardMemberRepo: boardMemberRepo,
		generateID:      generateID,
	}

}

func (s *boardService) Create(userID, name string) (domain.Board, error) {
	if name == "" || userID == "" {
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

	member := domain.BoardMember{
		ID:        s.generateID(),
		BoardID:   board.ID,
		UserID:    userID,
		Role:      domain.BoardRoleOwner,
		CreatedAt: time.Now(),
	}

	if err := s.boardMemberRepo.Add(member); err != nil {
		return domain.Board{}, err
	}

	return board, nil
}
func (s *boardService) GetAll(userID string) ([]domain.Board, error) {

	if userID == "" {
		return nil, domain.ErrInvalidInput
	}

	boards, err := s.boardRepo.GetAll()
	if err != nil {
		return nil, err
	}

	result := make([]domain.Board, 0, len(boards))

	for _, b := range boards {
		isMember, err := s.boardMemberRepo.IsMember(b.ID, userID)
		if err != nil {
			return nil, err
		}
		if isMember {
			result = append(result, b)
		}
	}

	return result, nil
}

func (s *boardService) Update(userID, boardID, name string) (domain.Board, error) {
	if boardID == "" || name == "" {
		return domain.Board{}, domain.ErrInvalidInput
	}

	role, err := s.requireMember(boardID, userID)
	if err != nil {
		return domain.Board{}, err
	}

	if role == domain.BoardRoleViewer {
		return domain.Board{}, domain.ErrForbidden
	}

	board, err := s.boardRepo.GetByID(boardID)
	if err != nil {
		return domain.Board{}, err
	}

	board.Name = name

	updated, err := s.boardRepo.Update(board)
	if err != nil {
		return domain.Board{}, err
	}
	return updated, nil

}

func (s *boardService) Delete(userID, boardID string) error {
	if boardID == "" {
		return domain.ErrInvalidInput
	}

	_, err := s.boardRepo.GetByID(boardID)
	if err != nil {
		return err
	}

	role, err := s.requireMember(boardID, userID)
	if err != nil {
		return err
	}
	if role != domain.BoardRoleOwner && role != domain.BoardRoleEditor {
		return domain.ErrForbidden
	}

	return s.boardRepo.Delete(boardID)
}

func (s *boardService) requireMember(boardID, userID string) (domain.BoardRole, error) {

	role, err := s.boardMemberRepo.GetRole(boardID, userID)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			return "", domain.ErrForbidden
		}
		return "", err
	}

	return role, nil
}
