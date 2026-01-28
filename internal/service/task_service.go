package service

import (
	"time"

	"github.com/ovk741/TasksStream/internal/domain"
	"github.com/ovk741/TasksStream/internal/storage"
)

type TaskService interface {
	Create(title string, description string, columnID string) (domain.Task, error)
	GetByColumnID(columnID string) ([]domain.Task, error)
	Update(taskID string, title string, description string) (domain.Task, error)
	Delete(taskID string) error
}

type taskService struct {
	taskRepo   storage.TaskRepository
	columnRepo storage.ColumnRepository
	generateID func() string
}

func NewTaskService(taskRepo storage.TaskRepository, columnRepo storage.ColumnRepository, generateID func() string) TaskService {
	return &taskService{
		taskRepo:   taskRepo,
		columnRepo: columnRepo,
		generateID: generateID,
	}
}

func (s *taskService) Create(title string, description string, columnID string) (domain.Task, error) {
	if title == "" {
		return domain.Task{}, ErrInvalidInput
	}

	_, err := s.columnRepo.GetByID(columnID)
	if err != nil {
		return domain.Task{}, ErrNotFound
	}

	task := domain.Task{
		ID:          s.generateID(),
		Title:       title,
		ColumnID:    columnID,
		Description: description,
		CreatedAt:   time.Now(),
	}

	s.taskRepo.Create(task)

	return task, nil
}
func (s *taskService) GetByColumnID(columnID string) ([]domain.Task, error) {

	_, err := s.columnRepo.GetByID(columnID)
	if err != nil {
		return nil, ErrNotFound
	}
	return s.taskRepo.GetByColumnID(columnID), nil
}

func (s *taskService) Update(taskID string, title string, description string) (domain.Task, error) {
	if taskID == "" || title == "" {
		return domain.Task{}, ErrInvalidInput
	}
	task, err := s.taskRepo.GetByID(taskID)
	if err != nil {
		return domain.Task{}, err
	}

	task.Title = title

	task.Description = description

	if err := s.taskRepo.Update(task); err != nil {
		return domain.Task{}, err
	}

	return task, nil
}

func (s *taskService) Delete(taskID string) error {
	if taskID == "" {
		return ErrInvalidInput
	}

	_, err := s.taskRepo.GetByID(taskID)
	if err != nil {
		return err

	}
	return s.taskRepo.Delete(taskID)

}
