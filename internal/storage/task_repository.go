package storage

import "github.com/ovk741/TasksStream/internal/domain"

type TaskRepository interface {
	Create(domain.Task) (domain.Task, error)
	GetByColumnID(ColumnID string) ([]domain.Task, error)
	GetByID(id string) (domain.Task, error)
	Update(task domain.Task) (domain.Task, error)
	Delete(id string) error
	Move(taskID, columnID string, position int) (domain.Task, error)
}
