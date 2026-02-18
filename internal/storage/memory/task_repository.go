package memory

import (
	"errors"

	"github.com/ovk741/TasksStream/internal/domain"
)

type TaskRepository struct {
	tasks map[string]domain.Task
}

var ErrNotFound = errors.New("not found")

func NewTaskRepository() *TaskRepository {
	return &TaskRepository{
		tasks: make(map[string]domain.Task),
	}
}

func (r *TaskRepository) Create(task domain.Task) {
	r.tasks[task.ID] = task
}

func (r *TaskRepository) GetByColumnID(columnID string) []domain.Task {
	result := []domain.Task{}

	for _, task := range r.tasks {
		if task.ColumnID == columnID {
			result = append(result, task)
		}
	}

	return result
}

func (r *TaskRepository) GetByID(id string) (domain.Task, error) {
	task, ok := r.tasks[id]
	if !ok {
		return domain.Task{}, ErrNotFound
	}
	return task, nil
}

func (r *TaskRepository) Update(task domain.Task) error {
	if _, ok := r.tasks[task.ID]; !ok {
		return ErrNotFound
	}

	r.tasks[task.ID] = task
	return nil
}

func (r *TaskRepository) Delete(id string) error {
	if _, ok := r.tasks[id]; !ok {
		return ErrNotFound
	}

	delete(r.tasks, id)
	return nil
}
