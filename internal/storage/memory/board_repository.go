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
