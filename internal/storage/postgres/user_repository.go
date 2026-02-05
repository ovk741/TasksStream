package postgres

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/ovk741/TasksStream/internal/domain"
)

type UserRepository struct {
	db *pgxpool.Pool
}

func NewUserRepository(db *pgxpool.Pool) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) Create(user domain.User) error {
	_, err := r.db.Exec(
		context.Background(),
		`INSERT INTO users (id, email, password_hash, created_at)
		 VALUES ($1, $2, $3, $4)`,
		user.ID,
		user.Email,
		user.PasswordHash,
		user.CreatedAt,
	)
	if err != nil {
		return domain.ErrInternal
	}

	return nil
}

func (r *UserRepository) GetByEmail(email string) (domain.User, error) {
	row := r.db.QueryRow(
		context.Background(),
		`SELECT id, email, password_hash, created_at
		 FROM users
		 WHERE email = $1`,
		email,
	)

	var u domain.User
	err := row.Scan(
		&u.ID,
		&u.Email,
		&u.PasswordHash,
		&u.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.User{}, domain.ErrNotFound
		}
		return domain.User{}, domain.ErrInternal
	}

	return u, nil
}

func (r *UserRepository) GetByID(id string) (domain.User, error) {
	row := r.db.QueryRow(
		context.Background(),
		`SELECT id, email, password_hash, created_at
		 FROM users
		 WHERE id = $1`,
		id,
	)

	var u domain.User
	err := row.Scan(
		&u.ID,
		&u.Email,
		&u.PasswordHash,
		&u.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.User{}, domain.ErrNotFound
		}
		return domain.User{}, domain.ErrInternal
	}

	return u, nil
}
