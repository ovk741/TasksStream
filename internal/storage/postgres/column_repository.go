package postgres

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/ovk741/TasksStream/internal/domain"
)

type ColumnRepository struct {
	db *pgxpool.Pool
}

func NewColumnRepository(db *pgxpool.Pool) *ColumnRepository {
	return &ColumnRepository{db: db}
}

func (r *ColumnRepository) Create(column domain.Column) (domain.Column, error) {
	row := r.db.QueryRow(
		context.Background(),
		`INSERT INTO columns (id, title, board_id, position, created_at)
		 VALUES ($1, $2, $3, $4, $5)
		 RETURNING id, title, board_id, position, created_at`,
		column.ID,
		column.Title,
		column.BoardID,
		column.Position,
		column.CreatedAt,
	)

	var created domain.Column
	if err := row.Scan(
		&created.ID,
		&created.Title,
		&created.BoardID,
		&created.Position,
		&created.CreatedAt,
	); err != nil {
		return domain.Column{}, domain.ErrInternal
	}

	return created, nil
}

func (r *ColumnRepository) GetByBoardID(boardID string) ([]domain.Column, error) {

	rows, err := r.db.Query(context.Background(),
		`SELECT id, title, position, board_id, created_at 
		FROM columns 
		WHERE board_id = $1 
		ORDER BY position`,
		boardID,
	)
	if err != nil {
		return nil, domain.ErrInternal
	}
	defer rows.Close()

	columns := make([]domain.Column, 0)

	for rows.Next() {
		var c domain.Column
		if err := rows.Scan(
			&c.ID,
			&c.Title,
			&c.Position,
			&c.BoardID,
			&c.CreatedAt,
		); err != nil {
			return nil, domain.ErrInternal
		}

		columns = append(columns, c)
	}

	if err := rows.Err(); err != nil {
		return nil, domain.ErrInternal
	}

	return columns, nil
}

func (r *ColumnRepository) Update(column domain.Column) (domain.Column, error) {
	row := r.db.QueryRow(
		context.Background(),
		`UPDATE columns
		 SET title = $1
		 WHERE id = $2
		 RETURNING id, title, position, board_id, created_at`,
		column.Title,
		column.ID,
	)

	var updated domain.Column
	err := row.Scan(
		&updated.ID,
		&updated.Title,
		&updated.Position,
		&updated.BoardID,
		&updated.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.Column{}, domain.ErrNotFound
		}
		return domain.Column{}, domain.ErrInternal
	}

	return updated, nil
}

func (r *ColumnRepository) Delete(columnID string) error {
	row := r.db.QueryRow(
		context.Background(),
		`DELETE FROM columns WHERE id = $1 RETURNING id`,
		columnID,
	)

	var deletedID string
	if err := row.Scan(&deletedID); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.ErrNotFound
		}
		return domain.ErrInternal
	}

	return nil
}

func (r *ColumnRepository) GetByID(columnID string) (domain.Column, error) {
	row := r.db.QueryRow(context.Background(),
		`SELECT id, title, position, board_id, created_at 
		FROM columns 
		WHERE id = $1`,
		columnID,
	)

	var c domain.Column
	err := row.Scan(
		&c.ID,
		&c.Title,
		&c.Position,
		&c.BoardID,
		&c.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.Column{}, domain.ErrNotFound
		}
		return domain.Column{}, domain.ErrInternal
	}

	return c, nil
}

func (r *ColumnRepository) Move(columnID string, position int) (domain.Column, error) {
	ctx := context.Background()

	tx, err := r.db.Begin(ctx)
	if err != nil {
		return domain.Column{}, err
	}
	defer tx.Rollback(ctx)

	// get current column position

	var column domain.Column
	err = tx.QueryRow(ctx,
		`SELECT id, board_id, position, title, created_at
		 FROM columns
		 WHERE id = $1`,
		columnID,
	).Scan(
		&column.ID,
		&column.BoardID,
		&column.Position,
		&column.Title,
		&column.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.Column{}, domain.ErrNotFound
		}
		return domain.Column{}, domain.ErrInternal
	}

	// count all columns

	var total int
	err = tx.QueryRow(ctx,
		`SELECT COUNT(*) FROM columns WHERE board_id = $1`,
		column.BoardID,
	).Scan(&total)
	if err != nil {
		return domain.Column{}, domain.ErrInternal
	}

	if position < 0 || position >= total {
		return domain.Column{}, domain.ErrInvalidInput
	}

	// if same position

	oldPosition := column.Position
	if oldPosition == position {
		return column, nil
	}

	// if move right

	if oldPosition < position {
		_, err = tx.Exec(ctx,
			`UPDATE columns
			 SET position = position - 1
			 WHERE board_id = $1
			   AND position > $2
			   AND position <= $3`,
			column.BoardID,
			oldPosition,
			position,
		)
		if err != nil {
			return domain.Column{}, domain.ErrInternal
		}
		// if move left
	} else {
		_, err = tx.Exec(ctx,
			`UPDATE columns
			 SET position = position + 1
			 WHERE board_id = $1
			   AND position >= $2
			   AND position < $3`,
			column.BoardID,
			position,
			oldPosition,
		)
		if err != nil {
			return domain.Column{}, domain.ErrInternal
		}
	}

	//update column
	err = tx.QueryRow(ctx,
		`UPDATE columns
		 SET position = $1
		 WHERE id = $2
		 RETURNING id, board_id, position, title, created_at`,
		position,
		columnID,
	).Scan(
		&column.ID,
		&column.BoardID,
		&column.Position,
		&column.Title,
		&column.CreatedAt,
	)
	if err != nil {
		return domain.Column{}, domain.ErrInternal
	}

	if err := tx.Commit(ctx); err != nil {
		return domain.Column{}, domain.ErrInternal
	}

	return column, nil
}
