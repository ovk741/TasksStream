package postgres

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/ovk741/TasksStream/internal/domain"
)

type BoardRepository struct {
	db *pgxpool.Pool
}

var ErrNotFound = errors.New("not found")
var ErrInvalidInput = errors.New("invalid input")

func NewBoardRepository(db *pgxpool.Pool) *BoardRepository {
	return &BoardRepository{db: db}
}

func (r *BoardRepository) Create(board domain.Board) error {
	_, err := r.db.Exec(
		context.Background(),
		`INSERT INTO boards (id, name, created_at)
		 VALUES ($1, $2, $3)`,
		board.ID,
		board.Name,
		board.CreatedAt,
	)
	return err
}

func (r *BoardRepository) GetAll() ([]domain.Board, error) {
	rows, err := r.db.Query(
		context.Background(),
		`SELECT id, name, created_at FROM boards`,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	boards := []domain.Board{}

	for rows.Next() {
		var b domain.Board
		if err := rows.Scan(&b.ID, &b.Name, &b.CreatedAt); err != nil {
			return nil, err
		}
		boards = append(boards, b)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return boards, nil
}

func (r *BoardRepository) GetByID(boardID string) (domain.Board, error) {
	row := r.db.QueryRow(context.Background(),
		`SELECT id, name, created_at FROM boards WHERE id = $1`, boardID)

	var b domain.Board
	err := row.Scan(&b.ID, &b.Name, &b.CreatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.Board{}, ErrNotFound
		}
		return domain.Board{}, err
	}

	return b, nil
}

func (r *BoardRepository) Update(board domain.Board) (domain.Board, error) {
	row := r.db.QueryRow(
		context.Background(),
		`UPDATE boards
		 SET name = $1
		 WHERE id = $2
		 RETURNING id, name, created_at`,
		board.Name,
		board.ID,
	)

	var updated domain.Board
	if err := row.Scan(&updated.ID, &updated.Name, &updated.CreatedAt); err != nil {
		return domain.Board{}, ErrNotFound
	}

	return updated, nil
}

func (r *BoardRepository) Delete(boardID string) error {
	row := r.db.QueryRow(
		context.Background(),
		`DELETE FROM boards
		 WHERE id = $1
		 RETURNING id`,
		boardID,
	)

	var deletedID string
	if err := row.Scan(&deletedID); err != nil {
		return ErrNotFound
	}

	return nil
}
