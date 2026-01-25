package service

import (
	"errors"
	"time"

	"github.com/ovk741/TasksStream/internal/domain"
	"github.com/ovk741/TasksStream/internal/storage"
)

type TaskService interface {
	Create(title string, description string, columnID string) (domain.Task, error)
	GetByColumnID(columnID string) ([]domain.Task, error)
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
		return domain.Task{}, errors.New("task title is empty")
	}

	if !s.columnRepo.Exists(columnID) {
		return domain.Task{}, errors.New("column not found")
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

	if !s.columnRepo.Exists(columnID) {
		return nil, errors.New("column not found")
	}
	return s.taskRepo.GetByColumnID(columnID), nil
}
