package memory

import "github.com/ovk741/TasksStream/internal/domain"

type TaskRepository struct {
	tasks map[string]domain.Task
}

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
