package memory

import (
	"github.com/ovk741/TasksStream/internal/domain"
)

type BoardRepository struct {
	boards map[string]domain.Board
}

func NewBoardRepository() *BoardRepository {
	return &BoardRepository{
		boards: make(map[string]domain.Board),
	}
}

func (r *BoardRepository) Create(board domain.Board) {
	r.boards[board.ID] = board
}

func (r *BoardRepository) GetAll() []domain.Board {
	result := make([]domain.Board, 0, len(r.boards))
	for _, board := range r.boards {
		result = append(result, board)
	}
	return result
}

func (r *BoardRepository) GetByID(boardID string) (domain.Board, error) {
	board, ok := r.boards[boardID]
	if !ok {
		return domain.Board{}, ErrNotFound
	}
	return board, nil
}

func (r *BoardRepository) Update(board domain.Board) error {
	if _, ok := r.boards[board.ID]; !ok {
		return ErrNotFound
	}
	r.boards[board.ID] = board
	return nil
}

func (r *BoardRepository) Delete(boardID string) error {
	if _, ok := r.boards[boardID]; !ok {
		return ErrNotFound
	}

	delete(r.boards, boardID)
	return nil
}
