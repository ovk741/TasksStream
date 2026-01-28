package storage

import "github.com/ovk741/TasksStream/internal/domain"

type TaskRepository interface {
	Create(domain.Task)
	GetByColumnID(ColumnID string) []domain.Task
	GetByID(id string) (domain.Task, error)
	Update(task domain.Task) error
	Delete(id string) error
}
