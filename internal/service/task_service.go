package service

import (
	"errors"
	"time"

	"github.com/ovk741/TasksStream/internal/domain"
	"github.com/ovk741/TasksStream/internal/storage"
)

type TaskService interface {
	Create(userID, title, description, columnID string) (domain.Task, error)
	GetByColumnID(userID, columnID string) ([]domain.Task, error)
	Update(userID, taskID string, title string, description string) (domain.Task, error)
	Move(userID, taskID string, columnID string, position int) (domain.Task, error)
	Delete(userID, taskID string) error
}

type taskService struct {
	taskRepo        storage.TaskRepository
	columnRepo      storage.ColumnRepository
	boardMemberRepo storage.BoardMemberRepository
	generateID      func() string
}

func NewTaskService(
	taskRepo storage.TaskRepository,
	columnRepo storage.ColumnRepository,
	boardMemberRepo storage.BoardMemberRepository,
	generateID func() string,
) TaskService {
	return &taskService{
		taskRepo:        taskRepo,
		columnRepo:      columnRepo,
		boardMemberRepo: boardMemberRepo,
		generateID:      generateID,
	}
}

func (s *taskService) Create(userID, title string, description string, columnID string) (domain.Task, error) {
	if title == "" || columnID == "" {
		return domain.Task{}, domain.ErrInvalidInput
	}

	column, err := s.columnRepo.GetByID(columnID)
	if err != nil {
		return domain.Task{}, domain.ErrNotFound
	}

	if err := s.requireBoardAccess(column.BoardID, userID); err != nil {
		return domain.Task{}, err
	}

	tasks, err := s.taskRepo.GetByColumnID(columnID)
	if err != nil {
		return domain.Task{}, err
	}

	task := domain.Task{
		ID:          s.generateID(),
		Title:       title,
		ColumnID:    columnID,
		Description: description,
		Position:    len(tasks),
		CreatedAt:   time.Now(),
	}

	if _, err := s.taskRepo.Create(task); err != nil {
		return domain.Task{}, err
	}

	return task, nil
}
func (s *taskService) GetByColumnID(userID, columnID string) ([]domain.Task, error) {
	if columnID == "" {
		return nil, domain.ErrInvalidInput
	}

	column, err := s.columnRepo.GetByID(columnID)
	if err != nil {
		return nil, domain.ErrNotFound
	}

	if err := s.requireBoardAccess(column.BoardID, userID); err != nil {
		return nil, err
	}

	task, err := s.taskRepo.GetByColumnID(columnID)
	if err != nil {
		return nil, err
	}
	return task, nil
}

func (s *taskService) Update(userID, taskID string, title string, description string) (domain.Task, error) {
	if taskID == "" || title == "" {
		return domain.Task{}, domain.ErrInvalidInput
	}
	task, err := s.taskRepo.GetByID(taskID)
	if err != nil {
		return domain.Task{}, err
	}

	column, err := s.columnRepo.GetByID(task.ColumnID)
	if err != nil {
		return domain.Task{}, err
	}

	if err := s.requireBoardAccess(column.BoardID, userID); err != nil {
		return domain.Task{}, err
	}

	task.Title = title

	task.Description = description

	updated, err := s.taskRepo.Update(task)
	if err != nil {
		return domain.Task{}, err
	}

	return updated, nil
}

func (s *taskService) Delete(userID, taskID string) error {
	if taskID == "" {
		return domain.ErrInvalidInput
	}

	task, err := s.taskRepo.GetByID(taskID)
	if err != nil {
		return err

	}

	column, err := s.columnRepo.GetByID(task.ColumnID)
	if err != nil {
		return err
	}

	if err := s.requireBoardAccess(column.BoardID, userID); err != nil {
		return err
	}

	return s.taskRepo.Delete(taskID)

}

func (s *taskService) Move(userID, taskID, columnID string, position int) (domain.Task, error) {

	if taskID == "" || columnID == "" || position < 0 {
		return domain.Task{}, domain.ErrInvalidInput
	}

	task, err := s.taskRepo.GetByID(taskID)
	if err != nil {
		return domain.Task{}, err
	}

	// исходная колонка задачи
	sourceColumn, err := s.columnRepo.GetByID(task.ColumnID)
	if err != nil {
		return domain.Task{}, err
	}

	// целевая колонка
	destColumn, err := s.columnRepo.GetByID(columnID)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			return domain.Task{}, domain.ErrNotFound
		}
		return domain.Task{}, err
	}

	// перемещать задачу можно только в пределах одной доски
	if sourceColumn.BoardID != destColumn.BoardID {
		return domain.Task{}, domain.ErrForbidden
	}

	if err := s.requireBoardAccess(sourceColumn.BoardID, userID); err != nil {
		return domain.Task{}, err
	}

	return s.taskRepo.Move(taskID, columnID, position)
}

func (s *taskService) requireBoardAccess(
	boardID, userID string,
) error {

	_, err := s.boardMemberRepo.GetRole(boardID, userID)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			return domain.ErrForbidden
		}
		return err
	}

	return nil
}
