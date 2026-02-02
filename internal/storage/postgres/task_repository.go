package postgres

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/ovk741/TasksStream/internal/domain"
)

type TaskRepository struct {
	db *pgxpool.Pool
}

func NewTaskRepository(db *pgxpool.Pool) *TaskRepository {
	return &TaskRepository{db: db}
}

func (r *TaskRepository) Create(task domain.Task) (domain.Task, error) {
	row := r.db.QueryRow(
		context.Background(),
		`INSERT INTO tasks (id, title, description, column_id, position, created_at)
		 VALUES ($1, $2, $3, $4, $5, $6)
		 RETURNING id, title, description, column_id, position, created_at`,
		task.ID,
		task.Title,
		task.Description,
		task.ColumnID,
		task.Position,
		task.CreatedAt,
	)

	var created domain.Task
	if err := row.Scan(
		&created.ID,
		&created.Title,
		&created.Description,
		&created.ColumnID,
		&created.Position,
		&created.CreatedAt,
	); err != nil {
		return domain.Task{}, domain.ErrInternal
	}

	return created, nil
}

func (r *TaskRepository) GetByColumnID(columnID string) ([]domain.Task, error) {
	rows, err := r.db.Query(
		context.Background(),
		`SELECT id, title, position, description, column_id, created_at 
		FROM tasks 
		WHERE column_id = $1 
		ORDER BY position`,
		columnID,
	)
	if err != nil {
		return nil, domain.ErrInternal
	}
	defer rows.Close()

	tasks := make([]domain.Task, 0)

	for rows.Next() {
		var t domain.Task
		if err := rows.Scan(
			&t.ID,
			&t.Title,
			&t.Position,
			&t.Description,
			&t.ColumnID,
			&t.CreatedAt,
		); err != nil {
			return nil, domain.ErrInternal
		}

		tasks = append(tasks, t)
	}

	if err := rows.Err(); err != nil {
		return nil, domain.ErrInternal
	}

	return tasks, nil
}

func (r *TaskRepository) GetByID(id string) (domain.Task, error) {
	row := r.db.QueryRow(context.Background(),
		`SELECT id, title, description, position, column_id, created_at 
		FROM tasks 
		WHERE id = $1`,
		id,
	)

	var t domain.Task
	err := row.Scan(
		&t.ID,
		&t.Title,
		&t.Description,
		&t.Position,
		&t.ColumnID,
		&t.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.Task{}, domain.ErrNotFound
		}
		return domain.Task{}, domain.ErrInternal
	}

	return t, nil
}

func (r *TaskRepository) Update(task domain.Task) (domain.Task, error) {
	row := r.db.QueryRow(
		context.Background(),
		`UPDATE tasks
		 SET title = $1, description = $2
		 WHERE id = $3
		 RETURNING id, column_id, title, description, position, created_at`,
		task.Title,
		task.Description,
		task.ID,
	)

	var updated domain.Task
	err := row.Scan(
		&updated.ID,
		&updated.ColumnID,
		&updated.Title,
		&updated.Description,
		&updated.Position,
		&updated.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.Task{}, domain.ErrNotFound
		}
		return domain.Task{}, domain.ErrInternal
	}

	return updated, nil
}

func (r *TaskRepository) Delete(id string) error {
	row := r.db.QueryRow(
		context.Background(),
		`DELETE FROM tasks
		 WHERE id = $1
		 RETURNING id`,
		id,
	)

	var deletedID string
	err := row.Scan(&deletedID)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.ErrNotFound
		}
		return domain.ErrInternal
	}

	return nil
}

func (r *TaskRepository) Move(
	taskID string,
	columnID string,
	position int,
) (domain.Task, error) {

	ctx := context.Background()

	tx, err := r.db.Begin(ctx)
	if err != nil {
		return domain.Task{}, domain.ErrInternal
	}
	defer tx.Rollback(ctx)

	var task domain.Task
	err = tx.QueryRow(ctx,
		`SELECT id, title, column_id, position, description, created_at
		 FROM tasks
		 WHERE id = $1`,
		taskID,
	).Scan(
		&task.ID,
		&task.Title,
		&task.ColumnID,
		&task.Position,
		&task.Description,
		&task.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.Task{}, domain.ErrNotFound
		}
		return domain.Task{}, domain.ErrInternal
	}

	_, err = tx.Exec(ctx,
		`UPDATE tasks
		 SET position = position + 1
		 WHERE column_id = $1
		   AND position >= $2`,
		columnID, position,
	)
	if err != nil {
		return domain.Task{}, domain.ErrInternal
	}

	err = tx.QueryRow(ctx,
		`UPDATE tasks
		 SET column_id = $1, position = $2
		 WHERE id = $3
		 RETURNING id, title, column_id, position, description, created_at`,
		columnID, position, taskID,
	).Scan(
		&task.ID,
		&task.Title,
		&task.ColumnID,
		&task.Position,
		&task.Description,
		&task.CreatedAt,
	)
	if err != nil {
		return domain.Task{}, domain.ErrInternal
	}

	if err := tx.Commit(ctx); err != nil {
		return domain.Task{}, domain.ErrInternal
	}

	return task, nil
}
