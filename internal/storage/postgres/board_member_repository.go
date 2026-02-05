package postgres

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/ovk741/TasksStream/internal/domain"
)

type BoardMemberRepository struct {
	db *pgxpool.Pool
}

func NewBoardMemberRepository(db *pgxpool.Pool) *BoardMemberRepository {
	return &BoardMemberRepository{db: db}
}

func (r *BoardMemberRepository) Add(member domain.BoardMember) error {
	_, err := r.db.Exec(
		context.Background(),
		`INSERT INTO board_members (id, board_id, user_id, role, created_at)
		 VALUES ($1, $2, $3, $4, $5)`,
		member.ID,
		member.BoardID,
		member.UserID,
		string(member.Role),
		member.CreatedAt,
	)
	if err != nil {
		return domain.ErrInternal
	}

	return nil
}

func (r *BoardMemberRepository) GetRole(boardID, userID string) (domain.BoardRole, error) {
	row := r.db.QueryRow(
		context.Background(),
		`SELECT role
		 FROM board_members
		 WHERE board_id = $1 AND user_id = $2`,
		boardID,
		userID,
	)

	var role string
	if err := row.Scan(&role); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return "", domain.ErrNotFound
		}
		return "", domain.ErrInternal
	}

	return domain.BoardRole(role), nil
}

func (r *BoardMemberRepository) IsMember(boardID, userID string) (bool, error) {
	row := r.db.QueryRow(
		context.Background(),
		`SELECT 1
		 FROM board_members
		 WHERE board_id = $1 AND user_id = $2`,
		boardID,
		userID,
	)

	var dummy int
	if err := row.Scan(&dummy); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return false, nil
		}
		return false, domain.ErrInternal
	}

	return true, nil
}
