package storage

import "github.com/ovk741/TasksStream/internal/domain"

type TaskRepository interface {
	Create(domain.Task)
	GetByColumnID(ColumnID string) []domain.Task
}
